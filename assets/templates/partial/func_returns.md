{{define "func_returns" -}}
**Returns:**

`{{.Returns.Type}}`{{if .Returns.Doc}}: {{.Returns.Doc}}{{end}}

{{end}}