package repl

import (
	"bufio"
	"fmt"
	"io"
	"lemon/evaluator"
	"lemon/lexer"
	"lemon/object"
	"lemon/parser"
)

const PROMPT = ">> "
const CONTINUE_PROMPT = ".. "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		var lines []string
		fmt.Printf(PROMPT)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				break
			}
			lines = append(lines, line)
			fmt.Printf(CONTINUE_PROMPT)
		}
		if len(lines) == 0 {
			return
		}

		input := ""
		for _, line := range lines {
			input += line + "\n"
		}

		l := lexer.New(input)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())

			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		evaluated := evaluator.Eval(expanded, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

		// io.WriteString(out, program.String())
		// io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
