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

var formats = map[string]document.Formatter{
	"":       &PlainFormatter{},
	"plain":  &PlainFormatter{},
	"mdbook": &MdBookFormatter{},
	"hugo":   &HugoFormatter{},
}

func GetFormatter(f string) (document.Formatter, error) {
	fm, ok := formats[f]
	if !ok {
		return nil, fmt.Errorf("unknown format '%s'. See flag --format", f)
	}
	return fm, nil
}
