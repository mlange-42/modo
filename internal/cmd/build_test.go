package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	cmd, err := buildCommand(nil)
	assert.Nil(t, err)

	cmd.SetArgs([]string{"../../test"})

	err = cmd.Execute()
	assert.Nil(t, err)
}

func TestBuildWatch(t *testing.T) {
	stop := make(chan struct{})

	cmd, err := buildCommand(stop)
	assert.Nil(t, err)

	cmd.SetArgs([]string{"../../test", "--watch"})

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		<-ticker.C
		stop <- struct{}{}
	}()

	err = cmd.Execute()
	assert.Nil(t, err)
}
