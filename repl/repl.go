// Read Eval Print Loop
package repl

import (
  "bufio"
  "fmt"
  "io"
  "Monkey/lexer"
  "Monkey/parser"
  "Monkey/evaluator"
  "Monkey/object"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
  var scanner *bufio.Scanner = bufio.NewScanner(in)
  env := object.NewEnvironment()
  for {
    fmt.Fprint(out, PROMPT)
    var scanned bool = scanner.Scan()
    if !scanned {
      return
    }

    line := scanner.Text()
    l := lexer.New(line)
    p := parser.New(l)
    program := p.ParseProgram()

    if len(p.Errors()) != 0 {
      printParserErrors(out, p.Errors())
      continue
    }

    evaluated := evaluator.Eval(program, env)
    if evaluated != nil {
      io.WriteString(out, evaluated.Inspect())
      io.WriteString(out, "\n")
    }
  }
}

func printParserErrors(out io.Writer, errors []string) {  
  for _, msg := range errors {
    io.WriteString(out, "\t"+msg+"\n")
 }
}
