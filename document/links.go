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
	err := proc.filterPackages()
	if err != nil {
		return err
	}
	proc.collectPaths()

	if !proc.UseExports {
		for k := range proc.linkTargets {
			proc.linkExports[k] = k
		}
	}
	if err := proc.processLinksPackage(proc.Docs.Decl, []string{}, true); err != nil {
		return err
	}
	if err := proc.processLinksPackage(proc.ExportDocs.Decl, []string{}, false); err != nil {
		return err
	}
	return nil
}

func (proc *Processor) processLinksPackage(p *Package, elems []string, firstPass bool) error {
	newElems := appendNew(elems, p.GetName())

	var err error
	p.Summary, err = proc.replaceRefs(p.Summary, newElems, len(newElems), firstPass)
	if err != nil {
		return err
	}
	p.Description, err = proc.replaceRefs(p.Description, newElems, len(newElems), firstPass)
	if err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		proc.processLinksPackage(pkg, newElems, firstPass)
	}
	for _, mod := range p.Modules {
		proc.processLinksModule(mod, newElems, firstPass)
	}

	for _, f := range p.Functions {
		err := proc.processLinksFunction(f, newElems, firstPass)
		if err != nil {
			return err
		}
	}
	for _, s := range p.Structs {
		err := proc.processLinksStruct(s, newElems, firstPass)
		if err != nil {
			return err
		}
	}
	for _, tr := range p.Traits {
		err := proc.processLinksTrait(tr, newElems, firstPass)
		if err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) processLinksModule(m *Module, elems []string, firstPass bool) error {
	newElems := appendNew(elems, m.GetName())

	var err error
	m.Summary, err = proc.replaceRefs(m.Summary, newElems, len(newElems), firstPass)
	if err != nil {
		return err
	}
	m.Description, err = proc.replaceRefs(m.Description, newElems, len(newElems), firstPass)
	if err != nil {
		return err
	}

	for _, f := range m.Functions {
		err := proc.processLinksFunction(f, newElems, firstPass)
		if err != nil {
			return err
		}
	}
	for _, s := range m.Structs {
		err := proc.processLinksStruct(s, newElems, firstPass)
		if err != nil {
			return err
		}
	}
	for _, tr := range m.Traits {
		err := proc.processLinksTrait(tr, newElems, firstPass)
		if err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) processLinksStruct(s *Struct, elems []string, firstPass bool) error {
	newElems := appendNew(elems, s.GetName())

	var err error
	s.Summary, err = proc.replaceRefs(s.Summary, newElems, len(elems), firstPass)
	if err != nil {
		return err
	}
	s.Description, err = proc.replaceRefs(s.Description, newElems, len(elems), firstPass)
	if err != nil {
		return err
	}

	for _, p := range s.Parameters {
		p.Description, err = proc.replaceRefs(p.Description, newElems, len(elems), firstPass)
		if err != nil {
			return err
		}
	}
	for _, f := range s.Fields {
		f.Summary, err = proc.replaceRefs(f.Summary, newElems, len(elems), firstPass)
		if err != nil {
			return err
		}
		f.Description, err = proc.replaceRefs(f.Description, newElems, len(elems), firstPass)
		if err != nil {
			return err
		}
	}
	for _, f := range s.Functions {
		if err := proc.processLinksMethod(f, elems, firstPass); err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) processLinksTrait(tr *Trait, elems []string, firstPass bool) error {
	newElems := appendNew(elems, tr.GetName())

	var err error
	tr.Summary, err = proc.replaceRefs(tr.Summary, newElems, len(elems), firstPass)
	if err != nil {
		return err
	}
	tr.Description, err = proc.replaceRefs(tr.Description, newElems, len(elems), firstPass)
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
		f.Summary, err = proc.replaceRefs(f.Summary, newElems, len(elems), firstPass)
		if err != nil {
			return err
		}
		f.Description, err = proc.replaceRefs(f.Description, newElems, len(elems), firstPass)
		if err != nil {
			return err
		}
	}
	for _, f := range tr.Functions {
		if err := proc.processLinksMethod(f, elems, firstPass); err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) processLinksFunction(f *Function, elems []string, firstPass bool) error {
	newElems := appendNew(elems, f.GetName())

	var err error
	f.Summary, err = proc.replaceRefs(f.Summary, newElems, len(elems), firstPass)
	if err != nil {
		return err
	}
	f.Description, err = proc.replaceRefs(f.Description, newElems, len(elems), firstPass)
	if err != nil {
		return err
	}
	f.ReturnsDoc, err = proc.replaceRefs(f.ReturnsDoc, newElems, len(elems), firstPass)
	if err != nil {
		return err
	}
	f.RaisesDoc, err = proc.replaceRefs(f.RaisesDoc, newElems, len(elems), firstPass)
	if err != nil {
		return err
	}

	for _, a := range f.Args {
		a.Description, err = proc.replaceRefs(a.Description, newElems, len(elems), firstPass)
		if err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		p.Description, err = proc.replaceRefs(p.Description, newElems, len(elems), firstPass)
		if err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := proc.processLinksFunction(o, elems, firstPass)
		if err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) processLinksMethod(f *Function, elems []string, firstPass bool) error {
	var err error
	f.Summary, err = proc.replaceRefs(f.Summary, elems, len(elems), firstPass)
	if err != nil {
		return err
	}
	f.Description, err = proc.replaceRefs(f.Description, elems, len(elems), firstPass)
	if err != nil {
		return err
	}
	f.ReturnsDoc, err = proc.replaceRefs(f.ReturnsDoc, elems, len(elems), firstPass)
	if err != nil {
		return err
	}
	f.RaisesDoc, err = proc.replaceRefs(f.RaisesDoc, elems, len(elems), firstPass)
	if err != nil {
		return err
	}

	for _, a := range f.Args {
		a.Description, err = proc.replaceRefs(a.Description, elems, len(elems), firstPass)
		if err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		p.Description, err = proc.replaceRefs(p.Description, elems, len(elems), firstPass)
		if err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := proc.processLinksMethod(o, elems, firstPass)
		if err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) replaceRefs(text string, elems []string, modElems int, firstPass bool) (string, error) {
	if firstPass {
		return proc.replaceRefsByPlaceholders(text, elems, modElems)
	}
	return proc.replacePlaceholdersByLinks(text, elems, modElems)
}

func (proc *Processor) replaceRefsByPlaceholders(text string, elems []string, modElems int) (string, error) {
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

		content, ok := proc.refToPlaceholder(link, elems, modElems)
		if !ok {
			continue
		}
		text = fmt.Sprintf("%s[%s]%s", text[:start], content, text[end:])
	}
	return text, nil
}

func (proc *Processor) replacePlaceholdersByLinks(text string, elems []string, modElems int) (string, error) {
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

		entry, linkText, parts, ok := proc.placeholderToLink(link, elems, modElems, proc.ShortLinks)
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

func (proc *Processor) placeholderToLink(link string, elems []string, modElems int, shorten bool) (entry *elemPath, text string, parts []string, ok bool) {
	linkParts := strings.SplitN(link, " ", 2)
	entry, text, parts, ok = proc.placeholderToAbsLink(linkParts[0], elems, modElems)
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

func (proc *Processor) placeholderToAbsLink(link string, elems []string, modElems int) (*elemPath, string, []string, bool) {
	elemPath, ok := proc.linkTargets[link]
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

func (proc *Processor) refToPlaceholder(link string, elems []string, modElems int) (string, bool) {
	linkParts := strings.SplitN(link, " ", 2)

	var placeholder string
	var ok bool
	if strings.HasPrefix(link, ".") {
		placeholder, ok = proc.refToPlaceholderRel(linkParts[0], elems, modElems)
	} else {
		placeholder, ok = proc.refToPlaceholderAbs(linkParts[0], elems)
	}
	if !ok {
		return "", false
	}
	if len(linkParts) > 1 {
		return fmt.Sprintf("%s %s", placeholder, linkParts[1]), true
	} else {
		return placeholder, true
	}
}

func (proc *Processor) refToPlaceholderRel(link string, elems []string, modElems int) (string, bool) {
	dots := 0
	for strings.HasPrefix(link[dots:], ".") {
		dots++
	}
	if dots > modElems {
		log.Printf("WARNING: Too many leading dots in cross ref '%s' in %s", link, strings.Join(elems, "."))
		return "", false
	}
	linkText := link[dots:]
	subElems := elems[:modElems-(dots-1)]
	var fullLink string
	if len(subElems) == 0 {
		fullLink = linkText
	} else {
		fullLink = strings.Join(subElems, ".") + "." + linkText
	}

	placeholder, ok := proc.linkExports[fullLink]
	if !ok {
		log.Printf("WARNING: Can't resolve cross ref '%s' (%s) in %s", link, fullLink, strings.Join(elems, "."))
		return "", false
	}
	return placeholder, true
}

func (proc *Processor) refToPlaceholderAbs(link string, elems []string) (string, bool) {
	placeholder, ok := proc.linkExports[link]
	if !ok {
		log.Printf("WARNING: Can't resolve cross ref '%s' in %s", link, strings.Join(elems, "."))
		return "", false
	}
	return placeholder, true
}
