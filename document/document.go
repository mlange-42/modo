package document

import (
	"bytes"
	"encoding/json"
	"unicode"
)

const capitalFileMarker = "-"

var CaseSensitiveSystem = true

type Docs struct {
	Decl    *Package
	Version string
}

type Package struct {
	MemberKind
	MemberName
	Description string
	Summary     string
	Modules     []*Module
	Packages    []*Package
}

type Module struct {
	MemberKind
	MemberName
	Summary     string
	Description string
	Aliases     []*Alias
	Functions   []*Function
	Structs     []*Struct
	Traits      []*Trait
}

type Alias struct {
	MemberKind
	MemberName
	Description string
	Summary     string
	Value       string
	Deprecated  string
}

type Struct struct {
	MemberKind
	MemberName
	Description  string
	Summary      string
	Aliases      []*Alias
	Constraints  string
	Convention   string
	Deprecated   string
	Fields       []*Field
	Functions    []*Function
	Parameters   []*Parameter
	ParentTraits []string
	Signature    string
}

type Function struct {
	MemberKind
	MemberName
	Description          string
	Summary              string
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
	MemberKind
	MemberName
	Description string
	Summary     string
	Type        string
}

type Trait struct {
	MemberKind
	MemberName
	Description  string
	Summary      string
	Fields       []*Field
	Functions    []*Function
	ParentTraits []string
	Deprecated   string
}

type Arg struct {
	MemberKind
	MemberName
	Description string
	Convention  string
	Type        string
	PassingKind string
	Default     string
}

type Parameter struct {
	MemberKind
	MemberName
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

type Kinded interface {
	GetKind() string
}

type Named interface {
	GetName() string
	GetFileName() string
}

type MemberKind struct {
	Kind string
}

func NewKind(kind string) MemberKind {
	return MemberKind{Kind: kind}
}

func (k *MemberKind) GetKind() string {
	return k.Kind
}

type MemberName struct {
	Name string
}

func NewName(name string) MemberName {
	return MemberName{Name: name}
}

func (k *MemberName) GetName() string {
	return k.Name
}

func (k *MemberName) GetFileName() string {
	if CaseSensitiveSystem {
		return k.Name
	}
	if isCap(k.Name) {
		return k.Name + capitalFileMarker
	}
	return k.Name
}

func isCap(s string) bool {
	if len(s) == 0 {
		return false
	}
	firstRune := []rune(s)[0]
	return unicode.IsUpper(firstRune)
}
