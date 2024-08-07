package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPop
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump
	OpNull
	OpSetGlobal
	OpGetGlobal
	OpArray
	OpIndex
	OpHash
	OpCall
	OpReturn
	OpReturnValue
	OpSetLocal
	OpGetLocal
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:      {Name: "OpConstant", OperandWidths: []int{2}},
	OpAdd:           {Name: "OpAdd"},
	OpSub:           {Name: "OpSub"},
	OpMul:           {Name: "OpMul"},
	OpDiv:           {Name: "OpDiv"},
	OpPop:           {Name: "OpPop"},
	OpTrue:          {Name: "OpTrue"},
	OpFalse:         {Name: "OpFalse"},
	OpEqual:         {Name: "OpEqual"},
	OpNotEqual:      {Name: "OpNotEqual"},
	OpGreaterThan:   {Name: "OpGreaterThan"},
	OpMinus:         {Name: "OpMinus"},
	OpBang:          {Name: "OpBang"},
	OpJumpNotTruthy: {Name: "OpJumpNotTruthy", OperandWidths: []int{2}},
	OpJump:          {Name: "OpJump", OperandWidths: []int{2}},
	OpNull:          {Name: "OpNull"},
	OpSetGlobal:     {Name: "OpSetGlobal", OperandWidths: []int{2}},
	OpGetGlobal:     {Name: "OpGetGlobal", OperandWidths: []int{2}},
	OpArray:         {Name: "OpArray", OperandWidths: []int{2}},
	OpIndex:         {Name: "OpIndex"},
	OpHash:          {Name: "OpHash", OperandWidths: []int{2}},
	OpCall:          {Name: "OpCall"},
	OpReturnValue:   {Name: "OpReturnValue"},
	OpReturn:        {Name: "OpReturn"},
	OpSetLocal:      {Name: "OpSetLocal", OperandWidths: []int{1}},
	OpGetLocal:      {Name: "OpGetLocal", OperandWidths: []int{1}},
}

func Lookup(op byte) (*Definition, error) {
	if def, ok := definitions[Opcode(op)]; ok {
		return def, nil
	}
	return nil, fmt.Errorf("Unknown opcode %d", op)
}

func Make(opcode Opcode, operands ...int) []byte {
	params, ok := definitions[opcode]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range params.OperandWidths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(opcode)

	assertInBounds := func(value int, maxBound int) {
		if value > maxBound {
			panic(fmt.Sprintf("Value %d is out of bounds (max=%d)", value, maxBound))
		}
	}

	offset := 1
	for i, o := range operands {
		width := params.OperandWidths[i]
		switch width {
		case 1:
			assertInBounds(o, 0xFF)
			instruction[offset] = byte(o)
		case 2:
			assertInBounds(o, 0xFFFF)
			binary.BigEndian.PutUint16(instruction[offset:offset+2], uint16(o))
		}
		offset += width
	}

	return instruction
}

func (inst Instructions) String() string {
	var out bytes.Buffer
	offset := 0

	for offset < len(inst) {
		def, err := Lookup(inst[offset])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			offset += 1
			continue
		}

		ops, consumed := ReadOperands(def, inst[offset+1:])
		fmt.Fprintf(&out, "%04d %s\n", offset, fmtInstruction(def, ops))
		offset += consumed + 1
	}

	return out.String()
}

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 1:
			operands[i] = int(ReadUint8(ins[offset : offset+1]))
		case 2:
			operands[i] = int(ReadUint16(ins[offset : offset+2]))
		}
		offset += width
	}

	return operands, offset
}

func ReadUint8(ins Instructions) byte {
	return ins[0]
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
