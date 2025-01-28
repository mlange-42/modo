package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
)

const srcDir = "src"
const docsDir = "docs"
const docsInDir = "src"
const docsOutDir = "site"
const initFile = "__init__.mojo"

type config struct {
	InputFiles   []string
	OutputDir    string
	RenderFormat string
	PreRun       []string
	PostTest     []string
}

func initCommand() (*cobra.Command, error) {
	var format string

	root := &cobra.Command{
		Use:   "init",
		Short: "Generate a Modo config file in the current directory",
		Long: `Generate a Modo config file in the current directory.

Complete documentation at https://mlange-42.github.io/modo/`,
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			file := configFile + ".yaml"
			exists, _, err := fileExists(file)
			if err != nil {
				return fmt.Errorf("error checking config file %s: %s", file, err.Error())
			}
			if exists {
				return fmt.Errorf("config file %s already exists", file)
			}
			if format == "" {
				format = "plain"
			}
			return initProject(format)
		},
	}

	root.Flags().StringVarP(&format, "format", "f", "plain", "Output format. One of (plain|mdbook|hugo)")

	return root, nil
}

func initProject(f string) error {
	_, err := format.GetFormatter(f)
	if err != nil {
		return err
	}

	file := configFile + ".yaml"

	templ := template.New("all")
	templ, err = templ.ParseFS(assets.Config, path.Join("config", file))
	if err != nil {
		return err
	}

	sources := []string{}
	srcExists, srcIsDir, err := fileExists(srcDir)
	if err != nil {
		return err
	}
	if srcExists && srcIsDir {
		infos, err := os.ReadDir(srcDir)
		if err != nil {
			return err
		}
		for _, info := range infos {
			if info.IsDir() {
				pkgFile := path.Join(srcDir, info.Name(), initFile)
				initExists, initIsDir, err := fileExists(pkgFile)
				if err != nil {
					return err
				}
				if initExists && !initIsDir {
					sources = append(sources, info.Name())
				}
			}
		}
	}
	if len(sources) == 0 {
		sources = []string{"mypkg"}
		fmt.Printf("WARNING: no package sources found; using %s\n", path.Join(srcDir, sources[0]))
	} else if f == "mdbook" && len(sources) > 1 {
		sources = sources[:1]
		fmt.Printf("WARNING: mdbook format can only use a single package but %d were found; using %s\n", len(sources), path.Join(srcDir, sources[0]))
	}

	inDir, outDir, err := createDocs(f, sources)
	if err != nil {
		return err
	}

	config := config{
		InputFiles:   []string{inDir},
		OutputDir:    outDir,
		RenderFormat: f,
	}

	b := bytes.Buffer{}
	if err := templ.ExecuteTemplate(&b, "modo.yaml", &config); err != nil {
		return err
	}
	if err := os.WriteFile(file, b.Bytes(), 0644); err != nil {
		return err
	}

	fmt.Println("Modo project initialized.")
	return nil
}

func createDocs(f string, sources []string) (inDir, outDir string, err error) {
	inDir = path.Join(docsDir, docsInDir)
	outDir = path.Join(docsDir, docsOutDir)
	if f == "hugo" {
		outDir = path.Join(outDir, "content")
	}
	if f == "mdbook" {
		inDir = path.Join(docsDir, sources[0]+".json")
		outDir = path.Join(docsDir)
	}

	docsExists, _, err := fileExists(docsDir)
	if err != nil {
		return
	}
	if docsExists {
		fmt.Printf("WARNING: folder %s already exists, skip creating\n", docsDir)
		return
	}

	if f != "mdbook" {
		if err = mkDirs(inDir); err != nil {
			return
		}
		if err = os.WriteFile(path.Join(inDir, "_index.md"), []byte("# Landing page\n"), 0644); err != nil {
			return
		}
	}
	if err = mkDirs(outDir); err != nil {
		return
	}
	return
}
