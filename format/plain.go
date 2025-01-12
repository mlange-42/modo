package format

import (
	"text/template"

	"github.com/mlange-42/modo/doc"
)

type PlainFormatter struct{}

func (f *PlainFormatter) WriteAuxiliary(p *doc.Package, dir string, t *template.Template) error {
	return nil
}
