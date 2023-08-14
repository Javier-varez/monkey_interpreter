package parser

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/javier-varez/monkey_interpreter/ast"
	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/token"
)

const (
	// Operator precedence is defined by this enumeration
	_ int = iota
	LOWEST
	EQUALS      // == or !=
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // - or !
	CALL        // fn(x)
	RANGE       // 1..2
	ARRAY_IDX   // array[idx]
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: ARRAY_IDX,
	token.TWO_DOTS: RANGE,
}

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

type parseError struct {
	span     token.Span
	errorMsg string
}

func (p *parseError) Error() string {
	return p.errorMsg
}

func (p *parseError) errLines() []string {
	lines := strings.Split(*p.span.Text, "\n")
	return lines[p.span.Start.Line : p.span.End.Line+1]
}

const UNDERLINE = "\x1b[4m"
const UNDERLINE_RESET = "\x1b[24m"
const RED = "\x1b[31m"
const RESET_COLOR = "\x1b[0m"

func (p *parseError) ContextualError() string {
	var buffer bytes.Buffer

	startLine := p.span.Start.Line
	endLine := p.span.End.Line

	for lineIdx, line := range p.errLines() {
		if lineIdx > startLine && lineIdx < endLine {
			buffer.WriteString(line)
		} else {
			if lineIdx == startLine && lineIdx == endLine {
				firstPart := line[:p.span.Start.Column]
				secondPart := line[p.span.Start.Column:p.span.End.Column]
				thirdPart := line[p.span.End.Column:]
				buffer.WriteString(firstPart)
				buffer.WriteString(UNDERLINE)
				buffer.WriteString(secondPart)
				buffer.WriteString(UNDERLINE_RESET)
				buffer.WriteString(thirdPart)
			} else if lineIdx == startLine {
				firstPart := line[:p.span.Start.Column]
				secondPart := line[p.span.Start.Column:p.span.End.Column]
				buffer.WriteString(firstPart)
				buffer.WriteString(UNDERLINE)
				buffer.WriteString(secondPart)
			} else if lineIdx == endLine {
				firstPart := line[p.span.Start.Column:p.span.End.Column]
				secondPart := line[p.span.End.Column:]
				buffer.WriteString(firstPart)
				buffer.WriteString(UNDERLINE_RESET)
				buffer.WriteString(secondPart)
			}
		}
		buffer.WriteByte('\n')
	}

	buffer.WriteString(fmt.Sprintf("\t%s%s%s\n", RED, p.errorMsg, RESET_COLOR))
	return buffer.String()
}

func (p *parseError) Span() token.Span {
	return p.span
}

func (p *Parser) mkError(s token.Span, msg string) {
	p.errors = append(p.errors, &parseError{
		span:     s,
		errorMsg: msg,
	})
}

type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
	errors         []ast.Error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.prefixParseFns[token.IDENT] = p.parseIdentExpr
	p.prefixParseFns[token.INT] = p.parseIntegerLiteralExpr
	p.prefixParseFns[token.BANG] = p.parsePrefixExpr
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpr
	p.prefixParseFns[token.TRUE] = p.parseBooleanLiteralExpr
	p.prefixParseFns[token.FALSE] = p.parseBooleanLiteralExpr
	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpr
	p.prefixParseFns[token.IF] = p.parseIfExpr
	p.prefixParseFns[token.FUNCTION] = p.parseFnLiteralExpr
	p.prefixParseFns[token.STRING] = p.parseStringLiteralExpr
	p.prefixParseFns[token.LBRACKET] = p.parseArrayLiteralExpr
	p.prefixParseFns[token.THREE_DOTS] = p.parseVarArgsLiteralExpr
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.infixParseFns[token.PLUS] = p.parseInfixExpr
	p.infixParseFns[token.MINUS] = p.parseInfixExpr
	p.infixParseFns[token.ASTERISK] = p.parseInfixExpr
	p.infixParseFns[token.SLASH] = p.parseInfixExpr
	p.infixParseFns[token.EQ] = p.parseInfixExpr
	p.infixParseFns[token.NOT_EQ] = p.parseInfixExpr
	p.infixParseFns[token.GT] = p.parseInfixExpr
	p.infixParseFns[token.LT] = p.parseInfixExpr
	p.infixParseFns[token.LPAREN] = p.parseCallExpr
	p.infixParseFns[token.LBRACKET] = p.parseArrayIndexExpr
	p.infixParseFns[token.TWO_DOTS] = p.parseRangeExpr

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statment{},
	}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}

	program.Diagnostics = p.errors
	return program
}

