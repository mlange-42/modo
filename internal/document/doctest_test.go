package document

import (
	"os"
	"path"
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
		{"```mojo {doctest=\"test\" hide=true, global=true}",
			"test", false, false, false, true},
		{"```mojo {doctest=\"test\" global=true, hide=true}",
			"test", false, false, false, true},
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

func TestExtractDocTests4Ticks(t *testing.T) {
	text := "Docstring\n" +
		"\n" +
		"```mojo {doctest=\"test1\" hide=true global=true}\n" +
		"````mojo {doctest=\"test2\"}\n" +
		"Test1\n" +
		"````\n" +
		"```\n" +
		"\n" +
		"````mojo {doctest=\"test3\" hide=true global=true}\n" +
		"```mojo {doctest=\"test4\"}\n" +
		"Test2\n" +
		"```\n" +
		"````\n"

	proc := NewProcessor(nil, nil, nil, &Config{})
	outText, err := proc.extractTests(text, []string{"pkg", "Struct"}, 1)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(strings.Split(outText, "\n")))

	assert.Equal(t, 2, len(proc.docTests))

	assert.Equal(t, &docTest{
		Name:   "test1",
		Path:   []string{"pkg", "Struct"},
		Code:   []string{},
		Global: []string{"````mojo {doctest=\"test2\"}", "Test1", "````"},
	}, proc.docTests[0])

	assert.Equal(t, &docTest{
		Name:   "test3",
		Path:   []string{"pkg", "Struct"},
		Code:   []string{},
		Global: []string{"```mojo {doctest=\"test4\"}", "Test2", "```"},
	}, proc.docTests[1])
}

func TestExtractTestsMarkdown(t *testing.T) {
	inDir := t.TempDir()
	outDir := t.TempDir()
	testDir := t.TempDir()
	_ = outDir

	config := Config{
		InputFiles: []string{inDir},
		OutputDir:  outDir,
		TestOutput: testDir,
		Strict:     true,
	}

	err := os.WriteFile(path.Join(inDir, "_index.md"), []byte(
		"```mojo {doctest=\"test1\"}\n"+
			"var a = 0\n"+
			"```\n",
	), 0666)
	assert.Nil(t, err)

	err = ExtractTestsMarkdown(&config, &TestFormatter{}, inDir, true)
	assert.Nil(t, err)

	_, err = os.Stat(path.Join(outDir, "_index.md"))
	assert.Nil(t, err)
	_, err = os.Stat(path.Join(testDir, "_index_test1_test.mojo"))
	assert.Nil(t, err)
}
