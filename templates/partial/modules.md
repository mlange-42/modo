{{define "modules" -}}
{{if .Modules}}## Modules

{{range .Modules -}}
 - [`{{.Name}}`]({{.GetFileName}}/_index.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}