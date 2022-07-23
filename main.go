package main

import (
	"fmt"

	"github.com/jrouviere/golox/parser"
)

func main() {
	// const input = `
	// 	var minDepth = 4;
	// 	var maxDepth = 14;
	// 	var stretchDepth = maxDepth + 1;
	// 	if (this.left == nil) {
	// 		return this.item;
	// 	}
	// `
	const input = `"hello" + ", " + "world" + "!"`
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
	fmt.Println(expr)

	res, err := expr.Evaluate()
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Println(res)
}
