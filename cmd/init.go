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
const testsDir = "doctest"
const initFile = "__init__.mojo"

type config struct {
	InputFiles   []string
	OutputDir    string
	TestsDir     string
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
	sources, err := findSources(f)
	if err != nil {
		return err
	}
	inDir, outDir, err := createDocs(f, sources)
	if err != nil {
		return err
	}
	preRun, err := createPreRun(f, sources)
	if err != nil {
		return err
	}
	config := config{
		InputFiles:   []string{inDir},
		OutputDir:    outDir,
		TestsDir:     testsDir,
		RenderFormat: f,
		PreRun:       []string{preRun},
		PostTest:     []string{createPostTest(sources)},
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

func findSources(f string) ([]string, error) {
	sources := []string{}
	srcExists, srcIsDir, err := fileExists(srcDir)
	if err != nil {
		return nil, err
	}

	if srcExists && srcIsDir {
		pkgFile := path.Join(srcDir, initFile)
		initExists, initIsDir, err := fileExists(pkgFile)
		if err != nil {
			return nil, err
		}
		if initExists && !initIsDir {
			sources = append(sources, "")
		} else {
			infos, err := os.ReadDir(srcDir)
			if err != nil {
				return nil, err
			}
			for _, info := range infos {
				if info.IsDir() {
					pkgFile := path.Join(srcDir, info.Name(), initFile)
					initExists, initIsDir, err := fileExists(pkgFile)
					if err != nil {
						return nil, err
					}
					if initExists && !initIsDir {
						sources = append(sources, info.Name())
					}
				}
			}
		}
	}
	if len(sources) == 0 {
		sources = []string{"mypkg"}
		fmt.Printf("WARNING: no package sources found; using %s\n", path.Join(srcDir, sources[0]))
	} else if f == "mdbook" && len(sources) > 1 {
		fmt.Printf("WARNING: mdbook format can only use a single package but %d were found; using %s\n", len(sources), path.Join(srcDir, sources[0]))
		sources = sources[:1]
	}
	return sources, nil
}

func createDocs(f string, sources []string) (inDir, outDir string, err error) {
	inDir = path.Join(docsDir, docsInDir)
	outDir = path.Join(docsDir, docsOutDir)
	if f == "hugo" {
		outDir = path.Join(outDir, "content")
	}
	if f == "mdbook" {
		file := sources[0]
		if file == "" {
			file, err = GetCwdName()
			if err != nil {
				return
			}
		}
		inDir = path.Join(docsDir, file+".json")
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
	if err = mkDirs(testsDir); err != nil {
		return
	}
	return
}

func createPreRun(f string, sources []string) (string, error) {
	s := "|\n    echo Running 'mojo test'...\n"

	inDir := docsDir
	if f != "mdbook" {
		inDir = path.Join(docsDir, docsInDir)
	}
	for _, src := range sources {
		file := src
		if file == "" {
			var err error
			file, err = GetCwdName()
			if err != nil {
				return "", err
			}
		}
		s += fmt.Sprintf("    magic run mojo doc -o %s.json %s\n", path.Join(inDir, file), path.Join(srcDir, src))
	}

	s += "    echo Done."
	return s, nil
}

func createPostTest(sources []string) string {
	src := srcDir
	if len(sources) == 1 && sources[0] == "" {
		src = "."
	}

	return fmt.Sprintf(`|
    echo Running 'mojo test'...
    magic run mojo test -I %s %s
    echo Done.`, src, testsDir)
}
