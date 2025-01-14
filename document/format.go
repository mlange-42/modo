package document

import "text/template"

type Formatter interface {
	ToFilePath(path string, kind string) (string, error)
	ToLinkPath(path string, kind string) (string, error)
	WriteAuxiliary(p *Package, dir string, t *template.Template) error
}
