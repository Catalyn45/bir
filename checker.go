package main

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

const (
	TYPE_MODULE   = iota
	TYPE_CONST    = iota
	TYPE_LITERAL  = iota
	TYPE_VARIABLE = iota
	TYPE_FUNCTION = iota
	TYPE_STRUCT   = iota
	TYPE_INTERFACE = iota
	TYPE_EXPRESSION = iota
)

type Parameter struct {
	name string
	paramType *SymbolType
	node *Node
}

type Signature struct {
	parameters[] *Parameter
	returnType *SymbolType
}

type SymbolType struct {
	kind int
	name string
	symbol *Symbol
	signature *Signature
}

type Symbol struct {
	name       string
	simbolType SymbolType
	node 	   *Node
	value      value.Value
	valueType  types.Type
}

type SymbolTable map[string]*Symbol

type Checker struct {
	modules map[string][]*Node
	symbolTables  *Stack[*SymbolTable]
	functionStack *Stack[*Symbol]
	imports []string
	currentStruct *SymbolType
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

func (this *Checker) addVariableSymbol(varName string, varType *SymbolType, node *Node) error {
	if this.symbolAlreadyExists(varName) {
		return fmt.Errorf("variable already declared in current scope")
	}

	lastScope := *this.symbolTables.peek()

	lastScope[varName] = &Symbol {
		name: varName,
		simbolType: *varType,
		node: node,
	}

	node.symbol = lastScope[varName]

	return nil
}

func (this *Checker) addFunctionSymbol(functionName string, returnType *SymbolType, parametersTypes []*Parameter, node *Node) (error, *Symbol) {
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
				parameters: parametersTypes,
				returnType: returnType,
			},
		},
		node: node,
	}

	node.symbol = lastScope[functionName]
	lastScope[functionName].simbolType.symbol = node.symbol

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
	lastScope[structName].simbolType.symbol = node.symbol

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
	lastScope[interfaceName].simbolType.symbol = node.symbol

	return nil
}

func (this *Checker) searchSymbolInModule(module string, symbolName string) (error, *Symbol) {
	for _, node := range this.modules[module] {
		value, ok := (*node.symbolTable)[symbolName]
		if ok {
			return nil, value
		}
	}

	return fmt.Errorf("couldn't find symbol in module"), nil
}

func (this *Checker) searchSymbol(symbolName string) (error, *Symbol) {
	// for _, imp := range this.imports {
	// 	if imp == symbolName {
	// 		return nil, &SymbolType {
	// 			kind: TYPE_MODULE,
	// 			name: imp,
	// 		}
	// 	}
	// }

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

func (this *Checker) getTypeFromNode(node *Node) (error, *SymbolType) {
	if node.nodeType == NODE_INT_TYPE {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "int"}
	}

	if node.nodeType == NODE_FLOAT_TYPE {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "float"}
	}

	if node.nodeType == NODE_STRING_TYPE {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "string"}
	}

	if node.nodeType == NODE_BOOL_TYPE {
		return nil, &SymbolType {kind: TYPE_LITERAL, name: "bool"}
	}

	if node.nodeType == NODE_CUSTOM_TYPE {
		err, symbol := this.searchSymbol(node.token.tokenValue)
		if err != nil {
			return err, nil
		}

		return nil, &symbol.simbolType
	}

	return fmt.Errorf("Invalid node"), nil
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
		node.token.tokenType == TOKEN_GREATER_EQUAL ||
		node.token.tokenType == TOKEN_PLUS ||
		node.token.tokenType == TOKEN_MINUS ||
		node.token.tokenType == TOKEN_MULTIPLY ||
		node.token.tokenType == TOKEN_DIVIDE {
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
		err, symbol := this.searchSymbol(node.token.tokenValue)
		if err != nil {
			return err, nil
		}

		node.symbol = symbol

		return nil, &symbol.simbolType
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

		var argumentTypes []*SymbolType
		for argument := node.right.right; argument != nil; argument = argument.next {
			err, argumentType := this.determineType(argument)
			if err != nil {
				return err, nil
			}

			argumentTypes = append(argumentTypes, argumentType)
		}

		toSubstract := 0
		if len(parameterTypes.parameters) > 0 && parameterTypes.parameters[0].name == "this" {
			toSubstract = 1
		}

		if (len(parameterTypes.parameters) - toSubstract) != len(argumentTypes) {
			return fmt.Errorf("Not the same number of arguments: %d, %d", len(parameterTypes.parameters), len(argumentTypes)), nil
		}

		for i := toSubstract; i < len(parameterTypes.parameters); i++ {
			if !this.isAssignable(parameterTypes.parameters[i].paramType, argumentTypes[i]) {
				return fmt.Errorf("Invalid argument type for parameter"), nil
			}
		}

		return nil, symbolType.signature.returnType
	}

	if node.nodeType == NODE_MEMBER_ACCESS {
		err, memberType := this.determineType(node.left)
		if err != nil {
			return err, nil
		}

		var symbol *Symbol
		if memberType.kind == TYPE_MODULE {
			err, symbol = this.searchSymbolInModule(memberType.name, node.token.tokenValue)
			if err != nil {
				return err, nil
			}

			return nil, &symbol.simbolType
		} else {
			err, symbol = this.searchSymbol(memberType.name)
			if err != nil {
				return err, nil
			}
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
		var initializationSymbolType *SymbolType
		if node.right != nil {
			err, initializationSymbol := this.determineType(node.right)
			if err != nil {
				return err, nil
			}

			initializationSymbolType = initializationSymbol
		}

		var variableSymbolType *SymbolType = nil
		if node.left != nil {
			err, symbolType := this.getTypeFromNode(node.left)
			if err != nil {
				return err, nil
			}

			variableSymbolType = symbolType
		}

		if variableSymbolType == nil {
			variableSymbolType = initializationSymbolType
		}

		if initializationSymbolType != nil && !this.isAssignable(variableSymbolType, initializationSymbolType) {
			return fmt.Errorf("can't initialize with different types"), nil
		}

		err := this.addVariableSymbol(node.token.tokenValue, variableSymbolType, node)
		if err != nil {
			return err, nil
		}

		return nil, &node.symbol.simbolType
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

		if len(value.simbolType.signature.parameters) != len(member.simbolType.signature.parameters) {
			return false
		}

		for index, parameterType := range value.simbolType.signature.parameters {
			if parameterType != member.simbolType.signature.parameters[index] {
				return false
			}
		}
	}

	return true
}

