package format_test

import (
	"testing"

	"github.com/mlange-42/modo/format"
	"github.com/stretchr/testify/assert"
)

func TestMdBookAccepts(t *testing.T) {
	f := format.MdBook{}

	err := f.Accepts([]string{"../test"})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "mdBook formatter can process only a single JSON file, but directory '../test' is given")

	err = f.Accepts([]string{"../main.go", "../go.mod"})
	assert.Equal(t, err.Error(), "mdBook formatter can process only a single JSON file, but 2 is given")
}

func TestMdBookToFilePath(t *testing.T) {
	f := format.MdBook{}

	text := f.ToFilePath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, "pkg/mod/Struct.md")

	text = f.ToFilePath("pkg/mod", "module")
	assert.Equal(t, text, "pkg/mod/_index.md")

	text = f.ToFilePath("pkg", "package")
	assert.Equal(t, text, "pkg/_index.md")
}

func TestMdBookToLinkPath(t *testing.T) {
	f := format.MdBook{}

	text := f.ToLinkPath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, "pkg/mod/Struct.md")

	text = f.ToLinkPath("pkg/mod", "module")
	assert.Equal(t, text, "pkg/mod/_index.md")

	text = f.ToLinkPath("pkg", "package")
	assert.Equal(t, text, "pkg/_index.md")
}
