package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"github.com/mlange-42/modo/format"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

const srcDir = "src"
const docsInDir = "src"
const docsOutDir = "site"
const testsDir = "test"
const gitignoreFile = ".gitignore"

const landingPageContent = `# Landing page

JSON created by mojo doc should be placed next to this file.

Additional documentation files go here, too.
They will be processed for doc-tests and copied to folder 'site'.
`

const landingPageContentHugo = `---
title: Landing page
type: docs
---

JSON created by mojo doc should be placed next to this file.

Additional documentation files go here, too.
They will be processed for doc-tests and copied to folder 'site/content'.
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

type hugoConfig struct {
	Title  string
	Repo   string
	Module string
	Pages  string
}

type mdBookConfig struct {
	Title string
}

type packageSource struct {
	Name string
	Path []string
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
	_, err := format.GetFormatter(initArgs.Format)
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
	inDir, outDir, err := createDocs(initArgs, templ, sources)
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

func createDocs(args *initArgs, templ *template.Template, sources []packageSource) (inDir, outDir string, err error) {
	var gitignore []string

	dir := args.DocsDirectory
	inDir = path.Join(dir, docsInDir)
	outDir = path.Join(dir, docsOutDir)
	testOurDir := path.Join(dir, testsDir)
	if args.Format == "hugo" {
		outDir = path.Join(outDir, "content")
		gitignore = append(gitignore,
			"# files generated by 'mojo doc'",
			fmt.Sprintf("/%s/*.json", docsInDir),
			"# files generated by Modo",
			fmt.Sprintf("/%s/%s/", docsOutDir, "content"),
			"# files generated by Hugo",
			fmt.Sprintf("/%s/%s/", docsOutDir, "public"),
			fmt.Sprintf("/%s/*.lock", docsOutDir))
	} else if args.Format == "mdbook" {
		file := sources[0].Name + ".json"
		inDir = path.Join(dir, file)
		outDir = dir
		gitignore = append(gitignore,
			"# file generated by 'mojo doc'",
			fmt.Sprintf("/%s", file),
			"# files generated by Modo",
			fmt.Sprintf("/%s/", sources[0].Name),
			"# files generated by MdBook",
			fmt.Sprintf("/%s/", "public"))
	} else {
		gitignore = append(gitignore,
			"# files generated by 'mojo doc'",
			fmt.Sprintf("/%s/*.json", docsInDir),
			"# files generated by Modo",
			fmt.Sprintf("/%s/", docsOutDir),
		)
	}

	if args.NoFolders {
		return
	}

	gitignore = append(gitignore,
		"# test file generated by Modo",
		fmt.Sprintf("/%s/", testsDir))

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
		var content string
		if args.Format == "hugo" {
			content = landingPageContentHugo
		} else {
			content = landingPageContent
		}
		if err = os.WriteFile(path.Join(inDir, "_index.md"), []byte(content), 0644); err != nil {
			return
		}
	}
	if err = mkDirs(outDir); err != nil {
		return
	}
	if err = mkDirs(testOurDir); err != nil {
		return
	}
	if err = writeGitIgnore(dir, gitignore); err != nil {
		return
	}
	if args.Format == "hugo" {
		if err = createHugoFiles(dir, path.Join(dir, docsOutDir), templ); err != nil {
			return
		}
	} else if args.Format == "mdbook" {
		if err = createMdBookFiles(sources[0].Name, outDir, templ); err != nil {
			return
		}
	}
	return
}

func createHugoFiles(docDir, hugoDir string, templ *template.Template) error {
	config, err := getGitOrigin(docDir)
	if err != nil {
		return err
	}

	files := [][]string{
		{"hugo.yaml", "hugo.yaml"},
		{"hugo.mod", "go.mod"},
		{"hugo.sum", "go.sum"},
	}
	for _, f := range files {
		outFile := path.Join(hugoDir, f[1])
		exists, _, err := fileExists(outFile)
		if err != nil {
			return err
		}
		if exists {
			fmt.Printf("WARNING: Hugo file %s already exists, skip creating\n", outFile)
			return nil
		}
	}

	for _, f := range files {
		outFile := path.Join(hugoDir, f[1])
		b := bytes.Buffer{}
		if err := templ.ExecuteTemplate(&b, f[0], &config); err != nil {
			return err
		}
		if err := os.WriteFile(outFile, b.Bytes(), 0644); err != nil {
			return err
		}
	}
	return nil
}

func createMdBookFiles(title, docDir string, templ *template.Template) error {
	outFile := path.Join(docDir, "book.toml")
	exists, _, err := fileExists(outFile)
	if err != nil {
		return err
	}
	if exists {
		fmt.Printf("WARNING: MdBook config file %s already exists, skip creating\n", outFile)
		return nil
	}

	config := mdBookConfig{Title: title}

	b := bytes.Buffer{}
	if err := templ.ExecuteTemplate(&b, "book.toml", &config); err != nil {
		return err
	}
	if err := os.WriteFile(outFile, b.Bytes(), 0644); err != nil {
		return err
	}

	cssDir := path.Join(docDir, "css")
	cssFile := path.Join(cssDir, "mdbook.css")
	exists, _, err = fileExists(cssFile)
	if err != nil {
		return err
	}
	if exists {
		fmt.Printf("WARNING: MdBook CSS file %s already exists, skip creating\n", cssFile)
		return nil
	}

	if err := os.MkdirAll(cssDir, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	css, err := fs.ReadFile(assets.CSS, "css/mdbook.css")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(cssDir, "custom.css"), css, 0644); err != nil {
		return err
	}
	return nil
}

func writeGitIgnore(dir string, gitignore []string) error {
	s := strings.Join(gitignore, "\n") + "\n"
	return os.WriteFile(path.Join(dir, gitignoreFile), []byte(s), 0644)
}

func createPreRun(docsDir, f string, sources []packageSource) (string, error) {
	s := "|\n    echo Running 'mojo doc'...\n"

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

func createPostTest(docsDir string, sources []packageSource) string {
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

func getGitOrigin(outDir string) (*hugoConfig, error) {
	gitFiles := []string{
		".git/config",
		"../.git/config",
	}

	var content *ini.File
	found := false
	for _, f := range gitFiles {
		exists, isDir, err := fileExists(f)
		if err != nil {
			return nil, err
		}
		if !exists || isDir {
			continue
		}
		content, err = ini.Load(f)
		if err != nil {
			return nil, err
		}
		found = true
		break
	}

	url := "https://github.com/your/package"
	ok := false
	if found {
		section := content.Section(`remote "origin"`)
		if section != nil {
			value := section.Key("url")
			if value != nil {
				url = strings.TrimSuffix(value.String(), ".git")
				ok = true
			}
		}
	}
	if !ok {
		fmt.Printf("WARNING: No Git repository or no remote 'origin' found.\n         Using dummy %s\n", url)
	}
	title, pages := repoToTitleEndPages(url)
	module := strings.ReplaceAll(strings.ReplaceAll(url, "https://", ""), "http://", "")
	module = fmt.Sprintf("%s/%s", module, outDir)

	return &hugoConfig{
		Title:  title,
		Repo:   url,
		Pages:  pages,
		Module: module,
	}, nil
}

func repoToTitleEndPages(repo string) (string, string) {
	if !strings.HasPrefix(repo, "https://github.com/") {
		parts := strings.Split(repo, "/")
		title := parts[len(parts)-1]
		return title, fmt.Sprintf("https://%s.com", title)
	}
	repo = strings.TrimPrefix(repo, "https://github.com/")
	parts := strings.Split(repo, "/")
	return parts[1], fmt.Sprintf("https://%s.github.io/%s/", parts[0], parts[1])
}
