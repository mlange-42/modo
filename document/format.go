package document

type Formatter interface {
	ToFilePath(path string, kind string) (string, error)
	ToLinkPath(path string, kind string) (string, error)
	ProcessMarkdown(element any, text string, proc *Processor) (string, error)
	WriteAuxiliary(p *Package, dir string, proc *Processor) error
}
