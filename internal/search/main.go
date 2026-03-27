// A helper tool to find newly added JSON fields.
// Usage example:
//
//	go run ./internal/search stdlib.json <new-field> kind
//
// Will list all element kinds in which <new-field> appears.
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "usage: search <file> <field> <neighbor>\n")
		os.Exit(1)
	}

	file := os.Args[1]
	target := os.Args[2]
	neighbor := os.Args[3]

	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var root any
	if err := json.Unmarshal(data, &root); err != nil {
		panic(err)
	}

	results := map[string]struct{}{}
	walk(root, target, neighbor, results)

	for v := range results {
		fmt.Println(v)
	}
}

func walk(v any, target, neighbor string, out map[string]struct{}) {
	switch x := v.(type) {

	case map[string]any:
		// If this object contains the target field, check neighbor
		if _, ok := x[target]; ok {
			if n, ok := x[neighbor]; ok {
				switch nv := n.(type) {
				case string:
					out[nv] = struct{}{}
				default:
					// stringify non-string values
					b, _ := json.Marshal(n)
					out[string(b)] = struct{}{}
				}
			}
		}

		// Recurse into all values
		for _, v2 := range x {
			walk(v2, target, neighbor, out)
		}

	case []any:
		for _, elem := range x {
			walk(elem, target, neighbor, out)
		}
	}
}
