package format

import (
	"text/template"

	"github.com/mlange-42/modo/doc"
)

type Format uint8

const (
	Plain Format = iota
	MdBook
)

var formatters = []Formatter{
	&PlainFormatter{},
	&MdBookFormatter{},
}

func GetFormatter(f Format) Formatter {
	return formatters[f]
}

type Formatter interface {
	WriteAuxiliary(p *doc.Package, dir string, t *template.Template) error
}
