package parser

import (
  "Monkey/ast"
  "Monkey/lexer"
  "Monkey/token"
  "fmt"
)

const (
  _ int = iota
  LOWEST
  EQUALS      // ==
  LESSGREATER // > OR <
  SUM         // +
  PRODUCT     // *
  PREFIX      // -X OR !X
  CALL        // myFunction(X)
)


type (
  prefixParseFn func() ast.Expression
  infixParseFn  func(ast.Expression) ast.Expression // argument resembles the left expression of a binary operation.
  )

type Parser struct {
  l*        lexer.Lexer

  curToken  token.Token
  peekToken token.Token   // successor of curToken

  // ParseFns take Token and return a function that parses it.
  // according to Pratt's Parsing algorithm
  prefixParseFns map[token.TokenType]prefixParseFn
  infixParseFns  map[token.TokenType]infixParseFn
  errors []string
}

func New(l* lexer.Lexer) *Parser {
  p := &Parser{
    l: l,
    errors: []string{},
  }
  
  p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
  p.registerPrefix(token.IDENT, p.parseIdentifier)

  // Read two tokens to set curToken and peekToken
  p.nextToken()
  p.nextToken()

  return p
}

func (p *Parser) nextToken() {
  p.curToken  = p.peekToken
  p.peekToken = p.l.NextToken()   // lexer calls nextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
  program := &ast.Program{}       // root of AST
  program.Statements = []ast.Statement{}

  for !p.curTokenIs(token.EOF) {
    stmt := p.parseStatement()
    if stmt != nil {
      program.Statements = append(program.Statements, stmt)
    }
    p.nextToken()
  }
  return program 
}

func (p *Parser) parseStatement() ast.Statement {
  switch p.curToken.Type {
    case token.LET:
      return p.parseLetStatement()
    case token.RETURN:
      return p.parseReturnStatement()
    default:
      return p.parseExpressionStatement()
  }
}

/* let statement structure - let <identifier> = <expression> */
func (p* Parser) parseLetStatement() *ast.LetStatement {
  stmt := &ast.LetStatement{Token: p.curToken}
  if !p.expectPeek(token.IDENT) {
    return nil
  }
  stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

  if !p.expectPeek(token.ASSIGN) {
    return nil
  }
  // TODO: We're skipping the expressions until we
  // encounter a semicolon
  for !p.curTokenIs(token.SEMICOLON) {
    p.nextToken()
  }
  return stmt
}

/* return statement structure - return <expression> */
func (p* Parser) parseReturnStatement() *ast.ReturnStatement {
  stmt := &ast.ReturnStatement{Token: p.curToken}

  p.nextToken()

  // TODO: We're skipping the expressions until we
  // encounter a semicolon
  if !p.curTokenIs(token.SEMICOLON) {
    p.nextToken()
  }
 return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
  stmt := &ast.ExpressionStatement{Token: p.curToken}
  stmt.Expression = p.parseExpression(LOWEST)

  if (p.peekTokenIs(token.SEMICOLON)) {
    p.nextToken()
  }
  return stmt
}

func (p *Parser) parseExpression(precendence int) ast.Expression {
  // prefix is a function that parses the given Token.
  prefix := p.prefixParseFns[p.curToken.Type]
  if prefix == nil {
    return nil
  }
  leftExp := prefix()
  return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
  return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}


func (p *Parser) curTokenIs(t token.TokenType) bool {
  return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
  return p.peekToken.Type == t
}

/* primary purpose is to enforce the correctness of the -..
* order of the tokens.
* */
func (p *Parser) expectPeek(t token.TokenType) bool {
  if p.peekTokenIs(t) {
    p.nextToken()
    return true
  }else {
    p.peekError(t)
    return false
  }
}

/***** Maps management *****/
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
  p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
  p.infixParseFns[tokenType] = fn
}

/***** Error Handling *****/
func (p *Parser) Errors() []string {
  return p.errors
}

func (p* Parser) peekError(t token.TokenType) {
  msg := fmt.Sprintf("expected next token to be %s, got %s instead",t, p.peekToken.Type)
  p.errors = append(p.errors, msg)
}
