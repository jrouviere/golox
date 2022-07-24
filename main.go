package main

import "github.com/jrouviere/golox/interpreter"

func main() {
	const input = `
		print "hello" + ", " + "world" + "!";
		var a = 12+23/3;
		var b = 8*5/2;
		a = a - 1;
		print a;
		print b;
		print a <= b;
	`

	interp := interpreter.New()

	interp.Run(input)
}
