package modo

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
)

//go:embed templates/* templates/**/*
var templates embed.FS
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

func Render(data document.Kinded) (string, error) {
	b := strings.Builder{}
	err := t.ExecuteTemplate(&b, data.GetKind()+".md", data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func RenderPackage(p *document.Package, dir string, rFormat format.Format, root bool) error {
	pkgPath := path.Join(dir, p.GetFileName())
	if err := mkDirs(pkgPath); err != nil {
		return err
	}
	p.SetPath(pkgPath)

	pkgFile := strings.Builder{}
	if err := t.ExecuteTemplate(&pkgFile, "package_path.md", pkgPath); err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		if err := RenderPackage(pkg, pkgPath, rFormat, false); err != nil {
			return err
		}
	}

	for _, mod := range p.Modules {
		modPath := path.Join(pkgPath, mod.GetFileName())
		if err := renderModule(mod, modPath); err != nil {
			return err
		}
	}

	text, err := Render(p)
	if err != nil {
		return err
	}
	if err := os.WriteFile(pkgFile.String(), []byte(text), 0666); err != nil {
		return err
	}

	if root {
		if err := format.GetFormatter(rFormat).WriteAuxiliary(p, dir, t); err != nil {
			return err
		}
	}

	return nil
}

func renderModule(mod *document.Module, dir string) error {
	if err := mkDirs(dir); err != nil {
		return err
	}
	mod.SetPath(dir)

	modFile := strings.Builder{}
	if err := t.ExecuteTemplate(&modFile, "module_path.md", dir); err != nil {
		return err
	}

	if err := renderList(mod.Structs, dir); err != nil {
		return err
	}
	if err := renderList(mod.Traits, dir); err != nil {
		return err
	}
	if err := renderList(mod.Functions, dir); err != nil {
		return err
	}

	text, err := Render(mod)
	if err != nil {
		return err
	}
	if err := os.WriteFile(modFile.String(), []byte(text), 0666); err != nil {
		return err
	}

	return nil
}

func renderList[T interface {
	document.Named
	document.Kinded
	document.Pathed
}](list []T, dir string) error {
	for _, elem := range list {
		text, err := Render(elem)
		if err != nil {
			return err
		}
		memberPath := path.Join(dir, elem.GetFileName())

		memberFile := strings.Builder{}
		if err := t.ExecuteTemplate(&memberFile, "member_path.md", memberPath); err != nil {
			return err
		}

		elem.SetPath(memberPath)
		if err := os.WriteFile(memberFile.String(), []byte(text), 0666); err != nil {
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
	templ, err := template.New("all").Funcs(functions).ParseFS(templates, allTemplates...)
	if err != nil {
		return nil, err
	}
	return templ, nil
}

func findTemplates() ([]string, error) {
	allTemplates := []string{}
	err := fs.WalkDir(templates, ".",
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
