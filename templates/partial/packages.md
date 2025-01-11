{{define "packages" -}}
{{if .Packages}}## Packages

{{range .Packages -}}
 - [`{{.GetName}}`]({{.GetName}}/_index.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}