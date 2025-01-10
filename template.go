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

func Render(data Kinder) (string, error) {
	b := strings.Builder{}
	err := t.ExecuteTemplate(&b, data.GetKind()+".md", data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func RenderPackage(p *Package, dir string) error {
	if err := mkDirs(dir); err != nil {
		return err
	}
	text, err := Render(p)
	if err != nil {
		return err
	}
	pkgPath := path.Join(dir, p.Name)
	if err := os.WriteFile(pkgPath+".md", []byte(text), 0666); err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		if err := RenderPackage(pkg, pkgPath); err != nil {
			return err
		}
	}

	for _, mod := range p.Modules {
		text, err := Render(mod)
		if err != nil {
			return err
		}
		modPath := path.Join(dir, mod.Name)
		if err := os.WriteFile(modPath+".md", []byte(text), 0666); err != nil {
			return err
		}
		if err := mkDirs(modPath); err != nil {
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
	)
}

func mkDirs(path string) error {
	if err := os.Mkdir(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}
