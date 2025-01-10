package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mlange-24/modo"
	"github.com/spf13/cobra"
)

func rootCommand() *cobra.Command {
	var file string

	start := &cobra.Command{
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

			docs, err := modo.FromJson(data)
			if err != nil {
				log.Fatal(err)
			}

			err = modo.RenderPackage(&docs.Decl, outDir)
			if err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	start.Flags().StringVarP(&file, "input", "I", "", "File to read. Reads from STDIN if not specified.")

	return start
}

func read(file string) ([]byte, error) {
	if file == "" {
		return io.ReadAll(os.Stdin)
	} else {
		return os.ReadFile(file)
	}
}
