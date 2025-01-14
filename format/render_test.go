package format

import (
	"fmt"
	"testing"

	"github.com/mlange-42/modo/document"
	"github.com/stretchr/testify/assert"
)

func TestRenderPackage(tt *testing.T) {
	pkg := document.Package{
		MemberKind:    document.NewKind("package"),
		MemberName:    document.NewName("Modo"),
		MemberSummary: document.NewSummary("Mojo documentation generator"),
		Description:   "Package description",
		Modules: []*document.Module{
			{
				MemberName:    document.NewName("mod1"),
				MemberSummary: document.NewSummary("Mod1 summary"),
			},
			{
				MemberName:    document.NewName("mod2"),
				MemberSummary: document.NewSummary("Mod2 summary"),
			},
		},
		Packages: []*document.Package{},
	}

	form := PlainFormatter{}
	templ, err := loadTemplates(&form)
	assert.Nil(tt, err)

	proc := document.NewProcessor(&form, templ)

	text, err := renderElement(&pkg, &proc)
	assert.Nil(tt, err)

	fmt.Println(text)
}

func TestRenderModule(tt *testing.T) {
	mod := document.Module{
		MemberKind:    document.NewKind("module"),
		MemberName:    document.NewName("modo"),
		Description:   "",
		MemberSummary: document.NewSummary("a test module"),
		Aliases:       []*document.Alias{},
		Structs: []*document.Struct{
			{
				MemberName:    document.NewName("TestStruct2"),
				MemberSummary: document.NewSummary("Struct summary..."),
			},
			{
				MemberName:    document.NewName("TestStruct"),
				MemberSummary: document.NewSummary("Struct summary 2..."),
			},
		},
		Traits:    []*document.Trait{},
		Functions: []*document.Function{},
	}

	form := PlainFormatter{}
	templ, err := loadTemplates(&form)
	assert.Nil(tt, err)

	proc := document.NewProcessor(&form, templ)

	text, err := renderElement(&mod, &proc)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}
