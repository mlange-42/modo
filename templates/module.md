Mojo module

# {{.GetName}}

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{template "aliases" .}}
{{template "structs" .}}
{{template "traits" .}}
{{template "functions" .}}