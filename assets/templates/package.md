Mojo package

# `{{.Name}}`

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{template "modules" .}}
{{template "packages" .}}