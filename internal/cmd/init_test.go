package cmd

import (
	"os"
	"path"
	"testing"

	"github.com/mlange-42/modo/internal/util"
	"github.com/stretchr/testify/assert"
)

func setupProject(t *testing.T, dir string) string {
	err := util.MkDirs(path.Join(dir, "src"))
	assert.Nil(t, err)
	err = util.MkDirs(path.Join(dir, "src", "test"))
	assert.Nil(t, err)

	err = os.WriteFile(path.Join(dir, "src", "test", "__init__.mojo"), []byte{}, 0644)
	assert.Nil(t, err)

	cwd, err := os.Getwd()
	assert.Nil(t, err)
	err = os.Chdir(dir)
	assert.Nil(t, err)

	return cwd
}

func TestInitHugo(t *testing.T) {
	dir := t.TempDir()

	cwd := setupProject(t, dir)

	cmd, err := initCommand()
	assert.Nil(t, err)

	cmd.SetArgs([]string{"hugo"})

	err = cmd.Execute()
	assert.Nil(t, err)

	err = os.Chdir(cwd)
	assert.Nil(t, err)
}
