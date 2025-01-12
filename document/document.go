package document

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

const capitalFileMarker = "-"

var CaseSensitiveSystem = true

type Docs struct {
	Decl    *Package
	Version string
}

type Package struct {
	Kind
	Name
	Description string
	Summary     string
	Modules     []*Module
	Packages    []*Package
}

type Module struct {
	Kind
	Name
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

	cleanup(&docs)

	return &docs, nil
}

func cleanup(doc *Docs) {
	cleanupPackage(doc.Decl)
}

func cleanupPackage(p *Package) {
	for _, pp := range p.Packages {
		cleanupPackage(pp)
	}
	newModules := make([]*Module, 0, len(p.Modules))
	for _, m := range p.Modules {
		cleanupModule(m)
		if m.GetName() != "__init__" {
			newModules = append(newModules, m)
		}
	}
	p.Modules = newModules
}

func cleanupModule(m *Module) {
	for _, s := range m.Structs {
		if s.Signature == "" {
			s.Signature = createSignature(s)
		}
	}
}

func createSignature(s *Struct) string {
	b := strings.Builder{}
	b.WriteString("struct ")
	b.WriteString(s.GetName())

	if len(s.Parameters) == 0 {
		return b.String()
	}

	b.WriteString("[")

	prevKind := ""
	for i, par := range s.Parameters {
		written := false
		if par.PassingKind == "kw" && prevKind != par.PassingKind {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString("*")
			written = true
		}
		if prevKind == "inferred" && par.PassingKind != prevKind {
			b.WriteString(", //")
			written = true
		}
		if prevKind == "pos" && par.PassingKind != prevKind {
			b.WriteString(", /")
			written = true
		}

		if i > 0 || written {
			b.WriteString(", ")
		}

		b.WriteString(fmt.Sprintf("%s: %s", par.GetName(), par.Type))
		if len(par.Default) > 0 {
			b.WriteString(fmt.Sprintf(" = %s", par.Default))
		}

		prevKind = par.PassingKind
	}
	if prevKind == "inferred" {
		b.WriteString(", //")
	}
	if prevKind == "pos" {
		b.WriteString(", /")
	}

	b.WriteString("]")

	return b.String()
}

type Kinded interface {
	GetKind() string
}

type Named interface {
	GetName() string
	GetFileName() string
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

func (k *Name) GetFileName() string {
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
