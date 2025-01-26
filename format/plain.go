package format

import (
	"path"

	"github.com/mlange-42/modo/document"
)

type Plain struct{}

func (f *Plain) Accepts(files []string) error {
	return nil
}

func (f *Plain) ProcessMarkdown(element any, text string, proc *document.Processor) (string, error) {
	return text, nil
}

func (f *Plain) WriteAuxiliary(p *document.Package, dir string, proc *document.Processor) error {
	return nil
}

func (f *Plain) ToFilePath(p string, kind string) string {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md")
	}
	if len(p) == 0 {
		return p
	}
	return p + ".md"
}

func (f *Plain) ToLinkPath(p string, kind string) string {
	return f.ToFilePath(p, kind)
}
