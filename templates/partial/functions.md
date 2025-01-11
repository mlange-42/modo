{{define "functions" -}}
{{if .Functions}}## Functions

{{range .Functions -}}
 - [`{{.GetName}}`]({{.GetName}}){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}