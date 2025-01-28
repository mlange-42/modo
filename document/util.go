package document

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

	"github.com/mlange-42/modo/assets"
)

// appends to a slice, but guaranties to return a new one and not alter the original.
func appendNew[T any](sl []T, elems ...T) []T {
	sl2 := make([]T, len(sl), len(sl)+len(elems))
	copy(sl2, sl)
	sl2 = append(sl2, elems...)
	return sl2
}

func warnOrError(strict bool, pattern string, args ...any) error {
	if strict {
		return fmt.Errorf(pattern, args...)
	}
	fmt.Printf("WARNING: "+pattern+"\n", args...)
	return nil
}

func LoadTemplates(f Formatter, additional ...string) (*template.Template, error) {
	allTemplates, err := findTemplatesFS()
	if err != nil {
		return nil, err
	}
	templ := template.New("all")
	templ = templ.Funcs(template.FuncMap{
		"toLink": f.ToLinkPath,
	})
	templ, err = templ.ParseFS(assets.Templates, allTemplates...)
	if err != nil {
		return nil, err
	}

	for _, dir := range additional {
		if dir == "" {
			continue
		}
		moreTemplates, err := findTemplates(dir)
		if err != nil {
			return nil, err
		}
		templ, err = templ.ParseFiles(moreTemplates...)
		if err != nil {
			return nil, err
		}
	}
	return templ, nil
}

func findTemplatesFS() ([]string, error) {
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

func findTemplates(dir string) ([]string, error) {
	allTemplates := []string{}
	err := filepath.WalkDir(dir,
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
