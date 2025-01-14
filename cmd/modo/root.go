package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mlange-42/modo"
	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
)

func rootCommand() *cobra.Command {
	var file string
	var renderFormat string
	var caseInsensitive bool

	root := &cobra.Command{
		Use:   "modo OUT-PATH",
		Short: "Mojo documentation generator",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			outDir := args[0]

			if outDir == "" {
				return fmt.Errorf("no output path given")
			}

			data, err := read(file)
			if err != nil {
				return err
			}

			docs, err := document.FromJson(data)
			if err != nil {
				return err
			}

			rFormat, err := getFormat(renderFormat)
			if err != nil {
				return err
			}
			if caseInsensitive {
				document.CaseSensitiveSystem = false
			}

			err = modo.Render(docs, outDir, rFormat)
			if err != nil {
				return err
			}

			return nil
		},
	}

	root.Flags().StringVarP(&file, "input", "i", "", "File to read. Reads from STDIN if not specified.")
	root.Flags().StringVarP(&renderFormat, "format", "f", "plain", "Output format. One of (plain|mdbook|hugo).")
	root.Flags().BoolVar(&caseInsensitive, "case-insensitive", false, "Build for systems that are not case-sensitive regarding file names.\nAppends hyphen (-) to capitalized file names.")

	root.Flags().SortFlags = false

	return root
}

func read(file string) ([]byte, error) {
	if file == "" {
		return io.ReadAll(os.Stdin)
	} else {
		return os.ReadFile(file)
	}
}

func getFormat(f string) (format.Format, error) {
	fm, ok := format.GetFormat(f)
	if !ok {
		return format.Plain, fmt.Errorf("unknown format '%s'. See flag --format", f)
	}
	return fm, nil
}
