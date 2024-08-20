package parser

import (
  "testing"
  "Monkey/ast"
  "Monkey/lexer"
)

/*****  Let Statement tests *****/

func TestLetStatements(t *testing.T){
  input :=`
  let x = 5;
  let y = 10;
  let z = 838383;
  `

  lex := lexer.New(input)
  parse := New(lex)

  program := parse.ParseProgram()
  checkParserErrors(t, parse)
  if program == nil {
    t.Fatalf("ParseProgram() returned nil")
  }

  if len(program.Statements) != 3 {
    t.Fatalf("Program statements does not contain 3 statements. got %d", len(program.Statements))
  }
 
  tests := []struct {
    expectedIdentifier string
  }{
    {"x"},
    {"y"},
    {"z"},
  }

  for i,tt := range tests {
    stmt := program.Statements[i]
    if !testLetStatement(t, stmt, tt.expectedIdentifier) {
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

func checkParserErrors(t *testing.T, p *Parser) {
  errors := p.Errors()
  if len(errors) == 0 {
    return
  }

  t.Errorf("parser has %d errors",len(errors))

  for _, msg := range errors {
    t.Errorf("parser errors: %q", msg)
  }

  t.FailNow()
}

