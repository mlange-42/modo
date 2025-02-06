package document

import "fmt"

func (proc *Processor) processTranscludes(docs *Docs) error {
	return proc.walkAllDocStrings(docs, proc.replaceTranscludes, func(elem Named) string {
		return elem.GetName()
	})
}

func (proc *Processor) replaceTranscludes(text string, elems []string, modElems int) (string, error) {
	indices, err := findLinks(text, transcludeRegex)
	if err != nil {
		return "", err
	}
	if len(indices) == 0 {
		return text, nil
	}
	for i := len(indices) - 2; i >= 0; i -= 2 {
		start, end := indices[i], indices[i+1]
		link := text[start+1 : end-1]

		content, ok, err := proc.refToPlaceholder(link, elems, modElems)
		if err != nil {
			return "", err
		}
		if !ok {
			continue
		}
		text = fmt.Sprintf("%s{%s}%s", text[:start], content, text[end:])
	}
	return text, nil
}
