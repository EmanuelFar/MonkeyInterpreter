package parser

import (
  "testing"
  "Monkey/ast"
  "Monkey/lexer"
  "fmt"
)

/***** Identifier expression tests *****/

func TestIdentifierExpression(t *testing.T) {
  input := "foobar;"
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

  ident, ok := stmt.Expression.(*ast.Identifier)
  if !ok {
    t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
  }
  if ident.Value != "foobar" {
    t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
  }
  if ident.TokenLiteral() != "foobar" {
    t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
    ident.TokenLiteral())
  }
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

  literal, ok := stmt.Expression.(*ast.IntegerLiteral)
  if !ok {
    t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
  }
  if literal.Value != 5 {
    t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
  }
  if literal.TokenLiteral() != "5" {
    t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
    literal.TokenLiteral())
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

/***** Prefix expressions e.g - !5, -1 *****/

func TestParsingPrefixExpressions(t *testing.T) {
  prefixTests := []struct {
    input        string
    operator     string
    integerValue int64
  }{
    {"!5;", "!", 5},
    {"-15;", "-", 15},
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
    if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
      return
    }
  }
}

/***** Infix expressions e.g - 5 - 5, 3 * 2 *****/

func TestParsingInfixExpressions(t *testing.T) {
  infixTests := []struct {
    input string
    leftValue int64
    operator string
    rightValue int64
  }{
    {"5 + 5;", 5, "+", 5},
    {"5 - 5;", 5, "-", 5},
    {"5 * 5;", 5, "*", 5},
    {"5 / 5;", 5, "/", 5},
    {"5 > 5;", 5, ">", 5},
    {"5 < 5;", 5, "<", 5},
    {"5 == 5;", 5, "==", 5},
    {"5 != 5;", 5, "!=", 5},
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
    if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
      return
    }
    if exp.Operator != tt.operator {
      t.Fatalf("exp.Operator is not '%s'. got=%s",
      tt.operator, exp.Operator)
    }
    if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
      return
    }
  }
}
