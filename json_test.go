package modo_test

import (
	"testing"

	"github.com/mlange-24/modo"
	"github.com/stretchr/testify/assert"
)

func TestFromJson(t *testing.T) {
	data := `{
	"decl": {
    	"kind": "package",
      	"name": "modo",
    	"description": "",
      	"summary": "",
      	"modules": [],
      	"packages": []
	},
    "version": "0.1.0"
}`

	docs, err := modo.FromJson([]byte(data))
	assert.Nil(t, err)
	assert.NotNil(t, docs)
}
