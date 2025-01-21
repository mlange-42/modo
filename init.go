package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/mlange-42/modo/assets"
	"github.com/spf13/cobra"
)

func initCommand() *cobra.Command {
	root := &cobra.Command{
		Use:          "init",
		Short:        "Generate Modo config file in the current directory.",
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			file := configFile + ".yaml"
			if _, err := os.Stat(file); err == nil {
				return fmt.Errorf("config file %s already exists", file)
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("error checking config file %s: %s", file, err.Error())
			}

			configData, err := fs.ReadFile(assets.Config, path.Join("config", file))
			if err != nil {
				return err
			}
			if err := os.WriteFile(file, configData, 0644); err != nil {
				return err
			}

			fmt.Println("Modo project initialized.")
			return nil
		},
	}

	return root
}
