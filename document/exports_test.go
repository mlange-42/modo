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
	proc := NewProcessor(nil, nil, false, false)
	exports := proc.parseExports(text, "pkg")

	assert.Equal(t, []PackageExport{
		{Short: "mod.Struct", Long: "pkg.mod.Struct"},
		{Short: "mod.Trait", Long: "pkg.mod.Trait"},
		{Short: "mod.func", Long: "pkg.mod.func"},
		{Short: "mod.submod", Long: "pkg.mod.submod"},
	}, exports)
}
