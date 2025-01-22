package cmd

import (
	"fmt"

	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func testCommand() (*cobra.Command, error) {
	v := viper.New()

	root := &cobra.Command{
		Use:   "test [PATH]",
		Short: "Extract doc-tests from 'mojo doc' JSON",
		Long: `Extract doc-tests from 'mojo doc' JSON.

Extracts tests based on the 'modo.yaml' file in the current directory if no path is given.
The flags listed below overwrite the settings from that file.

Complete documentation at https://mlange-42.github.io/modo/`,
		Example: `  modo init                      # set up a project
  mojo doc src/ -o api.json      # run 'mojo doc'
  modo test                      # extract tests
  mojo test -I . doctest/        # run the tests`,
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
			return runTest(cliArgs)
		},
	}

	root.Flags().StringSliceP("input", "i", []string{}, "'mojo doc' JSON file to process. Reads from STDIN if not specified")
	root.Flags().StringP("tests", "t", "", "Target folder to extract doctests for 'mojo test'")
	root.Flags().BoolP("case-insensitive", "C", false, "Build for systems that are not case-sensitive regarding file names.\nAppends hyphen (-) to capitalized file names")
	root.Flags().BoolP("strict", "S", false, "Strict mode. Errors instead of warnings")
	root.Flags().BoolP("dry-run", "D", false, "Dry-run without any file output")
	root.Flags().BoolP("bare", "B", false, "Don't run ore- and post-commands")
	root.Flags().StringSliceP("templates", "T", []string{}, "Optional directories with templates for (partial) overwrite.\nSee folder assets/templates in the repository")

	root.Flags().SortFlags = false
	root.MarkFlagFilename("input", "json")
	root.MarkFlagDirname("tests")
	root.MarkFlagDirname("templates")

	err := v.BindPFlags(root.Flags())
	if err != nil {
		return nil, err
	}
	return root, nil
}

func runTest(args *document.Config) error {
	if args.TestOutput == "" {
		return fmt.Errorf("no output path for tests given")
	}

	if !args.Bare {
		if err := runPreTestCommands(args); err != nil {
			return err
		}
	}

	if len(args.InputFiles) == 0 || (len(args.InputFiles) == 1 && args.InputFiles[0] == "") {
		if err := runTestOnce("", args); err != nil {
			return err
		}
	} else {
		for _, f := range args.InputFiles {
			if err := runTestOnce(f, args); err != nil {
				return err
			}
		}
	}

	if !args.Bare {
		if err := runPostTestCommands(args); err != nil {
			return err
		}
	}

	return nil
}

func runTestOnce(file string, args *document.Config) error {
	docs, err := readDocs(file)
	if err != nil {
		return err
	}
	if err := document.ExtractTests(docs, args, &format.PlainFormatter{}); err != nil {
		return err
	}
	return nil
}

func runPreTestCommands(cfg *document.Config) error {
	if err := runCommands(cfg.PreRun); err != nil {
		return err
	}
	if err := runCommands(cfg.PreTest); err != nil {
		return err
	}
	return nil
}

func runPostTestCommands(cfg *document.Config) error {
	if err := runCommands(cfg.PostTest); err != nil {
		return err
	}
	if err := runCommands(cfg.PostRun); err != nil {
		return err
	}
	return nil
}
