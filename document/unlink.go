package document

func (proc *Processor) unlinkData() {
	proc.ExportDocs.Decl = cp(proc.ExportDocs.Decl)
	proc.unlinkDataPackage(proc.ExportDocs.Decl)
}

func (proc *Processor) unlinkDataPackage(p *Package) {
	p.MemberSummary = cp(p.MemberSummary)
	p.MemberDescription = cp(p.MemberDescription)

	p.Packages = cpSlice(p.Packages)
	for i := range p.Packages {
		proc.unlinkDataPackage(p.Packages[i])
	}

	p.Modules = cpSlice(p.Modules)
	for i := range p.Modules {
		proc.unlinkDataModule(p.Modules[i])
	}

	p.Structs = cpSlice(p.Structs)
	for i := range p.Structs {
		proc.unlinkDataStruct(p.Structs[i])
	}

	p.Traits = cpSlice(p.Traits)
	for i := range p.Traits {
		proc.unlinkDataTrait(p.Traits[i])
	}

	p.Functions = cpSlice(p.Functions)
	for i := range p.Functions {
		proc.unlinkDataFunction(p.Functions[i])
	}
}

func (proc *Processor) unlinkDataModule(m *Module) {
	m.Aliases = cpSlice(m.Aliases)

	m.Structs = cpSlice(m.Structs)
	for i := range m.Structs {
		proc.unlinkDataStruct(m.Structs[i])
	}

	m.Traits = cpSlice(m.Traits)
	for i := range m.Traits {
		proc.unlinkDataTrait(m.Traits[i])
	}

	m.Functions = cpSlice(m.Functions)
	for i := range m.Functions {
		proc.unlinkDataFunction(m.Functions[i])
	}
}

func (proc *Processor) unlinkDataStruct(s *Struct) {
	s.Aliases = cpSlice(s.Aliases)
	s.Parameters = cpSlice(s.Parameters)
	s.Fields = cpSlice(s.Fields)

	s.Functions = cpSlice(s.Functions)
	for i := range s.Functions {
		proc.unlinkDataFunction(s.Functions[i])
	}
}

func (proc *Processor) unlinkDataTrait(t *Trait) {
	//t.Parameters = cpSlicePtr(t.Parameters)
	t.Fields = cpSlice(t.Fields)

	t.Functions = cpSlice(t.Functions)
	for i := range t.Functions {
		proc.unlinkDataFunction(t.Functions[i])
	}
}

func (proc *Processor) unlinkDataFunction(f *Function) {
	for i := range f.Parameters {
		f.Parameters[i] = cp(f.Parameters[i])
	}
	for i := range f.Args {
		f.Args[i] = cp(f.Args[i])
	}

	for i := range f.Overloads {
		f.Overloads[i] = cp(f.Overloads[i])
		proc.unlinkDataFunction(f.Overloads[i])
	}
}

func cp[T any](v *T) *T {
	cp := *v
	return &cp
}

func cpSlice[T any](sl []*T) []*T {
	newSl := make([]*T, len(sl))
	for i, v := range sl {
		newSl[i] = cp(v)
	}
	return newSl
}
