runtime::Object::makeHashMap(runtime::HashMap{
    {{range $k, $v := .Map}}
      runtime::HashMap::KvPair{ {{Transpile $k}}, {{Transpile $v}} },
    {{end}}
})
