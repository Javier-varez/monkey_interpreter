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
		expected []Symbol
	}{
		{
			name:  "local",
			table: local,
			expected: []Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			name:  "nestedLocal",
			table: local2,
			expected: []Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: FreeScope, Index: 0},
				{Name: "d", Scope: FreeScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, exp_content := range tt.expected {
				result, ok := tt.table.Resolve(exp_content.Name)
				if !ok {
					t.Errorf("Expected to find symbol %s", exp_content.Name)
				} else {
					if result != exp_content {
						t.Errorf("Invalid symbol %s. Expected=%q, Got=%q", exp_content.Name, exp_content, result)
					}
				}
			}
		})
	}
}

func TestResolveFree(t *testing.T) {
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
		name         string
		table        *SymbolTable
		expected     []Symbol
		expectedFree []Symbol
	}{
		{
			name:  "local",
			table: local,
			expected: []Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
		{
			name:  "nestedLocal",
			table: local2,
			expected: []Symbol{
				{Name: "a", Scope: GlobalScope, Index: 0},
				{Name: "b", Scope: GlobalScope, Index: 1},
				{Name: "c", Scope: FreeScope, Index: 0},
				{Name: "d", Scope: FreeScope, Index: 1},
				{Name: "e", Scope: LocalScope, Index: 0},
				{Name: "f", Scope: LocalScope, Index: 1},
			},
			expectedFree: []Symbol{
				{Name: "c", Scope: LocalScope, Index: 0},
				{Name: "d", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, exp_content := range tt.expected {
				result, ok := tt.table.Resolve(exp_content.Name)
				if !ok {
					t.Errorf("Expected to find symbol %s", exp_content.Name)
				} else {
					if result != exp_content {
						t.Errorf("Invalid symbol %s. Expected=%q, Got=%q", exp_content.Name, exp_content, result)
					}
				}
			}

			if len(tt.expectedFree) != len(tt.table.FreeSymbols) {
				t.Errorf("wrong number of free symbols. got=%d, want=%d",
					len(tt.table.FreeSymbols), len(tt.expectedFree))
			}

			for i, expected := range tt.expectedFree {
				actual := tt.table.FreeSymbols[i]
				if actual != expected {
					t.Errorf("wrong free symbol. got=%v, want=%v", actual, expected)
				}
			}
		})
	}
}

func TestResolveUnresolvableFree(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	expected := []Symbol{
		{Name: "a", Scope: GlobalScope, Index: 0},
		{Name: "c", Scope: FreeScope, Index: 0},
		{Name: "e", Scope: LocalScope, Index: 0},
		{Name: "f", Scope: LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := secondLocal.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}
		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v",
				sym.Name, sym, result)
		}
	}

	expectedUnresolvable := []string{
		"b",
		"d",
	}

	for _, name := range expectedUnresolvable {
		_, ok := secondLocal.Resolve(name)
		if ok {
			t.Errorf("name %s resolved, but was expected not to", name)
		}
	}
}
