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

func testErrorObject(expected *object.Error, actual object.Object) error {
	result, ok := actual.(*object.Error)
	if !ok {
		return fmt.Errorf("Object is not an *object.Error. got=%T (%+v)", actual, actual)
	}

	if result.Message != expected.Message {
		return fmt.Errorf("Object has wrong value. got=%q, want=%q", result.Message, expected.Message)
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
	case *object.Error:
		err := testErrorObject(expected, actual)
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

func TestFunctionCalls(t *testing.T) {
	tests := []vmTestCase{
		{`let a = fn() { 5 + 10 }; a()`, 15},
		{`let one = fn() { 1 }; let two = fn() { 2 }; one() + two()`, 3},
		{`let a = fn() { 1 }; let b = fn() { a() + 1 }; let c = fn() { b() + 1 }; c()`, 3},
		{`let a = fn() { return 10; 1; }; a()`, 10},
		{`let a = fn() { }; a()`, Null},
		{`let a = fn() { 1 }; let b = fn() { a }; b()()`, 1},
		{`let one = fn() { let one = 1; one }; one()`, 1},
		{
			`let firstFoobar = fn() { let foobar = 50; foobar; };
             let secondFoobar = fn() { let foobar = 100; foobar; };
             firstFoobar() + secondFoobar();`,
			150,
		},
		{
			`
			let globalSeed = 50;
            let minusOne = fn() {
                let num = 1;
                globalSeed - num;
            }
            let minusTwo = fn() {
                let num = 2;
                globalSeed - num;
            }
            minusOne() + minusTwo();
			`,
			97,
		},
		{
			`
			let a = 50;
            let clobberGlobal = fn() {
				let a = 10;
                a;
            }
            a;
			`,
			50,
		},
		{
			`
            let myFn = fn(a, b) {
				let c = 10;
                a + b + c;
            }
            myFn(2, 3);
			`,
			15,
		},
		{
			`
            let a = fn(a, b) {
				let c = 10;
                a + b + c;
            }
            a(2, 3) * a(5,7);
			`,
			15 * 22,
		},
		{`fn(a, ...) { let v = toArray(...); len(v) + a }(12, 1, 2, 3, 4)`, 16},
		{`fn(a, ...) { fn(a, b, c, d) { return a + b + c + d }(a, ...) }(1, 2, 3, 4)`, 10},
		{`fn(a, ...) { fn(a, b, c, d, ...) { return len(toArray(...)) }(a, ...) }(1, 2, 3, 4)`, 0},
		{`fn(a, ...) { fn(a, b, c, d, ...) { return len(toArray(...)) }(a, ...) }(1, 2, 3, 4, 5)`, 1},
		{`fn(a, ...) { fn(a, ...) { return len(toArray(...)) }(a, toArray(...)) }(1, 2, 3, 4, 5)`, 1},
		{`fn(a, ...) { fn(a, ...) { return len(toArray(...)[0]) }(a, toArray(...)) }(1, 2, 3, 4, 5)`, 4},
		{`fn(a, ...) { let v = toArray(...); last(v) + a }(12, 1, 2, 3, 4)`, 16},
	}

	runVmTests(t, tests)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `fn() { 1; }(1);`,
			expected: `wrong number of arguments: want=0, got=1`,
		},
		{
			input:    `fn(a) { a; }();`,
			expected: `wrong number of arguments: want=1, got=0`,
		},
		{
			input:    `fn(a, b) { a + b; }(1);`,
			expected: `wrong number of arguments: want=2, got=1`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program := parse(tt.input)
			comp := compiler.New()
			err := comp.Compile(program)
			if err != nil {
				t.Fatalf("compiler error: %s", err)
			}
			vm := New(comp.Bytecode())
			err = vm.Run()
			if err == nil {
				t.Fatalf("expected VM error but resulted in none.")
			}
			if err.Error() != tt.expected {
				t.Fatalf("wrong VM error: want=%q, got=%q", tt.expected, err)
			}
		})
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{
			`len(1)`,
			&object.Error{
				Message: "\"len\" builtin takes a single string or array argument",
			},
		},
		{`len("one", "two")`,
			&object.Error{
				Message: "\"len\" builtin takes a single string or array argument",
			},
		},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`puts("hello", "world!")`, Null},
		{`first([1, 2, 3])`, 1},
		{`first([])`,
			&object.Error{
				Message: "Array is empty",
			},
		},
		{`first(1)`,
			&object.Error{
				Message: "\"first\" builtin takes a single array argument",
			},
		},
		{`last([1, 2, 3])`, 3},
		{`last([])`,
			&object.Error{
				Message: "Array is empty",
			},
		},
		{`last(1)`,
			&object.Error{
				Message: "\"last\" builtin takes a single array argument",
			},
		},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`,
			&object.Error{
				Message: "Array is empty",
			},
		},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`,
			&object.Error{
				Message: "\"push\" builtin takes an array argument and a new object to push",
			},
		},
	}

	runVmTests(t, tests)

}

func TestClosures(t *testing.T) {
	tests := []vmTestCase{
		{`let a = fn(a) { let b = 10; fn(c) { 2 * b + 3 * a + c } }; a(40)(4)`, 144},
		{`let f = fn(...) { let b = 10; fn(c) { len(toArray(...)) + b + c } }; f(44,4,44,44,4)(4)`, 19},
		{`let newAdder = fn(a,b) { fn(c) {a+b+c} }; let adder = newAdder(1,2); adder(8)`, 11},
		{`let newAdder = fn(a,b) { let c = a + b; fn(d) {c+d} }; let adder = newAdder(1,2); adder(8)`, 11},
		{`let newAdderOuter = fn(a,b) { let c = a + b; fn(d) { let e = c+d; fn(f) {e+f} } }; let newAdderInner = newAdderOuter(1,2); let adder = newAdderInner(3); adder(8)`, 14},
		{`let newClosure = fn(a,b) { let one = fn() {a}; let two = fn() {b}; fn() { one() + two() }}; let closure = newClosure(9, 90); closure()`, 99},
	}

	runVmTests(t, tests)
}
