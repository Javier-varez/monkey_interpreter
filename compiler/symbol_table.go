package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "Global scope"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol)}
}

func (st *SymbolTable) Define(name string) Symbol {
	st.store[name] = Symbol{Name: name, Scope: GlobalScope, Index: st.numDefinitions}
	st.numDefinitions++
	return st.store[name]
}

func (st *SymbolTable) Resolve(name string) (Symbol, bool) {
	sym, ok := st.store[name]
	return sym, ok
}
