Mojo package

# `{{.Name}}`

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{template "structs" .}}
{{template "traits" .}}
{{template "functions" .}}
{{template "modules" .}}
{{template "packages" .}}
