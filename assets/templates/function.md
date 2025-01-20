Mojo function

# `{{.Name}}`

{{if .Overloads -}}
{{range .Overloads -}}
{{template "overload" . -}}
{{end -}}
{{else -}}
{{template "overload" . -}}
{{- end}}