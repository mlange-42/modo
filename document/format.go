package document

import "text/template"

type Formatter interface {
	Accepts(files []string) error
	ToFilePath(path string, kind string) string
	ToLinkPath(path string, kind string) string
	ProcessMarkdown(element any, text string, proc *Processor) (string, error)
	WriteAuxiliary(p *Package, dir string, proc *Processor) error
	Input(base, in string, sources []PackageSource) string
	Output(base, out string) string
	GitIgnore(in, out string, sources []PackageSource) []string
	CreateDirs(base, in, out string, sources []PackageSource, templ *template.Template) error
}

type PackageSource struct {
	Name string
	Path []string
}
