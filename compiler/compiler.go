package compiler

import (
	"fmt"

	"github.com/javier-varez/monkey_interpreter/ast"
	"github.com/javier-varez/monkey_interpreter/code"
	"github.com/javier-varez/monkey_interpreter/object"
	"github.com/javier-varez/monkey_interpreter/token"
)

type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expr)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
		return nil

	case *ast.InfixExpr:
		if node.OperatorToken.Type == token.LT {
			err := c.Compile(node.RightExpr)
			if err != nil {
				return err
			}

			err = c.Compile(node.LeftExpr)
			if err != nil {
				return err
			}

			c.emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.LeftExpr)
		if err != nil {
			return err
		}

		err = c.Compile(node.RightExpr)
		if err != nil {
			return err
		}

		switch node.OperatorToken.Type {
		case token.PLUS:
			c.emit(code.OpAdd)
		case token.MINUS:
			c.emit(code.OpSub)
		case token.ASTERISK:
			c.emit(code.OpMul)
		case token.SLASH:
			c.emit(code.OpDiv)
		case token.GT:
			c.emit(code.OpGreaterThan)
		case token.EQ:
			c.emit(code.OpEqual)
		case token.NOT_EQ:
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("Unhandled infix operator %s", node.OperatorToken.Type)
		}

		return nil

	case *ast.IntegerLiteralExpr:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.BoolLiteralExpr:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.PrefixExpr:
		err := c.Compile(node.InnerExpr)
		if err != nil {
			return err
		}

		switch node.OperatorToken.Type {
		case token.BANG:
			c.emit(code.OpBang)
		case token.MINUS:
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("Unhandled prefix operator %s", node.OperatorToken.Type)
		}
	}

	return nil
}

/// Operands is a list of operand offsets to the constants of the compiler
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	inst := code.Make(op, operands...)
	pos := c.addInstruction(inst)
	return pos
}

func (c *Compiler) addConstant(obj object.Object) int {
	off := len(c.constants)
	c.constants = append(c.constants, obj)
	return off
}

func (c *Compiler) addInstruction(inst code.Instructions) int {
	off := len(c.instructions)
	c.instructions = append(c.instructions, inst...)
	return off
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
