package util

import (
	"errors"
	"os"
	"path/filepath"
)

func MkDirs(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func GetCwdName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return cwd, err
	}
	return filepath.Base(cwd), nil
}

func FileExists(file string) (exists, isDir bool, err error) {
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
