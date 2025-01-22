package document

import (
	"bufio"
	"fmt"
	"path"
	"strings"
)

const docTestAttr = "doctest"
const hideAttr = "hide"
const globalAttr = "global"

func (proc *Processor) extractDocTests() error {
	proc.docTests = []*docTest{}
	return proc.walkDocs(proc.Docs, proc.extractTests, func(elem Named) string {
		return elem.GetFileName()
	})
}

func (proc *Processor) writeDocTests(dir string) error {
	if dir == "" {
		return nil
	}
	err := proc.mkDirs(dir)
	if err != nil {
		return err
	}
	for _, test := range proc.docTests {
		b := strings.Builder{}
		err := proc.Template.ExecuteTemplate(&b, "doctest.mojo", test)
		if err != nil {
			return err
		}
		filePath := strings.Join(test.Path, "_")
		filePath += "_" + test.Name + "_test.mojo"
		fullPath := path.Join(dir, filePath)

		err = proc.WriteFile(fullPath, b.String())
		if err != nil {
			return err
		}
	}
	fmt.Printf("Extracted %d tests.\n", len(proc.docTests))
	return nil
}

func (proc *Processor) extractTests(text string, elems []string, modElems int) (string, error) {
	_ = modElems
	scanner := bufio.NewScanner(strings.NewReader(text))
	outText := strings.Builder{}

	fenced := false
	blocks := map[string]*docTest{}
	var blockLines []string
	var globalLines []string
	var blockName string
	var excluded bool
	var global bool
	var count int
	for scanner.Scan() {
		origLine := scanner.Text()

		isStart := false
		isFence := strings.HasPrefix(origLine, codeFence3)
		if isFence && !fenced {
			var ok bool
			var err error
			blockName, excluded, global, ok, err = parseBlockAttr(origLine)
			if err != nil {
				if err := proc.warnOrError("%s in %s", err.Error(), strings.Join(elems, ".")); err != nil {
					return "", err
				}
			}
			if !ok {
				blockName = ""
			}
			fenced = true
			isStart = true
		}

		if !excluded {
			outText.WriteString(origLine)
			outText.WriteRune('\n')
		}

		if fenced && !isFence && blockName != "" {
			if global {
				globalLines = append(globalLines, origLine)
			} else {
				blockLines = append(blockLines, origLine)
			}
		}
		count++

		if isFence && fenced && !isStart {
			if blockName == "" {
				excluded = false
				global = false
				fenced = false
				continue
			}
			if dt, ok := blocks[blockName]; ok {
				dt.Code = append(dt.Code, blockLines...)
				dt.Global = append(dt.Global, globalLines...)
			} else {
				blocks[blockName] = &docTest{
					Name:   blockName,
					Path:   elems,
					Code:   append([]string{}, blockLines...),
					Global: append([]string{}, globalLines...),
				}
			}
			blockLines = blockLines[:0]
			globalLines = globalLines[:0]
			excluded = false
			global = false
			fenced = false
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if fenced {
		if err := proc.warnOrError("unbalanced code block in %s", strings.Join(elems, ".")); err != nil {
			return "", err
		}
	}

	for _, block := range blocks {
		proc.docTests = append(proc.docTests, block)
	}

	return strings.TrimSuffix(outText.String(), "\n"), nil
}

func parseBlockAttr(line string) (name string, hide bool, global bool, ok bool, err error) {
	parts := strings.SplitN(line, "{", 2)
	if len(parts) < 2 {
		return
	}
	attrString := strings.TrimSpace(parts[1])
	if !strings.HasSuffix(attrString, "}") {
		err = fmt.Errorf("missing closing parentheses in code block attributes")
		return
	}
	attrString = strings.TrimSuffix(attrString, "}")

	quoted := false
	attrPairs := strings.FieldsFunc(attrString, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == ' '
	})

	for _, pair := range attrPairs {
		elems := strings.Split(pair, "=")
		if len(elems) > 2 {
			err = fmt.Errorf("malformed code block attributes '%s'", pair)
			return
		}
		if len(elems) < 2 {
			continue
		}

		key := strings.TrimSpace(elems[0])
		if key == docTestAttr {
			name = strings.Trim(elems[1], "\"")
			continue
		}
		if key == hideAttr {
			h := strings.Trim(elems[1], "\" ")
			if h == "true" {
				hide = true
			}
			continue
		}
		if key == globalAttr {
			g := strings.Trim(elems[1], "\"")
			if g == "true" {
				global = true
			}
			continue
		}
	}
	ok = true
	return
}
