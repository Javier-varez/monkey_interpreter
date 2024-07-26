package vm

import (
	"fmt"

	"github.com/javier-varez/monkey_interpreter/code"
	"github.com/javier-varez/monkey_interpreter/compiler"
	"github.com/javier-varez/monkey_interpreter/object"
)

const STACK_SIZE = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // Always points to the next value. Top of stack is stack[sp-1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,

		stack: make([]object.Object, STACK_SIZE),
	}
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := vm.runBinaryOp(op)
			if err != nil {
				return err
			}
		case code.OpPop:
			_, err := vm.pop()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unhandled operation: %v", op)
		}
	}
	return nil
}

func (vm *VM) runBinaryOp(op code.Opcode) error {
	rhs, err := vm.pop()
	if err != nil {
		return err
	}

	lhs, err := vm.pop()
	if err != nil {
		return err
	}

	rhsVal := rhs.(*object.Integer).Value
	lhsVal := lhs.(*object.Integer).Value

	var result int64
	switch op {
	case code.OpAdd:
		result = lhsVal + rhsVal
	case code.OpSub:
		result = lhsVal - rhsVal
	case code.OpMul:
		result = lhsVal * rhsVal
	case code.OpDiv:
		result = lhsVal / rhsVal
	default:
		return fmt.Errorf("Invalid binary operation: %v", op)
	}
	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) push(ob object.Object) error {
	if vm.sp >= len(vm.stack) {
		return fmt.Errorf("Stack overflown")
	}
	vm.stack[vm.sp] = ob
	vm.sp++
	return nil
}

func (vm *VM) pop() (object.Object, error) {
	if vm.sp == 0 {
		return nil, fmt.Errorf("Stack underflown")
	}
	obj := vm.stack[vm.sp-1]
	vm.sp--
	return obj, nil
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}
