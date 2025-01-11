package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mlange-24/modo"
	"github.com/mlange-24/modo/doc"
	"github.com/mlange-24/modo/format"
	"github.com/spf13/cobra"
)

type renderFormats struct {
	mdBook bool
}

func rootCommand() *cobra.Command {
	var file string
	var formats renderFormats

	root := &cobra.Command{
		Use:   "modo <OUT>",
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

			docs, err := doc.FromJson(data)
			if err != nil {
				log.Fatal(err)
			}

			rFormat := getFormat(&formats)

			err = modo.RenderPackage(&docs.Decl, outDir, rFormat, true)
			if err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	root.Flags().StringVarP(&file, "input", "I", "", "File to read. Reads from STDIN if not specified.")
	root.Flags().BoolVar(&formats.mdBook, "mdbook", false, "Write in mdBook format.")

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
