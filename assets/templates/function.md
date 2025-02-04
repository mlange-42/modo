Mojo function{{template "source_link" .}}

# `{{.Name}}`

{{if .Overloads -}}
{{range .Overloads -}}
{{template "overload" . -}}
{{end -}}
{{else -}}
{{template "overload" . -}}
{{- end}}