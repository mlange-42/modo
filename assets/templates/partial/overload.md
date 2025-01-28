{{define "overload" -}}
{{template "signature_func" .}}

{{template "summary" . -}}
{{template "description" . -}}
{{template "func_parameters" . -}}
{{template "func_args" . -}}
{{template "func_returns" . -}}
{{template "func_raises" . -}}
{{end}}