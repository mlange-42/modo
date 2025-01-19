package document

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createFilterTestDocs() *Docs {
	return &Docs{
		Decl: &Package{
			MemberKind:        newKind("package"),
			MemberName:        newName("pkg"),
			MemberDescription: newDescription(""),
			Packages: []*Package{
				{
					MemberKind:        newKind("package"),
					MemberName:        newName("subpkg"),
					MemberDescription: newDescription(""),
					Modules: []*Module{
						{
							MemberKind: newKind("module"),
							MemberName: newName("mod3"),
							Structs: []*Struct{
								{
									MemberKind: newKind("struct"),
									MemberName: newName("Struct3"),
								},
							},
						},
					},
				},
			},
			Modules: []*Module{
				{
					MemberKind: newKind("module"),
					MemberName: newName("mod1"),
					Structs: []*Struct{
						{
							MemberKind: newKind("struct"),
							MemberName: newName("Struct1"),
						},
						{
							MemberKind: newKind("struct"),
							MemberName: newName("Struct2"),
						},
					},
					Traits: []*Trait{
						{
							MemberKind: newKind("trait"),
							MemberName: newName("Trait"),
						},
					},
					Functions: []*Function{
						{
							MemberKind: newKind("function"),
							MemberName: newName("func"),
						},
					},
				},
				{
					MemberKind: newKind("module"),
					MemberName: newName("mod2"),
					Structs: []*Struct{
						{
							MemberKind: newKind("struct"),
							MemberName: newName("Struct2"),
						},
					},
				},
			},
		},
	}
}

func TestFilterPackages(t *testing.T) {
	docs := createFilterTestDocs()

	docs.Decl.Description = `Package pkg

Exports:
 - mod1.Struct1
 - mod1.func
 - mod2
 - subpkg
 - subpkg.mod3
 - subpkg.mod3.Struct3
`

	docs.Decl.Packages[0].Description = `Package subpkg

Exports:
 - mod3.Struct3
`
	proc := NewProcessor(docs, nil, nil, &Config{UseExports: true, ShortLinks: true})
	proc.collectElementPaths()
	err := proc.filterPackages()
	assert.Nil(t, err)

	eDocs := proc.ExportDocs.Decl

	assert.Equal(t, 2, len(eDocs.Structs))
	assert.Equal(t, "Struct1", eDocs.Structs[0].Name)
	assert.Equal(t, "Struct3", eDocs.Structs[1].Name)

	assert.Equal(t, 0, len(eDocs.Traits))

	assert.Equal(t, 1, len(eDocs.Functions))
	assert.Equal(t, "func", eDocs.Functions[0].Name)

	assert.Equal(t, 2, len(eDocs.Modules))
	assert.Equal(t, "mod2", eDocs.Modules[0].Name)
	assert.Equal(t, "mod3", eDocs.Modules[1].Name)

	assert.Equal(t, 1, len(eDocs.Packages))
	assert.Equal(t, "subpkg", eDocs.Packages[0].Name)
	assert.Equal(t, 1, len(eDocs.Packages[0].Structs))
	assert.Equal(t, "Struct3", eDocs.Packages[0].Structs[0].Name)
}

func TestFilterPackagesLinks(t *testing.T) {
	docs := createFilterTestDocs()

	docs.Decl.Description = `Package pkg
Exports:
 - mod1.Struct1
 - mod1.func
 - mod2
 - subpkg
 - subpkg.mod3
 - subpkg.mod3.Struct3
`

	docs.Decl.Packages[0].Description = `Package subpkg
Exports:
 - mod3.Struct3
`

	proc := NewProcessor(docs, nil, nil, &Config{UseExports: true, ShortLinks: true})
	proc.collectElementPaths()
	err := proc.filterPackages()
	assert.Nil(t, err)

	for k, v := range proc.linkExports {
		fmt.Println(k, v)
	}
}
