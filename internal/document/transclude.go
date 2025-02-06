package document

import (
	"fmt"
	"strings"
)

func (proc *Processor) processTranscludes(docs *Docs) error {
	return proc.replaceTranscludesPackage(docs.Decl, []string{})
}

func (proc *Processor) replaceTranscludesPackage(p *Package, elems []string) error {
	newElems := appendNew(elems, p.Name)
	for _, pkg := range p.Packages {
		if err := proc.replaceTranscludesPackage(pkg, newElems); err != nil {
			return err
		}
	}
	for _, mod := range p.Modules {
		if err := proc.replaceTranscludesModule(mod, newElems); err != nil {
			return err
		}
	}
	return nil
}

func (proc *Processor) replaceTranscludesModule(m *Module, elems []string) error {
	newElems := appendNew(elems, m.Name)
	for _, s := range m.Structs {
		structElems := appendNew(newElems, s.Name)
		for _, f := range s.Functions {
			if len(f.Overloads) == 0 {
				if err := proc.replaceTranscludes(f, structElems, len(newElems)); err != nil {
					return err
				}
			} else {
				for _, o := range f.Overloads {
					if err := proc.replaceTranscludes(o, structElems, len(newElems)); err != nil {
						return err
					}
				}
			}
		}

	}
	return nil
}

func (proc *Processor) replaceTranscludes(elem *Function, elems []string, modElems int) error {
	indices, err := findLinks(elem.Summary, transcludeRegex, false)
	if err != nil {
		return err
	}
	if len(indices) == 0 {
		return nil
	}
	if len(indices) > 2 {
		return fmt.Errorf("multiple doc transclusions in %s", strings.Join(elems, "."))
	}
	start, end := indices[0], indices[1]
	link := elem.Summary[start+1 : end-1]

	content, ok, err := proc.refToPlaceholder(link, elems, modElems, false)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	from, ok := proc.allPaths[content]
	if !ok {
		return fmt.Errorf("doc transclusions source %s (%s) in %s not found", link, content, strings.Join(elems, "."))
	}
	trait, ok := from.(*Trait)
	if !ok {
		return fmt.Errorf("doc transclusions source %s (%s) in %s is not a trait", link, content, strings.Join(elems, "."))
	}

	sourceFunc, ok := findImplementedFunction(elem, trait)
	if !ok {
		return fmt.Errorf("no matching doc transclusion found for %s.%s in %s. Required signature: %s", link, elem.Name, strings.Join(elems, "."), elem.Signature)
	}

	elem.MemberSummary = sourceFunc.MemberSummary
	elem.Description = sourceFunc.Description
	elem.ReturnsDoc = sourceFunc.ReturnsDoc
	elem.RaisesDoc = sourceFunc.RaisesDoc

	for i, arg := range sourceFunc.Args {
		elem.Args[i].Description = arg.Description
	}
	for i, par := range sourceFunc.Parameters {
		elem.Parameters[i].Description = par.Description
	}

	return nil
}

func findImplementedFunction(structFunc *Function, trait *Trait) (*Function, bool) {
	var sourceFunc *Function
	for _, f := range trait.Functions {
		if f.Name != structFunc.Name {
			continue
		}
		if len(f.Overloads) == 0 {
			if functionsMatch(structFunc, f) {
				sourceFunc = f
			}
			break
		}
		for _, o := range f.Overloads {
			if functionsMatch(structFunc, o) {
				sourceFunc = o
				break
			}
		}
		break
	}
	return sourceFunc, sourceFunc != nil
}

func functionsMatch(s, t *Function) bool {
	if len(s.Args) != len(t.Args) {
		return false
	}
	if len(s.Parameters) != len(t.Parameters) {
		return false
	}

	for i := range s.Args {
		argS, argT := s.Args[i], t.Args[i]
		if argS.Type != argT.Type {
			if !(argS.Type == "Self" && argT.Type == "_Self") {
				return false
			}
		}
		if argS.PassingKind != argT.PassingKind {
			return false
		}
		if argS.Convention != argT.Convention {
			return false
		}
	}

	for i := range s.Parameters {
		parS, parT := s.Parameters[i], t.Parameters[i]
		if parS.Type != parT.Type {
			return false
		}
		if parS.PassingKind != parT.PassingKind {
			return false
		}
	}

	return true
}
