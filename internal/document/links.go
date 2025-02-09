package document

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

const linkRegexString = `(?s)(?:(` + "```.*?```)|(`.*?`" + `))|(\[[^\[].*?\])`
const transcludeRegexString = `(?s)(?:(` + "```.*?```)|(`.*?`" + `))|\[(\[.*?\])\]`

var linkRegex *regexp.Regexp
var transcludeRegex *regexp.Regexp

func init() {
	var err error
	linkRegex, err = regexp.Compile(linkRegexString)
	if err != nil {
		panic(err)
	}
	transcludeRegex, err = regexp.Compile(transcludeRegexString)
	if err != nil {
		panic(err)
	}
}

func (proc *Processor) processLinks(docs *Docs) error {
	w := walker{
		Func:     proc.replaceRefs,
		NameFunc: func(elem Named) string { return elem.GetName() },
	}
	return w.walkAllDocStrings(docs)
}

func (proc *Processor) replaceRefs(text string, elems []string, modElems int) (string, error) {
	indices, err := findLinks(text, linkRegex, true)
	if err != nil {
		return "", err
	}
	if len(indices) == 0 {
		return text, nil
	}
	for i := len(indices) - 2; i >= 0; i -= 2 {
		start, end := indices[i], indices[i+1]
		link := text[start+1 : end-1]

		content, ok, err := proc.refToPlaceholder(link, elems, modElems, true)
		if err != nil {
			return "", err
		}
		if !ok {
			continue
		}
		text = fmt.Sprintf("%s[%s]%s", text[:start], content, text[end:])
	}
	return text, nil
}

// ReplacePlaceholders replaces placeholders in the text with links to the corresponding elements.
func (proc *Processor) ReplacePlaceholders(text string, elems []string, modElems int) (string, error) {
	indices, err := findLinks(text, linkRegex, true)
	if err != nil {
		return "", err
	}
	if len(indices) == 0 {
		return text, nil
	}
	for i := len(indices) - 2; i >= 0; i -= 2 {
		start, end := indices[i], indices[i+1]
		link := text[start+1 : end-1]

		entry, linkText, parts, ok, err := proc.placeholderToLink(link, elems, modElems, proc.Config.ShortLinks)
		if err != nil {
			return "", err
		}
		if !ok {
			continue
		}

		var basePath string
		if entry.IsSection {
			basePath = path.Join(parts[:len(parts)-1]...)
		} else {
			basePath = path.Join(parts...)
		}

		pathStr := proc.Formatter.ToLinkPath(basePath, entry.Kind)
		if entry.IsSection {
			pathStr += parts[len(parts)-1]
		}
		text = fmt.Sprintf("%s[%s](%s)%s", text[:start], linkText, pathStr, text[end:])
	}
	return text, nil
}

func (proc *Processor) placeholderToLink(link string, elems []string, modElems int, shorten bool) (entry *elemPath, text string, parts []string, ok bool, err error) {
	linkParts := strings.SplitN(link, " ", 2)
	entry, text, parts, ok, err = proc.placeholderToRelLink(linkParts[0], elems, modElems)
	if err != nil {
		return
	}
	if !ok {
		return
	}
	if len(linkParts) > 1 {
		text = linkParts[1]
	} else {
		if shorten {
			textParts := strings.Split(text, ".")
			if entry.IsSection && entry.Kind != "package" && entry.Kind != "module" {
				text = strings.Join(textParts[len(textParts)-2:], ".")
			} else {
				text = textParts[len(textParts)-1]
			}
		}
		text = fmt.Sprintf("`%s`", text)
	}
	return
}

func (proc *Processor) placeholderToRelLink(link string, elems []string, modElems int) (*elemPath, string, []string, bool, error) {
	elemPath, ok := proc.linkTargets[link]
	if !ok {
		err := proc.warnOrError("Can't resolve cross ref placeholder '%s' in %s", link, strings.Join(elems, "."))
		return nil, "", nil, false, err
	}
	skip := 0
	for range modElems {
		if skip >= len(elemPath.Elements) {
			break
		}
		if elemPath.Elements[skip] == elems[skip] {
			skip++
		} else {
			break
		}
	}

	// redirect link to re-name by re-export
	link = proc.renameInLink(link, &elemPath)

	fullPath := []string{}
	for range modElems - skip {
		fullPath = append(fullPath, "..")
	}
	fullPath = append(fullPath, elemPath.Elements[skip:]...)
	if len(fullPath) == 0 {
		fullPath = []string{"."}
	}

	return &elemPath, link, fullPath, true, nil
}

