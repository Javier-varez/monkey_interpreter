package ast

import (
	"bytes"
	"fmt"

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

type Error interface {
	error
	ContextualError() string
	Span() token.Span
}

type Program struct {
	Statements  []Statment
	Diagnostics []Error
}

func (p *Program) Span() token.Span {
	if len(p.Statements) == 0 {
		return token.Span{}
	}

	return p.Statements[0].Span().Join(p.Statements[len(p.Statements)-1].Span())
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
	SemicolonToken *token.Token
}

func (stmt *LetStatement) statementNode() {}

func (stmt *LetStatement) Span() token.Span {
	var endSpan token.Span

	if stmt.SemicolonToken != nil {
		endSpan = stmt.SemicolonToken.Span
	} else {
		endSpan = stmt.Expr.Span()
	}

	return stmt.LetToken.Span.Join(endSpan)
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
	if stmt.SemicolonToken != nil {
		buf.WriteString(stmt.SemicolonToken.Literal)
	}

	return buf.String()
}

type ReturnStatement struct {
	ReturnToken    token.Token
	Expr           Expression
	SemicolonToken *token.Token
}

func (stmt *ReturnStatement) statementNode() {}

func (stmt *ReturnStatement) Span() token.Span {
	var endSpan token.Span

	if stmt.SemicolonToken != nil {
		endSpan = stmt.SemicolonToken.Span
	} else {
		endSpan = stmt.Expr.Span()
	}

	return stmt.ReturnToken.Span.Join(endSpan)
}

func (stmt *ReturnStatement) String() string {
	var buf bytes.Buffer

	buf.WriteString(stmt.ReturnToken.Literal + " ")
	if stmt.Expr != nil {
		buf.WriteString(stmt.Expr.String())
	}
	if stmt.SemicolonToken != nil {
		buf.WriteString(stmt.SemicolonToken.Literal)
	}

	return buf.String()
}

type ExpressionStatement struct {
	Expr           Expression
	SemicolonToken *token.Token
}

func (stmt *ExpressionStatement) statementNode() {}

func (stmt *ExpressionStatement) Span() token.Span {
	var endSpan token.Span

	if stmt.SemicolonToken != nil {
		endSpan = stmt.SemicolonToken.Span
	} else {
		endSpan = stmt.Expr.Span()
	}

	return stmt.Expr.Span().Join(endSpan)

}

func (stmt *ExpressionStatement) String() string {
	var buffer bytes.Buffer
	if stmt.Expr != nil {
		buffer.WriteString(stmt.Expr.String())
	}
	if stmt.SemicolonToken != nil {
		buffer.WriteString(stmt.SemicolonToken.Literal)
	}
	return buffer.String()
}

type BlockStatement struct {
	Lbrace, Rbrace token.Token
	Statements     []Statment
}

func (stmt *BlockStatement) statementNode() {}

func (stmt *BlockStatement) Span() token.Span {
	return stmt.Lbrace.Span.Join(stmt.Rbrace.Span)
}

func (stmt *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString(stmt.Lbrace.Literal)
	for _, stmt := range stmt.Statements {
		out.WriteString(stmt.String())
	}
	out.WriteString(stmt.Rbrace.Literal)

	return out.String()
}

type IdentifierExpr struct {
	IdentToken token.Token
}

func (expr *IdentifierExpr) expressionNode() {}

func (expr *IdentifierExpr) Span() token.Span {
	return expr.IdentToken.Span
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
	return expr.IntToken.Span
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
	return expr.OperatorToken.Span.Join(expr.InnerExpr.Span())
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
	return expr.LeftExpr.Span().Join(expr.RightExpr.Span())
}

func (expr *InfixExpr) String() string {
	return "(" + expr.LeftExpr.String() + expr.OperatorToken.Literal + expr.RightExpr.String() + ")"
}

type BoolLiteralExpr struct {
	Token token.Token
	Value bool
}

func (expr *BoolLiteralExpr) expressionNode() {}

func (expr *BoolLiteralExpr) Span() token.Span {
	return expr.Token.Span
}

func (expr *BoolLiteralExpr) String() string {
	return fmt.Sprint(expr.Value)
}

type IfExpr struct {
	IfToken     token.Token
	Condition   Expression
	Consequence *BlockStatement
	ElseToken   *token.Token
	Alternative *BlockStatement
}

func (expr *IfExpr) expressionNode() {}

func (expr *IfExpr) Span() token.Span {
	var endSpan token.Span
	if expr.Alternative != nil {
		endSpan = expr.Alternative.Span()
	} else {
		endSpan = expr.Condition.Span()
	}
	return expr.IfToken.Span.Join(endSpan)
}

func (expr *IfExpr) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(expr.Condition.String())
	out.WriteString(" ")
	out.WriteString(expr.Consequence.String())
	if expr.ElseToken != nil {
		out.WriteString(expr.ElseToken.Literal)
	}
	if expr.Alternative != nil {
		out.WriteString(expr.Alternative.String())
	}

	return out.String()
}

