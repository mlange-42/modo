{{range .Global}}{{.}}
{{end}}

{{if .Code -}}
fn test_{{.Name}}() raises:
{{range .Code}}    {{.}}
{{end}}
{{- end}}
