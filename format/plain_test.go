package format_test

import (
	"testing"

	"github.com/mlange-42/modo/format"
	"github.com/stretchr/testify/assert"
)

func TestPlainToFilePath(t *testing.T) {
	f := format.Plain{}

	text := f.ToFilePath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, "pkg/mod/Struct.md")

	text = f.ToFilePath("pkg/mod", "module")
	assert.Equal(t, text, "pkg/mod/_index.md")

	text = f.ToFilePath("pkg", "package")
	assert.Equal(t, text, "pkg/_index.md")
}

func TestPlainToLinkPath(t *testing.T) {
	f := format.Plain{}

	text := f.ToLinkPath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, "pkg/mod/Struct.md")

	text = f.ToLinkPath("pkg/mod", "module")
	assert.Equal(t, text, "pkg/mod/_index.md")

	text = f.ToLinkPath("pkg", "package")
	assert.Equal(t, text, "pkg/_index.md")
}
