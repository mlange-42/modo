package format_test

import (
	"testing"

	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/stretchr/testify/assert"
)

func TestHugoToFilePath(t *testing.T) {
	f := format.Hugo{}

	text := f.ToFilePath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, "pkg/mod/Struct.md")

	text = f.ToFilePath("pkg/mod", "module")
	assert.Equal(t, text, "pkg/mod/_index.md")

	text = f.ToFilePath("pkg", "package")
	assert.Equal(t, text, "pkg/_index.md")
}

func TestHugoToLinkPath(t *testing.T) {
	f := format.Hugo{}

	text := f.ToLinkPath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, `{{< ref "pkg/mod/Struct.md" >}}`)

	text = f.ToLinkPath("pkg/mod", "module")
	assert.Equal(t, text, `{{< ref "pkg/mod/_index.md" >}}`)

	text = f.ToLinkPath("pkg", "package")
	assert.Equal(t, text, `{{< ref "pkg/_index.md" >}}`)
}

func TestHugoProcessMarkdown(t *testing.T) {
	form := format.Hugo{}
	templ, err := document.LoadTemplates(&form)
	assert.Nil(t, err)

	proc := document.NewProcessor(nil, &form, templ, &document.Config{})
	_ = proc
}
