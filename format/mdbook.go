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

func (f *MdBookFormatter) WriteAuxiliary(p *document.Package, dir string, proc *document.Processor) error {
	if err := f.writeSummary(p, dir, proc); err != nil {
		return err
	}
	if err := f.writeToml(p, dir, proc.Template); err != nil {
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

func (f *MdBookFormatter) writeSummary(p *document.Package, dir string, proc *document.Processor) error {
	var summary string
	var err error
	if proc.UseExports {
		summary, err = f.renderSummaryExport(p, proc)
	} else {
		summary, err = f.renderSummaryNoExport(p, proc)
	}
	if err != nil {
		return err
	}
	summaryPath := path.Join(dir, p.GetFileName(), "SUMMARY.md")
	if err := os.WriteFile(summaryPath, []byte(summary), 0666); err != nil {
		return err
	}
	return nil
}

func (f *MdBookFormatter) renderSummaryExport(p *document.Package, proc *document.Processor) (string, error) {
	s := summary{}

	pkgFile, err := f.ToLinkPath("", "package")
	if err != nil {
		return "", err
	}
	s.Summary = fmt.Sprintf("[`%s`](%s)", p.GetName(), pkgFile)

	toCrawl := map[string]*members{}
	collectExportsPackage(p, toCrawl)

	pkgs := strings.Builder{}
	for _, p := range p.Packages {
		if mem, ok := toCrawl[p.Name]; ok {
			for _, m := range mem.Members {
				if err := f.renderPackageExports(p, proc.Template, []string{}, &m, &pkgs); err != nil {
					return "", err
				}
			}
		}
	}
	s.Packages = pkgs.String()

	mods := strings.Builder{}
	for _, m := range p.Modules {
		if mem, ok := toCrawl[m.Name]; ok {
			for _, mm := range mem.Members {
				if err := f.renderModule(m, []string{}, &mm, &mods); err != nil {
					return "", err
				}
			}
		}
	}
	s.Modules = mods.String()

	b := strings.Builder{}
	if err := proc.Template.ExecuteTemplate(&b, "mdbook_summary.md", &s); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (f *MdBookFormatter) renderSummaryNoExport(p *document.Package, proc *document.Processor) (string, error) {
	s := summary{}

	pkgFile, err := f.ToLinkPath("", "package")
	if err != nil {
		return "", err
	}
	s.Summary = fmt.Sprintf("[`%s`](%s)", p.GetName(), pkgFile)

	pkgs := strings.Builder{}
	for _, p := range p.Packages {
		if err := f.renderPackageNoExports(p, proc.Template, []string{}, &pkgs); err != nil {
			return "", err
		}
	}
	s.Packages = pkgs.String()

	mods := strings.Builder{}
	for _, m := range p.Modules {
		if err := f.renderModule(m, []string{}, nil, &mods); err != nil {
			return "", err
		}
	}
	s.Modules = mods.String()

	b := strings.Builder{}
	if err := proc.Template.ExecuteTemplate(&b, "mdbook_summary.md", &s); err != nil {
		return "", err
	}

	return b.String(), nil
}

func (f *MdBookFormatter) renderPackageExports(pkg *document.Package, t *template.Template, linkPath []string, parentMembers *member, out *strings.Builder) error {
	selfIncluded, toCrawl := collectExportMembers(parentMembers)
	collectExportsPackage(pkg, toCrawl)

	newPath := append([]string{}, linkPath...)
	if selfIncluded {
		newPath = append(newPath, pkg.GetFileName())
	}

	pkgFile, err := f.ToLinkPath(path.Join(newPath...), "package")
	if err != nil {
		return err
	}

	if selfIncluded {
		fmt.Fprintf(out, "%-*s- [`%s`](%s))\n", 2*(len(newPath)-1), "", pkg.GetName(), pkgFile)
	}
	for _, p := range pkg.Packages {
		if mem, ok := toCrawl[p.Name]; ok {
			for _, m := range mem.Members {
				if err := f.renderPackageExports(p, t, newPath, &m, out); err != nil {
					return err
				}
			}
		}
	}
	for _, m := range pkg.Modules {
		if mem, ok := toCrawl[m.Name]; ok {
			for _, mm := range mem.Members {
				if err := f.renderModule(m, newPath, &mm, out); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (f *MdBookFormatter) renderPackageNoExports(pkg *document.Package, t *template.Template, linkPath []string, out *strings.Builder) error {
	newPath := append([]string{}, linkPath...)
	newPath = append(newPath, pkg.GetFileName())

	pkgFile, err := f.ToLinkPath(path.Join(newPath...), "package")
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%-*s- [`%s`](%s))\n", 2*len(linkPath), "", pkg.GetName(), pkgFile)
	for _, p := range pkg.Packages {
		if err := f.renderPackageNoExports(p, t, newPath, out); err != nil {
			return err
		}
	}
	for _, m := range pkg.Modules {
		if err := f.renderModule(m, newPath, nil, out); err != nil {
			return err
		}
	}
	return nil
}

func (f *MdBookFormatter) renderModule(mod *document.Module, linkPath []string, parentMembers *member, out *strings.Builder) error {
	selfIncluded := true
	if parentMembers != nil {
		selfIncluded, _ = collectExportMembers(parentMembers)
	}

	newPath := append([]string{}, linkPath...)
	if selfIncluded {
		newPath = append(newPath, mod.GetFileName())
	}

	pathStr := path.Join(newPath...)

	if selfIncluded {
		modFile, err := f.ToLinkPath(pathStr, "module")
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*(len(newPath)-1), "", mod.GetName(), modFile)
	}

	for _, s := range mod.Structs {
		memPath, err := f.ToLinkPath(path.Join(pathStr, s.GetFileName(), ""), "")
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*(len(newPath)-1)+2, "", s.GetName(), memPath)
	}
	for _, tr := range mod.Traits {
		memPath, err := f.ToLinkPath(path.Join(pathStr, tr.GetFileName(), ""), "")
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*(len(newPath)-1)+2, "", tr.GetName(), memPath)
	}
	for _, ff := range mod.Functions {
		memPath, err := f.ToLinkPath(path.Join(pathStr, ff.GetFileName(), ""), "")
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*(len(newPath)-1)+2, "", ff.GetName(), memPath)
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
