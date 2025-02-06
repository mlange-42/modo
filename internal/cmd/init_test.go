package cmd

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mlange-42/modo/internal/util"
	"github.com/stretchr/testify/assert"
)

func setupProject(t *testing.T, dir string, packages []string) string {
	err := util.MkDirs(path.Join(dir, "src/test"))
	assert.Nil(t, err)
	err = util.MkDirs(path.Join(dir, "foo"))
	assert.Nil(t, err)
	err = util.MkDirs(path.Join(dir, "bar"))
	assert.Nil(t, err)

	for _, p := range packages {
		err := util.MkDirs(path.Join(dir, p))
		assert.Nil(t, err)
		err = os.WriteFile(path.Join(dir, p, "__init__.mojo"), []byte{}, 0644)
		assert.Nil(t, err)
	}

	cwd, err := os.Getwd()
	assert.Nil(t, err)
	err = os.Chdir(dir)
	assert.Nil(t, err)

	return cwd
}

func captureOutput(f func() error) (string, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out), err
}

func TestInitHugo(t *testing.T) {
	allPackages := [][]string{
		{"src"},
		{"package"},
		{"package/src"},
		{"src/package"},
	}

	for _, packages := range allPackages {
		dir := t.TempDir()
		cwd := setupProject(t, dir, packages)

		cmd, err := initCommand(nil)
		assert.Nil(t, err)

		cmd.SetArgs([]string{"hugo"})

		output, err := captureOutput(func() error {
			err = cmd.Execute()
			return err
		})
		assert.Nil(t, err)

		err = os.Chdir(cwd)
		assert.Nil(t, err)

		assert.Equal(t, 2, strings.Count(output, "WARNING"))
	}
}
