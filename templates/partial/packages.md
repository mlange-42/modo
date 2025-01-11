{{define "packages" -}}
{{if .Packages}}## Packages

{{range .Packages -}}
 - [`{{.GetName}}`]({{.GetName}}){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}