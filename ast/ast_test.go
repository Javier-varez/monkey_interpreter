package ast

import (
	"testing"

	"github.com/javier-varez/monkey_interpreter/token"
)

func TestString(t *testing.T) {
	program := Program{
		Statements: []Statment{
			&LetStatement{
				LetToken: token.Token{Type: token.LET, Literal: "let"},
				IdentExpr: &IdentifierExpr{
					IdentToken: token.Token{Type: token.IDENT, Literal: "myvar"},
				},
				AssignToken: token.Token{Type: token.ASSIGN, Literal: "="},
				Expr: &IdentifierExpr{
					IdentToken: token.Token{Type: token.IDENT, Literal: "anothervar"},
				},
				SemicolonToken: token.Token{Type: token.SEMICOLON, Literal: ";"},
			},
		},
	}

	if program.String() != "let myvar = anothervar;" {
		t.Fatalf("program String() is wrong: %q", program.String())
	}
}
