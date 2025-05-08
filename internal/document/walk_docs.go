package document

type walkFunc = func(text string, elems []string, modElems int) (string, error)
type nameFunc = func(elem Named) string

type walker struct {
	Func     walkFunc
	NameFunc nameFunc
}

func (w *walker) walkAllDocStrings(docs *Docs) error {
	return w.walkAllDocStringsPackage(docs.Decl, []string{})
}

func (w *walker) walkAllDocStringsPackage(p *Package, elems []string) error {
	newElems := appendNew(elems, w.NameFunc(p))

	var err error
	if p.Summary, err = w.Func(p.Summary, newElems, len(newElems)); err != nil {
		return err
	}
	if p.Description, err = w.Func(p.Description, newElems, len(newElems)); err != nil {
		return err
	}

	for _, pkg := range p.Packages {
		if err := w.walkAllDocStringsPackage(pkg, newElems); err != nil {
			return err
		}
	}
	for _, mod := range p.Modules {
		if err := w.walkAllDocStringsModule(mod, newElems); err != nil {
			return err
		}
	}

	for _, a := range p.Aliases {
		if err := w.walkAllDocStringsModuleAlias(a, newElems); err != nil {
			return err
		}
	}
	for _, f := range p.Functions {
		if err := w.walkAllDocStringsFunction(f, newElems); err != nil {
			return err
		}
	}
	for _, s := range p.Structs {
		if err := w.walkAllDocStringsStruct(s, newElems); err != nil {
			return err
		}
	}
	for _, tr := range p.Traits {
		if err := w.walkAllDocStringsTrait(tr, newElems); err != nil {
			return err
		}
	}
	return nil
}

func (w *walker) walkAllDocStringsModule(m *Module, elems []string) error {
	newElems := appendNew(elems, w.NameFunc(m))

	var err error
	if m.Summary, err = w.Func(m.Summary, newElems, len(newElems)); err != nil {
		return err
	}
	if m.Description, err = w.Func(m.Description, newElems, len(newElems)); err != nil {
		return err
	}

	for _, a := range m.Aliases {
		if err := w.walkAllDocStringsModuleAlias(a, newElems); err != nil {
			return err
		}
	}
	for _, f := range m.Functions {
		if err := w.walkAllDocStringsFunction(f, newElems); err != nil {
			return err
		}
	}
	for _, s := range m.Structs {
		if err := w.walkAllDocStringsStruct(s, newElems); err != nil {
			return err
		}
	}
	for _, tr := range m.Traits {
		if err := w.walkAllDocStringsTrait(tr, newElems); err != nil {
			return err
		}
	}
	return nil
}

func (w *walker) walkAllDocStringsStruct(s *Struct, elems []string) error {
	newElems := appendNew(elems, w.NameFunc(s))

	var err error
	if s.Summary, err = w.Func(s.Summary, newElems, len(elems)); err != nil {
		return err
	}
	if s.Description, err = w.Func(s.Description, newElems, len(elems)); err != nil {
		return err
	}
	if s.Deprecated, err = w.Func(s.Deprecated, newElems, len(elems)); err != nil {
		return err
	}

	for _, a := range s.Aliases {
		if a.Summary, err = w.Func(a.Summary, newElems, len(elems)); err != nil {
			return err
		}
		if a.Description, err = w.Func(a.Description, newElems, len(elems)); err != nil {
			return err
		}
		if a.Deprecated, err = w.Func(a.Deprecated, newElems, len(elems)); err != nil {
			return err
		}
	}
	for _, p := range s.Parameters {
		if p.Description, err = w.Func(p.Description, newElems, len(elems)); err != nil {
			return err
		}
	}
	for _, f := range s.Fields {
		if f.Summary, err = w.Func(f.Summary, newElems, len(elems)); err != nil {
			return err
		}
		if f.Description, err = w.Func(f.Description, newElems, len(elems)); err != nil {
			return err
		}
	}
	for _, f := range s.Functions {
		if err := w.walkAllDocStringsMethod(f, newElems); err != nil {
			return err
		}
	}

	return nil
}

