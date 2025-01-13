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
	for _, tr := range m.Traits {
		err := processLinksTrait(tr, newElems, lookup, t)
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

	for _, p := range s.Parameters {
		p.Description, err = replaceLinks(p.Description, newElems, len(elems), lookup, t)
		if err != nil {
			return err
		}
	}
	for _, f := range s.Fields {
		f.Summary, err = replaceLinks(f.Summary, newElems, len(elems), lookup, t)
		if err != nil {
			return err
		}
		f.Description, err = replaceLinks(f.Description, newElems, len(elems), lookup, t)
		if err != nil {
			return err
		}
	}
	for _, f := range s.Functions {
		if err := processLinksMethod(f, elems, lookup, t); err != nil {
			return err
		}
	}

	return nil
}

func processLinksTrait(tr *Trait, elems []string, lookup map[string]elemPath, t *template.Template) error {
	newElems := appendNew(elems, tr.GetName())

	var err error
	tr.Summary, err = replaceLinks(tr.Summary, newElems, len(elems), lookup, t)
	if err != nil {
		return err
	}
	tr.Description, err = replaceLinks(tr.Description, newElems, len(elems), lookup, t)
	if err != nil {
		return err
	}

	// TODO: add when traits support parameters
	/*for _, p := range tr.Parameters {
		p.Description, err = replaceLinks(p.Description, newElems, len(elems), lookup, t)
		if err != nil {
			return err
		}
	}*/
	for _, f := range tr.Fields {
		f.Summary, err = replaceLinks(f.Summary, newElems, len(elems), lookup, t)
		if err != nil {
			return err
		}
		f.Description, err = replaceLinks(f.Description, newElems, len(elems), lookup, t)
		if err != nil {
			return err
		}
	}
	for _, f := range tr.Functions {
		if err := processLinksMethod(f, elems, lookup, t); err != nil {
			return err
		}
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
	f.ReturnsDoc, err = replaceLinks(f.ReturnsDoc, newElems, len(elems), lookup, t)
	if err != nil {
		return err
	}
	f.RaisesDoc, err = replaceLinks(f.RaisesDoc, newElems, len(elems), lookup, t)
	if err != nil {
		return err
	}

	for _, a := range f.Args {
		a.Description, err = replaceLinks(a.Description, newElems, len(elems), lookup, t)
		if err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		p.Description, err = replaceLinks(p.Description, newElems, len(elems), lookup, t)
		if err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := processLinksFunction(o, elems, lookup, t)
		if err != nil {
			return err
		}
	}

	return nil
}

func processLinksMethod(f *Function, elems []string, lookup map[string]elemPath, t *template.Template) error {
	var err error
	f.Summary, err = replaceLinks(f.Summary, elems, len(elems), lookup, t)
	if err != nil {
		return err
	}
	f.Description, err = replaceLinks(f.Description, elems, len(elems), lookup, t)
	if err != nil {
		return err
	}
	f.ReturnsDoc, err = replaceLinks(f.ReturnsDoc, elems, len(elems), lookup, t)
	if err != nil {
		return err
	}
	f.RaisesDoc, err = replaceLinks(f.RaisesDoc, elems, len(elems), lookup, t)
	if err != nil {
		return err
	}

	for _, a := range f.Args {
		a.Description, err = replaceLinks(a.Description, elems, len(elems), lookup, t)
		if err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		p.Description, err = replaceLinks(p.Description, elems, len(elems), lookup, t)
		if err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := processLinksMethod(o, elems, lookup, t)
		if err != nil {
			return err
		}
	}

	return nil
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
	}
	return text, nil
}

func toLink(link string, elems []string, modElems int, lookup map[string]elemPath) (entry *elemPath, text string, parts []string, ok bool) {
	linkParts := strings.SplitN(link, " ", 2)
	if strings.HasPrefix(link, ".") {
		entry, text, parts, ok = toRelLink(linkParts[0], elems, modElems, lookup)
	} else {
		entry, text, parts, ok = toAbsLink(linkParts[0], elems, modElems, lookup)
	}
	if len(linkParts) > 1 {
		text = linkParts[1]
	} else {
		text = fmt.Sprintf("`%s`", text)
	}
	return
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
		log.Printf("WARNING: Too many leading dots in cross ref '%s' in %s", link, strings.Join(elems, "."))
		return nil, "", nil, false
	}
	linkText := link[dots:]
	subElems := elems[:modElems-(dots-1)]
	var fullLink string
	if len(subElems) == 0 {
		fullLink = linkText
	} else {
		fullLink = strings.Join(subElems, ".") + "." + linkText
	}

	elemPath, ok := lookup[fullLink]
	if !ok {
		log.Printf("WARNING: Can't resolve cross ref '%s' (%s) in %s", link, fullLink, strings.Join(elems, "."))
		return nil, "", nil, false
	}

	fullPath = append(fullPath, elemPath.Elements[len(subElems):]...)
	return &elemPath, link[dots:], fullPath, true
}

func toAbsLink(link string, elems []string, modElems int, lookup map[string]elemPath) (*elemPath, string, []string, bool) {
	elemPath, ok := lookup[link]
	if !ok {
		log.Printf("WARNING: Can't resolve cross ref '%s' in %s", link, strings.Join(elems, "."))
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
