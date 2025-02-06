package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/mlange-42/modo/internal/document"
	"github.com/mlange-42/modo/internal/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func buildCommand(stopWatch chan struct{}) (*cobra.Command, error) {
	v := viper.New()
	var config string
	var watch bool

	var cwd string

	root := &cobra.Command{
		Use:   "build [PATH]",
		Short: "Build documentation from 'mojo doc' JSON",
		Long: `Build documentation from 'mojo doc' JSON.

Builds based on the 'modo.yaml' file in the current directory if no path is given.
The flags listed below overwrite the settings from that file.

Complete documentation at https://mlange-42.github.io/modo/`,
		Example: `  modo init hugo                 # set up a project, e.g. for Hugo
  modo build                     # build the docs`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if err = checkConfigFile(config); err != nil {
				return err
			}
			if cwd, err = mountProject(v, config, args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if err := os.Chdir(cwd); err != nil {
					fmt.Println(err)
				}
			}()

			start := time.Now()

			cliArgs, err := document.ConfigFromViper(v)
			if err != nil {
				return err
			}
			if err := runBuild(cliArgs); err != nil {
				return err
			}
			if watch {
				return watchAndRun(cliArgs, runBuild, stopWatch)
			}

			fmt.Printf("Completed in %.1fms 🧯\n", float64(time.Since(start).Microseconds())/1000.0)
			return nil
		},
	}

	root.Flags().StringVarP(&config, "config", "c", defaultConfigFile, "Config file in the working directory to use")
	root.Flags().StringSliceP("input", "i", []string{}, "'mojo doc' JSON file to process. Reads from STDIN if not specified.\nIf a single directory is given, it is processed recursively")
	root.Flags().StringP("output", "o", "", "Output folder for generated Markdown files")
	root.Flags().StringP("tests", "t", "", "Target folder to extract doctests for 'mojo test'.\nSee also command 'modo test' (default no doctests)")
	root.Flags().StringP("format", "f", "plain", "Output format. One of (plain|mdbook|hugo)")
	root.Flags().BoolP("exports", "e", false, "Process according to 'Exports:' sections in packages")
	root.Flags().BoolP("short-links", "s", false, "Render shortened link labels, stripping packages and modules")
	root.Flags().BoolP("report-missing", "M", false, "Report missing docstings and coverage")
	root.Flags().BoolP("case-insensitive", "C", false, "Build for systems that are not case-sensitive regarding file names.\nAppends hyphen (-) to capitalized file names")
	root.Flags().BoolP("strict", "S", false, "Strict mode. Errors instead of warnings")
	root.Flags().BoolP("dry-run", "D", false, "Dry-run without any file output. Disables post-processing scripts")
	root.Flags().BoolP("bare", "B", false, "Don't run pre- and post-processing scripts")
	root.Flags().BoolVarP(&watch, "watch", "W", false, "Re-run on changes of sources and documentation files.\nDisables post-processing scripts after running them once")
	root.Flags().StringSliceP("templates", "T", []string{}, "Optional directories with templates for (partial) overwrite.\nSee folder assets/templates in the repository")

	root.Flags().SortFlags = false
	root.MarkFlagFilename("config", "yaml")
	root.MarkFlagFilename("input", "json")
	root.MarkFlagDirname("output")
	root.MarkFlagDirname("tests")
	root.MarkFlagDirname("templates")

	err := bindFlags(v, root.Flags())
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

	if err := runFilesOrDir(runBuildOnce, args, formatter); err != nil {
		return err
	}

	if !args.Bare && !args.DryRun {
		if err := runPostBuildCommands(args); err != nil {
			return err
		}
	}

	return nil
}

func runBuildOnce(file string, args *document.Config, form document.Formatter, subdir string, isFile, isDir bool) error {
	if isDir {
		if err := document.ExtractTestsMarkdown(args, form, file, true); err != nil {
			return err
		}
		return runDir(file, args, form, runBuildOnce)
	}
	docs, err := readDocs(file)
	if err != nil {
		return err
	}
	err = document.Render(docs, args, form, subdir)
	if err != nil {
		return err
	}
	return nil
}

func runPreBuildCommands(cfg *document.Config) error {
	if err := runCommands(cfg.PreRun); err != nil {
		return commandError("pre-run", err)
	}
	if err := runCommands(cfg.PreBuild); err != nil {
		return commandError("pre-build", err)
	}
	if cfg.TestOutput != "" {
		if err := runCommands(cfg.PreTest); err != nil {
			return commandError("pre-test", err)
		}
	}
	return nil
}

func runPostBuildCommands(cfg *document.Config) error {
	if cfg.TestOutput != "" {
		if err := runCommands(cfg.PostTest); err != nil {
			return commandError("post-test", err)
		}
	}
	if err := runCommands(cfg.PostBuild); err != nil {
		return commandError("post-build", err)
	}
	if err := runCommands(cfg.PostRun); err != nil {
		return commandError("post-run", err)
	}
	return nil
}
