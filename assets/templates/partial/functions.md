{{define "functions" -}}
{{if .Functions}}## Functions

{{range .Functions -}}
 - [`{{.Name}}`]({{.GetFileName}}.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}