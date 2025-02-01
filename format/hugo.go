package format

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/document"
	"gopkg.in/ini.v1"
)

const landingPageContentHugo = `---
title: Landing page
type: docs
---

JSON created by mojo doc should be placed next to this file.

Additional documentation files go here, too.
They will be processed for doc-tests and copied to folder 'site/content'.
`

type Hugo struct{}

type hugoConfig struct {
	Title  string
	Repo   string
	Module string
	Pages  string
}

func (f *Hugo) Accepts(files []string) error {
	return nil
}

func (f *Hugo) ProcessMarkdown(element any, text string, proc *document.Processor) (string, error) {
	b := strings.Builder{}
	err := proc.Template.ExecuteTemplate(&b, "hugo_front_matter.md", element)
	if err != nil {
		return "", err
	}
	b.WriteRune('\n')
	b.WriteString(text)
	return b.String(), nil
}

func (f *Hugo) WriteAuxiliary(p *document.Package, dir string, proc *document.Processor) error {
	return nil
}

func (f *Hugo) ToFilePath(p string, kind string) string {
	if kind == "package" || kind == "module" {
		return path.Join(p, "_index.md")
	}
	return p + ".md"
}

func (f *Hugo) ToLinkPath(p string, kind string) string {
	p = f.ToFilePath(p, kind)
	return fmt.Sprintf("{{< ref \"%s\" >}}", p)
}

func (f *Hugo) Input(in string, sources []document.PackageSource) string {
	return in
}

func (f *Hugo) Output(out string) string {
	return path.Join(out, "content")
}

func (f *Hugo) GitIgnore(in, out string, sources []document.PackageSource) []string {
	return []string{
		"# files generated by 'mojo doc'",
		fmt.Sprintf("/%s/*.json", in),
		"# files generated by Modo",
		fmt.Sprintf("/%s/%s/", out, "content"),
		"# files generated by Hugo",
		fmt.Sprintf("/%s/%s/", out, "public"),
		fmt.Sprintf("/%s/*.lock", out),
		"# test file generated by Modo",
		"/test/",
	}
}

func (f *Hugo) CreateDirs(base, in, out string, sources []document.PackageSource, templ *template.Template) error {
	inDir, outDir := path.Join(base, in), path.Join(base, f.Output(out))
	testDir := path.Join(base, "test")
	if err := mkDirs(inDir); err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(inDir, "_index.md"), []byte(landingPageContentHugo), 0644); err != nil {
		return err
	}
	if err := mkDirs(outDir); err != nil {
		return err
	}
	if err := mkDirs(testDir); err != nil {
		return err
	}
	return f.createInitialFiles(base, path.Join(base, out), templ)
}

func (f *Hugo) createInitialFiles(docDir, hugoDir string, templ *template.Template) error {
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

func (f *Hugo) Clean(out, tests string) error {
	if err := emptyDir(out); err != nil {
		return err
	}
	return emptyDir(tests)
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
	title, pages := repoToTitleAndPages(url)
	module := strings.ReplaceAll(strings.ReplaceAll(url, "https://", ""), "http://", "")
	module = fmt.Sprintf("%s/%s", module, outDir)

	return &hugoConfig{
		Title:  title,
		Repo:   url,
		Pages:  pages,
		Module: module,
	}, nil
}

func repoToTitleAndPages(repo string) (string, string) {
	if !strings.HasPrefix(repo, "https://github.com/") {
		parts := strings.Split(repo, "/")
		title := parts[len(parts)-1]
		return title, fmt.Sprintf("https://%s.com", title)
	}
	repo = strings.TrimPrefix(repo, "https://github.com/")
	parts := strings.Split(repo, "/")
	return parts[1], fmt.Sprintf("https://%s.github.io/%s/", parts[0], parts[1])
}
