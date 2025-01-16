package document

type Formatter interface {
	ToFilePath(path string, kind string) (string, error)
	ToLinkPath(path string, kind string) (string, error)
	ProcessMarkdown(name, summary, text string) (string, error)
	WriteAuxiliary(p *Package, dir string, proc *Processor) error
}
