package document

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderPackage(tt *testing.T) {
	pkg := Package{
		MemberKind:        newKind("package"),
		MemberName:        newName("Modo"),
		MemberSummary:     newSummary("Mojo documentation generator"),
		MemberDescription: newDescription("Package description"),
		Modules: []*Module{
			{
				MemberName:    newName("mod1"),
				MemberSummary: *newSummary("Mod1 summary"),
			},
			{
				MemberName:    newName("mod2"),
				MemberSummary: *newSummary("Mod2 summary"),
			},
		},
		Packages: []*Package{},
	}

	form := TestFormatter{}
	templ, err := LoadTemplates(&form, "")
	assert.Nil(tt, err)

	proc := NewProcessor(nil, &form, templ, &Config{})

	text, err := renderElement(&pkg, proc)
	assert.Nil(tt, err)

	fmt.Println(text)
}

func TestRenderModule(tt *testing.T) {
	mod := Module{
		MemberKind:    newKind("module"),
		MemberName:    newName("modo"),
		Description:   "",
		MemberSummary: *newSummary("a test module"),
		Aliases:       []*Alias{},
		Structs: []*Struct{
			{
				MemberName:    newName("TestStruct2"),
				MemberSummary: *newSummary("Struct summary..."),
			},
			{
				MemberName:    newName("TestStruct"),
				MemberSummary: *newSummary("Struct summary 2..."),
			},
		},
		Traits:    []*Trait{},
		Functions: []*Function{},
	}

	form := TestFormatter{}
	templ, err := LoadTemplates(&form, "")
	assert.Nil(tt, err)

	proc := NewProcessor(nil, &form, templ, &Config{})

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
	docs, err := FromYAML([]byte(yml))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	outDir := t.TempDir()
	files := map[string]string{}
	proc := createProcessor(t, docs, true, files)

	err = renderWith(&Config{OutputDir: outDir}, proc, "")
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
	docs, err := FromYAML([]byte(yml))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	outDir := t.TempDir()
	files := map[string]string{}
	proc := createProcessor(t, docs, false, files)

	err = renderWith(&Config{OutputDir: outDir}, proc, "")
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
	docs, err := FromYAML([]byte(yml))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	outDir := t.TempDir()
	files := map[string]string{}
	proc := createProcessor(t, docs, false, files)

	err = renderWith(&Config{OutputDir: outDir}, proc, "")
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
	docs, err := FromYAML([]byte(yml))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	outDir := t.TempDir()
	files := map[string]string{}
	proc := createProcessor(t, docs, false, files)

	err = renderWith(&Config{OutputDir: outDir}, proc, "")
	assert.Nil(t, err)
}

func createProcessor(t *testing.T, docs *Docs, useExports bool, files map[string]string) *Processor {
	formatter := TestFormatter{}
	templ, err := LoadTemplates(&formatter, "")
	assert.Nil(t, err)
	return NewProcessorWithWriter(docs, &formatter, templ, &Config{UseExports: useExports, ShortLinks: true}, func(file, text string) error {
		files[file] = text
		return nil
	})
}

func TestRenderDry(t *testing.T) {
	tmpDir := strings.ReplaceAll(t.TempDir(), "\\", "/")
	config := Config{
		InputFiles:      []string{"../../test/docs/test.json"},
		OutputDir:       tmpDir,
		UseExports:      true,
		ShortLinks:      true,
		CaseInsensitive: true,
		DryRun:          true,
	}
	formatter := TestFormatter{}

	data, err := os.ReadFile(config.InputFiles[0])
	assert.Nil(t, err)
	doc, err := FromJSON(data)
	assert.Nil(t, err)

	err = Render(doc, &config, &formatter, "")
	assert.Nil(t, err)

	tmpFiles, err := filterFiles(tmpDir)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(tmpFiles))
}

func TestRenderFiles(t *testing.T) {
	tmpDir := strings.ReplaceAll(t.TempDir(), "\\", "/")
	refDir := path.Join("..", "..", "test", "ref")
	config := Config{
		InputFiles:      []string{"../../test/docs"},
		SourceURLs:      map[string]string{"test": "https://github.com/mlange-42/modo/blob/main/test/src"},
		OutputDir:       tmpDir,
		UseExports:      true,
		ShortLinks:      true,
		CaseInsensitive: true,
		ReportMissing:   true,
	}
	formatter := TestFormatter{}

	data, err := os.ReadFile(path.Join(config.InputFiles[0], "test.json"))
	assert.Nil(t, err)
	doc, err := FromJSON(data)
	assert.Nil(t, err)

	err = Render(doc, &config, &formatter, "")
	assert.Nil(t, err)

	ExtractTestsMarkdown(&config, &formatter, config.InputFiles[0], true)

	refFiles, err := filterFiles(refDir)
	assert.Nil(t, err)
	tmpFiles, err := filterFiles(tmpDir)
	assert.Nil(t, err)

	assert.Equal(t, len(refFiles), len(tmpFiles))

	for i, ref := range refFiles {
		tmp := tmpFiles[i]
		refShort, tmpShort := strings.TrimPrefix(ref, refDir), strings.TrimPrefix(tmp, tmpDir)
		assert.Equal(t, refShort, tmpShort)

		refContent, err := os.ReadFile(ref)
		assert.Nil(t, err)
		tmpContent, err := os.ReadFile(tmp)
		assert.Nil(t, err)

		refStr, tmpStr := string(refContent), string(tmpContent)
		refStr = strings.ReplaceAll(refStr, "\r\n", "\n")
		tmpStr = strings.ReplaceAll(tmpStr, "\r\n", "\n")

		assert.Equal(t, refStr, tmpStr, "Mismatch in file content for %s", refShort)

		if refStr != tmpStr {
			fmt.Println("=====================")
			fmt.Println(refStr)
			fmt.Println("---------------------")
			fmt.Println(tmpStr)
			fmt.Println("---------------------")
		}
	}
}

func filterFiles(path string) ([]string, error) {
	files := []string{}
	err := filepath.WalkDir(path,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, strings.ReplaceAll(path, "\\", "/"))
			}
			return nil
		})
	return files, err
}
