package format

import (
	"fmt"

	"github.com/mlange-42/modo/document"
)

type Format uint8

const (
	Plain Format = iota
	MdBook
	Hugo
)

func GetFormat(f string) (Format, error) {
	fm, ok := formats[f]
	if !ok {
		return Plain, fmt.Errorf("unknown format '%s'. See flag --format", f)
	}
	return fm, nil
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

func GetFormatter(f Format) document.Formatter {
	return formatters[f]
}
