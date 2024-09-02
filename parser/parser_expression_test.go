package parser

import (
  "testing"
  "Monkey/ast"
  "Monkey/lexer"
  "fmt"
)

/***** Identifier expression tests *****/

func TestIdentifierExpression(t *testing.T) {
  input := "foobar"
  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()
  checkParserErrors(t, p)
  if len(program.Statements) != 1 {
    t.Fatalf("program has not enough statements. got=%d",
    len(program.Statements))
  }
  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
    program.Statements[0])
  }
  if !testIdentifier(t, stmt.Expression, input) {
    return
  }
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
  ident, ok := exp.(*ast.Identifier)
  if !ok {
    t.Errorf("exp not *ast.Identifier. got=%T", exp)
    return false
  }
  if ident.Value != value {
    t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
    return false
  }
  if ident.TokenLiteral() != value {
    t.Errorf("ident.TokenLiteral not %s. got=%s", value,
    ident.TokenLiteral())
    return false
  }
  return true
}

/***** Integer literal test *****/

func TestIntegerLiteralExpression(t *testing.T) {
  input := "5;"
  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()
  checkParserErrors(t, p)
  if len(program.Statements) != 1 {
    t.Fatalf("program has not enough statements. got=%d",
    len(program.Statements))
  }

  stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
  if !ok {
    t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
    program.Statements[0])
  }
  if !testIntegerLiteral(t, stmt.Expression, 5) {
    return
  }
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
  integ, ok := il.(*ast.IntegerLiteral)
  if !ok {
    t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
    return false
  }
  if integ.Value != value {
    t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
    return false
  }
  if integ.TokenLiteral() != fmt.Sprintf("%d", value) {

    t.Errorf("integ.TokenLiteral not %d. got=%s", value,
    integ.TokenLiteral())
    return false
  }
  return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
  switch v := expected.(type) {
  case int:
    return testIntegerLiteral(t, exp, int64(v))
  case int64:
    return testIntegerLiteral(t, exp, v)
  case string:
    return testIdentifier(t, exp, v)
  case bool:
    return testBooleanLiteral(t, exp, v)
  }
  t.Errorf("type of exp not handled. got=%T", exp)
  return false
}

/***** Prefix expressions e.g - !5, -1 *****/

func TestParsingPrefixExpressions(t *testing.T) {
  prefixTests := []struct {
    input        string
    operator     string
    integerValue interface{}
  }{
    {"!5;", "!", 5},
    {"-15;", "-", 15},
    {"!true", "!", true},
    {"!false", "!", false},
  }

  for _, tt := range prefixTests {
    l := lexer.New(tt.input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
      t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
      1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
      t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
      program.Statements[0])
    }

    exp, ok := stmt.Expression.(*ast.PrefixExpression)
    if !ok {
      t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
    }
    if exp.Operator != tt.operator {
      t.Fatalf("exp.Operator is not '%s'. got=%s",
      tt.operator, exp.Operator)
    }

    if !testLiteralExpression(t, exp.Right, tt.integerValue) {
      return
    }
  }
}

/***** Infix expressions e.g - 5 - 5, 3 * 2 *****/

func TestParsingInfixExpressions(t *testing.T) {
  infixTests := []struct {
    input string
    leftValue interface{}
    operator string
    rightValue interface{}
  }{
    {"5 + 5;", 5, "+", 5},
    {"5 - 5;", 5, "-", 5},
    {"5 * 5;", 5, "*", 5},
    {"5 / 5;", 5, "/", 5},
    {"5 > 5;", 5, ">", 5},
    {"5 < 5;", 5, "<", 5},
    {"5 == 5;", 5, "==", 5},
    {"5 != 5;", 5, "!=", 5},
    {"true == true", true, "==", true},
    {"true != false", true, "!=", false},
    {"false == false", false, "==", false},
  }
  for _, tt := range infixTests {
    l := lexer.New(tt.input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)
    if len(program.Statements) != 1 {
      t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
      1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
      t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
      program.Statements[0])
    }

    exp, ok := stmt.Expression.(*ast.InfixExpression)
    if !ok {
      t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
    }
    if !testInfixExpression(t, exp, tt.leftValue,
      tt.operator, tt.rightValue) {
      return
    }
  }
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
  operator string, right interface{}) bool {

  opExp, ok := exp.(*ast.InfixExpression)
  if !ok {
    t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
    return false
  }
  if !testLiteralExpression(t, opExp.Left, left) {
    return false
  }
  if opExp.Operator != operator {
    t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
    return false
  }
  if !testLiteralExpression(t, opExp.Right, right) {
    return false
  }
  return true
}


/***** boolean expressions *****/

func TestBooleanExpression(t *testing.T) {
  tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}
  for _, testCase := range tests {
    l := lexer.New(testCase.input)
    p := New(l)

    program := p.ParseProgram()
    checkParserErrors(t, p)
    if len(program.Statements) != 1 {
      t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
        1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
      t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
      program.Statements[0])
    }
    if !testBooleanLiteral(t, stmt.Expression, testCase.expectedBoolean) {
      return
    }
  }
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
  bool, ok := exp.(*ast.Boolean)
  if !ok {
    t.Errorf("exp not *ast.Boolean. got=%T", exp)
    return false
  }
  if bool.Value != value {
    t.Errorf("bool.Value not %t. got=%t", value, bool.Value)
    return false
  }
  if bool.TokenLiteral() != fmt.Sprintf("%t", value) {
    t.Errorf("bool.TokenLiteral not %t. got=%s",
    value, bool.TokenLiteral())
    return false
  }
  return true
}
