package format

import (
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"github.com/mlange-42/modo/document"
)

func Render(docs *document.Docs, dir string, rFormat Format, useExports bool, shortLinks bool) error {
	formatter := GetFormatter(rFormat)
	t, err := loadTemplates(formatter)
	if err != nil {
		return err
	}
	proc := document.NewProcessor(formatter, t, useExports, shortLinks)
	err = proc.ProcessLinks(docs)
	if err != nil {
		return err
	}

	err = renderPackage(docs.Decl, dir, &proc)
	if err != nil {
		return err
	}

	if err := proc.Formatter.WriteAuxiliary(docs.Decl, dir, t); err != nil {
		return err
	}

	return nil
}

func renderElement(data interface {
	document.Named
	document.Kinded
}, proc *document.Processor) (string, error) {
	b := strings.Builder{}
	err := proc.Template.ExecuteTemplate(&b, data.GetKind()+".md", data)
	if err != nil {
		return "", err
	}
	var summary string
	if d, ok := data.(document.Summarized); ok {
		summary = d.GetSummary()
	}
	return proc.Formatter.ProcessMarkdown(data.GetName(), summary, b.String())
}

func renderPackage(p *document.Package, dir string, proc *document.Processor) error {
	pkgPath := path.Join(dir, p.GetFileName())
	if err := mkDirs(pkgPath); err != nil {
		return err
	}

	pkgFile, err := proc.Formatter.ToFilePath(pkgPath, "package")
	if err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		if err := renderPackage(pkg, pkgPath, proc); err != nil {
			return err
		}
	}

	for _, mod := range p.Modules {
		modPath := path.Join(pkgPath, mod.GetFileName())
		if err := renderModule(mod, modPath, proc); err != nil {
			return err
		}
	}

	text, err := renderElement(p, proc)
	if err != nil {
		return err
	}
	if err := os.WriteFile(pkgFile, []byte(text), 0666); err != nil {
		return err
	}

	return nil
}

func renderModule(mod *document.Module, dir string, proc *document.Processor) error {
	if err := mkDirs(dir); err != nil {
		return err
	}

	modFile, err := proc.Formatter.ToFilePath(dir, "module")
	if err != nil {
		return err
	}

	if err := renderList(mod.Structs, dir, proc); err != nil {
		return err
	}
	if err := renderList(mod.Traits, dir, proc); err != nil {
		return err
	}
	if err := renderList(mod.Functions, dir, proc); err != nil {
		return err
	}

	text, err := renderElement(mod, proc)
	if err != nil {
		return err
	}
	if err := os.WriteFile(modFile, []byte(text), 0666); err != nil {
		return err
	}

	return nil
}

func renderList[T interface {
	document.Named
	document.Kinded
}](list []T, dir string, proc *document.Processor) error {
	for _, elem := range list {
		text, err := renderElement(elem, proc)
		if err != nil {
			return err
		}
		memberPath := path.Join(dir, elem.GetFileName())

		memberFile, err := proc.Formatter.ToFilePath(memberPath, elem.GetKind())
		if err != nil {
			return err
		}

		if err := os.WriteFile(memberFile, []byte(text), 0666); err != nil {
			return err
		}
	}
	return nil
}

func loadTemplates(f document.Formatter) (*template.Template, error) {
	allTemplates, err := findTemplates()
	if err != nil {
		return nil, err
	}
	templ, err := template.New("all").Funcs(template.FuncMap{
		"toLink": f.ToLinkPath,
	}).ParseFS(assets.Templates, allTemplates...)
	if err != nil {
		return nil, err
	}
	return templ, nil
}

func findTemplates() ([]string, error) {
	allTemplates := []string{}
	err := fs.WalkDir(assets.Templates, ".",
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				allTemplates = append(allTemplates, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return allTemplates, nil
}

func mkDirs(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
