package evaluator

import (
  "fmt"
  "Monkey/ast"
  "Monkey/object"
  )

var (
  // no need to consturct new variables since, boolean/null objects are the same.
  NULL  = &object.Null{}
  TRUE  = &object.Boolean{Value: true}
  FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
  switch nodeType := node.(type) {

  // Statements
  case *ast.Program:
    return evalProgram(nodeType, env)

  case *ast.LetStatement:
    val := Eval(nodeType.Value, env)
    if isError(val) {
      return val
    }
    env.Set(nodeType.Name.Value, val)

  case *ast.ExpressionStatement:
    return Eval(nodeType.Expression, env)

  case *ast.BlockStatement:
    return evalBlockStatement(nodeType, env)

  // Expressions
  case *ast.IntegerLiteral:
    return &object.Integer{Value: nodeType.Value}

  case *ast.Identifier:
    return evalIdentifier(nodeType, env)

  case *ast.Boolean:
    return nativeBoolToBooleanObject(nodeType.Value)

  case *ast.PrefixExpression:
    right := Eval(nodeType.Right, env)
    if isError(right) {
      return right
    }
    return evalPrefixExpression(nodeType.Operator, right)

  case *ast.InfixExpression:
    right := Eval(nodeType.Right, env)
    if isError(right){
      return right
    }
    left := Eval(nodeType.Left, env)
    if isError(left){
      return left
    }
    return evalInfixExpression(nodeType.Operator, left, right)

  case *ast.IfExpression:
    return evalIfExpression(nodeType, env)

  case *ast.ReturnStatement:
    val := Eval(nodeType.ReturnValue, env)
    if isError(val){
      return val
    }
    return &object.ReturnValue{Value: val}

  case *ast.FunctionLiteral:
    params := nodeType.Parameters
    body := nodeType.Body
    return &object.Function{Parameters: params, Env: env, Body: body}

  case *ast.CallExpression:
    function := Eval(nodeType.Function, env)
    if isError(function) {
      return function
    }
    args := evalExpressions(nodeType.Arguments, env)
    if len(args) == 1 && isError(args[0]) {
      return args[0]
    }
    return applyFunction(function, args)
  }
  return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
  var result object.Object
  for _, statement := range program.Statements {
    result = Eval(statement, env)

    switch result := result.(type) {
      case *object.ReturnValue:
        return result.Value
      case *object.Error:
        return result
      }
    }
  return result
}
// Usage: turn functionCall arguments into []object.Object.
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
  var result []object.Object

  for _, expression := range exps {
    evaluated := Eval(expression, env)
    if isError(evaluated) {
      return []object.Object{evaluated}
    }
    result = append(result, evaluated)
  }
  return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
  var result object.Object
  for _, statement := range block.Statements {
    result = Eval(statement, env)
    if result != nil {
      rt := result.Type()
      if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
        return result
      }
    }
  }
  return result
}
/***** Prefix Expressions *****/

func evalPrefixExpression(operator string, right object.Object) object.Object {
  switch operator {
    case "!":
      return evalBangOperatorExpression(right)
    case "-":
      return evalMinusPrefixOperatorExpression(right)
    default:
      return newError("unknown operator: %s%s", operator, right.Type())
  }
}

func evalBangOperatorExpression(right object.Object) object.Object {
  switch right := right.(type) {
    case *object.Boolean:
      if (right.Value){
        return FALSE
      }
      return TRUE

    case *object.Integer:
     if (right.Value == 0){
        return TRUE
      }
      return FALSE

    case *object.Null:
      return TRUE

    default:
      return FALSE
  }
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
  if right.Type() != object.INTEGER_OBJ {
    return newError("unknown operator: -%s", right.Type())
  }
  value := right.(*object.Integer).Value
  return &object.Integer{Value: -value}
}

/****** Infix Expressions ******/

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
  switch {
  case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
    return evalIntegerInfixExpression(operator, left, right)
  case operator == "==":
    return nativeBoolToBooleanObject(left == right)
  case operator == "!=":
    return nativeBoolToBooleanObject(left != right)
  case left.Type() != right.Type():
    return newError("type mismatch: %s %s %s",
    left.Type(), operator, right.Type())
  default:
    return newError("unknown operator: %s %s %s",
    left.Type(), operator, right.Type())  }
}


func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
  leftVal := left.(*object.Integer).Value
  rightVal := right.(*object.Integer).Value
  switch operator {
    case "+":
      return &object.Integer{Value: leftVal + rightVal}
    case "-":
      return &object.Integer{Value: leftVal - rightVal}
    case "*":
      return &object.Integer{Value: leftVal * rightVal}
    case "/":
      return &object.Integer{Value: leftVal / rightVal}
    case "<":
      return nativeBoolToBooleanObject(leftVal < rightVal)
    case ">":
      return nativeBoolToBooleanObject(leftVal > rightVal)
    case "==":
      return nativeBoolToBooleanObject(leftVal == rightVal)
    case "!=":
      return nativeBoolToBooleanObject(leftVal != rightVal)
    default:
     return newError("unknown operator: %s %s %s",
      left.Type(), operator, right.Type())
  }
}

/***** If - Else expressions ******/

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
  condition := Eval(ie.Condition, env)
  if isError(condition) {
    return condition
  }
  if isTruthy(condition) {
    return Eval(ie.Consequence, env)
  } else if ie.Alternative != nil {
    return Eval(ie.Alternative, env)
  } else {
    return NULL
  }
}

// If expression condition 
func isTruthy(obj object.Object) bool {
  switch obj {
  case NULL:
    return false
  case TRUE:
    return true
  case FALSE:
    return false
  default:
    if obj.Type() == object.INTEGER_OBJ{
      return obj.(*object.Integer).Value != 0
    }
  }
  return false
}

/***** Functions *****/

func applyFunction(fn object.Object, args []object.Object) object.Object {
  function, ok := fn.(*object.Function)
  if !ok {
    return newError("not a function: %s", fn.Type())
  }
  // open a new Environment and save the current one in Env.outer.
  extendedEnv := extendFunctionEnv(function, args)
  evaluated := Eval(function.Body, extendedEnv)
  return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
  env := object.NewEnclosedEnvironment(fn.Env)
  for paramIdx, param := range fn.Parameters {
    env.Set(param.Value, args[paramIdx])
  }
  return env
}

func unwrapReturnValue(obj object.Object) object.Object {
  if returnValue, ok := obj.(*object.ReturnValue); ok {
    return returnValue.Value
  }
  return obj
}

// -------


func nativeBoolToBooleanObject(input bool) *object.Boolean {
  if input {
    return TRUE
  }
  return FALSE
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
  val, ok := env.Get(node.Value)
  if !ok {
    return newError("identifier not found: " + node.Value)
  }
  return val
}

/****** Errors ******/

func newError(format string, a ...interface{}) *object.Error {
  return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
  if obj != nil {
    return obj.Type() == object.ERROR_OBJ
  }
  return false
}
