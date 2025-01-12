package modo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mlange-42/modo/document"
	"github.com/stretchr/testify/assert"
)

func TestRenderPackage(tt *testing.T) {
	pkg := document.Package{
		MemberKind:  document.NewKind("package"),
		MemberName:  document.NewName("Modo"),
		Summary:     "Mojo documentation generator",
		Description: "Package description",
		Modules: []*document.Module{
			{
				MemberName: document.NewName("mod1"),
				Summary:    "Mod1 summary",
			},
			{
				MemberName: document.NewName("mod2"),
				Summary:    "Mod2 summary",
			},
		},
		Packages: []*document.Package{},
	}

	text, err := Render(&pkg)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}

func TestRenderModule(tt *testing.T) {
	mod := document.Module{
		MemberKind:  document.NewKind("module"),
		MemberName:  document.NewName("modo"),
		Description: "",
		Summary:     "a test module",
		Aliases:     []*document.Alias{},
		Structs: []*document.Struct{
			{
				MemberName: document.NewName("TestStruct2"),
				Summary:    "Struct summary...",
			},
			{
				MemberName: document.NewName("TestStruct"),
				Summary:    "Struct summary 2...",
			},
		},
		Traits:    []*document.Trait{},
		Functions: []*document.Function{},
	}

	text, err := Render(&mod)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}

func TestPaths(tt *testing.T) {
	p := strings.Builder{}
	err := t.ExecuteTemplate(&p, "package_path.md", "a/b/c")
	assert.Nil(tt, err)
	assert.Equal(tt, "a/b/c/_index.md", p.String())

	p = strings.Builder{}
	err = t.ExecuteTemplate(&p, "module_path.md", "a/b/c")
	assert.Nil(tt, err)
	assert.Equal(tt, "a/b/c/_index.md", p.String())

	p = strings.Builder{}
	err = t.ExecuteTemplate(&p, "member_path.md", "a/b/c")
	assert.Nil(tt, err)
	assert.Equal(tt, "a/b/c.md", p.String())
}
