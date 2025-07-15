{{define "aliases" -}}
{{if .Aliases}}## Aliases

{{range .Aliases -}}
 - `{{.Signature}}{{if .Type}}: {{.Type}}{{end}} = {{.Value}}`{{if .Summary}}: {{.Summary}}{{end}}{{if .Description}} {{.Description}}{{end}}
{{end}}
{{end}}
{{- end}}