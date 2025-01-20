{{range .Global}}{{.}}
{{end}}

fn test_{{.Name}}():
{{range .Code}}    {{.}}
{{end}}
