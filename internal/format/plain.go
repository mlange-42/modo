package format

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/mlange-42/modo/internal/document"
	"github.com/mlange-42/modo/internal/util"
)

const landingPageContentPlain = `# Landing page

JSON created by mojo doc should be placed next to this file.

Additional documentation files go here, too.
They will be processed for doc-tests and copied to folder 'site'.
`

type Plain struct{}

func (f *Plain) Accepts(files []string) error {
	return nil
}

func (f *Plain) ProcessMarkdown(element any, text string, proc *document.Processor) (string, error) {
	return text, nil
}

func (f *Plain) WriteAuxiliary(p *document.Package, dir string, proc *document.Processor) error {
	return nil
}

func (f *Plain) ToFilePath(p string, kind string) string {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md")
	}
	if len(p) == 0 {
		return p
	}
	return p + ".md"
}

func (f *Plain) ToLinkPath(p string, kind string) string {
	return f.ToFilePath(p, kind)
}

func (f *Plain) Input(in string, sources []document.PackageSource) string {
	return in
}

func (f *Plain) Output(out string) string {
	return out
}

func (f *Plain) GitIgnore(in, out string, sources []document.PackageSource) []string {
	return []string{
		"# files generated by 'mojo doc'",
		fmt.Sprintf("/%s/*.json", in),
		"# files generated by Modo",
		fmt.Sprintf("/%s/", out),
		"# test file generated by Modo",
		"/test/",
	}
}

func (f *Plain) CreateDirs(base, in, out string, sources []document.PackageSource, _ *template.Template) error {
	inDir, outDir := path.Join(base, in), path.Join(base, out)
	testDir := path.Join(base, "test")
	if err := util.MkDirs(inDir); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(inDir, "_index.md"), []byte(landingPageContentPlain), 0644); err != nil {
		return err
	}
	if err := util.MkDirs(outDir); err != nil {
		return err
	}
	return util.MkDirs(testDir)
}

func (f *Plain) Clean(out, tests string) error {
	if err := emptyDir(out); err != nil {
		return err
	}
	return emptyDir(tests)
}
