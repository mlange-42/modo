package format

import (
	"path"
	"text/template"

	"github.com/mlange-42/modo/document"
)

type PlainFormatter struct{}

func (f *PlainFormatter) WriteAuxiliary(p *document.Package, dir string, t *template.Template) error {
	return nil
}

func (f *PlainFormatter) ToFilePath(p string, kind string) (string, error) {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md"), nil
	}
	return p + ".md", nil
}

func (f *PlainFormatter) ToLinkPath(p string, kind string) (string, error) {
	return f.ToFilePath(p, kind)
}
