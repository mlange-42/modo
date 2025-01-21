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
		Use:   "build",
		Short: "Build documentation from 'mojo doc' JSON.",
		Example: `  modo build -i api.json -o docs/        # from a file
  mojo doc ./src | modo build -o docs/    # from 'mojo doc'`,
		Args:         cobra.ExactArgs(0),
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
			return runBuild(&cliArgs)
		},
	}

	root.Flags().StringP("input", "i", "", "'mojo doc' JSON file to process. Reads from STDIN if not specified.")
	root.Flags().StringP("output", "o", "", "Output folder for generated Markdown files.")
	root.Flags().StringP("tests", "t", "", "Target folder to extract doctests for 'mojo test'.\nSee also command 'modo test'. (default no doctests)")
	root.Flags().StringP("format", "f", "plain", "Output format. One of (plain|mdbook|hugo).")
	root.Flags().BoolP("exports", "e", false, "Process according to 'Exports:' sections in packages.")
	root.Flags().BoolP("short-links", "s", false, "Render shortened link labels, stripping packages and modules.")
	root.Flags().BoolP("case-insensitive", "C", false, "Build for systems that are not case-sensitive regarding file names.\nAppends hyphen (-) to capitalized file names.")
	root.Flags().BoolP("strict", "S", false, "Strict mode. Errors instead of warnings.")
	root.Flags().BoolP("dry-run", "D", false, "Dry-run without any file output.")
	root.Flags().StringSliceP("templates", "T", []string{}, "Optional directories with templates for (partial) overwrite.\nSee folder assets/templates in the repository.")

	root.Flags().SortFlags = false
	root.MarkFlagFilename("input", "json")
	root.MarkFlagDirname("output")
	root.MarkFlagDirname("tests")
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

	formatter, err := format.GetFormatter(args.RenderFormat)
	if err != nil {
		return err
	}

	for _, f := range args.InputFiles {
		docs, err := readDocs(f)
		if err != nil {
			return err
		}
		err = formatter.Render(docs, args)
		if err != nil {
			return err
		}
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
	if cfg.TestOutput != "" {
		if err := runCommands(cfg.PreTest); err != nil {
			return err
		}
	}
	return nil
}

func runPostBuildCommands(cfg *document.Config) error {
	if cfg.TestOutput != "" {
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
