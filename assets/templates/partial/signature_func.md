{{define "signature_func" -}}
```mojo
{{if and (.IsStatic) (ne .Name "__init__")}}@staticmethod
{{end -}}
{{if .IsDef}}def{{else}}fn{{end}} {{if .Signature}}{{.Signature}}{{else}}{{.Name}}{{end}}
```
{{- end}}