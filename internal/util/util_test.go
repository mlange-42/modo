package util_test

import (
	"testing"

	"github.com/mlange-42/modo/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	tests := []struct {
		File     string
		Exists   bool
		IsDir    bool
		HasError bool
	}{
		{"util.go", true, false, false},
		{"foobar.go", false, false, false},
		{"../format", true, true, false},
	}

	for _, test := range tests {
		exists, isDir, err := util.FileExists(test.File)
		assert.Equal(t, test.Exists, exists)
		assert.Equal(t, test.IsDir, isDir)
		if test.HasError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestGetCwdName(t *testing.T) {
	name, err := util.GetCwdName()
	assert.Nil(t, err)
	assert.Equal(t, "util", name)
}
