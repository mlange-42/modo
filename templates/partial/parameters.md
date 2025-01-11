{{define "parameters" -}}
{{if .Parameters}}## Parameters

{{range .Parameters -}}
 - **{{.GetName}}** (`{{.Type}}`){{if .Description}}: {{.Description}}{{end}}
{{end -}}
{{end}}
{{- end}}