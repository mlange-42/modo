package main

import (
	"fmt"

	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func testCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "test [OUT-PATH]",
		Short: "Generate tests from 'mojo doc' JSON.",
		Example: `  modo test doctest -i docs.json        # from a file
  mojo doc ./src | modo test doctest    # from 'mojo doc'`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			viper.SetConfigName("modo")
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
				cliArgs.DocTests = args[0]
			} else {
				if cliArgs.OutputDir == "" {
					return fmt.Errorf("missing output directory argument")
				}
			}
			return runTest(&cliArgs)
		},
	}

	root.Flags().StringP("input", "i", "", "'mojo doc' JSON file to process. Reads from STDIN if not specified.")
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

func runTest(args *document.Config) error {
	if args.DocTests == "" {
		return fmt.Errorf("no output path for tests given")
	}
	for _, command := range args.PreTest {
		err := runCommand(command)
		if err != nil {
			return err
		}
	}

	docs, err := readDocs(args.InputFile)
	if err != nil {
		return err
	}
	if err := document.ExtractTests(docs, args, &format.PlainFormatter{}); err != nil {
		return err
	}

	for _, command := range args.PostTest {
		err := runCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
}
