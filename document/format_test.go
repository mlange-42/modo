package document

import (
	"path"
	"text/template"
)

type TestFormatter struct{}

func (f *TestFormatter) Accepts(files []string) error {
	return nil
}

func (f *TestFormatter) ProcessMarkdown(element any, text string, proc *Processor) (string, error) {
	return text, nil
}

func (f *TestFormatter) WriteAuxiliary(p *Package, dir string, proc *Processor) error {
	return nil
}

func (f *TestFormatter) ToFilePath(p string, kind string) string {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md")
	}
	if len(p) == 0 {
		return p
	}
	return p + ".md"
}

func (f *TestFormatter) ToLinkPath(p string, kind string) string {
	return f.ToFilePath(p, kind)
}

func (f *TestFormatter) Input(in string, sources []PackageSource) string {
	return in
}

func (f *TestFormatter) Output(out string) string {
	return out
}

func (f *TestFormatter) GitIgnore(in, out string, sources []PackageSource) []string {
	return []string{}
}

func (f *TestFormatter) CreateDirs(base, in, out string, sources []PackageSource, templ *template.Template) error {
	return nil
}

func (f *TestFormatter) Clean(out, tests string) error {
	return nil
}
