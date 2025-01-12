package format

import (
	"text/template"

	"github.com/mlange-42/modo/document"
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

type Config struct {
	Format          Format
	CaseInsensitive bool
}

func GetFormatter(f Format) Formatter {
	return formatters[f]
}

type Formatter interface {
	WriteAuxiliary(p *document.Package, dir string, t *template.Template) error
}
