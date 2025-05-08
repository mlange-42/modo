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

type pathHelper struct {
	AddPathFunc addPathFunc
	SetLink     bool
}

// Collects lookup for link target paths.
// Runs on the re-structured package.
func (proc *Processor) collectPaths() {
	pc := pathHelper{
		AddPathFunc: proc.addLinkTarget,
		SetLink:     false,
	}
	pc.collectPathsPackage(proc.ExportDocs.Decl, []string{}, []string{})
}

// Collects the paths of all (sub)-elements in the original structure.
func (proc *Processor) collectElementPaths() {
	pc := pathHelper{
		AddPathFunc: proc.addElementPath,
		SetLink:     true,
	}
	pc.collectPathsPackage(proc.Docs.Decl, []string{}, []string{})
}

func (pc *pathHelper) collectPathsPackage(p *Package, elems []string, pathElem []string) {
	newElems := appendNew(elems, p.GetName())
	newPath := appendNew(pathElem, p.GetFileName())
	if pc.SetLink {
		p.SetLink(newElems, p.Kind)
	}
	pc.AddPathFunc(p, newElems, newPath, "package", false)

	for _, pkg := range p.Packages {
		pc.collectPathsPackage(pkg, newElems, newPath)
	}
	for _, mod := range p.Modules {
		pc.collectPathsModule(mod, newElems, newPath)
	}

	for _, e := range p.Structs {
		pc.collectPathsStruct(e, newElems, newPath)
	}

	for _, e := range p.Traits {
		pc.collectPathsTrait(e, newElems, newPath)
	}
	for _, e := range p.Aliases {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#aliases")
		pc.AddPathFunc(e, newElems, newPath, "package", true) // kind=package for correct link paths
	}
	for _, e := range p.Functions {
		if pc.SetLink {
			e.SetLink(newElems, e.Kind)
		}
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, e.GetFileName())
		pc.AddPathFunc(e, newElems, newPath, "function", false)
	}
}

func (pc *pathHelper) collectPathsModule(m *Module, elems []string, pathElem []string) {
	newElems := appendNew(elems, m.GetName())
	newPath := appendNew(pathElem, m.GetFileName())
	if pc.SetLink {
		m.SetLink(newElems, m.Kind)
	}
	pc.AddPathFunc(m, newElems, newPath, "module", false)

	for _, e := range m.Structs {
		pc.collectPathsStruct(e, newElems, newPath)
	}
	for _, e := range m.Traits {
		pc.collectPathsTrait(e, newElems, newPath)
	}
	for _, e := range m.Aliases {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#aliases")
		pc.AddPathFunc(e, newElems, newPath, "module", true) // kind=module for correct link paths
	}
	for _, f := range m.Functions {
		if pc.SetLink {
			f.SetLink(newElems, f.Kind)
		}
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, f.GetFileName())
		pc.AddPathFunc(f, newElems, newPath, "function", false)
	}
}

func (pc *pathHelper) collectPathsStruct(s *Struct, elems []string, pathElem []string) {
	newElems := appendNew(elems, s.GetName())
	newPath := appendNew(pathElem, s.GetFileName())
	if pc.SetLink {
		s.SetLink(elems, s.Kind)
	}
	pc.AddPathFunc(s, newElems, newPath, "struct", false)

	for _, e := range s.Aliases {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#aliases")
		pc.AddPathFunc(e, newElems, newPath, "member", true)
	}
	for _, e := range s.Parameters {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#parameters")
		pc.AddPathFunc(e, newElems, newPath, "member", true)
	}
	for _, e := range s.Fields {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#fields")
		pc.AddPathFunc(e, newElems, newPath, "member", true)
	}
	for _, e := range s.Functions {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#"+strings.ToLower(e.GetName()))
		pc.AddPathFunc(e, newElems, newPath, "member", true)
	}
}

func (pc *pathHelper) collectPathsTrait(t *Trait, elems []string, pathElem []string) {
	newElems := appendNew(elems, t.GetName())
	newPath := appendNew(pathElem, t.GetFileName())
	if pc.SetLink {
		t.SetLink(elems, t.Kind)
	}
	pc.AddPathFunc(t, newElems, newPath, "trait", false)

	for _, e := range t.Aliases {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#aliases")
		pc.AddPathFunc(e, newElems, newPath, "member", true)
	}
	for _, e := range t.Fields {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#fields")
		pc.AddPathFunc(e, newElems, newPath, "member", true)
	}
	for _, e := range t.Functions {
		newElems := appendNew(newElems, e.GetName())
		newPath := appendNew(newPath, "#"+strings.ToLower(e.GetName()))
		pc.AddPathFunc(e, newElems, newPath, "member", true)
	}
}
