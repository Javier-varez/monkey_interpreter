package vm

import (
	"fmt"

	"github.com/javier-varez/monkey_interpreter/code"
	"github.com/javier-varez/monkey_interpreter/compiler"
	"github.com/javier-varez/monkey_interpreter/object"
	"github.com/javier-varez/monkey_interpreter/token"
)

const STACK_SIZE = 2048
const GLOBALS_SIZE = 65536
const MAX_FRAMES = 1024

type VM struct {
	constants []object.Object

	stack   []object.Object
	sp      int // Always points to the next value. Top of stack is stack[sp-1]
	globals []object.Object

	frames     []*Frame
	frameIndex int
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
	frames := make([]*Frame, MAX_FRAMES)
	frames[0] = NewFrame(&object.Closure{Fn: &object.CompiledFunction{Instructions: bytecode.Instructions}}, 0)

	return &VM{
		constants: bytecode.Constants,

		stack:   make([]object.Object, STACK_SIZE),
		globals: keyStore,

		frames:     frames,
		frameIndex: 1,
	}
}

func (vm *VM) Run() error {
	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip := vm.currentFrame().ip
		inst := vm.currentFrame().Instructions()
		op := code.Opcode(inst[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(inst[ip+1:])
			vm.currentFrame().ip += 2
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
			target := code.ReadUint16(inst[ip+1:])
			vm.currentFrame().ip += 2

			v, err := vm.pop()
			if err != nil {
				return err
			}

			if isTruthy := asBoolean(v); !isTruthy {
				vm.currentFrame().ip = int(target) - 1
			}

		case code.OpJump:
			target := code.ReadUint16(inst[ip+1:])
			vm.currentFrame().ip = int(target) - 1

		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.OpSetGlobal:
			idx := code.ReadUint16(inst[ip+1:])
			vm.currentFrame().ip += 2

			obj, err := vm.pop()
			if err != nil {
				return err
			}
			vm.globals[idx] = obj

		case code.OpGetGlobal:
			idx := code.ReadUint16(inst[ip+1:])
			vm.currentFrame().ip += 2

			obj := vm.globals[idx]
			assertNotNil(obj)
			err := vm.push(obj)
			if err != nil {
				return err
			}

		case code.OpArray:
			arrayLen := code.ReadUint16(inst[ip+1:])
			vm.currentFrame().ip += 2

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

		case code.OpHash:
			mapLen := code.ReadUint16(inst[ip+1:])
			vm.currentFrame().ip += 2

			hashmap := &object.HashMap{Elems: make(map[object.HashKey]object.HashEntry, mapLen)}
			for i := 0; i < int(mapLen); i++ {
				val, err := vm.pop()
				if err != nil {
					return err
				}

				key, err := vm.pop()
				if err != nil {
					return err
				}

				hashable, ok := key.(object.Hashable)
				if !ok {
					return fmt.Errorf("Key object is not hashable")
				}

				hashmap.Elems[hashable.HashKey()] = object.HashEntry{Key: key, Value: val}
			}

			err := vm.push(hashmap)
			if err != nil {
				return err
			}

		case code.OpIndex:
			indexObj, err := vm.pop()
			if err != nil {
				return err
			}

			indexedObj, err := vm.pop()
			if err != nil {
				return err
			}

			switch inner := indexedObj.(type) {
			case *object.Array:
				if indexObj.Type() != object.INTEGER_OBJ {
					return fmt.Errorf("Index to array must be an integral. Got=%T (%+v)", indexObj, indexObj)
				}

				var err error
				i := indexObj.(*object.Integer).Value
				if i < int64(len(inner.Elems)) {
					err = vm.push(inner.Elems[i])
				} else {
					err = vm.push(Null)
				}
				if err != nil {
					return err
				}

			case *object.HashMap:
				hashable, ok := indexObj.(object.Hashable)
				if !ok {
					return fmt.Errorf("Index of type %T (%+v) is not hashable", indexObj, indexObj)
				}
				key := hashable.HashKey()

				kv, ok := inner.Elems[key]
				var err error
				if ok {
					err = vm.push(kv.Value)
				} else {
					err = vm.push(Null)
				}
				if err != nil {
					return err
				}

			default:
				return fmt.Errorf("Cannot index object of type: %T", indexedObj)
			}

		case code.OpCall:
			fnObj, err := vm.pop()
			if err != nil {
				return err
			}

			numArgsObj, err := vm.pop()
			if err != nil {
				return err
			}

			numArgs, ok := numArgsObj.(*object.Integer)
			if !ok {
				return fmt.Errorf("Could not get number of arguments to function in the stack")
			}
			numArgsInCall := int(numArgs.Value)

			switch fn := fnObj.(type) {
			case *object.Closure:
				if err := vm.callCompiledFunction(fn, numArgsInCall); err != nil {
					return err
				}

			case *object.Builtin:
				if err := vm.executeBuiltin(fn, numArgsInCall); err != nil {
					return err
				}

			default:
				return fmt.Errorf("Not a callable, cannot be invoked")
			}

		case code.OpReturn:
			vm.popFrame()
			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.OpReturnValue:
			val, err := vm.pop()
			if err != nil {
				return err
			}

			vm.popFrame()

			err = vm.push(val)
			if err != nil {
				return err
			}

		case code.OpSetLocal:
			frame := vm.currentFrame()
			idx := frame.LocalsBase + int(code.ReadUint8(inst[ip+1:]))
			frame.ip += 1

			obj, err := vm.pop()
			if err != nil {
				return err
			}

			// TODO: This does not handle correctly accessing locals from a parent scope
			vm.stack[idx] = obj

		case code.OpGetLocal:
			frame := vm.currentFrame()
			idx := frame.LocalsBase + int(code.ReadUint8(inst[ip+1:]))
			frame.ip += 1

			// TODO: This does not handle correctly accessing locals from a parent scope
			obj := vm.stack[idx]
			assertNotNil(obj)
			err := vm.push(obj)
			if err != nil {
				return err
			}

		case code.OpGetBuiltin:
			idx := int(code.ReadUint8(inst[ip+1:]))
			vm.currentFrame().ip += 1

			if idx >= len(object.Builtins) {
				panic(fmt.Sprintf("Unknown builtin index %d", idx))
			}

			obj := object.Builtins[idx].Builtin
			err := vm.push(obj)
			if err != nil {
				return err
			}

		case code.OpClosure:
			constantIdx := int(code.ReadUint16(inst[ip+1:]))
			numFreeVars := int(code.ReadUint8(inst[ip+3:]))
			vm.currentFrame().ip += 3

			freeObjects := make([]object.Object, numFreeVars)
			for i := 0; i < numFreeVars; i++ {
				val, err := vm.pop()
				if err != nil {
					return err
				}
				freeObjects[numFreeVars-1-i] = val
			}

			fn, ok := vm.constants[constantIdx].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("Argument to the OpClosure is not a compiled function")
			}

			err := vm.push(&object.Closure{Fn: fn, FreeObjects: freeObjects})
			if err != nil {
				return err
			}

		case code.OpGetFree:
			freeIdx := int(code.ReadUint8(inst[ip+1:]))
			vm.currentFrame().ip += 1
			freeObjects := vm.currentFrame().closure.FreeObjects

			if freeIdx >= len(freeObjects) {
				return fmt.Errorf("Invalid free index: %d. Num free objects: %d", freeIdx, len(freeObjects))
			}

			if err := vm.push(freeObjects[freeIdx]); err != nil {
				return err
			}

		case code.OpRange:
			endObj, err := vm.pop()
			if err != nil {
				return err
			}

			if endObj.Type() != object.INTEGER_OBJ {
				return fmt.Errorf("Range start does not evaluate to an integer object: %T (%V)", endObj, endObj)
			}

			startObj, err := vm.pop()
			if err != nil {
				return err
			}

			if startObj.Type() != object.INTEGER_OBJ {
				return fmt.Errorf("Range start does not evaluate to an integer object: %T (%V)", startObj, startObj)
			}

			start := startObj.(*object.Integer).Value
			end := endObj.(*object.Integer).Value

			incr := int64(1)
			if start > end {
				// Decreasing range
				incr = -1
			}

			arrayObj := &object.Array{Elems: []object.Object{}}

			curValue := start
			for curValue != end {
				arrayObj.Elems = append(arrayObj.Elems, &object.Integer{Value: curValue})
				curValue = curValue + incr
			}

			err = vm.push(arrayObj)
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

func (vm *VM) callCompiledFunction(closure *object.Closure, numArgsInCall int) error {
	allArgs := make([]object.Object, numArgsInCall)
	copy(allArgs, vm.stack[vm.sp-numArgsInCall:vm.sp])

	// Free all the stack objects of the call, we will push them back now
	vm.sp -= numArgsInCall

	numArgsInCall = 0
	for _, arg := range allArgs {
		switch typedArg := arg.(type) {
		case *object.VarArgs:
			numArgsInCall += len(typedArg.Elems)
			for _, inner := range typedArg.Elems {
				if err := vm.push(inner); err != nil {
					return err
				}
			}
		default:
			numArgsInCall += 1
			if err := vm.push(arg); err != nil {
				return err
			}
		}
	}

	fn := closure.Fn
	if fn.VarArgs {
		if numArgsInCall < fn.NumArgs {
			return fmt.Errorf("wrong number of arguments: want>=%d, got=%d", fn.NumArgs, numArgsInCall)
		}

		numVarArgs := numArgsInCall - fn.NumArgs
		varArgs := &object.VarArgs{Elems: make([]object.Object, numVarArgs)}
		for i := 0; i < numVarArgs; i++ {
			v, err := vm.pop()
			if err != nil {
				return err
			}
			varArgs.Elems[numVarArgs-1-i] = v
		}
		err := vm.push(varArgs)
		if err != nil {
			return err
		}
	} else {
		if numArgsInCall != fn.NumArgs {
			return fmt.Errorf("wrong number of arguments: want=%d, got=%d", fn.NumArgs, numArgsInCall)
		}
	}

	vm.pushFrame(NewFrame(closure, vm.sp))
	return nil
}

func (vm *VM) executeBuiltin(fn *object.Builtin, numArgsInCall int) error {
	args := make([]object.Object, numArgsInCall)
	for i := 0; i < int(numArgsInCall); i++ {
		v, err := vm.pop()
		if err != nil {
			return err
		}
		args[int(numArgsInCall)-1-i] = v
	}

	val := fn.Function(token.Span{}, args...)
	if val == nil {
		val = Null
	}

	return vm.push(val)
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

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.frameIndex-1]
}

func (vm *VM) pushFrame(frame *Frame) {
	numArgs := frame.closure.Fn.NumArgs
	if frame.closure.Fn.VarArgs {
		numArgs += 1
	}
	// Arg locals have already been pushed to the stack, therefore we don't need to move them
	vm.sp += frame.closure.Fn.NumLocals - numArgs
	vm.frames[vm.frameIndex] = frame
	vm.frameIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.frameIndex--
	frame := vm.frames[vm.frameIndex]
	vm.sp -= frame.closure.Fn.NumLocals
	return frame
}
