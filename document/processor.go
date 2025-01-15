package document

import (
	"text/template"
)

type Processor struct {
	Template   *template.Template
	Formatter  Formatter
	UseExports bool
	ShortLinks bool
}

func NewProcessor(f Formatter, t *template.Template, useExports bool, shortLinks bool) Processor {
	return Processor{
		Template:   t,
		Formatter:  f,
		UseExports: useExports,
		ShortLinks: shortLinks,
	}
}
