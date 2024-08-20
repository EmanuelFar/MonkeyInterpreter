// Read Eval Print Loop
package repl

import (
  "bufio"
  "fmt"
  "io"
  "Monkey/lexer"
  "Monkey/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
  var scanner *bufio.Scanner = bufio.NewScanner(in)

  for {
    fmt.Fprint(out, PROMPT)
    var scanned bool = scanner.Scan()
    if !scanned {
      return
    }

    line := scanner.Text()
    l := lexer.New(line)

    for tok:= l.NextToken(); tok.Type != token.EOF ; tok = l.NextToken() {
      fmt.Fprintf(out, "%+v\n",tok)
    }
  }
}
