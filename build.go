package main

import (
	"fmt"

	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func buildCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "build [OUT-PATH]",
		Short: "Build documentation from 'mojo doc' JSON.",
		Example: `  modo build docs -i docs.json        # from a file
  mojo doc ./src | modo build docs    # from 'mojo doc'`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigName(configFile)
			viper.SetConfigType("yaml")
			viper.AddConfigPath(".")
			if err := viper.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
					return err
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs := document.Config{}
			err := viper.Unmarshal(&cliArgs)
			if err != nil {
				return err
			}
			if len(args) > 0 {
				cliArgs.OutputDir = args[0]
			} else {
				if cliArgs.OutputDir == "" {
					return fmt.Errorf("missing output directory argument")
				}
			}
			return runBuild(&cliArgs)
		},
	}

	root.Flags().StringP("input", "i", "", "'mojo doc' JSON file to process. Reads from STDIN if not specified.")
	root.Flags().StringP("doctest", "d", "", "Target folder to extract doctests for 'mojo test'. (default no doctests)")
	root.Flags().StringP("format", "f", "plain", "Output format. One of (plain|mdbook|hugo).")
	root.Flags().BoolP("exports", "e", false, "Process according to 'Exports:' sections in packages.")
	root.Flags().Bool("short-links", false, "Render shortened link labels, stripping packages and modules.")
	root.Flags().Bool("case-insensitive", false, "Build for systems that are not case-sensitive regarding file names.\nAppends hyphen (-) to capitalized file names.")
	root.Flags().Bool("strict", false, "Strict mode. Errors instead of warnings.")
	root.Flags().Bool("dry-run", false, "Dry-run without any file output.")
	root.Flags().StringSliceP("templates", "t", []string{}, "Optional directories with templates for (partial) overwrite.\nSee folder assets/templates in the repository.")

	root.Flags().SortFlags = false
	root.MarkFlagFilename("input", "json")
	root.MarkFlagDirname("templates")

	viper.BindPFlags(root.Flags())

	return root
}

func runBuild(args *document.Config) error {
	if args.OutputDir == "" {
		return fmt.Errorf("no output path given")
	}

	if err := runPreBuildCommands(args); err != nil {
		return err
	}

	docs, err := readDocs(args.InputFile)
	if err != nil {
		return err
	}
	formatter, err := format.GetFormatter(args.RenderFormat)
	if err != nil {
		return err
	}
	err = formatter.Render(docs, args)
	if err != nil {
		return err
	}

	if err := runPostBuildCommands(args); err != nil {
		return err
	}

	return nil
}

func runPreBuildCommands(cfg *document.Config) error {
	if err := runCommands(cfg.PreRun); err != nil {
		return err
	}
	if err := runCommands(cfg.PreBuild); err != nil {
		return err
	}
	if cfg.DocTests != "" {
		if err := runCommands(cfg.PreTest); err != nil {
			return err
		}
	}
	return nil
}

func runPostBuildCommands(cfg *document.Config) error {
	if cfg.DocTests != "" {
		if err := runCommands(cfg.PostTest); err != nil {
			return err
		}
	}
	if err := runCommands(cfg.PostBuild); err != nil {
		return err
	}
	if err := runCommands(cfg.PostRun); err != nil {
		return err
	}
	return nil
}
