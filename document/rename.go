package document

func (proc *Processor) renameAll(p *Package) {
	proc.renamePackage(p, p.Name)
}

func (proc *Processor) renamePackage(p *Package, ownPath string) {
	for i := range p.Packages {
		newPath := ownPath + "." + p.Packages[i].Name
		proc.renamePackage(p.Packages[i], newPath)
		if newName, ok := proc.renameExports[newPath]; ok {
			tempPkg := *p.Packages[i]
			tempPkg.MemberName = MemberName{Name: newName}
			p.Packages[i] = &tempPkg
		}
	}

	for i := range p.Modules {
		newPath := ownPath + "." + p.Modules[i].Name
		proc.renameModule(p.Modules[i], newPath)
		if newName, ok := proc.renameExports[newPath]; ok {
			tempMod := *p.Modules[i]
			tempMod.MemberName = MemberName{Name: newName}
			p.Modules[i] = &tempMod
		}
	}

	for i := range p.Aliases {
		newPath := ownPath + "." + p.Aliases[i].Name
		if newName, ok := proc.renameExports[newPath]; ok {
			tempMod := *p.Aliases[i]
			tempMod.MemberName = MemberName{Name: newName}
			p.Aliases[i] = &tempMod
		}
	}
	for i := range p.Structs {
		newPath := ownPath + "." + p.Structs[i].Name
		if newName, ok := proc.renameExports[newPath]; ok {
			//proc.renameStruct(p.Structs[i], newPath, newName)
			tempMod := *p.Structs[i]
			tempMod.MemberName = MemberName{Name: newName}
			p.Structs[i] = &tempMod
		}
	}
	for i := range p.Traits {
		newPath := ownPath + "." + p.Traits[i].Name
		if newName, ok := proc.renameExports[newPath]; ok {
			//proc.renameTrait(p.Traits[i], newPath, newName)
			tempMod := *p.Traits[i]
			tempMod.MemberName = MemberName{Name: newName}
			p.Traits[i] = &tempMod
		}
	}
	for i := range p.Functions {
		newPath := ownPath + "." + p.Functions[i].Name
		if newName, ok := proc.renameExports[newPath]; ok {
			tempMod := *p.Functions[i]
			tempMod.MemberName = MemberName{Name: newName}
			p.Functions[i] = &tempMod
		}
	}
}

func (proc *Processor) renameModule(m *Module, ownPath string) {
	for i := range m.Aliases {
		newPath := ownPath + "." + m.Aliases[i].Name
		if newName, ok := proc.renameExports[newPath]; ok {
			//proc.renameStruct(m.Structs[i], newPath, newName)
			tempMod := *m.Aliases[i]
			tempMod.MemberName = MemberName{Name: newName}
			m.Aliases[i] = &tempMod
		}
	}
	for i := range m.Structs {
		newPath := ownPath + "." + m.Structs[i].Name
		if newName, ok := proc.renameExports[newPath]; ok {
			//proc.renameTrait(m.Traits[i], newPath, newName)
			tempMod := *m.Structs[i]
			tempMod.MemberName = MemberName{Name: newName}
			m.Structs[i] = &tempMod
		}
	}
	for i := range m.Traits {
		newPath := ownPath + "." + m.Traits[i].Name
		if newName, ok := proc.renameExports[newPath]; ok {
			tempMod := *m.Traits[i]
			tempMod.MemberName = MemberName{Name: newName}
			m.Traits[i] = &tempMod
		}
	}
	for i := range m.Functions {
		newPath := ownPath + "." + m.Functions[i].Name
		if newName, ok := proc.renameExports[newPath]; ok {
			tempMod := *m.Functions[i]
			tempMod.MemberName = MemberName{Name: newName}
			m.Functions[i] = &tempMod
		}
	}
}
