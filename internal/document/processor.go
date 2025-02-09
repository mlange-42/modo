package document

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
)

// Processor is the main struct for processing and rendering documentation.
type Processor struct {
	Config             *Config
	Template           *template.Template
	Formatter          Formatter
	Docs               *Docs
	ExportDocs         *Docs
	allPaths           map[string]Named        // Full paths of all original members. Used to check whether all re-exports could be found.
	linkTargets        map[string]elemPath     // Mapping from full (new) member paths to link strings.
	linkExports        map[string]string       // Mapping from original to new member paths.
	linkExportsReverse map[string]*exportError // Used to check for name collisions through re-exports.
	renameExports      map[string]string       // Mapping from short to renamed member paths.
	docTests           []*docTest
	writer             func(file, text string) error
}

type exportError struct {
	NewPath  string
	OldPaths []string
}

type docTest struct {
	Name   string
	Path   []string
	Code   []string
	Global []string
}

// NewProcessor creates a new Processor instance.
func NewProcessor(docs *Docs, f Formatter, t *template.Template, config *Config) *Processor {
	return NewProcessorWithWriter(docs, f, t, config, func(file, text string) error {
		return os.WriteFile(file, []byte(text), 0644)
	})
}

// NewProcessorWithWriter creates a new Processor instance with a custom writer.
func NewProcessorWithWriter(docs *Docs, f Formatter, t *template.Template, config *Config, writer func(file, text string) error) *Processor {
	return &Processor{
		Config:      config,
		Template:    t,
		Formatter:   f,
		Docs:        docs,
		writer:      writer,
		allPaths:    map[string]Named{},
		linkTargets: map[string]elemPath{},
	}
}

// PrepareDocs processes the API docs for subsequent rendering.
func (proc *Processor) PrepareDocs(subdir string) error {
	err := proc.ExtractTests(subdir)
	if err != nil {
		return err
	}
	// Re-structure according to exports.
	err = proc.filterPackages()
	if err != nil {
		return err
	}
	// Collect all link target paths.
	proc.collectPaths()
	if !proc.Config.UseExports {
		for k := range proc.linkTargets {
			proc.linkExports[k] = k
		}
	}
	// Replaces cross-refs by placeholders.
	if err := proc.processLinks(proc.Docs); err != nil {
		return err
	}

	if err := proc.processTranscludes(proc.Docs); err != nil {
		return err
	}

	if proc.Config.UseExports {
		proc.renameAll(proc.ExportDocs.Decl)
	}

	return nil
}

// ExtractTests extracts and writes doc tests.
func (proc *Processor) ExtractTests(subdir string) error {
	// Collect the paths of all (sub)-elements in the original structure.
	proc.collectElementPaths()

	// Extract doc tests.
	err := proc.extractDocTests()
	if err != nil {
		return err
	}
	if proc.Config.TestOutput != "" {
		fmt.Printf("Extracted %d test(s) from package %s.\n", len(proc.docTests), proc.Docs.Decl.Name)
		outPath := path.Join(proc.Config.TestOutput, subdir, proc.Docs.Decl.Name)
		err = proc.writeDocTests(outPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (proc *Processor) writeFile(file, text string) error {
	return proc.writer(file, text)
}

func (proc *Processor) warnOrError(pattern string, args ...any) error {
	return warnOrError(proc.Config.Strict, pattern, args...)
}

func (proc *Processor) addLinkExport(oldPath, newPath []string) {
	pNew := strings.Join(newPath, ".")
	pOld := strings.Join(oldPath, ".")
	if present, ok := proc.linkExportsReverse[pNew]; ok {
		present.OldPaths = append(present.OldPaths, pOld)
	} else {
		proc.linkExportsReverse[pNew] = &exportError{
			NewPath:  pNew,
			OldPaths: []string{pOld},
		}
	}
	proc.linkExports[pOld] = pNew
}

func (proc *Processor) addLinkTarget(elem Named, elPath, filePath []string, kind string, isSection bool) {
	proc.linkTargets[strings.Join(elPath, ".")] = elemPath{Elements: filePath, Kind: kind, IsSection: isSection}
}

func (proc *Processor) addElementPath(elem Named, elPath, filePath []string, kind string, isSection bool) {
	if isSection && kind != "package" && kind != "module" { // actually, we are want to let aliases pass
		return
	}
	proc.allPaths[strings.Join(elPath, ".")] = elem
	_ = filePath
}

func (proc *Processor) mkDirs(path string) error {
	if proc.Config.DryRun {
		return nil
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
