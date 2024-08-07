package compiler

import (
	"testing"
)

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
		"c": {Name: "c", Scope: LocalScope, Index: 0},
		"d": {Name: "d", Scope: LocalScope, Index: 1},
		"e": {Name: "e", Scope: LocalScope, Index: 0},
		"f": {Name: "f", Scope: LocalScope, Index: 1},
	}

	st := NewSymbolTable()

	a := st.Define("a")
	if a != expected["a"] {
		t.Errorf("Invalid symbol a. Expected=%q, Got=%q", expected["a"], a)
	}

	b := st.Define("b")
	if b != expected["b"] {
		t.Errorf("Invalid symbol b. Expected=%q, Got=%q", expected["b"], b)
	}

	local := NewEnclosedSymbolTable(st)
	c := local.Define("c")
	if c != expected["c"] {
		t.Errorf("Invalid symbol c. Expected=%q, Got=%q", expected["c"], c)
	}

	d := local.Define("d")
	if d != expected["d"] {
		t.Errorf("Invalid symbol d. Expected=%q, Got=%q", expected["d"], d)
	}

	nestedLocal := NewEnclosedSymbolTable(st)
	e := nestedLocal.Define("e")
	if e != expected["e"] {
		t.Errorf("Invalid symbol e. Expected=%q, Got=%q", expected["e"], e)
	}

	f := nestedLocal.Define("f")
	if f != expected["f"] {
		t.Errorf("Invalid symbol f. Expected=%q, Got=%q", expected["f"], f)
	}
}

func TestResolve(t *testing.T) {
	st := NewSymbolTable()
	st.Define("a")
	st.Define("b")

	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
	}

	for sym, exp_content := range expected {
		result, ok := st.Resolve(sym)
		if !ok {
			t.Errorf("Expected to find symbol %s", sym)
		} else {
			if result != exp_content {
				t.Errorf("Invalid symbol %s. Expected=%q, Got=%q", sym, exp_content, result)
			}
		}
	}
}

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
		"c": {Name: "c", Scope: LocalScope, Index: 0},
		"d": {Name: "d", Scope: LocalScope, Index: 1},
	}

	for sym, exp_content := range expected {
		result, ok := local.Resolve(sym)
		if !ok {
			t.Errorf("Expected to find symbol %s", sym)
		} else {
			if result != exp_content {
				t.Errorf("Invalid symbol %s. Expected=%q, Got=%q", sym, exp_content, result)
			}
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	local2 := NewEnclosedSymbolTable(local)
	local2.Define("e")
	local2.Define("f")

	tests := []struct {
		name     string
		table    *SymbolTable
		expected map[string]Symbol
	}{
		{
			name:  "local",
			table: local,
			expected: map[string]Symbol{
				"a": {Name: "a", Scope: GlobalScope, Index: 0},
				"b": {Name: "b", Scope: GlobalScope, Index: 1},
				"c": {Name: "c", Scope: LocalScope, Index: 0},
				"d": {Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			name:  "nestedLocal",
			table: local2,
			expected: map[string]Symbol{
				"a": {Name: "a", Scope: GlobalScope, Index: 0},
				"b": {Name: "b", Scope: GlobalScope, Index: 1},
				"c": {Name: "c", Scope: LocalScope, Index: 0},
				"d": {Name: "d", Scope: LocalScope, Index: 1},
				"e": {Name: "e", Scope: LocalScope, Index: 0},
				"f": {Name: "f", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for sym, exp_content := range tt.expected {
				result, ok := tt.table.Resolve(sym)
				if !ok {
					t.Errorf("Expected to find symbol %s", sym)
				} else {
					if result != exp_content {
						t.Errorf("Invalid symbol %s. Expected=%q, Got=%q", sym, exp_content, result)
					}
				}
			}
		})
	}
}
