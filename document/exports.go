package document

import (
	"bufio"
	"strings"
)

const exportsMarker = "Exports:"
const exportsPrefix = "- "

func (proc *Processor) collectExports(p *Package) {
	for _, pkg := range p.Packages {
		proc.collectExports(pkg)
	}
	p.Exports = proc.parseExports(p.Description)
}

func (proc *Processor) parseExports(pkgDocs string) []string {
	scanner := bufio.NewScanner(strings.NewReader(pkgDocs))

	exports := []string{}
	isExport := false
	exportIndex := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if isExport {
			if exportIndex == 0 && line == "" {
				continue
			}
			if !strings.HasPrefix(line, exportsPrefix) {
				isExport = false
				continue
			}
			exports = append(exports, line[len(exportsPrefix):])
			exportIndex++
		} else {
			if line == exportsMarker {
				isExport = true
				exportIndex = 0
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return exports
}
