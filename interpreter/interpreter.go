package interpreter

import (
	"fmt"
	"time"

	"github.com/jrouviere/golox/parser"
)

type Interpreter struct {
	env *parser.Env
}

func New() *Interpreter {
	globals := parser.NewEnv(nil)
	globals.Define("clock", nativeClock{})
	return &Interpreter{
		env: globals,
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

type nativeClock struct{}

func (nativeClock) Call(env *parser.Env, args []interface{}) (interface{}, error) {
	return float64(time.Now().UnixMilli()) / 1000.0, nil
}
func (nativeClock) Arity() int {
	return 0
}

func (nativeClock) String() string {
	return "<nativeFn>"
}
