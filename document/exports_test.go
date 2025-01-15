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
	proc := NewProcessor(nil, nil, false)
	exports := proc.parseExports(text)

	assert.Equal(t, []string{"mod.Struct", "mod.Trait", "mod.func", "mod.submod"}, exports)
}
