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
	Elements  []string
	Kind      string
	IsSection bool
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
	//for k, v := range lookup {
	//	fmt.Println(k, v.Elements)
	//}
	return processLinksPackage(doc.Decl, []string{}, lookup, t)
}

func processLinksPackage(p *Package, elems []string, lookup map[string]elemPath, t *template.Template) error {
	newElems := appendNew(elems, p.GetName())

	var err error
	p.Summary, err = replaceLinks(p.Summary, newElems, len(newElems), lookup, t)
	if err != nil {
		return err
	}
	p.Description, err = replaceLinks(p.Description, newElems, len(newElems), lookup, t)
	if err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		processLinksPackage(pkg, newElems, lookup, t)
	}
	for _, mod := range p.Modules {
		processLinksModule(mod, newElems, lookup, t)
	}

	return nil
}

func processLinksModule(m *Module, elems []string, lookup map[string]elemPath, t *template.Template) error {
	newElems := appendNew(elems, m.GetName())

	var err error
	m.Summary, err = replaceLinks(m.Summary, newElems, len(newElems), lookup, t)
	if err != nil {
		return err
	}
	m.Description, err = replaceLinks(m.Description, newElems, len(newElems), lookup, t)
	if err != nil {
		return err
	}

	for _, f := range m.Functions {
		err := processLinksFunction(f, newElems, lookup, t)
		if err != nil {
			return err
		}
	}
	for _, s := range m.Structs {
		err := processLinksStruct(s, newElems, lookup, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func processLinksStruct(s *Struct, elems []string, lookup map[string]elemPath, t *template.Template) error {
	newElems := appendNew(elems, s.GetName())

	var err error
	s.Summary, err = replaceLinks(s.Summary, newElems, len(elems), lookup, t)
	if err != nil {
		return err
	}
	s.Description, err = replaceLinks(s.Description, newElems, len(elems), lookup, t)
	if err != nil {
		return err
	}

	return nil
}

func processLinksFunction(f *Function, elems []string, lookup map[string]elemPath, t *template.Template) error {
	newElems := appendNew(elems, f.GetName())

	var err error
	f.Summary, err = replaceLinks(f.Summary, newElems, len(elems), lookup, t)
	if err != nil {
		return err
	}
	f.Description, err = replaceLinks(f.Description, newElems, len(elems), lookup, t)
	if err != nil {
		return err
	}

	for _, o := range f.Overloads {
		err := processLinksFunction(o, elems, lookup, t)
		if err != nil {
			return err
		}
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
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "package", IsSection: false}

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
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "module", IsSection: false}

	for _, s := range m.Structs {
		collectPathsStruct(s, newElems, newPath, out)
	}
	for _, t := range m.Traits {
		collectPathsTrait(t, newElems, newPath, out)
	}
	for _, f := range m.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, f.GetFileName())
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: false}
	}
}

func collectPathsStruct(s *Struct, elems []string, pathElem []string, out map[string]elemPath) {
	newElems := appendNew(elems, s.GetName())
	newPath := appendNew(pathElem, s.GetFileName())
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: false}

	for _, f := range s.Parameters {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#Parameters")
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: true}
	}
	for _, f := range s.Fields {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#Fields")
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: true}
	}
	for _, f := range s.Functions {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#"+f.GetName())
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: true}
	}
}

func collectPathsTrait(t *Trait, elems []string, pathElem []string, out map[string]elemPath) {
	newElems := appendNew(elems, t.GetName())
	newPath := appendNew(pathElem, t.GetFileName())
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: false}

	for _, f := range t.Fields {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#Fields")
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: true}
	}
	for _, f := range t.Functions {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#"+f.GetName())
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: true}
	}
}

func replaceLinks(text string, elems []string, modElems int, lookup map[string]elemPath, t *template.Template) (string, error) {
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

		entry, linkText, parts, ok := toLink(link, elems, modElems, lookup)
		if !ok {
			continue
		}

		var basePath string
		if entry.IsSection {
			basePath = path.Join(parts[:len(parts)-1]...)
		} else {
			basePath = path.Join(parts...)
		}
		pathStr := strings.Builder{}
		err := t.ExecuteTemplate(&pathStr, entry.Kind+"_path.md", basePath)
		if err != nil {
			return "", err
		}
		if entry.IsSection {
			pathStr.WriteString(parts[len(parts)-1])
		}
		text = fmt.Sprintf("%s[%s](%s)%s", text[:start], linkText, pathStr.String(), text[end:])
		fmt.Println(link, "-->", pathStr.String())
	}
	return text, nil
}

func toLink(link string, elems []string, modElems int, lookup map[string]elemPath) (*elemPath, string, []string, bool) {
	if strings.HasPrefix(link, ".") {
		return toRelLink(link, elems, modElems, lookup)
	}
	return toAbsLink(link, elems, modElems, lookup)
}

func toRelLink(link string, elems []string, modElems int, lookup map[string]elemPath) (*elemPath, string, []string, bool) {
	dots := 0
	fullPath := []string{}
	for strings.HasPrefix(link[dots:], ".") {
		if dots > 0 {
			fullPath = append(fullPath, "..")
		}
		dots++
	}
	if dots > modElems {
		log.Printf("Too many leading dots in cross ref: %s", link)
		return nil, "", nil, false
	}
	fullLink := link
	linkText := link[dots:]
	subElems := elems[:modElems-(dots-1)]

	if len(subElems) == 0 {
		fullLink = link[dots:]
	} else {
		fullLink = strings.Join(subElems, ".") + "." + linkText
	}

	elemPath, ok := lookup[fullLink]
	if !ok {
		log.Printf("Can't resolve cross ref: %s (%s)", link, fullLink)
		return nil, "", nil, false
	}

	fullPath = append(fullPath, elemPath.Elements[len(subElems):]...)
	return &elemPath, link[dots:], fullPath, true
}

func toAbsLink(link string, elems []string, modElems int, lookup map[string]elemPath) (*elemPath, string, []string, bool) {
	elemPath, ok := lookup[link]
	if !ok {
		log.Printf("Can't resolve cross ref: %s", link)
		return nil, "", nil, false
	}
	skip := 0
	for range modElems {
		if len(elemPath.Elements) <= skip {
			break
		}
		if elemPath.Elements[skip] == elems[skip] {
			skip++
		} else {
			break
		}
	}
	fullPath := []string{}
	for range modElems - skip {
		fullPath = append(fullPath, "..")
	}
	fullPath = append(fullPath, elemPath.Elements[skip:]...)
	return &elemPath, link, fullPath, true
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
