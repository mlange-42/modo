package document

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"gopkg.in/ini.v1"
)

const codeFence3 = "```"
const codeFence4 = "````"

type fenceType uint8

const (
	fenceNone fenceType = iota
	fenceThree
	fenceFour
)

type GitInfo struct {
	Title    string
	Repo     string
	Pages    string
	GoModule string
	BasePath string
}

func getFenceType(line string) fenceType {
	isFence4 := strings.HasPrefix(line, codeFence4)
	if strings.HasPrefix(line, codeFence3) && !isFence4 {
		return fenceThree
	}
	if isFence4 {
		return fenceFour
	}
	return fenceNone
}

// appends to a slice, but guaranties to return a new one and not alter the original.
func appendNew[T any](sl []T, elems ...T) []T {
	sl2 := make([]T, len(sl), len(sl)+len(elems))
	copy(sl2, sl)
	sl2 = append(sl2, elems...)
	return sl2
}

func warnOrError(strict bool, pattern string, args ...any) error {
	if strict {
		return fmt.Errorf(pattern, args...)
	}
	fmt.Printf("WARNING: "+pattern+"\n", args...)
	return nil
}

func LoadTemplates(f Formatter, sourceURL string, additional ...string) (*template.Template, error) {
	templ := template.New("all")
	templ = templ.Funcs(template.FuncMap{
		"toLink":    f.ToLinkPath,
		"sourceUrl": func() string { return sourceURL },
	})
	templ, err := templ.ParseFS(assets.Templates, "templates/*.*", "templates/**/*.*")
	if err != nil {
		return nil, err
	}

	for _, dir := range additional {
		if dir == "" {
			continue
		}
		exists, isDir, err := fileExists(dir)
		if err != nil {
			return nil, err
		}
		if !exists || !isDir {
			return nil, fmt.Errorf("template directory '%s' does not exist", dir)
		}
		moreTemplates, err := findTemplates(dir)
		if err != nil {
			return nil, err
		}
		templ, err = templ.ParseFiles(moreTemplates...)
		if err != nil {
			return nil, err
		}
	}
	return templ, nil
}

func fileExists(file string) (exists, isDir bool, err error) {
	var s os.FileInfo
	if s, err = os.Stat(file); err == nil {
		exists = true
		isDir = s.IsDir()
		return
	} else if !errors.Is(err, os.ErrNotExist) {
		return
	}
	err = nil
	return
}

func findTemplates(dir string) ([]string, error) {
	allTemplates := []string{}
	err := filepath.WalkDir(dir,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				allTemplates = append(allTemplates, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return allTemplates, nil
}

func GetCwdName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return cwd, err
	}
	return filepath.Base(cwd), nil
}

func GetGitOrigin(outDir string) (*GitInfo, error) {
	gitFiles := []string{
		".git/config",
		"../.git/config",
	}

	var content *ini.File
	found := false
	basePath := ""
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

		if strings.HasPrefix(f, "..") {
			basePath, err = GetCwdName()
			if err != nil {
				return nil, err
			}
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

	return &GitInfo{
		Title:    title,
		Repo:     url,
		Pages:    pages,
		GoModule: module,
		BasePath: basePath,
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
