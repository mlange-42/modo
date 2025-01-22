package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseExports(t *testing.T) {
	text := `Text.
Text
   indented text

Exports:
 - mod.Struct
 - mod.Trait
 - mod.func

` +
		"```python\n" +
		`Exports:
 - xxx.Struct
 - xxx.Trait
 - xxx.func
` +
		"```\n" +
		`
Text

Exports:

 - mod.submod

Text
`

	proc := NewProcessor(nil, nil, nil, &Config{})
	exports, newText, anyExp := proc.parseExports(text, []string{"pkg"}, true)

	assert.True(t, anyExp)

	assert.Equal(t, []*packageExport{
		{Short: []string{"mod", "Struct"}, Long: []string{"pkg", "mod", "Struct"}},
		{Short: []string{"mod", "Trait"}, Long: []string{"pkg", "mod", "Trait"}},
		{Short: []string{"mod", "func"}, Long: []string{"pkg", "mod", "func"}},
		{Short: []string{"mod", "submod"}, Long: []string{"pkg", "mod", "submod"}},
	}, exports)

	assert.Equal(t, newText, `Text.
Text
   indented text


`+
		"```python\n"+
		`Exports:
 - xxx.Struct
 - xxx.Trait
 - xxx.func
`+
		"```\n"+
		`
Text


Text`)
}
