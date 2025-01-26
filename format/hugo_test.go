package format_test

import (
	"testing"

	"github.com/mlange-42/modo/format"
	"github.com/stretchr/testify/assert"
)

func TestHugoToFilePath(t *testing.T) {
	f := format.Hugo{}

	text, err := f.ToFilePath("pkg/mod/Struct", "struct")
	assert.Nil(t, err)
	assert.Equal(t, text, "pkg/mod/Struct.md")

	text, err = f.ToFilePath("pkg/mod", "module")
	assert.Nil(t, err)
	assert.Equal(t, text, "pkg/mod/_index.md")

	text, err = f.ToFilePath("pkg", "package")
	assert.Nil(t, err)
	assert.Equal(t, text, "pkg/_index.md")
}

func TestHugoToLinkPath(t *testing.T) {
	f := format.Hugo{}

	text, err := f.ToLinkPath("pkg/mod/Struct", "struct")
	assert.Nil(t, err)
	assert.Equal(t, text, `{{< ref "pkg/mod/Struct.md" >}}`)

	text, err = f.ToLinkPath("pkg/mod", "module")
	assert.Nil(t, err)
	assert.Equal(t, text, `{{< ref "pkg/mod/_index.md" >}}`)

	text, err = f.ToLinkPath("pkg", "package")
	assert.Nil(t, err)
	assert.Equal(t, text, `{{< ref "pkg/_index.md" >}}`)
}
