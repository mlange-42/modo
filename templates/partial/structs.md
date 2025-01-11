{{define "structs" -}}
{{if .Structs}}## Structs

{{range .Structs -}}
 - [`{{.GetName}}`]({{.GetName}}.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}