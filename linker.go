package main

import (
	"os"
	"os/exec"
)

type Linker struct {
	programName string
	compilers []*Compiler
}

func newLinker(asts []*Node, programName string) *Linker {
	// map asts by module name
	modules := make(map[string][]*Node)
	for _, ast := range asts {
		medatadaNode := ast.left
		moduleNode := medatadaNode.left
		pathNode := moduleNode.left

		modules[pathNode.token.tokenValue] = append(modules[pathNode.token.tokenValue], ast)
	}

	// create compilers
	var compilers []*Compiler
	for moduleName, asts := range modules {
		compiler := newCompiler(asts, moduleName)
		compilers = append(compilers, compiler)
	}

	return &Linker{
		programName: programName,
		compilers: compilers,
	}
}

func (this *Linker) Link() error {
	var outputs []string
	// run the compilers
	for _, compiler := range this.compilers {
		err := compiler.Compile()
		if err != nil {
			return err
		}

		outputs = append(outputs, *compiler.moduleName+".obj")
	}

	arguments := append([]string{}, outputs...)
	arguments = append(arguments, "-o")
	arguments = append(arguments, this.programName)

	// run clang link
	cmd := exec.Command("clang", arguments...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
