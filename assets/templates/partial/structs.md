{{define "structs" -}}
{{if .Structs}}## Structs

{{range .Structs -}}
 - [`{{.Name}}`]({{toLink .GetFileName "struct"}}){{if .Summary}}: {{.Summary}}{{end}}
{{end}}
{{end}}
{{- end}}