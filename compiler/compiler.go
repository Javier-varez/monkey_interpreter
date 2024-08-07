package compiler

import (
	"fmt"
	"sort"

	"github.com/javier-varez/monkey_interpreter/ast"
	"github.com/javier-varez/monkey_interpreter/code"
	"github.com/javier-varez/monkey_interpreter/object"
	"github.com/javier-varez/monkey_interpreter/token"
)

type CompilationScope struct {
	instructions code.Instructions

	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type Compiler struct {
	constants   []object.Object
	symbolTable *SymbolTable

	scopes   []CompilationScope
	curScope int
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func New() *Compiler {
	st := NewSymbolTable()
	for i, builtin := range object.Builtins {
		st.DefineBuiltin(i, builtin.Name)
	}

	return NewWithState([]object.Object{}, st)
}

func NewWithState(constants []object.Object, symbolTable *SymbolTable) *Compiler {
	mainScope := CompilationScope{}
	return &Compiler{
		constants:   constants,
		symbolTable: symbolTable,

		scopes: []CompilationScope{mainScope},
	}
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.curScope].instructions
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
		c.changeOperand(notTrutyInst, len(c.currentInstructions()))

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

		c.changeOperand(endTruthyJumpPos, len(c.currentInstructions()))
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
		if sym.Scope == LocalScope {
			c.emit(code.OpSetLocal, sym.Index)
		} else {
			c.emit(code.OpSetGlobal, sym.Index)
		}

	case *ast.IdentifierExpr:
		sym, ok := c.symbolTable.Resolve(node.IdentToken.Literal)
		if !ok {
			return fmt.Errorf("Unknown identifier %s", node.IdentToken.Literal)
		}

		switch sym.Scope {
		case LocalScope:
			c.emit(code.OpGetLocal, sym.Index)
		case GlobalScope:
			c.emit(code.OpGetGlobal, sym.Index)
		case BuiltinScope:
			c.emit(code.OpGetBuiltin, sym.Index)
		}

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

	case *ast.MapLiteralExpr:
		keys := []ast.Expression{}
		for k := range node.Map {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}

			err = c.Compile(node.Map[k])
			if err != nil {
				return err
			}
		}
		c.emit(code.OpHash, len(node.Map))

	case *ast.FnLiteralExpr:
		c.enterScope()
		// Define arguments
		for _, arg := range node.Args {
			c.symbolTable.Define(arg.IdentToken.Literal)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		if c.lastInstructionIsPop() {
			c.removeLastPop()
			c.emit(code.OpReturnValue)
		} else if !c.lastInstructionIsReturnValue() {
			c.emit(code.OpReturn)
		}

		insts, numLocals := c.exitScope()

		c.emit(code.OpConstant, c.addConstant(&object.CompiledFunction{
			Instructions: insts,
			NumLocals:    numLocals,
			NumArgs:      len(node.Args),
			VarArgs:      node.VarArgs,
		}))

	case *ast.ReturnStatement:
		err := c.Compile(node.Expr)
		if err != nil {
			return err
		}
		c.emit(code.OpReturnValue)

	case *ast.CallExpr:
		for _, arg := range node.Args {
			err := c.Compile(arg)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpConstant, c.addConstant(&object.Integer{Value: int64(len(node.Args))}))

		err := c.Compile(node.CallableExpr)
		if err != nil {
			return err
		}

		c.emit(code.OpCall)

	default:
		return fmt.Errorf("Unhandled node type %T", untypedNode)
	}

	return nil
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.scopes[c.curScope].lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) lastInstructionIsReturnValue() bool {
	return c.scopes[c.curScope].lastInstruction.Opcode == code.OpReturnValue
}

func (c *Compiler) removeLastPop() {
	if !c.lastInstructionIsPop() {
		panic("Last instruction was not pop")
	}
	c.scopes[c.curScope].instructions = c.scopes[c.curScope].instructions[:c.scopes[c.curScope].lastInstruction.Position]
	c.scopes[c.curScope].lastInstruction = c.scopes[c.curScope].previousInstruction
}

func (c *Compiler) replaceInstruction(pos int, newInstr []byte) {
	for i := range newInstr {
		c.scopes[c.curScope].instructions[pos+i] = newInstr[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	opcode := code.Opcode(c.scopes[c.curScope].instructions[opPos])
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
	currentInsts := c.currentInstructions()
	off := len(currentInsts)
	c.scopes[c.curScope].instructions = append(currentInsts, inst...)
	return off
}

func (c *Compiler) setLastInstruction(op code.Opcode, position int) {
	c.scopes[c.curScope].previousInstruction = c.scopes[c.curScope].lastInstruction
	c.scopes[c.curScope].lastInstruction = EmittedInstruction{
		Opcode:   op,
		Position: position,
	}
}

func (c *Compiler) enterScope() {
	c.scopes = append(c.scopes, CompilationScope{})
	c.curScope++

	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) exitScope() (code.Instructions, int) {
	insts := c.currentInstructions()
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.curScope--

	numLocals := c.symbolTable.NumDefinitions
	c.symbolTable = c.symbolTable.Parent

	return insts, numLocals
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
