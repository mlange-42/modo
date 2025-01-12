{{define "packages" -}}
{{if .Packages}}## Packages

{{range .Packages -}}
 - [`{{.GetName}}`]({{.GetFileName}}/_index.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}