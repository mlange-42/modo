package document

type walkFunc = func(text string, elems []string, modElems int) (string, error)

func (proc *Processor) walkDocs(docs *Docs, fn walkFunc) error {
	return proc.walkDocsPackage(docs.Decl, []string{}, fn)
}

func (proc *Processor) walkDocsPackage(p *Package, elems []string, fn walkFunc) error {
	newElems := appendNew(elems, p.GetName())

	var err error
	p.Summary, err = fn(p.Summary, newElems, len(newElems))
	if err != nil {
		return err
	}
	p.Description, err = fn(p.Description, newElems, len(newElems))
	if err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		proc.walkDocsPackage(pkg, newElems, fn)
	}
	for _, mod := range p.Modules {
		proc.walkDocsModule(mod, newElems, fn)
	}
	// Runs on the original docs, to packages can't have structs, traits or function yet.
	return nil
}

func (proc *Processor) walkDocsModule(m *Module, elems []string, fn walkFunc) error {
	newElems := appendNew(elems, m.GetName())

	var err error
	m.Summary, err = fn(m.Summary, newElems, len(newElems))
	if err != nil {
		return err
	}
	m.Description, err = fn(m.Description, newElems, len(newElems))
	if err != nil {
		return err
	}

	for _, a := range m.Aliases {
		err := proc.walkDocsModuleAlias(a, newElems, fn)
		if err != nil {
			return err
		}
	}
	for _, f := range m.Functions {
		err := proc.walkDocsFunction(f, newElems, fn)
		if err != nil {
			return err
		}
	}
	for _, s := range m.Structs {
		err := proc.walkDocsStruct(s, newElems, fn)
		if err != nil {
			return err
		}
	}
	for _, tr := range m.Traits {
		err := proc.walkDocsTrait(tr, newElems, fn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (proc *Processor) walkDocsStruct(s *Struct, elems []string, fn walkFunc) error {
	newElems := appendNew(elems, s.GetName())

	var err error
	s.Summary, err = fn(s.Summary, newElems, len(elems))
	if err != nil {
		return err
	}
	s.Description, err = fn(s.Description, newElems, len(elems))
	if err != nil {
		return err
	}
	s.Deprecated, err = fn(s.Deprecated, newElems, len(elems))
	if err != nil {
		return err
	}

	for _, a := range s.Aliases {
		a.Summary, err = fn(a.Summary, newElems, len(elems))
		if err != nil {
			return err
		}
		a.Description, err = fn(a.Description, newElems, len(elems))
		if err != nil {
			return err
		}
		a.Deprecated, err = fn(a.Deprecated, newElems, len(elems))
		if err != nil {
			return err
		}
	}
	for _, p := range s.Parameters {
		p.Description, err = fn(p.Description, newElems, len(elems))
		if err != nil {
			return err
		}
	}
	for _, f := range s.Fields {
		f.Summary, err = fn(f.Summary, newElems, len(elems))
		if err != nil {
			return err
		}
		f.Description, err = fn(f.Description, newElems, len(elems))
		if err != nil {
			return err
		}
	}
	for _, f := range s.Functions {
		if err := proc.walkDocsMethod(f, elems, fn); err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) walkDocsTrait(tr *Trait, elems []string, fn walkFunc) error {
	newElems := appendNew(elems, tr.GetName())

	var err error
	tr.Summary, err = fn(tr.Summary, newElems, len(elems))
	if err != nil {
		return err
	}
	tr.Description, err = fn(tr.Description, newElems, len(elems))
	if err != nil {
		return err
	}
	tr.Deprecated, err = fn(tr.Deprecated, newElems, len(elems))
	if err != nil {
		return err
	}

	// TODO: add when traits support parameters
	/*for _, p := range tr.Parameters {
		p.Description, err = replaceLinks(p.Description, newElems, len(elems), lookup, t)
		if err != nil {
			return err
		}
	}*/
	for _, f := range tr.Fields {
		f.Summary, err = fn(f.Summary, newElems, len(elems))
		if err != nil {
			return err
		}
		f.Description, err = fn(f.Description, newElems, len(elems))
		if err != nil {
			return err
		}
	}
	for _, f := range tr.Functions {
		if err := proc.walkDocsMethod(f, elems, fn); err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) walkDocsFunction(f *Function, elems []string, fn walkFunc) error {
	newElems := appendNew(elems, f.GetName())

	var err error
	f.Summary, err = fn(f.Summary, newElems, len(elems))
	if err != nil {
		return err
	}
	f.Description, err = fn(f.Description, newElems, len(elems))
	if err != nil {
		return err
	}
	f.Deprecated, err = fn(f.Deprecated, newElems, len(elems))
	if err != nil {
		return err
	}
	f.ReturnsDoc, err = fn(f.ReturnsDoc, newElems, len(elems))
	if err != nil {
		return err
	}
	f.RaisesDoc, err = fn(f.RaisesDoc, newElems, len(elems))
	if err != nil {
		return err
	}

	for _, a := range f.Args {
		a.Description, err = fn(a.Description, newElems, len(elems))
		if err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		p.Description, err = fn(p.Description, newElems, len(elems))
		if err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := proc.walkDocsFunction(o, elems, fn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) walkDocsModuleAlias(a *Alias, elems []string, fn walkFunc) error {
	newElems := appendNew(elems, a.GetName())

	var err error
	a.Summary, err = fn(a.Summary, newElems, len(elems))
	if err != nil {
		return err
	}
	a.Description, err = fn(a.Description, newElems, len(elems))
	if err != nil {
		return err
	}
	a.Deprecated, err = fn(a.Deprecated, newElems, len(elems))
	if err != nil {
		return err
	}
	return nil
}

func (proc *Processor) walkDocsMethod(f *Function, elems []string, fn walkFunc) error {
	var err error
	f.Summary, err = fn(f.Summary, elems, len(elems))
	if err != nil {
		return err
	}
	f.Description, err = fn(f.Description, elems, len(elems))
	if err != nil {
		return err
	}
	f.Deprecated, err = fn(f.Deprecated, elems, len(elems))
	if err != nil {
		return err
	}
	f.ReturnsDoc, err = fn(f.ReturnsDoc, elems, len(elems))
	if err != nil {
		return err
	}
	f.RaisesDoc, err = fn(f.RaisesDoc, elems, len(elems))
	if err != nil {
		return err
	}

	for _, a := range f.Args {
		a.Description, err = fn(a.Description, elems, len(elems))
		if err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		p.Description, err = fn(p.Description, elems, len(elems))
		if err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := proc.walkDocsMethod(o, elems, fn)
		if err != nil {
			return err
		}
	}

	return nil
}
