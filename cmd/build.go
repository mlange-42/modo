package cmd

import (
	"fmt"

	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func buildCommand() (*cobra.Command, error) {
	v := viper.New()

	root := &cobra.Command{
		Use:   "build [PATH]",
		Short: "Build documentation from 'mojo doc' JSON",
		Long: `Build documentation from 'mojo doc' JSON.

Builds based on the 'modo.yaml' file in the current directory if no path is given.
The flags listed below overwrite the settings from that file.

Complete documentation at https://mlange-42.github.io/modo/`,
		Example: `  modo init                      # set up a project
  mojo doc src/ -o api.json      # run 'mojo doc'
  modo build                     # build the docs`,
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
			return runBuild(cliArgs)
		},
	}

	root.Flags().StringSliceP("input", "i", []string{}, "'mojo doc' JSON file to process. Reads from STDIN if not specified")
	root.Flags().StringP("output", "o", "", "Output folder for generated Markdown files")
	root.Flags().StringP("tests", "t", "", "Target folder to extract doctests for 'mojo test'.\nSee also command 'modo test' (default no doctests)")
	root.Flags().StringP("format", "f", "plain", "Output format. One of (plain|mdbook|hugo)")
	root.Flags().BoolP("exports", "e", false, "Process according to 'Exports:' sections in packages")
	root.Flags().BoolP("short-links", "s", false, "Render shortened link labels, stripping packages and modules")
	root.Flags().BoolP("case-insensitive", "C", false, "Build for systems that are not case-sensitive regarding file names.\nAppends hyphen (-) to capitalized file names")
	root.Flags().BoolP("strict", "S", false, "Strict mode. Errors instead of warnings")
	root.Flags().BoolP("dry-run", "D", false, "Dry-run without any file output")
	root.Flags().BoolP("bare", "B", false, "Don't run ore- and post-commands")
	root.Flags().StringSliceP("templates", "T", []string{}, "Optional directories with templates for (partial) overwrite.\nSee folder assets/templates in the repository")

	root.Flags().SortFlags = false
	root.MarkFlagFilename("input", "json")
	root.MarkFlagDirname("output")
	root.MarkFlagDirname("tests")
	root.MarkFlagDirname("templates")

	err := v.BindPFlags(root.Flags())
	if err != nil {
		return nil, err
	}
	return root, nil
}

func runBuild(args *document.Config) error {
	if args.OutputDir == "" {
		return fmt.Errorf("no output path given")
	}

	if !args.Bare {
		if err := runPreBuildCommands(args); err != nil {
			return err
		}
	}

	formatter, err := format.GetFormatter(args.RenderFormat)
	if err != nil {
		return err
	}

	if len(args.InputFiles) == 0 || (len(args.InputFiles) == 1 && args.InputFiles[0] == "") {
		if err := runBuildOnce("", args, formatter); err != nil {
			return err
		}
	} else {
		for _, f := range args.InputFiles {
			if err := runBuildOnce(f, args, formatter); err != nil {
				return err
			}
		}
	}

	if !args.Bare {
		if err := runPostBuildCommands(args); err != nil {
			return err
		}
	}

	return nil
}

func runBuildOnce(file string, args *document.Config, form document.Formatter) error {
	docs, err := readDocs(file)
	if err != nil {
		return err
	}
	err = form.Render(docs, args)
	if err != nil {
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
