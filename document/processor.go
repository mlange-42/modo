package document

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

type Processor struct {
	Template    *template.Template
	Formatter   Formatter
	UseExports  bool
	ShortLinks  bool
	Strict      bool
	Docs        *Docs
	ExportDocs  *Docs
	allPaths    map[string]bool
	linkTargets map[string]elemPath
	linkExports map[string]string
	writer      func(file, text string) error
}

func NewProcessor(docs *Docs, f Formatter, t *template.Template, useExports bool, shortLinks bool) *Processor {
	return NewProcessorWithWriter(docs, f, t, useExports, shortLinks, func(file, text string) error {
		return os.WriteFile(file, []byte(text), 0666)
	})
}

func NewProcessorWithWriter(docs *Docs, f Formatter, t *template.Template, useExports bool, shortLinks bool, writer func(text, file string) error) *Processor {
	return &Processor{
		Template:   t,
		Formatter:  f,
		UseExports: useExports,
		ShortLinks: shortLinks,
		Docs:       docs,
		writer:     writer,
	}
}

// PrepareDocs processes the API docs for subsequent rendering.
func (proc *Processor) PrepareDocs() error {
	// Collect the paths of all (sub)-elements in the original structure.
	proc.collectElementPaths()
	// Restructure according to exports.
	err := proc.filterPackages()
	if err != nil {
		return err
	}
	// Collect all link target paths.
	proc.collectPaths()
	if !proc.UseExports {
		for k := range proc.linkTargets {
			proc.linkExports[k] = k
		}
	}
	// Replaces cross-refs by placeholders.
	if err := proc.processLinksPackage(proc.Docs.Decl, []string{}); err != nil {
		return err
	}
	return nil
}

func (proc *Processor) WriteFile(file, text string) error {
	return proc.writer(file, text)
}

func (proc *Processor) warnOrError(pattern string, args ...any) error {
	if proc.Strict {
		return fmt.Errorf(pattern, args...)
	}
	log.Printf("WARNING: "+pattern+"\n", args...)
	return nil
}

func (proc *Processor) addLinkExport(oldPath, newPath []string) {
	proc.linkExports[strings.Join(oldPath, ".")] = strings.Join(newPath, ".")
}

func (proc *Processor) addLinkTarget(elPath, filePath []string, kind string, isSection bool) {
	proc.linkTargets[strings.Join(elPath, ".")] = elemPath{Elements: filePath, Kind: kind, IsSection: isSection}
}

func (proc *Processor) addElementPath(elPath, filePath []string, kind string, isSection bool) {
	if isSection {
		return
	}
	proc.allPaths[strings.Join(elPath, ".")] = true
	_, _ = filePath, kind
}
