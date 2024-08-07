package object

import (
	"fmt"

	"github.com/javier-varez/monkey_interpreter/token"
)

func mkError(s token.Span, msg string) *Error {
	return &Error{Span: s, Message: msg}
}

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		Name: "len",
		Builtin: &Builtin{
			Function: func(span token.Span, objects ...Object) Object {
				if len(objects) != 1 {
					return mkError(span, "\"len\" builtin takes a single string or array argument")
				}

				strObj, ok := objects[0].(*String)
				if ok {
					return &Integer{Value: int64(len(strObj.Value))}
				}

				arrObj, ok := objects[0].(*Array)
				if ok {
					return &Integer{Value: int64(len(arrObj.Elems))}
				}

				return mkError(span, "\"len\" builtin takes a single string or array argument")
			},
		},
	},
	{
		Name: "first",
		Builtin: &Builtin{
			Function: func(span token.Span, objects ...Object) Object {
				if len(objects) != 1 {
					return mkError(span, "\"first\" builtin takes a single array argument")
				}

				arrObj, ok := objects[0].(*Array)
				if ok {
					if len(arrObj.Elems) == 0 {
						return mkError(span, "Array is empty")
					}
					return arrObj.Elems[0]
				}

				return mkError(span, "\"first\" builtin takes a single array argument")
			},
		},
	},
	{
		Name: "last",
		Builtin: &Builtin{
			Function: func(span token.Span, objects ...Object) Object {
				if len(objects) != 1 {
					return mkError(span, "\"last\" builtin takes a single array argument")
				}

				arrObj, ok := objects[0].(*Array)
				if ok {
					if len(arrObj.Elems) == 0 {
						return mkError(span, "Array is empty")
					}
					return arrObj.Elems[len(arrObj.Elems)-1]
				}

				return mkError(span, "\"last\" builtin takes a single array argument")
			},
		},
	},
	{
		Name: "rest",
		Builtin: &Builtin{
			Function: func(span token.Span, objects ...Object) Object {
				if len(objects) != 1 {
					return mkError(span, "\"rest\" builtin takes a single array argument")
				}

				arrObj, ok := objects[0].(*Array)
				if ok {
					if len(arrObj.Elems) == 0 {
						return mkError(span, "Array is empty")
					}
					newArr := &Array{Elems: make([]Object, len(arrObj.Elems)-1)}
					copy(newArr.Elems[:], arrObj.Elems[1:])
					return newArr
				}

				return mkError(span, "\"rest\" builtin takes a single array argument")
			},
		},
	},
	{
		Name: "push",
		Builtin: &Builtin{
			Function: func(span token.Span, objects ...Object) Object {
				if len(objects) != 2 {
					return mkError(span, "\"push\" builtin takes an array argument and a new object to push")
				}

				arrObj, ok := objects[0].(*Array)
				if !ok {
					return mkError(span, "\"push\" builtin takes an array argument and a new object to push")
				}

				oldLen := len(arrObj.Elems)
				newArr := &Array{Elems: make([]Object, oldLen+1)}
				copy(newArr.Elems[:oldLen], arrObj.Elems[:])
				newArr.Elems[oldLen] = objects[1]

				return newArr
			},
		},
	},
	{
		Name: "puts",
		Builtin: &Builtin{
			Function: func(span token.Span, objects ...Object) Object {
				for _, object := range objects {
					fmt.Print(object.Inspect())
				}
				fmt.Println()
				return &Null{}
			},
		},
	},
	{
		Name: "toArray",
		Builtin: &Builtin{
			Function: func(span token.Span, objects ...Object) Object {
				if len(objects) != 1 {
					return mkError(span, "\"toArray\" builtin takes a VarArg argument")
				}

				varArgObj, ok := objects[0].(*VarArgs)
				if !ok {
					return mkError(span, "\"toArray\" builtin takes a VarArg argument")
				}

				return &Array{Elems: varArgObj.Elems}
			},
		},
	},
	{
		Name: "contains",
		Builtin: &Builtin{
			Function: func(span token.Span, objects ...Object) Object {
				if len(objects) != 2 {
					return mkError(span, "\"contains\" builtin takes a HashMap argument and a key")
				}

				hashMapObj, ok := objects[0].(*HashMap)
				if !ok {
					return mkError(span, "First argument is not a hash map")
				}

				keyObj, ok := objects[1].(Hashable)
				if !ok {
					return mkError(span, "Second argument is not a hashable object")
				}

				elem, ok := hashMapObj.Elems[keyObj.HashKey()]
				if ok {
					if elem.Key.Inspect() == keyObj.Inspect() {
						return &Boolean{Value: true}
					}
				}

				return &Boolean{Value: false}
			},
		},
	},
}

func GetBuiltinByName(name string) *Builtin {
	for _, builtin := range Builtins {
		if builtin.Name == name {
			return builtin.Builtin
		}
	}
	return nil
}
