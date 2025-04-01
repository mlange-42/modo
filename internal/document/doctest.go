package document

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const docTestAttr = "doctest"
const hideAttr = "hide"
const globalAttr = "global"

func (proc *Processor) extractDocTests() error {
	proc.docTests = []*docTest{}
	w := walker{
		Func:     proc.extractTests,
		NameFunc: func(elem Named) string { return elem.GetFileName() },
	}
	return w.walkAllDocStrings(proc.Docs)
}

func (proc *Processor) extractDocTestsMarkdown(baseDir string, build bool) error {
	proc.docTests = []*docTest{}
	outDir := filepath.Clean(proc.Config.OutputDir)
	baseDir = filepath.Clean(baseDir)
	err := filepath.WalkDir(baseDir,
		func(p string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			return proc.extractMarkdown(p, baseDir, outDir, build)
		})
	if err != nil {
		return err
	}
	if proc.Config.TestOutput != "" {
		fmt.Printf("Extracted %d test(s) from Markdown files.\n", len(proc.docTests))
		err = proc.writeDocTests(proc.Config.TestOutput)
		if err != nil {
			return err
		}
	}
	return nil
}

func (proc *Processor) extractMarkdown(file, baseDir, outDir string, build bool) error {
	if strings.HasSuffix(strings.ToLower(file), ".json") {
		return nil
	}

	cleanPath := path.Clean(file)
	relPath := filepath.Clean(strings.TrimPrefix(cleanPath, baseDir))
	targetPath := filepath.Join(outDir, relPath)
	targetDir, _ := filepath.Split(targetPath)

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return err
	}
	contentStr := string(content)
	if strings.HasSuffix(strings.ToLower(file), ".md") {
		var err error
		contentStr, err = proc.extractTests(contentStr, []string{strings.TrimSuffix(relPath, ".md")}, 1)
		if err != nil {
			return err
		}
	}

	if build {
		err = proc.mkDirs(targetDir)
		if err != nil {
			return err
		}
		return proc.writeFile(targetPath, contentStr)
	}
	return nil
}

func (proc *Processor) writeDocTests(dir string) error {
	if dir == "" {
		return nil
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

		parentDir, _ := filepath.Split(filepath.Clean(fullPath))
		err = proc.mkDirs(parentDir)
		if err != nil {
			return err
		}

		err = proc.writeFile(fullPath, b.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func (proc *Processor) extractTests(text string, elems []string, modElems int) (string, error) {
	t, tests, err := extractTestsText(text, elems, proc.Config.Strict)
	if err != nil {
		return "", err
	}
	proc.docTests = append(proc.docTests, tests...)
	return t, nil
}

func extractTestsText(text string, elems []string, strict bool) (string, []*docTest, error) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	outText := strings.Builder{}

	fenced := fenceNone
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
		currFence := getFenceType(origLine)
		if currFence != fenceNone && fenced == fenceNone {
			var ok bool
			var err error
			blockName, excluded, global, ok, err = parseBlockAttr(origLine)
			if err != nil {
				if err := warnOrError(strict, "%s in %s", err.Error(), strings.Join(elems, ".")); err != nil {
					return "", nil, err
				}
			}
			if !ok {
				blockName = ""
			}
			fenced = currFence
			isStart = true
		}

		if !excluded {
			outText.WriteString(origLine)
			outText.WriteRune('\n')
		}

		if fenced != fenceNone && currFence != fenced && blockName != "" {
			if global {
				globalLines = append(globalLines, origLine)
			} else {
				blockLines = append(blockLines, origLine)
			}
		}
		count++

		if fenced != fenceNone && currFence == fenced && !isStart {
			if blockName == "" {
				excluded = false
				global = false
				fenced = fenceNone
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
			fenced = fenceNone
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if fenced != fenceNone {
		if err := warnOrError(strict, "unbalanced code block in %s", strings.Join(elems, ".")); err != nil {
			return "", nil, err
		}
	}

	tests := make([]*docTest, 0, len(blocks))
	for _, block := range blocks {
		tests = append(tests, block)
	}
	sort.Slice(tests, func(i, j int) bool { return tests[i].Name < tests[j].Name })

	return strings.TrimSuffix(outText.String(), "\n"), tests, nil
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
			} else if h != "false" {
				err = fmt.Errorf("invalid argument in code block attribute 'hide': '%s'", h)
				return
			}
			continue
		}
		if key == globalAttr {
			g := strings.Trim(elems[1], "\"")
			if g == "true" {
				global = true
			} else if g != "false" {
				err = fmt.Errorf("invalid argument in code block attribute 'global': '%s'", g)
				return
			}
			continue
		}
	}
	ok = true
	return
}
