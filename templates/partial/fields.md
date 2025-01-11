{{define "fields" -}}
{{if .Fields}}## Fields

{{range .Fields -}}
 - **{{.GetName}}** (`{{.Type}}`){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}