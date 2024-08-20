package ast

import (
  "Monkey/token"
  "testing"
)

func TestString(t *testing.T) {
  // construction of a program that contains only "let myVar = anotherVar" statement
  program := &Program{
    Statements: []Statement{
      &LetStatement{
        Token: token.Token{Type: token.LET, Literal: "let"},
        Name: &Identifier{
          Token: token.Token{Type: token.IDENT, Literal: "myVar"},
          Value: "myVar",
        },
        Value: &Identifier{
          Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
          Value: "anotherVar",
       },
      },
    },
  }
  if program.String() != "let myVar = anotherVar;" {
    t.Errorf("program.String() wrong. got=%q", program.String())
  }
}
