{{define "signature" -}}
```mojo
{{if .Convention}}@{{.Convention}}{{end}}
{{if .Signature}}{{.Signature}}{{else}}{{.Name}}{{end}}
```
{{- end}}