Mojo package

# {{.GetName}}

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{if .Modules}}## Modules

{{range .Modules -}}
 - [`{{.GetName}}`]({{.GetName}}){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{if .Packages}}## Packages

{{range .Packages -}}
 - [`{{.GetName}}`]({{.GetName}}){{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
