// token/token.go
package token


type TokenType string
// Note: String type is easy to debug -.. but expensive comparing to int/Byte

type Token struct {
  Type TokenType
  Literal string
}

func LookupIdent(ident string) TokenType {
  if tok, exists := keywords[ident]; exists {
    return tok
  }
  return IDENT
}


var keywords = map[string]TokenType {
  "fn"    : FUNCTION,
  "let"   : LET,
  "true"  : TRUE,
  "false" : FALSE,
  "if"    : IF,
  "else"  : ELSE,
  "return": RETURN,
}
// Token types (In monkey we've limited tokens comparing to other languages)
const (
  ILLEGAL = "ILLEGAL" // Unknown token/character
  EOF     = "EOF"

  // Identifiers & Literals
  IDENT   = "IDENT" // x, y...
  INT     = "INT"   // 1,2,3
  STRING  = "STRING"
  // Operators
  ASSIGN  = "="
  PLUS    = "+"
  MINUS   = "-"
  BANG    = "!"
  ASTERISK = "*"
  SLASH   = "/"
  LT      = "<"
  GT      = ">"
  EQ      = "=="
  NOT_EQ  = "!="

  // Delimiters
  COMMA     = ","
  SEMICOLON = ";"

  LPAREN  = "("
  RPAREN  = ")"
  LBRACE  = "{"
  RBRACE  = "}"

  // Keywords
  FUNCTION = "FUNCTION" // Function declaration
  LET      = "LET"      // Variable declaration
  TRUE     = "TRUE"
  FALSE    = "FALSE"
  IF       = "IF"
  ELSE     = "ELSE"
  RETURN   = "RETURN"
)


