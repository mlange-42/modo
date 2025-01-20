{{define "methods" -}}
{{if .Functions}}## Methods

{{range .Functions -}}
{{template "method" . -}}
{{end}}
{{end}}
{{- end}}