package compiler

import (
	"testing"
)

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Scope: GlobalScope, Index: 1},
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
