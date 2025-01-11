{{define "traits" -}}
{{if .Traits}}## Traits

{{range .Traits -}}
 - [`{{.GetName}}`]({{.GetName}}){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}