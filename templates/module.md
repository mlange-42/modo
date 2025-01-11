Mojo module

# {{.GetName}}

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{if .Aliases}}## Aliases

{{range .Aliases -}}
 - `{{.GetName}} = {{.Value}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{if .Structs}}## Structs

{{range .Structs -}}
 - `{{.GetName}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{if .Traits}}## Traits

{{range .Traits -}}
 - `{{.GetName}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
{{if .Functions}}## Functions

{{range .Functions -}}
 - `{{.GetName}}`{{if .Summary}}: {{.Summary}}{{end}}
{{end -}}
{{end}}
