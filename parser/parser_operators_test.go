package parser

import (
  "testing"
  "Monkey/lexer"
)

func TestOperatorPrecedenceParsing(t *testing.T) {
  tests := []struct {
    input string
    expected string
  }{
    {
      "a + add(b * c) + d",
      "((a + add((b * c))) + d)",
    },
    {
      "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
      "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
    },
    {
      "add(a + b + c * d / f + g)",
      "add((((a + b) + ((c * d) / f)) + g))",
    },
    {
      " 1 + (2 + 3) + 4",
      "((1 + (2 + 3)) + 4)",
    },
    {
      "(5 + 5) * 2",
      "((5 + 5) * 2)",
    },
    {
      "2 / (5 + 5)",
      "(2 / (5 + 5))",
    },
    {
      "-(5 + 5)",
      "(-(5 + 5))",
    },
    {
      "!(true == true)",
      "(!(true == true))",
    },
    {
      "false",
      "false",
    },
    {
      "true",
      "true",
    },
    {
      "3 > 5 == false",
      "((3 > 5) == false)",
    },
    {
      "3 < 5 == true",
      "((3 < 5) == true)",
    },
    {
      "-a * b",
      "((-a) * b)",
    },
    {
      "!-a",
      "(!(-a))",
    },
    {
      "a + b + c",
      "((a + b) + c)",
    },
    {
      "a + b - c",
      "((a + b) - c)",
    },
    {
      "a * b * c",
      "((a * b) * c)",
    },
    {
      "a * b / c",
      "((a * b) / c)",
    },
    {
      "a + b / c",
      "(a + (b / c))",
    },
    {
      "a + b * c + d / e - f",
      "(((a + (b * c)) + (d / e)) - f)",
    },
    {
      "3 + 4; -5 * 5",
      "(3 + 4)((-5) * 5)",
    },
    {
      "5 > 4 == 3 < 4",
      "((5 > 4) == (3 < 4))",
    },
    {
      "5 < 4 != 3 > 4",
      "((5 < 4) != (3 > 4))",
    },
    {
      "3 + 4 * 5 == 3 * 1 + 4 * 5",
      "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
    },
  }
  for _, tt := range tests {
    l := lexer.New(tt.input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)
    actual := program.String()
    if actual != tt.expected {
      t.Errorf("expected=%q, got=%q", tt.expected, actual)
    }
  }
}
