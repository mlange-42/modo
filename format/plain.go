package format

import (
	"text/template"

	"github.com/mlange-42/modo/document"
)

type PlainFormatter struct{}

func (f *PlainFormatter) WriteAuxiliary(p *document.Package, dir string, t *template.Template) error {
	return nil
}
