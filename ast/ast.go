package ast

import "github.com/javier-varez/monkey_interpreter/token"

type Node interface {
	Span() token.Span
}

type Expression interface {
	Node
	expressionNode()
}

type Statment interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statment
}

func (p *Program) Span() token.Span {
	// TODO(ja): Implement
	return token.Span{}
}

type LetStatement struct {
	LetToken       token.Token
	IdentExpr      *IdentifierExpr
	AssignToken    token.Token
	Expr           Expression
	SemicolonToken token.Token
}

func (stmt *LetStatement) statementNode() {}
func (stmt *LetStatement) Span() token.Span {
	// TODO(ja): Implement
	return token.Span{}
}

type ReturnStatement struct {
	ReturnToken    token.Token
	Expr           Expression
	SemicolonToken token.Token
}

func (stmt *ReturnStatement) statementNode() {}
func (stmt *ReturnStatement) Span() token.Span {
	// TODO(ja): Implement
	return token.Span{}
}

type IdentifierExpr struct {
	IdentToken token.Token
}

func (expr *IdentifierExpr) expressionNode() {}
func (expr *IdentifierExpr) Span() token.Span {
	// TODO(ja): Implement
	return token.Span{}
}
