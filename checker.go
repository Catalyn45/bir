package main

import "fmt"

const (
	TYPE_CONST    = iota
	TYPE_LITERAL  = iota
	TYPE_VARIABLE = iota
	TYPE_FUNCTION = iota
	TYPE_STRUCT   = iota
	TYPE_INTERFACE = iota
)

type Signature struct {
	parameterTypes[] string
	returnType string
}

type SymbolType struct {
	kind int
	name string
	signature *Signature
}

type Symbol struct {
	name       string
	simbolType SymbolType
	node 	   *Node
}

type SymbolTable map[string]*Symbol

type Checker struct {
	modules map[string][]*Node
	symbolTables  *Stack[*SymbolTable]
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
		symbolTables:  &Stack[*SymbolTable]{},
		functionStack: &Stack[*Symbol]{},
	}
}

func (this *Checker) symbolAlreadyExists(symbolName string) bool {
	lastScope := *this.symbolTables.peek()
	
	_, ok := lastScope[symbolName]

	return ok
}

func (this *Checker) addVariableSymbol(varName string, varType string, node *Node) error {
	if this.symbolAlreadyExists(varName) {
		return fmt.Errorf("variable already declared in current scope")
	}

	lastScope := *this.symbolTables.peek()

	lastScope[varName] = &Symbol {
		name: varName,
		simbolType: SymbolType {
			kind: TYPE_VARIABLE,
			name: varType,
		},
		node: node,
	}

	node.symbol = lastScope[varName]

	return nil
}

func (this *Checker) addFunctionSymbol(functionName string, returnType string, parametersTypes []string, node *Node) (error, *Symbol) {
	if this.symbolAlreadyExists(functionName) {
		return fmt.Errorf("function already declared in current scope"), nil
	}

	lastScope := *this.symbolTables.peek()

	lastScope[functionName] = &Symbol {
		name: functionName,
		simbolType: SymbolType {
			kind: TYPE_FUNCTION,
			name: functionName,
			signature: &Signature {
				parameterTypes: parametersTypes,
				returnType: returnType,
			},
		},
		node: node,
	}

	node.symbol = lastScope[functionName]

	return nil, lastScope[functionName]
}

func (this *Checker) addStructSymbol(structName string, node *Node) error {
	lastScope := *this.symbolTables.peek()

	lastScope[structName] = &Symbol {
		name: structName,
		simbolType: SymbolType {
			kind: TYPE_STRUCT,
			name: structName,
		},
		node: node,
	}

	node.symbol = lastScope[structName]

	return nil
}

