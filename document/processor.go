package document

import (
	"os"
	"strings"
	"text/template"
)

type Processor struct {
	Template    *template.Template
	Formatter   Formatter
	UseExports  bool
	ShortLinks  bool
	Docs        *Docs
	ExportDocs  *Docs
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
	err := proc.filterPackages()
	if err != nil {
		return err
	}
	proc.collectPaths()

	if !proc.UseExports {
		for k := range proc.linkTargets {
			proc.linkExports[k] = k
		}
	}
	if err := proc.processLinksPackage(proc.Docs.Decl, []string{}); err != nil {
		return err
	}
	return nil
}

func (proc *Processor) WriteFile(file, text string) error {
	return proc.writer(file, text)
}

func (proc *Processor) addLinkExport(oldPath, newPath []string) {
	proc.linkExports[strings.Join(oldPath, ".")] = strings.Join(newPath, ".")
}

func (proc *Processor) addLinkTarget(elPath, filePath []string, kind string, isSection bool) {
	proc.linkTargets[strings.Join(elPath, ".")] = elemPath{Elements: filePath, Kind: kind, IsSection: isSection}
}
