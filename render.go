package modo

import (
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
)

var t *template.Template

var functions = template.FuncMap{
	"pathJoin": path.Join,
}

func init() {
	var err error
	t, err = loadTemplates()
	if err != nil {
		log.Fatal(err)
	}
}

func Render(docs *document.Docs, dir string, rFormat format.Format) error {
	proc := document.NewProcessor(format.GetFormatter(rFormat))
	err := proc.ProcessLinks(docs)
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

func renderElement(data document.Kinded) (string, error) {
	b := strings.Builder{}
	err := t.ExecuteTemplate(&b, data.GetKind()+".md", data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
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

	text, err := renderElement(p)
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

	text, err := renderElement(mod)
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
		text, err := renderElement(elem)
		if err != nil {
			return err
		}
		memberPath := path.Join(dir, elem.GetFileName())

		memberFile, err := proc.Formatter.ToFilePath(memberPath, "")
		if err != nil {
			return err
		}

		if err := os.WriteFile(memberFile, []byte(text), 0666); err != nil {
			return err
		}
	}
	return nil
}

func loadTemplates() (*template.Template, error) {
	allTemplates, err := findTemplates()
	if err != nil {
		return nil, err
	}
	templ, err := template.New("all").Funcs(functions).ParseFS(assets.Templates, allTemplates...)
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
