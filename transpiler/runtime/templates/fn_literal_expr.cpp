Object::makeFunction(
  Function{
    ConstexprLit<size_t, {{ len .Args }}>{},
    ConstexprLit<bool, {{ .VarArgs }}>{},
    [=]({{range $i, $el := .Args}} {{if $i}},{{end}} const Object {{$el}} {{end}} {{if .VarArgs}}{{if len .Args}},{{end}} const Object _varArgs{{end}}) noexcept -> Object {
      return ({ {{Transpile .Body}} });
    }
  }
)
