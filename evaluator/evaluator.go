package evaluator

import (
	"fmt"
	"log"

	"github.com/javier-varez/monkey_interpreter/ast"
	"github.com/javier-varez/monkey_interpreter/object"
	"github.com/javier-varez/monkey_interpreter/token"
)

func mkError(s token.Span, msg string) *object.Error {
	return &object.Error{Span: s, Message: msg}
}

func evalAdd(leftObject, rightObject object.Object) object.Object {
	left := leftObject.(*object.Integer)
	right := rightObject.(*object.Integer)

	return &object.Integer{
		Value: left.Value + right.Value,
	}
}

func evalSub(leftObject, rightObject object.Object) object.Object {
	left := leftObject.(*object.Integer)
	right := rightObject.(*object.Integer)

	return &object.Integer{
		Value: left.Value - right.Value,
	}
}

func evalMul(leftObject, rightObject object.Object) object.Object {
	left := leftObject.(*object.Integer)
	right := rightObject.(*object.Integer)

	return &object.Integer{
		Value: left.Value * right.Value,
	}
}

func evalDiv(leftObject, rightObject object.Object) object.Object {
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
	} else {
		panic("Unsupported operands.")
	}

	return &object.Boolean{Value: result}
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
	} else {
		panic("Unsupported operands.")
	}

	return &object.Boolean{Value: result}
}

func evalLess(leftObject, rightObject object.Object) object.Object {
	left := leftObject.(*object.Integer)
	right := rightObject.(*object.Integer)
	result := left.Value < right.Value

	return &object.Boolean{Value: result}
}

func evalGreater(leftObject, rightObject object.Object) object.Object {
	left := leftObject.(*object.Integer)
	right := rightObject.(*object.Integer)
	result := left.Value > right.Value

	return &object.Boolean{Value: result}
}

func evalInfixExpr(expr *ast.InfixExpr) object.Object {
	left := Eval(expr.LeftExpr)
	if left.Type() == object.ERROR_VALUE_OBJ {
		return left
	}

	right := Eval(expr.RightExpr)
	if right.Type() == object.ERROR_VALUE_OBJ {
		return right
	}

	switch expr.OperatorToken.Type {
	case token.PLUS:
		fallthrough
	case token.MINUS:
		fallthrough
	case token.ASTERISK:
		fallthrough
	case token.SLASH:
		fallthrough
	case token.LT:
		fallthrough
	case token.GT:
		if left.Type() != object.INTEGER_OBJ {
			return mkError(expr.LeftExpr.Span(), "Expression does not evaluate to an integer object")
		}

		if right.Type() != object.INTEGER_OBJ {
			return mkError(expr.RightExpr.Span(), "Expression does not evaluate to an integer object")
		}
	case token.EQ:
		fallthrough
	case token.NOT_EQ:
		if left.Type() != object.INTEGER_OBJ && left.Type() != object.BOOLEAN_OBJ {
			return mkError(expr.LeftExpr.Span(), "Expression does not evaluate to an integer or boolean object")
		}

		if right.Type() != object.INTEGER_OBJ && right.Type() != object.BOOLEAN_OBJ {
			return mkError(expr.RightExpr.Span(), "Expression does not evaluate to an integer or boolean object")
		}

		if right.Type() != left.Type() {
			return mkError(expr.Span(), "Left and right arguments to the infix operator do not have the same type")
		}
	default:
		log.Fatalf("Unsupported infix operator: %v", expr.OperatorToken)
	}

	switch expr.OperatorToken.Type {
	case token.PLUS:
		return evalAdd(left, right)
	case token.MINUS:
		return evalSub(left, right)
	case token.ASTERISK:
		return evalMul(left, right)
	case token.SLASH:
		return evalDiv(left, right)
	case token.EQ:
		return evalEq(left, right)
	case token.NOT_EQ:
		return evalNeq(left, right)
	case token.LT:
		return evalLess(left, right)
	case token.GT:
		return evalGreater(left, right)
	default:
		log.Fatalf("Unsupported infix operator: %v", expr.OperatorToken)
	}

	return nil
}

func evalBang(obj object.Object) object.Object {
	boolObj := obj.(*object.Boolean)
	return &object.Boolean{Value: !boolObj.Value}
}

func evalMinus(obj object.Object) object.Object {
	intObj := obj.(*object.Integer)
	return &object.Integer{Value: -intObj.Value}
}

func evalPrefixExpr(expr *ast.PrefixExpr) object.Object {
	innerResult := Eval(expr.InnerExpr)
	if innerResult.Type() == object.ERROR_VALUE_OBJ {
		return innerResult
	}

	switch expr.OperatorToken.Type {
	case token.BANG:
		if innerResult.Type() != object.BOOLEAN_OBJ {
			return mkError(expr.Span(), fmt.Sprintf("%q requires a boolean argument", token.BANG))
		}
		return evalBang(innerResult)
	case token.MINUS:
		if innerResult.Type() != object.INTEGER_OBJ {
			return mkError(expr.Span(), fmt.Sprintf("%q requires an integer argument", token.MINUS))
		}
		return evalMinus(innerResult)
	default:
		log.Fatalf("Unsupported prefix operator: %v", expr.OperatorToken.Type)
	}
	return nil
}

func evalIfExpr(expr *ast.IfExpr) object.Object {
	condition := Eval(expr.Condition)
	if condition.Type() == object.ERROR_VALUE_OBJ {
		return condition
	}

	if condition.Type() != object.BOOLEAN_OBJ {
		return mkError(expr.Condition.Span(), "Condition must evaluate to a boolean object")
	}

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

		if result.Type() == object.ERROR_VALUE_OBJ {
			return result
		}

		if result.Type() == object.RETURN_VALUE_OBJ {
			returnObj := result.(*object.Return)
			return returnObj
		}
	}
	return result
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)

		if result.Type() == object.ERROR_VALUE_OBJ {
			return result
		}

		if result.Type() == object.RETURN_VALUE_OBJ {
			returnObj := result.(*object.Return)
			return returnObj.Value
		}
	}
	return result
}

func evalReturnStatement(stmt *ast.ReturnStatement) object.Object {
	var result object.Object = &object.Null{}
	if stmt.Expr != nil {
		result = Eval(stmt.Expr)
		if result.Type() == object.ERROR_VALUE_OBJ {
			return result
		}
	}

	return &object.Return{Value: result}
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
		return evalPrefixExpr(node)

	case *ast.InfixExpr:
		return evalInfixExpr(node)

	case *ast.IfExpr:
		return evalIfExpr(node)

	case *ast.BlockStatement:
		return evalBlockStatement(node)

	case *ast.ReturnStatement:
		return evalReturnStatement(node)

	default:
		log.Fatalf("Unimplemented evaluation of node type: %T\n", node)
	}

	return nil
}
