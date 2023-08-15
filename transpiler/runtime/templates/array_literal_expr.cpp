Object::makeArray(Array{ {{range $i, $el := .Elems}}{{if $i}},{{end}}{{Transpile .}}{{end}} })
