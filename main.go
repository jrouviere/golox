package main

import (
	"fmt"

	"github.com/jrouviere/golox/parser"
)

func main() {
	const input = `
		var minDepth = 4;
		var maxDepth = 14;
		var stretchDepth = maxDepth + 1;
		if (this.left == nil) {
			return this.item;
		}
		*+,;
		{== != 
		!}-
		<>
		<= >=
		"a long string" == "" != "3.14"
		3.14==6*5+2.5/1234.0
		// this is a comment
		()++
	/3.2`
	scanner := parser.NewScanner(input)
	tokens, err := scanner.Scan()
	if err != nil {
		fmt.Println("Error", err)
	}
	for _, t := range tokens {
		fmt.Println(t.Line, t)
	}
}
