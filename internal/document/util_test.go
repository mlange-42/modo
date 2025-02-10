package document

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendNew(t *testing.T) {
	sl1 := make([]int, 0, 32)
	sl1 = append(sl1, 1, 2)

	sl2 := appendNew(sl1, 3, 4)

	assert.Equal(t, []int{1, 2}, sl1)
	assert.Equal(t, []int{1, 2, 3, 4}, sl2)
}

func TestWarnOrError(t *testing.T) {
	assert.Nil(t, warnOrError(false, "%s", "test"))
	assert.NotNil(t, warnOrError(true, "%s", "test"))
}

func TestLoadTemplates(t *testing.T) {
	f := TestFormatter{}
	templ, err := LoadTemplates(&f, "https://example.com", "../../docs/docs/templates")
	assert.Nil(t, err)

	assert.NotNil(t, templ.Lookup("package.md"))
}

func TestGetGitOrigin(t *testing.T) {
	oldDir, err := os.Getwd()
	assert.Nil(t, err)

	err = os.Chdir("..")
	assert.Nil(t, err)

	conf, err := GetGitOrigin("docs")

	assert.Nil(t, err)
	assert.Equal(t, conf.Repo, "https://github.com/mlange-42/modo")
	assert.Equal(t, conf.Title, "modo")
	assert.Equal(t, conf.Pages, "https://mlange-42.github.io/modo/")
	assert.Equal(t, conf.GoModule, "github.com/mlange-42/modo/docs")

	err = os.Chdir(oldDir)
	assert.Nil(t, err)
}

func TestGetRepoUrl(t *testing.T) {
	url, err := getRepoURL("https://github.com/mlange-42/modo.git")
	assert.Nil(t, err)
	assert.Equal(t, "https://github.com/mlange-42/modo", url)

	url, err = getRepoURL("https://github.com/mlange-42/modo.git/")
	assert.Nil(t, err)
	assert.Equal(t, "https://github.com/mlange-42/modo", url)

	url, err = getRepoURL("git@github.com:mlange-42/modo")
	assert.Nil(t, err)
	assert.Equal(t, "https://github.com/mlange-42/modo", url)
}

func TestRepoToTitleAndPages(t *testing.T) {
	title, pages := repoToTitleAndPages("https://github.com/user/repo/")
	assert.Equal(t, title, "repo")
	assert.Equal(t, pages, "https://user.github.io/repo/")

	title, pages = repoToTitleAndPages("https://gitlab.com/user/repo")
	assert.Equal(t, title, "repo")
	assert.Equal(t, pages, "https://user.gitlab.io/repo/")

	title, pages = repoToTitleAndPages("https://my-git.com/user/repo")
	assert.Equal(t, title, "repo")
	assert.Equal(t, pages, "https://user.my-git.io/repo/")

	title, pages = repoToTitleAndPages("https://my-git.com/repo")
	assert.Equal(t, title, "repo")
	assert.Equal(t, pages, "https://my-git.io/repo/")

	title, pages = repoToTitleAndPages("https://my-git.com")
	assert.Equal(t, title, "unknown")
	assert.Equal(t, pages, "https://my-git.io/")
}
