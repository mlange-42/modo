package document

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	InputFiles      []string `mapstructure:"input" yaml:"input"`
	OutputDir       string   `mapstructure:"output" yaml:"output"`
	TestOutput      string   `mapstructure:"tests" yaml:"tests"`
	RenderFormat    string   `mapstructure:"format" yaml:"format"`
	UseExports      bool     `mapstructure:"exports" yaml:"exports"`
	ShortLinks      bool     `mapstructure:"short-links" yaml:"short-links"`
	Strict          bool     `mapstructure:"strict" yaml:"strict"`
	DryRun          bool     `mapstructure:"dry-run" yaml:"dry-run"`
	CaseInsensitive bool     `mapstructure:"case-insensitive" yaml:"case-insensitive"`
	Bare            bool     `mapstructure:"bare" yaml:"bare"`
	TemplateDirs    []string `mapstructure:"templates" yaml:"templates"`
	PreRun          []string `mapstructure:"pre-run" yaml:"pre-run"`
	PreBuild        []string `mapstructure:"pre-build" yaml:"pre-build"`
	PreTest         []string `mapstructure:"pre-test" yaml:"pre-test"`
	PostTest        []string `mapstructure:"post-test" yaml:"post-test"`
	PostBuild       []string `mapstructure:"post-build" yaml:"post-build"`
	PostRun         []string `mapstructure:"post-run" yaml:"post-run"`
}

func ConfigFromViper(v *viper.Viper) (*Config, error) {
	c := Config{}
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}
	if err := c.check(v); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) check(v *viper.Viper) error {
	fields := map[string]bool{}
	tp := reflect.TypeOf(c).Elem()
	for i := range tp.NumField() {
		field := tp.Field(i)
		if tag, ok := field.Tag.Lookup("mapstructure"); ok {
			if _, ok := fields[tag]; ok {
				return fmt.Errorf("duplicate field name '%s'", tag)
			}
			fields[tag] = true
		}
	}
	for _, key := range v.AllKeys() {
		if _, ok := fields[key]; !ok {
			return fmt.Errorf("unknown field '%s' in config file", key)
		}
	}
	return nil
}
