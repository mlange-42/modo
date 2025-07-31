{{define "func_returns_old" -}}
{{if .ReturnType}}**Returns:**

`{{.ReturnType}}`{{if .ReturnsDoc}}: {{.ReturnsDoc}}{{end}}

{{end}}
{{- end}}