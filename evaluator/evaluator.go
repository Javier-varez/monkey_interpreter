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
			return mkError(span, "\"len\" builtin takes a single string or array argument")
		}

		strObj, ok := objects[0].(*object.String)
		if ok {
			return &object.Integer{Value: int64(len(strObj.Value))}
		}

		arrObj, ok := objects[0].(*object.Array)
		if ok {
			return &object.Integer{Value: int64(len(arrObj.Elems))}
		}

		return mkError(span, "\"len\" builtin takes a single string or array argument")
	},
	"first": func(span token.Span, objects ...object.Object) object.Object {
		if len(objects) != 1 {
			return mkError(span, "\"first\" builtin takes a single array argument")
		}

		arrObj, ok := objects[0].(*object.Array)
		if ok {
			if len(arrObj.Elems) == 0 {
				return mkError(span, "Array is empty")
			}
			return arrObj.Elems[0]
		}

		return mkError(span, "\"first\" builtin takes a single array argument")
	},
	"last": func(span token.Span, objects ...object.Object) object.Object {
		if len(objects) != 1 {
			return mkError(span, "\"last\" builtin takes a single array argument")
		}

		arrObj, ok := objects[0].(*object.Array)
		if ok {
			if len(arrObj.Elems) == 0 {
				return mkError(span, "Array is empty")
			}
			return arrObj.Elems[len(arrObj.Elems)-1]
		}

		return mkError(span, "\"last\" builtin takes a single array argument")
	},
	"rest": func(span token.Span, objects ...object.Object) object.Object {
		if len(objects) != 1 {
			return mkError(span, "\"rest\" builtin takes a single array argument")
		}

		arrObj, ok := objects[0].(*object.Array)
		if ok {
			if len(arrObj.Elems) == 0 {
				return mkError(span, "Array is empty")
			}
			newArr := &object.Array{Elems: make([]object.Object, len(arrObj.Elems)-1)}
			copy(newArr.Elems[:], arrObj.Elems[1:])
			return newArr
		}

		return mkError(span, "\"rest\" builtin takes a single array argument")
	},
	"push": func(span token.Span, objects ...object.Object) object.Object {
		if len(objects) != 2 {
			return mkError(span, "\"push\" builtin takes an array argument and a new object to push")
		}

		arrObj, ok := objects[0].(*object.Array)
		if !ok {
			return mkError(span, "\"push\" builtin takes an array argument and a new object to push")
		}

		oldLen := len(arrObj.Elems)
		newArr := &object.Array{Elems: make([]object.Object, oldLen+1)}
		copy(newArr.Elems[:oldLen], arrObj.Elems[:])
		newArr.Elems[oldLen] = objects[1]

		return newArr
	},
	"puts": func(span token.Span, objects ...object.Object) object.Object {
		for _, object := range objects {
			fmt.Print(object.Inspect())
		}
		fmt.Println()
		return &object.Null{}
	},
	"toArray": func(span token.Span, objects ...object.Object) object.Object {
		if len(objects) != 1 {
			return mkError(span, "\"toArray\" builtin takes a VarArg argument")
		}

		varArgObj, ok := objects[0].(*object.VarArgs)
		if !ok {
			return mkError(span, "\"toArray\" builtin takes a VarArg argument")
		}

		return &object.Array{Elems: varArgObj.Elems}
	},
	"contains": func(span token.Span, objects ...object.Object) object.Object {
		if len(objects) != 2 {
			return mkError(span, "\"contains\" builtin takes a HashMap argument and a key")
		}

		hashMapObj, ok := objects[0].(*object.HashMap)
		if !ok {
			return mkError(span, "First argument is not a hash map")
		}

		keyObj, ok := objects[1].(object.Hashable)
		if !ok {
			return mkError(span, "Second argument is not a hashable object")
		}

		elem, ok := hashMapObj.Elems[keyObj.HashKey()]
		if ok {
			if elem.Key.Inspect() == keyObj.Inspect() {
				return &object.Boolean{Value: true}
			}
		}

		return &object.Boolean{Value: false}
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
		Args:    expr.Args,
		VarArgs: expr.VarArgs,
		Body:    expr.Body,
		Env:     env.Copy(),
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
	// Eval args
	var args []object.Object
	for _, arg := range expr.Args {
		res := Eval(arg, env)
		if res.Type() == object.ERROR_VALUE_OBJ {
			return res
		}

		if res.Type() == object.VAR_ARGS_OBJ {
			// Expanded on call site
			varArgsObj := res.(*object.VarArgs)
			args = append(args, varArgsObj.Elems...)
		} else {
			args = append(args, res)
		}
	}

	if !fnObj.VarArgs && len(fnObj.Args) != len(args) {
		return mkError(expr.Span(), fmt.Sprintf("Callable takes %d arguments, but %d were supplied", len(fnObj.Args), len(args)))
	}

	if fnObj.VarArgs && len(fnObj.Args) > len(args) {
		return mkError(expr.Span(), fmt.Sprintf("Callable takes at least %d arguments, but only %d were supplied", len(fnObj.Args), len(args)))
	}

	// Bind args to new environment
	newEnv := object.NewEnclosedEnvironment(fnObj.Env)
	for i := range fnObj.Args {
		newEnv.Set(fnObj.Args[i].IdentToken.Literal, args[i])
	}
	if fnObj.VarArgs {
		varArgs := args[len(fnObj.Args):]
		newEnv.SetVarArgs(varArgs)
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

func evalIndexOperatorExpr(expr *ast.IndexOperatorExpr, env *object.Environment) object.Object {
	indexedObj := Eval(expr.ObjExpr, env)
	if indexedObj.Type() == object.ERROR_VALUE_OBJ {
		return indexedObj
	}

	indexObj := Eval(expr.IndexExpr, env)
	if indexObj.Type() == object.ERROR_VALUE_OBJ {
		return indexObj
	}

	if indexedObj.Type() == object.ARRAY_OBJ {
		if indexObj.Type() != object.INTEGER_OBJ {
			return mkError(expr.IndexExpr.Span(), "Expression must evaluate to an integer object")
		}

		indexValue := indexObj.(*object.Integer).Value
		elems := &indexedObj.(*object.Array).Elems

		if indexValue >= int64(len(*elems)) {
			return mkError(expr.IndexExpr.Span(), fmt.Sprintf("Index %d exceeds length of the array (%d)", indexValue, len(*elems)))
		}

		return (*elems)[indexValue]
	} else if indexedObj.Type() == object.MAP_OBJ {
		hashable, ok := indexObj.(object.Hashable)
		if !ok {
			return mkError(expr.IndexExpr.Span(), "Expression must evaluate to a hashable object")
		}

		hashKey := hashable.HashKey()
		mapObj := indexedObj.(*object.HashMap)
		value, ok := mapObj.Elems[hashKey]
		if !ok {
			return mkError(expr.Span(), fmt.Sprintf("Key %q not found", indexObj.Inspect()))
		}

		if value.Key.Inspect() != indexObj.Inspect() {
			return mkError(expr.Span(), fmt.Sprintf("Key %q not found", indexObj.Inspect()))
		}
		return value.Value
	}

	return mkError(expr.ObjExpr.Span(), "Expression must evaluate to an array or map object")

}

func evalVarArgsLiteralExpr(node *ast.VarArgsLiteralExpr, env *object.Environment) object.Object {
	if varArgs, ok := env.GetVarArgs(); ok {
		return &object.VarArgs{Elems: varArgs}
	}

	return mkError(node.Span(), "Function has no var args to expand")
}

func evalRangeExpr(node *ast.RangeExpr, env *object.Environment) object.Object {
	startObj := Eval(node.StartExpr, env)
	if startObj.Type() != object.INTEGER_OBJ {
		return mkError(node.StartExpr.Span(), "Expression does not evaluate to an integer object")
	}

	endObj := Eval(node.EndExpr, env)
	if endObj.Type() != object.INTEGER_OBJ {
		return mkError(node.StartExpr.Span(), "Expression does not evaluate to an integer object")
	}

	startIntObj := startObj.(*object.Integer)
	endIntObj := endObj.(*object.Integer)

	incr := int64(1)
	if startIntObj.Value > endIntObj.Value {
		// Decreasing range
		incr = -1
	}

	arrayObj := &object.Array{Elems: []object.Object{}}

	curValue := startIntObj.Value
	for curValue != endIntObj.Value {
		arrayObj.Elems = append(arrayObj.Elems, &object.Integer{Value: curValue})
		curValue = curValue + incr
	}

	return arrayObj
}

func evalMapLiteralExpr(node *ast.MapLiteralExpr, env *object.Environment) object.Object {
	mapObj := &object.HashMap{
		Elems: map[object.HashKey]object.HashEntry{},
	}

	for kExpr, vExpr := range node.Map {
		kObj := Eval(kExpr, env)
		if kObj.Type() == object.ERROR_VALUE_OBJ {
			return kObj
		}

		hashableK, ok := kObj.(object.Hashable)
		if !ok {
			return mkError(kExpr.Span(), "Expression is not hashable")
		}
		hashKey := hashableK.HashKey()

		vObj := Eval(vExpr, env)
		if vObj.Type() == object.ERROR_VALUE_OBJ {
			return vObj
		}

		hashEntry := object.HashEntry{
			Key:   kObj,
			Value: vObj,
		}

		mapObj.Elems[hashKey] = hashEntry
	}

	return mapObj
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

	case *ast.IndexOperatorExpr:
		return evalIndexOperatorExpr(node, env)

	case *ast.VarArgsLiteralExpr:
		return evalVarArgsLiteralExpr(node, env)

	case *ast.RangeExpr:
		return evalRangeExpr(node, env)

	case *ast.MapLiteralExpr:
		return evalMapLiteralExpr(node, env)

	default:
		log.Fatalf("Unimplemented evaluation of node type: %T\n", node)
	}

	return nil
}
