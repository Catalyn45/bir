package main

func main ()  {
	var text = `
	var a : int = 10
	a += 10
	b += 20

	struct a {
		a: int
		b: bool
		b: float
	}

	interface I {
		f()
		fu(a: bool, b: bool)
	}

	func functie (a: int, b: int) :int {
		return a + b
	}
	`

	lexer := newLexer(text)
	for {
		err, token := lexer.next()
		if err != nil {
			panic(err)
		}

		print(token.toString(), "\n")

		if token.tokenType == TOKEN_EOF {
			break
		}
	}
}