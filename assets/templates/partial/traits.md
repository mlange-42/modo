{{define "traits" -}}
{{if .Traits}}## Traits

{{range .Traits -}}
 - [`{{.Name}}`]({{toLink .GetFileName "trait"}}){{if .Summary}}: {{.Summary}}{{end}}
{{end}}
{{end}}
{{- end}}