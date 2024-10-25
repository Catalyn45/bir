package main

import "fmt"

type Symbol struct {
	name       string
	simbolType string
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

func (this *Checker) addSymbol(symbolName string, symbolType string) error {
	for _, symbol := range this.symbolTable[this.currentScope] {
		if symbol.name == symbolName {
			return fmt.Errorf("Symbol already declared in current scope")
		}
	}

	this.symbolTable[this.currentScope] = append(this.symbolTable[this.currentScope], Symbol {
		name: symbolName,
		simbolType: symbolType,
	})

	return nil
}

func (this *Checker) searchSymbolType(symbolName string) (error, string) {
	for i := this.currentScope; i >= 0; i-- {
		symbols := this.symbolTable[this.currentScope]

		for _, symbol := range symbols {
			if symbol.name == symbolName {
				return nil, symbol.simbolType
			}
		}
	}

	return fmt.Errorf("Symbol not declarated"), ""
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

	// if node.nodeType == NODE_BOOL {
	//		return nil, "bool"
	// }

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
	}

	return fmt.Errorf("Can't check type"), ""
}

func (this *Checker) walk(node *Node) error {
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
			return fmt.Errorf("can't assign different types")
		}

		err := this.addSymbol(node.token.tokenValue, variableSymbolType)
		if err != nil {
			return err
		}
	} else if node.nodeType == NODE_FUNCTION {
		symbolName := node.left.token.tokenValue
		err, symbolType := this.getTypeFromNode(node.left.left)
		if err != nil {
			return err
		}

		err = this.addSymbol(symbolName, symbolType)
		if err != nil {
			return err
		}
		
		this.currentScope++

		err = this.walk(node.right)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Invalid node")
	}

	if node.next != nil {
		err := this.walk(node.next)
		if err != nil {
			return err
		}
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