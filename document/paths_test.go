package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectPaths(t *testing.T) {
	docs := Docs{
		Decl: &Package{
			MemberKind: NewKind("package"),
			MemberName: NewName("pkg"),
			Packages: []*Package{
				{
					MemberKind: NewKind("package"),
					MemberName: NewName("subpkg"),
				},
			},
			Modules: []*Module{
				{
					MemberKind: NewKind("module"),
					MemberName: NewName("mod"),
					Structs: []*Struct{
						{
							MemberKind: NewKind("struct"),
							MemberName: NewName("Struct"),
							Parameters: []*Parameter{
								{
									MemberKind: NewKind("parameter"),
									MemberName: NewName("par"),
								},
							},
							Fields: []*Field{
								{
									MemberKind: NewKind("field"),
									MemberName: NewName("f"),
								},
							},
							Functions: []*Function{
								{
									MemberKind: NewKind("function"),
									MemberName: NewName("func"),
								},
							},
						},
					},
					Traits: []*Trait{
						{
							MemberKind: NewKind("trait"),
							MemberName: NewName("Trait"),
							Fields: []*Field{
								{
									MemberKind: NewKind("field"),
									MemberName: NewName("f"),
								},
							},
							Functions: []*Function{
								{
									MemberKind: NewKind("function"),
									MemberName: NewName("func"),
								},
							},
						},
					},
					Functions: []*Function{
						{
							MemberKind: NewKind("function"),
							MemberName: NewName("func"),
							Overloads: []*Function{
								{
									MemberKind: NewKind("function"),
									MemberName: NewName("func"),
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

	p := collectPaths(&docs)
	assert.Equal(t, 11, len(p))

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
		obs, ok := p[tt.mem]
		assert.Equal(t, tt.ok, ok)
		assert.Equal(t, tt.exp, obs.Elements)
	}
}
