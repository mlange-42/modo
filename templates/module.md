Mojo module

# {{.Name}}

{{if .Summary}}{{.Summary}}
{{end}}
{{if .Description}}{{.Description}}
{{end}}
{{if .Structs}}## Structs

{{range .Structs -}}
 - {{.Name}}: {{.Summary}}
{{end -}}
{{end}}
{{if .Traits}}## Traits

{{range .Traits -}}
 - {{.Name}}: {{.Summary}}
{{end -}}
{{end}}
{{if .Functions}}## Functions

{{range .Functions -}}
 - {{.Name}}: {{.Summary}}
{{end -}}
{{end}}
