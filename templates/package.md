Mojo package

# {{.Name}}

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{if .Modules}}## Modules

{{range .Modules -}}
 - `{{.Name}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{if .Packages}}## Packages

{{range .Packages -}}
 - `{{.Name}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
