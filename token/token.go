package token

type TokenType string

type Location struct {
	Line   int
	Column int
}

type Span struct {
	Start, End Location
}

type Token struct {
	Type    TokenType
	Literal string
	Span    Span
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var keywords = map[string]TokenType{
	"let": LET,
	"fn":  FUNCTION,
}

func LookupIdentifier(ident string) TokenType {
	if tt, ok := keywords[ident]; ok {
		return tt
	}
	return IDENT
}
