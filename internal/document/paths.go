package document

import (
	"strings"
)

type elemPath struct {
	Elements  []string
	Kind      string
	IsSection bool
}

type addPathFunc = func(Named, []string, []string, string, bool)

// Collects lookup for link target paths.
// Runs on the re-structured package.
func (proc *Processor) collectPaths() {
	proc.linkTargets = map[string]elemPath{}
	proc.collectPathsPackage(proc.ExportDocs.Decl, []string{}, []string{}, proc.addLinkTarget, false)
}

// Collects the paths of all (sub)-elements in the original structure.
func (proc *Processor) collectElementPaths() {
	proc.allPaths = map[string]Named{}
	proc.collectPathsPackage(proc.Docs.Decl, []string{}, []string{}, proc.addElementPath, true)
}

func (proc *Processor) collectPathsPackage(p *Package, elems []string, pathElem []string, add addPathFunc, isOriginal bool) {
	newElems := appendNew(elems, p.GetName())
	newPath := appendNew(pathElem, p.GetFileName())
	if isOriginal {
		p.SetLink(newElems, p.Kind)
	}
	add(p, newElems, newPath, "package", false)

	for _, pkg := range p.Packages {
		proc.collectPathsPackage(pkg, newElems, newPath, add, isOriginal)
	}
	for _, mod := range p.Modules {
		proc.collectPathsModule(mod, newElems, newPath, add, isOriginal)
	}

	for _, e := range p.Structs {
		proc.collectPathsStruct(e, newElems, newPath, add, isOriginal)
	}

	for _, e := range p.Traits {
		proc.collectPathsTrait(e, newElems, newPath, add, isOriginal)
	}
	for _, e := range p.Aliases {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#aliases")
		add(e, newElems, newPath, "package", true) // kind=package for correct link paths
	}
	for _, e := range p.Functions {
		if isOriginal {
			e.SetLink(newElems, e.Kind)
		}
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, e.GetFileName())
		add(e, newElems, newPath, "function", false)
	}
}

func (proc *Processor) collectPathsModule(m *Module, elems []string, pathElem []string, add addPathFunc, isOriginal bool) {
	newElems := appendNew(elems, m.GetName())
	newPath := appendNew(pathElem, m.GetFileName())
	m.SetLink(newElems, m.Kind)
	add(m, newElems, newPath, "module", false)

	for _, e := range m.Structs {
		proc.collectPathsStruct(e, newElems, newPath, add, isOriginal)
	}
	for _, e := range m.Traits {
		proc.collectPathsTrait(e, newElems, newPath, add, isOriginal)
	}
	for _, e := range m.Aliases {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#aliases")
		add(e, newElems, newPath, "module", true) // kind=module for correct link paths
	}
	for _, f := range m.Functions {
		if isOriginal {
			f.SetLink(newElems, f.Kind)
		}
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, f.GetFileName())
		add(f, newElems, newPath, "function", false)
	}
}

func (proc *Processor) collectPathsStruct(s *Struct, elems []string, pathElem []string, add addPathFunc, isOriginal bool) {
	newElems := appendNew(elems, s.GetName())
	newPath := appendNew(pathElem, s.GetFileName())
	if isOriginal {
		s.SetLink(elems, s.Kind)
	}
	add(s, newElems, newPath, "struct", false)

	for _, e := range s.Aliases {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#aliases")
		add(e, newElems, newPath, "member", true)
	}
	for _, e := range s.Parameters {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#parameters")
		add(e, newElems, newPath, "member", true)
	}
	for _, e := range s.Fields {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#fields")
		add(e, newElems, newPath, "member", true)
	}
	for _, e := range s.Functions {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#"+strings.ToLower(e.GetName()))
		add(e, newElems, newPath, "member", true)
	}
}

func (proc *Processor) collectPathsTrait(t *Trait, elems []string, pathElem []string, add addPathFunc, isOriginal bool) {
	newElems := appendNew(elems, t.GetName())
	newPath := appendNew(pathElem, t.GetFileName())
	if isOriginal {
		t.SetLink(elems, t.Kind)
	}
	add(t, newElems, newPath, "trait", false)

	for _, e := range t.Fields {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#fields")
		add(e, newElems, newPath, "member", true)
	}
	for _, e := range t.Functions {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#"+strings.ToLower(e.GetName()))
		add(e, newElems, newPath, "member", true)
	}
}
