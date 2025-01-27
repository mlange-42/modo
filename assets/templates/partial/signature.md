{{define "signature" -}}
```mojo
{{if .Signature}}{{.Signature}}{{else}}{{.Name}}{{end}}
```
{{- end}}