package main

import "fmt"

const (
	TYPE_LITERAL  = iota
	TYPE_VARIABLE = iota
	TYPE_FUNCTION = iota
	TYPE_STRUCT   = iota
	TYPE_INTERFACE = iota
)

type SymbolType struct {
	kind int
	name string
	signature []string
}

type Symbol struct {
	name       string
	simbolType SymbolType
	node 	   *Node
	children   []Symbol
}

type Checker struct {
	symbolTable  map[int][]Symbol
	currentScope int
	modules map[string][]*Node
	functionStack *Stack[*Symbol]
}

func newChecker(asts []*Node) *Checker {
	modules := make(map[string][]*Node)

	for _, ast := range asts {
		medatadaNode := ast.left
		moduleNode := medatadaNode.left
		pathNode := moduleNode.left

		modules[pathNode.token.tokenValue] = append(modules[pathNode.token.tokenValue], ast)
	}

	return &Checker {
		modules:      modules,
		currentScope: 0,
		symbolTable:  make(map[int][]Symbol),
		functionStack: &Stack[*Symbol]{},
	}
}

func (this *Checker) symbolAlreadyExists(symbolName string) bool {
	for _, symbol := range this.symbolTable[this.currentScope] {
		if symbol.name == symbolName {
			return true
		}
	}

	return false
}

func (this *Checker) addVariableSymbol(varName string, varType string, node *Node) error {
	if this.symbolAlreadyExists(varName) {
		return fmt.Errorf("Symbol already declared in current scope")
	}

	this.symbolTable[this.currentScope] = append(this.symbolTable[this.currentScope], Symbol {
		name: varName,
		simbolType: SymbolType {
			kind: TYPE_VARIABLE,
			name: varType,
		},
		node: node,
	})

	return nil
}

func (this *Checker) addFunctionSymbol(functionName string, returnType string, signature []string, node *Node) (error, *Symbol) {
	if this.symbolAlreadyExists(functionName) {
		return fmt.Errorf("Symbol already declared in current scope"), nil
	}

	this.symbolTable[this.currentScope] = append(this.symbolTable[this.currentScope], Symbol {
		name: functionName,
		simbolType: SymbolType {
			kind: TYPE_FUNCTION,
			name: returnType,
			signature: signature,
		},
		node: node,
	})

	return nil, &this.symbolTable[this.currentScope][len(this.symbolTable[this.currentScope]) - 1]
}

func (this *Checker) addStructSymbol(structName string, node *Node, children []Symbol) error {
	if this.symbolAlreadyExists(structName) {
		return fmt.Errorf("Symbol already declared in current scope")
	}

	this.symbolTable[this.currentScope] = append(this.symbolTable[this.currentScope], Symbol {
		name: structName,
		simbolType: SymbolType {
			kind: TYPE_STRUCT,
			name: structName,
		},
		node: node,
		children: children,
	})

	return nil
}

func (this *Checker) searchSymbolType(symbolName string) (error, *SymbolType) {
	err, symbol := this.searchSymbol(symbolName)
	if err != nil {
		return err, nil
	}

	return nil, &symbol.simbolType
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

func (this *Checker) expressionAllowed(node *Node, expressionType string) bool {
	if node.token.tokenType == TOKEN_EQUAL || node.token.tokenType == TOKEN_DIFFERENT {
		switch expressionType {
		case "int":
			fallthrough
		case "float":
			fallthrough
		case "string":
			fallthrough
		case "bool":
			return true
		}

		return false
	}

 	if node.token.tokenType == TOKEN_LESS ||
		node.token.tokenType == TOKEN_LESS_EQUAL ||
		node.token.tokenType == TOKEN_GREATER ||
		node.token.tokenType == TOKEN_GREATER_EQUAL {
		switch expressionType {
		case "int":
			fallthrough
		case "float":
			return true
		}

		return false
	}

 	if node.token.tokenType == TOKEN_AND ||
		node.token.tokenType == TOKEN_OR {
		switch expressionType {
		case "bool":
			return true
		}

		return false
	}

	return true
}

func (this *Checker) expressionResultType(node *Node, expressionType *SymbolType) (error, *SymbolType) {
	if node.token.tokenType == TOKEN_EQUAL {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "bool"}
	}

	if node.token.tokenType == TOKEN_DIFFERENT {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "bool"}
	}

	if node.token.tokenType == TOKEN_GREATER {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "bool"}
	}

	if node.token.tokenType == TOKEN_GREATER_EQUAL {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "bool"}
	}

	if node.token.tokenType == TOKEN_LESS {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "bool"}
	}

	if node.token.tokenType == TOKEN_LESS_EQUAL {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "bool"}
	}

	if node.token.tokenType == TOKEN_AND {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "bool"}
	}

	if node.token.tokenType == TOKEN_OR {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "bool"}
	}

	return nil, expressionType
}

