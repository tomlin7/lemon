package cli

import (
	"bufio"
	"io"
	"lemon/evaluator"
	"lemon/lexer"
	"lemon/object"
	"lemon/parser"
)

func Exec(fileIn io.Reader, out io.Writer) {
	var input string
	scanner := bufio.NewScanner(fileIn)
	for scanner.Scan() {
		input += scanner.Text()
	}

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	evaluator.Eval(program, env)

	// evaluated := evaluator.Eval(program, env)
	// if evaluated != nil {
	// 	io.WriteString(out, evaluated.Inspect())
	// }
}
