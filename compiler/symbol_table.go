package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "Global scope"
	LocalScope   SymbolScope = "Local scope"
	BuiltinScope SymbolScope = "Builtin scope"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Parent         *SymbolTable
	store          map[string]Symbol
	NumDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol)}
}

func NewEnclosedSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{Parent: parent, store: make(map[string]Symbol)}
}

func (st *SymbolTable) scope() SymbolScope {
	if st.Parent == nil {
		return GlobalScope
	} else {
		return LocalScope
	}
}

func (st *SymbolTable) Define(name string) Symbol {
	st.store[name] = Symbol{Name: name, Scope: st.scope(), Index: st.NumDefinitions}
	st.NumDefinitions++
	return st.store[name]
}

func (st *SymbolTable) DefineBuiltin(index int, name string) {
	st.store[name] = Symbol{Name: name, Scope: BuiltinScope, Index: index}
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	if sym, ok := st.store[name]; ok || st.Parent == nil {
		return sym, ok
	}

	return st.Parent.Resolve(name)
}
