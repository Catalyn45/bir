package main

func main() {
	var text = `
	module a.b

	import b.c
	import c

	struct a {
		bla: int
		bli: bool
		blu: float
	}

	interface I {
		function f()
		function fu(a: bool, b: bool)
	}

	function functie (a: int, b: int) :int {
	}

	function main() {
		var a: int
	}
`

	lexer := newLexer(text)
	parser := newParser(lexer)

	err, root := parser.Parse()
	if err != nil {
		panic(err)
	}

	root.Dump(4)
}
