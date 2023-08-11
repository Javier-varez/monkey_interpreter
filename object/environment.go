package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{store: map[string]Object{}}
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Get(name string) (Object, bool) {
	val, ok := e.store[name]
	return val, ok
}
