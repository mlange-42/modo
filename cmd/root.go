package cmd

import (
	"github.com/spf13/cobra"
)

func RootCommand() (*cobra.Command, error) {
	root := &cobra.Command{
		Use:   "modo",
		Short: "Modo -- DocGen for Mojo.",
		Long: `Modo -- DocGen for Mojo.

Modo generates Markdown for static site generators (SSGs) from 'mojo doc' JSON output.

Complete documentation at https://mlange-42.github.io/modo/`,
		Example: `  modo init                      # set up a project
  mojo doc src/ -o api.json      # run 'mojo doc'
  modo build                     # build the docs`,
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
	}

	root.CompletionOptions.HiddenDefaultCmd = true

	for _, fn := range []func() (*cobra.Command, error){initCommand, buildCommand, testCommand} {
		cmd, err := fn()
		if err != nil {
			return nil, err
		}
		root.AddCommand(cmd)
	}

	return root, nil
}
