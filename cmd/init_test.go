package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGitOrigin(t *testing.T) {
	conf, err := getGitOrigin("docs")

	assert.Nil(t, err)
	assert.Equal(t, conf.Repo, "https://github.com/mlange-42/modo")
	assert.Equal(t, conf.Title, "modo")
	assert.Equal(t, conf.Pages, "https://mlange-42.github.io/modo/")
	assert.Equal(t, conf.Module, "github.com/mlange-42/modo/docs")
}
