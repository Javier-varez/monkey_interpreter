Object::makeFunction(
  Function{
    ConstexprLit<size_t, {{ len .Args }}>{},
    ConstexprLit<bool, {{ .VarArgs }}>{},
    [=]({{range $i, $el := .Args}} {{if $i}},{{end}} const Object {{$el}} {{end}}) noexcept -> Object {
      return ({ {{Transpile .Body}} });
    }
  }
)
