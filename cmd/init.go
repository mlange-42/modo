package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
)

const srcDir = "src"
const docsInDir = "src"
const docsOutDir = "site"
const testsDir = "test"
const gitignoreFile = ".gitignore"

type config struct {
	Warning      string
	InputFiles   []string
	Sources      []string
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

func initCommand() (*cobra.Command, error) {
	initArgs := initArgs{}

	root := &cobra.Command{
		Use:   "init FORMAT",
		Short: "Set up a Modo project in the current directory",
		Long: `Set up a Modo project in the current directory.

The format argument is required and must be one of (plain|mdbook|hugo).
Complete documentation at https://mlange-42.github.io/modo/`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			initArgs.Format = args[0]

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

	root.Flags().StringVarP(&initArgs.DocsDirectory, "docs", "d", "docs", "Folder for documentation")
	root.Flags().BoolVarP(&initArgs.NoFolders, "no-folders", "F", false, "Don't create any folders")

	root.Flags().SortFlags = false
	root.MarkFlagDirname("docs")

	return root, nil
}

func initProject(initArgs *initArgs) error {
	form, err := format.GetFormatter(initArgs.Format)
	if err != nil {
		return err
	}

	file := configFile + ".yaml"

	templ := template.New("all")
	templ, err = templ.ParseFS(assets.Config, "**/*")
	if err != nil {
		return err
	}
	sources, warning, err := findSources(initArgs.Format)
	if err != nil {
		return err
	}
	inDir, outDir, err := createDocs(initArgs, form, templ, sources)
	if err != nil {
		return err
	}
	preRun, err := createPreRun(initArgs.DocsDirectory, sources, form)
	if err != nil {
		return err
	}

	sourceDirs := make([]string, 0, len(sources))
	for _, s := range sources {
		sourceDirs = append(sourceDirs, path.Join(s.Path...))
	}
	config := config{
		Warning:      warning,
		InputFiles:   []string{inDir},
		Sources:      sourceDirs,
		OutputDir:    outDir,
		TestsDir:     path.Join(initArgs.DocsDirectory, testsDir),
		RenderFormat: initArgs.Format,
		PreRun:       []string{preRun},
		PostTest:     []string{createPostTest(initArgs.DocsDirectory, sources)},
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

func findSources(f string) ([]document.PackageSource, string, error) {
	warning := ""
	sources := []document.PackageSource{}
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
			sources = append(sources, document.PackageSource{Name: file, Path: []string{dir}})
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
				sources = append(sources, document.PackageSource{Name: dir, Path: []string{dir, srcDir}})
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
					sources = append(sources, document.PackageSource{Name: info.Name(), Path: []string{dir, info.Name()}})
				}
			}
		}
	}

	if nestedSrc && len(sources) > 1 {
		warning = "WARNING: with folder structure <pkg>/src/__init__.mojo, only a single package is supported"
		fmt.Println(warning)
	}

	if len(sources) == 0 {
		sources = []document.PackageSource{{Name: "mypkg", Path: []string{srcDir, "mypkg"}}}
		warning = fmt.Sprintf("WARNING: no package sources found; using %s", path.Join(sources[0].Path...))
		fmt.Println(warning)
	} else if f == "mdbook" && len(sources) > 1 {
		warning = fmt.Sprintf("WARNING: mdbook format can only use a single package but %d were found; using %s", len(sources), path.Join(sources[0].Path...))
		sources = sources[:1]
		fmt.Println(warning)
	}
	return sources, warning, nil
}

func createDocs(args *initArgs, form document.Formatter, templ *template.Template, sources []document.PackageSource) (inDir, outDir string, err error) {
	dir := args.DocsDirectory
	inDir, outDir = form.Input(docsInDir, sources), form.Output(docsOutDir)
	inDir, outDir = path.Join(dir, inDir), path.Join(dir, outDir)

	if args.NoFolders {
		return
	}

	docsExists, _, err := fileExists(dir)
	if err != nil {
		return
	}
	if docsExists {
		err = fmt.Errorf("documentation folder '%s' already exists.\n"+
			"Use flag --docs to use a different folder, --no-folders to skip folders setup", dir)
		return
	}

	if err = form.CreateDirs(dir, docsInDir, docsOutDir, sources, templ); err != nil {
		return
	}

	gitignore := form.GitIgnore(docsInDir, docsOutDir, sources)
	if err = writeGitIgnore(dir, gitignore); err != nil {
		return
	}
	return
}

func writeGitIgnore(dir string, gitignore []string) error {
	s := strings.Join(gitignore, "\n") + "\n"
	return os.WriteFile(path.Join(dir, gitignoreFile), []byte(s), 0644)
}

func createPreRun(docsDir string, sources []document.PackageSource, form document.Formatter) (string, error) {
	s := "|\n    echo Running 'mojo doc'...\n"

	inDir := path.Join(docsDir, form.Input(docsInDir, sources))
	for _, src := range sources {
		outFile := inDir
		if !strings.HasSuffix(outFile, ".json") {
			outFile = path.Join(inDir, src.Name+".json")
		}
		s += fmt.Sprintf("    magic run mojo doc -o %s %s\n", outFile, path.Join(src.Path...))
	}

	s += "    echo Done."
	return s, nil
}

func createPostTest(docsDir string, sources []document.PackageSource) string {
	testOurDir := path.Join(docsDir, testsDir)
	var src string
	if len(sources[0].Path) == 1 {
		src = "."
	} else {
		src = sources[0].Path[0]
	}

	return fmt.Sprintf(`|
    echo Running 'mojo test'...
    magic run mojo test -I %s %s
    echo Done.`, src, testOurDir)
}
