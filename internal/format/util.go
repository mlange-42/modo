package format

import (
	"os"

	"github.com/mlange-42/modo/internal/util"
)

func emptyDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return util.MkDirs(dir)
}
