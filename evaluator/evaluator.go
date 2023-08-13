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

var builtins map[string]object.BuiltinFunction = map[string]object.BuiltinFunction{
	"len": func(span token.Span, objects ...object.Object) object.Object {
		if len(objects) != 1 {
			return mkError(span, "\"len\" builtin takes a single string argument")
		}

		strObj, ok := objects[0].(*object.String)
		if ok {
			return &object.Integer{Value: int64(len(strObj.Value))}
		}

		return mkError(span, "\"len\" builtin takes a single string argument")
	},
	"puts": func(span token.Span, objects ...object.Object) object.Object {
		for _, object := range objects {
			fmt.Print(object.Inspect())
		}
		fmt.Println()
		return &object.Null{}
	},
}

func evalAdd(leftObject, rightObject object.Object) object.Object {
	if leftObject.Type() == object.INTEGER_OBJ && rightObject.Type() == object.INTEGER_OBJ {
		left := leftObject.(*object.Integer)
		right := rightObject.(*object.Integer)
		return &object.Integer{Value: left.Value + right.Value}
	} else if leftObject.Type() == object.STRING_OBJ && rightObject.Type() == object.STRING_OBJ {
		left := leftObject.(*object.String)
		right := rightObject.(*object.String)
		return &object.String{Value: left.Value + right.Value}
	}

	panic("Invalid types for call to evalAdd")
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
	} else if leftObject.Type() == object.STRING_OBJ && rightObject.Type() == object.STRING_OBJ {
		left := leftObject.(*object.String)
		right := rightObject.(*object.String)
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
	} else if leftObject.Type() == object.STRING_OBJ && rightObject.Type() == object.STRING_OBJ {
		left := leftObject.(*object.String)
		right := rightObject.(*object.String)
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

func evalInfixExpr(expr *ast.InfixExpr, env *object.Environment) object.Object {
	left := Eval(expr.LeftExpr, env)
	if left.Type() == object.ERROR_VALUE_OBJ {
		return left
	}

	right := Eval(expr.RightExpr, env)
	if right.Type() == object.ERROR_VALUE_OBJ {
		return right
	}

	switch expr.OperatorToken.Type {
	case token.PLUS:
		if left.Type() != object.INTEGER_OBJ && left.Type() != object.STRING_OBJ {
			return mkError(expr.LeftExpr.Span(), "Expression does not evaluate to an integer or string object")
		}

		if right.Type() != object.INTEGER_OBJ && right.Type() != object.STRING_OBJ {
			return mkError(expr.RightExpr.Span(), "Expression does not evaluate to an integer or string object")
		}

		if right.Type() != left.Type() {
			return mkError(expr.Span(), "Left and right arguments to the infix operator do not have the same type")
		}
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
		if left.Type() != object.INTEGER_OBJ && left.Type() != object.BOOLEAN_OBJ && left.Type() != object.STRING_OBJ {
			return mkError(expr.LeftExpr.Span(), "Expression does not evaluate to an integer, boolean or string object")
		}

		if right.Type() != object.INTEGER_OBJ && right.Type() != object.BOOLEAN_OBJ && right.Type() != object.STRING_OBJ {
			return mkError(expr.RightExpr.Span(), "Expression does not evaluate to an integer, boolean or string object")
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

func evalPrefixExpr(expr *ast.PrefixExpr, env *object.Environment) object.Object {
	innerResult := Eval(expr.InnerExpr, env)
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

func evalIfExpr(expr *ast.IfExpr, env *object.Environment) object.Object {
	condition := Eval(expr.Condition, env)
	if condition.Type() == object.ERROR_VALUE_OBJ {
		return condition
	}

	if condition.Type() != object.BOOLEAN_OBJ {
		return mkError(expr.Condition.Span(), "Condition must evaluate to a boolean object")
	}

	boolCondition := condition.(*object.Boolean)

	if boolCondition.Value {
		return Eval(expr.Consequence, env)
	}

	if expr.Alternative != nil {
		return Eval(expr.Alternative, env)
	}

	return &object.Null{}
}

func evalBlockStatement(stmt *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range stmt.Statements {
		result = Eval(stmt, env)

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

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)

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

func evalReturnStatement(stmt *ast.ReturnStatement, env *object.Environment) object.Object {
	var result object.Object = &object.Null{}
	if stmt.Expr != nil {
		result = Eval(stmt.Expr, env)
		if result.Type() == object.ERROR_VALUE_OBJ {
			return result
		}
	}

	return &object.Return{Value: result}
}

func evalLetStatement(stmt *ast.LetStatement, env *object.Environment) object.Object {
	obj := Eval(stmt.Expr, env)
	if obj.Type() == object.ERROR_VALUE_OBJ {
		return obj
	}

	return env.Set(stmt.IdentExpr.(*ast.IdentifierExpr).IdentToken.Literal, obj)
}

func evalIdentifierExpr(expr *ast.IdentifierExpr, env *object.Environment) object.Object {
	if obj, ok := env.Get(expr.IdentToken.Literal); ok {
		return obj
	}

	// Try builtin identifiers
	if builtin, ok := builtins[expr.IdentToken.Literal]; ok {
		return &object.Builtin{Function: builtin}
	}

	return mkError(expr.Span(), "Identifier not found")
}

func evalFnLiteralExpr(expr *ast.FnLiteralExpr, env *object.Environment) object.Object {
	return &object.Function{
		Args: expr.Args,
		Body: expr.Body,
		Env:  env.Copy(),
	}
}

func evalCallBuiltin(builtin *object.Builtin, expr *ast.CallExpr, env *object.Environment) object.Object {
	// Eval args
	var args []object.Object
	for _, arg := range expr.Args {
		res := Eval(arg, env)
		if res.Type() == object.ERROR_VALUE_OBJ {
			return res
		}
		args = append(args, res)
	}

	return builtin.Function(expr.Span(), args...)
}

func evalCallFnObject(fnObj *object.Function, expr *ast.CallExpr, env *object.Environment) object.Object {
	if len(fnObj.Args) != len(expr.Args) {
		return mkError(expr.Span(), fmt.Sprintf("Callable takes %d arguments, but %d were supplied", len(fnObj.Args), len(expr.Args)))
	}

	// Eval args
	var args []object.Object
	for _, arg := range expr.Args {
		res := Eval(arg, env)
		if res.Type() == object.ERROR_VALUE_OBJ {
			return res
		}
		args = append(args, res)
	}

	// Bound args to new environment
	newEnv := object.NewEnclosedEnvironment(fnObj.Env)
	for i := range expr.Args {
		newEnv.Set(fnObj.Args[i].IdentToken.Literal, args[i])
	}

	result := Eval(fnObj.Body, newEnv)

	// Unwrap return so that it does not cross the boundary of the function
	if result.Type() == object.RETURN_VALUE_OBJ {
		returnObject := result.(*object.Return)
		result = returnObject.Value
	}
	return result
}

func evalCallExpr(expr *ast.CallExpr, env *object.Environment) object.Object {
	fn := Eval(expr.CallableExpr, env)
	if fn.Type() == object.ERROR_VALUE_OBJ {
		return fn
	}

	if fn.Type() == object.BUILTIN_OBJ {
		builtinObj := fn.(*object.Builtin)
		return evalCallBuiltin(builtinObj, expr, env)
	}

	if fn.Type() == object.FUNCTION_OBJ {
		fnObj := fn.(*object.Function)
		return evalCallFnObject(fnObj, expr, env)
	}

	return mkError(expr.CallableExpr.Span(), "Call expression must have a callable type (function literal or identifier bounded to a function)")
}

func evalArrayLiteralExpr(expr *ast.ArrayLiteralExpr, env *object.Environment) object.Object {
	result := &object.Array{
		Elems: []object.Object{},
	}

	for _, inner := range expr.Elems {
		innerEval := Eval(inner, env)
		if innerEval.Type() == object.ERROR_VALUE_OBJ {
			return innerEval
		}

		result.Elems = append(result.Elems, innerEval)
	}

	return result
}

func evalArrayIndexOperatorExpr(expr *ast.ArrayIndexOperatorExpr, env *object.Environment) object.Object {
	arrayObj := Eval(expr.ArrayExpr, env)
	if arrayObj.Type() == object.ERROR_VALUE_OBJ {
		return arrayObj
	}

	if arrayObj.Type() != object.ARRAY_OBJ {
		return mkError(expr.ArrayExpr.Span(), "Expression must evaluate to an array object")
	}

	indexObj := Eval(expr.IndexExpr, env)
	if indexObj.Type() == object.ERROR_VALUE_OBJ {
		return indexObj
	}

	if indexObj.Type() != object.INTEGER_OBJ {
		return mkError(expr.IndexExpr.Span(), "Expression must evaluate to an integer object")
	}

	indexValue := indexObj.(*object.Integer).Value
	elems := &arrayObj.(*object.Array).Elems

	if indexValue >= int64(len(*elems)) {
		return mkError(expr.IndexExpr.Span(), fmt.Sprintf("Index %d exceeds length of the array (%d)", indexValue, len(*elems)))
	}

	return (*elems)[indexValue]
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expr, env)

	case *ast.IntegerLiteralExpr:
		return &object.Integer{Value: node.Value}

	case *ast.BoolLiteralExpr:
		return &object.Boolean{Value: node.Value}

	case *ast.PrefixExpr:
		return evalPrefixExpr(node, env)

	case *ast.InfixExpr:
		return evalInfixExpr(node, env)

	case *ast.IfExpr:
		return evalIfExpr(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)

	case *ast.LetStatement:
		return evalLetStatement(node, env)

	case *ast.IdentifierExpr:
		return evalIdentifierExpr(node, env)

	case *ast.FnLiteralExpr:
		return evalFnLiteralExpr(node, env)

	case *ast.CallExpr:
		return evalCallExpr(node, env)

	case *ast.StringLiteralExpr:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteralExpr:
		return evalArrayLiteralExpr(node, env)

	case *ast.ArrayIndexOperatorExpr:
		return evalArrayIndexOperatorExpr(node, env)

	default:
		log.Fatalf("Unimplemented evaluation of node type: %T\n", node)
	}

	return nil
}
