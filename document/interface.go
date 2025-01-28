package document

import "unicode"

type missingChecker interface {
	CheckMissing() (missing []string)
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

func (m *MemberSummary) CheckMissing() (missing []string) {
	if m.Summary == "" {
		missing = append(missing, "summary")
	}
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
