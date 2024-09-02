package parser

import (
  "testing"
  "Monkey/ast"
  "Monkey/lexer"
)

/*****  Let Statement tests *****/

func TestLetStatements(t *testing.T) {
  tests := []struct {
    input   string
    expectedIdentifier string
    expectedValue interface{}
  }{
    {"let x = 5;", "x", 5},
    {"let y = true;", "y", true},
    {"let foobar = y;", "foobar", "y"},
  }
  for _, tt := range tests {
    l := lexer.New(tt.input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)
    if len(program.Statements) != 1 {
      t.Fatalf("program.Statements does not contain 1 statements. got=%d",
      len(program.Statements))
    }
    stmt := program.Statements[0]
    if !testLetStatement(t, stmt, tt.expectedIdentifier) {
      return
    }
    val := stmt.(*ast.LetStatement).Value
    if !testLiteralExpression(t, val, tt.expectedValue) {
      return
    }
  }
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
  if s.TokenLiteral() != "let" {
    t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
    return false
  }
  letStmt, ok := s.(*ast.LetStatement)
  if !ok {
    t.Errorf("s not *ast.LetStatement. got=%T", s)
    return false
  }
  // Name is Identifier struct.
  if letStmt.Name.Value != name {
    t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
    return false
  }
  if letStmt.Name.TokenLiteral() != name {
    t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
    name, letStmt.Name.TokenLiteral())
    return false
  }
  return true
}

/***** return statement tests ******/
func TestReturnStatements(t *testing.T) {
  input := `
  return 5;
  return 10;
  return 993322;
  `
  l := lexer.New(input)
  p := New(l)
  program := p.ParseProgram()
  checkParserErrors(t, p)
  if len(program.Statements) != 3 {
    t.Fatalf("program.Statements does not contain 3 statements. got=%d",
    len(program.Statements))
  }
  for _, stmt := range program.Statements {
    returnStmt, ok := stmt.(*ast.ReturnStatement)
    if !ok {
      t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
      continue
    }
    if returnStmt.TokenLiteral() != "return" {
      t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
      returnStmt.TokenLiteral())
    }
  }
}
