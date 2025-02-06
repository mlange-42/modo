package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	cmd, err := buildCommand()
	assert.Nil(t, err)

	cmd.SetArgs([]string{"../../test"})

	err = cmd.Execute()
	assert.Nil(t, err)
}
