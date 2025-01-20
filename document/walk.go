package document

type walkFunc = func(text string, elems []string, modElems int) (string, error)
type nameFunc = func(elem Named) string

func (proc *Processor) walkDocs(docs *Docs, fn walkFunc, nameFn nameFunc) error {
	return proc.walkDocsPackage(docs.Decl, []string{}, fn, nameFn)
}

func (proc *Processor) walkDocsPackage(p *Package, elems []string, fn walkFunc, nameFn nameFunc) error {
	newElems := appendNew(elems, nameFn(p))

	var err error
	if p.Summary, err = fn(p.Summary, newElems, len(newElems)); err != nil {
		return err
	}
	if p.Description, err = fn(p.Description, newElems, len(newElems)); err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		if err := proc.walkDocsPackage(pkg, newElems, fn, nameFn); err != nil {
			return err
		}
	}
	for _, mod := range p.Modules {
		if err := proc.walkDocsModule(mod, newElems, fn, nameFn); err != nil {
			return err
		}
	}

	for _, a := range p.Aliases {
		if err := proc.walkDocsModuleAlias(a, newElems, fn, nameFn); err != nil {
			return err
		}
	}
	for _, f := range p.Functions {
		if err := proc.walkDocsFunction(f, newElems, fn, nameFn); err != nil {
			return err
		}
	}
	for _, s := range p.Structs {
		if err := proc.walkDocsStruct(s, newElems, fn, nameFn); err != nil {
			return err
		}
	}
	for _, tr := range p.Traits {
		if err := proc.walkDocsTrait(tr, newElems, fn, nameFn); err != nil {
			return err
		}
	}
	return nil
}

func (proc *Processor) walkDocsModule(m *Module, elems []string, fn walkFunc, nameFn nameFunc) error {
	newElems := appendNew(elems, nameFn(m))

	var err error
	if m.Summary, err = fn(m.Summary, newElems, len(newElems)); err != nil {
		return err
	}
	if m.Description, err = fn(m.Description, newElems, len(newElems)); err != nil {
		return err
	}

	for _, a := range m.Aliases {
		if err := proc.walkDocsModuleAlias(a, newElems, fn, nameFn); err != nil {
			return err
		}
	}
	for _, f := range m.Functions {
		if err := proc.walkDocsFunction(f, newElems, fn, nameFn); err != nil {
			return err
		}
	}
	for _, s := range m.Structs {
		if err := proc.walkDocsStruct(s, newElems, fn, nameFn); err != nil {
			return err
		}
	}
	for _, tr := range m.Traits {
		if err := proc.walkDocsTrait(tr, newElems, fn, nameFn); err != nil {
			return err
		}
	}
	return nil
}

func (proc *Processor) walkDocsStruct(s *Struct, elems []string, fn walkFunc, nameFn nameFunc) error {
	newElems := appendNew(elems, nameFn(s))

	var err error
	if s.Summary, err = fn(s.Summary, newElems, len(elems)); err != nil {
		return err
	}
	if s.Description, err = fn(s.Description, newElems, len(elems)); err != nil {
		return err
	}
	if s.Deprecated, err = fn(s.Deprecated, newElems, len(elems)); err != nil {
		return err
	}

	for _, a := range s.Aliases {
		if a.Summary, err = fn(a.Summary, newElems, len(elems)); err != nil {
			return err
		}
		if a.Description, err = fn(a.Description, newElems, len(elems)); err != nil {
			return err
		}
		if a.Deprecated, err = fn(a.Deprecated, newElems, len(elems)); err != nil {
			return err
		}
	}
	for _, p := range s.Parameters {
		if p.Description, err = fn(p.Description, newElems, len(elems)); err != nil {
			return err
		}
	}
	for _, f := range s.Fields {
		if f.Summary, err = fn(f.Summary, newElems, len(elems)); err != nil {
			return err
		}
		if f.Description, err = fn(f.Description, newElems, len(elems)); err != nil {
			return err
		}
	}
	for _, f := range s.Functions {
		if err := proc.walkDocsMethod(f, newElems, fn, nameFn); err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) walkDocsTrait(tr *Trait, elems []string, fn walkFunc, nameFn nameFunc) error {
	newElems := appendNew(elems, nameFn(tr))

	var err error
	if tr.Summary, err = fn(tr.Summary, newElems, len(elems)); err != nil {
		return err
	}
	if tr.Description, err = fn(tr.Description, newElems, len(elems)); err != nil {
		return err
	}
	if tr.Deprecated, err = fn(tr.Deprecated, newElems, len(elems)); err != nil {
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
		if f.Summary, err = fn(f.Summary, newElems, len(elems)); err != nil {
			return err
		}
		if f.Description, err = fn(f.Description, newElems, len(elems)); err != nil {
			return err
		}
	}
	for _, f := range tr.Functions {
		if err := proc.walkDocsMethod(f, newElems, fn, nameFn); err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) walkDocsFunction(f *Function, elems []string, fn walkFunc, nameFn nameFunc) error {
	newElems := appendNew(elems, nameFn(f))

	var err error
	if f.Summary, err = fn(f.Summary, newElems, len(elems)); err != nil {
		return err
	}
	if f.Description, err = fn(f.Description, newElems, len(elems)); err != nil {
		return err
	}
	if f.Deprecated, err = fn(f.Deprecated, newElems, len(elems)); err != nil {
		return err
	}
	if f.ReturnsDoc, err = fn(f.ReturnsDoc, newElems, len(elems)); err != nil {
		return err
	}
	if f.RaisesDoc, err = fn(f.RaisesDoc, newElems, len(elems)); err != nil {
		return err
	}

	for _, a := range f.Args {
		if a.Description, err = fn(a.Description, newElems, len(elems)); err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		if p.Description, err = fn(p.Description, newElems, len(elems)); err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := proc.walkDocsFunction(o, elems, fn, nameFn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (proc *Processor) walkDocsModuleAlias(a *Alias, elems []string, fn walkFunc, nameFn nameFunc) error {
	newElems := appendNew(elems, nameFn(a))

	var err error
	if a.Summary, err = fn(a.Summary, newElems, len(elems)); err != nil {
		return err
	}
	if a.Description, err = fn(a.Description, newElems, len(elems)); err != nil {
		return err
	}
	if a.Deprecated, err = fn(a.Deprecated, newElems, len(elems)); err != nil {
		return err
	}
	return nil
}

func (proc *Processor) walkDocsMethod(f *Function, elems []string, fn walkFunc, nameFn nameFunc) error {
	var err error
	if f.Summary, err = fn(f.Summary, elems, len(elems)-1); err != nil {
		return err
	}
	if f.Description, err = fn(f.Description, elems, len(elems)-1); err != nil {
		return err
	}
	if f.Deprecated, err = fn(f.Deprecated, elems, len(elems)-1); err != nil {
		return err
	}
	if f.ReturnsDoc, err = fn(f.ReturnsDoc, elems, len(elems)-1); err != nil {
		return err
	}
	if f.RaisesDoc, err = fn(f.RaisesDoc, elems, len(elems)-1); err != nil {
		return err
	}

	for _, a := range f.Args {
		if a.Description, err = fn(a.Description, elems, len(elems)-1); err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		if p.Description, err = fn(p.Description, elems, len(elems)-1); err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := proc.walkDocsMethod(o, elems, fn, nameFn)
		if err != nil {
			return err
		}
	}

	return nil
}
