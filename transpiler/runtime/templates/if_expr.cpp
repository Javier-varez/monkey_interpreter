({
  Object _if_expr_result{};

  if (({{Transpile .Condition}}).getBool()) {
    _if_expr_result = ({
      {{Transpile .Consequence}}
    });
{{if .Alternative}}
  } else {
    _if_expr_result = ({
      {{Transpile .Alternative}}
    });
{{end}}
  }

  _if_expr_result;
})
