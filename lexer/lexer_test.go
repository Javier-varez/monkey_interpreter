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

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;
"test string that also has 123 numbers and -;/\\ special chars +=-<>";
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
		{Type: token.IF, Literal: "if", Span: newSpan(12, 0, 2)},
		{Type: token.LPAREN, Literal: "(", Span: newSpan(12, 3, 1)},
		{Type: token.INT, Literal: "5", Span: newSpan(12, 4, 1)},
		{Type: token.LT, Literal: "<", Span: newSpan(12, 6, 1)},
		{Type: token.INT, Literal: "10", Span: newSpan(12, 8, 2)},
		{Type: token.RPAREN, Literal: ")", Span: newSpan(12, 10, 1)},
		{Type: token.LBRACE, Literal: "{", Span: newSpan(12, 12, 1)},
		{Type: token.RETURN, Literal: "return", Span: newSpan(13, 1, 6)},
		{Type: token.TRUE, Literal: "true", Span: newSpan(13, 8, 4)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(13, 12, 1)},
		{Type: token.RBRACE, Literal: "}", Span: newSpan(14, 0, 1)},
		{Type: token.ELSE, Literal: "else", Span: newSpan(14, 2, 4)},
		{Type: token.LBRACE, Literal: "{", Span: newSpan(14, 7, 1)},
		{Type: token.RETURN, Literal: "return", Span: newSpan(15, 1, 6)},
		{Type: token.FALSE, Literal: "false", Span: newSpan(15, 8, 5)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(15, 13, 1)},
		{Type: token.RBRACE, Literal: "}", Span: newSpan(16, 0, 1)},
		{Type: token.INT, Literal: "10", Span: newSpan(18, 0, 2)},
		{Type: token.EQ, Literal: "==", Span: newSpan(18, 3, 2)},
		{Type: token.INT, Literal: "10", Span: newSpan(18, 6, 2)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(18, 8, 1)},
		{Type: token.INT, Literal: "10", Span: newSpan(19, 0, 2)},
		{Type: token.NOT_EQ, Literal: "!=", Span: newSpan(19, 3, 2)},
		{Type: token.INT, Literal: "9", Span: newSpan(19, 6, 1)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(19, 7, 1)},
		{Type: token.STRING, Literal: `"test string that also has 123 numbers and -;/\\ special chars +=-<>"`, Span: newSpan(20, 0, 69)},
		{Type: token.SEMICOLON, Literal: ";", Span: newSpan(20, 69, 1)},
		{Type: token.EOF, Literal: ``, Span: newSpan(21, 0, 0)},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Literal != tt.Literal {
			t.Fatalf("tests[%d] - tokenliteral wrong. expected=%v, got=%v", i, tt, tok)
		}
		if tok.Type != tt.Type {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%v, got=%v", i, tt, tok)
		}
		if tok.Span.Start != tt.Span.Start {
			t.Fatalf("tests[%d] - tokenspan wrong. expected=%v, got=%v", i, tt, tok)
		}

		if tok.Span.End != tt.Span.End {
			t.Fatalf("tests[%d] - tokenspan wrong. expected=%v, got=%v", i, tt, tok)
		}
	}
}
