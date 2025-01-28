package document

import (
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

//type missingChecker interface {
//	CheckMissing(path string) (missing []missingDocs)
//}

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

func (m *MemberKind) GetKind() string {
	return m.Kind
}

type MemberName struct {
	Name string
}

func newName(name string) MemberName {
	return MemberName{Name: name}
}

func (m *MemberName) GetName() string {
	return m.Name
}

func (m *MemberName) GetFileName() string {
	if caseSensitiveSystem {
		return m.Name
	}
	if isCap(m.Name) {
		return m.Name + capitalFileMarker
	}
	return m.Name
}

type MemberSummary struct {
	Summary string
}

func newSummary(summary string) *MemberSummary {
	return &MemberSummary{Summary: summary}
}

func (m *MemberSummary) GetSummary() string {
	return m.Summary
}

func (m *MemberSummary) CheckMissing(path string, stats *missingStats) (missing []missingDocs) {
	if m.Summary == "" {
		missing = append(missing, missingDocs{path, "description"})
		stats.Missing++
	}
	stats.Total++
	return missing
}

type MemberDescription struct {
	Description string
}

func newDescription(description string) *MemberDescription {
	return &MemberDescription{Description: description}
}

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
