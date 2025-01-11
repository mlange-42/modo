{{define "structs" -}}
{{if .Structs}}## Structs

{{range .Structs -}}
 - [`{{.GetName}}`]({{.GetName}}){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}