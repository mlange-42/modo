package document

import (
	"errors"
	"fmt"
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
	templ := template.New("all")
	templ = templ.Funcs(template.FuncMap{
		"toLink": f.ToLinkPath,
	})
	templ, err := templ.ParseFS(assets.Templates, "templates/*.*", "templates/**/*.*")
	if err != nil {
		return nil, err
	}

	for _, dir := range additional {
		if dir == "" {
			continue
		}
		exists, isDir, err := fileExists(dir)
		if err != nil {
			return nil, err
		}
		if !exists || !isDir {
			return nil, fmt.Errorf("template directory '%s' does not exist", dir)
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

func fileExists(file string) (exists, isDir bool, err error) {
	var s os.FileInfo
	if s, err = os.Stat(file); err == nil {
		exists = true
		isDir = s.IsDir()
		return
	} else if !errors.Is(err, os.ErrNotExist) {
		return
	}
	err = nil
	return
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
