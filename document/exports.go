package document

import (
	"bufio"
	"strings"
)

const exportsMarker = "Exports:"
const exportsPrefix = "- "

type packageExport struct {
	Short []string
	Long  []string
}

func (proc *Processor) collectExports(p *Package, elems []string) bool {
	anyExports := false

	newElems := appendNew(elems, p.Name)
	for _, pkg := range p.Packages {
		if anyHere := proc.collectExports(pkg, newElems); anyHere {
			anyExports = true
		}
	}

	if proc.UseExports {
		var anyHere bool
		p.exports, p.Description, anyHere = proc.parseExports(p.Description, newElems, true)
		if anyHere {
			anyExports = true
		}
		return anyExports
	}

	p.exports = make([]*packageExport, 0, len(p.Packages)+len(p.Modules))
	for _, pkg := range p.Packages {
		p.exports = append(p.exports, &packageExport{Short: []string{pkg.Name}, Long: appendNew(newElems, pkg.Name)})
	}
	for _, mod := range p.Modules {
		p.exports = append(p.exports, &packageExport{Short: []string{mod.Name}, Long: appendNew(newElems, mod.Name)})
	}

	return anyExports
}

func (proc *Processor) parseExports(pkgDocs string, basePath []string, remove bool) ([]*packageExport, string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(pkgDocs))

	outText := strings.Builder{}
	exports := []*packageExport{}
	anyExports := false
	isExport := false
	exportIndex := 0
	for scanner.Scan() {
		origLine := scanner.Text()
		line := strings.TrimSpace(origLine)
		if isExport {
			if exportIndex == 0 && line == "" {
				continue
			}
			if !strings.HasPrefix(line, exportsPrefix) {
				outText.WriteString(line)
				outText.WriteRune('\n')
				isExport = false
				continue
			}
			short := line[len(exportsPrefix):]
			parts := strings.Split(short, ".")
			exports = append(exports, &packageExport{Short: parts, Long: appendNew(basePath, parts...)})
			anyExports = true
			exportIndex++
		} else {
			if line == exportsMarker {
				isExport = true
				exportIndex = 0
				continue
			}
			outText.WriteString(line)
			outText.WriteRune('\n')
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if remove {
		return exports, outText.String(), anyExports
	}
	return exports, pkgDocs, anyExports
}
