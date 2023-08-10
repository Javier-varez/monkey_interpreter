package ast

import (
	"bytes"

	"github.com/javier-varez/monkey_interpreter/token"
)

type Node interface {
	Span() token.Span
	String() string
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

func (p *Program) String() string {
	var buf bytes.Buffer

	for _, s := range p.Statements {
		buf.WriteString(s.String())
	}

	return buf.String()
}

type LetStatement struct {
	LetToken       token.Token
	IdentExpr      Expression
	AssignToken    token.Token
	Expr           Expression
	SemicolonToken token.Token
}

func (stmt *LetStatement) statementNode() {}
func (stmt *LetStatement) Span() token.Span {
	// TODO(ja): Implement
	return token.Span{}
}
func (stmt *LetStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(stmt.LetToken.Literal + " ")
	if stmt.IdentExpr != nil {
		buf.WriteString(stmt.IdentExpr.String())
	}
	buf.WriteString(" " + stmt.AssignToken.Literal + " ")
	if stmt.Expr != nil {
		buf.WriteString(stmt.Expr.String())
	}
	buf.WriteString(stmt.SemicolonToken.Literal)

	return buf.String()
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

func (stmt *ReturnStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(stmt.ReturnToken.Literal + " ")
	if stmt.Expr != nil {
		buf.WriteString(stmt.Expr.String())
	}
	buf.WriteString(stmt.SemicolonToken.Literal)

	return buf.String()
}

type ExpressionStatement struct {
	Token          token.Token // First token in the expression
	Expr           Expression
	SemicolonToken *token.Token // Is optional
}

func (stmt *ExpressionStatement) statementNode() {}
func (stmt *ExpressionStatement) Span() token.Span {
	// TODO(ja): Implement
	return token.Span{}
}

func (stmt *ExpressionStatement) String() string {
	if stmt.Expr != nil {
		return stmt.Expr.String()
	}
	return ""
}

type IdentifierExpr struct {
	IdentToken token.Token
}

func (expr *IdentifierExpr) expressionNode() {}
func (expr *IdentifierExpr) Span() token.Span {
	// TODO(ja): Implement
	return token.Span{}
}

func (expr *IdentifierExpr) String() string {
	return expr.IdentToken.Literal
}

type IntegerLiteralExpr struct {
	IntToken token.Token
	Value    int64
}

func (expr *IntegerLiteralExpr) expressionNode() {}
func (expr *IntegerLiteralExpr) Span() token.Span {
	// TODO(ja): Implement
	return token.Span{}
}

func (expr *IntegerLiteralExpr) String() string {
	return expr.IntToken.Literal
}

type PrefixExpr struct {
	OperatorToken token.Token
	InnerExpr     Expression
}

func (expr *PrefixExpr) expressionNode() {}
func (expr *PrefixExpr) Span() token.Span {
	// TODO(ja): Implement
	return token.Span{}
}

func (expr *PrefixExpr) String() string {
	return "(" + expr.OperatorToken.Literal + expr.InnerExpr.String() + ")"
}

type InfixExpr struct {
	OperatorToken token.Token
	LeftExpr      Expression
	RightExpr     Expression
}

func (expr *InfixExpr) expressionNode() {}
func (expr *InfixExpr) Span() token.Span {
	// TODO(ja): Implement
	return token.Span{}
}

func (expr *InfixExpr) String() string {
	return "(" + expr.LeftExpr.String() + expr.OperatorToken.Literal + expr.RightExpr.String() + ")"
}
