{{range .Global}}{{.}}
{{end}}

fn test_{{.Name}}() raises:
{{range .Code}}    {{.}}
{{end}}