func (p *Parser) parseIdentExpr() ast.Expression {
	return &ast.IdentifierExpr{IdentToken: p.curToken}
}

func (p *Parser) parseIntegerLiteralExpr() ast.Expression {
	expr := &ast.IntegerLiteralExpr{IntToken: p.curToken}

	val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.mkError(p.curToken.Span, "Invalid integer literal. Could not be converted to a 64-bit 10-base integer.")
		return nil
	}
	expr.Value = val

	return expr
}

func (p *Parser) parseBooleanLiteralExpr() ast.Expression {
	if p.curToken.Type == token.FALSE {
		return &ast.BoolLiteralExpr{
			Token: p.curToken,
			Value: false,
		}
	} else if p.curToken.Type == token.TRUE {
		return &ast.BoolLiteralExpr{
			Token: p.curToken,
			Value: true,
		}
	}
	log.Fatalf("Unparsable boolean literal expression: %v\n", p.curToken)
	return nil
}

func (p *Parser) parseGroupedExpr() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if p.peekToken.Type != token.RPAREN {
		return nil
	}

	p.nextToken()
	return exp
}

func (p *Parser) parseIfExpr() ast.Expression {
	expr := ast.IfExpr{
		IfToken: p.curToken,
	}

	if p.peekToken.Type != token.LPAREN {
		return nil
	}

	p.nextToken()
	p.nextToken()

	expr.Condition = p.parseExpression(LOWEST)

	if p.peekToken.Type != token.RPAREN {
		return nil
	}
	p.nextToken()

	if p.peekToken.Type != token.LBRACE {
		return nil
	}
	p.nextToken()

	expr.Consequence = p.parseBlockStatement()

	if p.peekToken.Type != token.ELSE {
		// No else condition, return now
		return &expr
	}
	p.nextToken()

	elseToken := p.curToken
	expr.ElseToken = &elseToken

	if p.peekToken.Type != token.LBRACE {
		return nil
	}
	p.nextToken()

	expr.Alternative = p.parseBlockStatement()

	return &expr
}

