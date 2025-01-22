package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mlange-42/modo/cmd"
)

func main() {
	start := time.Now()
	root, err := cmd.RootCommand()
	if err != nil {
		panic(err)
	}
	if err := root.Execute(); err != nil {
		fmt.Println("Use 'modo --help' for help.")
		os.Exit(1)
	}
	fmt.Printf("Completed in %.1fms 🧯\n", float64(time.Since(start).Microseconds())/1000.0)
}
