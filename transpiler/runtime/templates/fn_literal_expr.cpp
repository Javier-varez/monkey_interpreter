runtime::Object::makeFunction(
  runtime::Function{
    runtime::ConstexprLit<size_t, {{ len .Args }}>{},
    runtime::ConstexprLit<bool, {{ .VarArgs }}>{},
    [=]({{range $i, $el := .Args}} {{if $i}},{{end}} const runtime::Object {{$el}} {{end}} {{if .VarArgs}}{{if len .Args}},{{end}} const runtime::Object _varArgs{{end}}) noexcept -> runtime::Object {
      return ({ {{Transpile .Body}} });
    }
  }
)
