package parser

import (
	"testing"

	"github.com/javier-varez/monkey_interpreter/ast"
	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/token"
)

func checkDiagnostics(t *testing.T, program *ast.Program) {
	if len(program.Diagnostics) != 0 {
		t.Errorf("Diagnostics in program:")
		for _, err := range program.Diagnostics {
			t.Errorf("%s", err.ContextualError())
		}
		t.Fatalf("Unrecoverable program diagnostics")
	}
}

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = true;
let foobar = x;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d\n", len(program.Statements))
	}

	tests := []struct {
		ident         string
		expectedValue interface{}
	}{
		{"x", 5}, {"y", true}, {"foobar", "x"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.ident, tt.expectedValue) {
			return
		}
	}
}

func isIdentExpr(i ast.Expression, txt string) bool {
	identExpr, ok := i.(*ast.IdentifierExpr)
	if !ok {
		return false
	}
	return identExpr.IdentToken.Literal == txt && identExpr.IdentToken.Type == token.IDENT
}

func testLetStatement(t *testing.T, statement ast.Statment, expectedIdent string, expectedValue interface{}) bool {
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement is not a let statement: %+v", statement)
		return false
	}

	if !letStatement.LetToken.IsLet() {
		t.Errorf("Let statement should contain a let token: %+v", statement)
		return false
	}

	if !isIdentExpr(letStatement.IdentExpr, expectedIdent) {
		t.Errorf("Let statement should contain a let token: %+v", statement)
		return false
	}

	if !testLiteralExpression(t, letStatement.Expr, expectedValue) {
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return false;
return ident;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d\n", len(program.Statements))
	}

	expected := []interface{}{
		5, false, "ident",
	}

	for idx, tt := range program.Statements {
		stmt, ok := tt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Not a return statement: %+v", tt)
			continue
		}

		if !stmt.ReturnToken.IsReturn() {
			t.Errorf("Return statement does not start with a return token: %+v", tt)
		}

		testLiteralExpression(t, stmt.Expr, expected[idx])
	}
}

func TestIdentifierExpressions(t *testing.T) {
	input := `
foobar;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got %d\n", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Not an expression statement: %+v", program.Statements[0])
	}

	expr, ok := stmt.Expr.(*ast.IdentifierExpr)
	if !ok {
		t.Fatalf("Expression is not an IdentifierExpr: %+v", stmt.Expr)
	}

	if expr.IdentToken.Literal != "foobar" {
		t.Fatalf("Identifier is not foobar: %q", expr.IdentToken.Literal)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `
512;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got %d\n", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Not an expression statement: %+v", program.Statements[0])
	}

	expr, ok := stmt.Expr.(*ast.IntegerLiteralExpr)
	if !ok {
		t.Fatalf("Expression is not an IntegerLiteralExpr: %T", stmt.Expr)
	}

	if expr.IntToken.Literal != "512" {
		t.Fatalf("literal is not 512: %q", expr.IntToken.Literal)
	}

	if expr.Value != 512 {
		t.Fatalf("value is not 512: %q", expr.IntToken.Literal)
	}
}

func TestPrefixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		operator   string
		intLiteral int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Errorf("ParseProgram() returned a nil program")
			continue
		}

		checkDiagnostics(t, program)

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements does not contain 1 statement. got %d\n", len(program.Statements))
			continue
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Not an expression statement: %+v", program.Statements[0])
			continue
		}

		prefixExpr, ok := stmt.Expr.(*ast.PrefixExpr)
		if !ok {
			t.Errorf("Expression is not a PrefixExpr: %T", stmt.Expr)
			continue
		}

		if prefixExpr.OperatorToken.Literal != test.operator {
			t.Errorf("Unexpected operator. got %q, expected %q", prefixExpr.OperatorToken.Literal, test.operator)
			continue
		}

		intLiteralExpr, ok := prefixExpr.InnerExpr.(*ast.IntegerLiteralExpr)
		if !ok {
			t.Errorf("Expression is not an IntegerLiteralExpr: %T", prefixExpr.InnerExpr)
			continue
		}

		if intLiteralExpr.Value != test.intLiteral {
			t.Errorf("Uexpected literal. got %q, expected %q", intLiteralExpr.Value, test.intLiteral)
			continue
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		leftLiteral  int64
		rightLiteral int64
	}{
		{"5 + 5;", "+", 5, 5},
		{"5/5;", "/", 5, 5},
		{"5*5;", "*", 5, 5},
		{"5-5;", "-", 5, 5},
		{"5<5;", "<", 5, 5},
		{"5>5;", ">", 5, 5},
		{"5== 5;", "==", 5, 5},
		{"5 !=5 ;", "!=", 5, 5},
	}

	for pIdx, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Errorf("ParseProgram() returned a nil program")
			continue
		}

		checkDiagnostics(t, program)

		if len(program.Statements) != 1 {
			t.Errorf("program[%d].Statements does not contain 1 statement. got %d\n", pIdx, len(program.Statements))
			continue
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Not an expression statement: %+v", program.Statements[0])
			continue
		}

		infixExpr, ok := stmt.Expr.(*ast.InfixExpr)
		if !ok {
			t.Errorf("Expression is not an InfixExpr: %T", stmt.Expr)
			continue
		}

		if infixExpr.OperatorToken.Literal != test.operator {
			t.Errorf("Unexpected operator. got %q, expected %q", infixExpr.OperatorToken.Literal, test.operator)
			continue
		}

		leftLiteralExpr, ok := infixExpr.LeftExpr.(*ast.IntegerLiteralExpr)
		if !ok {
			t.Errorf("Expression is not an IntegerLiteralExpr: %T", infixExpr.LeftExpr)
			continue
		}

		if leftLiteralExpr.Value != test.leftLiteral {
			t.Errorf("Uexpected left literal. got %q, expected %q", leftLiteralExpr.Value, test.leftLiteral)
			continue
		}

		rightLiteralExpr, ok := infixExpr.RightExpr.(*ast.IntegerLiteralExpr)
		if !ok {
			t.Errorf("Expression is not an IntegerLiteralExpr: %T", infixExpr.RightExpr)
			continue
		}

		if rightLiteralExpr.Value != test.rightLiteral {
			t.Errorf("Uexpected right literal. got %q, expected %q", rightLiteralExpr.Value, test.rightLiteral)
			continue
		}
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"-a * b", "((-a)*b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a+b)+c)"},
		{"a + b - c", "((a+b)-c)"},
		{"a * b * c", "((a*b)*c)"},
		{"a * b / c", "((a*b)/c)"},
		{"a + b / c", "(a+(b/c))"},
		{"a + b * c + d / e - f", "(((a+(b*c))+(d/e))-f)"},
		{"3 + 4; -5 * 5", "(3+4);((-5)*5)"},
		{"5 > 4 == 3 < 4", "((5>4)==(3<4))"},
		{"5 < 4 != 3 > 4", "((5<4)!=(3>4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3+(4*5))==((3*1)+(4*5)))"},
		{"3 + 4 * 5 == true", "((3+(4*5))==true)"},
		{"false != true", "(false!=true)"},
		{"1 + (2 + 3) + 4", "((1+(2+3))+4)"},
		{"(5 + 5) * 2", "((5+5)*2)"},
		{"2 / (5 + 5)", "(2/(5+5))"},
		{"-(5 + 5)", "(-(5+5))"},
		{"!(true==true)", "(!(true==true))"},
		{"a + add(b * c) + d", "((a+add((b*c)))+d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a,b,1,(2*3),(4+5),add(6,(7*8)))"},
		{"let a = 300;", "let a = 300;"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Errorf("ParseProgram() returned a nil program")
			continue
		}

		checkDiagnostics(t, program)

		if program.String() != test.output {
			t.Errorf("Unexpected output for input %q. Expected %q, got %q", test.input, test.output, program.String())
			continue
		}
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, expected int64) bool {
	intLiteralExpr, ok := exp.(*ast.IntegerLiteralExpr)
	if !ok {
		t.Errorf("Expression is not an integer literal: %T", exp)
		return false
	}

	if intLiteralExpr.Value != expected {
		t.Errorf("Integer literal does not match. Expected %d, got %d", expected, intLiteralExpr.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, expected string) bool {
	identExpr, ok := exp.(*ast.IdentifierExpr)
	if !ok {
		t.Errorf("Not an IdentifierExpr: %T", exp)
		return false
	}

	if identExpr.IdentToken.Literal != expected {
		t.Errorf("IdentifierExpr literal does not match. Expected %q, got %q", expected, identExpr.IdentToken.Literal)
		return false
	}

	return true
}

func testBoolLiteral(t *testing.T, exp ast.Expression, expected bool) bool {
	boolLiteralExpr, ok := exp.(*ast.BoolLiteralExpr)
	if !ok {
		t.Errorf("Expression is not a bool literal: %T", exp)
		return false
	}

	if boolLiteralExpr.Value != expected {
		t.Errorf("Bool literal does not match. Expected %v, got %v", expected, boolLiteralExpr.Value)
		return false
	}

	return true
}

func testStringLiteralExpression(t *testing.T, exp ast.Expression, expected string) bool {
	strLiteralExpr, ok := exp.(*ast.StringLiteralExpr)
	if !ok {
		t.Errorf("Expression is not a string literal: %T", exp)
		return false
	}

	if strLiteralExpr.Value != expected {
		t.Errorf("String literal does not match. Expected %v, got %v", expected, strLiteralExpr.Value)
		return false
	}

	return true
}

func testArrayLiteralExpression(t *testing.T, expr ast.Expression, elems []interface{}) bool {
	arrayLiteralExpr, ok := expr.(*ast.ArrayLiteralExpr)
	if !ok {
		t.Errorf("Expression is not an array literal")
		return false
	}

	if len(arrayLiteralExpr.Elems) != len(elems) {
		t.Errorf("Number of elements in array does not match. Expected %d, got %d", len(elems), len(arrayLiteralExpr.Elems))
		return false
	}

	for i := range arrayLiteralExpr.Elems {
		if !testLiteralExpression(t, arrayLiteralExpr.Elems[i], elems[i]) {
			t.Errorf("Error found in element with index: %d", i)
			return false
		}
	}

	return true
}

func testRangeExpression(t *testing.T, expr ast.Expression, startExpr, endExpr int64) bool {
	rangeExpr, ok := expr.(*ast.RangeExpr)
	if !ok {
		t.Errorf("Expression is not an array literal")
		return false
	}

	if !testIntegerLiteral(t, rangeExpr.StartExpr, startExpr) {
		return false
	}

	if !testIntegerLiteral(t, rangeExpr.EndExpr, endExpr) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expr ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expr, int64(v))
	case int64:
		return testIntegerLiteral(t, expr, v)
	case bool:
		return testBoolLiteral(t, expr, v)
	case string:
		return testIdentifier(t, expr, v)
	case []interface{}:
		return testArrayLiteralExpression(t, expr, v)
	}
	t.Errorf("type of exp not handled. got=%T", expr)
	return false
}

