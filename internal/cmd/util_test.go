package cmd

import (
	"os"
	"testing"

	"github.com/mlange-42/modo/internal/document"
	"github.com/stretchr/testify/assert"
)

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
