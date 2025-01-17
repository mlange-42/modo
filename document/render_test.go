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
    See [.mod1.Struct1] and [.mod1.Struct1.field field]
    
    Exports:
     - mod1.Struct1
     - mod2
  modules:
    - name: mod1
      kind: module
      structs:
        - name: Struct1
          kind: struct
          fields:
            - name: field
              kind: field
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

	outDir := t.TempDir()
	files := map[string]string{}
	proc := createProcessor(t, docs, true, files)

	err = renderWith(outDir, proc)
	assert.Nil(t, err)
}

func TestRenderStruct(t *testing.T) {
	yml := `
decl:
  name: modo
  kind: package
  modules:
    - name: mod1
      kind: module
      structs:
        - name: Struct1
          kind: struct
          aliases:
            - name: A
              kind: alias
              summary: A summary
              description: A description
          parameters:
            - name: T
              kind: parameter
              description: A description
          fields:
            - name: fld
              kind: field
              summary: A summary
              description: A description
          functions:
            - name: fld
              kind: function
              overloads:
                - name: fld
                  kind: function
                  summary: A summary
                  description: A description
                  parameters:
                    - name: T
                      kind: parameter
                      description: A description
                  args:
                    - name: arg
                      kind: argument
                      description: A description
`
	docs, err := FromYaml([]byte(yml))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	outDir := t.TempDir()
	files := map[string]string{}
	proc := createProcessor(t, docs, false, files)

	err = renderWith(outDir, proc)
	assert.Nil(t, err)
}

func TestRenderTrait(t *testing.T) {
	yml := `
decl:
  name: modo
  kind: package
  modules:
    - name: mod1
      kind: module
      traits:
        - name: Trait1
          kind: trait
          fields:
            - name: fld
              kind: field
              summary: A summary
              description: A description
          functions:
            - name: fld
              kind: function
              overloads:
                - name: fld
                  kind: function
                  summary: A summary
                  description: A description
                  parameters:
                    - name: T
                      kind: parameter
                      description: A description
                  args:
                    - name: arg
                      kind: argument
                      description: A description
`
	docs, err := FromYaml([]byte(yml))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	outDir := t.TempDir()
	files := map[string]string{}
	proc := createProcessor(t, docs, false, files)

	err = renderWith(outDir, proc)
	assert.Nil(t, err)
}

func TestRenderFunction(t *testing.T) {
	yml := `
decl:
  name: modo
  kind: package
  modules:
    - name: mod1
      kind: module
      functions:
        - name: fld
          kind: function
          overloads:
            - name: fld
              kind: function
              summary: A summary
              description: A description
              parameters:
                - name: T
                  kind: parameter
                  description: A description
              args:
                - name: arg
                  kind: argument
                  description: A description
`
	docs, err := FromYaml([]byte(yml))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	outDir := t.TempDir()
	files := map[string]string{}
	proc := createProcessor(t, docs, false, files)

	err = renderWith(outDir, proc)
	assert.Nil(t, err)
}

func createProcessor(t *testing.T, docs *Docs, useExports bool, files map[string]string) *Processor {
	formatter := TestFormatter{}
	templ, err := loadTemplates(&formatter)
	assert.Nil(t, err)
	return NewProcessorWithWriter(docs, &formatter, templ, useExports, true, func(file, text string) error {
		files[file] = text
		return nil
	})
}
