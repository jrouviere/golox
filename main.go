package main

import "github.com/jrouviere/golox/interpreter"

func main() {
	const input = `
		var t1 = clock();
		fun fib(n) {
			if (n <= 1) return n;
			return fib(n-2)+fib(n-1);
		}

		for (var i=0; i<30; i=i+1) {
			print fib(i);
		}

		print clock() - t1;
	`

	interp := interpreter.New()

	interp.Run(input)
}
