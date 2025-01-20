package document

import (
	"bytes"
	"encoding/json"
	"unicode"

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

type Alias struct {
	MemberKind    `yaml:",inline"`
	MemberName    `yaml:",inline"`
	MemberSummary `yaml:",inline"`
	Description   string
	Value         string
	Deprecated    string
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

type Field struct {
	MemberKind    `yaml:",inline"`
	MemberName    `yaml:",inline"`
	MemberSummary `yaml:",inline"`
	Description   string
	Type          string
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

type Arg struct {
	MemberKind  `yaml:",inline"`
	MemberName  `yaml:",inline"`
	Description string
	Convention  string
	Type        string
	PassingKind string
	Default     string
}

type Parameter struct {
	MemberKind  `yaml:",inline"`
	MemberName  `yaml:",inline"`
	Description string
	Type        string
	PassingKind string
	Default     string
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

type Kinded interface {
	GetKind() string
}

type Named interface {
	GetName() string
	GetFileName() string
}

type Summarized interface {
	GetSummary() string
}

type MemberKind struct {
	Kind string
}

func newKind(kind string) MemberKind {
	return MemberKind{Kind: kind}
}

func (k *MemberKind) GetKind() string {
	return k.Kind
}

type MemberName struct {
	Name string
}

func newName(name string) MemberName {
	return MemberName{Name: name}
}

func (k *MemberName) GetName() string {
	return k.Name
}

func (k *MemberName) GetFileName() string {
	if caseSensitiveSystem {
		return k.Name
	}
	if isCap(k.Name) {
		return k.Name + capitalFileMarker
	}
	return k.Name
}

type MemberSummary struct {
	Summary string
}

func newSummary(summary string) *MemberSummary {
	return &MemberSummary{Summary: summary}
}

func (k *MemberSummary) GetSummary() string {
	return k.Summary
}

type MemberDescription struct {
	Description string
}

func newDescription(description string) *MemberDescription {
	return &MemberDescription{Description: description}
}

func (k *MemberDescription) GetDescription() string {
	return k.Description
}

func isCap(s string) bool {
	if len(s) == 0 {
		return false
	}
	firstRune := []rune(s)[0]
	return unicode.IsUpper(firstRune)
}
