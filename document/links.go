package document

import (
	"fmt"
	"log"
	"path"
	"regexp"
	"strings"
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

func (proc *Processor) ProcessLinks() error {
	proc.filterPackages()
	lookup := proc.collectPaths()
	//for k, v := range lookup {
	//	fmt.Println(k, v.Elements)
	//}
	return proc.processLinksPackage(proc.Docs.Decl, []string{}, lookup)
}

func (proc *Processor) processLinksPackage(p *Package, elems []string, lookup map[string]elemPath) error {
	newElems := appendNew(elems, p.GetName())

	var err error
	p.Summary, err = proc.replaceLinks(p.Summary, newElems, len(newElems), lookup)
	if err != nil {
		return err
	}
	p.Description, err = proc.replaceLinks(p.Description, newElems, len(newElems), lookup)
	if err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		proc.processLinksPackage(pkg, newElems, lookup)
	}
	for _, mod := range p.Modules {
		proc.processLinksModule(mod, newElems, lookup)
	}

	return nil
}

func (proc *Processor) processLinksModule(m *Module, elems []string, lookup map[string]elemPath) error {
	newElems := appendNew(elems, m.GetName())

	var err error
	m.Summary, err = proc.replaceLinks(m.Summary, newElems, len(newElems), lookup)
	if err != nil {
		return err
	}
	m.Description, err = proc.replaceLinks(m.Description, newElems, len(newElems), lookup)
	if err != nil {
		return err
	}

	for _, f := range m.Functions {
		err := proc.processLinksFunction(f, newElems, lookup)
		if err != nil {
			return err
		}
	}
	for _, s := range m.Structs {
		err := proc.processLinksStruct(s, newElems, lookup)
		if err != nil {
			return err
		}
	}
	for _, tr := range m.Traits {
		err := proc.processLinksTrait(tr, newElems, lookup)
		if err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) processLinksStruct(s *Struct, elems []string, lookup map[string]elemPath) error {
	newElems := appendNew(elems, s.GetName())

	var err error
	s.Summary, err = proc.replaceLinks(s.Summary, newElems, len(elems), lookup)
	if err != nil {
		return err
	}
	s.Description, err = proc.replaceLinks(s.Description, newElems, len(elems), lookup)
	if err != nil {
		return err
	}

	for _, p := range s.Parameters {
		p.Description, err = proc.replaceLinks(p.Description, newElems, len(elems), lookup)
		if err != nil {
			return err
		}
	}
	for _, f := range s.Fields {
		f.Summary, err = proc.replaceLinks(f.Summary, newElems, len(elems), lookup)
		if err != nil {
			return err
		}
		f.Description, err = proc.replaceLinks(f.Description, newElems, len(elems), lookup)
		if err != nil {
			return err
		}
	}
	for _, f := range s.Functions {
		if err := proc.processLinksMethod(f, elems, lookup); err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) processLinksTrait(tr *Trait, elems []string, lookup map[string]elemPath) error {
	newElems := appendNew(elems, tr.GetName())

	var err error
	tr.Summary, err = proc.replaceLinks(tr.Summary, newElems, len(elems), lookup)
	if err != nil {
		return err
	}
	tr.Description, err = proc.replaceLinks(tr.Description, newElems, len(elems), lookup)
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
		f.Summary, err = proc.replaceLinks(f.Summary, newElems, len(elems), lookup)
		if err != nil {
			return err
		}
		f.Description, err = proc.replaceLinks(f.Description, newElems, len(elems), lookup)
		if err != nil {
			return err
		}
	}
	for _, f := range tr.Functions {
		if err := proc.processLinksMethod(f, elems, lookup); err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) processLinksFunction(f *Function, elems []string, lookup map[string]elemPath) error {
	newElems := appendNew(elems, f.GetName())

	var err error
	f.Summary, err = proc.replaceLinks(f.Summary, newElems, len(elems), lookup)
	if err != nil {
		return err
	}
	f.Description, err = proc.replaceLinks(f.Description, newElems, len(elems), lookup)
	if err != nil {
		return err
	}
	f.ReturnsDoc, err = proc.replaceLinks(f.ReturnsDoc, newElems, len(elems), lookup)
	if err != nil {
		return err
	}
	f.RaisesDoc, err = proc.replaceLinks(f.RaisesDoc, newElems, len(elems), lookup)
	if err != nil {
		return err
	}

	for _, a := range f.Args {
		a.Description, err = proc.replaceLinks(a.Description, newElems, len(elems), lookup)
		if err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		p.Description, err = proc.replaceLinks(p.Description, newElems, len(elems), lookup)
		if err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := proc.processLinksFunction(o, elems, lookup)
		if err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) processLinksMethod(f *Function, elems []string, lookup map[string]elemPath) error {
	var err error
	f.Summary, err = proc.replaceLinks(f.Summary, elems, len(elems), lookup)
	if err != nil {
		return err
	}
	f.Description, err = proc.replaceLinks(f.Description, elems, len(elems), lookup)
	if err != nil {
		return err
	}
	f.ReturnsDoc, err = proc.replaceLinks(f.ReturnsDoc, elems, len(elems), lookup)
	if err != nil {
		return err
	}
	f.RaisesDoc, err = proc.replaceLinks(f.RaisesDoc, elems, len(elems), lookup)
	if err != nil {
		return err
	}

	for _, a := range f.Args {
		a.Description, err = proc.replaceLinks(a.Description, elems, len(elems), lookup)
		if err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		p.Description, err = proc.replaceLinks(p.Description, elems, len(elems), lookup)
		if err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := proc.processLinksMethod(o, elems, lookup)
		if err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) replaceLinks(text string, elems []string, modElems int, lookup map[string]elemPath) (string, error) {
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

		entry, linkText, parts, ok := toLink(link, elems, modElems, lookup, proc.ShortLinks)
		if !ok {
			continue
		}

		var basePath string
		if entry.IsSection {
			basePath = path.Join(parts[:len(parts)-1]...)
		} else {
			basePath = path.Join(parts...)
		}
		pathStr, err := proc.Formatter.ToLinkPath(basePath, entry.Kind)
		if err != nil {
			return "", err
		}
		if entry.IsSection {
			pathStr += parts[len(parts)-1]
		}
		text = fmt.Sprintf("%s[%s](%s)%s", text[:start], linkText, pathStr, text[end:])
	}
	return text, nil
}

func toLink(link string, elems []string, modElems int, lookup map[string]elemPath, shorten bool) (entry *elemPath, text string, parts []string, ok bool) {
	linkParts := strings.SplitN(link, " ", 2)
	if strings.HasPrefix(link, ".") {
		entry, text, parts, ok = toRelLink(linkParts[0], elems, modElems, lookup)
	} else {
		entry, text, parts, ok = toAbsLink(linkParts[0], elems, modElems, lookup)
	}
	if !ok {
		return
	}
	if len(linkParts) > 1 {
		text = linkParts[1]
	} else {
		if shorten {
			textParts := strings.Split(text, ".")
			if entry.IsSection {
				text = strings.Join(textParts[len(textParts)-2:], ".")
			} else {
				text = textParts[len(textParts)-1]
			}
		}
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
	return &elemPath, linkText, fullPath, true
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
