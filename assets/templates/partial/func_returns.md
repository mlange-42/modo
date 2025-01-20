{{define "func_returns" -}}
{{if .ReturnType}}**Returns:**

`{{.ReturnType}}`{{if .ReturnsDoc}}: {{.ReturnsDoc}}{{end}}

{{end}}
{{- end}}