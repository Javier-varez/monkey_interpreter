package vm

import (
	"fmt"

	"github.com/javier-varez/monkey_interpreter/code"
	"github.com/javier-varez/monkey_interpreter/compiler"
	"github.com/javier-varez/monkey_interpreter/object"
)

const STACK_SIZE = 2048
const GLOBALS_SIZE = 65536

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack   []object.Object
	sp      int // Always points to the next value. Top of stack is stack[sp-1]
	globals []object.Object
}

var Null = &object.Null{}
var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

func assertNotNil(obj any) {
	if obj == nil {
		panic("Unexpected nil object")
	}
}

func New(bytecode *compiler.Bytecode) *VM {
	return NewWithGlobalKeyStore(bytecode, make([]object.Object, GLOBALS_SIZE))
}

func NewWithGlobalKeyStore(bytecode *compiler.Bytecode, keyStore []object.Object) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,

		stack:   make([]object.Object, STACK_SIZE),
		globals: keyStore,
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
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
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
			} else if v.Type() == object.NULL_OBJ {
				asBool = false
			} else {
				asBool = true
			}

			vm.push(&object.Boolean{Value: !asBool})

		case code.OpJumpNotTruthy:
			target := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			v, err := vm.pop()
			if err != nil {
				return err
			}

			if isTruthy := asBoolean(v); !isTruthy {
				ip = int(target) - 1
			}

		case code.OpJump:
			target := code.ReadUint16(vm.instructions[ip+1:])
			ip = int(target) - 1

		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.OpSetGlobal:
			idx := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			obj, err := vm.pop()
			if err != nil {
				return err
			}
			vm.globals[idx] = obj

		case code.OpGetGlobal:
			idx := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			obj := vm.globals[idx]
			assertNotNil(obj)
			err := vm.push(obj)
			if err != nil {
				return err
			}

		case code.OpArray:
			arrayLen := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			arr := &object.Array{Elems: make([]object.Object, arrayLen)}
			for i := int(arrayLen) - 1; i >= 0; i-- {
				val, err := vm.pop()
				if err != nil {
					return err
				}
				arr.Elems[i] = val
			}

			err := vm.push(arr)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("Unhandled operation: %v", op)
		}
	}
	return nil
}

func asBoolean(o object.Object) bool {
	switch o := o.(type) {
	case *object.Integer:
		return o.Value != 0
	case *object.Boolean:
		return o.Value
	case *object.Null:
		return false
	}
	return true
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

	if rhs.Type() == object.STRING_OBJ && lhs.Type() == object.STRING_OBJ {
		return vm.runStringBinaryOp(op, lhs.(*object.String), rhs.(*object.String))
	}
	if rhs.Type() == object.INTEGER_OBJ && lhs.Type() == object.INTEGER_OBJ {
		return vm.runIntBinaryOp(op, lhs.(*object.Integer), rhs.(*object.Integer))
	}
	return fmt.Errorf("Invalid binary operation %d for types %T and %T", op, lhs, rhs)
}

func (vm *VM) runIntBinaryOp(op code.Opcode, lhs, rhs *object.Integer) error {
	var result int64
	switch op {
	case code.OpAdd:
		result = lhs.Value + rhs.Value
	case code.OpSub:
		result = lhs.Value - rhs.Value
	case code.OpMul:
		result = lhs.Value * rhs.Value
	case code.OpDiv:
		result = lhs.Value / rhs.Value
	default:
		return fmt.Errorf("Invalid binary operation: %v", op)
	}
	return vm.push(&object.Integer{Value: result})
}

func (vm *VM) runStringBinaryOp(op code.Opcode, lhs, rhs *object.String) error {
	var result string
	switch op {
	case code.OpAdd:
		result = lhs.Value + rhs.Value
	default:
		return fmt.Errorf("Invalid string binary operation: %v", op)
	}
	return vm.push(&object.String{Value: result})
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
