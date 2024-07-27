package vm

import (
	"fmt"
	"testing"

	"github.com/javier-varez/monkey_interpreter/ast"
	"github.com/javier-varez/monkey_interpreter/compiler"
	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/object"
	"github.com/javier-varez/monkey_interpreter/parser"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("Object is not an integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("Object has wrong value. got=%d, want=%d", result.Value, expected)
	}
	return nil
}

func testBoolObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("Object is not a bool. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("Object has wrong value. got=%v, want=%v", result.Value, expected)
	}
	return nil
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v_%v", tt.input, tt.expected), func(t *testing.T) {
			program := parse(tt.input)

			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("compiler error: %s", err)
			}

			vm := New(comp.Bytecode())
			err = vm.Run()
			if err != nil {
				t.Fatalf("vm error: %s", err)
			}

			stackElem := vm.LastPoppedStackElem()
			testExpectedObject(t, tt.expected, stackElem)
		})
	}
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Fatalf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBoolObject(expected, actual)
		if err != nil {
			t.Fatalf("testBoolObject failed: %s", err)
		}
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"5 * 2", 10},
		{"5 / 2", 2},
	}

	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	runVmTests(t, tests)
}
