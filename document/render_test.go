package document

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderPackage(tt *testing.T) {
	pkg := Package{
		MemberKind:        NewKind("package"),
		MemberName:        NewName("Modo"),
		MemberSummary:     NewSummary("Mojo documentation generator"),
		MemberDescription: NewDescription("Package description"),
		Modules: []*Module{
			{
				MemberName:    NewName("mod1"),
				MemberSummary: *NewSummary("Mod1 summary"),
			},
			{
				MemberName:    NewName("mod2"),
				MemberSummary: *NewSummary("Mod2 summary"),
			},
		},
		Packages: []*Package{},
	}

	form := TestFormatter{}
	templ, err := loadTemplates(&form)
	assert.Nil(tt, err)

	proc := NewProcessor(nil, &form, templ, false, false)

	text, err := renderElement(&pkg, proc)
	assert.Nil(tt, err)

	fmt.Println(text)
}

func TestRenderModule(tt *testing.T) {
	mod := Module{
		MemberKind:    NewKind("module"),
		MemberName:    NewName("modo"),
		Description:   "",
		MemberSummary: *NewSummary("a test module"),
		Aliases:       []*Alias{},
		Structs: []*Struct{
			{
				MemberName:    NewName("TestStruct2"),
				MemberSummary: *NewSummary("Struct summary..."),
			},
			{
				MemberName:    NewName("TestStruct"),
				MemberSummary: *NewSummary("Struct summary 2..."),
			},
		},
		Traits:    []*Trait{},
		Functions: []*Function{},
	}

	form := TestFormatter{}
	templ, err := loadTemplates(&form)
	assert.Nil(tt, err)

	proc := NewProcessor(nil, &form, templ, false, false)

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
	docs, err := FromYaml([]byte(yml))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	formatter := TestFormatter{}
	files := map[string]string{}
	templ, err := loadTemplates(&formatter)
	assert.Nil(t, err)
	proc := NewProcessorWithWriter(docs, &formatter, templ, true, true, func(file, text string) error {
		files[file] = text
		return nil
	})

	err = renderWith("out", proc)
	assert.Nil(t, err)

	for f := range files {
		fmt.Println(f)
	}
}