func (p *Parser) parseFnLiteralExpr() ast.Expression {
	expr := &ast.FnLiteralExpr{
		FnToken: p.curToken,
	}

	if p.peekToken.Type != token.LPAREN {
		p.mkError(p.peekToken.Span, "fn literal must be followed by argument list")
		return nil
	}
	p.nextToken()
	p.nextToken()

	for p.curToken.Type != token.RPAREN {
		if (p.curToken.Type != token.IDENT) && (p.curToken.Type != token.THREE_DOTS) {
			p.mkError(p.curToken.Span, "Parameters to an fn literal must be identifier expressions or \"...\"")
			return nil
		}

		if p.peekToken.Type != token.COMMA && p.peekToken.Type != token.RPAREN {
			p.mkError(p.peekToken.Span, "Invalid token found in argument list of fn literal expression")
			return nil
		}

		if p.curToken.Type == token.THREE_DOTS {
			expr.VarArgs = true
			p.nextToken()
			if p.curToken.Type != token.RPAREN {
				p.mkError(p.curToken.Span, "`...` var args must be the last argument to a function")
				return nil
			}
		} else {
			identExpr := p.parseIdentExpr().(*ast.IdentifierExpr)
			expr.Args = append(expr.Args, identExpr)
			p.nextToken()
		}

		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	if p.peekToken.Type != token.LBRACE {
		p.mkError(p.peekToken.Span, "Expected body of fn literal")
		return nil
	}
	p.nextToken()

	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseStringLiteralExpr() ast.Expression {
	return &ast.StringLiteralExpr{
		StringLitToken: p.curToken,
		Value:          strings.Trim(p.curToken.Literal, "\""),
	}
}

func (p *Parser) parseVarArgsLiteralExpr() ast.Expression {
	return &ast.VarArgsLiteralExpr{Token: p.curToken}
}

func (p *Parser) parseArrayLiteralExpr() ast.Expression {
	expr := &ast.ArrayLiteralExpr{
		Lbracket: p.curToken,
		Elems:    []ast.Expression{},
	}

	p.nextToken()

	for p.curToken.Type != token.RBRACKET {
		inner := p.parseExpression(LOWEST)
		expr.Elems = append(expr.Elems, inner)

		p.nextToken()
		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	expr.Rbracket = p.curToken

	return expr
}

func (p *Parser) parsePrefixExpr() ast.Expression {
	expr := &ast.PrefixExpr{
		OperatorToken: p.curToken,
	}
	p.nextToken()

	expr.InnerExpr = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpr(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpr{
		LeftExpr:      left,
		OperatorToken: p.curToken,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	expr.RightExpr = p.parseExpression(precedence)
	return expr
}

func (p *Parser) parseCallExpr(left ast.Expression) ast.Expression {
	expr := &ast.CallExpr{
		CallableExpr: left,
		Lparen:       p.curToken,
		Args:         []ast.Expression{},
	}

	p.nextToken()

	for p.curToken.Type != token.RPAREN {
		argExpr := p.parseExpression(LOWEST)

		if p.peekToken.Type != token.COMMA && p.peekToken.Type != token.RPAREN {
			p.mkError(p.peekToken.Span, "Invalid delimiter token found in call expression argument list")
			return nil
		}

		expr.Args = append(expr.Args, argExpr)
		p.nextToken()
		if p.curToken.Type == token.COMMA {
			p.nextToken()
		}
	}

	expr.Rparen = p.curToken

	return expr
}

func (p *Parser) parseArrayIndexExpr(left ast.Expression) ast.Expression {
	expr := &ast.ArrayIndexOperatorExpr{
		ArrayExpr: left,
		Lbracket:  p.curToken,
	}

	p.nextToken()

	expr.IndexExpr = p.parseExpression(LOWEST)
	if p.peekToken.Type != token.RBRACKET {
		p.mkError(p.peekToken.Span, "Expected ] delimiter to close array index expression")
		return nil
	}
	p.nextToken()

	expr.Rbracket = p.curToken

	return expr
}

func (p *Parser) parseRangeExpr(left ast.Expression) ast.Expression {
	expr := &ast.RangeExpr{
		StartExpr: left,
		DotsToken: p.curToken,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.EndExpr = p.parseExpression(precedence)

	return expr
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{}

	if !p.curToken.IsLet() {
		p.mkError(p.curToken.Span, "Let statement does not start with \"let\"")
		return nil
	}
	stmt.LetToken = p.curToken
	p.nextToken()

	if p.curToken.Type != token.IDENT {
		p.mkError(p.curToken.Span, "Let statement expected an identifier")
		return nil
	}
	stmt.IdentExpr = p.parseIdentExpr()
	p.nextToken()

	if p.curToken.Type != token.ASSIGN {
		p.mkError(p.curToken.Span, "Expected \"=\" in let statement")
		return nil
	}
	stmt.AssignToken = p.curToken
	p.nextToken()

	stmt.Expr = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
		token := p.curToken
		stmt.SemicolonToken = &token
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{}

	if !p.curToken.IsReturn() {
		p.mkError(p.curToken.Span, "Return statement does not begin with \"return\"")
		return nil
	}
	stmt.ReturnToken = p.curToken
	p.nextToken()

	stmt.Expr = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
		token := p.curToken
		stmt.SemicolonToken = &token
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	stmt := &ast.BlockStatement{
		Lbrace: p.curToken,
	}

	for p.peekToken.Type != token.RBRACE {
		p.nextToken()
		s := p.parseStatement()
		stmt.Statements = append(stmt.Statements, s)
	}

	p.nextToken()
	stmt.Rbrace = p.curToken

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok || prefix == nil {
		p.mkError(p.curToken.Span, "Invalid token")
		return nil
	}

	leftExp := prefix()

	for (p.peekToken.Type != token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok || infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{}

	stmt.Expr = p.parseExpression(LOWEST)
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
		token := p.curToken
		stmt.SemicolonToken = &token
	}
	return stmt
}

func (p *Parser) parseStatement() ast.Statment {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.curToken.Type]; ok {
		return precedence
	}
	return LOWEST
}
