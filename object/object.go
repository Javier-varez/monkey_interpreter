package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/javier-varez/monkey_interpreter/ast"
	"github.com/javier-varez/monkey_interpreter/code"
	"github.com/javier-varez/monkey_interpreter/token"
)

type ObjectType string

const (
	INTEGER_OBJ           = "INTEGER"
	BOOLEAN_OBJ           = "BOOLEAN"
	STRING_OBJ            = "STRING"
	NULL_OBJ              = "NULL"
	RETURN_VALUE_OBJ      = "RETURN_VALUE"
	ERROR_VALUE_OBJ       = "ERROR_VALUE"
	FUNCTION_OBJ          = "FUNCTION"
	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION"
	CLOSURE_OBJ           = "CLOSURE"
	BUILTIN_OBJ           = "BUILTIN"
	ARRAY_OBJ             = "ARRAY"
	VAR_ARGS_OBJ          = "VAR_ARGS"
	MAP_OBJ               = "MAP"
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
	Args    []*ast.IdentifierExpr
	VarArgs bool
	Body    *ast.BlockStatement
	Env     *Environment
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
	Function BuiltinFunction
}

func (f *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (f *Builtin) Inspect() string {
	return "<Builtin>"
}

type CompiledFunction struct {
	Instructions code.Instructions
	NumLocals    int
	NumArgs      int
	VarArgs      bool
}

func (f *CompiledFunction) Type() ObjectType {
	return COMPILED_FUNCTION_OBJ
}

func (f *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", f)
}

type Closure struct {
	Fn          *CompiledFunction
	FreeObjects []Object
}

func (f *Closure) Type() ObjectType {
	return CLOSURE_OBJ
}

func (f *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", f)
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

type VarArgs struct {
	Elems []Object
}

func (a *VarArgs) Type() ObjectType {
	return VAR_ARGS_OBJ
}

func (a *VarArgs) Inspect() string {
	var buffer bytes.Buffer

	buffer.WriteString("VA[")
	for i, obj := range a.Elems {
		buffer.WriteString(obj.Inspect())
		if i != len(a.Elems)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("]")

	return buffer.String()
}

type HashKey struct {
	Type string
	Hash uint64
}

type Hashable interface {
	Object
	HashKey() HashKey
}

func (b *Boolean) HashKey() HashKey {
	hash := uint64(0)
	if b.Value {
		hash = 1
	}

	return HashKey{
		Type: BOOLEAN_OBJ,
		Hash: hash,
	}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type: INTEGER_OBJ,
		Hash: uint64(i.Value),
	}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64()
	h.Write([]byte(s.Value))
	hash := h.Sum64()

	return HashKey{
		Type: STRING_OBJ,
		Hash: hash,
	}
}

type HashEntry struct {
	Key   Object
	Value Object
	// TODO(ja): Deal with collisions. Implement chaining.
}

type HashMap struct {
	Elems map[HashKey]HashEntry
}

func (a *HashMap) Type() ObjectType {
	return MAP_OBJ
}

func (a *HashMap) Inspect() string {
	var buffer bytes.Buffer

	buffer.WriteString("{")
	for _, v := range a.Elems {
		buffer.WriteString(v.Key.Inspect())
		buffer.WriteString(":")
		buffer.WriteString(v.Value.Inspect())
		buffer.WriteString(",")
	}
	buffer.WriteString("}")

	return buffer.String()
}
