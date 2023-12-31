package lexer

import (
	"github.com/javier-varez/monkey_interpreter/token"
)

type Lexer struct {
	input          string
	position       int  // current position in input (points to current char)
	ch             byte // current char under examination
	currentLine    int  // Keeps track of the current line
	lineByteOffset int  // Keeps track of the offset to the start of the line in input
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.ch = l.input[l.position]
	return l
}

func (l *Lexer) readChar() {
	next := l.position + 1
	if next >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[next]
	}
	l.position = next
}

func newToken(tokenType token.TokenType, literal byte, line, col int, text *string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(literal),
		Span: token.Span{
			Text:  text,
			Start: token.Location{Line: line, Column: col},
			End:   token.Location{Line: line, Column: col + 1},
		},
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token
	switch l.ch {
	case '=':
		if l.peekChar(1) == '=' {
			tok.Type = token.EQ
			ch := l.ch
			position := l.position
			l.readChar()
			tok.Literal = string(ch) + string(l.ch)
			tok.Span = token.Span{
				Text:  &l.input,
				Start: token.Location{Line: l.currentLine, Column: position - l.lineByteOffset},
				End:   token.Location{Line: l.currentLine, Column: position + 2 - l.lineByteOffset},
			}
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case ':':
		tok = newToken(token.COLON, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case '!':
		if l.peekChar(1) == '=' {
			tok.Type = token.NOT_EQ
			ch := l.ch
			position := l.position
			l.readChar()
			tok.Literal = string(ch) + string(l.ch)
			tok.Span = token.Span{
				Text:  &l.input,
				Start: token.Location{Line: l.currentLine, Column: position - l.lineByteOffset},
				End:   token.Location{Line: l.currentLine, Column: position + 2 - l.lineByteOffset},
			}
		} else {
			tok = newToken(token.BANG, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case '/':
		tok = newToken(token.SLASH, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case '>':
		tok = newToken(token.GT, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case '<':
		tok = newToken(token.LT, l.ch, l.currentLine, l.position-l.lineByteOffset, &l.input)
	case '.':
		return l.readDots()
	case '"':
		tok.Literal, tok.Span = l.readString()
		tok.Type = token.STRING
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Span = token.Span{
			Text:  &l.input,
			Start: token.Location{Line: l.currentLine, Column: l.position - l.lineByteOffset},
			End:   token.Location{Line: l.currentLine, Column: l.position - l.lineByteOffset},
		}
	default:
		if isLetter(l.ch) {
			tok.Literal, tok.Span = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal, tok.Span = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = l.illegalToken()
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) illegalToken() token.Token {
	return token.Token{
		Type:    token.ILLEGAL,
		Literal: string(l.ch),
		Span: token.Span{
			Text:  &l.input,
			Start: token.Location{Line: l.currentLine, Column: l.position - l.lineByteOffset},
			End:   token.Location{Line: l.currentLine, Column: l.position + 1 - l.lineByteOffset},
		},
	}
}

func (l *Lexer) readIdentifier() (string, token.Span) {
	line := l.currentLine
	startPos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position], token.Span{
		Text:  &l.input,
		Start: token.Location{Line: line, Column: startPos - l.lineByteOffset},
		End:   token.Location{Line: line, Column: l.position - l.lineByteOffset},
	}
}

func (l *Lexer) readNumber() (string, token.Span) {
	line := l.currentLine
	startPos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position], token.Span{
		Text:  &l.input,
		Start: token.Location{Line: line, Column: startPos - l.lineByteOffset},
		End:   token.Location{Line: line, Column: l.position - l.lineByteOffset},
	}
}

func (l *Lexer) readString() (string, token.Span) {
	line := l.currentLine
	startPos := l.position
	// Skip the initial "
	l.readChar()
	for l.ch != '"' {
		l.readChar()
	}

	// Skip the last "
	l.readChar()
	return l.input[startPos:l.position], token.Span{
		Text:  &l.input,
		Start: token.Location{Line: line, Column: startPos - l.lineByteOffset},
		End:   token.Location{Line: line, Column: l.position - l.lineByteOffset},
	}
}

func (l *Lexer) readDots() token.Token {
	firstDot := l.ch
	secondDot := l.peekChar(1)
	thirdDot := l.peekChar(2)

	if firstDot != '.' || secondDot != '.' {
		return l.illegalToken()
	}

	startPos := l.position
	literal := string(firstDot) + string(secondDot)
	tokType := token.TWO_DOTS

	l.readChar()
	l.readChar()

	if thirdDot == '.' {
		literal = literal + string(thirdDot)
		tokType = token.THREE_DOTS
		l.readChar()
	}

	return token.Token{
		Type:    token.TokenType(tokType),
		Literal: literal,
		Span: token.Span{
			Text:  &l.input,
			Start: token.Location{Line: l.currentLine, Column: startPos - l.lineByteOffset},
			End:   token.Location{Line: l.currentLine, Column: l.position - l.lineByteOffset},
		},
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == '\n' || l.ch == '\r' || l.ch == ' ' || l.ch == '\t' {
		if l.ch == '\n' {
			l.currentLine += 1
			l.lineByteOffset = l.position + 1
		}
		l.readChar()
	}
}

func (l *Lexer) peekChar(offset int) byte {
	next := l.position + offset
	if next >= len(l.input) {
		return 0
	}
	return l.input[next]
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') || (ch == '_')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
