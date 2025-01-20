{{define "aliases" -}}
{{if .Aliases}}## Aliases

{{range .Aliases -}}
 - `{{.Name}} = {{.Value}}`{{if .Summary}}: {{.Summary}}{{end}}{{if .Description}} {{.Description}}{{end}}
{{end}}
{{end}}
{{- end}}