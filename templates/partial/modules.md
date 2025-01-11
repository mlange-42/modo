{{define "modules" -}}
{{if .Modules}}## Modules

{{range .Modules -}}
 - [`{{.GetName}}`]({{.GetName}}){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{- end}}