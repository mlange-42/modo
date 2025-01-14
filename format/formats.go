package format

import "github.com/mlange-42/modo/document"

type Format uint8

const (
	Plain Format = iota
	MdBook
	Hugo
)

var formatters = []document.Formatter{
	&PlainFormatter{},
	&MdBookFormatter{},
	&HugoFormatter{},
}

type Config struct {
	Format          Format
	CaseInsensitive bool
}

func GetFormatter(f Format) document.Formatter {
	return formatters[f]
}
