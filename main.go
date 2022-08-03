package main

import "github.com/jrouviere/golox/interpreter"

func main() {
	const input = `
		print "hello" + ", " + "world" + "!";
		var a = 42;
		var b = 2*21;
		if (true) {
			if (a == b) {
				print "equals";
			} else {
				print "not equals";
			}
		}
		print nil or a;
		print "hello" and "world" or 42==12;
		print 1==2 or 33;
		var i = 0;
		while (i < 10) {
			print i;
			i = i + 1;
		}
	`

	interp := interpreter.New()

	interp.Run(input)
}