func testInfixExpression(t *testing.T, expr ast.Expression, left interface{}, operator string, right interface{}) bool {
	infixExpr, ok := expr.(*ast.InfixExpr)
	if !ok {
		t.Errorf("Expression is not an InfixExpr: %T", expr)
		return false
	}

	if infixExpr.OperatorToken.Literal != operator {
		t.Errorf("Unexpected operator. got %q, expected %q", infixExpr.OperatorToken.Literal, operator)
		return false
	}

	if !testLiteralExpression(t, infixExpr.LeftExpr, left) {
		t.Errorf("Error validating left literal expression. got %v", infixExpr.RightExpr)
		return false
	}
	if !testLiteralExpression(t, infixExpr.RightExpr, right) {
		t.Errorf("Error validating right literal expression. got %v", infixExpr.RightExpr)
		return false
	}

	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression: %T", program.Statements[0])
	}

	ifExpr, ok := stmt.Expr.(*ast.IfExpr)
	if !ok {
		t.Fatalf("Not an if expression: %T", stmt.Expr)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		t.Fatalf("Error in condition of if expr")
	}

	if len(ifExpr.Consequence.Statements) != 1 {
		t.Fatalf("Unexpected length for ifExpr.Consequence.Statements: %d", len(ifExpr.Consequence.Statements))
	}

	consequenceExprStmt, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence does not contain an expression statement")
	}

	if !testIdentifier(t, consequenceExprStmt.Expr, "x") {
		t.Fatalf("Error in condition of if expr")
	}

	if ifExpr.Alternative != nil {
		t.Fatalf("Has an unexpected alternative")
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression: %T", program.Statements[0])
	}

	ifExpr, ok := stmt.Expr.(*ast.IfExpr)
	if !ok {
		t.Fatalf("Not an if expression: %T", stmt.Expr)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		t.Fatalf("Error in condition of if expr")
	}

	if len(ifExpr.Consequence.Statements) != 1 {
		t.Fatalf("Unexpected length for ifExpr.Consequence.Statements: %d", len(ifExpr.Consequence.Statements))
	}

	consequenceExprStmt, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Consequence does not contain an expression statement")
	}

	if !testIdentifier(t, consequenceExprStmt.Expr, "x") {
		t.Fatalf("Error in consequence of if expr")
	}

	if ifExpr.Alternative == nil {
		t.Fatalf("No alternative")
	}

	if len(ifExpr.Alternative.Statements) != 1 {
		t.Fatalf("Unexpected length for ifExpr.Alternative.Statements: %d", len(ifExpr.Alternative.Statements))
	}

	alternativeExprStmt, ok := ifExpr.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Alternative does not contain an expression statement")
	}

	if !testIdentifier(t, alternativeExprStmt.Expr, "y") {
		t.Fatalf("Error in alternative of if expr")
	}
}

