package document

import (
	"bufio"
	"fmt"
	"path"
	"strings"
)

const globalSuffix = "-global"
const setupSuffix = "-setup"
const teardownSuffix = "-teardown"
const docTestAttr = "doctest"

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
		filePath += "_" + test.Name + ".mojo"
		fullPath := path.Join(dir, filePath)

		err = proc.WriteFile(fullPath, b.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (proc *Processor) extractTests(text string, elems []string, modElems int) (string, error) {
	_ = modElems
	scanner := bufio.NewScanner(strings.NewReader(text))
	outText := strings.Builder{}

	fenced := false
	blocks := map[string]*docTest{}
	var blockLines []string
	var blockName string
	var excluded bool
	for scanner.Scan() {
		origLine := scanner.Text()

		isStart := false
		isFence := strings.HasPrefix(origLine, codeFence3)
		if isFence && !fenced {
			var ok bool
			var err error
			blockName, ok, err = parseBlockAttr(origLine)
			if err != nil {
				return "", fmt.Errorf("%s in %s", err.Error(), strings.Join(elems, "."))
			}
			if !ok {
				blockName = ""
			}
			if strings.HasSuffix(blockName, globalSuffix) ||
				strings.HasSuffix(blockName, setupSuffix) ||
				strings.HasSuffix(blockName, teardownSuffix) {
				excluded = true
			}
			fenced = true
			isStart = true
		}

		if !excluded {
			outText.WriteString(origLine)
			outText.WriteRune('\n')
		}

		if fenced && !isFence && blockName != "" {
			blockLines = append(blockLines, origLine)
		}

		if isFence && fenced && !isStart {
			if blockName == "" {
				excluded = false
				fenced = false
				continue
			}
			if dt, ok := blocks[blockName]; ok {
				dt.Code = append(dt.Code, blockLines...)
			} else {
				blocks[blockName] = &docTest{
					Name: blockName,
					Path: elems,
					Code: append([]string{}, blockLines...)}
			}
			blockLines = blockLines[:0]
			excluded = false
			fenced = false
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if fenced {
		return "", fmt.Errorf("unbalanced code block in %s", strings.Join(elems, "."))
	}

	for name, block := range blocks {
		if strings.HasSuffix(name, globalSuffix) ||
			strings.HasSuffix(name, setupSuffix) ||
			strings.HasSuffix(name, teardownSuffix) {
			continue
		}
		if global, ok := blocks[name+globalSuffix]; ok {
			block.Global = global.Code
		}
		if setup, ok := blocks[name+setupSuffix]; ok {
			block.Code = append(setup.Code, block.Code...)
		}
		if teardown, ok := blocks[name+teardownSuffix]; ok {
			block.Code = append(block.Code, teardown.Code...)
		}
		proc.docTests = append(proc.docTests, block)
	}

	return outText.String(), nil
}

func parseBlockAttr(line string) (string, bool, error) {
	parts := strings.SplitN(line, "{", 2)
	if len(parts) < 2 {
		return "", false, nil
	}
	attrString := strings.TrimSpace(parts[1])
	if !strings.HasSuffix(attrString, "}") {
		return "", false, fmt.Errorf("missing closing parentheses in code block attributes")
	}
	attrString = strings.TrimSuffix(attrString, "}")
	attrPairs := strings.Split(attrString, " ")

	for _, pair := range attrPairs {
		elems := strings.Split(pair, "=")
		if len(elems) != 2 {
			return "", false, fmt.Errorf("malformed code block attributes '%s'", pair)
		}
		if strings.TrimSpace(elems[0]) != docTestAttr {
			continue
		}
		return strings.Trim(elems[1], "\""), true, nil
	}
	return "", false, nil
}
