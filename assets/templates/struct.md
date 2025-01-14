Mojo struct

# `{{.Name}}`

```mojo
{{if .Convention}}@{{.Convention}}{{end}}
{{if .Signature}}{{.Signature}}{{else}}{{.Name}}{{end}}
```

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{template "parameters" .}}
{{template "aliases" .}}
{{template "fields" .}}
{{template "parentTraits" .}}
{{template "methods" .}}
