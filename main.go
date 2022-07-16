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
	const input = `5*(1+3)/2.1 == 1+false`
	scanner := parser.NewScanner(input)
	tokens, err := scanner.Scan()
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, t := range tokens {
		fmt.Println(t.Line, t)
	}

	expr := parser.New(tokens).Parse()
	fmt.Println(expr)
}
