{{define "signature_struct" -}}
```mojo
{{if .Convention}}@{{.Convention}}
{{end -}}
{{if .Signature}}{{.Signature}}{{else}}{{.Name}}{{end}}
```
{{- end}}