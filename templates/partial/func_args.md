{{define "func_args" -}}
{{if .Args}}**Args:**

{{range .Args -}}
 - **{{.GetName}}** (`{{.Type}}`){{if .Description}}: {{.Description}}{{end}}
{{end -}}
{{end}}
{{- end}}