{{define "functions" -}}
{{if .Functions}}## Functions

{{range .Functions -}}
 - [`{{.Name}}`]({{toLink .GetFileName "function"}}){{if .Summary}}: {{.Summary}}{{end}}
{{end}}
{{end}}
{{- end}}