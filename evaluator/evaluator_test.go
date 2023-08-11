package evaluator

import (
	"testing"

	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/object"
	"github.com/javier-varez/monkey_interpreter/parser"
)

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	intRes, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Result is not an integral object: %v", obj)
		return false
	}

	if intRes.Value != expected {
		t.Errorf("Unexpected value: expected %d, got %d", expected, intRes.Value)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	boolRes, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Result is not a boolean object: %v", obj)
		return false
	}

	if boolRes.Value != expected {
		t.Errorf("Unexpected value: expected %t, got %t", expected, boolRes.Value)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.Null)
	if !ok {
		t.Errorf("Result is not a null object: %v", obj)
		return false
	}

	return true
}

func testObject(t *testing.T, obj object.Object, inner interface{}) bool {
	switch inner := inner.(type) {
	case int64:
		return testIntegerObject(t, obj, inner)
	case bool:
		return testBooleanObject(t, obj, inner)
	case nil:
		return testNullObject(t, obj)
	default:
		panic("Unhandled type in testObject")
	}

	return false
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	return Eval(program)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input  string
		output int64
	}{
		{"5", 5},
		{"123", 123},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testIntegerObject(t, result, tt.output)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.output)
	}
}

func TestEvalBangOperator(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{"!false", true},
		{"!true", false},
		{"!!true", true},
		{"!!false", false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.output)
	}
}

func TestEvalMinuxOperator(t *testing.T) {
	tests := []struct {
		input  string
		output int64
	}{
		{"--123", 123},
		{"-123", -123},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testIntegerObject(t, result, tt.output)
	}
}

func TestEvalInfixOperators(t *testing.T) {
	tests := []struct {
		input  string
		output int64
	}{
		{"123 + 123", 246},
		{"12 * 2", 24},
		{"16 / 2", 8},
		{"16 - 2", 14},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testIntegerObject(t, result, tt.output)
	}
}

func TestEvalInfixBoolOperators(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{"123 == 123", true},
		{"123 == 124", false},
		{"123 != 123", false},
		{"123 != 124", true},
		{"true == true", true},
		{"true == false", false},
		{"true != true", false},
		{"true != false", true},
		{"12 < 123", true},
		{"12 > 123", false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.output)
	}
}

func TestEvalIfExpression(t *testing.T) {
	tests := []struct {
		input  string
		output int64
	}{
		{"if (100 > 200) { 1 } else { 2 }", 2},
		{"if (100 > 20) { 1 } else { 2 }", 1},
		{"if (100 > 20) { 1 }", 1},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testIntegerObject(t, result, tt.output)
	}
}

func TestEvalIfExpressionWithoutAlternative(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"if (100 > 200) { 1 }"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testNullObject(t, result)
	}
}

func TestEvalReturnStatements(t *testing.T) {
	tests := []struct {
		input  string
		output interface{}
	}{
		{"return 10", int64(10)},
		{"return true", true},
		{"return", nil},
		{"return 10; 2;", int64(10)},
		{"return false; 2;", false},
		{"return; 2;", nil},
		{"if (100 < 200) { 2 * 2; return 33; 22; }; 2;", int64(33)},
		{"if (200 < 200) { 2 * 2; return 33; 22; }; 2;", int64(2)},
		{"if (100 < 200) { 2 * 2; if (1 != 2) { return 33; }; return 22; }; 2;", int64(33)},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testObject(t, result, tt.output)
	}
}
