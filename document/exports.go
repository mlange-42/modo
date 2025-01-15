package document

import (
	"bufio"
	"strings"
)

const exportsMarker = "Exports:"
const exportsPrefix = "- "

type PackageExport struct {
	Short []string
	Long  []string
}

func (proc *Processor) collectExports(p *Package, elems []string) {
	newElems := appendNew(elems, p.Name)
	for _, pkg := range p.Packages {
		proc.collectExports(pkg, newElems)
	}

	if proc.UseExports {
		p.Exports = proc.parseExports(p.Description, newElems)
		return
	}

	p.Exports = make([]*PackageExport, 0, len(p.Packages)+len(p.Modules))
	for _, pkg := range p.Packages {
		p.Exports = append(p.Exports, &PackageExport{Short: []string{pkg.Name}, Long: appendNew(newElems, pkg.Name)})
	}
	for _, mod := range p.Modules {
		p.Exports = append(p.Exports, &PackageExport{Short: []string{mod.Name}, Long: appendNew(newElems, mod.Name)})
	}

}

func (proc *Processor) parseExports(pkgDocs string, basePath []string) []*PackageExport {
	scanner := bufio.NewScanner(strings.NewReader(pkgDocs))

	exports := []*PackageExport{}
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
			short := line[len(exportsPrefix):]
			parts := strings.Split(short, ".")
			exports = append(exports, &PackageExport{Short: parts, Long: appendNew(basePath, parts...)})
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
