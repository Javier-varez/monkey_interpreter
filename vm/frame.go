package vm

import (
	"github.com/javier-varez/monkey_interpreter/code"
	"github.com/javier-varez/monkey_interpreter/object"
)

type Frame struct {
	closure *object.Closure
	ip      int

	// points to the local vars on the stack
	LocalsBase int
}

func NewFrame(closure *object.Closure, sp int) *Frame {
	args := closure.Fn.NumArgs
	if closure.Fn.VarArgs {
		args += 1
	}
	return &Frame{
		closure: closure,
		ip:      -1,

		LocalsBase: sp - args,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.closure.Fn.Instructions
}
