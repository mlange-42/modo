package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// RootCommand creates the root command.
func RootCommand() (*cobra.Command, error) {
	var showVersion bool

	root := &cobra.Command{
		Use:   "modo",
		Short: "Modo -- DocGen for Mojo.",
		Long: `Modo -- DocGen for Mojo.

Modo generates Markdown for static site generators (SSGs) from 'mojo doc' JSON output.

Complete documentation at https://mlange-42.github.io/modo/`,
		Example: `  modo init hugo                 # set up a project, e.g. for Hugo
  modo build                     # build the docs`,
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if showVersion {
				info, _ := debug.ReadBuildInfo()
				fmt.Printf("Modo %s\n", info.Main.Version)
				return nil
			}
			return cmd.Help()
		},
	}

	root.CompletionOptions.HiddenDefaultCmd = true

	for _, fn := range []func(chan struct{}) (*cobra.Command, error){initCommand, buildCommand, testCommand, cleanCommand} {
		cmd, err := fn(nil)
		if err != nil {
			return nil, err
		}
		root.AddCommand(cmd)
	}

	root.Flags().BoolVarP(&showVersion, "version", "V", false, "")

	return root, nil
}
