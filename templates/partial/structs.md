{{define "structs" -}}
{{if .Structs}}## Structs

{{range .Structs -}}
 - [`{{.GetName}}`]({{.GetFileName}}.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}