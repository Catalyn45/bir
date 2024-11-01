package main

import (
	"fmt"
	"os"
	"os/exec"
)

func parseFile(fileName string) (error, *Node) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	text := string(data)

	lexer := newLexer(text)
	parser := newParser(lexer)

	err, root := parser.Parse()
	if err != nil {
		return err, nil
	}

	fmt.Println("==============================================================================================================")
	root.Dump(0, &[]int{}, "")
	fmt.Println("==============================================================================================================")

	return nil, root
}

func main() {
	// programName := os.Args[0]
	// args := os.Args[1:]
	//
	// if len(args) == 0 {
	// 	fmt.Printf("Usage: %s ./file1.bir [./file2.bir ...]", programName)
	// 	args = []string {"test.bir", "test2.bir"}
	// }
	//
	// var roots []*Node
	// for _, arg := range args {
	// 	err, root := parseFile(arg)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	roots = append(roots, root)
	// }
	//
	// checker := newChecker(roots)
	//
	// err := checker.Check()
	// if err != nil {
	// 	panic(err)
	// }

	err, node := parseFile("testCompile.bir")
	if err != nil {
		panic(err)
	}

	checker := newChecker([]*Node{node})

	err = checker.Check()
	if err != nil {
		panic(err)
	}

	compiler := newCompile([]*Node{node})
	err = compiler.Compile()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("./output.exe")

	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Command finished with exit code: %d\n", exitError.ExitCode())
		} else {
			panic(err)
		}
	}
}
