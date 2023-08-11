package evaluator

import (
	"log"

	"github.com/javier-varez/monkey_interpreter/ast"
	"github.com/javier-varez/monkey_interpreter/object"
	"github.com/javier-varez/monkey_interpreter/token"
)

func evalAdd(leftObject, rightObject object.Object) object.Object {
	// TODO(ja): Handle other types
	left := leftObject.(*object.Integer)
	right := rightObject.(*object.Integer)

	return &object.Integer{
		Value: left.Value + right.Value,
	}
}

func evalSub(leftObject, rightObject object.Object) object.Object {
	// TODO(ja): Handle other types
	left := leftObject.(*object.Integer)
	right := rightObject.(*object.Integer)

	return &object.Integer{
		Value: left.Value - right.Value,
	}
}

func evalMul(leftObject, rightObject object.Object) object.Object {
	// TODO(ja): Handle other types
	left := leftObject.(*object.Integer)
	right := rightObject.(*object.Integer)

	return &object.Integer{
		Value: left.Value * right.Value,
	}
}

func evalDiv(leftObject, rightObject object.Object) object.Object {
	// TODO(ja): Handle other types
	left := leftObject.(*object.Integer)
	right := rightObject.(*object.Integer)

	return &object.Integer{
		Value: left.Value / right.Value,
	}
}

func evalEq(leftObject, rightObject object.Object) object.Object {
	result := false
	if leftObject.Type() == object.INTEGER_OBJ && rightObject.Type() == object.INTEGER_OBJ {
		left := leftObject.(*object.Integer)
		right := rightObject.(*object.Integer)
		result = left.Value == right.Value
	} else if leftObject.Type() == object.BOOLEAN_OBJ && rightObject.Type() == object.BOOLEAN_OBJ {
		left := leftObject.(*object.Boolean)
		right := rightObject.(*object.Boolean)
		result = left.Value == right.Value
	}

	return &object.Boolean{
		Value: result,
	}
}

func evalNeq(leftObject, rightObject object.Object) object.Object {
	result := false
	if leftObject.Type() == object.INTEGER_OBJ && rightObject.Type() == object.INTEGER_OBJ {
		left := leftObject.(*object.Integer)
		right := rightObject.(*object.Integer)
		result = left.Value != right.Value
	} else if leftObject.Type() == object.BOOLEAN_OBJ && rightObject.Type() == object.BOOLEAN_OBJ {
		left := leftObject.(*object.Boolean)
		right := rightObject.(*object.Boolean)
		result = left.Value != right.Value
	}

	return &object.Boolean{
		Value: result,
	}
}

func evalLess(leftObject, rightObject object.Object) object.Object {
	result := false
	if leftObject.Type() == object.INTEGER_OBJ && rightObject.Type() == object.INTEGER_OBJ {
		left := leftObject.(*object.Integer)
		right := rightObject.(*object.Integer)
		result = left.Value < right.Value
	}

	return &object.Boolean{
		Value: result,
	}
}

func evalGreater(leftObject, rightObject object.Object) object.Object {
	result := false
	if leftObject.Type() == object.INTEGER_OBJ && rightObject.Type() == object.INTEGER_OBJ {
		left := leftObject.(*object.Integer)
		right := rightObject.(*object.Integer)
		result = left.Value > right.Value
	}

	return &object.Boolean{
		Value: result,
	}
}

func evalInfixExpr(leftObject object.Object, operator token.Token, rightObject object.Object) object.Object {
	switch operator.Type {
	case token.PLUS:
		return evalAdd(leftObject, rightObject)
	case token.MINUS:
		return evalSub(leftObject, rightObject)
	case token.ASTERISK:
		return evalMul(leftObject, rightObject)
	case token.SLASH:
		return evalDiv(leftObject, rightObject)
	case token.EQ:
		return evalEq(leftObject, rightObject)
	case token.NOT_EQ:
		return evalNeq(leftObject, rightObject)
	case token.LT:
		return evalLess(leftObject, rightObject)
	case token.GT:
		return evalGreater(leftObject, rightObject)
	default:
		log.Fatalf("Unsupported infix operator: %v", operator)
	}

	return nil
}

func evalBang(obj object.Object) object.Object {
	// TODO(ja): Handle other types
	boolObj := obj.(*object.Boolean)

	return &object.Boolean{Value: !boolObj.Value}
}

func evalMinus(obj object.Object) object.Object {
	// TODO(ja): Handle other types
	intObj := obj.(*object.Integer)

	return &object.Integer{Value: -intObj.Value}
}

func evalPrefixExpr(obj object.Object, operator token.Token) object.Object {
	switch operator.Type {
	case token.BANG:
		return evalBang(obj)
	case token.MINUS:
		return evalMinus(obj)
	default:
		log.Fatalf("Unsupported prefix operator: %v", operator)
	}

	return nil
}

func evalIfExpr(expr *ast.IfExpr) object.Object {
	condition := Eval(expr.Condition)
	// TODO(ja): Handle error gracefully
	boolCondition := condition.(*object.Boolean)

	if boolCondition.Value {
		return Eval(expr.Consequence)
	}

	if expr.Alternative != nil {
		return Eval(expr.Alternative)
	}

	return &object.Null{}
}

func evalBlockStatement(stmt *ast.BlockStatement) object.Object {
	var result object.Object
	for _, stmt := range stmt.Statements {
		result = Eval(stmt)

		if returnObj, ok := result.(*object.Return); ok {
			return returnObj
		}
	}
	return result
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)

		if returnObj, ok := result.(*object.Return); ok {
			return returnObj.Value
		}
	}
	return result
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expr)

	case *ast.IntegerLiteralExpr:
		return &object.Integer{Value: node.Value}

	case *ast.BoolLiteralExpr:
		return &object.Boolean{Value: node.Value}

	case *ast.PrefixExpr:
		return evalPrefixExpr(Eval(node.InnerExpr), node.OperatorToken)

	case *ast.InfixExpr:
		return evalInfixExpr(Eval(node.LeftExpr), node.OperatorToken, Eval(node.RightExpr))

	case *ast.IfExpr:
		return evalIfExpr(node)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.ReturnStatement:
		var result object.Object = &object.Null{}
		if node.Expr != nil {
			result = Eval(node.Expr)
		}
		return &object.Return{Value: result}

	default:
		log.Fatalf("Unimplemented evaluation of node type: %T\n", node)
	}

	return nil
}
