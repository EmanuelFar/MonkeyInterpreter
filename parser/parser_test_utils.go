package parser

import (
  "testing"
)


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

