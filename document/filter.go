package document

func (proc *Processor) filterPackages() {
	proc.collectExports(proc.Docs.Decl, nil)

	if !proc.UseExports {
		proc.ExportDocs = &Docs{
			Version: proc.Docs.Version,
			Decl:    proc.Docs.Decl,
		}
		return
	}

	proc.ExportDocs = &Docs{
		Version: proc.Docs.Version,
		Decl:    proc.Docs.Decl.copyEmpty(),
	}
	proc.filterPackage(proc.Docs.Decl, proc.ExportDocs.Decl)
}

func (proc *Processor) filterPackage(src, rootOut *Package) {
	rootExports := map[string]*members{}
	collectExportsPackage(src, rootExports)

	for _, mod := range src.Modules {
		if mems, ok := rootExports[mod.Name]; ok {
			proc.collectModuleExports(mod, rootOut, mems)
		}
	}

	for _, pkg := range src.Packages {
		if mems, ok := rootExports[pkg.Name]; ok {
			proc.collectPackageExports(pkg, rootOut, mems)
		}
	}
}

func (proc *Processor) collectPackageExports(src, rootOut *Package, rootMembers *members) {
	selfIncluded, toCrawl := collectExportMembers(rootMembers)
	if selfIncluded {
		rootOut.Packages = append(rootOut.Packages, src)
	}

	for _, mod := range src.Modules {
		if mems, ok := toCrawl[mod.Name]; ok {
			proc.collectModuleExports(mod, rootOut, mems)
		}
	}

	for _, pkg := range src.Packages {
		if mems, ok := toCrawl[pkg.Name]; ok {
			proc.collectPackageExports(pkg, rootOut, mems)
		}
	}
}

func (proc *Processor) collectModuleExports(src *Module, rootOut *Package, rootMembers *members) {
	selfIncluded, toCrawl := collectExportMembers(rootMembers)
	if selfIncluded {
		rootOut.Modules = append(rootOut.Modules, src)
	}

	for _, elem := range src.Structs {
		if _, ok := toCrawl[elem.Name]; ok {
			rootOut.Structs = append(rootOut.Structs, elem)
		}
	}
	for _, elem := range src.Traits {
		if _, ok := toCrawl[elem.Name]; ok {
			rootOut.Traits = append(rootOut.Traits, elem)
		}
	}
	for _, elem := range src.Functions {
		if _, ok := toCrawl[elem.Name]; ok {
			rootOut.Functions = append(rootOut.Functions, elem)
		}
	}
}

type members struct {
	Members []member
}

type member struct {
	Include bool
	RelPath []string
}

func collectExportMembers(parentMember *members) (selfIncluded bool, toCrawl map[string]*members) {
	selfIncluded = false
	toCrawl = map[string]*members{}

	for _, mem := range parentMember.Members {
		if mem.Include {
			selfIncluded = true
			continue
		}
		var newMember member
		if len(mem.RelPath) == 1 {
			newMember = member{Include: true}
		} else {
			newMember = member{RelPath: mem.RelPath[1:]}
		}
		if members, ok := toCrawl[mem.RelPath[0]]; ok {
			members.Members = append(members.Members, newMember)
			continue
		}
		toCrawl[mem.RelPath[0]] = &members{Members: []member{newMember}}
	}

	return
}

func collectExportsPackage(p *Package, out map[string]*members) {
	for _, ex := range p.Exports {
		var newMember member
		if len(ex.Short) == 1 {
			newMember = member{Include: true}
		} else {
			newMember = member{RelPath: ex.Short[1:]}
		}
		if members, ok := out[ex.Short[0]]; ok {
			members.Members = append(members.Members, newMember)
			continue
		}
		out[ex.Short[0]] = &members{Members: []member{newMember}}
	}
}
