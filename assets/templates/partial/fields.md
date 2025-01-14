{{define "fields" -}}
{{if .Fields}}## Fields

{{range .Fields -}}
 - **{{.Name}}** (`{{.Type}}`){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}