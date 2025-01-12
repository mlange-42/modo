{{define "traits" -}}
{{if .Traits}}## Traits

{{range .Traits -}}
 - [`{{.GetName}}`]({{.GetFileName}}.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}