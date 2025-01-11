Mojo trait

# {{.GetName}}

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{template "fields" .}}
{{template "parentTraits" .}}
{{template "methods" .}}
