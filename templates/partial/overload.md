{{define "overload" -}}
```mojo
{{.Signature}}
```

{{if .Summary}}{{.Summary}}

{{end -}}
{{if .Description}}{{.Description}}

{{end -}}
{{end}}