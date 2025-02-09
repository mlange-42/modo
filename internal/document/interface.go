package document

import (
	"path"
	"unicode"
)

type missingDocs struct {
	Who  string
	What string
}

type missingStats struct {
	Total   int
	Missing int
}

// Kinded is an interface for types that have a kind.
type Kinded interface {
	GetKind() string
}

// Named is an interface for types that have a name.
type Named interface {
	GetName() string
	GetFileName() string
}

// Summarized is an interface for types that have a summary.
type Summarized interface {
	GetSummary() string
}

// Linked is an interface for types that can be linked/referenced.
type Linked interface {
	SetLink(path []string, kind string)
	GetLink() string
}

// MemberKind holds the kind of a member.
type MemberKind struct {
	Kind string
}

func newKind(kind string) MemberKind {
	return MemberKind{Kind: kind}
}

// GetKind returns the kind of the member.
func (m *MemberKind) GetKind() string {
	return m.Kind
}

// MemberName holds the name of a member.
type MemberName struct {
	Name string
}

func newName(name string) MemberName {
	return MemberName{Name: name}
}

// GetName returns the name of the member.
func (m *MemberName) GetName() string {
	return m.Name
}

// GetFileName returns the file name of the member.
func (m *MemberName) GetFileName() string {
	return toFileName(m.Name)
}

// MemberSummary holds the summary of a member.
type MemberSummary struct {
	Summary string
}

func newSummary(summary string) *MemberSummary {
	return &MemberSummary{Summary: summary}
}

// GetSummary returns the summary of the member.
func (m *MemberSummary) GetSummary() string {
	return m.Summary
}

func (m *MemberSummary) checkMissing(path string, stats *missingStats) (missing []missingDocs) {
	if m.Summary == "" {
		missing = append(missing, missingDocs{path, "description"})
		stats.Missing++
	}
	stats.Total++
	return missing
}

// MemberDescription holds the description of a member.
type MemberDescription struct {
	Description string
}

func newDescription(description string) *MemberDescription {
	return &MemberDescription{Description: description}
}

// GetDescription returns the description of the member.
func (m *MemberDescription) GetDescription() string {
	return m.Description
}

func isCap(s string) bool {
	if len(s) == 0 {
		return false
	}
	firstRune := []rune(s)[0]
	return unicode.IsUpper(firstRune)
}

// MemberLink holds the link of a member.
type MemberLink struct {
	Link string
}

// SetLink sets the link of the member.
func (m *MemberLink) SetLink(p []string, kind string) {
	if kind == "package" {
		m.Link = path.Join(path.Join(p[1:]...), "__init__.mojo")
	} else {
		m.Link = path.Join(p[1:]...) + ".mojo"
	}
}

// GetLink returns the link of the member.
func (m *MemberLink) GetLink() string {
	return m.Link
}
