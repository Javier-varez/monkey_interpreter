package vm

import (
	"github.com/javier-varez/monkey_interpreter/code"
	"github.com/javier-varez/monkey_interpreter/object"
)

type Frame struct {
	fn *object.CompiledFunction
	ip int

	// points to the local vars on the stack
	LocalsBase int
}

func NewFrame(fn *object.CompiledFunction, sp int) *Frame {
	return &Frame{
		fn: fn,
		ip: -1,

		LocalsBase: sp,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
