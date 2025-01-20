package document

import (
	"strings"
)

type elemPath struct {
	Elements  []string
	Kind      string
	IsSection bool
}

// Collects lookup for link target paths.
// Runs on the re-structured package.
func (proc *Processor) collectPaths() {
	proc.linkTargets = map[string]elemPath{}
	proc.collectPathsPackage(proc.ExportDocs.Decl, []string{}, []string{}, proc.addLinkTarget)
}

// Collects the paths of all (sub)-elements in the original structure.
func (proc *Processor) collectElementPaths() {
	proc.allPaths = map[string]bool{}
	proc.collectPathsPackage(proc.Docs.Decl, []string{}, []string{}, proc.addElementPath)
}

func (proc *Processor) collectPathsPackage(p *Package, elems []string, pathElem []string, add func([]string, []string, string, bool)) {
	newElems := appendNew(elems, p.GetName())
	newPath := appendNew(pathElem, p.GetFileName())
	add(newElems, newPath, "package", false)

	for _, pkg := range p.Packages {
		proc.collectPathsPackage(pkg, newElems, newPath, add)
	}
	for _, mod := range p.Modules {
		proc.collectPathsModule(mod, newElems, newPath, add)
	}

	for _, s := range p.Structs {
		proc.collectPathsStruct(s, newElems, newPath, add)
	}
	for _, t := range p.Traits {
		proc.collectPathsTrait(t, newElems, newPath, add)
	}
	for _, a := range p.Aliases {
		newElems := appendNew(newElems, a.GetName())
		newPath := appendNew(newPath, "#aliases")
		add(newElems, newPath, "package", true) // kind=package for correct link paths
	}
	for _, f := range p.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, f.GetFileName())
		add(newElems, newPath, "function", false)
	}
}

func (proc *Processor) collectPathsModule(m *Module, elems []string, pathElem []string, add func([]string, []string, string, bool)) {
	newElems := appendNew(elems, m.GetName())
	newPath := appendNew(pathElem, m.GetFileName())
	add(newElems, newPath, "module", false)

	for _, s := range m.Structs {
		proc.collectPathsStruct(s, newElems, newPath, add)
	}
	for _, t := range m.Traits {
		proc.collectPathsTrait(t, newElems, newPath, add)
	}
	for _, a := range m.Aliases {
		newElems := appendNew(newElems, a.GetName())
		newPath := appendNew(newPath, "#aliases")
		add(newElems, newPath, "module", true) // kind=module for correct link paths
	}
	for _, f := range m.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, f.GetFileName())
		add(newElems, newPath, "function", false)
	}
}

func (proc *Processor) collectPathsStruct(s *Struct, elems []string, pathElem []string, add func([]string, []string, string, bool)) {
	newElems := appendNew(elems, s.GetName())
	newPath := appendNew(pathElem, s.GetFileName())
	add(newElems, newPath, "struct", false)

	for _, f := range s.Aliases {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#aliases")
		add(newElems, newPath, "member", true)
	}
	for _, f := range s.Parameters {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#parameters")
		add(newElems, newPath, "member", true)
	}
	for _, f := range s.Fields {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#fields")
		add(newElems, newPath, "member", true)
	}
	for _, f := range s.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#"+strings.ToLower(f.GetName()))
		add(newElems, newPath, "member", true)
	}
}

func (proc *Processor) collectPathsTrait(t *Trait, elems []string, pathElem []string, add func([]string, []string, string, bool)) {
	newElems := appendNew(elems, t.GetName())
	newPath := appendNew(pathElem, t.GetFileName())
	add(newElems, newPath, "trait", false)

	for _, f := range t.Fields {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#fields")
		add(newElems, newPath, "member", true)
	}
	for _, f := range t.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#"+strings.ToLower(f.GetName()))
		add(newElems, newPath, "member", true)
	}
}
