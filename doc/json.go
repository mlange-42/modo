package doc

import (
	"bytes"
	"encoding/json"
)

type Docs struct {
	Decl    Package
	Version string
}

type Package struct {
	Kind
	Name
	Path
	Description string
	Summary     string
	Modules     []*Module
	Packages    []*Package
}

type Module struct {
	Kind
	Name
	Path
	Summary     string
	Description string
	Aliases     []*Alias
	Functions   []*Function
	Structs     []*Struct
	Traits      []*Trait
}

type Alias struct {
	Kind
	Name
	Description string
	Summary     string
	Value       string
	Deprecated  string
}

type Struct struct {
	Kind
	Name
	Path
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
	Kind
	Name
	Path
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
	Kind
	Name
	Description string
	Summary     string
	Type        string
}

type Trait struct {
	Kind
	Name
	Path
	Description  string
	Summary      string
	Fields       []*Field
	Functions    []*Function
	ParentTraits []string
	Deprecated   string
}

type Arg struct {
	Kind
	Name
	Description string
	Convention  string
	Type        string
	PassingKind string
	Default     string
}

type Parameter struct {
	Kind
	Name
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

	return &docs, nil
}

type Kinded interface {
	GetKind() string
}

type Named interface {
	GetName() string
}

type Pathed interface {
	GetPath() string
	SetPath(p string)
}

type Kind struct {
	Kind string
}

func NewKind(kind string) Kind {
	return Kind{Kind: kind}
}

func (k *Kind) GetKind() string {
	return k.Kind
}

type Name struct {
	Name string
}

func NewName(name string) Name {
	return Name{Name: name}
}

func (k *Name) GetName() string {
	return k.Name
}

type Path struct {
	Path string
}

func (p *Path) GetPath() string {
	return p.Path
}

func (p *Path) SetPath(path string) {
	p.Path = path
}
