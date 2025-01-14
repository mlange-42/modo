package document

import (
	"text/template"
)

type Processor struct {
	Template  *template.Template
	Formatter Formatter
}

func NewProcessor(f Formatter, t *template.Template) Processor {
	return Processor{
		Template:  t,
		Formatter: f,
	}
}
