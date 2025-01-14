package document

import "text/template"

type Processor struct {
	t *template.Template
}

func NewProcessor(t *template.Template) Processor {
	return Processor{
		t: t,
	}
}
