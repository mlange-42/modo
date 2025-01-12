package format

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/document"
)

type MdBookFormatter struct{}

func (f *MdBookFormatter) WriteAuxiliary(p *document.Package, dir string, t *template.Template) error {
	if err := f.writeSummary(p, dir, t); err != nil {
		return err
	}
	if err := f.writeToml(p, dir, t); err != nil {
		return err
	}
	return nil
}

type summary struct {
	Summary  string
	Packages string
	Modules  string
}

func (f *MdBookFormatter) writeSummary(p *document.Package, dir string, t *template.Template) error {
	summary, err := f.renderSummary(p, t)
	if err != nil {
		return err
	}
	summaryPath := path.Join(dir, p.GetFileName(), "SUMMARY.md")
	if err := os.WriteFile(summaryPath, []byte(summary), 0666); err != nil {
		return err
	}
	return nil
}

func (f *MdBookFormatter) renderSummary(p *document.Package, t *template.Template) (string, error) {
	s := summary{}

	s.Summary = fmt.Sprintf("[`%s`](./_index.md)", p.GetName())

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
		return "", err
	}

	return b.String(), nil
}

func (f *MdBookFormatter) renderPackage(pkg *document.Package, linkPath []string, out *strings.Builder) {
	newPath := append([]string{}, linkPath...)
	newPath = append(newPath, pkg.GetFileName())
	fmt.Fprintf(out, "%-*s- [`%s`](./%s/_index.md))\n", 2*len(linkPath), "", pkg.GetName(), path.Join(newPath...))
	for _, p := range pkg.Packages {
		f.renderPackage(p, newPath, out)
	}
	for _, m := range pkg.Modules {
		f.renderModule(m, newPath, out)
	}
}

func (f *MdBookFormatter) renderModule(mod *document.Module, linkPath []string, out *strings.Builder) {
	newPath := append([]string{}, linkPath...)
	newPath = append(newPath, mod.GetFileName())
	pathStr := path.Join(newPath...)
	fmt.Fprintf(out, "%-*s- [`%s`](./%s/_index.md)\n", 2*len(linkPath), "", mod.GetName(), pathStr)

	for _, s := range mod.Structs {
		fmt.Fprintf(out, "%-*s- [`%s`](./%s/%s.md)\n", 2*len(linkPath)+2, "", s.GetName(), pathStr, s.GetFileName())
	}
	for _, t := range mod.Traits {
		fmt.Fprintf(out, "%-*s- [`%s`](./%s/%s.md)\n", 2*len(linkPath)+2, "", t.GetName(), pathStr, t.GetFileName())
	}
	for _, f := range mod.Functions {
		fmt.Fprintf(out, "%-*s- [`%s`](./%s/%s.md)\n", 2*len(linkPath)+2, "", f.GetName(), pathStr, f.GetFileName())
	}
}

func (f *MdBookFormatter) writeToml(p *document.Package, dir string, t *template.Template) error {
	toml, err := f.renderToml(p, t)
	if err != nil {
		return err
	}
	tomlPath := path.Join(dir, "book.toml")
	if err := os.WriteFile(tomlPath, []byte(toml), 0666); err != nil {
		return err
	}
	return nil
}

func (f *MdBookFormatter) renderToml(p *document.Package, t *template.Template) (string, error) {
	b := strings.Builder{}
	if err := t.ExecuteTemplate(&b, "book.toml", p); err != nil {
		return "", err
	}
	return b.String(), nil
}
