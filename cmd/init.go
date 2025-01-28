package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
)

const srcDir = "src"
const docsInDir = "src"
const docsOutDir = "site"
const testsDir = "doctest"

const landingPageContent = `# Landing page

JSON created by mojo doc should be placed next to this file.
`

type config struct {
	Warning      string
	InputFiles   []string
	OutputDir    string
	TestsDir     string
	RenderFormat string
	PreRun       []string
	PostTest     []string
}

type initArgs struct {
	Format        string
	DocsDirectory string
	NoFolders     bool
}

type packageSource struct {
	Name string
	Path []string
}

func initCommand() (*cobra.Command, error) {
	initArgs := initArgs{}

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
			if initArgs.Format == "" {
				initArgs.Format = "plain"
			}
			initArgs.DocsDirectory = strings.ReplaceAll(initArgs.DocsDirectory, "\\", "/")
			return initProject(&initArgs)
		},
	}

	root.Flags().StringVarP(&initArgs.Format, "format", "f", "plain", "Output format. One of (plain|mdbook|hugo)")
	root.Flags().StringVarP(&initArgs.DocsDirectory, "docs", "d", "docs", "Folder for documentation")
	root.Flags().BoolVarP(&initArgs.NoFolders, "no-folders", "F", false, "Don't create any folders")

	return root, nil
}

func initProject(initArgs *initArgs) error {
	_, err := format.GetFormatter(initArgs.Format)
	if err != nil {
		return err
	}

	file := configFile + ".yaml"

	templ := template.New("all")
	templ, err = templ.ParseFS(assets.Config, path.Join("config", file))
	if err != nil {
		return err
	}
	sources, warning, err := findSources(initArgs.Format)
	if err != nil {
		return err
	}
	inDir, outDir, err := createDocs(initArgs, sources)
	if err != nil {
		return err
	}
	preRun, err := createPreRun(initArgs.DocsDirectory, initArgs.Format, sources)
	if err != nil {
		return err
	}
	config := config{
		Warning:      warning,
		InputFiles:   []string{inDir},
		OutputDir:    outDir,
		TestsDir:     testsDir,
		RenderFormat: initArgs.Format,
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

	fmt.Println("Modo project initialized.\nSee file 'modo.yaml' for configuration.")
	return nil
}

func findSources(f string) ([]packageSource, string, error) {
	warning := ""
	sources := []packageSource{}
	srcExists, srcIsDir, err := fileExists(srcDir)
	if err != nil {
		return nil, warning, err
	}

	var allDirs []string
	if srcExists && srcIsDir {
		allDirs = append(allDirs, srcDir)
	} else {
		infos, err := os.ReadDir(".")
		if err != nil {
			return nil, warning, err
		}
		for _, info := range infos {
			if info.IsDir() {
				allDirs = append(allDirs, info.Name())
			}
		}
	}

	nestedSrc := false

	for _, dir := range allDirs {
		isPkg, err := isPackage(dir)
		if err != nil {
			return nil, warning, err
		}
		if isPkg {
			// Package is `<dir>/__init__.mojo`
			file := dir
			if file == srcDir {
				// Package is `src/__init__.mojo`
				file, err = GetCwdName()
				if err != nil {
					return nil, warning, err
				}
			}
			sources = append(sources, packageSource{file, []string{dir}})
			continue
		}
		if dir != srcDir {
			isPkg, err := isPackage(path.Join(dir, srcDir))
			if err != nil {
				return nil, warning, err
			}
			if isPkg {
				// Package is `<dir>/src/__init__.mojo`
				nestedSrc = true
				sources = append(sources, packageSource{dir, []string{dir, srcDir}})
			}
			continue
		}
		infos, err := os.ReadDir(dir)
		if err != nil {
			return nil, warning, err
		}
		for _, info := range infos {
			if info.IsDir() {
				isPkg, err := isPackage(path.Join(dir, info.Name()))
				if err != nil {
					return nil, warning, err
				}
				if isPkg {
					// Package is `src/<dir>/__init__.mojo`
					sources = append(sources, packageSource{info.Name(), []string{dir, info.Name()}})
				}
			}
		}
	}

	if nestedSrc && len(sources) > 1 {
		warning = "WARNING: with folder structure <pkg>/src/__init__.mojo, only a single package is supported"
		fmt.Println(warning)
	}

	if len(sources) == 0 {
		sources = []packageSource{{"mypkg", []string{srcDir, "mypkg"}}}
		warning = fmt.Sprintf("WARNING: no package sources found; using %s", path.Join(sources[0].Path...))
		fmt.Println(warning)
	} else if f == "mdbook" && len(sources) > 1 {
		warning = fmt.Sprintf("WARNING: mdbook format can only use a single package but %d were found; using %s", len(sources), path.Join(sources[0].Path...))
		sources = sources[:1]
		fmt.Println(warning)
	}
	return sources, warning, nil
}

func createDocs(args *initArgs, sources []packageSource) (inDir, outDir string, err error) {
	dir := args.DocsDirectory
	inDir = path.Join(dir, docsInDir)
	outDir = path.Join(dir, docsOutDir)
	if args.Format == "hugo" {
		outDir = path.Join(outDir, "content")
	}
	if args.Format == "mdbook" {
		inDir = path.Join(dir, sources[0].Name+".json")
		outDir = path.Join(dir)
	}

	if args.NoFolders {
		return
	}

	docsExists, _, err := fileExists(dir)
	if err != nil {
		return
	}
	if docsExists {
		fmt.Printf("WARNING: folder %s already exists, skip creating\n", dir)
		return
	}

	if args.Format != "mdbook" {
		if err = mkDirs(inDir); err != nil {
			return
		}
		if err = os.WriteFile(path.Join(inDir, "_index.md"), []byte(landingPageContent), 0644); err != nil {
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

func createPreRun(docsDir, f string, sources []packageSource) (string, error) {
	s := "|\n    echo Running 'mojo test'...\n"

	inDir := docsDir
	if f != "mdbook" {
		inDir = path.Join(docsDir, docsInDir)
	}
	for _, src := range sources {
		s += fmt.Sprintf("    magic run mojo doc -o %s.json %s\n", path.Join(inDir, src.Name), path.Join(src.Path...))
	}

	s += "    echo Done."
	return s, nil
}

func createPostTest(sources []packageSource) string {
	var src string
	if len(sources[0].Path) == 1 {
		src = "."
	} else {
		src = sources[0].Path[0]
	}

	return fmt.Sprintf(`|
    echo Running 'mojo test'...
    magic run mojo test -I %s %s
    echo Done.`, src, testsDir)
}
