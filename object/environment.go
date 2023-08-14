package object

type Environment struct {
	store      map[string]Object
	outer      *Environment
	hasVarArgs bool
	varArgs    []Object
}

func NewEnvironment() *Environment {
	return &Environment{store: map[string]Object{}}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	return &Environment{store: map[string]Object{}, outer: outer}
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) SetVarArgs(varArgs []Object) {
	e.hasVarArgs = true
	e.varArgs = varArgs
}

func (e *Environment) Get(name string) (Object, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		// Fall back to outer environment
		val, ok = e.outer.Get(name)
	}
	return val, ok
}

func (e *Environment) GetVarArgs() ([]Object, bool) {
	return e.varArgs, e.hasVarArgs
}

func (e *Environment) Copy() *Environment {
	newE := NewEnvironment()
	newE.outer = e.outer
	for k, v := range e.store {
		newE.Set(k, v)
	}
	return newE
}
