package main

import (
	"fmt"

	"github.com/jrouviere/golox/parser"
)

func main() {
	const input = `(
		*+,;
		{== != 
		!}-
		// this is a comment
		()++
	`
	scanner := parser.NewScanner(input)

	for _, t := range scanner.Scan() {
		fmt.Println(t.Line, t)
	}
}
