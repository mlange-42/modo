{{define "aliases" -}}
{{if .Aliases}}## Aliases

{{range .Aliases -}}
 - `{{.Name}} = {{.Value}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}