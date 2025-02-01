package format

import (
	"testing"

	"github.com/mlange-42/modo/document"
	"github.com/stretchr/testify/assert"
)

func TestPlainToFilePath(t *testing.T) {
	f := Plain{}

	text := f.ToFilePath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, "pkg/mod/Struct.md")

	text = f.ToFilePath("pkg/mod", "module")
	assert.Equal(t, text, "pkg/mod/_index.md")

	text = f.ToFilePath("pkg", "package")
	assert.Equal(t, text, "pkg/_index.md")
}

func TestPlainToLinkPath(t *testing.T) {
	f := Plain{}

	text := f.ToLinkPath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, "pkg/mod/Struct.md")

	text = f.ToLinkPath("pkg/mod", "module")
	assert.Equal(t, text, "pkg/mod/_index.md")

	text = f.ToLinkPath("pkg", "package")
	assert.Equal(t, text, "pkg/_index.md")
}

func TestPlainInput(t *testing.T) {
	f := Plain{}
	assert.Equal(t, f.Input("src", []document.PackageSource{
		{Name: "pkg", Path: []string{"src", "pkg"}},
	}), "src")
}

func TestPlainOutput(t *testing.T) {
	f := Plain{}
	assert.Equal(t, f.Output("site"), "site")
}
