package document

import (
	"fmt"
	"strings"
)

// Removes `__init__` packages, creates missing struct signatures (Mojo bug?).
func cleanup(doc *Docs) {
	cleanupPackage(doc.Decl)
}

func cleanupPackage(p *Package) {
	if p.MemberSummary == nil {
		p.MemberSummary = newSummary("")
	}
	if p.MemberDescription == nil {
		p.MemberDescription = newDescription("")
	}

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
			s.Signature = createSignature(fmt.Sprintf("struct %s", s.GetName()), s.Parameters)
		}
	}
}

func fixAliasSignatures(doc *Docs) {
	fixAliasSignaturesPackage(doc.Decl)
}

func fixAliasSignaturesPackage(p *Package) {
	for _, pp := range p.Packages {
		fixAliasSignaturesPackage(pp)
	}
	for _, m := range p.Modules {
		fixAliasSignaturesModule(m)
	}

	for _, a := range p.Aliases {
		a.Signature = createSignature(a.GetName(), a.Parameters)
	}
	for _, s := range p.Structs {
		for _, a := range s.Aliases {
			a.Signature = createSignature(a.GetName(), a.Parameters)
		}
	}
	for _, t := range p.Traits {
		for _, a := range t.Aliases {
			a.Signature = createSignature(a.GetName(), a.Parameters)
		}
	}
}

func fixAliasSignaturesModule(m *Module) {
	for _, a := range m.Aliases {
		a.Signature = createSignature(a.GetName(), a.Parameters)
	}
	for _, s := range m.Structs {
		for _, a := range s.Aliases {
			a.Signature = createSignature(a.GetName(), a.Parameters)
		}
	}
	for _, t := range m.Traits {
		for _, a := range t.Aliases {
			a.Signature = createSignature(a.GetName(), a.Parameters)
		}
	}
}

func createSignature(prefix string, pars []*Parameter) string {
	b := strings.Builder{}
	b.WriteString(prefix)

	if len(pars) == 0 {
		return b.String()
	}

	b.WriteString("[")

	prevKind := ""
	for i, par := range pars {
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
