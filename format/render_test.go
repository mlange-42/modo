package format

import (
	"fmt"
	"testing"

	"github.com/mlange-42/modo/document"
	"github.com/stretchr/testify/assert"
)

func TestRenderPackage(tt *testing.T) {
	pkg := document.Package{
		MemberKind:        document.NewKind("package"),
		MemberName:        document.NewName("Modo"),
		MemberSummary:     document.NewSummary("Mojo documentation generator"),
		MemberDescription: document.NewDescription("Package description"),
		Modules: []*document.Module{
			{
				MemberName:    document.NewName("mod1"),
				MemberSummary: *document.NewSummary("Mod1 summary"),
			},
			{
				MemberName:    document.NewName("mod2"),
				MemberSummary: *document.NewSummary("Mod2 summary"),
			},
		},
		Packages: []*document.Package{},
	}

	form := PlainFormatter{}
	templ, err := loadTemplates(&form)
	assert.Nil(tt, err)

	proc := document.NewProcessor(nil, &form, templ, false, false)

	text, err := renderElement(&pkg, proc)
	assert.Nil(tt, err)

	fmt.Println(text)
}

func TestRenderModule(tt *testing.T) {
	mod := document.Module{
		MemberKind:    document.NewKind("module"),
		MemberName:    document.NewName("modo"),
		Description:   "",
		MemberSummary: *document.NewSummary("a test module"),
		Aliases:       []*document.Alias{},
		Structs: []*document.Struct{
			{
				MemberName:    document.NewName("TestStruct2"),
				MemberSummary: *document.NewSummary("Struct summary..."),
			},
			{
				MemberName:    document.NewName("TestStruct"),
				MemberSummary: *document.NewSummary("Struct summary 2..."),
			},
		},
		Traits:    []*document.Trait{},
		Functions: []*document.Function{},
	}

	form := PlainFormatter{}
	templ, err := loadTemplates(&form)
	assert.Nil(tt, err)

	proc := document.NewProcessor(nil, &form, templ, false, false)

	text, err := renderElement(&mod, proc)
	if err != nil {
		panic(err)
	}
	fmt.Println(text)
}

func TestRenderAll(t *testing.T) {
	yml := `
decl:
  name: modo
  kind: package
  summary: Package modo
  description: |
    Exports:
     - mod1.Struct1
     - mod2
  modules:
    - name: mod1
      kind: module
      structs:
        - name: Struct1
          kind: struct
    - name: mod2
      kind: module
      aliases:
        - name: Alias1
          kind: alias
      structs:
        - name: Struct2
          kind: struct
      traits:
        - name: Trait2
          kind: trait
      functions:
        - name: func2
          kind: function
          overloads:
            - name: func2
              kind: function
`
	docs, err := document.FromYaml([]byte(yml))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	files := map[string]string{}
	err = renderWithWriter(docs, &Config{
		OutputDir:    "out",
		TemplateDirs: []string{},
		RenderFormat: Plain,
		UseExports:   true,
		ShortLinks:   true,
	}, func(file, text string) error {
		files[file] = text
		return nil
	})
	assert.Nil(t, err)

	for f := range files {
		fmt.Println(f)
	}
}