func (this *Checker) addInterfaceSymbol(interfaceName string, node *Node) error {
	lastScope := *this.symbolTables.peek()

	lastScope[interfaceName] = &Symbol {
		name: interfaceName,
		simbolType: SymbolType {
			kind: TYPE_INTERFACE,
			name: interfaceName,
		},
		node: node,
	}

	node.symbol = lastScope[interfaceName]

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
	var foundSymbol *Symbol = nil

	this.symbolTables.foreach(func (item *SymbolTable) (stop bool) {
		val, ok := (*item)[symbolName]
		if !ok {
			return false
		}

		foundSymbol = val

		return true
	})

	if foundSymbol == nil {
		return fmt.Errorf("Symbol not declarated"), nil
	}

	return nil, foundSymbol
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

		if len(parameterTypes.parameterTypes) != len(argumentTypes) {
			return fmt.Errorf("Not the same number of arguments"), nil
		}

		for i := 0; i < len(parameterTypes.parameterTypes); i++ {
			if !this.isAssignable(parameterTypes.parameterTypes[i], argumentTypes[i]) {
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

		if symbol.simbolType.kind != TYPE_STRUCT && symbol.simbolType.kind != TYPE_INTERFACE {
			return fmt.Errorf("Can only access field of struct or interface"), nil
		}

		symbol, ok := (*symbol.node.symbolTable)[node.token.tokenValue]
		if !ok {
			return fmt.Errorf("member does not exist in struct or interface"), nil
		}

		return nil, &symbol.simbolType
	}

	if node.nodeType == NODE_VARIABLE_DECLARATION {
		var initializationSymbolType string
		if node.right != nil {
			err, initializationSymbol := this.determineType(node.right)
			if err != nil {
				return err, nil
			}

			initializationSymbolType = initializationSymbol.name
		}

		var variableSymbolType string
		if node.left != nil {
			err, symbolType := this.getTypeFromNode(node.left)
			if err != nil {
				return err, nil
			}

			variableSymbolType = symbolType
		}

		if variableSymbolType == "" {
			variableSymbolType = initializationSymbolType
		}

		if initializationSymbolType != "" && !this.isAssignable(variableSymbolType, initializationSymbolType) {
			return fmt.Errorf("can't initialize with different types"), nil
		}

		err := this.addVariableSymbol(node.token.tokenValue, variableSymbolType, node)
		if err != nil {
			return err, nil
		}

		err, symbol := this.searchSymbol(node.token.tokenValue)
		
		if err != nil {
			return err, nil
		}

		return nil, &symbol.simbolType
	}

	return fmt.Errorf("Can't check type"), nil
}

func (this *Checker) enterScope(node *Node) {
	symbolTable := make(SymbolTable)
	node.symbolTable = &symbolTable
	this.symbolTables.push(&symbolTable)
}

func (this *Checker) leaveScope() {
	this.symbolTables.pop()
}

func (this *Checker) implementsInterface(leftSymbol *Symbol, rightSymbol *Symbol) bool {
	for _, member := range (*leftSymbol.node.symbolTable) {
		value, ok := (*rightSymbol.node.symbolTable)[member.name]
		if !ok {
			return false
		}

		if value.simbolType.kind != TYPE_FUNCTION {
			continue
		}

		if value.simbolType.signature.returnType != member.simbolType.signature.returnType {
			return false
		}

		if len(value.simbolType.signature.parameterTypes) != len(member.simbolType.signature.parameterTypes) {
			return false
		}

		for index, parameterType := range value.simbolType.signature.parameterTypes {
			if parameterType != member.simbolType.signature.parameterTypes[index] {
				return false
			}
		}
	}

	return true
}

func (this *Checker) isAssignable(leftType string, rightType string) bool {
	if leftType == rightType {
		return true
	}

	err, leftSymbol := this.searchSymbol(leftType)
	if err != nil {
		return false
	}

	if leftSymbol.simbolType.kind != TYPE_INTERFACE {
		return false
	}

	err, rightSymbol := this.searchSymbol(rightType)
	if err != nil {
		return false
	}

	if rightSymbol.simbolType.kind != TYPE_STRUCT {
		return false
	}

	return this.implementsInterface(leftSymbol, rightSymbol)
}

func (this *Checker) addFunctionDeclaration(node *Node) (error, *Symbol) {
	symbolName := node.token.tokenValue

	var symbolType string
	if node.left.left == nil {
		symbolType = "void"
	} else {
		var err error
		err, symbolType = this.getTypeFromNode(node.left.left)
		if err != nil {
			return err, nil
		}
	}

	var signature []string
	for parameter := node.right; parameter != nil; parameter = parameter.next {
		err, parameterType := this.getTypeFromNode(parameter.left)
		if err != nil {
			return err, nil
		}

		signature = append(signature, parameterType)
	}

	return this.addFunctionSymbol(symbolName, symbolType, signature, node)
}

func (this *Checker) walkStatements(node *Node) error {
	err, _ := this.walkGetLastStatement(node)
	
	return err
}

func (this *Checker) walkGetLastStatement(node *Node) (error, *Node) {
	var lastNode *Node = nil
	for node != nil {
		if node.nodeType == NODE_IF || node.nodeType == NODE_WHILE {
			this.enterScope(node)

			err, symbolType := this.determineType(node.left)
			if err != nil {
				return err, nil
			}

			if symbolType.name != "bool" {
				return fmt.Errorf("Can't have non-bool in if"), nil
			}

			branchNode := node.right

			if branchNode.left != nil {
				this.enterScope(branchNode.left)

				err = this.walkStatements(branchNode.left)
				if err != nil {
					return err, nil
				}

				this.leaveScope()
			}

			if branchNode.right != nil {
				this.enterScope(branchNode.right)

				err = this.walkStatements(branchNode.right)
				if err != nil {
					return err, nil
				}

				this.leaveScope()
			}

			this.leaveScope()
		} else if node.nodeType == NODE_ASSIGNMENT {
			err, leftSymbolType := this.determineType(node.left)
			if err != nil {
				return err, nil
			}

			err, rightSymbolType := this.determineType(node.right)
			if err != nil {
				return err, nil
			}

			if !this.isAssignable(leftSymbolType.name, rightSymbolType.name) {
				return fmt.Errorf("Can't assign different types"), nil
			}
		} else if node.nodeType == NODE_RETURN {
			err, symbolType := this.determineType(node.left)
			if err != nil {
				return err, nil
			}

			currentFunction := this.functionStack.peek()
			if currentFunction == nil {
				return fmt.Errorf("Return can only be inside a function"), nil
			}

			if currentFunction.simbolType.signature.returnType != symbolType.name {
				return fmt.Errorf("Invalid return type"), nil
			}
		} else {
			err, _ := this.determineType(node)
			if err != nil {
				return err, nil
			}
		}

		lastNode = node
		node = node.next
	}

	return nil, lastNode
}

func (this *Checker) addSymbolHeader (value string, typeType int, node *Node) error {
	if this.symbolAlreadyExists(value) {
		return fmt.Errorf("struct already declared in current scope")
	}

	lastScope := *this.symbolTables.peek()

	lastScope[value] = &Symbol {
		name: value,
		simbolType: SymbolType {
			kind: typeType,
		},
		node: node,
	}

	node.symbol = lastScope[value]

	return nil
}

func (this *Checker) walkRootTypes (node *Node) error {
	for node != nil {
		if node.nodeType == NODE_STRUCT {
			err := this.addSymbolHeader(node.token.tokenValue, TYPE_STRUCT, node)
			if err != nil {
				return err
			}

			// create the symbol table early in case implement statement appear before struct declaration statement
			this.enterScope(node)
			this.leaveScope()
		} else if node.nodeType == NODE_INTERFACE {
			err := this.addSymbolHeader(node.token.tokenValue, TYPE_INTERFACE, node)
			if err != nil {
				return err
			}
		}

		node = node.next
	}

	return nil
}

func (this *Checker) walkRootDeclarations (node *Node) error {
	for node != nil {
		if node.nodeType == NODE_STRUCT {
			structName := node.token.tokenValue

			this.symbolTables.push(node.symbolTable)
			
			err := this.walkStatements(node.right)
			if err != nil {
				return err
			}

			this.symbolTables.pop()

			err = this.addStructSymbol(structName, node)
			if err != nil {
				return err
			}
		} else if node.nodeType == NODE_IMPLEMENT {
			structName := node.token.tokenValue

			err, symbol := this.searchSymbol(structName)
			if err != nil {
				return err
			}

			if symbol.simbolType.kind != TYPE_STRUCT {
				return fmt.Errorf("Only structs can be implemented")
			}

			// push the struct symbol table
			this.symbolTables.push(symbol.node.symbolTable)

			err = this.walkRootDeclarations(node.right)
			if err != nil {
				return err
			}

			this.symbolTables.pop()
		} else if node.nodeType == NODE_INTERFACE {
			interfaceName := node.token.tokenValue

			this.enterScope(node)

			err := this.walkRootDeclarations(node.right)
			if err != nil {
				return err
			}

			this.leaveScope()

			err = this.addInterfaceSymbol(interfaceName, node)
			if err != nil {
				return err
			}
		} else if node.nodeType == NODE_FUNCTION {
			err, _ := this.addFunctionDeclaration(node.left)
			if err != nil {
				return err
			}
		} else if node.nodeType == NODE_FUNCTION_DECLARATION {
			err, _ := this.addFunctionDeclaration(node)
			if err != nil {
				return err
			}
		}

		node = node.next
	}

	return nil
}

func (this *Checker) walk (node *Node) error {
	for node != nil {
		if node.nodeType == NODE_IMPLEMENT {
			err := this.walk(node.right)
			if err != nil {
				return err
			}
		} else if node.nodeType == NODE_FUNCTION {
			symbol := node.left.symbol

			this.functionStack.push(symbol)
			this.enterScope(node)

			err := this.walkStatements(node.left.right)
			if err != nil {
				return err
			}

			err, lastStatement := this.walkGetLastStatement(node.right)
			if err != nil {
				return err
			}

			if symbol.simbolType.signature.returnType != "void" && (lastStatement == nil || lastStatement.nodeType != NODE_RETURN) {
				return fmt.Errorf("function needs to end with return")
			}

			this.leaveScope()
			this.functionStack.pop()
		} else {
			this.walkStatements(node)
		}

		node = node.next
	}

	return nil
}

func (this *Checker) walkRoot(node *Node) error {
	err := this.walkRootTypes(node)
	if err != nil {
		return err
	}

	return this.walkRootDeclarations(node)
}

func (this *Checker) Check() error {
	for _, asts := range this.modules {
		symbolTable := make(SymbolTable)
		this.symbolTables.push(&symbolTable)

		for _, ast := range asts {
			ast.symbolTable = &symbolTable

			err := this.walkRoot(ast.right)
			if err != nil {
				return err
			}
		}

		for _, ast := range asts {
			err := this.walk(ast.right)
			if err != nil {
				return err
			}
		}

		this.symbolTables.pop()
	}

	return nil
}