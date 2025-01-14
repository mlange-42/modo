package format

import "github.com/mlange-42/modo/document"

type Format uint8

const (
	Plain Format = iota
	MdBook
	Hugo
)

func GetFormat(f string) (Format, bool) {
	fm, ok := formats[f]
	return fm, ok
}

var formats = map[string]Format{
	"":       Plain,
	"plain":  Plain,
	"mdbook": MdBook,
	"hugo":   Hugo,
}

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
