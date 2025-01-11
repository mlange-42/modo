package format

import "github.com/mlange-24/modo/doc"

type Format uint8

const (
	Plain Format = iota
	MdBook
)

var formatters = []Formatter{}

type Formatter interface {
	WriteAuxiliary(p *doc.Package, dir string) error
}
