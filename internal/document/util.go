package document

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mlange-42/modo/assets"
	"github.com/mlange-42/modo/internal/util"
	"gopkg.in/ini.v1"
)

const codeFence3 = "```"
const codeFence4 = "````"

var initializers = [3]string{
	"__init__", "__moveinit__", "__copyinit__",
}

type fenceType uint8

const (
	fenceNone fenceType = iota
	fenceThree
	fenceFour
)

// GitInfo contains information about a Git repository.
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

// LoadTemplates loads all templates from the assets and additional directories.
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
		exists, isDir, err := util.FileExists(dir)
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

func toFileName(name string) string {
	if caseSensitiveSystem {
		return name
	}
	if isCap(name) {
		return name + capitalFileMarker
	}
	return name
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

// GetGitOrigin tries to determine the `origin` remote repository.
func GetGitOrigin(outDir string) (*GitInfo, error) {
	gitFiles := []string{
		".git/config",
		"../.git/config",
	}

	var content *ini.File
	found := false
	basePath := ""
	for _, f := range gitFiles {
		exists, isDir, err := util.FileExists(f)
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
			basePath, err = util.GetCwdName()
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
				url = value.String()
				ok = true
			}
		}
	}
	if !ok {
		fmt.Printf("WARNING: No Git repository or no remote 'origin' found.\n         Using dummy %s\n", url)
	}
	url, err := getRepoURL(url)
	if err != nil {
		url = "https://github.com/your/package"
		fmt.Printf("WARNING: Git remote 'origin' could not be parsed.\n         Using dummy %s\n", url)
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

func getRepoURL(url string) (string, error) {
	url = strings.TrimSuffix(url, "/")
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return strings.TrimSuffix(url, ".git"), nil
	}
	if !strings.HasPrefix(url, "git@") {
		return "", fmt.Errorf("git remote 'origin' could not be parsed")
	}
	url = strings.TrimPrefix(url, "git@")
	domainRepo := strings.SplitN(url, ":", 2)
	return fmt.Sprintf("https://%s/%s", domainRepo[0], domainRepo[1]), nil
}

func repoToTitleAndPages(repo string) (string, string) {
	repo = strings.TrimSuffix(repo, "/")
	protocolAddress := strings.Split(repo, "//")
	if len(protocolAddress) < 2 {
		return "unknown", "https://example.com"
	}
	parts := strings.Split(protocolAddress[1], "/")
	domainParts := strings.Split(parts[0], ".")
	domain := strings.Join(domainParts[:len(domainParts)-1], ".")

	var title, user, pages string
	switch len(parts) {
	case 1:
		title = "unknown"
		pages = fmt.Sprintf("%s//%s.io/", protocolAddress[0], domain)
	case 2:
		title = parts[1]
		pages = fmt.Sprintf("%s//%s.io/%s/", protocolAddress[0], domain, title)
	default:
		user = parts[1]
		title = parts[2]
		pages = fmt.Sprintf("%s//%s.%s.io/%s/", protocolAddress[0], user, domain, title)
	}
	return title, pages
}
