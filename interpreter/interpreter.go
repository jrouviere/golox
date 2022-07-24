package interpreter

import (
	"fmt"

	"github.com/jrouviere/golox/parser"
)

type Interpreter struct {
	env *parser.Env
}

func New() *Interpreter {
	return &Interpreter{
		env: parser.NewEnv(),
	}
}

func (i *Interpreter) Run(input string) {
	scanner := parser.NewScanner(input)
	tokens, err := scanner.Scan()
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	for _, t := range tokens {
		fmt.Println(t.Line, t)
	}

	expr, err := parser.New(tokens).Parse()
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	for _, e := range expr {
		fmt.Println(e)
	}

	for _, e := range expr {
		err := e.Evaluate(i.env)
		if err != nil {
			fmt.Println("Error", err)
			return
		}
	}
}
