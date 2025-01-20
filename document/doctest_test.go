package document

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractDocTests(t *testing.T) {
	text := "Docstring\n" +
		"\n" +
		"```mojo {doctest=\"test-setup\"}\n" +
		"import b\n" +
		"```\n" +
		"\n" +
		"Some text\n" +
		"\n" +
		"```mojo {doctest=\"test\"}\n" +
		"var a = b\n" +
		"```\n" +
		"\n" +
		"Some text\n" +
		"\n" +
		"```mojo {doctest=\"test-teardown\"}\n" +
		"assert(b == 0)\n" +
		"```\n"

	proc := NewProcessor(nil, nil, nil, &Config{})
	outText, err := proc.extractTests(text, []string{"pkg", "Struct"}, 1)
	assert.Nil(t, err)
	assert.Equal(t, 12, len(strings.Split(outText, "\n")))

	assert.Equal(t, 1, len(proc.docTests))
	assert.Equal(t, proc.docTests[0], &docTest{
		Name: "test",
		Path: []string{"pkg", "Struct"},
		Code: []string{
			"import b",
			"var a = b",
			"assert(b == 0)",
		},
	})
}
