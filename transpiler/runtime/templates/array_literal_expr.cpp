runtime::Object::makeArray(runtime::Array{ {{range $i, $el := .Elems}}{{if $i}},{{end}}{{Transpile .}}{{end}} })
