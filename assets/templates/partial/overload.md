{{define "overload" -}}
{{template "signature_func" .}}

{{template "summary" . -}}
{{template "description" . -}}
{{template "func_parameters" . -}}
{{template "func_args" . -}}
{{if .Returns}}{{template "func_returns" .}}{{else}}{{template "func_returns_old" .}}{{end -}}
{{template "func_raises" . -}}
{{end}}