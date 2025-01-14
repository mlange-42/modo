package format

import (
	"fmt"
	"path"
	"text/template"

	"github.com/mlange-42/modo/document"
)

type HugoFormatter struct{}

const hugoFrontMatter = `+++
type = "docs"
title = "%s"
#summary = """%s"""
+++

%s
`

func (f *HugoFormatter) ProcessMarkdown(name, summary, text string) (string, error) {
	return fmt.Sprintf(hugoFrontMatter, name, summary, text), nil
}

func (f *HugoFormatter) WriteAuxiliary(p *document.Package, dir string, t *template.Template) error {
	return nil
}

func (f *HugoFormatter) ToFilePath(p string, kind string) (string, error) {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md"), nil
	}
	return p + ".md", nil
}

func (f *HugoFormatter) ToLinkPath(p string, kind string) (string, error) {
	return fmt.Sprintf("{{< ref \"%s\" >}}", p), nil
}
