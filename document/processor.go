package document

import (
	"log"
	"strings"
	"text/template"
)

type Processor struct {
	Template    *template.Template
	Formatter   Formatter
	UseExports  bool
	ShortLinks  bool
	Docs        *Docs
	ExportDocs  *Docs
	linkTargets map[string]elemPath
	linkExports map[string]string
}

func NewProcessor(docs *Docs, f Formatter, t *template.Template, useExports bool, shortLinks bool) Processor {
	return Processor{
		Template:   t,
		Formatter:  f,
		UseExports: useExports,
		ShortLinks: shortLinks,
		Docs:       docs,
	}
}

func (proc *Processor) getLinkPath(link string) (elemPath, bool) {
	newLink := link

	if proc.UseExports {
		var ok bool
		newLink, ok = proc.linkExports[link]
		if !ok {
			log.Printf("link export lookup failed for %s", link)
			return elemPath{}, false
		}
	}

	elem, ok := proc.linkTargets[newLink]
	if !ok {
		log.Printf("link target lookup failed for %s (%s)", newLink, link)
		return elemPath{}, false
	}

	return elem, true
}

func (proc *Processor) addLinkExport(oldPath, newPath []string) {
	proc.linkExports[strings.Join(oldPath, ".")] = strings.Join(newPath, ".")
}
