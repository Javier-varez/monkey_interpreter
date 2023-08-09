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

func isIdentExpr(i ast.IdentifierExpr, txt string) bool {
	return i.IdentToken.Literal == txt && i.IdentToken.Type == token.IDENT
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

	if !isIdentExpr(*letStatement.IdentExpr, expectedIdent) {
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
