package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mlange-24/modo"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Error reading stdin:", err)
	}

	docs, err := modo.FromJson(data)
	if err != nil {
		log.Fatal(err)
	}
	_ = docs

	fmt.Println("StdIn read...")

	err = modo.RenderPackage(&docs.Decl, "out")
	if err != nil {
		log.Fatal(err)
	}
}
