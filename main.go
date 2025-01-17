package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCommand().Execute(); err != nil {
		fmt.Println("Use 'modo --help' for help.")
		os.Exit(1)
	}
}

type args struct {
	file            string
	renderFormat    string
	caseInsensitive bool
	useExports      bool
	shortLinks      bool
	outDir          string
	templateDirs    []string
}

func rootCommand() *cobra.Command {
	var cliArgs args

	root := &cobra.Command{
		Use:   "modo OUT-PATH",
		Short: "Modo -- DocGen for Mojo.",
		Long: `Modo -- DocGen for Mojo.

Modo generates Markdown for static site generators (SSGs) from 'mojo doc' JSON output.

Usage:
  modo docs -i docs.json        # from a file
  mojo doc ./src | modo docs    # from 'mojo doc'
`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs.outDir = args[0]
			return run(&cliArgs)
		},
	}

	root.Flags().StringVarP(&cliArgs.file, "input", "i", "", "'mojo doc' JSON file to process. Reads from STDIN if not specified.")
	root.Flags().StringVarP(&cliArgs.renderFormat, "format", "f", "plain", "Output format. One of (plain|mdbook|hugo).")
	root.Flags().BoolVarP(&cliArgs.useExports, "exports", "e", false, "Process according to 'Exports:' sections in packages.")
	root.Flags().BoolVar(&cliArgs.shortLinks, "short-links", false, "Render shortened link labels, stripping packages and modules.")
	root.Flags().BoolVar(&cliArgs.caseInsensitive, "case-insensitive", false, "Build for systems that are not case-sensitive regarding file names.\nAppends hyphen (-) to capitalized file names.")
	root.Flags().StringSliceVarP(&cliArgs.templateDirs, "templates", "t", []string{}, "Optional directories with templates for (partial) overwrite.\nSee folder assets/templates in the repository.")

	root.Flags().SortFlags = false
	root.MarkFlagFilename("input", "json")
	root.MarkFlagDirname("templates")

	return root
}

func run(args *args) error {
	if args.outDir == "" {
		return fmt.Errorf("no output path given")
	}

	docs, err := readDocs(args.file)
	if err != nil {
		return err
	}

	rFormat, err := format.GetFormat(args.renderFormat)
	if err != nil {
		return err
	}
	formatter := format.GetFormatter(rFormat)

	err = document.Render(docs, &document.Config{
		OutputDir:     args.outDir,
		TemplateDirs:  args.templateDirs,
		RenderFormat:  formatter,
		UseExports:    args.useExports,
		ShortLinks:    args.shortLinks,
		CaseSensitive: !args.caseInsensitive,
	})
	if err != nil {
		return err
	}

	return nil
}

func readDocs(file string) (*document.Docs, error) {
	data, err := read(file)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
		return document.FromYaml(data)
	}

	return document.FromJson(data)
}

func read(file string) ([]byte, error) {
	if file == "" {
		return io.ReadAll(os.Stdin)
	} else {
		return os.ReadFile(file)
	}
}
