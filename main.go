package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

type args struct {
	file            string
	renderFormat    string
	caseInsensitive bool
	outDir          string
}

func rootCommand() *cobra.Command {
	var cliArgs args

	root := &cobra.Command{
		Use:   "modo OUT-PATH",
		Short: "Mojo documentation generator",
		Long:  ``,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs.outDir = args[0]
			return run(&cliArgs)
		},
	}

	root.Flags().StringVarP(&cliArgs.file, "input", "i", "", "File to read. Reads from STDIN if not specified.")
	root.Flags().StringVarP(&cliArgs.renderFormat, "format", "f", "plain", "Output format. One of (plain|mdbook|hugo).")
	root.Flags().BoolVar(&cliArgs.caseInsensitive, "case-insensitive", false, "Build for systems that are not case-sensitive regarding file names.\nAppends hyphen (-) to capitalized file names.")

	root.Flags().SortFlags = false

	return root
}

func run(args *args) error {
	if args.outDir == "" {
		return fmt.Errorf("no output path given")
	}

	data, err := read(args.file)
	if err != nil {
		return err
	}

	docs, err := document.FromJson(data)
	if err != nil {
		return err
	}

	rFormat, err := format.GetFormat(args.renderFormat)
	if err != nil {
		return err
	}
	if args.caseInsensitive {
		document.CaseSensitiveSystem = false
	}

	err = format.Render(docs, args.outDir, rFormat)
	if err != nil {
		return err
	}

	return nil
}

func read(file string) ([]byte, error) {
	if file == "" {
		return io.ReadAll(os.Stdin)
	} else {
		return os.ReadFile(file)
	}
}
