package parser

import (
  "Monkey/ast"
  "Monkey/lexer"
  "Monkey/token"
  "fmt"
  "strconv"
)

var precendences = map[token.TokenType]int {
  token.EQ:       EQUALS,
  token.NOT_EQ:   EQUALS,
  token.LT:       LESSGREATER,
  token.GT:       LESSGREATER,
  token.PLUS:     SUM,
  token.MINUS:    SUM,
  token.SLASH:    PRODUCT,
  token.ASTERISK: PRODUCT,
  token.LPAREN:   CALL,
}

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
  peekToken token.Token   // successor of curToken.

  // ParseFns take Token and return a function that parses it.
  // according to Pratt's Parsing algorithm.
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
  p.registerPrefix(token.INT, p.parseIntegerLiteral)
  p.registerPrefix(token.BANG, p.parsePrefixExpression)
  p.registerPrefix(token.MINUS, p.parsePrefixExpression)
  p.registerPrefix(token.TRUE, p.parseBoolean)
  p.registerPrefix(token.FALSE, p.parseBoolean)
  p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
  p.registerPrefix(token.IF, p.parseIfExpression)
  p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
  p.registerPrefix(token.STRING,  p.parseStringLiteral)

  p.infixParseFns = make(map[token.TokenType]infixParseFn)
  p.registerInfix(token.PLUS, p.parseInfixExpression)
  p.registerInfix(token.MINUS, p.parseInfixExpression)
  p.registerInfix(token.SLASH, p.parseInfixExpression)
  p.registerInfix(token.ASTERISK, p.parseInfixExpression)
  p.registerInfix(token.EQ, p.parseInfixExpression)
  p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
  p.registerInfix(token.LT, p.parseInfixExpression)
  p.registerInfix(token.GT, p.parseInfixExpression)
  p.registerInfix(token.LPAREN, p.parseCallExpression)

  // Read two tokens to set curToken and peekToken.
  p.nextToken()
  p.nextToken()

  return p
}

func (p *Parser) nextToken() {
  p.curToken  = p.peekToken
  p.peekToken = p.l.NextToken()   // lexer calls nextToken().
}

func (p *Parser) ParseProgram() *ast.Program {
  program := &ast.Program{}       // root of AST.
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
  p.nextToken()
  stmt.Value = p.parseExpression(LOWEST)

  if _, ok := stmt.Value.(*ast.FunctionLiteral); ok {
    return stmt
  }

  for !p.curTokenIs(token.SEMICOLON) {
    p.nextToken()
  }
  return stmt
}