func TestFnLiteralExpression(t *testing.T) {
	input := `
fn(x,y) {
	x + y;
}`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression: %T", program.Statements[0])
	}

	fnLitExpr, ok := stmt.Expr.(*ast.FnLiteralExpr)
	if !ok {
		t.Fatalf("Not an fn literal expression: %T", stmt.Expr)
	}

	if len(fnLitExpr.Args) != 2 {
		t.Fatalf("Unexpected number of function args: %d", len(fnLitExpr.Args))
	}

	if !testIdentifier(t, fnLitExpr.Args[0], "x") {
		return
	}

	if !testIdentifier(t, fnLitExpr.Args[1], "y") {
		return
	}

	if fnLitExpr.Body == nil {
		t.Fatalf("No body for function")
	}

	if len(fnLitExpr.Body.Statements) != 1 {
		t.Fatalf("Unexpected number of statements for fn body: %d", len(fnLitExpr.Body.Statements))
	}

	bodyStmt, ok := fnLitExpr.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expr statement")
	}

	testInfixExpression(t, bodyStmt.Expr, "x", "+", "y")
}

func TestFnLiteralExpressionParams(t *testing.T) {
	tests := []struct {
		input   string
		args    []string
		VarArgs bool
	}{
		{"fn() {}", []string{}, false},
		{"fn(...) {}", []string{}, true},
		{"fn(x) {}", []string{"x"}, false},
		{"fn(x,...) {}", []string{"x"}, true},
		{"fn(x,y,z) {}", []string{"x", "y", "z"}, false},
		{"fn(abc,def) {}", []string{"abc", "def"}, false},
		{"fn(abc,def,...) {}", []string{"abc", "def"}, true},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		if program == nil {
			t.Fatalf("ParseProgram() returned a nil program")
		}

		checkDiagnostics(t, program)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is not an expression: %T", program.Statements[0])
		}

		fnLitExpr, ok := stmt.Expr.(*ast.FnLiteralExpr)
		if !ok {
			t.Fatalf("Not an fn literal expression: %T", stmt.Expr)
		}

		if len(fnLitExpr.Args) != len(test.args) {
			t.Fatalf("Unexpected number of function args: %d", len(fnLitExpr.Args))
		}

		if fnLitExpr.VarArgs != test.VarArgs {
			if fnLitExpr.VarArgs {
				t.Fatalf("Function has variable arguments, but they were not expected")
			} else {
				t.Fatalf("Function does not have variable arguments, but they were expected")
			}
		}

		for idx, expectedArg := range test.args {
			if !testIdentifier(t, fnLitExpr.Args[idx], expectedArg) {
				return
			}
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression: %T", program.Statements[0])
	}

	callExpr, ok := stmt.Expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("Not a call expression: %T", stmt.Expr)
	}

	if !testIdentifier(t, callExpr.CallableExpr, "add") {
		return
	}

	if len(callExpr.Args) != 3 {
		t.Fatalf("Unexpected number of arguments in callExpr.Args: %d", len(callExpr.Args))
	}

	if !testLiteralExpression(t, callExpr.Args[0], 1) {
		return
	}

	if !testInfixExpression(t, callExpr.Args[1], 2, "*", 3) {
		return
	}

	if !testInfixExpression(t, callExpr.Args[2], 4, "+", 5) {
		return
	}
}

