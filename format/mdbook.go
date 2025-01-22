package format

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"github.com/mlange-42/modo/document"
)

type MdBookFormatter struct{}

func (f *MdBookFormatter) Render(docs *document.Docs, config *document.Config) error {
	return document.Render(docs, config, f)
}

func (f *MdBookFormatter) ProcessMarkdown(element any, text string, proc *document.Processor) (string, error) {
	return text, nil
}

func (f *MdBookFormatter) WriteAuxiliary(p *document.Package, dir string, proc *document.Processor) error {
	if err := f.writeSummary(p, dir, proc); err != nil {
		return err
	}
	if err := f.writeToml(p, dir, proc); err != nil {
		return err
	}
	if err := f.writeCss(dir, proc); err != nil {
		return err
	}
	return nil
}

func (f *MdBookFormatter) ToFilePath(p string, kind string) (string, error) {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md"), nil
	}
	if len(p) == 0 {
		return p, nil
	}
	return p + ".md", nil
}

func (f *MdBookFormatter) ToLinkPath(p string, kind string) (string, error) {
	return f.ToFilePath(p, kind)
}

type summary struct {
	Summary   string
	Packages  string
	Modules   string
	Structs   string
	Traits    string
	Functions string
}

func (f *MdBookFormatter) writeSummary(p *document.Package, dir string, proc *document.Processor) error {
	summary, err := f.renderSummary(p, proc)
	if err != nil {
		return err
	}
	summaryPath := path.Join(dir, p.GetFileName(), "SUMMARY.md")
	if proc.Config.DryRun {
		return nil
	}
	if err := os.WriteFile(summaryPath, []byte(summary), 0644); err != nil {
		return err
	}
	return nil
}

func (f *MdBookFormatter) renderSummary(p *document.Package, proc *document.Processor) (string, error) {
	s := summary{}

	pkgFile, err := f.ToLinkPath("", "package")
	if err != nil {
		return "", err
	}
	s.Summary = fmt.Sprintf("[`%s`](%s)", p.GetName(), pkgFile)

	pkgs := strings.Builder{}
	for _, p := range p.Packages {
		if err := f.renderPackage(p, proc.Template, nil, &pkgs); err != nil {
			return "", err
		}
	}
	s.Packages = pkgs.String()

	mods := strings.Builder{}
	for _, m := range p.Modules {
		if err := f.renderModule(m, nil, &mods); err != nil {
			return "", err
		}
	}
	s.Modules = mods.String()

	elems := strings.Builder{}
	for _, elem := range p.Structs {
		if err := f.renderModuleMember(elem, "", 0, &elems); err != nil {
			return "", err
		}
	}
	s.Structs = elems.String()
	elems = strings.Builder{}
	for _, elem := range p.Traits {
		if err := f.renderModuleMember(elem, "", 0, &elems); err != nil {
			return "", err
		}
	}
	s.Traits = elems.String()
	elems = strings.Builder{}
	for _, elem := range p.Functions {
		if err := f.renderModuleMember(elem, "", 0, &elems); err != nil {
			return "", err
		}
	}
	s.Functions = elems.String()

	b := strings.Builder{}
	if err := proc.Template.ExecuteTemplate(&b, "mdbook_summary.md", &s); err != nil {
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

	pathStr := path.Join(newPath...)
	childDepth := 2*(len(newPath)-1) + 2
	for _, elem := range pkg.Structs {
		if err := f.renderModuleMember(elem, pathStr, childDepth, out); err != nil {
			return err
		}
	}
	for _, elem := range pkg.Traits {
		if err := f.renderModuleMember(elem, pathStr, childDepth, out); err != nil {
			return err
		}
	}
	for _, elem := range pkg.Functions {
		if err := f.renderModuleMember(elem, pathStr, childDepth, out); err != nil {
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
	fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", 2*(len(newPath)-1), "", mod.GetName(), modFile)

	childDepth := 2*(len(newPath)-1) + 2
	for _, elem := range mod.Structs {
		if err := f.renderModuleMember(elem, pathStr, childDepth, out); err != nil {
			return err
		}
	}
	for _, elem := range mod.Traits {
		if err := f.renderModuleMember(elem, pathStr, childDepth, out); err != nil {
			return err
		}
	}
	for _, elem := range mod.Functions {
		if err := f.renderModuleMember(elem, pathStr, childDepth, out); err != nil {
			return err
		}
	}
	return nil
}

func (f *MdBookFormatter) renderModuleMember(mem document.Named, pathStr string, depth int, out io.Writer) error {
	memPath, err := f.ToLinkPath(path.Join(pathStr, mem.GetFileName(), ""), "")
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", depth, "", mem.GetName(), memPath)
	return nil
}

func (f *MdBookFormatter) writeToml(p *document.Package, dir string, proc *document.Processor) error {
	toml, err := f.renderToml(p, proc.Template)
	if err != nil {
		return err
	}
	if proc.Config.DryRun {
		return nil
	}
	tomlPath := path.Join(dir, "book.toml")
	if err := os.WriteFile(tomlPath, []byte(toml), 0644); err != nil {
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

func (f *MdBookFormatter) writeCss(dir string, proc *document.Processor) error {
	cssDir := path.Join(dir, "css")
	if !proc.Config.DryRun {
		if err := os.MkdirAll(cssDir, os.ModePerm); err != nil && !os.IsExist(err) {
			return err
		}
	}
	css, err := fs.ReadFile(assets.CSS, "css/mdbook.css")
	if err != nil {
		return err
	}
	if !proc.Config.DryRun {
		if err := os.WriteFile(path.Join(cssDir, "custom.css"), css, 0644); err != nil {
			return err
		}
	}
	return nil
}
