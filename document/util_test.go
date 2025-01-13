package document

import (
	"fmt"
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
	text := "A [Struct], a [Struct.member], a [.Trait], a [q.func]. And a [Markdown](link)."
	lookup := map[string]elemPath{
		"stdlib.Trait":           {Elements: []string{"stdlib", "Trait"}, Kind: "member"},
		"stdlib.p.Struct":        {Elements: []string{"stdlib", "p", "Struct"}, Kind: "member"},
		"stdlib.p.Struct.member": {Elements: []string{"stdlib", "p", "Struct", "#member"}, Kind: "member"},
		"stdlib.p.q.func":        {Elements: []string{"stdlib", "p", "q", "func"}, Kind: "member"},
	}
	elems := []string{"stdlib", "p"}
	templ := template.New("all").Funcs(template.FuncMap{"pathJoin": path.Join})
	templ, err := templ.Parse(`{{define "member_path.md"}}{{.}}.md{{end}}`)
	assert.Nil(t, err)

	out, err := replaceLinks(text, elems, lookup, templ)
	assert.Nil(t, err)

	fmt.Println(out)
	assert.Equal(t, "A [Struct](Struct.md), a [Struct.member](Struct.md#member), a [Trait](../Trait.md), a [q.func](q/func.md). And a [Markdown](link).", out)
}