type FnLiteralExpr struct {
	FnToken token.Token
	Args    []*IdentifierExpr
	VarArgs bool
	Body    *BlockStatement
}

func (expr *FnLiteralExpr) expressionNode() {}

func (expr *FnLiteralExpr) Span() token.Span {
	return expr.FnToken.Span.Join(expr.Body.Span())
}

func (expr *FnLiteralExpr) String() string {
	var out bytes.Buffer

	out.WriteString(expr.FnToken.Literal)
	out.WriteString("(")
	for i, arg := range expr.Args {
		out.WriteString(arg.String())
		if i != len(expr.Args)-1 || expr.VarArgs {
			out.WriteString(",")
		}
	}
	if expr.VarArgs {
		out.WriteString("...")
	}
	out.WriteString(") ")
	out.WriteString(expr.Body.String())

	return out.String()
}

type CallExpr struct {
	CallableExpr Expression
	Lparen       token.Token
	Args         []Expression
	Rparen       token.Token
}

func (expr *CallExpr) expressionNode() {}

func (expr *CallExpr) Span() token.Span {
	return expr.CallableExpr.Span().Join(expr.Rparen.Span)
}

func (expr *CallExpr) String() string {
	var out bytes.Buffer

	out.WriteString(expr.CallableExpr.String())
	out.WriteString(expr.Lparen.Literal)
	for i, arg := range expr.Args {
		out.WriteString(arg.String())
		if i != len(expr.Args)-1 {
			out.WriteString(",")
		}
	}
	out.WriteString(expr.Rparen.Literal)

	return out.String()
}

type StringLiteralExpr struct {
	StringLitToken token.Token
	Value          string
}

func (expr *StringLiteralExpr) expressionNode() {}

func (expr *StringLiteralExpr) Span() token.Span {
	return expr.StringLitToken.Span
}

func (expr *StringLiteralExpr) String() string {
	return expr.StringLitToken.Literal
}

type ArrayLiteralExpr struct {
	Lbracket, Rbracket token.Token
	Elems              []Expression
}

func (expr *ArrayLiteralExpr) expressionNode() {}

func (expr *ArrayLiteralExpr) Span() token.Span {
	return expr.Lbracket.Span.Join(expr.Rbracket.Span)
}

func (expr *ArrayLiteralExpr) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(expr.Lbracket.Literal)
	for i, obj := range expr.Elems {
		buffer.WriteString(obj.String())
		if i != len(expr.Elems)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString(expr.Rbracket.Literal)

	return buffer.String()
}

type ArrayIndexOperatorExpr struct {
	ArrayExpr          Expression
	Lbracket, Rbracket token.Token
	IndexExpr          Expression
}

func (expr *ArrayIndexOperatorExpr) expressionNode() {}

func (expr *ArrayIndexOperatorExpr) Span() token.Span {
	return expr.ArrayExpr.Span().Join(expr.Rbracket.Span)
}

func (expr *ArrayIndexOperatorExpr) String() string {
	var out bytes.Buffer

	out.WriteString(expr.ArrayExpr.String())
	out.WriteString(expr.Lbracket.Literal)
	out.WriteString(expr.IndexExpr.String())
	out.WriteString(expr.Rbracket.Literal)

	return out.String()
}

type VarArgsLiteralExpr struct {
	Token token.Token
}

func (expr *VarArgsLiteralExpr) expressionNode() {}

func (expr *VarArgsLiteralExpr) Span() token.Span {
	return expr.Token.Span
}

func (expr *VarArgsLiteralExpr) String() string {
	return expr.Token.Literal
}

type RangeExpr struct {
	StartExpr, EndExpr Expression
	DotsToken          token.Token
	Start, End         int64
}

func (expr *RangeExpr) expressionNode() {}

func (expr *RangeExpr) Span() token.Span {
	return expr.StartExpr.Span().Join(expr.EndExpr.Span())
}

func (expr *RangeExpr) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("(")
	buffer.WriteString(expr.StartExpr.String())
	buffer.WriteString(expr.DotsToken.Literal)
	buffer.WriteString(expr.EndExpr.String())
	buffer.WriteString(")")
	return buffer.String()
}
