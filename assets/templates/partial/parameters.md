{{define "parameters" -}}
{{if .Parameters}}## Parameters

{{range .Parameters -}}
 - **{{.Name}}** (`{{.Type}}`){{if .Description}}: {{.Description}}{{end}}
{{end}}
{{end}}
{{- end}}