/* return statement structure - return <expression> */
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
  stmt := &ast.ReturnStatement{Token: p.curToken}
  p.nextToken()
  stmt.ReturnValue = p.parseExpression(LOWEST)

  if p.peekTokenIs(token.SEMICOLON) {
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

// this method that kicks off expression parsing
func (p *Parser) parseExpression(precendence int) ast.Expression {
  // prefix is a function that parses the given Token.
  var prefix prefixParseFn = p.prefixParseFns[p.curToken.Type]
  if prefix == nil {
    p.noPrefixParseFnError(p.curToken.Type)
    return nil
  }
  var leftExp ast.Expression = prefix()
  
  for !p.peekTokenIs(token.SEMICOLON) && precendence < p.peekPrecedence() {
    var infix infixParseFn = p.infixParseFns[p.peekToken.Type]
    if infix == nil {        // no infix expression found
      return leftExp
    }

    p.nextToken()

    leftExp = infix(leftExp) // leftExp holds left expression of the infix expression
  }

  return leftExp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
  expression := &ast.PrefixExpression{
    Token: p.curToken,
    Operator: p.curToken.Literal,
  }
  // consume another token to parse prefix expression.
  p.nextToken()

  // integer parsing.
  expression.Right = p.parseExpression(PREFIX)

  return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
  expression := &ast.InfixExpression{
    Token: p.curToken,
    Operator: p.curToken.Literal,
    Left: left,
  }
  // before curPrecedence is called, curToken is the operator.
  precedence := p.curPrecedence()
  p.nextToken()
  expression.Right = p.parseExpression(precedence)

  return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
  p.nextToken()

  exp := p.parseExpression(LOWEST)

  if !p.expectPeek(token.RPAREN) {
    return nil
  }
  return exp
}

/***** if else expressions parsing *****/

func (p *Parser) parseIfExpression() ast.Expression {
  expression := &ast.IfExpression{Token: p.curToken}
  if !p.expectPeek(token.LPAREN) {
    return nil
  }
  p.nextToken()
  expression.Condition = p.parseExpression(LOWEST)
  if !p.expectPeek(token.RPAREN) {
    return nil
  }
  if !p.expectPeek(token.LBRACE) {
    return nil
  }
  expression.Consequence = p.parseBlockStatement()

  if p.peekTokenIs(token.ELSE) {
    p.nextToken()

    if !p.expectPeek(token.LBRACE) {
      return nil
    }
    expression.Alternative = p.parseBlockStatement()
  }
  return expression
}

/***** Functions parsing *****/

func (p *Parser) parseFunctionLiteral() ast.Expression {
  lit := &ast.FunctionLiteral{Token: p.curToken}
  if !p.expectPeek(token.LPAREN) {
    return nil
  }
  lit.Parameters = p.parseFunctionParameters()
  if !p.expectPeek(token.LBRACE) {
    return nil
  }
  lit.Body = p.parseBlockStatement()
    return lit
  }

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
  identifiers := []*ast.Identifier{}

  if p.peekTokenIs(token.RPAREN) {
    p.nextToken()
    return identifiers
  }
  p.nextToken()

  ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
  identifiers = append(identifiers, ident)

  for p.peekTokenIs(token.COMMA) {
    p.nextToken()
    p.nextToken()

    ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
    identifiers = append(identifiers, ident)
  }

  if !p.expectPeek(token.RPAREN) {
    return nil
  }
  return identifiers
}

// Function calls expressions
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
  exp := &ast.CallExpression{Token: p.curToken, Function: function}
  exp.Arguments = p.parseCallArguments()
  return exp
  }

func (p *Parser) parseCallArguments() []ast.Expression {
  args := []ast.Expression{}
  if p.peekTokenIs(token.RPAREN) {
    p.nextToken()
    return args
  }
  p.nextToken()
  args = append(args, p.parseExpression(LOWEST))
  for p.peekTokenIs(token.COMMA) {
    p.nextToken()
    p.nextToken()
    args = append(args, p.parseExpression(LOWEST))
  }
  if !p.expectPeek(token.RPAREN) {
    return nil
  }
  return args
}

/***** if-else and functions body are represented as BlockStatement  ******/

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
  block := &ast.BlockStatement{Token: p.curToken}
  block.Statements = []ast.Statement{}
  p.nextToken()
  for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
    stmt := p.parseStatement()
    if stmt != nil {
      block.Statements = append(block.Statements, stmt)
    }
    p.nextToken()
  }
  return block
}
/***** Tokens parsing *****/

func (p *Parser) parseStringLiteral() ast.Expression {
  return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIdentifier() ast.Expression {
  return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
  return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}
func (p *Parser) parseIntegerLiteral() ast.Expression {
  lit := &ast.IntegerLiteral{Token: p.curToken}

  value, error := strconv.ParseInt(p.curToken.Literal, 0, 64)
  if error != nil {
    msg := fmt.Sprintf("count not parse %q as integer", p.curToken.Literal)
    p.errors = append(p.errors, msg)
    return nil
  }
  lit.Value = value

  return lit
}

/***** Tokens type check *****/

func (p *Parser) curTokenIs(t token.TokenType) bool {
  return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
  return p.peekToken.Type == t
}

/* primary purpose is to enforce the correctness of the -..
* order of the tokens.
*  - expectPeek advances token.
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

/***** precendence check ******/

func (p *Parser) peekPrecedence() int {
  if p, ok := precendences[p.peekToken.Type]; ok {
    return p
  }
  return LOWEST
}

func (p *Parser) curPrecedence() int {
  if p, ok := precendences[p.curToken.Type]; ok {
    return p
  }
  return LOWEST
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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
  msg := fmt.Sprintf("no prefix parse function for %s found", t)
  p.errors = append(p.errors, msg)
}


