{{if .Expr}}
return {{Transpile .Expr}};
{{else}}
return Object{};
{{end}}
Object{};
