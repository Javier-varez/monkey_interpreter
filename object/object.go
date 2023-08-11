package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/javier-varez/monkey_interpreter/token"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_VALUE_OBJ  = "ERROR_VALUE"
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

type Null struct{}

func (i *Null) Type() ObjectType {
	return NULL_OBJ
}

func (i *Null) Inspect() string {
	return "null"
}

type Return struct {
	Value Object
}

func (i *Return) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (i *Return) Inspect() string {
	return i.Value.Inspect()
}

type Error struct {
	Message string
	Span    token.Span
}

func (i *Error) Type() ObjectType {
	return ERROR_VALUE_OBJ
}

func (i *Error) Inspect() string {
	return i.Message
}

const UNDERLINE = "\x1b[4m"
const UNDERLINE_RESET = "\x1b[24m"
const RED = "\x1b[31m"
const RESET_COLOR = "\x1b[0m"

func (e *Error) errLines(input string) []string {
	lines := strings.Split(input, "\n")
	return lines[e.Span.Start.Line : e.Span.End.Line+1]
}

func (e *Error) ContextualError(input string) string {
	var buffer bytes.Buffer

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