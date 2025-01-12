package modo_test

import (
	"fmt"
	"testing"

	"github.com/mlange-42/modo"
	"github.com/mlange-42/modo/document"
)

func TestTemplatePackage(t *testing.T) {
	pkg := document.Package{
		Kind:        document.NewKind("package"),
		Name:        document.NewName("Modo"),
		Summary:     "Mojo documentation generator",
		Description: "Package description",
		Modules: []*document.Module{
			{
				Name:    document.NewName("mod1"),
				Summary: "Mod1 summary",
			},
			{
				Name:    document.NewName("mod2"),
				Summary: "Mod2 summary",
			},
		},
		Packages: []*document.Package{},
	}

	text, err := modo.Render(&pkg)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}

func TestTemplateModule(t *testing.T) {
	mod := document.Module{
		Kind:        document.NewKind("module"),
		Name:        document.NewName("modo"),
		Description: "",
		Summary:     "a test module",
		Aliases:     []*document.Alias{},
		Structs: []*document.Struct{
			{
				Name:    document.NewName("TestStruct2"),
				Summary: "Struct summary...",
			},
			{
				Name:    document.NewName("TestStruct"),
				Summary: "Struct summary 2...",
			},
		},
		Traits:    []*document.Trait{},
		Functions: []*document.Function{},
	}

	text, err := modo.Render(&mod)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
