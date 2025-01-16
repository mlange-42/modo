package document

import (
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
}

func NewProcessor(docs *Docs, f Formatter, t *template.Template, useExports bool, shortLinks bool) Processor {
	return Processor{
		Template:   t,
		Formatter:  f,
		UseExports: useExports,
		ShortLinks: shortLinks,
		Docs:       docs,
	}
}

func (proc *Processor) addLinkExport(oldPath, newPath []string) {
	proc.linkExports[strings.Join(oldPath, ".")] = strings.Join(newPath, ".")
}
