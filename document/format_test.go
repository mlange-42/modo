package document

import "path"

type TestFormatter struct{}

func (f *TestFormatter) Render(docs *Docs, config *Config) error {
	return Render(docs, config, f)
}

func (f *TestFormatter) ProcessMarkdown(element any, text string, proc *Processor) (string, error) {
	return text, nil
}

func (f *TestFormatter) WriteAuxiliary(p *Package, dir string, proc *Processor) error {
	return nil
}

func (f *TestFormatter) ToFilePath(p string, kind string) (string, error) {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md"), nil
	}
	if len(p) == 0 {
		return p, nil
	}
	return p + ".md", nil
}

func (f *TestFormatter) ToLinkPath(p string, kind string) (string, error) {
	return f.ToFilePath(p, kind)
}
