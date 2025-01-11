Mojo module

# {{.Name}}

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{if .Aliases}}## Aliases

{{range .Aliases -}}
 - `{{.Name}} = {{.Value}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{if .Structs}}## Structs

{{range .Structs -}}
 - `{{.Name}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{if .Traits}}## Traits

{{range .Traits -}}
 - `{{.Name}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{if .Functions}}## Functions

{{range .Functions -}}
 - `{{.Name}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
