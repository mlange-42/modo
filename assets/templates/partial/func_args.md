{{define "func_args" -}}
{{if .Args}}**Args:**

{{range .Args -}}
 - **{{.Name}}** (`{{.Type}}`){{if .Description}}: {{.Description}}{{end}}
{{end}}
{{end}}
{{- end}}