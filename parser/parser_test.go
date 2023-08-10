package parser

import (
	"testing"

	"github.com/javier-varez/monkey_interpreter/ast"
	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/token"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program, error := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	if error != nil {
		t.Fatalf("ParseProgram() returned an error: %v", error.Error())
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d\n", len(program.Statements))
	}

	tests := []struct {
		ident string
	}{
		{"x"}, {"y"}, {"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.ident) {
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

func testLetStatement(t *testing.T, statement ast.Statment, expectedIdent string) bool {
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

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 123456;
`
	l := lexer.New(input)
	p := New(l)

	program, error := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	if error != nil {
		t.Fatalf("ParseProgram() returned an error: %v", error.Error())
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d\n", len(program.Statements))
	}

	for _, tt := range program.Statements {
		stmt, ok := tt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Not a return statement: %+v", tt)
			continue
		}

		if !stmt.ReturnToken.IsReturn() {
			t.Errorf("Return statement does not start with a return token: %+v", tt)
		}
	}
}

func TestIdentifierExpressions(t *testing.T) {
	input := `
foobar;
`
	l := lexer.New(input)
	p := New(l)

	program, error := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	if error != nil {
		t.Fatalf("ParseProgram() returned an error: %v", error.Error())
	}

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

	program, error := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned a nil program")
	}

	if error != nil {
		t.Fatalf("ParseProgram() returned an error: %v", error.Error())
	}

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

		program, error := p.ParseProgram()
		if program == nil {
			t.Errorf("ParseProgram() returned a nil program")
			continue
		}

		if error != nil {
			t.Errorf("ParseProgram() returned an error: %v", error.Error())
			continue
		}

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

		program, error := p.ParseProgram()
		if program == nil {
			t.Errorf("ParseProgram() returned a nil program")
			continue
		}

		if error != nil {
			t.Errorf("ParseProgram() returned an error: %v", error.Error())
			continue
		}

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
