package format

import (
	"strings"
	"testing"

	"github.com/mlange-42/modo/document"
	"github.com/stretchr/testify/assert"
)

func TestHugoToFilePath(t *testing.T) {
	f := Hugo{}

	text := f.ToFilePath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, "pkg/mod/Struct.md")

	text = f.ToFilePath("pkg/mod", "module")
	assert.Equal(t, text, "pkg/mod/_index.md")

	text = f.ToFilePath("pkg", "package")
	assert.Equal(t, text, "pkg/_index.md")
}

func TestHugoToLinkPath(t *testing.T) {
	f := Hugo{}

	text := f.ToLinkPath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, `{{< ref "pkg/mod/Struct.md" >}}`)

	text = f.ToLinkPath("pkg/mod", "module")
	assert.Equal(t, text, `{{< ref "pkg/mod/_index.md" >}}`)

	text = f.ToLinkPath("pkg", "package")
	assert.Equal(t, text, `{{< ref "pkg/_index.md" >}}`)
}

func TestHugoProcessMarkdown(t *testing.T) {
	form := Hugo{}
	templ, err := document.LoadTemplates(&form)
	assert.Nil(t, err)

	proc := document.NewProcessor(nil, &form, templ, &document.Config{})

	text, err := form.ProcessMarkdown(document.Struct{
		MemberName: document.MemberName{Name: "Struct"},
		MemberKind: document.MemberKind{Kind: "struct"},
	}, "test", proc)
	assert.Nil(t, err)

	assert.Equal(t,
		strings.ReplaceAll(text, "\r\n", "\n"),
		`---
type: docs
title: Struct
weight: 100
---

test`)
}

func TestGetGitOrigin(t *testing.T) {
	conf, err := getGitOrigin("docs")

	assert.Nil(t, err)
	assert.Equal(t, conf.Repo, "https://github.com/mlange-42/modo")
	assert.Equal(t, conf.Title, "modo")
	assert.Equal(t, conf.Pages, "https://mlange-42.github.io/modo/")
	assert.Equal(t, conf.Module, "github.com/mlange-42/modo/docs")
}
