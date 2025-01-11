{{define "traits" -}}
{{if .Traits}}## Traits

{{range .Traits -}}
 - [`{{.GetName}}`]({{.GetName}}.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}