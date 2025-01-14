package format

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"github.com/mlange-42/modo/document"
)

type MdBookFormatter struct{}

func (f *MdBookFormatter) ProcessMarkdown(name, summary, text string) (string, error) {
	return text, nil
}

func (f *MdBookFormatter) WriteAuxiliary(p *document.Package, dir string, t *template.Template) error {
	if err := f.writeSummary(p, dir, t); err != nil {
		return err
	}
	if err := f.writeToml(p, dir, t); err != nil {
		return err
	}
	if err := f.writeCss(dir); err != nil {
		return err
	}
	return nil
}

func (f *MdBookFormatter) ToFilePath(p string, kind string) (string, error) {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md"), nil
	}
	return p + ".md", nil
}

func (f *MdBookFormatter) ToLinkPath(p string, kind string) (string, error) {
	return f.ToFilePath(p, kind)
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

	pkgFile, err := f.ToLinkPath("", "package")
	if err != nil {
		return "", err
	}
	s.Summary = fmt.Sprintf("[`%s`](%s)", p.GetName(), pkgFile)

	pkgs := strings.Builder{}
	for _, p := range p.Packages {
		if err := f.renderPackage(p, t, []string{}, &pkgs); err != nil {
			return "", err
		}
	}
	s.Packages = pkgs.String()

	mods := strings.Builder{}
	for _, m := range p.Modules {
		if err := f.renderModule(m, []string{}, &mods); err != nil {
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

	pkgFile, err := f.ToLinkPath(path.Join(newPath...), "package")
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%-*s- [`%s`](%s))\n", 2*len(linkPath), "", pkg.GetName(), pkgFile)
	for _, p := range pkg.Packages {
		if err := f.renderPackage(p, t, newPath, out); err != nil {
			return err
		}
	}
	for _, m := range pkg.Modules {
		if err := f.renderModule(m, newPath, out); err != nil {
			return err
		}
	}
	return nil
}

func (f *MdBookFormatter) renderModule(mod *document.Module, linkPath []string, out *strings.Builder) error {
	newPath := append([]string{}, linkPath...)
	newPath = append(newPath, mod.GetFileName())

	pathStr := path.Join(newPath...)

	modFile, err := f.ToLinkPath(pathStr, "module")
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*len(linkPath), "", mod.GetName(), modFile)

	for _, s := range mod.Structs {
		memPath, err := f.ToLinkPath(path.Join(pathStr, s.GetFileName(), ""), "")
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*len(linkPath)+2, "", s.GetName(), memPath)
	}
	for _, tr := range mod.Traits {
		memPath, err := f.ToLinkPath(path.Join(pathStr, tr.GetFileName(), ""), "")
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*len(linkPath)+2, "", tr.GetName(), memPath)
	}
	for _, ff := range mod.Functions {
		memPath, err := f.ToLinkPath(path.Join(pathStr, ff.GetFileName(), ""), "")
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*len(linkPath)+2, "", ff.GetName(), memPath)
	}
	return nil
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

func (f *MdBookFormatter) writeCss(dir string) error {
	cssDir := path.Join(dir, "css")
	if err := os.MkdirAll(cssDir, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	css, err := fs.ReadFile(assets.CSS, "css/mdbook.css")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(cssDir, "custom.css"), css, 0666); err != nil {
		return err
	}
	return nil
}
