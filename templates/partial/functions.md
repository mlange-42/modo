{{define "functions" -}}
{{if .Functions}}## Functions

{{range .Functions -}}
 - [`{{.GetName}}`]({{.GetFileName}}.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}