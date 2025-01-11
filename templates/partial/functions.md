{{define "functions" -}}
{{if .Functions}}## Functions

{{range .Functions -}}
 - [`{{.GetName}}`]({{.GetName}}.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}