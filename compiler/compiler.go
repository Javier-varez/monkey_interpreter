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

	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	symbolTable         *SymbolTable
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func New() *Compiler {
	return NewWithState([]object.Object{}, NewSymbolTable())
}

func NewWithState(constants []object.Object, symbolTable *SymbolTable) *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    constants,
		symbolTable:  symbolTable,
	}
}

func (c *Compiler) Compile(untypedNode ast.Node) error {
	switch node := untypedNode.(type) {
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

	case *ast.IfExpr:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		notTrutyInst := c.emit(code.OpJumpNotTruthy, 1234)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		endTruthyJumpPos := c.emit(code.OpJump, 1234)
		c.changeOperand(notTrutyInst, len(c.instructions))

		if node.Alternative != nil {
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		} else {
			c.emit(code.OpNull)
		}

		c.changeOperand(endTruthyJumpPos, len(c.instructions))
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.LetStatement:
		err := c.Compile(node.Expr)
		if err != nil {
			return err
		}

		sym := c.symbolTable.Define(node.IdentExpr.(*ast.IdentifierExpr).IdentToken.Literal)
		c.emit(code.OpSetGlobal, sym.Index)

	case *ast.IdentifierExpr:
		sym, ok := c.symbolTable.Resolve(node.IdentToken.Literal)
		if !ok {
			return fmt.Errorf("Unknown identifier %s", node.IdentToken.Literal)
		}

		c.emit(code.OpGetGlobal, sym.Index)

	case *ast.StringLiteralExpr:
		idx := c.addConstant(&object.String{Value: node.Value})
		c.emit(code.OpConstant, idx)

	case *ast.ArrayLiteralExpr:
		for _, elemExpr := range node.Elems {
			err := c.Compile(elemExpr)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elems))

	case *ast.IndexOperatorExpr:
		err := c.Compile(node.ObjExpr)
		if err != nil {
			return err
		}

		err = c.Compile(node.IndexExpr)
		if err != nil {
			return err
		}

		c.emit(code.OpIndex)

	default:
		return fmt.Errorf("Unhandled node type %T", untypedNode)
	}

	return nil
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstr []byte) {
	for i := range newInstr {
		c.instructions[pos+i] = newInstr[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	opcode := code.Opcode(c.instructions[opPos])
	instrs := code.Make(opcode, operand)
	c.replaceInstruction(opPos, instrs)
}

/// Operands is a list of operand offsets to the constants of the compiler
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	inst := code.Make(op, operands...)
	pos := c.addInstruction(inst)
	c.setLastInstruction(op, pos)
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

func (c *Compiler) setLastInstruction(op code.Opcode, position int) {
	c.previousInstruction = c.lastInstruction
	c.lastInstruction = EmittedInstruction{
		Opcode:   op,
		Position: position,
	}
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
