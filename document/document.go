package document

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

const capitalFileMarker = "-"

// Global variable for file case sensitivity.
//
// TODO: find another way to handle this, without using a global variable.
var caseSensitiveSystem = true

type Docs struct {
	Decl    *Package
	Version string
}

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
}

func (p *Package) CheckMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := p.Name
	if len(path) > 0 {
		newPath = fmt.Sprintf("%s.%s", path, p.Name)
	}
	missing = p.MemberSummary.CheckMissing(newPath, stats)
	for _, e := range p.Packages {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range p.Modules {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range p.Aliases {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range p.Structs {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range p.Traits {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range p.Functions {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
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
	}
}

type Module struct {
	MemberKind    `yaml:",inline"`
	MemberName    `yaml:",inline"`
	MemberSummary `yaml:",inline"`
	Description   string
	Aliases       []*Alias
	Functions     []*Function
	Structs       []*Struct
	Traits        []*Trait
}

func (m *Module) CheckMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := fmt.Sprintf("%s.%s", path, m.Name)
	missing = m.MemberSummary.CheckMissing(newPath, stats)
	for _, e := range m.Aliases {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range m.Structs {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range m.Traits {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range m.Functions {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	return missing
}

type Alias struct {
	MemberKind    `yaml:",inline"`
	MemberName    `yaml:",inline"`
	MemberSummary `yaml:",inline"`
	Description   string
	Value         string
	Deprecated    string
}

func (a *Alias) CheckMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := fmt.Sprintf("%s.%s", path, a.Name)
	return a.MemberSummary.CheckMissing(newPath, stats)
}

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
}

func (s *Struct) CheckMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := fmt.Sprintf("%s.%s", path, s.Name)
	missing = s.MemberSummary.CheckMissing(newPath, stats)
	for _, e := range s.Aliases {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range s.Fields {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range s.Parameters {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range s.Functions {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	return missing
}

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
}

func (f *Function) CheckMissing(path string, stats *missingStats) (missing []missingDocs) {
	if len(f.Overloads) == 0 {
		newPath := fmt.Sprintf("%s.%s", path, f.Name)
		missing = f.MemberSummary.CheckMissing(newPath, stats)
		if f.Raises && f.RaisesDoc == "" {
			missing = append(missing, missingDocs{newPath, "raises docs"})
			stats.Missing++
		}
		if f.ReturnType != "" && f.ReturnsDoc == "" {
			missing = append(missing, missingDocs{newPath, "return docs"})
			stats.Missing++
		}
		stats.Total += 2

		for _, e := range f.Parameters {
			missing = append(missing, e.CheckMissing(newPath, stats)...)
		}
		for _, e := range f.Args {
			missing = append(missing, e.CheckMissing(newPath, stats)...)
		}
		return missing
	}
	for _, o := range f.Overloads {
		missing = append(missing, o.CheckMissing(path, stats)...)
	}
	return missing
}

type Field struct {
	MemberKind    `yaml:",inline"`
	MemberName    `yaml:",inline"`
	MemberSummary `yaml:",inline"`
	Description   string
	Type          string
}

func (f *Field) CheckMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := fmt.Sprintf("%s.%s", path, f.Name)
	return f.MemberSummary.CheckMissing(newPath, stats)
}

type Trait struct {
	MemberKind    `yaml:",inline"`
	MemberName    `yaml:",inline"`
	MemberSummary `yaml:",inline"`
	Description   string
	Fields        []*Field
	Functions     []*Function
	ParentTraits  []string
	Deprecated    string
}

func (t *Trait) CheckMissing(path string, stats *missingStats) (missing []missingDocs) {
	newPath := fmt.Sprintf("%s.%s", path, t.Name)
	missing = t.MemberSummary.CheckMissing(newPath, stats)
	for _, e := range t.Fields {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	for _, e := range t.Functions {
		missing = append(missing, e.CheckMissing(newPath, stats)...)
	}
	return missing
}

type Arg struct {
	MemberKind  `yaml:",inline"`
	MemberName  `yaml:",inline"`
	Description string
	Convention  string
	Type        string
	PassingKind string
	Default     string
}

func (a *Arg) CheckMissing(path string, stats *missingStats) (missing []missingDocs) {
	if a.Name == "self" && a.Type == "Self" {
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

type Parameter struct {
	MemberKind  `yaml:",inline"`
	MemberName  `yaml:",inline"`
	Description string
	Type        string
	PassingKind string
	Default     string
}

func (p *Parameter) CheckMissing(path string, stats *missingStats) (missing []missingDocs) {
	if p.Description == "" {
		missing = append(missing, missingDocs{fmt.Sprintf("%s.%s", path, p.Name), "description"})
		stats.Missing++
	}
	stats.Total++
	return missing
}

func FromJson(data []byte) (*Docs, error) {
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

func (d *Docs) ToJson() ([]byte, error) {
	b := bytes.Buffer{}
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")

	if err := enc.Encode(d); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func FromYaml(data []byte) (*Docs, error) {
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

func (d *Docs) ToYaml() ([]byte, error) {
	b := bytes.Buffer{}
	enc := yaml.NewEncoder(&b)

	if err := enc.Encode(d); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
