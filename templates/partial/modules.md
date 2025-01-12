{{define "modules" -}}
{{if .Modules}}## Modules

{{range .Modules -}}
 - [`{{.GetName}}`]({{.GetFileName}}/_index.md){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}