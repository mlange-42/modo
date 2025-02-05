package format

import (
	"fmt"
	"testing"

	"github.com/mlange-42/modo/internal/document"
	"github.com/stretchr/testify/assert"
)

func TestMdBookAccepts(t *testing.T) {
	f := MdBook{}

	err := f.Accepts([]string{"../../test"})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "mdBook formatter can process only a single JSON file, but directory '../../test' is given")

	err = f.Accepts([]string{"../../main.go", "../../go.mod"})
	assert.Equal(t, err.Error(), "mdBook formatter can process only a single JSON file, but 2 is given")
}

func TestMdBookToFilePath(t *testing.T) {
	f := MdBook{}

	text := f.ToFilePath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, "pkg/mod/Struct.md")

	text = f.ToFilePath("pkg/mod", "module")
	assert.Equal(t, text, "pkg/mod/_index.md")

	text = f.ToFilePath("pkg", "package")
	assert.Equal(t, text, "pkg/_index.md")
}

func TestMdBookToLinkPath(t *testing.T) {
	f := MdBook{}

	text := f.ToLinkPath("pkg/mod/Struct", "struct")
	assert.Equal(t, text, "pkg/mod/Struct.md")

	text = f.ToLinkPath("pkg/mod", "module")
	assert.Equal(t, text, "pkg/mod/_index.md")

	text = f.ToLinkPath("pkg", "package")
	assert.Equal(t, text, "pkg/_index.md")
}

func TestMdBookInput(t *testing.T) {
	f := MdBook{}

	assert.Equal(t, f.Input("src", []document.PackageSource{
		{Name: "pkg", Path: []string{"src", "pkg"}},
	}), "pkg.json")
}

func TestMdBookOutput(t *testing.T) {
	f := MdBook{}

	assert.Equal(t, f.Output("site"), "")
}

func TestMdBookGitIgnore(t *testing.T) {
	f := MdBook{}
	gi := f.GitIgnore("src", "site", []document.PackageSource{{Name: "pkg"}})

	assert.Contains(t, gi, "/pkg.json")
	assert.Contains(t, gi, "/pkg/")
	assert.Contains(t, gi, "/public/")
	assert.Contains(t, gi, "/test/")
}

func TestMdBookCreateDirs(t *testing.T) {
	f := MdBook{}
	testCreateDirs(&f, t, "docs", []string{
		".",
		"docs",
		"docs/book.toml",
		"docs/css",
		"docs/css/custom.css",
		"docs/test",
	})
}

func TestMdBookRenderSummary(t *testing.T) {
	f := MdBook{}

	docs := document.Docs{
		Decl: &document.Package{
			MemberName: document.MemberName{Name: "pkg"},
			MemberKind: document.MemberKind{Kind: "package"},

			Modules: []*document.Module{
				{
					MemberName: document.MemberName{Name: "mod"},
					MemberKind: document.MemberKind{Kind: "module"},
					Structs: []*document.Struct{
						{
							MemberName: document.MemberName{Name: "Struct"},
							MemberKind: document.MemberKind{Kind: "struct"},
						},
					},
				},
			},
			Packages: []*document.Package{
				{
					MemberName: document.MemberName{Name: "subpkg"},
					MemberKind: document.MemberKind{Kind: "package"},
				},
			},
		},
	}

	templ, err := document.LoadTemplates(&f, "")
	assert.Nil(t, err)

	proc := document.NewProcessor(&docs, &f, templ, &document.Config{})

	text, err := f.renderSummary(docs.Decl, proc)
	assert.Nil(t, err)

	fmt.Println(text)

	assert.Contains(t, text, "[`pkg`](_index.md)")
	assert.Contains(t, text, "- [`subpkg`](subpkg/_index.md))")
	assert.Contains(t, text, "- [`mod`](mod/_index.md)")
	assert.Contains(t, text, "  - [`Struct`](mod/Struct.md)")
}
