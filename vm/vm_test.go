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

func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("Object is not a string. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("Object has wrong value. got=%q, want=%q", result.Value, expected)
	}
	return nil
}

func testArrayObject(t *testing.T, expected []interface{}, actual object.Object) error {
	result, ok := actual.(*object.Array)
	if !ok {
		return fmt.Errorf("Object is not an array. got=%T (%+v)", actual, actual)
	}

	if len(expected) != len(result.Elems) {
		return fmt.Errorf("Unexpected number of elements: got=%+v, expected=%+v", result.Elems, expected)
	}

	for i, elem := range result.Elems {
		testExpectedObject(t, expected[i], elem)
	}
	return nil
}

func testMapObject(t *testing.T, expected map[interface{}]interface{}, actual object.Object) error {
	result, ok := actual.(*object.HashMap)
	if !ok {
		return fmt.Errorf("Object is not a map. got=%T (%+v)", actual, actual)
	}

	if len(expected) != len(result.Elems) {
		return fmt.Errorf("Unexpected number of elements: got=%+v, expected=%+v", result.Elems, expected)
	}

	for _, v := range result.Elems {
		var key interface{}
		switch k := v.Key.(type) {
		case *object.Boolean:
			key = k.Value
		case *object.Integer:
			key = int(k.Value)
		case *object.String:
			key = k.Value
		default:
			return fmt.Errorf("Invalid key type used in map")
		}

		if expectedVal, ok := expected[key]; ok {
			testExpectedObject(t, expectedVal, v.Value)
		} else {
			return fmt.Errorf("Key %+v (%T) not found", v.Key, v.Key)
		}
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
	case *object.Null:
		if actual != Null {
			t.Fatalf("object is not Null: %T (%+v)", actual, actual)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Fatalf("testStringObject failed: %s", err)
		}
	case []interface{}:
		err := testArrayObject(t, expected, actual)
		if err != nil {
			t.Fatalf("testArrayObject failed: %s", err)
		}
	case map[interface{}]interface{}:
		err := testMapObject(t, expected, actual)
		if err != nil {
			t.Fatalf("testMapObject failed: %s", err)
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
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 } ", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
		{"!(if (false) { 10 })", true},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}

	runVmTests(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let a = 10; a;", 10},
		{"let a = 10; let b = 24; a + b;", 34},
		{"let a = 10; let b = 24; let c = 12; a + b + c;", 46},
	}

	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runVmTests(t, tests)
}

func TestArrayExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`[]`, []interface{}{}},
		{`[1, 2, 3]`, []interface{}{1, 2, 3}},
		{`[1, "hi", 3, 4]`, []interface{}{1, "hi", 3, 4}},
		{`[1 + 2, 3 * 4, 5 + 6]`, []interface{}{3, 12, 11}},
		{`[1 + 2, 3 * 4, 5 + 6][0]`, 3},
		{`[1 + 2, 3 * 4, 5 + 6][0 + 1 + 32 * 0]`, 12},
		{`[1 + 2, 3 * 4, 5 + 6][0 + 1 + 1 + 32 * 0]`, 11},
		{`[1 + 2, 3 * 4, 5 + 6][0 + 1 + 1 + 32]`, Null},
	}

	runVmTests(t, tests)
}

func TestHashExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`{}`, map[interface{}]interface{}{}},
		{`{ 1: 3, "3": 4 }`, map[interface{}]interface{}{1: 3, "3": 4}},
		{`{ 1: 3, "34": 4 }[0 + 1 + 123 * 0]`, 3},
		{`{ 1: 3, "34": 4 }["3" + "4"]`, 4},
		{`{ 1: 3, "34": 4 }["4" + "4"]`, Null},
	}

	runVmTests(t, tests)
}
