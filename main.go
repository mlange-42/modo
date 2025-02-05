package main

import (
	"fmt"
	"os"

	"github.com/mlange-42/modo/internal/cmd"
)

func main() {
	root, err := cmd.RootCommand()
	if err != nil {
		panic(err)
	}
	if err := root.Execute(); err != nil {
		fmt.Println("Use 'modo --help' for help.")
		os.Exit(1)
	}
}
