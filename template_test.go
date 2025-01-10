package modo_test

import (
	"testing"

	"github.com/mlange-24/modo"
)

func TestTemplatePackage(t *testing.T) {
	pkg := modo.Package{
		Kind:        modo.NewKind("package"),
		Name:        "Modo",
		Description: "",
		Summary:     "Mojo documentation generator",
		Modules: []*modo.Module{
			{
				Name:    "mod1",
				Summary: "Mod1 summary",
			},
			{
				Name:    "mod2",
				Summary: "Mod2 summary",
			},
		},
		Packages: []*modo.Package{},
	}

	err := modo.Render(&pkg)
	if err != nil {
		panic(err)
	}
}

func TestTemplateModule(t *testing.T) {
	pkg := modo.Module{
		Kind:        modo.NewKind("module"),
		Name:        "modo",
		Description: "",
		Summary:     "a test module",
		Aliases:     []*modo.Alias{},
		Structs: []*modo.Struct{
			{
				Name:    "TestStruct2",
				Summary: "Struct summary...",
			},
			{
				Name:    "TestStruct",
				Summary: "Struct summary 2...",
			},
		},
		Traits:    []*modo.Trait{},
		Functions: []*modo.Function{},
	}

	err := modo.Render(&pkg)
	if err != nil {
		panic(err)
	}
}
