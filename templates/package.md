Mojo package

# {{.Name}}

{{if .Summary}}{{.Summary}}
{{end}}
{{if .Description}}{{.Description}}
{{end}}
{{if .Modules}}## Modules

{{range .Modules -}}
 - {{.Name}}: {{.Summary}}
{{end -}}
{{end}}
{{if .Packages}}## Packages

{{range .Packages -}}
 - {{.Name}}: {{.Summary}}
{{end -}}
{{end}}
