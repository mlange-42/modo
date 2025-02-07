package cmd

import (
	"os"
	"testing"

	"github.com/mlange-42/modo/internal/document"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	v := NewVersion(1, 2, 3, true)
	assert.Equal(t, "v1.2.3", v.Version())

	v = NewVersion(1, 2, 3, false)
	assert.Equal(t, "v1.2.3-dev", v.Version())
}

func TestGetWatchPaths(t *testing.T) {
	config := document.Config{
		Sources:    []string{"src/mypkg"},
		InputFiles: []string{"docs/src"},
	}

	cwd, err := os.Getwd()
	assert.Nil(t, err)

	err = os.Chdir("../../docs")
	assert.Nil(t, err)

	watch, err := getWatchPaths(&config)
	assert.Nil(t, err)
	assert.Equal(t, []string{"src/mypkg/...", "docs/src/..."}, watch)

	err = os.Chdir(cwd)
	assert.Nil(t, err)
}
