{{if .Expr}}
return {{Transpile .Expr}};
{{else}}
return runtime::Object{};
{{end}}
runtime::Object{};
