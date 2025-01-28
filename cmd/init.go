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

type packageSource struct {
	Name string
	Path string
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

func findSources(f string) ([]packageSource, error) {
	sources := []packageSource{}
	srcExists, srcIsDir, err := fileExists(srcDir)
	if err != nil {
		return nil, err
	}

	var allDirs []string
	if srcExists && srcIsDir {
		allDirs = append(allDirs, srcDir)
	} else {
		infos, err := os.ReadDir(".")
		if err != nil {
			return nil, err
		}
		for _, info := range infos {
			if info.IsDir() {
				allDirs = append(allDirs, info.Name())
			}
		}
	}

	nestedSrc := false

	for _, dir := range allDirs {
		pkgFile := path.Join(dir, initFile)
		initExists, initIsDir, err := fileExists(pkgFile)
		if err != nil {
			return nil, err
		}
		if initExists && !initIsDir {
			// Package is `<dir>/__init__.mojo`
			file, err := GetCwdName()
			if err != nil {
				return nil, err
			}
			sources = append(sources, packageSource{file, dir})
			continue
		}
		if dir != srcDir {
			pkgFile := path.Join(dir, srcDir, initFile)
			initExists, initIsDir, err := fileExists(pkgFile)
			if err != nil {
				return nil, err
			}
			if initExists && !initIsDir {
				// Package is `<dir>/src/__init__.mojo`
				nestedSrc = true
				sources = append(sources, packageSource{dir, path.Join(dir, srcDir)})
			}
			continue
		}
		infos, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, info := range infos {
			if info.IsDir() {
				pkgFile := path.Join(dir, info.Name(), initFile)
				initExists, initIsDir, err := fileExists(pkgFile)
				if err != nil {
					return nil, err
				}
				if initExists && !initIsDir {
					// Package is `src/<dir>/__init__.mojo`
					sources = append(sources, packageSource{info.Name(), path.Join(dir, info.Name())})
				}
			}
		}
	}

	if nestedSrc && len(sources) > 1 {
		fmt.Println("WARNING: with folder structure <pkg>/src/__init__.mojo, only a single package is supported")
	}

	if len(sources) == 0 {
		sources = []packageSource{{"mypkg", path.Join(srcDir, "mypkg")}}
		fmt.Printf("WARNING: no package sources found; using %s\n", sources[0].Path)
	} else if f == "mdbook" && len(sources) > 1 {
		fmt.Printf("WARNING: mdbook format can only use a single package but %d were found; using %s\n", len(sources), sources[0].Path)
		sources = sources[:1]
	}
	return sources, nil
}

func createDocs(f string, sources []packageSource) (inDir, outDir string, err error) {
	inDir = path.Join(docsDir, docsInDir)
	outDir = path.Join(docsDir, docsOutDir)
	if f == "hugo" {
		outDir = path.Join(outDir, "content")
	}
	if f == "mdbook" {
		inDir = path.Join(docsDir, sources[0].Name+".json")
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

func createPreRun(f string, sources []packageSource) (string, error) {
	s := "|\n    echo Running 'mojo test'...\n"

	inDir := docsDir
	if f != "mdbook" {
		inDir = path.Join(docsDir, docsInDir)
	}
	for _, src := range sources {
		s += fmt.Sprintf("    magic run mojo doc -o %s.json %s\n", path.Join(inDir, src.Name), src.Path)
	}

	s += "    echo Done."
	return s, nil
}

func createPostTest(sources []packageSource) string {
	src := srcDir
	if len(sources) == 1 && sources[0].Path == srcDir {
		src = "."
	}

	return fmt.Sprintf(`|
    echo Running 'mojo test'...
    magic run mojo test -I %s %s
    echo Done.`, src, testsDir)
}
