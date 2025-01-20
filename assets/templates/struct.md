Mojo struct

# `{{.Name}}`

```mojo
{{if .Convention}}@{{.Convention}}{{end}}
{{if .Signature}}{{.Signature}}{{else}}{{.Name}}{{end}}
```

{{template "summary" . -}}
{{template "description" . -}}
{{template "aliases" . -}}
{{template "parameters" . -}}
{{template "fields" . -}}
{{template "parentTraits" . -}}
{{template "methods" . -}}