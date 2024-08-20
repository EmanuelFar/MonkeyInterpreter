package lexer

import "Monkey/token"

type Lexer struct {
  input         string
  position      int     // current position  in input (current char)
  readPosition  int     // current reading position   (after current char)
  ch            byte    // current char that's being examined
}

func New(input string) *Lexer {
  var l *Lexer = &Lexer{input: input}
  l.readChar()
  return l
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
  return token.Token{Type: tokenType, Literal: string(ch)}
}

/* reads a character as it advances position and readPosition. */
func (l *Lexer) readChar(){
  if l.readPosition >= len(l.input) {
    l.ch = 0
  } else {
    l.ch = l.input[l.readPosition]
  }
  l.position = l.readPosition
  l.readPosition += 1
}

/* returns next character in the input without incrementing position */
func (l* Lexer) peekChar() byte{
  if (l.readPosition >= len(l.input)) {
    return 0;
  }
  return l.input[l.readPosition]
}

/* reads identifier and advances lexer's position until it -..
* encouters a non-letter-character.
* @return a string represeting the identifier. */
func (l* Lexer) readIdentifier() string {
  position := l.position
  for isLetter(l.ch){
    l.readChar()
  }
  return l.input[position:l.position]
}

func (l* Lexer) readNumber() string {
  position := l.position
  for isDigit(l.ch) {
    l.readChar()
  }
  return l.input[position:l.position]
}

/* turns l.ch into it's compatible token. */
func (l* Lexer) NextToken() token.Token {
  var tok token.Token
  l.skipWhitespace()

  // Token classification
  switch l.ch {
    case '=':
        if l.peekChar() == '='{
          char := l.ch
          l.readChar()
          literal := string(char) + string(l.ch)
          tok = token.Token{Type : token.EQ, Literal: literal}
        }else{
          tok = newToken(token.ASSIGN, l.ch)
        }
    case '!':
        if l.peekChar() == '='{
          char := l.ch
          l.readChar()
          literal := string(char) + string(l.ch)
          tok = token.Token{Type : token.NOT_EQ, Literal: literal}
        }else{
          tok = newToken(token.BANG, l.ch)
        }
    case ';':
        tok = newToken(token.SEMICOLON, l.ch)
    case '(':
        tok = newToken(token.LPAREN, l.ch)
    case ')':
        tok = newToken(token.RPAREN, l.ch)
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case '+':
        tok = newToken(token.PLUS, l.ch)
    case '{':
        tok = newToken(token.LBRACE, l.ch)
    case '}':
        tok = newToken(token.RBRACE, l.ch)
    case '-':
        tok = newToken(token.MINUS, l.ch)
    case '/':
        tok = newToken(token.SLASH, l.ch)
    case '*':
        tok = newToken(token.ASTERISK, l.ch)
    case '<':
        tok = newToken(token.LT, l.ch)
    case '>':
        tok = newToken(token.GT, l.ch)
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
    default:
        if isLetter(l.ch){
          tok.Literal = l.readIdentifier()
          tok.Type = token.LookupIdent(tok.Literal)
          return tok
        }else if isDigit(l.ch){
          tok.Type = token.INT 
          tok.Literal = l.readNumber()
          return tok
        }else {
          tok = newToken(token.ILLEGAL, l.ch)
        }
  }
  l.readChar()
  return tok
}

func (l* Lexer) skipWhitespace() {
  for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r'{
    l.readChar()
  }
}

// NOTE: Adding/Removing checks would re-arrange identifier and keywords span.
func isLetter(ch byte) bool {
  return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(char byte) bool {
  return '0' <= char && char <= '9'
}
