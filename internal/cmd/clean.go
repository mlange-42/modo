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

func cleanCommand(_ chan struct{}) (*cobra.Command, error) {
	v := viper.New()
	var config string

	var cwd string

	root := &cobra.Command{
		Use:   "clean [PATH]",
		Short: "Remove Markdown and test files generated by Modo",
		Long: `Remove Markdown and test files generated by Modo.

Complete documentation at https://mlange-42.github.io/modo/`,
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
			if err := runClean(cliArgs); err != nil {
				return err
			}

			fmt.Printf("Completed in %.1fms 🧯\n", float64(time.Since(start).Microseconds())/1000.0)
			return nil
		},
	}

	root.Flags().StringVarP(&config, "config", "c", defaultConfigFile, "Config file in the working directory to use")

	root.Flags().SortFlags = false
	root.MarkFlagFilename("config", "yaml")

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

	return formatter.Clean(args.OutputDir, args.TestOutput)
}
