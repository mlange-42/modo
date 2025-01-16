package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseExports(t *testing.T) {
	text := `Text.
Text

Exports:
 - mod.Struct
 - mod.Trait
 - mod.func

Text

Exports:

 - mod.submod

Text
`
	proc := NewProcessor(nil, nil, nil, false, false)
	exports, newText, anyExp := proc.parseExports(text, []string{"pkg"}, true)

	assert.True(t, anyExp)

	assert.Equal(t, []*PackageExport{
		{Short: []string{"mod", "Struct"}, Long: []string{"pkg", "mod", "Struct"}},
		{Short: []string{"mod", "Trait"}, Long: []string{"pkg", "mod", "Trait"}},
		{Short: []string{"mod", "func"}, Long: []string{"pkg", "mod", "func"}},
		{Short: []string{"mod", "submod"}, Long: []string{"pkg", "mod", "submod"}},
	}, exports)

	assert.Equal(t, `Text.
Text


Text


Text
`, newText)
}
