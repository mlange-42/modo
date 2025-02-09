package document

type renameHelper struct {
	rename map[string]string
}

func (proc *Processor) renameAll(p *Package) {
	r := renameHelper{rename: proc.renameExports}
	r.renamePackage(p, p.Name)
}

func (r *renameHelper) renamePackage(p *Package, ownPath string) {
	for i := range p.Packages {
		newPath := ownPath + "." + p.Packages[i].Name
		r.renamePackage(p.Packages[i], newPath)
		if newName, ok := r.rename[newPath]; ok {
			tempPkg := *p.Packages[i]
			tempPkg.MemberName = MemberName{Name: newName}
			p.Packages[i] = &tempPkg
		}
	}

	for i := range p.Modules {
		newPath := ownPath + "." + p.Modules[i].Name
		r.renameModule(p.Modules[i], newPath)
		if newName, ok := r.rename[newPath]; ok {
			tempMod := *p.Modules[i]
			tempMod.MemberName = MemberName{Name: newName}
			p.Modules[i] = &tempMod
		}
	}

	for i := range p.Aliases {
		newPath := ownPath + "." + p.Aliases[i].Name
		if newName, ok := r.rename[newPath]; ok {
			tempMod := *p.Aliases[i]
			tempMod.MemberName = MemberName{Name: newName}
			p.Aliases[i] = &tempMod
		}
	}
	for i := range p.Structs {
		newPath := ownPath + "." + p.Structs[i].Name
		if newName, ok := r.rename[newPath]; ok {
			tempMod := *p.Structs[i]
			tempMod.MemberName = MemberName{Name: newName}
			p.Structs[i] = &tempMod
		}
	}
	for i := range p.Traits {
		newPath := ownPath + "." + p.Traits[i].Name
		if newName, ok := r.rename[newPath]; ok {
			tempMod := *p.Traits[i]
			tempMod.MemberName = MemberName{Name: newName}
			p.Traits[i] = &tempMod
		}
	}
	for i := range p.Functions {
		newPath := ownPath + "." + p.Functions[i].Name
		if newName, ok := r.rename[newPath]; ok {
			tempMod := *p.Functions[i]
			tempMod.MemberName = MemberName{Name: newName}
			p.Functions[i] = &tempMod
		}
	}
}

func (r *renameHelper) renameModule(m *Module, ownPath string) {
	for i := range m.Aliases {
		newPath := ownPath + "." + m.Aliases[i].Name
		if newName, ok := r.rename[newPath]; ok {
			tempMod := *m.Aliases[i]
			tempMod.MemberName = MemberName{Name: newName}
			m.Aliases[i] = &tempMod
		}
	}
	for i := range m.Structs {
		newPath := ownPath + "." + m.Structs[i].Name
		if newName, ok := r.rename[newPath]; ok {
			tempMod := *m.Structs[i]
			tempMod.MemberName = MemberName{Name: newName}
			m.Structs[i] = &tempMod
		}
	}
	for i := range m.Traits {
		newPath := ownPath + "." + m.Traits[i].Name
		if newName, ok := r.rename[newPath]; ok {
			tempMod := *m.Traits[i]
			tempMod.MemberName = MemberName{Name: newName}
			m.Traits[i] = &tempMod
		}
	}
	for i := range m.Functions {
		newPath := ownPath + "." + m.Functions[i].Name
		if newName, ok := r.rename[newPath]; ok {
			tempMod := *m.Functions[i]
			tempMod.MemberName = MemberName{Name: newName}
			m.Functions[i] = &tempMod
		}
	}
}
