package document

import "strings"

type elemPath struct {
	Elements  []string
	Kind      string
	IsSection bool
}

func (proc *Processor) collectPaths() {
	proc.linkTargets = map[string]elemPath{}
	proc.collectPathsPackage(proc.ExportDocs.Decl, []string{}, []string{})
}

func (proc *Processor) collectPathsPackage(p *Package, elems []string, pathElem []string) {
	newElems := appendNew(elems, p.GetName())
	newPath := appendNew(pathElem, p.GetFileName())
	proc.addLinkTarget(newElems, newPath, "package", false)

	for _, pkg := range p.Packages {
		proc.collectPathsPackage(pkg, newElems, newPath)
	}
	for _, mod := range p.Modules {
		proc.collectPathsModule(mod, newElems, newPath)
	}

	for _, s := range p.Structs {
		proc.collectPathsStruct(s, newElems, newPath)
	}
	for _, t := range p.Traits {
		proc.collectPathsTrait(t, newElems, newPath)
	}
	for _, f := range p.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, f.GetFileName())
		proc.addLinkTarget(newElems, newPath, "function", false)
	}
}

func (proc *Processor) collectPathsModule(m *Module, elems []string, pathElem []string) {
	newElems := appendNew(elems, m.GetName())
	newPath := appendNew(pathElem, m.GetFileName())
	proc.addLinkTarget(newElems, newPath, "module", false)

	for _, s := range m.Structs {
		proc.collectPathsStruct(s, newElems, newPath)
	}
	for _, t := range m.Traits {
		proc.collectPathsTrait(t, newElems, newPath)
	}
	for _, f := range m.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, f.GetFileName())
		proc.addLinkTarget(newElems, newPath, "function", false)
	}
}

func (proc *Processor) collectPathsStruct(s *Struct, elems []string, pathElem []string) {
	newElems := appendNew(elems, s.GetName())
	newPath := appendNew(pathElem, s.GetFileName())
	proc.addLinkTarget(newElems, newPath, "struct", false)

	for _, f := range s.Parameters {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#parameters")
		proc.addLinkTarget(newElems, newPath, "member", true)
	}
	for _, f := range s.Fields {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#fields")
		proc.addLinkTarget(newElems, newPath, "member", true)
	}
	for _, f := range s.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#"+strings.ToLower(f.GetName()))
		proc.addLinkTarget(newElems, newPath, "member", true)
	}
}

func (proc *Processor) collectPathsTrait(t *Trait, elems []string, pathElem []string) {
	newElems := appendNew(elems, t.GetName())
	newPath := appendNew(pathElem, t.GetFileName())
	proc.addLinkTarget(newElems, newPath, "trait", false)

	for _, f := range t.Fields {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#fields")
		proc.addLinkTarget(newElems, newPath, "member", true)
	}
	for _, f := range t.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#"+strings.ToLower(f.GetName()))
		proc.addLinkTarget(newElems, newPath, "member", true)
	}
}
