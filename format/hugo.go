package format

import (
	"fmt"
	"path"
	"strings"

	"github.com/mlange-42/modo/document"
)

type HugoFormatter struct{}

func (f *HugoFormatter) ProcessMarkdown(element any, text string, proc *document.Processor) (string, error) {
	b := strings.Builder{}
	err := proc.Template.ExecuteTemplate(&b, "hugo_front_matter.md", element)
	if err != nil {
		return "", err
	}
	b.WriteRune('\n')
	b.WriteString(text)
	return b.String(), nil
}

func (f *HugoFormatter) WriteAuxiliary(p *document.Package, dir string, proc *document.Processor) error {
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
