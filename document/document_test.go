package document

import (
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
}

func TestCleanup(t *testing.T) {
	doc := Docs{
		Decl: &Package{
			Modules: []*Module{
				{Name: NewName("__init__")},
				{Name: NewName("modname")},
			},
		},
	}
	cleanup(&doc)

	assert.Equal(t, 1, len(doc.Decl.Modules))
}

func TestCreateSignature(t *testing.T) {
	s := Struct{
		Name: NewName("Struct"),
		Parameters: []*Parameter{
			{Name: NewName("A"), Type: "TypeA", PassingKind: "inferred"},
			{Name: NewName("B"), Type: "TypeB", PassingKind: "pos"},
			{Name: NewName("C"), Type: "TypeC", PassingKind: "pos_or_kw"},
			{Name: NewName("D"), Type: "TypeD", PassingKind: "kw"},
		},
	}

	assert.Equal(t, "struct Struct[A: TypeA, //, B: TypeB, /, C: TypeC, *, D: TypeD]", createSignature(&s))

	s = Struct{
		Name: NewName("Struct"),
		Parameters: []*Parameter{
			{Name: NewName("A"), Type: "TypeA", PassingKind: "inferred"},
		},
	}

	assert.Equal(t, "struct Struct[A: TypeA, //]", createSignature(&s))

	s = Struct{
		Name: NewName("Struct"),
		Parameters: []*Parameter{
			{Name: NewName("B"), Type: "TypeB", PassingKind: "pos"},
		},
	}

	assert.Equal(t, "struct Struct[B: TypeB, /]", createSignature(&s))

	s = Struct{
		Name: NewName("Struct"),
		Parameters: []*Parameter{
			{Name: NewName("D"), Type: "TypeD", PassingKind: "kw"},
		},
	}

	assert.Equal(t, "struct Struct[*, D: TypeD]", createSignature(&s))

}
