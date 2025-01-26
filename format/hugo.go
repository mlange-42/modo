package format

import (
	"fmt"
	"path"
	"strings"

	"github.com/mlange-42/modo/document"
)

type Hugo struct{}

func (f *Hugo) Accepts(files []string) error {
	return nil
}

func (f *Hugo) ProcessMarkdown(element any, text string, proc *document.Processor) (string, error) {
	b := strings.Builder{}
	err := proc.Template.ExecuteTemplate(&b, "hugo_front_matter.md", element)
	if err != nil {
		return "", err
	}
	b.WriteRune('\n')
	b.WriteString(text)
	return b.String(), nil
}

func (f *Hugo) WriteAuxiliary(p *document.Package, dir string, proc *document.Processor) error {
	return nil
}

func (f *Hugo) ToFilePath(p string, kind string) string {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md")
	}
	return p + ".md"
}

func (f *Hugo) ToLinkPath(p string, kind string) string {
	p = f.ToFilePath(p, kind)
	return fmt.Sprintf("{{< ref \"%s\" >}}", p)
}
