{{define "modules" -}}
{{if .Modules}}## Modules

{{range .Modules -}}
 - [`{{.Name}}`]({{toLink .GetFileName "module"}}){{if .Summary}}: {{.Summary}}{{end}}
{{end}}
{{end}}
{{- end}}