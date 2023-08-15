package token

type TokenType string

type Location struct {
	Line   int
	Column int
}

type Span struct {
	Text       *string
	Start, End Location
}

func (first Span) Join(second Span) Span {
	if first.Text != second.Text {
		panic("Attempted to join spans from two different inputs")
	}
	return Span{
		Text:  first.Text,
		Start: first.Start,
		End:   second.End,
	}
}

const (
	CMP_LESS comparisonResult = iota
	CMP_EQ
	CMP_GREATER
)

type comparisonResult int

func (l Location) compare(other Location) comparisonResult {
	if l.Line < other.Line {
		return CMP_LESS
	} else if l.Line > other.Line {
		return CMP_GREATER
	}

	if l.Column < other.Column {
		return CMP_LESS
	} else if l.Column > other.Column {
		return CMP_GREATER
	} else {
		return CMP_EQ
	}
}

func (s Span) Contains(l Location) bool {
	startComparison := s.Start.compare(l)
	endComparison := s.End.compare(l)
	return (startComparison == CMP_LESS || startComparison == CMP_EQ) && endComparison == CMP_GREATER
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

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

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

	COMMA      = ","
	COLON      = ":"
	SEMICOLON  = ";"
	TWO_DOTS   = ".."
	THREE_DOTS = "..."

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

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
