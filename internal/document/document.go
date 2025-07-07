package document

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"

	"gopkg.in/yaml.v3"
)

const capitalFileMarker = "-"

// Global variable for file case sensitivity.
//
// TODO: find another way to handle this, without using a global variable.
var caseSensitiveSystem = true

// Docs holds the document for a package.
type Docs struct {
	Decl    *Package
	Version string
}

// Package holds the document for a package.
type Package struct {
	MemberKind         `yaml:",inline"`
	MemberName         `yaml:",inline"`
	*MemberSummary     `yaml:",inline"`
	*MemberDescription `yaml:",inline"`
	Modules            []*Module
	Packages           []*Package
	Aliases            []*Alias         `yaml:",omitempty" json:",omitempty"` // Additional field for package re-exports
	Functions          []*Function      `yaml:",omitempty" json:",omitempty"` // Additional field for package re-exports
	Structs            []*Struct        `yaml:",omitempty" json:",omitempty"` // Additional field for package re-exports
	Traits             []*Trait         `yaml:",omitempty" json:",omitempty"` // Additional field for package re-exports
	exports            []*packageExport `yaml:"-" json:"-"`                   // Additional field for package re-exports
	MemberLink         `yaml:"-" json:"-"`
}

// checkMissing checks for missing documentation.
func (p *Package) checkMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := p.Name
	if len(path) > 0 {
		newPath = fmt.Sprintf("%s.%s", path, p.Name)
	}
	missing = p.MemberSummary.checkMissing(newPath, stats)
	for _, e := range p.Packages {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range p.Modules {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range p.Aliases {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range p.Structs {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range p.Traits {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range p.Functions {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	return missing
}

func (p *Package) linkedCopy() *Package {
	return &Package{
		MemberName:        newName(p.Name),
		MemberKind:        newKind(p.Kind),
		MemberSummary:     p.MemberSummary,
		MemberDescription: p.MemberDescription,
		exports:           p.exports,
		MemberLink:        p.MemberLink,
	}
}

// Module holds the document for a module.
type Module struct {
	MemberKind    `yaml:",inline"`
	MemberName    `yaml:",inline"`
	MemberSummary `yaml:",inline"`
	Description   string
	Aliases       []*Alias
	Functions     []*Function
	Structs       []*Struct
	Traits        []*Trait
	MemberLink    `yaml:"-" json:"-"`
}

func (m *Module) checkMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := fmt.Sprintf("%s.%s", path, m.Name)
	missing = m.MemberSummary.checkMissing(newPath, stats)
	for _, e := range m.Aliases {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range m.Structs {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range m.Traits {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range m.Functions {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	return missing
}

// Alias holds the document for an alias.
type Alias struct {
	MemberKind    `yaml:",inline"`
	MemberName    `yaml:",inline"`
	MemberSummary `yaml:",inline"`
	Description   string
	Value         string
	Deprecated    string
	Parameters    []*Parameter
}

func (a *Alias) checkMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := fmt.Sprintf("%s.%s", path, a.Name)
	missing = a.MemberSummary.checkMissing(newPath, stats)
	for _, e := range a.Parameters {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	return missing
}

// Struct holds the document for a struct.
type Struct struct {
	MemberKind    `yaml:",inline"`
	MemberName    `yaml:",inline"`
	MemberSummary `yaml:",inline"`
	Description   string
	Aliases       []*Alias
	Constraints   string
	Convention    string
	Deprecated    string
	Fields        []*Field
	Functions     []*Function
	Parameters    []*Parameter
	ParentTraits  []string
	Signature     string
	MemberLink    `yaml:"-" json:"-"`
}

func (s *Struct) checkMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := fmt.Sprintf("%s.%s", path, s.Name)
	missing = s.MemberSummary.checkMissing(newPath, stats)
	for _, e := range s.Aliases {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range s.Fields {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range s.Parameters {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range s.Functions {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	return missing
}

// Function holds the document for a function.
type Function struct {
	MemberKind           `yaml:",inline"`
	MemberName           `yaml:",inline"`
	MemberSummary        `yaml:",inline"`
	Description          string
	Args                 []*Arg
	Overloads            []*Function
	Async                bool
	Constraints          string
	Deprecated           string
	IsDef                bool
	IsStatic             bool
	IsImplicitConversion bool
	Raises               bool
	RaisesDoc            string
	ReturnType           string
	ReturnsDoc           string
	Signature            string
	Parameters           []*Parameter
	MemberLink           `yaml:"-" json:"-"`
}

func (f *Function) checkMissing(path string, stats *missingStats) (missing []missingDocs) {
	if len(f.Overloads) == 0 {
		newPath := fmt.Sprintf("%s.%s", path, f.Name)
		missing = f.MemberSummary.checkMissing(newPath, stats)
		if f.Raises && f.RaisesDoc == "" {
			missing = append(missing, missingDocs{newPath, "raises docs"})
			stats.Missing++
		}
		stats.Total++

		if !slices.Contains(initializers[:], f.Name) {
			if f.ReturnType != "" && f.ReturnsDoc == "" {
				missing = append(missing, missingDocs{newPath, "return docs"})
				stats.Missing++
			}
			stats.Total++
		}

		for _, e := range f.Parameters {
			missing = append(missing, e.checkMissing(newPath, stats)...)
		}
		for _, e := range f.Args {
			missing = append(missing, e.checkMissing(newPath, stats)...)
		}
		return missing
	}
	for _, o := range f.Overloads {
		missing = append(missing, o.checkMissing(path, stats)...)
	}
	return missing
}

// Field holds the document for a field.
type Field struct {
	MemberKind    `yaml:",inline"`
	MemberName    `yaml:",inline"`
	MemberSummary `yaml:",inline"`
	Description   string
	Type          string
}

func (f *Field) checkMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := fmt.Sprintf("%s.%s", path, f.Name)
	return f.MemberSummary.checkMissing(newPath, stats)
}

// Trait holds the document for a trait.
type Trait struct {
	MemberKind    `yaml:",inline"`
	MemberName    `yaml:",inline"`
	MemberSummary `yaml:",inline"`
	Description   string
	Aliases       []*Alias
	Fields        []*Field
	Functions     []*Function
	ParentTraits  []string
	Deprecated    string
	MemberLink    `yaml:"-" json:"-"`
}

func (t *Trait) checkMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := fmt.Sprintf("%s.%s", path, t.Name)
	missing = t.MemberSummary.checkMissing(newPath, stats)
	for _, e := range t.Aliases {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range t.Fields {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	for _, e := range t.Functions {
		missing = append(missing, e.checkMissing(newPath, stats)...)
	}
	return missing
}

// Arg holds the document for a function argument.
type Arg struct {
	MemberKind  `yaml:",inline"`
	MemberName  `yaml:",inline"`
	Description string
	Convention  string
	Type        string
	PassingKind string
	Default     string
}

func (a *Arg) checkMissing(path string, stats *missingStats) (missing []missingDocs) {
	if a.Name == "self" {
		return nil
	}
	if a.Convention == "out" {
		return nil
	}
	if a.Description == "" {
		missing = append(missing, missingDocs{fmt.Sprintf("%s.%s", path, a.Name), "description"})
		stats.Missing++
	}
	stats.Total++
	return missing
}

// Parameter holds the document for a parameter.
type Parameter struct {
	MemberKind  `yaml:",inline"`
	MemberName  `yaml:",inline"`
	Description string
	Type        string
	PassingKind string
	Default     string
}

func (p *Parameter) checkMissing(path string, stats *missingStats) (missing []missingDocs) {
	if p.Description == "" {
		missing = append(missing, missingDocs{fmt.Sprintf("%s.%s", path, p.Name), "description"})
		stats.Missing++
	}
	stats.Total++
	return missing
}

// FromJSON parses JSON documentation.
func FromJSON(data []byte) (*Docs, error) {
	reader := bytes.NewReader(data)
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()

	var docs Docs

	if err := dec.Decode(&docs); err != nil {
		return nil, err
	}

	cleanup(&docs)

	return &docs, nil
}

// ToJSON converts the documentation to JSON.
func (d *Docs) ToJSON() ([]byte, error) {
	b := bytes.Buffer{}
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")

	if err := enc.Encode(d); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// FromYAML parses YAML documentation.
func FromYAML(data []byte) (*Docs, error) {
	reader := bytes.NewReader(data)
	dec := yaml.NewDecoder(reader)
	dec.KnownFields(true)

	var docs Docs

	if err := dec.Decode(&docs); err != nil {
		return nil, err
	}

	cleanup(&docs)

	return &docs, nil
}

// ToYAML converts the documentation to YAML.
func (d *Docs) ToYAML() ([]byte, error) {
	b := bytes.Buffer{}
	enc := yaml.NewEncoder(&b)

	if err := enc.Encode(d); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
