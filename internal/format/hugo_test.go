package format

import (
	"strings"
	"testing"

	"github.com/mlange-42/modo/internal/document"
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
	templ, err := document.LoadTemplates(&form, "")
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

func TestHugoInput(t *testing.T) {
	f := Hugo{}
	assert.Equal(t, f.Input("src", []document.PackageSource{
		{Name: "pkg", Path: []string{"src", "pkg"}},
	}), "src")
}

func TestHugoOutput(t *testing.T) {
	f := Hugo{}
	assert.Equal(t, f.Output("site"), "site/content")
}

func TestHugoGitIgnore(t *testing.T) {
	f := Hugo{}
	gi := f.GitIgnore("src", "site", []document.PackageSource{})

	assert.Contains(t, gi, "/src/*.json")
	assert.Contains(t, gi, "/site/content/")
	assert.Contains(t, gi, "/site/public/")
	assert.Contains(t, gi, "/site/*.lock")
	assert.Contains(t, gi, "/test/")
}

func TestHugoCreateDirs(t *testing.T) {
	f := Hugo{}
	testCreateDirs(&f, t, "docs/site/content", []string{
		".",
		"docs",
		"docs/site",
		"docs/site/content",
		"docs/site/go.mod",
		"docs/site/go.sum",
		"docs/site/hugo.yaml",
		"docs/src",
		"docs/src/_index.md",
		"docs/test",
	})
}
