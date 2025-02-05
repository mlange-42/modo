package format

import (
	"fmt"

	"github.com/mlange-42/modo/internal/document"
)

var formats = map[string]document.Formatter{
	"":       &Plain{},
	"plain":  &Plain{},
	"mdbook": &MdBook{},
	"hugo":   &Hugo{},
}

func GetFormatter(f string) (document.Formatter, error) {
	fm, ok := formats[f]
	if !ok {
		return nil, fmt.Errorf("unknown format '%s'. See flag --format", f)
	}
	return fm, nil
}
