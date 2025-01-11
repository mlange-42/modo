{{define "overload" -}}
```mojo
{{.Signature}}
```

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{template "func_parameters" .}}
{{template "func_args" .}}
{{template "func_returns" .}}
{{template "func_raises" .}}
{{end}}
