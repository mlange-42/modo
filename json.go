package modo

import (
	"bytes"
	"encoding/json"
)

type Docs struct {
	Decl    Package
	Version string
}

type Package struct {
	Name        string
	Description string
	Summary     string
	Kind        string
	Modules     []*Module
	Packages    []*Package
}

type Module struct {
	Kind        string
	Name        string
	Summary     string
	Description string
	Aliases     []*Alias
	Functions   []*Function
	Structs     []*Struct
	Traits      []*Trait
}

type Alias struct {
	Kind        string
	Name        string
	Description string
	Summary     string
	Value       string
	Deprecated  string
}

type Struct struct {
	Kind         string
	Name         string
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
	Kind                 string
	Name                 string
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
	Kind        string
	Name        string
	Description string
	Summary     string
	Type        string
}

type Trait struct {
	Kind         string
	Name         string
	Description  string
	Summary      string
	Fields       []*Field
	Functions    []*Function
	ParentTraits []string
	Deprecated   string
}

type Arg struct {
	Kind        string
	Name        string
	Description string
	Convention  string
	Type        string
	PassingKind string
	Default     string
}

type Parameter struct {
	Kind        string
	Name        string
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
