package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected Instructions
	}{
		{OpConstant, []int{65534}, Instructions{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, Instructions{byte(OpAdd)}},
		{OpSub, []int{}, Instructions{byte(OpSub)}},
		{OpMul, []int{}, Instructions{byte(OpMul)}},
		{OpDiv, []int{}, Instructions{byte(OpDiv)}},
		{OpPop, []int{}, Instructions{byte(OpPop)}},
		{OpTrue, []int{}, Instructions{byte(OpTrue)}},
		{OpFalse, []int{}, Instructions{byte(OpFalse)}},
		{OpEqual, []int{}, Instructions{byte(OpEqual)}},
		{OpNotEqual, []int{}, Instructions{byte(OpNotEqual)}},
		{OpGreaterThan, []int{}, Instructions{byte(OpGreaterThan)}},
		{OpMinus, []int{}, Instructions{byte(OpMinus)}},
		{OpBang, []int{}, Instructions{byte(OpBang)}},
		{OpJump, []int{123}, Instructions{byte(OpJump), 0, 123}},
		{OpJumpNotTruthy, []int{126}, Instructions{byte(OpJumpNotTruthy), 0, 126}},
		{OpNull, []int{}, Instructions{byte(OpNull)}},
		{OpGetGlobal, []int{123}, Instructions{byte(OpGetGlobal), 0, 123}},
		{OpSetGlobal, []int{123}, Instructions{byte(OpSetGlobal), 0, 123}},
		{OpArray, []int{2}, Instructions{byte(OpArray), 0, 2}},
		{OpIndex, []int{}, Instructions{byte(OpIndex)}},
		{OpHash, []int{2}, Instructions{byte(OpHash), 0, 2}},
		{OpCall, []int{}, Instructions{byte(OpCall)}},
		{OpReturn, []int{}, Instructions{byte(OpReturn)}},
		{OpReturnValue, []int{}, Instructions{byte(OpReturnValue)}},
		{OpGetLocal, []int{254}, Instructions{byte(OpGetLocal), 254}},
		{OpSetLocal, []int{254}, Instructions{byte(OpSetLocal), 254}},
		{OpGetBuiltin, []int{254}, Instructions{byte(OpGetBuiltin), 254}},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("Generated instruction has the wrong length: want=%d, got=%d", len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != b {
				t.Errorf("Wrong byte at pos %d: want=%d, got=%d", i, b, instruction[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpAdd),
		Make(OpSub),
		Make(OpMul),
		Make(OpDiv),
		Make(OpPop),
		Make(OpTrue),
		Make(OpFalse),
		Make(OpEqual),
		Make(OpNotEqual),
		Make(OpGreaterThan),
		Make(OpMinus),
		Make(OpBang),
		Make(OpJump, 1234),
		Make(OpJumpNotTruthy, 1234),
		Make(OpNull),
		Make(OpGetGlobal, 1234),
		Make(OpSetGlobal, 1234),
		Make(OpArray, 2),
		Make(OpIndex),
		Make(OpHash, 3),
		Make(OpCall),
		Make(OpReturn),
		Make(OpReturnValue),
		Make(OpGetLocal, 254),
		Make(OpSetLocal, 254),
		Make(OpGetBuiltin, 254),
	}

	expected := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpAdd
0010 OpSub
0011 OpMul
0012 OpDiv
0013 OpPop
0014 OpTrue
0015 OpFalse
0016 OpEqual
0017 OpNotEqual
0018 OpGreaterThan
0019 OpMinus
0020 OpBang
0021 OpJump 1234
0024 OpJumpNotTruthy 1234
0027 OpNull
0028 OpGetGlobal 1234
0031 OpSetGlobal 1234
0034 OpArray 2
0037 OpIndex
0038 OpHash 3
0041 OpCall
0042 OpReturn
0043 OpReturnValue
0044 OpGetLocal 254
0046 OpSetLocal 254
0048 OpGetBuiltin 254
`
	concatted := Instructions{}
	for _, ins := range instructions {
		concatted = append(concatted, ins...)
	}

	asString := concatted.String()
	if asString != expected {
		t.Errorf("instructions wrongly formatted.\nwant=%q\ngot=%q", expected, asString)
	}

}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. want=%q, got=%q", tt.bytesRead, n)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong. want=%d, got=%d", want, operandsRead[i])
			}
		}
	}
}
