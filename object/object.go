package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
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
