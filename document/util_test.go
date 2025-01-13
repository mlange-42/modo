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
	text := "A [Struct] and a [Struct.member]. And a [Markdown](link)."
	lookup := map[string]elemPath{
		"stdlib.Struct":        {Elements: []string{"stdlib", "Struct"}, Kind: "member"},
		"stdlib.Struct.member": {Elements: []string{"stdlib", "Struct", "#member"}, Kind: "member"},
	}
	elems := []string{"stdlib"}
	templ := template.New("all").Funcs(template.FuncMap{"pathJoin": path.Join})
	templ, err := templ.Parse(`{{define "member_path.md"}}{{.}}.md{{end}}`)
	assert.Nil(t, err)

	out, err := replaceLinks(text, elems, lookup, templ)
	assert.Nil(t, err)

	assert.Equal(t, "A [Struct](Struct.md) and a [Struct.member](Struct.md#member). And a [Markdown](link).", out)
}
