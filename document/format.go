package document

import "text/template"

type Formatter interface {
	Accepts(files []string) error
	ToFilePath(path string, kind string) string
	ToLinkPath(path string, kind string) string
	ProcessMarkdown(element any, text string, proc *Processor) (string, error)
	WriteAuxiliary(p *Package, dir string, proc *Processor) error
	Input(in string, sources []PackageSource) string
	Output(out string) string
	GitIgnore(in, out string, sources []PackageSource) []string
	CreateDirs(base, in, out string, sources []PackageSource, templ *template.Template) error
	Clean(out, tests string) error
}

type PackageSource struct {
	Name string
	Path []string
}
