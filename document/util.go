package document

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"
	"text/template"
)

const regexString = `(?s)(?:(` + "```.*?```)|(`.*?`" + `))|(\[.*?\])`

type elemPath struct {
	Elements []string
	Kind     string
}

func findLinks(text string) ([]int, error) {
	re, err := regexp.Compile(regexString)
	if err != nil {
		return nil, err
	}
	links := []int{}
	results := re.FindAllStringSubmatchIndex(text, -1)
	for _, r := range results {
		if r[6] >= 0 {
			if len(text) > r[7] && string(text[r[7]]) == "(" {
				continue
			}
			links = append(links, r[6], r[7])
		}
	}

	return links, nil
}

func ProcessLinks(doc *Docs, t *template.Template) error {
	lookup := collectPaths(doc)
	return processLinksPackage(doc.Decl, []string{}, lookup, t)
}

func processLinksPackage(p *Package, elems []string, lookup map[string]elemPath, t *template.Template) error {
	newElems := appendNew(elems, p.GetName())

	var err error
	p.Summary, err = replaceLinks(p.Summary, newElems, lookup, t)
	if err != nil {
		return err
	}
	p.Description, err = replaceLinks(p.Summary, newElems, lookup, t)
	if err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		processLinksPackage(pkg, newElems, lookup, t)
	}

	return nil
}

func collectPaths(doc *Docs) map[string]elemPath {
	out := map[string]elemPath{}
	collectPathsPackage(doc.Decl, []string{}, []string{}, out)
	return out
}

func collectPathsPackage(p *Package, elems []string, pathElem []string, out map[string]elemPath) {
	newElems := appendNew(elems, p.GetName())
	newPath := appendNew(pathElem, p.GetFileName())
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "package"}

	for _, pkg := range p.Packages {
		collectPathsPackage(pkg, newElems, newPath, out)
	}
	for _, mod := range p.Modules {
		collectPathsModule(mod, newElems, newPath, out)
	}
}

func collectPathsModule(m *Module, elems []string, pathElem []string, out map[string]elemPath) {
	newElems := appendNew(elems, m.GetName())
	newPath := appendNew(pathElem, m.GetFileName())
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "module"}

	for _, s := range m.Structs {
		collectPathsStruct(s, newElems, newPath, out)
	}
	for _, t := range m.Traits {
		collectPathsTrait(t, newElems, newPath, out)
	}
	for _, f := range m.Functions {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, f.GetFileName())
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member"}
	}
}

func collectPathsStruct(s *Struct, elems []string, pathElem []string, out map[string]elemPath) {
	newElems := appendNew(elems, s.GetName())
	newPath := appendNew(pathElem, s.GetFileName())
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member"}

	for _, f := range s.Parameters {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#Parameters")
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member"}
	}
	for _, f := range s.Fields {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#Fields")
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member"}
	}
	for _, f := range s.Functions {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#"+f.GetName())
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member"}
	}
}

func collectPathsTrait(t *Trait, elems []string, pathElem []string, out map[string]elemPath) {
	newElems := appendNew(elems, t.GetName())
	newPath := appendNew(pathElem, t.GetFileName())
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member"}

	for _, f := range t.Fields {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#Fields")
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member"}
	}
	for _, f := range t.Functions {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#"+f.GetName())
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member"}
	}
}

func replaceLinks(text string, elems []string, lookup map[string]elemPath, t *template.Template) (string, error) {
	indices, err := findLinks(text)
	if err != nil {
		return "", err
	}
	if len(indices) == 0 {
		return text, nil
	}
	for i := len(indices) - 2; i >= 0; i -= 2 {
		start, end := indices[i], indices[i+1]
		link := text[start+1 : end-1]

		dots := 0
		pathPrefix := []string{}
		for strings.HasPrefix(link[dots:], ".") {
			dots++
			pathPrefix = append(pathPrefix, "..")
		}
		if dots > len(elems) {
			log.Printf("Too many leading dots in cross ref: %s", link)
			continue
		}

		var fullLink string
		linkText := link[dots:]
		subElems := elems[:len(elems)-dots]
		if len(subElems) == 0 {
			fullLink = link[dots:]
		} else {
			fullLink = strings.Join(subElems, ".") + "." + linkText
		}
		elemPath, ok := lookup[fullLink]
		if !ok {
			log.Printf("Can't resolve cross ref: %s (%s)", link, fullLink)
			continue
		}

		pathPrefixStr := path.Join(pathPrefix...)
		pathStr := strings.Builder{}
		if strings.HasPrefix(elemPath.Elements[len(elemPath.Elements)-1], "#") {
			fullPath := path.Join(pathPrefixStr, path.Join(elemPath.Elements[len(subElems):len(elemPath.Elements)-1]...))
			err := t.ExecuteTemplate(&pathStr, elemPath.Kind+"_path.md", fullPath)
			if err != nil {
				return "", err
			}
			pathStr.WriteString(elemPath.Elements[len(elemPath.Elements)-1])
		} else {
			fullPath := path.Join(pathPrefixStr, path.Join(elemPath.Elements[len(subElems):]...))
			err := t.ExecuteTemplate(&pathStr, elemPath.Kind+"_path.md", fullPath)
			if err != nil {
				return "", err
			}
		}
		text = fmt.Sprintf("%s[%s](%s)%s", text[:start], linkText, pathStr.String(), text[end:])
	}
	return text, nil
}

func appendNew[T any](sl []T, elems ...T) []T {
	sl2 := make([]T, len(sl), len(sl)+len(elems))
	copy(sl2, sl)
	sl2 = append(sl2, elems...)
	return sl2
}

func cleanup(doc *Docs) {
	cleanupPackage(doc.Decl)
}

func cleanupPackage(p *Package) {
	for _, pp := range p.Packages {
		cleanupPackage(pp)
	}
	newModules := make([]*Module, 0, len(p.Modules))
	for _, m := range p.Modules {
		cleanupModule(m)
		if m.GetName() != "__init__" {
			newModules = append(newModules, m)
		}
	}
	p.Modules = newModules
}

func cleanupModule(m *Module) {
	for _, s := range m.Structs {
		if s.Signature == "" {
			s.Signature = createSignature(s)
		}
	}
}

func createSignature(s *Struct) string {
	b := strings.Builder{}
	b.WriteString("struct ")
	b.WriteString(s.GetName())

	if len(s.Parameters) == 0 {
		return b.String()
	}

	b.WriteString("[")

	prevKind := ""
	for i, par := range s.Parameters {
		written := false
		if par.PassingKind == "kw" && prevKind != par.PassingKind {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString("*")
			written = true
		}
		if prevKind == "inferred" && par.PassingKind != prevKind {
			b.WriteString(", //")
			written = true
		}
		if prevKind == "pos" && par.PassingKind != prevKind {
			b.WriteString(", /")
			written = true
		}

		if i > 0 || written {
			b.WriteString(", ")
		}

		b.WriteString(fmt.Sprintf("%s: %s", par.GetName(), par.Type))
		if len(par.Default) > 0 {
			b.WriteString(fmt.Sprintf(" = %s", par.Default))
		}

		prevKind = par.PassingKind
	}
	if prevKind == "inferred" {
		b.WriteString(", //")
	}
	if prevKind == "pos" {
		b.WriteString(", /")
	}

	b.WriteString("]")

	return b.String()
}
