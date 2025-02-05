package main

//go:generate go run .
//go:generate go fmt ../version.go

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("Updating version...")

	cmd := exec.Command("git", "describe", "--tags")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	version := strings.TrimSpace(string(out))

	versionFile := "../version.go"
	content := fmt.Sprintf("package cmd\n\nconst version = \"%s\"\n", version)
	if err := os.WriteFile(versionFile, []byte(content), 0666); err != nil {
		panic(err)
	}
}
