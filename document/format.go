package document

type Formatter interface {
	Render(docs *Docs, config *Config) error
	ToFilePath(path string, kind string) (string, error)
	ToLinkPath(path string, kind string) (string, error)
	ProcessMarkdown(element any, text string, proc *Processor) (string, error)
	WriteAuxiliary(p *Package, dir string, proc *Processor) error
}
