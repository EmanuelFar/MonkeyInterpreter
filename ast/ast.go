package ast

import (
  "Monkey/token"
  "bytes"
)

/* - statementNode() and expressionNode() are dummy functions so @GO compiler -..
* can easily destiguish between both interfaces.
*
* - #identifier in a let statement is represented as an -..
* expression for simplification.
*
* */


/* AST contains a Node, Node contains either a Statement or an Expression */
type Node interface {
  TokenLiteral() string
  String()       string
}

type Statement interface {
  Node
  statementNode()
}

type Expression interface {
  Node
  expressionNode()
}

type Program struct {
  Statements []Statement
}

func (p *Program) TokenLiteral() string {
  if len(p.Statements) > 0 {
    return p.Statements[0].TokenLiteral()
  }else{
    return ""
  }
}

func (p* Program) String() string {
  var out bytes.Buffer

  for _, s := range p.Statements {
    out.WriteString(s.String())
  }
  return out.String()
}
/*****   let statement    *****/

// <token.LET> <identifier> = <expression>
type LetStatement struct {
  Token token.Token       // token.LET
  Name *Identifier
  Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string  { return ls.Token.Literal }

func (ls *LetStatement) String() string {
  var out bytes.Buffer

  out.WriteString(ls.TokenLiteral() + " ")
  out.WriteString(ls.Name.String())
  out.WriteString(" = ")

  if ls.Value != nil {
    out.WriteString(ls.Value.String())
  }
  out.WriteString(";")
  return out.String()
}

type Identifier struct {
  Token token.Token       // token.IDENT
  Value string
}

// Identifier can possibly be an expression, e.g 'x + x.'
func (id *Identifier) expressionNode() {}
func (id *Identifier) TokenLiteral() string  { return id.Token.Literal }

func (id *Identifier) String() string {
  return id.Value
}


/*****   return statement    *****/

// <token.RETURN> <statement>
type ReturnStatement struct {
  Token       token.Token
  ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {return rs.Token.Literal}

func (rs *ReturnStatement) String() string {
  var out bytes.Buffer
  out.WriteString(rs.TokenLiteral() + " ")
  if rs.ReturnValue != nil {
    out.WriteString(rs.ReturnValue.String())
  }
  out.WriteString(";")
  return out.String()
}
/*****   expression statement   *****/

// a statement that consists solely one expression , e.g 'x + 10;'
type ExpressionStatement struct {
  Token       token.Token
  Expression  Expression
}

// since we implement statement interface, ES can be added to Program.Statements 
func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {return es.Token.Literal}

func (es *ExpressionStatement) String() string {  
  if es.Expression != nil {
    return es.Expression.String()
  }
  return ""
}

/***** prefix expression ******/

type PrefixExpression struct {
  Token     token.Token   // e.g !, -
  Operator  string
  Right     Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
var out bytes.Buffer
  out.WriteString("(")
  out.WriteString(pe.Operator)
  out.WriteString(pe.Right.String())
  out.WriteString(")")
  return out.String()
}


type IntegerLiteral struct {
  Token token.Token
  Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {return il.Token.Literal }
func (il *IntegerLiteral) String() string {return il.Token.Literal }

/***** infix expression *****/

type InfixExpression struct {
  Token   token.Token
  Left    Expression
  Operator  string
  Right   Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
  var out bytes.Buffer
  out.WriteString("(")
  out.WriteString(ie.Left.String())
  out.WriteString(" " + ie.Operator + " ")
  out.WriteString(ie.Right.String())
  out.WriteString(")")
  return out.String()
}

