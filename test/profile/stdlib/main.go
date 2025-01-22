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

	"github.com/mlange-42/modo/document"
	"github.com/mlange-42/modo/format"
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
		formatter := format.PlainFormatter{}

		doc, err := document.FromJson(data)
		if err != nil {
			panic(err)
		}
		err = formatter.Render(doc, &config)
		if err != nil {
			panic(err)
		}
	}
}
