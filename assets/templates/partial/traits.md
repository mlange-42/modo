{{define "traits" -}}
{{if .Traits}}## Traits

{{range .Traits -}}
 - [`{{.Name}}`]({{.GetFileName}}.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}