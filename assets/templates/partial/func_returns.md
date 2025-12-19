{{define "func_returns" -}}
{{if .Returns -}}
**Returns:**

`{{.Returns.Type}}`{{if .Returns.Doc}}: {{.Returns.Doc}}{{end}}

{{end}}
{{- end}}