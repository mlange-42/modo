package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectPaths(t *testing.T) {
	docs := Docs{
		Decl: &Package{
			MemberKind: newKind("package"),
			MemberName: newName("pkg"),
			Packages: []*Package{
				{
					MemberKind: newKind("package"),
					MemberName: newName("subpkg"),
				},
			},
			Modules: []*Module{
				{
					MemberKind: newKind("module"),
					MemberName: newName("mod"),
					Structs: []*Struct{
						{
							MemberKind: newKind("struct"),
							MemberName: newName("Struct"),
							Parameters: []*Parameter{
								{
									MemberKind: newKind("parameter"),
									MemberName: newName("par"),
								},
							},
							Fields: []*Field{
								{
									MemberKind: newKind("field"),
									MemberName: newName("f"),
								},
							},
							Functions: []*Function{
								{
									MemberKind: newKind("function"),
									MemberName: newName("func"),
								},
							},
						},
					},
					Traits: []*Trait{
						{
							MemberKind: newKind("trait"),
							MemberName: newName("Trait"),
							Fields: []*Field{
								{
									MemberKind: newKind("field"),
									MemberName: newName("f"),
								},
							},
							Functions: []*Function{
								{
									MemberKind: newKind("function"),
									MemberName: newName("func"),
								},
							},
						},
					},
					Functions: []*Function{
						{
							MemberKind: newKind("function"),
							MemberName: newName("func"),
							Overloads: []*Function{
								{
									MemberKind: newKind("function"),
									MemberName: newName("func"),
									Parameters: []*Parameter{},
									Args:       []*Arg{},
								},
							},
						},
					},
				},
			},
		},
	}

	proc := NewProcessor(&docs, nil, nil, &Config{})
	proc.ExportDocs = &docs
	proc.collectPaths()
	assert.Equal(t, 11, len(proc.linkTargets))

	tests := []struct {
		mem string
		exp []string
		ok  bool
	}{
		{"pkg", []string{"pkg"}, true},
		{"pkg.subpkg", []string{"pkg", "subpkg"}, true},
		{"pkg.mod.func", []string{"pkg", "mod", "func"}, true},
		{"pkg.mod.Struct.f", []string{"pkg", "mod", "Struct", "#fields"}, true},
		{"pkg.mod.Struct.func", []string{"pkg", "mod", "Struct", "#func"}, true},
	}

	for _, tt := range tests {
		obs, ok := proc.linkTargets[tt.mem]
		assert.Equal(t, tt.ok, ok)
		assert.Equal(t, tt.exp, obs.Elements)
	}
}
