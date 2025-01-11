package modo

import (
	"embed"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
)

//go:embed templates/*.md
var templates embed.FS
var t *template.Template

func init() {
	var err error
	t, err = loadTemplates()
	if err != nil {
		log.Fatal(err)
	}
}

func Render(data Kinded) (string, error) {
	b := strings.Builder{}
	err := t.ExecuteTemplate(&b, data.GetKind()+".md", data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func RenderPackage(p *Package, dir string) error {
	pkgPath := path.Join(dir, p.GetName())
	if err := mkDirs(pkgPath); err != nil {
		return err
	}
	p.SetPath(pkgPath)
	pkgFile := path.Join(pkgPath, "_index.md")

	for _, pkg := range p.Packages {
		if err := RenderPackage(pkg, pkgPath); err != nil {
			return err
		}
	}

	for _, mod := range p.Modules {
		modPath := path.Join(pkgPath, mod.GetName())
		if err := renderModule(mod, modPath); err != nil {
			return err
		}
	}

	text, err := Render(p)
	if err != nil {
		return err
	}
	if err := os.WriteFile(pkgFile, []byte(text), 0666); err != nil {
		return err
	}

	return nil
}

func renderModule(mod *Module, dir string) error {
	if err := mkDirs(dir); err != nil {
		return err
	}
	mod.SetPath(dir)
	modFile := path.Join(dir, "_index.md")

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
	if err := os.WriteFile(modFile, []byte(text), 0666); err != nil {
		return err
	}

	return nil
}

func renderList[T interface {
	Named
	Kinded
	Pathed
}](list []T, dir string) error {
	for _, elem := range list {
		text, err := Render(elem)
		if err != nil {
			return err
		}
		strPath := path.Join(dir, elem.GetName())
		elem.SetPath(strPath)
		if err := os.WriteFile(strPath+".md", []byte(text), 0666); err != nil {
			return err
		}
	}
	return nil
}

func loadTemplates() (*template.Template, error) {
	return template.New("all").ParseFS(
		templates,
		"templates/package.md",
		"templates/module.md",
		"templates/struct.md",
		"templates/trait.md",
		"templates/function.md",
	)
}

func mkDirs(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
