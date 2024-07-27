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
		case code.OpTrue:
			err := vm.push(&object.Boolean{Value: true})
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(&object.Boolean{Value: false})
			if err != nil {
				return err
			}
		case code.OpGreaterThan, code.OpEqual, code.OpNotEqual:
			err := vm.runComparisonOp(op)
			if err != nil {
				return err
			}
		case code.OpMinus:
			v, err := vm.pop()
			if err != nil {
				return err
			}

			if v.Type() != object.INTEGER_OBJ {
				return fmt.Errorf("Cannot apply minus operator on type %T", v)
			}

			asInt := v.(*object.Integer)
			vm.push(&object.Integer{Value: -asInt.Value})
		case code.OpBang:
			v, err := vm.pop()
			if err != nil {
				return err
			}

			var asBool bool
			if v.Type() == object.BOOLEAN_OBJ {
				asBool = v.(*object.Boolean).Value
			} else if v.Type() == object.INTEGER_OBJ {
				asBool = v.(*object.Integer).Value != 0
			} else {
				return fmt.Errorf("Cannot apply bang operator on type %T", v)
			}

			vm.push(&object.Boolean{Value: !asBool})
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

func (vm *VM) runComparisonOp(op code.Opcode) error {
	rhs, err := vm.pop()
	if err != nil {
		return err
	}

	lhs, err := vm.pop()
	if err != nil {
		return err
	}

	if rhs.Type() == object.INTEGER_OBJ && lhs.Type() == object.INTEGER_OBJ {
		return vm.runIntComparisonOp(op, lhs.(*object.Integer), rhs.(*object.Integer))
	}

	if rhs.Type() != object.BOOLEAN_OBJ || lhs.Type() != object.BOOLEAN_OBJ {
		return fmt.Errorf("Cannot apply comparison operator on types %T and %T", lhs, rhs)
	}

	lhsVal := lhs.(*object.Boolean)
	rhsVal := rhs.(*object.Boolean)

	var result bool
	switch op {
	case code.OpEqual:
		result = lhsVal.Value == rhsVal.Value
	case code.OpNotEqual:
		result = lhsVal.Value != rhsVal.Value
	default:
		return fmt.Errorf("Invalid comparison operation: %v", op)
	}
	return vm.push(&object.Boolean{Value: result})
}

func (vm *VM) runIntComparisonOp(op code.Opcode, lhs, rhs *object.Integer) error {
	var result bool
	switch op {
	case code.OpEqual:
		result = lhs.Value == rhs.Value
	case code.OpNotEqual:
		result = lhs.Value != rhs.Value
	case code.OpGreaterThan:
		result = lhs.Value > rhs.Value
	default:
		return fmt.Errorf("Invalid integer comparison operation: %v", op)
	}
	return vm.push(&object.Boolean{Value: result})
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
