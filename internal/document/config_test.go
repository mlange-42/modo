package document

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfigRemovePostScripts(t *testing.T) {
	config := Config{
		PreRun:    []string{"test"},
		PreBuild:  []string{"test"},
		PreTest:   []string{"test"},
		PostRun:   []string{"test"},
		PostBuild: []string{"test"},
		PostTest:  []string{"test"},
	}
	config.RemovePostScripts()

	assert.Nil(t, config.PostRun)
	assert.Nil(t, config.PostBuild)
	assert.Nil(t, config.PostTest)

	assert.NotNil(t, config.PreRun)
	assert.NotNil(t, config.PreBuild)
	assert.NotNil(t, config.PreTest)
}

func TestConfigFromViper(t *testing.T) {
	v := viper.New()

	v.Set("output", "test")
	config, err := ConfigFromViper(v)
	assert.Nil(t, err)
	assert.Equal(t, "test", config.OutputDir)

	v.Set("unknown", "test")
	config, err = ConfigFromViper(v)
	assert.NotNil(t, err)
	assert.Nil(t, config)
}
