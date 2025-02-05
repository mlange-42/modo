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
 - mod.func as f

` +
		"```mojo\n" +
		`Exports:
 - xxx.Struct
 - xxx.Trait
 - xxx.func as f
` +
		"```\n" +
		`
Text

Exports:

 - mod.submod

Text
`

	proc := NewProcessor(nil, nil, nil, &Config{Strict: true})
	exports, newText, anyExp, err := proc.parseExports(text, []string{"pkg"}, true)

	assert.Nil(t, err)
	assert.True(t, anyExp)

	assert.Equal(t, []*packageExport{
		{Short: []string{"mod", "Struct"}, Exported: []string{"pkg", "Struct"}, Renamed: "Struct", Long: []string{"pkg", "mod", "Struct"}},
		{Short: []string{"mod", "Trait"}, Exported: []string{"pkg", "Trait"}, Renamed: "Trait", Long: []string{"pkg", "mod", "Trait"}},
		{Short: []string{"mod", "func"}, Exported: []string{"pkg", "func"}, Renamed: "f", Long: []string{"pkg", "mod", "func"}},
		{Short: []string{"mod", "submod"}, Exported: []string{"pkg", "submod"}, Renamed: "submod", Long: []string{"pkg", "mod", "submod"}},
	}, exports)

	assert.Equal(t, newText, `Text.
Text
   indented text


`+
		"```mojo\n"+
		`Exports:
 - xxx.Struct
 - xxx.Trait
 - xxx.func as f
`+
		"```\n"+
		`
Text


Text`)
}
