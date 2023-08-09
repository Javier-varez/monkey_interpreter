package lexer

import (
	"testing"

	"github.com/javier-varez/monkey_interpreter/token"
)

func newSpan(line, columnStart, numChars int) token.Span {
	return token.Span{
		Start: token.Location{
			Line:   line,
			Column: columnStart,
		},
		End: token.Location{
			Line:   line,
			Column: columnStart + numChars,
		},
	}
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
    x + y;
};

let result = add(five, ten);

!-/*5;
5 < 10 > 5;
`

	tests := []token.Token{
		{Type: token.LET, Literal: "let", Span: newSpan(0, 0, 3)},
		{Type: token.IDENT, Literal: "five", Span: newSpan(0, 4, 4)},
		{Type: token.ASSIGN, Literal: "=", Span: newSpan(0, 9, 1)},
		{Type: token.INT, Literal: "5", Span: newSpan(0, 11, 1)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(0, 12, 1)},
		{Type: token.LET, Literal: "let", Span: newSpan(1, 0, 3)},
		{Type: token.IDENT, Literal: "ten", Span: newSpan(1, 4, 3)},
		{Type: token.ASSIGN, Literal: "=", Span: newSpan(1, 8, 1)},
		{Type: token.INT, Literal: "10", Span: newSpan(1, 10, 2)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(1, 12, 1)},
		{Type: token.LET, Literal: "let", Span: newSpan(3, 0, 3)},
		{Type: token.IDENT, Literal: "add", Span: newSpan(3, 4, 3)},
		{Type: token.ASSIGN, Literal: "=", Span: newSpan(3, 8, 1)},
		{Type: token.FUNCTION, Literal: "fn", Span: newSpan(3, 10, 2)},
		{Type: token.LPAREN, Literal: "(", Span: newSpan(3, 12, 1)},
		{Type: token.IDENT, Literal: "x", Span: newSpan(3, 13, 1)},
		{Type: token.COMMA, Literal: ",", Span: newSpan(3, 14, 1)},
		{Type: token.IDENT, Literal: "y", Span: newSpan(3, 16, 1)},
		{Type: token.RPAREN, Literal: ")", Span: newSpan(3, 17, 1)},
		{Type: token.LBRACE, Literal: "{", Span: newSpan(3, 19, 1)},
		{Type: token.IDENT, Literal: "x", Span: newSpan(4, 4, 1)},
		{Type: token.PLUS, Literal: "+", Span: newSpan(4, 6, 1)},
		{Type: token.IDENT, Literal: "y", Span: newSpan(4, 8, 1)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(4, 9, 1)},
		{Type: token.RBRACE, Literal: "}", Span: newSpan(5, 0, 1)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(5, 1, 1)},
		{Type: token.LET, Literal: "let", Span: newSpan(7, 0, 3)},
		{Type: token.IDENT, Literal: "result", Span: newSpan(7, 4, 6)},
		{Type: token.ASSIGN, Literal: "=", Span: newSpan(7, 11, 1)},
		{Type: token.IDENT, Literal: "add", Span: newSpan(7, 13, 3)},
		{Type: token.LPAREN, Literal: "(", Span: newSpan(7, 16, 1)},
		{Type: token.IDENT, Literal: "five", Span: newSpan(7, 17, 4)},
		{Type: token.COMMA, Literal: ",", Span: newSpan(7, 21, 1)},
		{Type: token.IDENT, Literal: "ten", Span: newSpan(7, 23, 3)},
		{Type: token.RPAREN, Literal: ")", Span: newSpan(7, 26, 1)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(7, 27, 1)},
		{Type: token.BANG, Literal: "!", Span: newSpan(9, 0, 1)},
		{Type: token.MINUS, Literal: "-", Span: newSpan(9, 1, 1)},
		{Type: token.SLASH, Literal: "/", Span: newSpan(9, 2, 1)},
		{Type: token.ASTERISK, Literal: "*", Span: newSpan(9, 3, 1)},
		{Type: token.INT, Literal: "5", Span: newSpan(9, 4, 1)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(9, 5, 1)},
		{Type: token.INT, Literal: "5", Span: newSpan(10, 0, 1)},
		{Type: token.LT, Literal: "<", Span: newSpan(10, 2, 1)},
		{Type: token.INT, Literal: "10", Span: newSpan(10, 4, 2)},
		{Type: token.GT, Literal: ">", Span: newSpan(10, 7, 1)},
		{Type: token.INT, Literal: "5", Span: newSpan(10, 9, 1)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(10, 10, 1)},
		{Type: token.EOF, Literal: "", Span: newSpan(11, 0, 0)},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Literal != tt.Literal {
			t.Fatalf("tests[%d] - tokenliteral wrong. expected=%q, got=%q", i, tt.Literal, tok.Literal)
		}
		if tok.Type != tt.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.Type, tok.Type)
		}
		if tok.Span != tt.Span {
			t.Fatalf("tests[%d] - tokenspan wrong. expected=%v, got=%v", i, tt.Span, tok.Span)
		}
	}
}
