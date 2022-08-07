package main

import "github.com/jrouviere/golox/interpreter"

func main() {
	const input = `
		var t1 = clock();
		print "hello" + ", " + "world" + "!";
		var a = 0;
		var temp;
		for (var b = 1; a < 10000000; b = temp + b) {
			print a;
			temp = a;
			a = b;
		}
		var t2 = clock();
		print t2 - t1;
	`

	interp := interpreter.New()

	interp.Run(input)
}
