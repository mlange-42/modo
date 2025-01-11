package modo_test

import (
	"fmt"
	"testing"

	"github.com/mlange-24/modo"
)

func TestTemplatePackage(t *testing.T) {
	pkg := modo.Package{
		Kind:        modo.NewKind("package"),
		Name:        modo.NewName("Modo"),
		Summary:     "Mojo documentation generator",
		Description: "Package description",
		Modules: []*modo.Module{
			{
				Name:    modo.NewName("mod1"),
				Summary: "Mod1 summary",
			},
			{
				Name:    modo.NewName("mod2"),
				Summary: "Mod2 summary",
			},
		},
		Packages: []*modo.Package{},
	}

	text, err := modo.Render(&pkg)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}

func TestTemplateModule(t *testing.T) {
	mod := modo.Module{
		Kind:        modo.NewKind("module"),
		Name:        modo.NewName("modo"),
		Description: "",
		Summary:     "a test module",
		Aliases:     []*modo.Alias{},
		Structs: []*modo.Struct{
			{
				Name:    modo.NewName("TestStruct2"),
				Summary: "Struct summary...",
			},
			{
				Name:    modo.NewName("TestStruct"),
				Summary: "Struct summary 2...",
			},
		},
		Traits:    []*modo.Trait{},
		Functions: []*modo.Function{},
	}

	text, err := modo.Render(&mod)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
