package document

import "strings"

type elemPath struct {
	Elements  []string
	Kind      string
	IsSection bool
}

func (proc *Processor) collectPaths() map[string]elemPath {
	out := map[string]elemPath{}
	proc.collectPathsPackage(proc.Docs.Decl, []string{}, []string{}, out)
	return out
}

func (proc *Processor) collectPathsPackage(p *Package, elems []string, pathElem []string, out map[string]elemPath) {
	newElems := appendNew(elems, p.GetName())
	newPath := appendNew(pathElem, p.GetFileName())
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "package", IsSection: false}

	for _, pkg := range p.Packages {
		proc.collectPathsPackage(pkg, newElems, newPath, out)
	}
	for _, mod := range p.Modules {
		proc.collectPathsModule(mod, newElems, newPath, out)
	}
}

func (proc *Processor) collectPathsModule(m *Module, elems []string, pathElem []string, out map[string]elemPath) {
	newElems := appendNew(elems, m.GetName())
	newPath := appendNew(pathElem, m.GetFileName())
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "module", IsSection: false}

	for _, s := range m.Structs {
		proc.collectPathsStruct(s, newElems, newPath, out)
	}
	for _, t := range m.Traits {
		proc.collectPathsTrait(t, newElems, newPath, out)
	}
	for _, f := range m.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, f.GetFileName())
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: false}
	}
}

func (proc *Processor) collectPathsStruct(s *Struct, elems []string, pathElem []string, out map[string]elemPath) {
	newElems := appendNew(elems, s.GetName())
	newPath := appendNew(pathElem, s.GetFileName())
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: false}

	for _, f := range s.Parameters {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#parameters")
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: true}
	}
	for _, f := range s.Fields {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#fields")
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: true}
	}
	for _, f := range s.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#"+strings.ToLower(f.GetName()))
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: true}
	}
}

func (proc *Processor) collectPathsTrait(t *Trait, elems []string, pathElem []string, out map[string]elemPath) {
	newElems := appendNew(elems, t.GetName())
	newPath := appendNew(pathElem, t.GetFileName())
	out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: false}

	for _, f := range t.Fields {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#fields")
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: true}
	}
	for _, f := range t.Functions {
		newElems := appendNew(newElems, f.GetName())
		newPath := appendNew(newPath, "#"+strings.ToLower(f.GetName()))
		out[strings.Join(newElems, ".")] = elemPath{Elements: newPath, Kind: "member", IsSection: true}
	}
}