func TestEmptyCallExpression(t *testing.T) {
	input := `add();`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression: %T", program.Statements[0])
	}

	callExpr, ok := stmt.Expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("Not a call expression: %T", stmt.Expr)
	}

	if !testIdentifier(t, callExpr.CallableExpr, "add") {
		return
	}

	if len(callExpr.Args) != 0 {
		t.Fatalf("Unexpected number of arguments in callExpr.Args: %d", len(callExpr.Args))
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"Hello world!"`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression: %T", program.Statements[0])
	}

	if !testStringLiteralExpression(t, stmt.Expr, "Hello world!") {
		return
	}
}

func TestArrayLiteralExpression(t *testing.T) {
	input := `[123, test, true]`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression: %T", program.Statements[0])
	}

	if !testArrayLiteralExpression(t, stmt.Expr, []interface{}{123, "test", true}) {
		return
	}
}

func TestArrayIndexOperator(t *testing.T) {
	input := `a[123]`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression: %T", program.Statements[0])
	}

	arrayIndexOperatorExpr, ok := stmt.Expr.(*ast.ArrayIndexOperatorExpr)
	if !ok {
		t.Fatalf("Not an array index operator expression: %v", stmt.Expr)
	}

	if !testLiteralExpression(t, arrayIndexOperatorExpr.ArrayExpr, "a") {
		return
	}

	if !testLiteralExpression(t, arrayIndexOperatorExpr.IndexExpr, 123) {
		return
	}
}

func TestVarArgExpr(t *testing.T) {
	input := `...`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression: %T", program.Statements[0])
	}

	varArgsLiteralExpr, ok := stmt.Expr.(*ast.VarArgsLiteralExpr)
	if !ok {
		t.Fatalf("Not an array index operator expression: %v", stmt.Expr)
	}

	if varArgsLiteralExpr.Token.Type != token.THREE_DOTS {
		t.Fatalf("Unexpected token in var args literal")
	}
}

func TestRangeExpression(t *testing.T) {
	input := `0..3`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	checkDiagnostics(t, program)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is not an expression: %T", program.Statements[0])
	}

	if !testRangeExpression(t, stmt.Expr, 0, 3) {
		return
	}
}
