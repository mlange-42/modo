package document

import (
	"fmt"
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
			links = append(links, r[6], r[7])
		}
	}

	return links, nil
}

func collectPaths(doc *Docs) map[string][]string {
	out := map[string][]string{}
	collectPathsPackage(doc.Decl, []string{}, []string{}, out)
	return out
}

func collectPathsPackage(p *Package, elems []string, pathElem []string, out map[string][]string) {
	newElems := appendNew(elems, p.GetName())
	newPath := appendNew(pathElem, p.GetFileName())
	out[strings.Join(newElems, ".")] = newPath

	for _, pkg := range p.Packages {
		collectPathsPackage(pkg, newElems, newPath, out)
	}
	for _, mod := range p.Modules {
		collectPathsModule(mod, newElems, newPath, out)
	}
}

func collectPathsModule(m *Module, elems []string, pathElem []string, out map[string][]string) {
	newElems := appendNew(elems, m.GetName())
	newPath := appendNew(pathElem, m.GetFileName())
	out[strings.Join(newElems, ".")] = newPath

	for _, s := range m.Structs {
		collectPathsStruct(s, newElems, newPath, out)
	}
	for _, t := range m.Traits {
		collectPathsTrait(t, newElems, newPath, out)
	}
	for _, f := range m.Functions {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, f.GetFileName())
		out[strings.Join(newElems, ".")] = newPath
	}
}

func collectPathsStruct(s *Struct, elems []string, pathElem []string, out map[string][]string) {
	newElems := appendNew(elems, s.GetName())
	newPath := appendNew(pathElem, s.GetFileName())
	out[strings.Join(newElems, ".")] = newPath

	for _, f := range s.Parameters {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#Parameters")
		out[strings.Join(newElems, ".")] = newPath
	}
	for _, f := range s.Fields {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#Fields")
		out[strings.Join(newElems, ".")] = newPath
	}
	for _, f := range s.Functions {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#"+f.GetName())
		out[strings.Join(newElems, ".")] = newPath
	}
}

func collectPathsTrait(t *Trait, elems []string, pathElem []string, out map[string][]string) {
	newElems := appendNew(elems, t.GetName())
	newPath := appendNew(pathElem, t.GetFileName())
	out[strings.Join(newElems, ".")] = newPath

	for _, f := range t.Fields {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#Fields")
		out[strings.Join(newElems, ".")] = newPath
	}
	for _, f := range t.Functions {
		newElems := appendNew(elems, f.GetName())
		newPath := appendNew(pathElem, "#"+f.GetName())
		out[strings.Join(newElems, ".")] = newPath
	}
}

func appendNew[T any](sl []T, elem T) []T {
	sl2 := make([]T, len(sl)+1)
	copy(sl2, sl)
	sl2[len(sl)] = elem
	return sl2
}

func cleanup(doc *Docs) {
	cleanupPackage(doc.Decl)
}

func cleanupPackage(p *Package) {
	for _, pp := range p.Packages {
		cleanupPackage(pp)
	}
	newModules := make([]*Module, 0, len(p.Modules))
	for _, m := range p.Modules {
		cleanupModule(m)
		if m.GetName() != "__init__" {
			newModules = append(newModules, m)
		}
	}
	p.Modules = newModules
}

func cleanupModule(m *Module) {
	for _, s := range m.Structs {
		if s.Signature == "" {
			s.Signature = createSignature(s)
		}
	}
}

func createSignature(s *Struct) string {
	b := strings.Builder{}
	b.WriteString("struct ")
	b.WriteString(s.GetName())

	if len(s.Parameters) == 0 {
		return b.String()
	}

	b.WriteString("[")

	prevKind := ""
	for i, par := range s.Parameters {
		written := false
		if par.PassingKind == "kw" && prevKind != par.PassingKind {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString("*")
			written = true
		}
		if prevKind == "inferred" && par.PassingKind != prevKind {
			b.WriteString(", //")
			written = true
		}
		if prevKind == "pos" && par.PassingKind != prevKind {
			b.WriteString(", /")
			written = true
		}

		if i > 0 || written {
			b.WriteString(", ")
		}

		b.WriteString(fmt.Sprintf("%s: %s", par.GetName(), par.Type))
		if len(par.Default) > 0 {
			b.WriteString(fmt.Sprintf(" = %s", par.Default))
		}

		prevKind = par.PassingKind
	}
	if prevKind == "inferred" {
		b.WriteString(", //")
	}
	if prevKind == "pos" {
		b.WriteString(", /")
	}

	b.WriteString("]")

	return b.String()
}
