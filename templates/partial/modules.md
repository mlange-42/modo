{{define "modules" -}}
{{if .Modules}}## Modules

{{range .Modules -}}
 - [`{{.GetName}}`]({{.GetName}}/_index.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}