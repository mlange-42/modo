package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mlange-42/modo"
	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
)

type renderFormats struct {
	mdBook bool
}

func rootCommand() *cobra.Command {
	var file string
	var formats renderFormats
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
				log.Fatal(err)
			}

			rFormat := getFormat(&formats)
			if caseInsensitive {
				document.CaseSensitiveSystem = false
			}

			err = modo.RenderDocs(docs, outDir, rFormat)
			if err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	root.Flags().StringVarP(&file, "input", "i", "", "File to read. Reads from STDIN if not specified.")
	root.Flags().BoolVar(&formats.mdBook, "mdbook", false, "Writes in mdBook format.")
	root.Flags().BoolVar(&caseInsensitive, "case-insensitive", false, "Build for systems that are not case-sensitive regarding file names.\nAppends hyphen (-) to capitalized file names.")

	root.MarkFlagsMutuallyExclusive("mdbook")

	return root
}

func read(file string) ([]byte, error) {
	if file == "" {
		return io.ReadAll(os.Stdin)
	} else {
		return os.ReadFile(file)
	}
}

func getFormat(f *renderFormats) format.Format {
	if f.mdBook {
		return format.MdBook
	}
	return format.Plain
}
