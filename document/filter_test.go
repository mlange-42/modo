package document

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createFilterTestDocs() *Docs {
	return &Docs{
		Decl: &Package{
			MemberKind:        NewKind("package"),
			MemberName:        NewName("pkg"),
			MemberDescription: NewDescription(""),
			Packages: []*Package{
				{
					MemberKind:        NewKind("package"),
					MemberName:        NewName("subpkg"),
					MemberDescription: NewDescription(""),
					Modules: []*Module{
						{
							MemberKind: NewKind("module"),
							MemberName: NewName("mod3"),
							Structs: []*Struct{
								{
									MemberKind: NewKind("struct"),
									MemberName: NewName("Struct3"),
								},
							},
						},
					},
				},
			},
			Modules: []*Module{
				{
					MemberKind: NewKind("module"),
					MemberName: NewName("mod1"),
					Structs: []*Struct{
						{
							MemberKind: NewKind("struct"),
							MemberName: NewName("Struct1"),
						},
						{
							MemberKind: NewKind("struct"),
							MemberName: NewName("Struct2"),
						},
					},
					Traits: []*Trait{
						{
							MemberKind: NewKind("trait"),
							MemberName: NewName("Trait"),
						},
					},
					Functions: []*Function{
						{
							MemberKind: NewKind("function"),
							MemberName: NewName("func"),
						},
					},
				},
				{
					MemberKind: NewKind("module"),
					MemberName: NewName("mod2"),
					Structs: []*Struct{
						{
							MemberKind: NewKind("struct"),
							MemberName: NewName("Struct2"),
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
	proc := NewProcessor(docs, nil, nil, true, true)
	proc.filterPackages()
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

	proc := NewProcessor(docs, nil, nil, true, true)
	proc.filterPackages()

	for k, v := range proc.linkExports {
		fmt.Println(k, v)
	}
}
