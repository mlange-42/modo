package document

type Formatter interface {
	Accepts(files []string) error
	ToFilePath(path string, kind string) string
	ToLinkPath(path string, kind string) string
	ProcessMarkdown(element any, text string, proc *Processor) (string, error)
	WriteAuxiliary(p *Package, dir string, proc *Processor) error
}
