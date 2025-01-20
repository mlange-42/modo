Mojo package

# `{{.Name}}`

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{template "aliases" .}}
{{template "structs" .}}
{{template "traits" .}}
{{template "functions" .}}
{{template "modules" .}}
{{template "packages" .}}
