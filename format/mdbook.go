package format

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-24/modo/doc"
)

type MdBookFormatter struct{}

func (f *MdBookFormatter) WriteAuxiliary(p *doc.Package, dir string, t *template.Template) error {
	if err := f.writeSummary(p, dir, t); err != nil {
		return err
	}
	return nil
}

type summary struct {
	Summary  string
	Packages string
	Modules  string
}

func (f *MdBookFormatter) writeSummary(p *doc.Package, dir string, t *template.Template) error {
	summaryPath := path.Join(dir, "SUMMARY.md")

	s := summary{}

	s.Summary = fmt.Sprintf("[%s](./%s/_index.md)", p.GetName(), p.GetName())

	pkgs := strings.Builder{}
	for _, p := range p.Packages {
		f.renderPackage(p, []string{}, &pkgs)
	}
	s.Packages = pkgs.String()

	mods := strings.Builder{}
	for _, m := range p.Modules {
		f.renderModule(m, []string{}, &mods)
	}
	s.Modules = mods.String()

	b := strings.Builder{}
	if err := t.ExecuteTemplate(&b, "mdbook_summary.md", &s); err != nil {
		return err
	}

	if err := os.WriteFile(summaryPath, []byte(b.String()), 0666); err != nil {
		return err
	}
	return nil
}

func (f *MdBookFormatter) renderPackage(pkg *doc.Package, linkPath []string, out *strings.Builder) {
	newPath := append([]string{}, linkPath...)
	newPath = append(newPath, pkg.GetName())
	fmt.Fprintf(out, "%-*s- [%s](./%s/_index.md))\n", 2*len(linkPath), "", pkg.GetName(), path.Join(newPath...))
	for _, p := range pkg.Packages {
		f.renderPackage(p, newPath, out)
	}
	for _, m := range pkg.Modules {
		f.renderModule(m, newPath, out)
	}
}

func (f *MdBookFormatter) renderModule(mod *doc.Module, linkPath []string, out *strings.Builder) {
	newPath := append([]string{}, linkPath...)
	newPath = append(newPath, mod.GetName())
	fmt.Fprintf(out, "%-*s- [%s](./%s/_index.md)\n", 2*len(linkPath), "", mod.GetName(), path.Join(newPath...))
}
