package main

import "fmt"

type Symbol struct {
	name       string
	simbolType string
	node 	   *Node
}

type Checker struct {
	symbolTable  map[int][]Symbol
	currentScope int
	modules map[string][]*Node
}

func newChecker(asts []*Node) *Checker {
	modules := make(map[string][]*Node)

	for _, ast := range asts {
		medatadaNode := ast.left
		moduleNode := medatadaNode.left
		pathNode := moduleNode.left

		modules[pathNode.token.tokenValue] = append(modules[pathNode.token.tokenValue], ast)
	}

	return &Checker{
		modules:      modules,
		currentScope: 0,
		symbolTable:  make(map[int][]Symbol),
	}
}

func (this *Checker) addSymbol(symbolName string, symbolType string, node *Node) error {
	for _, symbol := range this.symbolTable[this.currentScope] {
		if symbol.name == symbolName {
			return fmt.Errorf("Symbol already declared in current scope")
		}
	}

	this.symbolTable[this.currentScope] = append(this.symbolTable[this.currentScope], Symbol {
		name: symbolName,
		simbolType: symbolType,
		node: node,
	})

	return nil
}

func (this *Checker) searchSymbolType(symbolName string) (error, string) {
	err, symbol := this.searchSymbol(symbolName)
	if err != nil {
		return err, ""
	}

	return nil, symbol.simbolType
}

func (this *Checker) searchSymbol(symbolName string) (error, *Symbol) {
	for i := this.currentScope; i >= 0; i-- {
		symbols := this.symbolTable[i]

		for _, symbol := range symbols {
			if symbol.name == symbolName {
				return nil, &symbol
			}
		}
	}

	return fmt.Errorf("Symbol not declarated"), nil
}

func (this *Checker) getTypeFromNode(node *Node) (error, string) {
	if node.nodeType == NODE_INT_TYPE {
		return nil, "int"
	}

	if node.nodeType == NODE_FLOAT_TYPE {
		return nil, "float"
	}

	if node.nodeType == NODE_STRING_TYPE {
		return nil, "string"
	}

	if node.nodeType == NODE_BOOL_TYPE {
		return nil, "bool"
	}

	if node.nodeType == NODE_CUSTOM_TYPE {
		return nil, node.token.tokenValue
	}

	return fmt.Errorf("Invalid node"), ""
}

func (this *Checker) determineType(node *Node) (error, string) {
	if node.nodeType == NODE_INT {
		return nil, "int"
	}

	if node.nodeType == NODE_FLOAT {
		return nil, "float"
	}

	if node.nodeType == NODE_BOOL {
		return nil, "bool"
	}

	if node.nodeType == NODE_STRING {
		return nil, "string"
	}

	if node.nodeType == NODE_VARIABLE {
		err, symbolType := this.searchSymbolType(node.token.tokenValue)
		if err != nil {
			return err, ""
		}

		return nil, symbolType
	}

	if node.nodeType == NODE_NOT {
		err, symbolType := this.determineType(node.left)
		if err != nil {
			return err, ""
		}

		if symbolType != "bool" {
			return fmt.Errorf("Can't apply not on non bool type"), ""
		}

		return nil, symbolType
	}

	if node.nodeType == NODE_BINARY_EXPRESSION {
		err, typeLeft := this.determineType(node.left)
		if err != nil {
			return err, ""
		}

		err, typeRight := this.determineType(node.right)
		if err != nil {
			return err, ""
		}

		if typeLeft != typeRight {
			return fmt.Errorf("invalid operation between different types: %s and %s", typeLeft, typeRight), ""
		}

		return nil, typeLeft
	}

	if node.nodeType == NODE_CALL {
		err, functionSymbol := this.searchSymbol(node.left.token.tokenValue)
		if err != nil {
			return err, ""
		}

		functionNode := functionSymbol.node

		var parameterTypes []string
		for parameter := functionNode.left.right; parameter != nil; parameter = parameter.next {
			err, parameterType := this.getTypeFromNode(parameter.left)
			if err != nil {
				return err, ""
			}

			parameterTypes = append(parameterTypes, parameterType)
		}

		var argumentTypes []string
		for argument := node.right.right; argument != nil; argument = argument.next {
			err, argumentType := this.determineType(argument)
			if err != nil {
				return err, ""
			}

			argumentTypes = append(argumentTypes, argumentType)
		}

		if len(parameterTypes) != len(argumentTypes) {
			return fmt.Errorf("Not the same number of arguments"), ""
		}

		for i := 0; i < len(parameterTypes); i++ {
			if parameterTypes[i] != argumentTypes[i] {
				return fmt.Errorf("Invalid argument type for parameter"), ""
			}
		}

		return nil, functionSymbol.simbolType
	}

	if node.nodeType == NODE_RETURN {
		return this.determineType(node.left)
	}

	return fmt.Errorf("Can't check type"), ""
}

func (this *Checker) enterScope() {
	this.currentScope ++
}

func (this *Checker) leaveScope() {
	delete(this.symbolTable, this.currentScope)

	this.currentScope --
}

func (this *Checker) walk(node *Node) error {
	this.enterScope()
	defer this.leaveScope()

	for node != nil {
		if node.nodeType == NODE_VARIABLE_DECLARATION {
			var initializationSymbolType string
			if node.right != nil {
				var err error
				err, initializationSymbolType = this.determineType(node.right)
				if err != nil {
					return err
				}
			}

			var variableSymbolType string
			if node.left != nil {
				err, symbolType := this.getTypeFromNode(node.left)
				if err != nil {
					return err
				}

				variableSymbolType = symbolType
			}

			if variableSymbolType == "" {
				variableSymbolType = initializationSymbolType
			}

			if variableSymbolType != initializationSymbolType {
				return fmt.Errorf("can't initialize with different types")
			}

			err := this.addSymbol(node.token.tokenValue, variableSymbolType, node)
			if err != nil {
				return err
			}
		} else if node.nodeType == NODE_FUNCTION {
			symbolName := node.left.token.tokenValue
			err, symbolType := this.getTypeFromNode(node.left.left.left)
			if err != nil {
				return err
			}

			err = this.addSymbol(symbolName, symbolType, node)
			if err != nil {
				return err
			}
			
			err = this.walk(node.right)
			if err != nil {
				return err
			}
		} else if node.nodeType == NODE_IF {
			err, symbolType := this.determineType(node.left)
			if err != nil {
				return err
			}

			if symbolType != "bool" {
				return fmt.Errorf("Can't have non-bool in if")
			}

			branchNode := node.right

			err = this.walk(branchNode.left)
			if err != nil {
				return err
			}

			err = this.walk(branchNode.right)
			if err != nil {
				return err
			}
		} else if node.nodeType == NODE_ASSIGNMENT {
			err, leftSymbolType := this.determineType(node.left)
			if err != nil {
				return err
			}

			err, rightSymbolType := this.determineType(node.right)
			if err != nil {
				return err
			}

			if leftSymbolType != rightSymbolType {
				return fmt.Errorf("Can't assign different types")
			}
		} else if node.nodeType == NODE_RETURN {
			err, _ := this.determineType(node)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Invalid node")
		}

		node = node.next
	}

	return nil
}

func (this *Checker) Check() error {
	for _, asts := range this.modules {
		err := this.walk(asts[0].right)
		if err != nil {
			return err
		}
	}

	return nil
}