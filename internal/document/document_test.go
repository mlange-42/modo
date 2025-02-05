package document

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromJson(t *testing.T) {
	data := `{
	"decl": {
    	"kind": "package",
      	"name": "modo",
    	"description": "",
      	"summary": "",
      	"modules": [],
      	"packages": []
	},
    "version": "0.1.0"
}`

	docs, err := FromJson([]byte(data))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	outJson, err := docs.ToJson()
	assert.Nil(t, err)
	fmt.Println(string(outJson))
}

func TestFromYaml(t *testing.T) {
	data := `
decl:
  name: modo
  kind: package
  modules:
    - name: mod
      kind: module
      structs:
        - name: Struct
          kind: struct
version: 0.1.0
`

	docs, err := FromYaml([]byte(data))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	assert.Equal(t, "Struct", docs.Decl.Modules[0].Structs[0].Name)

	outYaml, err := docs.ToYaml()
	assert.Nil(t, err)
	fmt.Println(string(outYaml))
}

func TestCleanup(t *testing.T) {
	doc := Docs{
		Decl: &Package{
			Modules: []*Module{
				{MemberName: newName("__init__")},
				{MemberName: newName("modname")},
			},
		},
	}
	cleanup(&doc)

	assert.Equal(t, 1, len(doc.Decl.Modules))
}

func TestCreateSignature(t *testing.T) {
	s := Struct{
		MemberName: newName("Struct"),
		Parameters: []*Parameter{
			{MemberName: newName("A"), Type: "TypeA", PassingKind: "inferred"},
			{MemberName: newName("B"), Type: "TypeB", PassingKind: "pos"},
			{MemberName: newName("C"), Type: "TypeC", PassingKind: "pos_or_kw"},
			{MemberName: newName("D"), Type: "TypeD", PassingKind: "kw"},
		},
	}

	assert.Equal(t, "struct Struct[A: TypeA, //, B: TypeB, /, C: TypeC, *, D: TypeD]", createSignature(&s))

	s = Struct{
		MemberName: newName("Struct"),
		Parameters: []*Parameter{
			{MemberName: newName("A"), Type: "TypeA", PassingKind: "inferred"},
		},
	}

	assert.Equal(t, "struct Struct[A: TypeA, //]", createSignature(&s))

	s = Struct{
		MemberName: newName("Struct"),
		Parameters: []*Parameter{
			{MemberName: newName("B"), Type: "TypeB", PassingKind: "pos"},
		},
	}

	assert.Equal(t, "struct Struct[B: TypeB, /]", createSignature(&s))

	s = Struct{
		MemberName: newName("Struct"),
		Parameters: []*Parameter{
			{MemberName: newName("D"), Type: "TypeD", PassingKind: "kw"},
		},
	}

	assert.Equal(t, "struct Struct[*, D: TypeD]", createSignature(&s))

}
