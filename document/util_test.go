package document

import (
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
	templ, err := LoadTemplates(&f, "../docs/docs/templates")
	assert.Nil(t, err)

	assert.NotNil(t, templ.Lookup("package.md"))
}

func TestGetGitOrigin(t *testing.T) {
	conf, err := GetGitOrigin("docs")

	assert.Nil(t, err)
	assert.Equal(t, conf.Repo, "https://github.com/mlange-42/modo")
	assert.Equal(t, conf.Title, "modo")
	assert.Equal(t, conf.Pages, "https://mlange-42.github.io/modo/")
	assert.Equal(t, conf.GoModule, "github.com/mlange-42/modo/docs")
}

func TestRepoToTitleAndPages(t *testing.T) {
	title, pages := repoToTitleAndPages("https://github.com/user/repo")
	assert.Equal(t, title, "repo")
	assert.Equal(t, pages, "https://user.github.io/repo/")

	title, pages = repoToTitleAndPages("https://gitlab.com/user/repo")
	assert.Equal(t, title, "repo")
	assert.Equal(t, pages, "https://repo.com")
}
