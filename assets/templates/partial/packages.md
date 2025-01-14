{{define "packages" -}}
{{if .Packages}}## Packages

{{range .Packages -}}
 - [`{{.Name}}`]({{.GetFileName}}/_index.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}