func (w *walker) walkAllDocStringsTrait(tr *Trait, elems []string) error {
	newElems := appendNew(elems, w.NameFunc(tr))

	var err error
	if tr.Summary, err = w.Func(tr.Summary, newElems, len(elems)); err != nil {
		return err
	}
	if tr.Description, err = w.Func(tr.Description, newElems, len(elems)); err != nil {
		return err
	}
	if tr.Deprecated, err = w.Func(tr.Deprecated, newElems, len(elems)); err != nil {
		return err
	}

	for _, a := range tr.Aliases {
		if a.Summary, err = w.Func(a.Summary, newElems, len(elems)); err != nil {
			return err
		}
		if a.Description, err = w.Func(a.Description, newElems, len(elems)); err != nil {
			return err
		}
		if a.Deprecated, err = w.Func(a.Deprecated, newElems, len(elems)); err != nil {
			return err
		}
	}
	// TODO: add when traits support parameters
	/*for _, p := range tr.Parameters {
		p.Description, err = replaceLinks(p.Description, newElems, len(elems), lookup, t)
		if err != nil {
			return err
		}
	}*/
	for _, f := range tr.Fields {
		if f.Summary, err = w.Func(f.Summary, newElems, len(elems)); err != nil {
			return err
		}
		if f.Description, err = w.Func(f.Description, newElems, len(elems)); err != nil {
			return err
		}
	}
	for _, f := range tr.Functions {
		if err := w.walkAllDocStringsMethod(f, newElems); err != nil {
			return err
		}
	}

	return nil
}

func (w *walker) walkAllDocStringsFunction(f *Function, elems []string) error {
	newElems := appendNew(elems, w.NameFunc(f))

	var err error
	if f.Summary, err = w.Func(f.Summary, newElems, len(elems)); err != nil {
		return err
	}
	if f.Description, err = w.Func(f.Description, newElems, len(elems)); err != nil {
		return err
	}
	if f.Deprecated, err = w.Func(f.Deprecated, newElems, len(elems)); err != nil {
		return err
	}
	if f.ReturnsDoc, err = w.Func(f.ReturnsDoc, newElems, len(elems)); err != nil {
		return err
	}
	if f.RaisesDoc, err = w.Func(f.RaisesDoc, newElems, len(elems)); err != nil {
		return err
	}

	for _, a := range f.Args {
		if a.Description, err = w.Func(a.Description, newElems, len(elems)); err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		if p.Description, err = w.Func(p.Description, newElems, len(elems)); err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := w.walkAllDocStringsFunction(o, elems)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *walker) walkAllDocStringsModuleAlias(a *Alias, elems []string) error {
	newElems := appendNew(elems, w.NameFunc(a))

	var err error
	if a.Summary, err = w.Func(a.Summary, newElems, len(elems)); err != nil {
		return err
	}
	if a.Description, err = w.Func(a.Description, newElems, len(elems)); err != nil {
		return err
	}
	if a.Deprecated, err = w.Func(a.Deprecated, newElems, len(elems)); err != nil {
		return err
	}
	return nil
}

func (w *walker) walkAllDocStringsMethod(f *Function, elems []string) error {
	var err error
	if f.Summary, err = w.Func(f.Summary, elems, len(elems)-1); err != nil {
		return err
	}
	if f.Description, err = w.Func(f.Description, elems, len(elems)-1); err != nil {
		return err
	}
	if f.Deprecated, err = w.Func(f.Deprecated, elems, len(elems)-1); err != nil {
		return err
	}
	if f.ReturnsDoc, err = w.Func(f.ReturnsDoc, elems, len(elems)-1); err != nil {
		return err
	}
	if f.RaisesDoc, err = w.Func(f.RaisesDoc, elems, len(elems)-1); err != nil {
		return err
	}

	for _, a := range f.Args {
		if a.Description, err = w.Func(a.Description, elems, len(elems)-1); err != nil {
			return err
		}
	}
	for _, p := range f.Parameters {
		if p.Description, err = w.Func(p.Description, elems, len(elems)-1); err != nil {
			return err
		}
	}

	for _, o := range f.Overloads {
		err := w.walkAllDocStringsMethod(o, elems)
		if err != nil {
			return err
		}
	}

	return nil
}
