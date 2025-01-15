package document

import (
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
