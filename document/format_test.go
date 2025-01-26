package document

import "path"

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
