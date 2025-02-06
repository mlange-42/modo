package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClean(t *testing.T) {
	cwd, err := os.Getwd()
	assert.Nil(t, err)

	cmd, err := buildCommand()
	assert.Nil(t, err)
	cmd.SetArgs([]string{"../../test"})

	err = cmd.Execute()
	assert.Nil(t, err)

	err = os.Chdir(cwd)
	assert.Nil(t, err)

	cmd, err = cleanCommand()
	assert.Nil(t, err)
	cmd.SetArgs([]string{"../../test"})

	err = cmd.Execute()
	assert.Nil(t, err)
}
