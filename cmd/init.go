package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/mlange-42/modo/assets"
	"github.com/spf13/cobra"
)

func initCommand() (*cobra.Command, error) {
	root := &cobra.Command{
		Use:   "init",
		Short: "Generate a Modo config file in the current directory",
		Long: `Generate a Modo config file in the current directory.

Complete documentation at https://mlange-42.github.io/modo/`,
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			file := configFile + ".yaml"
			exists, err := fileExists(file)
			if err != nil {
				return fmt.Errorf("error checking config file %s: %s", file, err.Error())
			}
			if exists {
				return fmt.Errorf("config file %s already exists", file)
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

	return root, nil
}
