package document

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindLinks(t *testing.T) {
	text := "⌘a [link1].\n" +
		"a `[link2] in inline` code\n" +
		"and finally...\n" +
		"```mojo\n" +
		"a [link3] in a code block\n" +
		"```\n" +
		"and a normal [link4] again.\n" +
		"and a doc inheritance [[link4]].\n"
	indices, err := findLinks(text, linkRegex, true)
	assert.Nil(t, err)
	assert.NotNil(t, indices)
	assert.Equal(t, 4, len(indices))
	assert.Equal(t, "[link1]", text[indices[0]:indices[1]])
	assert.Equal(t, "[link4]", text[indices[2]:indices[3]])
}

func TestTranscludesLinks(t *testing.T) {
	text := "⌘a [[link1]].\n" +
		"a `[[link2]] in inline` code\n" +
		"and finally...\n" +
		"```mojo\n" +
		"a [[link3]] in a code block\n" +
		"```\n" +
		"and a normal [[link4]] again.\n" +
		"and a doc inheritance [link5].\n"
	indices, err := findLinks(text, transcludeRegex, false)
	assert.Nil(t, err)
	assert.NotNil(t, indices)
	assert.Equal(t, 4, len(indices))
	assert.Equal(t, "[link1]", text[indices[0]:indices[1]])
	assert.Equal(t, "[link4]", text[indices[2]:indices[3]])
}

func TestReplaceRefs(t *testing.T) {
	text := "A [.Struct], a [.Struct.member], a [..Trait], a [.q.func], abs [stdlib.Trait], [stdlib]. And a [Markdown](link)."
	lookup := map[string]string{
		"stdlib.Trait":           "stdlib.Trait",
		"stdlib.p.Struct":        "stdlib.p.Struct",
		"stdlib.p.Struct.member": "stdlib.p.Struct.member",
		"stdlib.p.q.func":        "stdlib.p.q.func",
		"stdlib":                 "stdlib",
	}
	elems := []string{"stdlib", "p", "Struct"}

	proc := NewProcessor(nil, &TestFormatter{}, nil, &Config{})
	proc.linkExports = lookup
	out, err := proc.replaceRefs(text, elems, 2)
	assert.Nil(t, err)

	assert.Equal(t, "A [stdlib.p.Struct], a [stdlib.p.Struct.member], a [stdlib.Trait], a [stdlib.p.q.func], abs [stdlib.Trait], [stdlib]. And a [Markdown](link).", out)
}

func TestReplacePlaceholders(t *testing.T) {
	text := "A [stdlib.p.Struct], a [stdlib.p.Struct.member], a [stdlib.Trait], a [stdlib.p.q.func], abs [stdlib.Trait], [stdlib]. And a [Markdown](link)."
	lookup := map[string]elemPath{
		"stdlib.Trait":           {Elements: []string{"stdlib", "Trait"}, Kind: "member"},
		"stdlib.p.Struct":        {Elements: []string{"stdlib", "p", "Struct"}, Kind: "member"},
		"stdlib.p.Struct.member": {Elements: []string{"stdlib", "p", "Struct", "#member"}, Kind: "member", IsSection: true},
		"stdlib.p.q.func":        {Elements: []string{"stdlib", "p", "q", "func"}, Kind: "member"},
		"stdlib":                 {Elements: []string{"stdlib"}, Kind: "package"},
	}
	elems := []string{"stdlib", "p", "Struct"}

	proc := NewProcessor(nil, &TestFormatter{}, nil, &Config{ShortLinks: true})
	proc.linkTargets = lookup
	out, err := proc.ReplacePlaceholders(text, elems, 2)
	assert.Nil(t, err)

	fmt.Println(out)
	assert.Equal(t, "A [`Struct`](Struct.md), a [`Struct.member`](Struct.md#member), a [`Trait`](../Trait.md), a [`func`](q/func.md), abs [`Trait`](../Trait.md), [`stdlib`](../_index.md). And a [Markdown](link).", out)
}
