package document

import (
	"path"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestFindLinks(t *testing.T) {
	text := "⌘a [link1].\n" +
		"a `[link2] in inline` code\n" +
		"and finally...\n" +
		"```mojo\n" +
		"a [link3] in a code block\n" +
		"```\n" +
		"and a normal [link4] again.\n"
	indices, err := findLinks(text)
	assert.Nil(t, err)
	assert.NotNil(t, indices)
	assert.Equal(t, 4, len(indices))
	assert.Equal(t, "[link1]", text[indices[0]:indices[1]])
	assert.Equal(t, "[link4]", text[indices[2]:indices[3]])
}

func TestReplaceLinks(t *testing.T) {
	text := "A [.Struct], a [.Struct.member], a [..Trait], a [.q.func], abs [stdlib.Trait]. And a [Markdown](link)."
	lookup := map[string]elemPath{
		"stdlib.Trait":           {Elements: []string{"stdlib", "Trait"}, Kind: "member"},
		"stdlib.p.Struct":        {Elements: []string{"stdlib", "p", "Struct"}, Kind: "member"},
		"stdlib.p.Struct.member": {Elements: []string{"stdlib", "p", "Struct", "#member"}, Kind: "member", IsSection: true},
		"stdlib.p.q.func":        {Elements: []string{"stdlib", "p", "q", "func"}, Kind: "member"},
	}
	elems := []string{"stdlib", "p", "Struct"}
	templ := template.New("all").Funcs(template.FuncMap{"pathJoin": path.Join})
	templ, err := templ.Parse(`{{define "member_path.md"}}{{.}}.md{{end}}`)
	assert.Nil(t, err)

	out, err := replaceLinks(text, elems, 2, lookup, templ)
	assert.Nil(t, err)

	assert.Equal(t, "A [`Struct`](Struct.md), a [`Struct.member`](Struct.md#member), a [`Trait`](../Trait.md), a [`q.func`](q/func.md), abs [`stdlib.Trait`](../Trait.md). And a [Markdown](link).", out)
}

func TestToRelLink(t *testing.T) {
	lookup := map[string]elemPath{
		"stdlib.Trait":           {Elements: []string{"stdlib", "Trait"}, Kind: "member", IsSection: false},
		"stdlib.p.Struct":        {Elements: []string{"stdlib", "p", "Struct"}, Kind: "member", IsSection: false},
		"stdlib.p.Struct2":       {Elements: []string{"stdlib", "p", "Struct2"}, Kind: "member", IsSection: false},
		"stdlib.q.Struct3":       {Elements: []string{"stdlib", "q", "Struct"}, Kind: "member", IsSection: false},
		"stdlib.q.Struct4":       {Elements: []string{"stdlib", "q", "Struct2"}, Kind: "member", IsSection: false},
		"stdlib.p.Struct.member": {Elements: []string{"stdlib", "p", "Struct", "#member"}, Kind: "member", IsSection: true},
		"stdlib.p.q.func":        {Elements: []string{"stdlib", "p", "q", "func"}, Kind: "member", IsSection: false},
	}
	elems := []string{"stdlib", "p", "Struct"}

	_, text, link, ok := toLink("..q.Struct3", elems, 2, lookup)
	assert.True(t, ok)
	assert.Equal(t, "`q.Struct3`", text)
	assert.Equal(t, []string{"..", "q", "Struct"}, link)

	_, _, link, ok = toLink(".Struct2", elems, 2, lookup)
	assert.True(t, ok)
	assert.Equal(t, []string{"Struct2"}, link)

	_, _, link, ok = toLink(".Struct.member", elems, 2, lookup)
	assert.True(t, ok)
	assert.Equal(t, []string{"Struct", "#member"}, link)

	_, _, link, ok = toLink("..Trait", elems, 2, lookup)
	assert.True(t, ok)
	assert.Equal(t, []string{"..", "Trait"}, link)

	_, text, link, ok = toLink(".q.func", elems, 2, lookup)
	assert.True(t, ok)
	assert.Equal(t, "`q.func`", text)
	assert.Equal(t, []string{"q", "func"}, link)
}

func TestToAbsLink(t *testing.T) {
	lookup := map[string]elemPath{
		"stdlib.Trait":           {Elements: []string{"stdlib", "Trait"}, Kind: "member", IsSection: false},
		"stdlib.p.Struct":        {Elements: []string{"stdlib", "p", "Struct"}, Kind: "member", IsSection: false},
		"stdlib.p.Struct2":       {Elements: []string{"stdlib", "p", "Struct2"}, Kind: "member", IsSection: false},
		"stdlib.q.Struct3":       {Elements: []string{"stdlib", "q", "Struct"}, Kind: "member", IsSection: false},
		"stdlib.q.Struct4":       {Elements: []string{"stdlib", "q", "Struct2"}, Kind: "member", IsSection: false},
		"stdlib.p.Struct.member": {Elements: []string{"stdlib", "p", "Struct", "#member"}, Kind: "member", IsSection: true},
		"stdlib.p.q.func":        {Elements: []string{"stdlib", "p", "q", "func"}, Kind: "member", IsSection: false},
	}
	elems := []string{"stdlib", "p", "Struct"}

	_, text, link, ok := toLink("stdlib.q.Struct3 S3", elems, 2, lookup)
	assert.True(t, ok)
	assert.Equal(t, "S3", text)
	assert.Equal(t, []string{"..", "q", "Struct"}, link)

	_, _, _, ok = toLink("", elems, 2, lookup)
	assert.False(t, ok)

	_, _, link, ok = toLink("stdlib.p.Struct2", elems, 2, lookup)
	assert.True(t, ok)
	assert.Equal(t, []string{"Struct2"}, link)

	_, _, link, ok = toLink("stdlib.p.Struct.member", elems, 2, lookup)
	assert.True(t, ok)
	assert.Equal(t, []string{"Struct", "#member"}, link)

	_, _, link, ok = toLink("stdlib.Trait", elems, 2, lookup)
	assert.True(t, ok)
	assert.Equal(t, []string{"..", "Trait"}, link)

	_, text, link, ok = toLink("stdlib.p.q.func", elems, 2, lookup)
	assert.True(t, ok)
	assert.Equal(t, "`stdlib.p.q.func`", text)
	assert.Equal(t, []string{"q", "func"}, link)
}