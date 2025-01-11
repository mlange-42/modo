{{define "aliases" -}}
{{if .Aliases}}## Aliases

{{range .Aliases -}}
 - `{{.GetName}} = {{.Value}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}