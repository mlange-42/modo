{{define "fields" -}}
{{if .Fields}}## Fields

{{range .Fields -}}
 - **{{.Name}}** (`{{.Type}}`){{if .Summary}}: {{.Summary}}{{end}}{{if .Description}} {{.Description}}{{end}}
{{end}}
{{end}}
{{- end}}