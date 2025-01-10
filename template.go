package modo

import (
	"embed"
	"os"
	"text/template"
)

//go:embed templates/*.md
var templates embed.FS

func Render(data Kinder) error {
	templ, err := template.New("all").ParseFS(
		templates,
		"templates/package.md",
		"templates/module.md",
	)
	if err != nil {
		return err
	}

	err = templ.ExecuteTemplate(os.Stdout, data.GetKind()+".md", data)
	if err != nil {
		return err
	}
	return nil
}
