{{define "func_raises" -}}
{{if .Raises}}**Raises:**

{{if .RaisesDoc}}{{.RaisesDoc}}

{{end -}}
{{end}}
{{- end}}