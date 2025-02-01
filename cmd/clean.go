package cmd

import (
	"fmt"

	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func cleanCommand() (*cobra.Command, error) {
	v := viper.New()

	root := &cobra.Command{
		Use:   "clean [PATH]",
		Short: "Remove Modo's output files",
		Long: `Remove Modo's output files.

Complete documentation at https://mlange-42.github.io/modo/`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := mountProject(v, args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs, err := document.ConfigFromViper(v)
			if err != nil {
				return err
			}
			return runClean(cliArgs)
		},
	}

	return root, nil
}

func runClean(args *document.Config) error {
	if args.OutputDir == "" {
		return fmt.Errorf("no output path given")
	}

	formatter, err := format.GetFormatter(args.RenderFormat)
	if err != nil {
		return err
	}

	return formatter.Clean(args)
}
