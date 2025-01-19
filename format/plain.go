package format

import (
	"path"

	"github.com/mlange-42/modo/document"
)

type PlainFormatter struct{}

func (f *PlainFormatter) Render(docs *document.Docs, config *document.Config) error {
	return document.Render(docs, config, f)
}

func (f *PlainFormatter) ProcessMarkdown(element any, text string, proc *document.Processor) (string, error) {
	return text, nil
}

func (f *PlainFormatter) WriteAuxiliary(p *document.Package, dir string, proc *document.Processor) error {
	return nil
}

func (f *PlainFormatter) ToFilePath(p string, kind string) (string, error) {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md"), nil
	}
	if len(p) == 0 {
		return p, nil
	}
	return p + ".md", nil
}

func (f *PlainFormatter) ToLinkPath(p string, kind string) (string, error) {
	return f.ToFilePath(p, kind)
}
