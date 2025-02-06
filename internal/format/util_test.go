package format

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"github.com/mlange-42/modo/internal/document"
	"github.com/stretchr/testify/assert"
)

func listDir(dir string) ([]string, error) {
	paths := []string{}
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			paths = append(paths, strings.ReplaceAll(path, "\\", "/"))
			return nil
		})
	return paths, err
}

func testCreateDirs(f document.Formatter, t *testing.T, outDir string, expected []string) {
	templ := template.New("all")
	templ, err := templ.ParseFS(assets.Config, "**/*")
	assert.Nil(t, err)

	curr, err := os.Getwd()
	assert.Nil(t, err)

	root := t.TempDir()
	err = os.Chdir(root)
	assert.Nil(t, err)

	err = f.CreateDirs("docs", "src", "site", []document.PackageSource{
		{Name: "pkg1", Path: []string{"src", "pkg1"}},
		{Name: "pkg2", Path: []string{"src", "pkg2"}},
	}, templ)
	assert.Nil(t, err)

	files, err := listDir(".")
	assert.Nil(t, err)

	assert.Equal(t, expected, files)

	err = f.Clean(outDir, "test")
	assert.Nil(t, err)

	err = os.Chdir(curr)
	assert.Nil(t, err)
}
