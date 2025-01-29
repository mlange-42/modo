package format

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/document"
)

type MdBook struct{}

func (f *MdBook) Accepts(files []string) error {
	if len(files) > 1 {
		return fmt.Errorf("mdBook formatter can process only a single JSON file, but %d is given", len(files))
	}
	if len(files) == 0 || files[0] == "" {
		return nil
	}
	if s, err := os.Stat(files[0]); err == nil {
		if s.IsDir() {
			return fmt.Errorf("mdBook formatter can process only a single JSON file, but directory '%s' is given", files[0])
		}
	} else {
		return err
	}
	return nil
}

func (f *MdBook) ProcessMarkdown(element any, text string, proc *document.Processor) (string, error) {
	return text, nil
}

func (f *MdBook) WriteAuxiliary(p *document.Package, dir string, proc *document.Processor) error {
	if err := f.writeSummary(p, dir, proc); err != nil {
		return err
	}
	return nil
}

func (f *MdBook) ToFilePath(p string, kind string) string {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md")
	}
	if len(p) == 0 {
		return p
	}
	return p + ".md"
}

func (f *MdBook) ToLinkPath(p string, kind string) string {
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

func (f *MdBook) writeSummary(p *document.Package, dir string, proc *document.Processor) error {
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

func (f *MdBook) renderSummary(p *document.Package, proc *document.Processor) (string, error) {
	s := summary{}

	pkgFile := f.ToLinkPath("", "package")
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

func (f *MdBook) renderPackage(pkg *document.Package, t *template.Template, linkPath []string, out *strings.Builder) error {
	newPath := append([]string{}, linkPath...)
	newPath = append(newPath, pkg.GetFileName())

	pkgFile := f.ToLinkPath(path.Join(newPath...), "package")
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

func (f *MdBook) renderModule(mod *document.Module, linkPath []string, out *strings.Builder) error {
	newPath := append([]string{}, linkPath...)
	newPath = append(newPath, mod.GetFileName())

	pathStr := path.Join(newPath...)

	modFile := f.ToLinkPath(pathStr, "module")
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

func (f *MdBook) renderModuleMember(mem document.Named, pathStr string, depth int, out io.Writer) error {
	memPath := f.ToLinkPath(path.Join(pathStr, mem.GetFileName(), ""), "")
	fmt.Fprintf(out, "%-*s- [`%s`](%s)\n", depth, "", mem.GetName(), memPath)
	return nil
}