func (proc *Processor) renameInLink(link string, elems *elemPath) string {
	if len(proc.renameExports) == 0 {
		return link
	}

	maxDepth := len(elems.Elements)
	if elems.IsSection {
		maxDepth--
	}

	newLink := strings.Split(link, ".")
	dotPos := 0
	currDepth := 0
	changed := false
	for currDepth < len(elems.Elements) {
		idx := strings.IndexRune(link[dotPos:], '.')
		if idx < 0 {
			idx = len(link) - dotPos
		}
		subLink := link[:dotPos+idx]
		if renamed, ok := proc.renameExports[subLink]; ok {
			if currDepth < maxDepth {
				newName := toFileName(renamed)
				elems.Elements[currDepth] = newName
			}
			newLink[currDepth] = renamed
			changed = true
		}
		dotPos += idx + 1
		currDepth++
	}
	if changed {
		return strings.Join(newLink, ".")
	}
	return link
}

func (proc *Processor) refToPlaceholder(link string, elems []string, modElems int, redirect bool) (string, bool, error) {
	linkParts := strings.SplitN(link, " ", 2)

	var placeholder string
	var ok bool
	var err error
	if strings.HasPrefix(link, ".") {
		placeholder, ok, err = proc.refToPlaceholderRel(linkParts[0], elems, modElems, redirect)
	} else {
		placeholder, ok, err = proc.refToPlaceholderAbs(linkParts[0], elems, redirect)
	}
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}

	if len(linkParts) > 1 {
		return fmt.Sprintf("%s %s", placeholder, linkParts[1]), true, nil
	}
	return placeholder, true, nil
}

func (proc *Processor) refToPlaceholderRel(link string, elems []string, modElems int, redirect bool) (string, bool, error) {
	dots := 0
	for strings.HasPrefix(link[dots:], ".") {
		dots++
	}
	if dots > modElems {
		err := proc.warnOrError("Too many leading dots in cross ref '%s' in %s", link, strings.Join(elems, "."))
		return "", false, err
	}
	linkText := link[dots:]
	subElems := elems[:modElems-(dots-1)]
	var fullLink string
	if len(subElems) == 0 {
		fullLink = linkText
	} else {
		fullLink = strings.Join(subElems, ".") + "." + linkText
	}

	if !redirect {
		return fullLink, true, nil
	}

	placeholder, ok := proc.linkExports[fullLink]
	if !ok {
		err := proc.warnOrError("Can't resolve cross ref (rel) '%s' (%s) in %s", link, fullLink, strings.Join(elems, "."))
		return "", false, err
	}
	return placeholder, true, nil
}

func (proc *Processor) refToPlaceholderAbs(link string, elems []string, redirect bool) (string, bool, error) {
	if !redirect {
		return link, true, nil
	}
	placeholder, ok := proc.linkExports[link]
	if !ok {
		err := proc.warnOrError("Can't resolve cross ref (abs) '%s' in %s", link, strings.Join(elems, "."))
		return "", false, err
	}
	return placeholder, true, nil
}

func findLinks(text string, regex *regexp.Regexp, isReference bool) ([]int, error) {
	links := []int{}
	results := regex.FindAllStringSubmatchIndex(text, -1)
	for _, r := range results {
		if r[6] >= 0 {
			// TODO: this can probably be moved into the REGEXP somehow.
			if isReference {
				// Excludes markdown links
				if len(text) > r[7] && string(text[r[7]]) == "(" {
					continue
				}
				// Excludes double square brackets
				if r[6] > 0 && string(text[r[6]-1]) == "[" {
					continue
				}
			}
			links = append(links, r[6], r[7])
		}
	}

	return links, nil
}
