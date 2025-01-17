package document

import (
	"fmt"
)

// Filters and re-structures docs for package re-exports.
//
// Also collects a lookup, mapping from original to altered cross-refs.
func (proc *Processor) filterPackages() error {
	proc.linkExports = map[string]string{}
	anyExports := proc.collectExports(proc.Docs.Decl, nil)

	if !proc.UseExports {
		proc.ExportDocs = &Docs{
			Version: proc.Docs.Version,
			Decl:    proc.Docs.Decl,
		}
		return nil
	}

	if !anyExports {
		return fmt.Errorf("no package re-exports found. Given flag '--exports', there will be no output.\n       Add exports or run without flag '--exports'")
	}

	proc.ExportDocs = &Docs{
		Version: proc.Docs.Version,
		Decl:    proc.Docs.Decl.linkedCopy(),
	}
	proc.filterPackage(proc.Docs.Decl, proc.ExportDocs.Decl, nil, nil)
	return nil
}

func (proc *Processor) filterPackage(src, rootOut *Package, oldPath, newPath []string) {
	rootExports := map[string]*members{}
	collectExportsPackage(src, rootExports)

	oldPath = appendNew(oldPath, src.Name)
	newPath = appendNew(newPath, src.Name)

	for _, mod := range src.Modules {
		if mems, ok := rootExports[mod.Name]; ok {
			proc.collectModuleExports(mod, rootOut, oldPath, newPath, mems)
		}
	}

	for _, pkg := range src.Packages {
		if mems, ok := rootExports[pkg.Name]; ok {
			proc.collectPackageExports(pkg, rootOut, oldPath, newPath, mems)
		}
	}
}

func (proc *Processor) collectPackageExports(src, rootOut *Package, oldPath, newPath []string, rootMembers *members) {
	selfIncluded, toCrawl := collectExportMembers(rootMembers)

	if selfIncluded {
		newPkg := src.linkedCopy()
		proc.filterPackage(src, newPkg, oldPath, newPath)
		rootOut.Packages = append(rootOut.Packages, newPkg)

		tempOldPath := appendNew(oldPath, src.Name)
		tempNewPath := appendNew(newPath, src.Name)
		proc.addLinkExport(tempOldPath, tempNewPath)
	}
	oldPath = appendNew(oldPath, src.Name)

	for _, mod := range src.Modules {
		if mems, ok := toCrawl[mod.Name]; ok {
			proc.collectModuleExports(mod, rootOut, oldPath, newPath, mems)
		}
	}

	for _, pkg := range src.Packages {
		if mems, ok := toCrawl[pkg.Name]; ok {
			proc.collectPackageExports(pkg, rootOut, oldPath, newPath, mems)
		}
	}

	for _, elem := range src.Structs {
		proc.collectExportsStruct(elem, oldPath, newPath)
	}
	for _, elem := range src.Traits {
		proc.collectExportsTrait(elem, oldPath, newPath)
	}
	for _, elem := range src.Functions {
		proc.collectExportsFunction(elem, oldPath, newPath)
	}
}

func (proc *Processor) collectModuleExports(src *Module, rootOut *Package, oldPath, newPath []string, rootMembers *members) {
	selfIncluded, toCrawl := collectExportMembers(rootMembers)

	oldPath = appendNew(oldPath, src.Name)
	if selfIncluded {
		rootOut.Modules = append(rootOut.Modules, src)

		tempNewPath := appendNew(newPath, src.Name)
		proc.addLinkExport(oldPath, tempNewPath)
	}

	for _, elem := range src.Structs {
		if _, ok := toCrawl[elem.Name]; ok {
			rootOut.Structs = append(rootOut.Structs, elem)
			proc.collectExportsStruct(elem, oldPath, newPath)
		}
	}
	for _, elem := range src.Traits {
		if _, ok := toCrawl[elem.Name]; ok {
			rootOut.Traits = append(rootOut.Traits, elem)
			proc.collectExportsTrait(elem, oldPath, newPath)
		}
	}
	for _, elem := range src.Functions {
		if _, ok := toCrawl[elem.Name]; ok {
			rootOut.Functions = append(rootOut.Functions, elem)
			proc.collectExportsFunction(elem, oldPath, newPath)
		}
	}
}

func (proc *Processor) collectExportsStruct(s *Struct, oldPath []string, newPath []string) {
	oldPath = appendNew(oldPath, s.Name)
	newPath = appendNew(newPath, s.Name)

	proc.addLinkExport(oldPath, newPath)

	collectExportsList(proc, s.Parameters, oldPath, newPath)
	collectExportsList(proc, s.Fields, oldPath, newPath)
	collectExportsList(proc, s.Functions, oldPath, newPath)
}

func (proc *Processor) collectExportsTrait(s *Trait, oldPath []string, newPath []string) {
	oldPath = appendNew(oldPath, s.Name)
	newPath = appendNew(newPath, s.Name)

	proc.addLinkExport(oldPath, newPath)

	collectExportsList(proc, s.Fields, oldPath, newPath)
	collectExportsList(proc, s.Functions, oldPath, newPath)
}

func (proc *Processor) collectExportsFunction(s *Function, oldPath []string, newPath []string) {
	oldPath = appendNew(oldPath, s.Name)
	newPath = appendNew(newPath, s.Name)

	proc.addLinkExport(oldPath, newPath)
}

func collectExportsList[T Named](proc *Processor, sl []T, oldPath, newPath []string) {
	for _, elem := range sl {
		oldPath := appendNew(oldPath, elem.GetName())
		newPath := appendNew(newPath, elem.GetName())
		proc.addLinkExport(oldPath, newPath)
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
