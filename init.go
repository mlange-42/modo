package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/mlange-42/modo/document"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var defaultConfig = document.Config{
	InputFile:    "docs.json",
	OutputDir:    "docs/",
	DocTests:     "doctest/",
	RenderFormat: "plain",
	ShortLinks:   true,
	UseExports:   true,
	PreBuild: []string{`echo Use pre- and post-commands...
echo to run 'mojo doc' or 'mojo test'.
`},
}

func initCommand() *cobra.Command {
	root := &cobra.Command{
		Use:          "init",
		Short:        "Generate Modo config file in the current directory.",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			file := configFile + ".yaml"
			if _, err := os.Stat(file); err == nil {
				return fmt.Errorf("config file %s already exists", file)
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("error checking config file %s: %s", file, err.Error())
			}

			b := bytes.Buffer{}
			enc := yaml.NewEncoder(&b)
			enc.SetIndent(2)

			if err := enc.Encode(&defaultConfig); err != nil {
				return err
			}
			if err := os.WriteFile(file, b.Bytes(), 0644); err != nil {
				return err
			}
			fmt.Println("Modo project initialized.")
			return nil
		},
	}

	return root
}
