package main

// Profiling:
// go build ./benchmark/profile/stdlib
// ./stdlib stdlib.json
// go tool pprof -http=":8000" -nodefraction=0.001 ./stdlib cpu.pprof
// go tool pprof -http=":8000" -nodefraction=0.001 ./stdlib mem.pprof

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/mlange-42/modo/internal/document"
	"github.com/mlange-42/modo/internal/format"
	"github.com/pkg/profile"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("no JSON file given")
	}
	iters := 10

	stop := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
	run(os.Args[1], iters)
	stop.Stop()

	//stop = profile.Start(profile.MemProfileAllocs, profile.ProfilePath("."))
	//run(os.Args[1], iters)
	//stop.Stop()
}

func run(file string, iters int) {
	t := testing.T{}
	tmpDir := t.TempDir()
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	for range iters {
		config := document.Config{
			InputFiles:      []string{file},
			OutputDir:       tmpDir,
			UseExports:      false,
			ShortLinks:      true,
			CaseInsensitive: true,
			TestOutput:      path.Join(tmpDir, "doctest"),
			DryRun:          true,
		}
		formatter := format.Plain{}

		doc, err := document.FromJSON(data)
		if err != nil {
			panic(err)
		}

		err = document.Render(doc, &config, &formatter, "")
		if err != nil {
			panic(err)
		}
	}
}