func (this *Checker) determineType(node *Node) (error, *SymbolType) {
	if node.nodeType == NODE_INT {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "int"}
	}

	if node.nodeType == NODE_FLOAT {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "float"}
	}

	if node.nodeType == NODE_BOOL {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "bool"}
	}

	if node.nodeType == NODE_STRING {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "string"}
	}

	if node.nodeType == NODE_VARIABLE {
		err, symbolType := this.searchSymbolType(node.token.tokenValue)
		if err != nil {
			return err, nil
		}

		return nil, symbolType
	}

	if node.nodeType == NODE_NOT {
		err, symbolType := this.determineType(node.left)
		if err != nil {
			return err, nil
		}

		if symbolType.name != "bool" {
			return fmt.Errorf("Can't apply not on non bool type"), nil
		}

		return nil, symbolType
	}

	if node.nodeType == NODE_BINARY_EXPRESSION {
		err, typeLeft := this.determineType(node.left)
		if err != nil {
			return err, nil
		}

		err, typeRight := this.determineType(node.right)
		if err != nil {
			return err, nil
		}

		if typeLeft.name != typeRight.name {
			return fmt.Errorf("invalid operation between different types: %s and %s", typeLeft.name, typeRight.name), nil
		}

		if !this.expressionAllowed(node, typeLeft.name) {
			return fmt.Errorf("expression not allowed for type %s", typeLeft.name), nil
		}

		return this.expressionResultType(node, typeLeft)
	}

	if node.nodeType == NODE_CALL {
		err, symbolType := this.determineType(node.left)
		if err != nil {
			return err, nil
		}

		if symbolType.kind == TYPE_STRUCT {
			// constructor
			return nil, symbolType
		}

		if symbolType.kind != TYPE_FUNCTION {
			return fmt.Errorf("Only functions can be called"), nil
		}

		parameterTypes := symbolType.signature

		var argumentTypes []string
		for argument := node.right.right; argument != nil; argument = argument.next {
			err, argumentType := this.determineType(argument)
			if err != nil {
				return err, nil
			}

			argumentTypes = append(argumentTypes, argumentType.name)
		}

		if len(parameterTypes) != len(argumentTypes) {
			return fmt.Errorf("Not the same number of arguments"), nil
		}

		for i := 0; i < len(parameterTypes); i++ {
			if parameterTypes[i] != argumentTypes[i] {
				return fmt.Errorf("Invalid argument type for parameter"), nil
			}
		}

		return nil, symbolType
	}

	if node.nodeType == NODE_MEMBER_ACCESS {
		err, memberType := this.determineType(node.left)
		if err != nil {
			return err, nil
		}

		err, symbol := this.searchSymbol(memberType.name)
		if err != nil {
			return err, nil
		}

		if symbol.simbolType.kind != TYPE_STRUCT {
			return fmt.Errorf("Can only access field of struct"), nil
		}

		for _, child := range symbol.children {
			if child.name == node.token.tokenValue {
				return nil, &child.simbolType
			}
		}

		return fmt.Errorf("member does not exist in struct"), nil
	}

	return fmt.Errorf("Can't check type"), nil
}

func (this *Checker) enterScope() {
	this.currentScope ++
}

func (this *Checker) leaveScope() {
	delete(this.symbolTable, this.currentScope)

	this.currentScope --
}

func (this *Checker) walk(node *Node) error {
	for node != nil {
		if node.nodeType == NODE_STRUCT {
			structName := node.token.tokenValue

			this.enterScope()
			err := this.walk(node.right)
			if err != nil {
				return err
			}

			children := this.symbolTable[this.currentScope]
			this.leaveScope()

			err = this.addStructSymbol(structName, node, children)
			if err != nil {
				return err
			}
		} else if node.nodeType == NODE_VARIABLE_DECLARATION {
			var initializationSymbolType string
			if node.right != nil {
				err, initializationSymbol := this.determineType(node.right)
				if err != nil {
					return err
				}

				initializationSymbolType = initializationSymbol.name
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

			if initializationSymbolType != "" && variableSymbolType != initializationSymbolType {
				return fmt.Errorf("can't initialize with different types")
			}

			err := this.addVariableSymbol(node.token.tokenValue, variableSymbolType, node)
			if err != nil {
				return err
			}
		} else if node.nodeType == NODE_FUNCTION {
			symbolName := node.left.token.tokenValue
			err, symbolType := this.getTypeFromNode(node.left.left.left)
			if err != nil {
				return err
			}

			var signature []string
			for parameter := node.left.right; parameter != nil; parameter = parameter.next {
				err, parameterType := this.getTypeFromNode(parameter.left)
				if err != nil {
					return err
				}

				signature = append(signature, parameterType)
			}

			err, symbol := this.addFunctionSymbol(symbolName, symbolType, signature, node)
			if err != nil {
				return err
			}

			this.functionStack.push(symbol)
			this.enterScope()

			err = this.walk(node.left.right)
			if err != nil {
				return err
			}

			err = this.walk(node.right)
			if err != nil {
				return err
			}

			this.leaveScope()
			this.functionStack.pop()

		} else if node.nodeType == NODE_IF {
			err, symbolType := this.determineType(node.left)
			if err != nil {
				return err
			}

			if symbolType.name != "bool" {
				return fmt.Errorf("Can't have non-bool in if")
			}

			branchNode := node.right

			this.enterScope()

			err = this.walk(branchNode.left)
			if err != nil {
				return err
			}
			this.leaveScope()

			this.enterScope()
			err = this.walk(branchNode.right)
			if err != nil {
				return err
			}
			this.leaveScope()
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
			err, symbolType := this.determineType(node.left)
			if err != nil {
				return err
			}

			currentFunction := this.functionStack.peek()
			if currentFunction == nil {
				return fmt.Errorf("Return can only be inside a function")
			}

			if currentFunction.simbolType.name != symbolType.name {
				return fmt.Errorf("Invalid return type")
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