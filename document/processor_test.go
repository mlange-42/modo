package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderAll(t *testing.T) {
	yml := `
decl:
  name: modo
  kind: package
  summary: Package modo
  description: |
    Exports:
     - mod1.Struct1
     - mod2
  modules:
    - name: mod1
      kind: module
      structs:
        - name: Struct1
          kind: struct
    - name: mod2
      kind: module
      aliases:
        - name: Alias1
          kind: alias
      structs:
        - name: Struct2
          kind: struct
      traits:
        - name: Trait2
          kind: trait
      functions:
        - name: func2
          kind: function
          overloads:
            - name: func2
              kind: function
`
	docs, err := FromYaml([]byte(yml))
	assert.Nil(t, err)
	assert.NotNil(t, docs)

	formatter := TestFormatter{}
	proc := NewProcessor(docs, &formatter, nil, true, true)
	err = proc.PrepareDocs()
	assert.Nil(t, err)
}
