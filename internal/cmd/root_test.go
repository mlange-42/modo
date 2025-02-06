package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	cmd, err := RootCommand()
	assert.Nil(t, err)

	err = cmd.Execute()
	assert.Nil(t, err)
}

func TestRootVersion(t *testing.T) {
	cmd, err := RootCommand()
	assert.Nil(t, err)

	cmd.SetArgs([]string{"--version"})

	err = cmd.Execute()
	assert.Nil(t, err)
}
