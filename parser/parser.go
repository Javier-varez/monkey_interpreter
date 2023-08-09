package parser

import (
	"fmt"

	"github.com/javier-varez/monkey_interpreter/ast"
	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
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
	}

	return program, nil
}

func (p *Parser) parseIdentExpr() (*ast.IdentifierExpr, error) {
	expr := &ast.IdentifierExpr{}
	if !p.curToken.IsIdent() {
		return nil, fmt.Errorf("Not an identifier: %+v", p.curToken)
	}
	expr.IdentToken = p.curToken
	p.nextToken()
	return expr, nil
}

func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{}

	if !p.curToken.IsLet() {
		return nil, fmt.Errorf("Let expression does not start with `let`: %+v", p.curToken)
	}
	stmt.LetToken = p.curToken
	p.nextToken()

	identExpr, err := p.parseIdentExpr()
	if err != nil {
		return nil, err
	}
	stmt.IdentExpr = identExpr

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
	p.nextToken()

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
	p.nextToken()

	return stmt, nil
}

func (p *Parser) parseStatement() (ast.Statment, error) {
	if p.curToken.IsLet() {
		return p.parseLetStatement()
	} else if p.curToken.IsReturn() {
		return p.parseReturnStatement()
	}

	// Unknown statement
	return nil, fmt.Errorf("Unknown statement starting with token: %v\n", p.curToken)
}
