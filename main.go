package main

import "github.com/jrouviere/golox/interpreter"

func main() {
	const input = `
		print "hello" + ", " + "world" + "!";
		var a = "global a";
		var b = "global b";
		{
			var a = "outer a";
			{
				var a = "inner a";
				b = "global b2";
				print a;
				print b;
			}
			print a;
			print b;
		}
		print a;
		print b;
	`

	interp := interpreter.New()

	interp.Run(input)
}
