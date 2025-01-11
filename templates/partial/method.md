{{define "method" -}}
### {{.GetName}}

{{if .Overloads}}
{{range .Overloads}}
{{template "overload" .}}
{{end}}
{{else}}
{{template "overload" .}}
{{end}}
{{- end}}