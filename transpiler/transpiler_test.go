package transpiler

import (
	"testing"

	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/parser"
)

func testTranspile(input string) string {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	transpiled := Transpile(program)
	return Compile(transpiled)
}

func TestTranspileStringLiteral(t *testing.T) {
	test := []struct {
		input          string
		expectedOutput string
	}{
		{`puts("Hello, world!")`, "Hello, world!\n"},
		{`puts("Hello, world! ", "More strings")`, "Hello, world! More strings\n"},
	}
	for i, tt := range test {
		out := testTranspile(tt.input)
		if out != tt.expectedOutput {
			t.Errorf("[%d] Test failed. expected %q, got %q", i, tt.expectedOutput, out)
		}
	}
}

func TestTranspileIntLiteral(t *testing.T) {
	test := []struct {
		input          string
		expectedOutput string
	}{
		{`puts(123)`, "123\n"},
		{`puts(100)`, "100\n"},
	}
	for i, tt := range test {
		out := testTranspile(tt.input)
		if out != tt.expectedOutput {
			t.Errorf("[%d] Test failed. expected %q, got %q", i, tt.expectedOutput, out)
		}
	}
}

func TestTranspileBoolLiteral(t *testing.T) {
	test := []struct {
		input          string
		expectedOutput string
	}{
		{`puts(true)`, "true\n"},
		{`puts(false)`, "false\n"},
	}
	for i, tt := range test {
		out := testTranspile(tt.input)
		if out != tt.expectedOutput {
			t.Errorf("[%d] Test failed. expected %q, got %q", i, tt.expectedOutput, out)
		}
	}
}

func TestTranspilePrefixExpression(t *testing.T) {
	test := []struct {
		input          string
		expectedOutput string
	}{
		{`puts(-(10))`, "-10\n"},
		{`puts(-(-10))`, "10\n"},
		{`puts(!(false))`, "true\n"},
		{`puts(!(true))`, "false\n"},
	}
	for i, tt := range test {
		out := testTranspile(tt.input)
		if out != tt.expectedOutput {
			t.Errorf("[%d] Test failed. expected %q, got %q", i, tt.expectedOutput, out)
		}
	}
}

func TestTranspileInfixExpression(t *testing.T) {
	test := []struct {
		input          string
		expectedOutput string
	}{
		{`puts(10 + 12)`, "22\n"},
		{`puts(10 - 13)`, "-3\n"},
		{`puts(12 * 3)`, "36\n"},
		{`puts(12 / 3)`, "4\n"},
		{`puts(12 == 3)`, "false\n"},
		{`puts(12 == 12)`, "true\n"},
		{`puts(12 > 12)`, "false\n"},
		{`puts(12 < 12)`, "false\n"},
		{`puts(12 > 11)`, "true\n"},
		{`puts(12 < 13)`, "true\n"},
		{`puts(12 > 13)`, "false\n"},
		{`puts(12 < 11)`, "false\n"},
	}
	for i, tt := range test {
		out := testTranspile(tt.input)
		if out != tt.expectedOutput {
			t.Errorf("[%d] Test failed. expected %q, got %q", i, tt.expectedOutput, out)
		}
	}
}

func TestFunctionCall(t *testing.T) {
	test := []struct {
		input          string
		expectedOutput string
	}{
		{`let add = fn(x, y) { x + y }; puts(add(10, 12))`, "22\n"},
	}
	for i, tt := range test {
		out := testTranspile(tt.input)
		if out != tt.expectedOutput {
			t.Errorf("[%d] Test failed. expected %q, got %q", i, tt.expectedOutput, out)
		}
	}
}

func TestClosure(t *testing.T) {
	test := []struct {
		input          string
		expectedOutput string
	}{
		{"let makeAddN = fn(x) { fn(y) { x + y } }; let addTwo = makeAddN(2); puts(addTwo(123))", "125\n"},
		{"let makeAddN = fn(x) { fn(y) { x + y } }; let addTwo = makeAddN(2); let addThree = makeAddN(3); puts(addThree(123))", "126\n"},
	}
	for i, tt := range test {
		out := testTranspile(tt.input)
		if out != tt.expectedOutput {
			t.Errorf("[%d] Test failed. expected %q, got %q", i, tt.expectedOutput, out)
		}
	}
}

func TestCallSelf(t *testing.T) {
	test := []struct {
		input          string
		expectedOutput string
	}{
		{"let fib = fn(self, n) { if (n < 2) { n; } else { self(self, n-1) + self(self, n-2) } }; puts(fib(fib, 10))", "55\n"},
	}
	for i, tt := range test {
		out := testTranspile(tt.input)
		if out != tt.expectedOutput {
			t.Errorf("[%d] Test failed. expected %q, got %q", i, tt.expectedOutput, out)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	test := []struct {
		input          string
		expectedOutput string
	}{
		{"puts(fn() { return 10 }())", "10\n"},
		{"puts(fn() { return true }())", "true\n"},
		{"puts(fn() { return; }())", "nil\n"},
		{"puts(fn() { return 10; 2; }())", "10\n"},
		{"puts(fn() { return false; 2; }())", "false\n"},
		{"puts(fn() { return; 2; }())", "nil\n"},
		{"puts(fn() { if (100 < 200) { 2 * 2; return 33; 22; }; 2; }())", "33\n"},
		{"puts(fn() { if (200 < 200) { 2 * 2; return 33; 22; }; 2; }())", "2\n"},
		{"puts(fn() { if (100 < 200) { 2 * 2; if (1 != 2) { return 33; }; return 22; }; 2; }())", "33\n"},
	}
	for i, tt := range test {
		out := testTranspile(tt.input)
		if out != tt.expectedOutput {
			t.Errorf("[%d] Test failed. expected %q, got %q", i, tt.expectedOutput, out)
		}
	}
}
