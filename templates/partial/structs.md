{{define "structs" -}}
{{if .Structs}}## Structs

{{range .Structs -}}
 - [`{{.Name}}`]({{.GetFileName}}.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}