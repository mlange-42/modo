{{define "packages" -}}
{{if .Packages}}## Packages

{{range .Packages -}}
 - [`{{.Name}}`]({{toLink .GetFileName "module"}}){{if .Summary}}: {{.Summary}}{{end}}
{{end}}
{{end}}
{{- end}}