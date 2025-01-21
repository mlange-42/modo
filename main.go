package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mlange-42/modo/document"
	"github.com/spf13/cobra"
)

const configFile = "modo"

func main() {
	start := time.Now()
	root, err := rootCommand()
	if err != nil {
		panic(err)
	}
	if err := root.Execute(); err != nil {
		fmt.Println("Use 'modo --help' for help.")
		os.Exit(1)
	}
	fmt.Printf("Completed in %.1fms 🧯\n", float64(time.Since(start).Microseconds())/1000.0)
}

func rootCommand() (*cobra.Command, error) {
	root := &cobra.Command{
		Use:   "modo",
		Short: "Modo -- DocGen for Mojo.",
		Long: `Modo -- DocGen for Mojo.

Modo generates Markdown for static site generators (SSGs) from 'mojo doc' JSON output.`,
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
	}

	root.CompletionOptions.HiddenDefaultCmd = true

	for _, fn := range []func() (*cobra.Command, error){initCommand, buildCommand, testCommand} {
		cmd, err := fn()
		if err != nil {
			return nil, err
		}
		root.AddCommand(cmd)
	}

	return root, nil
}

func runCommand(command string) error {
	cmd := exec.Command("bash", "-c", command)
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
