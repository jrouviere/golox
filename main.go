package main

import "github.com/jrouviere/golox/interpreter"

func main() {
	const input = `
		print "hello" + ", " + "world" + "!";
		var a = 0;
		var temp;
		for (var b = 1; a < 10000; b = temp + b) {
			print a;
			temp = a;
			a = b;
		}
	`

	interp := interpreter.New()

	interp.Run(input)
}
