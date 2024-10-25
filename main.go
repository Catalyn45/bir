package main

import "os"

func main() {
	data, err := os.ReadFile("test.bir")
	if err != nil {
		panic(err)
	}

	text := string(data)

	lexer := newLexer(text)
	parser := newParser(lexer)

	err, root := parser.Parse()
	if err != nil {
		panic(err)
	}

	root.Dump(0, &[]int{}, "")

	checker := newChecker([]*Node {root})

	err = checker.Check()
	if err != nil {
		panic(err)
	}
}
