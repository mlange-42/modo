package format_test

import (
	"testing"

	"github.com/mlange-42/modo/internal/format"
	"github.com/stretchr/testify/assert"
)

func TestGetFormatter(t *testing.T) {
	f, err := format.GetFormatter("hugo")
	assert.Nil(t, err)
	assert.Empty(t, f, &format.Hugo{})

	f, err = format.GetFormatter("mdbook")
	assert.Nil(t, err)
	assert.Empty(t, f, &format.MdBook{})

	f, err = format.GetFormatter("plain")
	assert.Nil(t, err)
	assert.Empty(t, f, &format.Plain{})

	f, err = format.GetFormatter("foobar")
	assert.NotNil(t, err)
	assert.Nil(t, f)
}
