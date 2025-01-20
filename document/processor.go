package document

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

type Config struct {
	InputFile       string
	OutputDir       string
	TemplateDirs    []string
	UseExports      bool
	ShortLinks      bool
	CaseInsensitive bool
	Strict          bool
	DryRun          bool
}

type Processor struct {
	Config             *Config
	Template           *template.Template
	Formatter          Formatter
	Docs               *Docs
	ExportDocs         *Docs
	allPaths           map[string]bool
	linkTargets        map[string]elemPath
	linkExports        map[string]string
	linkExportsReverse map[string]*exportError
	writer             func(file, text string) error
}

func NewProcessor(docs *Docs, f Formatter, t *template.Template, config *Config) *Processor {
	return NewProcessorWithWriter(docs, f, t, config, func(file, text string) error {
		return os.WriteFile(file, []byte(text), 0666)
	})
}

func NewProcessorWithWriter(docs *Docs, f Formatter, t *template.Template, config *Config, writer func(file, text string) error) *Processor {
	return &Processor{
		Config:    config,
		Template:  t,
		Formatter: f,
		Docs:      docs,
		writer:    writer,
	}
}

// PrepareDocs processes the API docs for subsequent rendering.
func (proc *Processor) PrepareDocs() error {
	// Collect the paths of all (sub)-elements in the original structure.
	proc.collectElementPaths()

	// Re-structure according to exports.
	err := proc.filterPackages()
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
	return nil
}

func (proc *Processor) WriteFile(file, text string) error {
	return proc.writer(file, text)
}

func (proc *Processor) warnOrError(pattern string, args ...any) error {
	if proc.Config.Strict {
		return fmt.Errorf(pattern, args...)
	}
	log.Printf("WARNING: "+pattern+"\n", args...)
	return nil
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

func (proc *Processor) addLinkTarget(elPath, filePath []string, kind string, isSection bool) {
	proc.linkTargets[strings.Join(elPath, ".")] = elemPath{Elements: filePath, Kind: kind, IsSection: isSection}
}

func (proc *Processor) addElementPath(elPath, filePath []string, kind string, isSection bool) {
	if isSection && kind != "package" && kind != "module" { // actually, we are catching aliases here
		return
	}
	proc.allPaths[strings.Join(elPath, ".")] = true
	_, _ = filePath, kind
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
