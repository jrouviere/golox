package main

import (
	"fmt"

	"github.com/jrouviere/golox/parser"
)

func main() {
	const input = `
		print "hello" + ", " + "world" + "!";
		print 12+23/3;
		print 8*5/2;
		print 12+23/3 <= 8*5/2;
	`
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
		err := e.Evaluate()
		if err != nil {
			fmt.Println("Error", err)
			return
		}
	}
}
