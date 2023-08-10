package parser

import (
	"fmt"
	"strconv"

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
}

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	l              *lexer.Lexer
	curToken       token.Token
	peekToken      token.Token
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.prefixParseFns[token.IDENT] = p.parseIdentExpr
	p.prefixParseFns[token.INT] = p.parseIntegerLiteralExpr
	p.prefixParseFns[token.BANG] = p.parsePrefixExpr
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpr
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.infixParseFns[token.PLUS] = p.parseInfixExpr
	p.infixParseFns[token.MINUS] = p.parseInfixExpr
	p.infixParseFns[token.ASTERISK] = p.parseInfixExpr
	p.infixParseFns[token.SLASH] = p.parseInfixExpr
	p.infixParseFns[token.EQ] = p.parseInfixExpr
	p.infixParseFns[token.NOT_EQ] = p.parseInfixExpr
	p.infixParseFns[token.GT] = p.parseInfixExpr
	p.infixParseFns[token.LT] = p.parseInfixExpr

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{
		Statements: []ast.Statment{},
	}

	for p.curToken.Type != token.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return program, err
		}
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}

	return program, nil
}

func (p *Parser) parseIdentExpr() ast.Expression {
	return &ast.IdentifierExpr{IdentToken: p.curToken}
}

func (p *Parser) parseIntegerLiteralExpr() ast.Expression {
	expr := &ast.IntegerLiteralExpr{IntToken: p.curToken}

	val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		fmt.Printf("could not parse %q as integer", p.curToken.Literal)
		return nil
	}
	expr.Value = val

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

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{}

	if !p.curToken.IsLet() {
		return nil, fmt.Errorf("Let expression does not start with `let`: %+v", p.curToken)
	}
	stmt.LetToken = p.curToken
	p.nextToken()

	if p.curToken.Type != token.IDENT {
		return nil, fmt.Errorf("Let expected an identifier: %+v", p.curToken)
	}
	stmt.IdentExpr = p.parseIdentExpr()
	p.nextToken()

	if p.curToken.Type != token.ASSIGN {
		return nil, fmt.Errorf("Expected `=` in assignment operator: %+v", p.curToken)
	}
	stmt.AssignToken = p.curToken
	p.nextToken()

	// TODO(ja): Parse actual expression instead of ignoring everything until semicolon
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}
	stmt.SemicolonToken = p.curToken

	return stmt, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{}

	if !p.curToken.IsReturn() {
		return nil, fmt.Errorf("Return expression does not start with `return`: %+v", p.curToken)
	}
	stmt.ReturnToken = p.curToken
	p.nextToken()

	// TODO(ja): Parse actual expression instead of ignoring everything until semicolon
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	stmt.SemicolonToken = p.curToken

	return stmt, nil
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix, ok := p.prefixParseFns[p.curToken.Type]
	if !ok || prefix == nil {
		// TODO(ja): log error
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

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expr = p.parseExpression(LOWEST)
	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}
	return stmt, nil
}

func (p *Parser) parseStatement() (ast.Statment, error) {
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
