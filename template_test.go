package modo_test

import (
	"fmt"
	"testing"

	"github.com/mlange-24/modo"
	"github.com/mlange-24/modo/doc"
)

func TestTemplatePackage(t *testing.T) {
	pkg := doc.Package{
		Kind:        doc.NewKind("package"),
		Name:        doc.NewName("Modo"),
		Summary:     "Mojo documentation generator",
		Description: "Package description",
		Modules: []*doc.Module{
			{
				Name:    doc.NewName("mod1"),
				Summary: "Mod1 summary",
			},
			{
				Name:    doc.NewName("mod2"),
				Summary: "Mod2 summary",
			},
		},
		Packages: []*doc.Package{},
	}

	text, err := modo.Render(&pkg)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}

func TestTemplateModule(t *testing.T) {
	mod := doc.Module{
		Kind:        doc.NewKind("module"),
		Name:        doc.NewName("modo"),
		Description: "",
		Summary:     "a test module",
		Aliases:     []*doc.Alias{},
		Structs: []*doc.Struct{
			{
				Name:    doc.NewName("TestStruct2"),
				Summary: "Struct summary...",
			},
			{
				Name:    doc.NewName("TestStruct"),
				Summary: "Struct summary 2...",
			},
		},
		Traits:    []*doc.Trait{},
		Functions: []*doc.Function{},
	}

	text, err := modo.Render(&mod)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
