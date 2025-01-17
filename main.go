package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
)

func main() {
	start := time.Now()
	if err := rootCommand().Execute(); err != nil {
		fmt.Println("Use 'modo --help' for help.")
		os.Exit(1)
	}
	fmt.Printf("Completed in %.1fms 🧯\n", float64(time.Since(start).Microseconds())/1000.0)
}

func rootCommand() *cobra.Command {
	var cliArgs document.Config
	var renderFormat string

	root := &cobra.Command{
		Use:   "modo OUT-PATH",
		Short: "Modo -- DocGen for Mojo.",
		Long: `Modo -- DocGen for Mojo.

Modo generates Markdown for static site generators (SSGs) from 'mojo doc' JSON output.`,
		Example: `  modo docs -i docs.json        # from a file
  mojo doc ./src | modo docs    # from 'mojo doc'`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliArgs.OutputDir = args[0]
			return run(&cliArgs, renderFormat)
		},
	}

	root.Flags().StringVarP(&cliArgs.InputFile, "input", "i", "", "'mojo doc' JSON file to process. Reads from STDIN if not specified.")
	root.Flags().StringVarP(&renderFormat, "format", "f", "plain", "Output format. One of (plain|mdbook|hugo).")
	root.Flags().BoolVarP(&cliArgs.UseExports, "exports", "e", false, "Process according to 'Exports:' sections in packages.")
	root.Flags().BoolVar(&cliArgs.ShortLinks, "short-links", false, "Render shortened link labels, stripping packages and modules.")
	root.Flags().BoolVar(&cliArgs.CaseInsensitive, "case-insensitive", false, "Build for systems that are not case-sensitive regarding file names.\nAppends hyphen (-) to capitalized file names.")
	root.Flags().BoolVar(&cliArgs.Strict, "strict", false, "Strict mode. Errors instead of warnings.")
	root.Flags().BoolVar(&cliArgs.DryRun, "dry-run", false, "Dry-run without any file output.")
	root.Flags().StringSliceVarP(&cliArgs.TemplateDirs, "templates", "t", []string{}, "Optional directories with templates for (partial) overwrite.\nSee folder assets/templates in the repository.")

	root.Flags().SortFlags = false
	root.MarkFlagFilename("input", "json")
	root.MarkFlagDirname("templates")

	return root
}

func run(args *document.Config, renderFormat string) error {
	if args.OutputDir == "" {
		return fmt.Errorf("no output path given")
	}

	docs, err := readDocs(args.InputFile)
	if err != nil {
		return err
	}

	rFormat, err := format.GetFormat(renderFormat)
	if err != nil {
		return err
	}
	formatter := format.GetFormatter(rFormat)
	err = formatter.Render(docs, args)
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
