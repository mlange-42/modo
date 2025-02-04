package document

import (
	"bufio"
	"fmt"
	"strings"
)

const exportsMarker = "Exports:"
const exportsPrefix = "- "

type packageExport struct {
	Short   []string
	Renamed []string
	Long    []string
}

// Parses and collects project re-exports, recursively.
func (proc *Processor) collectExports(p *Package, elems []string) (bool, error) {
	proc.renameExports = map[string]string{}
	anyExports := false

	newElems := appendNew(elems, p.Name)
	for _, pkg := range p.Packages {
		anyHere, err := proc.collectExports(pkg, newElems)
		if err != nil {
			return anyExports, err
		}
		if anyHere {
			anyExports = true
		}
	}

	if proc.Config.UseExports {
		var anyHere bool
		var err error
		p.exports, p.Description, anyHere, err = proc.parseExports(p.Description, newElems, true)
		if err != nil {
			return anyExports, err
		}
		if anyHere {
			anyExports = true
		}
		for _, ex := range p.exports {
			longPath := strings.Join(ex.Long, ".")
			if _, ok := proc.allPaths[longPath]; !ok {
				return anyExports, fmt.Errorf("unresolved package re-export '%s' in %s", longPath, strings.Join(newElems, "."))
			}
			proc.renameExports[strings.Join(ex.Short, ".")] = strings.Join(ex.Renamed, ".")
		}
		return anyExports, nil
	}

	p.exports = make([]*packageExport, 0, len(p.Packages)+len(p.Modules))
	for _, pkg := range p.Packages {
		p.exports = append(p.exports, &packageExport{
			Short:   []string{pkg.Name},
			Renamed: []string{pkg.Name},
			Long:    appendNew(newElems, pkg.Name),
		})
	}
	for _, mod := range p.Modules {
		p.exports = append(p.exports, &packageExport{
			Short:   []string{mod.Name},
			Renamed: []string{mod.Name},
			Long:    appendNew(newElems, mod.Name),
		})
	}

	return anyExports, nil
}

func (proc *Processor) parseExports(pkgDocs string, basePath []string, remove bool) ([]*packageExport, string, bool, error) {
	scanner := bufio.NewScanner(strings.NewReader(pkgDocs))

	outText := strings.Builder{}
	exports := []*packageExport{}
	anyExports := false
	isExport := false
	fenced3 := false
	fenced4 := false

	exportIndex := 0
	for scanner.Scan() {
		origLine := scanner.Text()
		line := strings.TrimSpace(origLine)

		fenced := false
		if strings.HasPrefix(origLine, codeFence3) {
			fenced3 = !fenced3
			fenced = true
		}
		if strings.HasPrefix(origLine, codeFence4) {
			fenced4 = !fenced4
			fenced = true
		}
		if fenced || fenced3 || fenced4 {
			isExport = false
			outText.WriteString(origLine)
			outText.WriteRune('\n')
			continue
		}

		if isExport {
			if exportIndex == 0 && line == "" {
				continue
			}
			if !strings.HasPrefix(line, exportsPrefix) {
				outText.WriteString(origLine)
				outText.WriteRune('\n')
				isExport = false
				continue
			}
			exportsAs := strings.Split(line[len(exportsPrefix):], " ")
			short := exportsAs[0]
			partsShort := strings.Split(short, ".")
			renamed := partsShort[len(partsShort)-1]
			if len(exportsAs) == 3 && exportsAs[1] == "as" {
				renamed = exportsAs[2]
			} else if len(exportsAs) != 1 {
				if err := proc.warnOrError("invalid syntax in package re-export '%s' in %s", line[len(exportsPrefix):], strings.Join(basePath, ".")); err != nil {
					return nil, "", false, err
				}
			}
			exports = append(exports, &packageExport{
				Short:   partsShort,
				Renamed: appendNew(partsShort[:len(partsShort)-1], renamed),
				Long:    appendNew(basePath, partsShort...)})
			anyExports = true
			exportIndex++
		} else {
			if line == exportsMarker {
				isExport = true
				exportIndex = 0
				continue
			}
			outText.WriteString(origLine)
			outText.WriteRune('\n')
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if remove {
		return exports, strings.TrimSuffix(outText.String(), "\n"), anyExports, nil
	}
	return exports, pkgDocs, anyExports, nil
}
