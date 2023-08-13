package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/javier-varez/monkey_interpreter/ast"
	"github.com/javier-varez/monkey_interpreter/token"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_VALUE_OBJ  = "ERROR_VALUE"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (i *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (i *Boolean) Inspect() string {
	return fmt.Sprintf("%t", i.Value)
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) Inspect() string {
	return fmt.Sprintf("%s", s.Value)
}

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (r *Return) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
	Span    token.Span
}

func (e *Error) Type() ObjectType {
	return ERROR_VALUE_OBJ
}

func (e *Error) Inspect() string {
	return e.Message
}

const UNDERLINE = "\x1b[4m"
const UNDERLINE_RESET = "\x1b[24m"
const RED = "\x1b[31m"
const RESET_COLOR = "\x1b[0m"

func (e *Error) errLines(input string) []string {
	lines := strings.Split(input, "\n")
	return lines[e.Span.Start.Line : e.Span.End.Line+1]
}

func (e *Error) ContextualError() string {
	var buffer bytes.Buffer

	input := *e.Span.Text
	startLine := e.Span.Start.Line
	endLine := e.Span.End.Line

	for lineIdx, line := range e.errLines(input) {
		if lineIdx > startLine && lineIdx < endLine {
			buffer.WriteString(line)
		} else {
			if lineIdx == startLine && lineIdx == endLine {
				firstPart := line[:e.Span.Start.Column]
				secondPart := line[e.Span.Start.Column:e.Span.End.Column]
				thirdPart := line[e.Span.End.Column:]
				buffer.WriteString(firstPart)
				buffer.WriteString(UNDERLINE)
				buffer.WriteString(secondPart)
				buffer.WriteString(UNDERLINE_RESET)
				buffer.WriteString(thirdPart)
			} else if lineIdx == startLine {
				firstPart := line[:e.Span.Start.Column]
				secondPart := line[e.Span.Start.Column:e.Span.End.Column]
				buffer.WriteString(firstPart)
				buffer.WriteString(UNDERLINE)
				buffer.WriteString(secondPart)
			} else if lineIdx == endLine {
				firstPart := line[e.Span.Start.Column:e.Span.End.Column]
				secondPart := line[e.Span.End.Column:]
				buffer.WriteString(firstPart)
				buffer.WriteString(UNDERLINE_RESET)
				buffer.WriteString(secondPart)
			}
		}
		buffer.WriteByte('\n')
	}

	buffer.WriteString(fmt.Sprintf("\t%s%s%s\n", RED, e.Message, RESET_COLOR))
	return buffer.String()
}

type Function struct {
	Args []*ast.IdentifierExpr
	Body *ast.BlockStatement
	Env  *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	out.WriteString("fn(")
	for i, arg := range f.Args {
		out.WriteString(arg.String())
		if i != len(f.Args)-1 {
			out.WriteString(",")
		}
	}
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}

type BuiltinFunction func(span token.Span, objects ...Object) Object

type Builtin struct {
	Function       BuiltinFunction
	NumArgs        int // -1 for any
	SupportedTypes []ObjectType
}

func (f *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (f *Builtin) Inspect() string {
	return "<Builtin>"
}

type Array struct {
	Elems []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (a *Array) Inspect() string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	for i, obj := range a.Elems {
		buffer.WriteString(obj.Inspect())
		if i != len(a.Elems)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("]")

	return buffer.String()
}
