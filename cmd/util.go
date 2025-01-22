package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/mlange-42/modo/document"
	"github.com/spf13/viper"
)

const configFile = "modo"
const setExitOnError = "set -e"

func runCommand(command string) error {
	commandWithExit := fmt.Sprintf("%s\n%s", setExitOnError, command)
	cmd := exec.Command("bash", "-c", commandWithExit)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommands(commands []string) error {
	for _, command := range commands {
		err := runCommand(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func readDocs(file string) (*document.Docs, error) {
	data, err := read(file)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
		return document.FromYaml(data)
	}

	return document.FromJson(data)
}

func read(file string) ([]byte, error) {
	if file == "" {
		return io.ReadAll(os.Stdin)
	} else {
		return os.ReadFile(file)
	}
}

func fileExists(file string) (bool, error) {
	if _, err := os.Stat(file); err == nil {
		return true, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}
	return false, nil
}

func mountProject(v *viper.Viper, paths []string) error {
	withConfig := len(paths) > 0
	p := "."
	if withConfig {
		p = paths[0]
		if err := os.Chdir(p); err != nil {
			return err
		}
	}

	filePath := configFile + ".yaml"
	exists, err := fileExists(filePath)
	if err != nil {
		return err
	}
	if !exists {
		if withConfig {
			return fmt.Errorf("no '%s' found in path '%s'", filePath, p)
		}
		return nil
	}

	v.SetConfigName(configFile)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		_, notFound := err.(viper.ConfigFileNotFoundError)
		if !notFound {
			return err
		}
		if withConfig {
			return err
		}
	}
	return nil
}
