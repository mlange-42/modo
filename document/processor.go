package document

import (
	"text/template"
)

type Processor struct {
	Template   *template.Template
	Formatter  Formatter
	ShortLinks bool
}

func NewProcessor(f Formatter, t *template.Template, shortLinks bool) Processor {
	return Processor{
		Template:   t,
		Formatter:  f,
		ShortLinks: shortLinks,
	}
}
