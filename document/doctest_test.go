package document

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBlockAttributes(t *testing.T) {
	tests := []struct {
		Text, Name              string
		Hide, Global, Ok, Error bool
	}{
		{"```mojo",
			"", false, false, false, false},
		{"```mojo {doctest=\"test\" hide=true global=true}",
			"test", true, true, true, false},
		{"```mojo { doctest=\"test\" hide=true global=true }",
			"test", true, true, true, false},
		{"```mojo {doctest=\"test\"}",
			"test", false, false, true, false},
		{"```mojo {other=\"abc\" doctest=\"test\"}",
			"test", false, false, true, false},
		{"```mojo {.class1 doctest=\"test\" class2}",
			"test", false, false, true, false},
		{"```mojo {hide=true}",
			"", true, false, true, false},
	}

	for _, test := range tests {
		name, hide, global, ok, err := parseBlockAttr(test.Text)
		assert.Equal(t, name, test.Name, "Name %s", test.Text)
		assert.Equal(t, hide, test.Hide, "Hide %s", test.Text)
		assert.Equal(t, global, test.Global, "Global %s", test.Text)
		assert.Equal(t, ok, test.Ok, "Ok %s", test.Text)
		assert.Equal(t, err != nil, test.Error, "Err %s %s", test.Text, err)
	}
}

func TestExtractDocTests(t *testing.T) {
	text := "Docstring\n" +
		"\n" +
		"```mojo {doctest=\"test\" global=true hide=true}\n" +
		"struct Test:\n" +
		"    pass\n" +
		"```\n" +
		"\n" +
		"Some text\n" +
		"\n" +
		"```mojo {doctest=\"test\" hide=true}\n" +
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
		"```mojo {doctest=\"test\" hide=true}\n" +
		"assert(b == 0)\n" +
		"```\n"

	proc := NewProcessor(nil, nil, nil, &Config{})
	outText, err := proc.extractTests(text, []string{"pkg", "Struct"}, 1)
	assert.Nil(t, err)
	assert.Equal(t, 14, len(strings.Split(outText, "\n")))

	assert.Equal(t, 1, len(proc.docTests))
	assert.Equal(t, proc.docTests[0], &docTest{
		Name: "test",
		Path: []string{"pkg", "Struct"},
		Code: []string{
			"import b",
			"var a = b",
			"assert(b == 0)",
		},
		Global: []string{"struct Test:", "    pass"},
	})
}