func (this *Checker) isAssignable(leftSymbolType *SymbolType, rightSymbolType *SymbolType) bool {
	if leftSymbolType.name == rightSymbolType.name {
		return true
	}

	if leftSymbolType.kind != TYPE_INTERFACE {
		return false
	}

	if leftSymbolType.kind != TYPE_STRUCT {
		return false
	}

	return this.implementsInterface(leftSymbolType.symbol, rightSymbolType.symbol)
}

func (this *Checker) addFunctionDeclaration(node *Node) (error, *Symbol) {
	symbolName := node.token.tokenValue

	var symbolType *SymbolType
	if node.left.left == nil {
		symbolType = &SymbolType {
			kind: TYPE_LITERAL,
			name: "void",
		}
	} else {
		var err error
		err, symbolType = this.getTypeFromNode(node.left.left)
		if err != nil {
			return err, nil
		}
	}

	var signature []*Parameter

	if this.currentStruct != nil {
		// add this
		signature = append(signature, &Parameter{
			name: "this",
			paramType: this.currentStruct,
			node: &Node{},
		})
	}

	for parameter := node.right; parameter != nil; parameter = parameter.next {
		err, parameterType := this.getTypeFromNode(parameter.left)
		if err != nil {
			return err, nil
		}

		signature = append(signature, &Parameter {
				name: parameter.token.tokenValue,
				paramType: parameterType,
				node: parameter,
			},
		)
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

			if !this.isAssignable(leftSymbolType, rightSymbolType) {
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

			if currentFunction.simbolType.signature.returnType.name != symbolType.name {
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

			this.currentStruct = &symbol.simbolType

			err = this.walkRootDeclarations(node.right)
			if err != nil {
				return err
			}

			this.currentStruct = nil

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
		} else if node.nodeType == NODE_FUNCTION || node.nodeType == NODE_CONSTRUCTOR {
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

func (this *Checker) walk(node *Node) error {
	for node != nil {
		if node.nodeType == NODE_IMPLEMENT {
			err, symbol := this.searchSymbol(node.token.tokenValue)
			if err != nil {
				return err
			}

			this.currentStruct = &symbol.simbolType

			err = this.walk(node.right)
			if err != nil {
				return err
			}

			this.currentStruct = nil
		} else if node.nodeType == NODE_FUNCTION || node.nodeType == NODE_CONSTRUCTOR {
			symbol := node.left.symbol

			this.functionStack.push(symbol)
			this.enterScope(node.left)

			// declare this
			if this.currentStruct != nil {
				currentFunction := this.functionStack.peek()

				err := this.addVariableSymbol("this", this.currentStruct, currentFunction.simbolType.signature.parameters[0].node)
				if err != nil {
					return err
				}
			}

			err := this.walkStatements(node.left.right)
			if err != nil {
				return err
			}

			err, lastStatement := this.walkGetLastStatement(node.right)
			if err != nil {
				return err
			}

			if symbol.simbolType.signature.returnType.name != "void" && (lastStatement == nil || lastStatement.nodeType != NODE_RETURN) {
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

func (this *Checker) walkImports(node *Node) error {
	var imports []string = nil

	for imp := node.right; imp != nil; imp = imp.next {
		imports = append(imports, imp.left.token.tokenValue)
	}

	node.imports = imports

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

		this.symbolTables.pop()
	}

	for _, asts := range this.modules {
		for _, ast := range asts {
			this.symbolTables.push(ast.symbolTable)

			err := this.walkImports(ast.left)
			if err != nil {
				return err
			}

			this.imports = ast.left.imports

			err = this.walk(ast.right)
			if err != nil {
				return err
			}

			this.imports = nil

			this.symbolTables.pop()
		}
	}

	return nil
}