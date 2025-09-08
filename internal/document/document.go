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
	Type          string
	Path          string
	Value         string
	Deprecated    string
	Signature     string
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
	MemberKind         `yaml:",inline"`
	MemberName         `yaml:",inline"`
	MemberSummary      `yaml:",inline"`
	Description        string
	Aliases            []*Alias
	Constraints        string
	Convention         string
	Deprecated         string
	Fields             []*Field
	Functions          []*Function
	Parameters         []*Parameter
	ParentTraits       []*ParentTrait `yaml:"-" json:"-"`                       // remove tag on next stable release of Mojo.
	ParentTraitsHelper ParentTraits   `yaml:"parentTraits" json:"parentTraits"` // remove on next stable release of Mojo.
	Signature          string
	MemberLink         `yaml:"-" json:"-"`
}

// Custom logic after deserialization
//
// TODO: remove on next stable release of Mojo.
func (o *Struct) UnmarshalJSON(data []byte) error {
	type Alias Struct // Avoid recursion
	var temp Alias
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.ParentTraitsHelper.IsStruct {
		temp.ParentTraits = temp.ParentTraitsHelper.Structs
	} else {
		for _, name := range temp.ParentTraitsHelper.Strings {
			temp.ParentTraits = append(temp.ParentTraits, &ParentTrait{Name: name})
		}
	}

	*o = Struct(temp)
	return nil
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
	MemberKind               `yaml:",inline"`
	MemberName               `yaml:",inline"`
	MemberSummary            `yaml:",inline"`
	Description              string
	Args                     []*Arg
	Overloads                []*Function
	Async                    bool
	Constraints              string
	Deprecated               string
	IsDef                    bool
	IsStatic                 bool
	IsImplicitConversion     bool
	Raises                   bool
	RaisesDoc                string
	ReturnType               string // TODO: remove
	ReturnsDoc               string // TODO: remove
	Returns                  *Returns
	Signature                string
	Parameters               []*Parameter
	HasDefaultImplementation bool
	MemberLink               `yaml:"-" json:"-"`
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
			if f.Returns == nil { // old version
				if f.ReturnType != "" && f.ReturnsDoc == "" {
					missing = append(missing, missingDocs{newPath, "return docs"})
					stats.Missing++
				}
			} else { // new version
				if f.Returns.Doc == "" {
					missing = append(missing, missingDocs{newPath, "return docs"})
					stats.Missing++
				}
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

// Returns holds information on function return type and docs
type Returns struct {
	Type string
	Doc  string
	Path string
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
	MemberKind         `yaml:",inline"`
	MemberName         `yaml:",inline"`
	MemberSummary      `yaml:",inline"`
	Description        string
	Aliases            []*Alias
	Fields             []*Field
	Functions          []*Function
	ParentTraits       []*ParentTrait `yaml:"-" json:"-"`                       // remove tag on next stable release of Mojo.
	ParentTraitsHelper ParentTraits   `yaml:"parentTraits" json:"parentTraits"` // remove on next stable release of Mojo.
	Deprecated         string
	MemberLink         `yaml:"-" json:"-"`
}

// Custom logic after deserialization
//
// TODO: remove on next stable release of Mojo.
func (o *Trait) UnmarshalJSON(data []byte) error {
	type Alias Trait // Avoid recursion
	var temp Alias
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.ParentTraitsHelper.IsStruct {
		temp.ParentTraits = temp.ParentTraitsHelper.Structs
	} else {
		for _, name := range temp.ParentTraitsHelper.Strings {
			temp.ParentTraits = append(temp.ParentTraits, &ParentTrait{Name: name})
		}
	}

	*o = Trait(temp)
	return nil
}

// ParentTrait holds name and path information for a parent trait.
type ParentTrait struct {
	Name string
	Path string
}

// ParentTraits is a temporal wrapper to handle different versions of parent traits in JSON.
//
// TODO: remove on next stable release of Mojo.
type ParentTraits struct {
	Strings  []string
	Structs  []*ParentTrait
	IsStruct bool
}

func (s *ParentTraits) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as []string
	var strSlice []string
	if err := json.Unmarshal(data, &strSlice); err == nil {
		s.Strings = strSlice
		s.IsStruct = false
		return nil
	}

	// Try to unmarshal as []MyStruct
	var structSlice []*ParentTrait
	if err := json.Unmarshal(data, &structSlice); err == nil {
		s.Structs = structSlice
		s.IsStruct = true
		return nil
	}

	return fmt.Errorf("stat field is neither []string nor []ParentTrait")
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
	Path        string
	Traits      []*TraitEntry
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
	Traits      []*TraitEntry
	Path        string
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

// TraitEntry holds type and path information for a trait of a parameter or arg.
type TraitEntry struct {
	Type string
	Path string
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
