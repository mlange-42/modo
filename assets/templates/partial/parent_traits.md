{{define "parent_traits" -}}
{{if .ParentTraits}}## Implemented traits

{{range $i, $e := .ParentTraits -}}
{{ if $i }}, {{ end }}`{{ $e }}`{{end}}

{{end -}}
{{end}}
