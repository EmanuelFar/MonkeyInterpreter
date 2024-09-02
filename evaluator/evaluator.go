package evaluator

import (
  "Monkey/ast"
  "Monkey/object"
  )

var (
  // no need to consturct new variables since, boolean/null objects are the same.
  NULL  = &object.Null{}
  TRUE  = &object.Boolean{Value: true}
  FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
  switch nodeType := node.(type) {

  // Statements
  case *ast.Program:
    return evalProgram(nodeType.Statements)

  case *ast.ExpressionStatement:
    return Eval(nodeType.Expression)

  case *ast.BlockStatement:
    return evalBlockStatement(nodeType)

  // Expressions
  case *ast.IntegerLiteral:
    return &object.Integer{Value: nodeType.Value}

  case *ast.Boolean:
    return nativeBoolToBooleanObject(nodeType.Value)

  case *ast.PrefixExpression:
    right := Eval(nodeType.Right)
    return evalPrefixExpression(nodeType.Operator, right)

  case *ast.InfixExpression:
    right := Eval(nodeType.Right)
    left := Eval(nodeType.Left)
    return evalInfixExpression(nodeType.Operator, left, right)

  case *ast.IfExpression:
    return evalIfExpression(nodeType)

  case *ast.ReturnStatement:
    val := Eval(nodeType.ReturnValue)
    return &object.ReturnValue{Value: val}
  }
  return nil
}

func evalProgram(stmts []ast.Statement) object.Object {
  var result object.Object

  for _,statement := range stmts {
    result = Eval(statement)

    if returnValue, ok := result.(*object.ReturnValue); ok {
      return returnValue.Value
    }
  }
  return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
  var result object.Object
  for _, statement := range block.Statements {
    result = Eval(statement)
    if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
      return result
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
      return NULL
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
    return NULL
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
  default:
    return NULL
  }
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
      return NULL
  }
}

/***** If - Else expressions ******/

func evalIfExpression(ie *ast.IfExpression) object.Object {
  condition := Eval(ie.Condition)
  if isTruthy(condition) {
    return Eval(ie.Consequence)
  } else if ie.Alternative != nil {
    return Eval(ie.Alternative)
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

func nativeBoolToBooleanObject(input bool) *object.Boolean {
  if input {
    return TRUE
  }
  return FALSE
}


