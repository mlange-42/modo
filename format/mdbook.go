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

	pkgFile := strings.Builder{}
	if err := t.ExecuteTemplate(&pkgFile, "package_path.md", ""); err != nil {
		return "", err
	}
	s.Summary = fmt.Sprintf("[`%s`](%s)", p.GetName(), pkgFile.String())

	pkgs := strings.Builder{}
	for _, p := range p.Packages {
		if err := f.renderPackage(p, t, []string{}, &pkgs); err != nil {
			return "", err
		}
	}
	s.Packages = pkgs.String()

	mods := strings.Builder{}
	for _, m := range p.Modules {
		if err := f.renderModule(m, t, []string{}, &mods); err != nil {
			return "", err
		}
	}
	s.Modules = mods.String()

	b := strings.Builder{}
	if err := t.ExecuteTemplate(&b, "mdbook_summary.md", &s); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (f *MdBookFormatter) renderPackage(pkg *document.Package, t *template.Template, linkPath []string, out *strings.Builder) error {
	newPath := append([]string{}, linkPath...)
	newPath = append(newPath, pkg.GetFileName())

	pkgFile := strings.Builder{}
	if err := t.ExecuteTemplate(&pkgFile, "package_path.md", path.Join(newPath...)); err != nil {
		return err
	}

	fmt.Fprintf(out, "%-*s- [`%s`](%s))\n", 2*len(linkPath), "", pkg.GetName(), pkgFile.String())
	for _, p := range pkg.Packages {
		if err := f.renderPackage(p, t, newPath, out); err != nil {
			return err
		}
	}
	for _, m := range pkg.Modules {
		if err := f.renderModule(m, t, newPath, out); err != nil {
			return err
		}
	}
	return nil
}

func (f *MdBookFormatter) renderModule(mod *document.Module, t *template.Template, linkPath []string, out *strings.Builder) error {
	newPath := append([]string{}, linkPath...)
	newPath = append(newPath, mod.GetFileName())

	pathStr := path.Join(newPath...)

	pkgFile := strings.Builder{}
	if err := t.ExecuteTemplate(&pkgFile, "module_path.md", pathStr); err != nil {
		return err
	}

	fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*len(linkPath), "", mod.GetName(), pkgFile.String())

	for _, s := range mod.Structs {
		memPath, err := memberPath(t, pathStr, s.GetFileName())
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*len(linkPath)+2, "", s.GetName(), memPath)
	}
	for _, tr := range mod.Traits {
		memPath, err := memberPath(t, pathStr, tr.GetFileName())
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*len(linkPath)+2, "", tr.GetName(), memPath)
	}
	for _, f := range mod.Functions {
		memPath, err := memberPath(t, pathStr, f.GetFileName())
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*len(linkPath)+2, "", f.GetName(), memPath)
	}
	return nil
}

func memberPath(t *template.Template, p string, fn string) (string, error) {
	pathStr := path.Join(p, fn)
	b := strings.Builder{}
	if err := t.ExecuteTemplate(&b, "member_path.md", pathStr); err != nil {
		return "", err
	}
	return b.String(), nil
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
