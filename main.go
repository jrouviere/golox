package main

import "github.com/jrouviere/golox/interpreter"

func main() {
	const input = `
		var t1 = clock();
		fun add(a, b) {
			var c = a + b;
			print c;
		}
		fun sayHello(name) {
			print "Hello, " + name;
		}
		print add;
		add(1,2);

		sayHello("world!");
	`

	interp := interpreter.New()

	interp.Run(input)
}
