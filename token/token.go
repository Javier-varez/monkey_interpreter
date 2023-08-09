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

func (t *Token) IsLet() bool {
	return t.Literal == "let" && t.Type == LET
}

func (t *Token) IsIdent() bool {
	return t.Type == IDENT
}

func (t *Token) IsReturn() bool {
	return t.Type == RETURN && t.Literal == "return"
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
	EQ       = "=="
	NOT_EQ   = "!="

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdentifier(ident string) TokenType {
	if tt, ok := keywords[ident]; ok {
		return tt
	}
	return IDENT
}
