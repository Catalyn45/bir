package main

import "fmt"

const (
	NODE_PROGRAM     = iota
	NODE_ADD         = iota
	NODE_SUBSTRACT   = iota
	NODE_MULTIPLY    = iota
	NODE_DIVIDE      = iota
	NODE_EXPRESSION  = iota
	NODE_ASSIGNMENT  = iota
	NODE_PATH        = iota
	NODE_DECLARATION = iota
	NODE_PARAMETER   = iota
	NODE_STATEMENT   = iota
	NODE_RETURN      = iota
	NODE_IF          = iota
	NODE_BRANCH      = iota
	NODE_VARIABLE    = iota
	NODE_WHILE       = iota
	NODE_FOR         = iota
	NODE_METADATA    = iota
	NODE_MEMBER      = iota
	NODE_MODULE      = iota
	NODE_CONST       = iota
	NODE_ITERATION   = iota
	NODE_STRUCT      = iota
	NODE_FUNCTION_DECLARATION      = iota
	NODE_FUNCTION = iota
	NODE_IMPLEMENT = iota
	NODE_INTERFACE = iota
	NODE_IMPORT = iota
	NODE_TYPE_BOOL = iota
	NODE_TYPE_INT = iota
	NODE_TYPE_FLOAT = iota
	NODE_TYPE_STRING = iota
	NODE_TYPE_CUSTOM = iota
)

var nodeStrings = []string {
	"NODE_PROGRAM",
	"NODE_ADD",
	"NODE_SUBSTRACT",
	"NODE_MULTIPLY",
	"NODE_DIVIDE",
	"NODE_EXPRESSION",
	"NODE_ASSIGNMENT",
	"NODE_PATH",
	"NODE_DECLARATION",
	"NODE_PARAMETER",
	"NODE_STATEMENT",
	"NODE_RETURN",
	"NODE_IF",
	"NODE_BRANCH",
	"NODE_VARIABLE",
	"NODE_WHILE",
	"NODE_FOR",
	"NODE_METADATA",
	"NODE_MEMBER",
	"NODE_MODULE",
	"NODE_CONST",
	"NODE_ITERATION",
	"NODE_STRUCT",
	"NODE_FUNCTION_DECLARATION",
	"NODE_FUNCTION",
	"NODE_IMPLEMENT",
	"NODE_INTERFACE",
	"NODE_IMPORT",
	"NODE_TYPE_BOOL",
	"NODE_TYPE_INT",
	"NODE_TYPE_FLOAT",
	"NODE_TYPE_STRING",
	"NODE_TYPE_CUSTOM",
}

type Node struct {
	nodeType int
	token    *Token
	left     *Node
	right    *Node
	next     *Node
}

func (this *Node) ToString() string {
	tokenString := ""
	if this.token != nil {
		tokenString = this.token.toString()
	}

	return fmt.Sprintf("node: %s, token: %s", nodeStrings[this.nodeType], tokenString)
}

func (this *Node) Dump(indent int) {
	for i := 0; i < indent - 4; i++ {
		fmt.Print(" ")
	}

    if indent > 0 {
        fmt.Print("└── ")
    }

	fmt.Println(this.ToString())

	if this.left != nil {
		this.left.Dump(indent + 4)
	}

	if this.right != nil {
		this.right.Dump(indent + 4)
	}

	if this.next != nil {
		this.next.Dump(indent)
	}
}

type Parser struct {
	lexer        *Lexer
	currentToken *Token
}

func newParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer:        lexer,
		currentToken: nil,
	}
}

func (this *Parser) invalidTokenError(expectedTokenType int) error {
	err := fmt.Errorf("Invalid token: %s, expected: %s", this.currentToken.toString(), tokenTypesString[expectedTokenType])

	panic(err)
}

func  (this *Parser) unexpectedTokenError() error {
	err := fmt.Errorf("Unexpected token: %s", this.currentToken.toString())
	panic(err)
}

func (this *Parser) expectToken(expectedTokenType int) error {
	if this.currentToken.tokenType != expectedTokenType {
		return this.invalidTokenError(expectedTokenType)
	}

	return nil
}

func (this *Parser) eat(expectedTokenType int) error {
	err := this.expectToken(expectedTokenType)
	if err != nil {
		return err
	}

	this.advance()

	return nil
}

func (this *Parser) advance() error {
	err, token := this.lexer.next()
	if err != nil {
		return err
	}

	this.currentToken = token
	return nil
}

func (this *Parser) parseExpression() (error, *Node) {
	return fmt.Errorf("Expression parsing not implemented yet"), nil
}

func (this *Parser) parseExpressionStatement() (error, *Node) {
	err, expression := this.parseExpression()
	if err != nil {
		return err, nil
	}

	if this.currentToken.tokenType != TOKEN_ASSIGN {
		return nil, expression
	}

	this.advance()

	// TODO: Check if expression is lvalue
	err, rightExpression := this.parseExpression()
	if err != nil {
		return err, nil
	}

	assignmentNode := &Node {
		nodeType: NODE_ASSIGNMENT,
		left: expression,
		right: rightExpression,
	}

	return nil, assignmentNode
}

func (this *Parser) parseTypeSpecification() (error, *Node) {
	err := this.eat(TOKEN_COLONS)
	if err != nil {
		return err, nil
	}

	var node *Node
	if this.currentToken.tokenType == TOKEN_TYPE_BOOL {
		node = &Node { nodeType: NODE_TYPE_BOOL }
	} else if this.currentToken.tokenType == TOKEN_TYPE_INT {
		node = &Node { nodeType: NODE_TYPE_INT }
	} else if this.currentToken.tokenType == TOKEN_TYPE_FLOAT {
		node = &Node { nodeType: NODE_TYPE_FLOAT }
	} else if this.currentToken.tokenType == TOKEN_TYPE_STRING {
		node = &Node { nodeType: NODE_TYPE_STRING }
	} else if this.currentToken.tokenType == TOKEN_IDENTIFIER {
		node = &Node { nodeType: NODE_TYPE_CUSTOM }
	} else {
		return this.unexpectedTokenError(), nil
	}

	this.advance()

	return nil, node
}

func (this *Parser) parseTypedIdentifier() (error, *Node) {
	err := this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	variableNode := &Node {
		nodeType: NODE_VARIABLE,
		token: this.currentToken,
	}

	this.advance()

	err, typeNode := this.parseTypeSpecification()
	if err != nil {
		return err, nil
	}

	variableNode.left = typeNode

	return nil, variableNode
}

func (this *Parser) parseVariableDeclaration() (error, *Node) {
	err := this.eat(TOKEN_TYPE_VAR)
	if err != nil {
		return err, nil
	}

	err = this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	variableNode := &Node{
		nodeType: NODE_VARIABLE,
		token: this.currentToken,
	}

	this.advance()

	var variableTypeNode *Node = nil
	if this.currentToken.tokenType == TOKEN_COLONS {
		var err error
		err, variableTypeNode = this.parseTypeSpecification()
		if err != nil {
			return err, nil
		}
	}

	var expressionNode *Node = nil
	if this.currentToken.tokenType == TOKEN_EQUAL {
		this.advance()

		var err error
		err, expressionNode = this.parseExpression()
		if err != nil {
			return err, nil
		}
	}

	if variableTypeNode == nil && expressionNode == nil {
		return fmt.Errorf("Variable needs to be either typed or initialized"), nil
	}

	variableNode.left = variableTypeNode
	variableNode.right = expressionNode

	return nil, variableNode
}

func (this *Parser) parseConstant() (error, *Node) {
	err := this.eat(TOKEN_CONST)
	if err != nil {
		return err, nil
	}

	err = this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	constNode := &Node{
		nodeType: NODE_CONST,
		token: this.currentToken,
	}

	this.advance()

	var constTypeNode *Node = nil
	if this.currentToken.tokenType == TOKEN_COLONS {
		var err error
		err, constTypeNode = this.parseTypeSpecification()
		if err != nil {
			return err, nil
		}
	}

	err = this.eat(TOKEN_EQUAL)
	if err != nil {
		return err, nil
	}

	err, expressionNode := this.parseExpression()
	if err != nil {
		return err, nil
	}

	constNode.left = constTypeNode
	constNode.right = expressionNode

	return nil, constNode
}

func (this *Parser) parseIf() (error, *Node) {
	err := this.eat(TOKEN_IF)
	if err != nil {
		return err, nil
	}

	err, expressionNode := this.parseExpression()
	if err != nil {
		return err, nil
	}

	err, statementsNode := this.parseBlock()
	if err != nil {
		return err, nil
	}

	branchNode := Node {
		nodeType: NODE_BRANCH,
		left: statementsNode,
	}

	if this.currentToken.tokenType == TOKEN_ELSE {
		this.advance()

		var elseNode *Node
		var err error

		if this.currentToken.tokenType == TOKEN_IF {
			err, elseNode = this.parseIf()
		} else {
			err, elseNode = this.parseBlock()
		}

		if err != nil {
			return err, nil
		}

		branchNode.right = elseNode
	}

	ifNode := &Node {
		nodeType: NODE_IF,
		left: expressionNode,
		right: &branchNode,
	}

	return nil, ifNode
}

func (this *Parser) parseWhile() (error, *Node) {
	err := this.eat(TOKEN_WHILE)
	if err != nil {
		return err, nil
	}

	err, expressionNode := this.parseExpression()
	if err != nil {
		return err, nil
	}

	err, statementsNode := this.parseBlock()
	if err != nil {
		return err, nil
	}

	branchNode := Node {
		nodeType: NODE_BRANCH,
		left: statementsNode,
	}

	if this.currentToken.tokenType == TOKEN_ELSE {
		this.advance()

		err, elseNode := this.parseBlock()
		if err != nil {
			return err, nil
		}

		branchNode.right = elseNode
	}

	whileNode := &Node {
		nodeType: NODE_IF,
		left: expressionNode,
		right: &branchNode,
	}

	return nil, whileNode
}

func (this *Parser) parseFor() (error, *Node) {
	err := this.eat(TOKEN_FOR)
	if err != nil {
		return err, nil
	}

	err = this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	iteratorNode := &Node{
		nodeType: NODE_VARIABLE,
		token: this.currentToken,
	}

	this.advance()

	if this.currentToken.tokenType == TOKEN_COLONS {
		var err error
		err, iteratorNode.left = this.parseTypeSpecification()
		if err != nil {
			return err, nil
		}
	}

	err = this.eat(TOKEN_IN)
	if err != nil {
		return err, nil
	}

	err, iterableNode := this.parseExpression()
	if err != nil {
		return err, nil
	}
	this.advance()

	err, statementsNode := this.parseBlock()
	if err != nil {
		return err, nil
	}

	iterationNode := &Node {
		nodeType: NODE_ITERATION,
		left: iteratorNode,
		right: iterableNode,
	}

	forNode := &Node {
		nodeType: NODE_FOR,
		left: iterationNode,
		right: statementsNode,
	}

	return nil, forNode
}

func (this *Parser) parseReturn() (error, *Node) {
	err := this.eat(TOKEN_RETURN)
	if err != nil {
		return err, nil
	}

	err, expressionNode := this.parseExpression()
	if err != nil {
		return err, nil
	}

	returnNode := &Node {
		nodeType: NODE_RETURN,
		left: expressionNode,
	}

	return nil, returnNode
}

func (this *Parser) parseStatement() (error, *Node) {
	if this.currentToken.tokenType == TOKEN_CONST {
		return this.parseConstant()
	}

	if this.currentToken.tokenType == TOKEN_TYPE_VAR {
		return this.parseVariableDeclaration()
	}

	if this.currentToken.tokenType == TOKEN_IF {
		return this.parseIf()
	}

	if this.currentToken.tokenType == TOKEN_WHILE {
		return this.parseWhile()
	}

	if this.currentToken.tokenType == TOKEN_FOR {
		return this.parseFor()
	}

	if this.currentToken.tokenType == TOKEN_RETURN {
		return this.parseReturn()
	}

	return this.parseExpressionStatement()
}

func (this *Parser) parseBlock() (error, *Node) {
	err := this.eat(TOKEN_OPEN_BRACKET)
	if err != nil {
		return err, nil
	}

	var statementsNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_BRACKET; {
		err, node := this.parseStatement()
		if err != nil {
			return err, nil
		}

		if currentNode == nil {
			statementsNode = node
		} else {
			currentNode.next = node
		}

		currentNode = node
	}
	this.advance()

	return nil, statementsNode
}

func (this *Parser) parseMembers() (error, *Node) {
	err := this.eat(TOKEN_OPEN_BRACKET)
	if err != nil {
		return err, nil
	}

	var membersNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_BRACKET; {
		err, node := this.parseTypedIdentifier()
		if err != nil {
			return err, nil
		}

		if currentNode == nil {
			membersNode = node
		} else {
			currentNode.next = node
		}

		currentNode = node
	}
	this.advance()

	return nil, membersNode
}

func (this *Parser) parseStruct() (error, *Node) {
	err := this.eat(TOKEN_TYPE_STRUCT)
	if err != nil {
		return err, nil
	}

	err = this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	structNode := &Node {
		nodeType: NODE_STRUCT,
		token: this.currentToken,
	}

	this.advance()

	err, membersNode := this.parseMembers()
	if err != nil {
		return err, nil
	}

	structNode.left = membersNode

	return nil, structNode
}

func (this *Parser) parseFunctionParameters() (error, *Node) {
	err := this.eat(TOKEN_OPEN_PARANTHESIS)
	if err != nil {
		return err, nil
	}

	var parametersNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_PARANTHESIS; {
		err, node := this.parseTypedIdentifier()
		if err != nil {
			return err, nil
		}

		if currentNode == nil {
			parametersNode = node
		} else {
			currentNode.next = node
		}

		currentNode = node

		if this.currentToken.tokenType == TOKEN_COMMA {
			this.advance()
		}
	}
	this.advance()

	return nil, parametersNode
}

func (this *Parser) parseFunctionDeclaration() (error, *Node) {
	err := this.eat(TOKEN_FUNCTION)
	if err != nil {
		return err, nil
	}

	err = this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	functionDeclarationNode := &Node {
		nodeType: NODE_FUNCTION_DECLARATION,
		token: this.currentToken,
	}

	this.advance()

	err, parametersNode := this.parseFunctionParameters()
	if err != nil {
		return err, nil
	}

	var functionTypeNode *Node = nil
	if this.currentToken.tokenType == TOKEN_COLONS {
		err, functionTypeNode = this.parseTypeSpecification()
		if err != nil {
			return err, nil
		}
	}

	functionDeclarationNode.left = functionTypeNode
	functionDeclarationNode.right = parametersNode

	return nil, functionDeclarationNode
}

func (this *Parser) parseInterfaceFunctionDeclarations() (error, *Node) {
	err := this.eat(TOKEN_OPEN_BRACKET)
	if err != nil {
		return err, nil
	}

	var functionDeclarationsNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_BRACKET; {
		err, node := this.parseFunctionDeclaration()
		if err != nil {
			return err, nil
		}

		if currentNode == nil {
			functionDeclarationsNode = node
		} else {
			currentNode.next = node
		}

		currentNode = node
	}
	this.advance()

	return nil, functionDeclarationsNode
}

func (this *Parser) parseFunction() (error, *Node) {
	err, functionDeclarationNode := this.parseFunctionDeclaration()
	if err != nil {
		return err, nil
	}

	err, statementsNode := this.parseBlock()
	if err != nil {
		return err, nil
	}

	functionNode := &Node {
		nodeType: NODE_FUNCTION,
		left: functionDeclarationNode,
		right: statementsNode,
	}

	return nil, functionNode
}

func (this *Parser) parseImplementFunctionDefinitions() (error, *Node) {
	err := this.eat(TOKEN_OPEN_BRACKET)
	if err != nil {
		return err, nil
	}

	var functionsNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_BRACKET; {
		err, node := this.parseFunction()
		if err != nil {
			return err, nil
		}

		if currentNode == nil {
			functionsNode = node
		} else {
			currentNode.next = node
		}

		currentNode = node
	}
	this.advance()

	return nil, functionsNode
}

func (this *Parser) parseInterface() (error, *Node) {
	err := this.eat(TOKEN_TYPE_INTERFACE)
	if err != nil {
		return err, nil
	}

	err = this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	interfaceNode := &Node {
		nodeType: NODE_INTERFACE,
		token: this.currentToken,
	}

	this.advance()

	err, functionDeclarationsNode := this.parseInterfaceFunctionDeclarations()
	if err != nil {
		return err, nil
	}

	interfaceNode.left = functionDeclarationsNode

	return nil, interfaceNode
}

func (this *Parser) parseImplement() (error, *Node) {
	err := this.eat(TOKEN_IMPLEMENT)
	if err != nil {
		return err, nil
	}

	err = this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	implementNode := &Node {
		nodeType: NODE_IMPLEMENT,
		token: this.currentToken,
	}

	this.advance()

	err, functionDefinitionsNode := this.parseImplementFunctionDefinitions()
	if err != nil {
		return err, nil
	}

	implementNode.left = functionDefinitionsNode

	return nil, implementNode
}

func (this *Parser) parseDeclaration() (error, *Node) {
	if this.currentToken.tokenType == TOKEN_TYPE_STRUCT {
		return this.parseStruct()
	}

	if this.currentToken.tokenType == TOKEN_TYPE_INTERFACE {
		return this.parseInterface()
	}

	if this.currentToken.tokenType == TOKEN_IMPLEMENT {
		return this.parseImplement()
	}

	if this.currentToken.tokenType == TOKEN_FUNCTION {
		return this.parseFunction()
	}

	if this.currentToken.tokenType == TOKEN_CONST {
		return this.parseConstant()
	}

	return this.unexpectedTokenError(), nil
}

func (this *Parser) parseDeclarations() (error, *Node) {
	var declarationsNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_BRACKET; {
		err, node := this.parseDeclaration()
		if err != nil {
			return err, nil
		}

		if currentNode == nil {
			declarationsNode = node
		} else {
			currentNode.next = node
		}

		currentNode = node
	}
	this.advance()

	return nil, declarationsNode
}

func (this *Parser) parseExport() (error, *Node) {
	if this.currentToken.tokenType == TOKEN_TYPE_STRUCT {
		return this.parseStruct()
	}

	if this.currentToken.tokenType == TOKEN_TYPE_INTERFACE {
		return this.parseInterface()
	}

	if this.currentToken.tokenType == TOKEN_IMPLEMENT {
		return this.parseImplement()
	}

	if this.currentToken.tokenType == TOKEN_FUNCTION {
		return this.parseFunction()
	}

	if this.currentToken.tokenType == TOKEN_CONST {
		return this.parseConstant()
	}

	if this.currentToken.tokenType == TOKEN_OPEN_BRACKET {
		return this.parseDeclarations()
	}

	return this.unexpectedTokenError(), nil
}

func (this *Parser) parseRootStatement() (error, *Node) {
	if this.currentToken.tokenType == TOKEN_TYPE_STRUCT {
		return this.parseStruct()
	}

	if this.currentToken.tokenType == TOKEN_TYPE_INTERFACE {
		return this.parseInterface()
	}

	if this.currentToken.tokenType == TOKEN_IMPLEMENT {
		return this.parseImplement()
	}

	if this.currentToken.tokenType == TOKEN_FUNCTION {
		return this.parseFunction()
	}

	if this.currentToken.tokenType == TOKEN_CONST {
		return this.parseConstant()
	}

	if this.currentToken.tokenType == TOKEN_EXPORT {
		return this.parseExport()
	}

	return this.unexpectedTokenError(), nil
}

func (this *Parser) parsePath() (error, *Node) {
	var pathNode *Node = nil
	var currentNode *Node = nil
	for {
		err := this.expectToken(TOKEN_IDENTIFIER)
		if err != nil {
			return err, nil
		}

		node := &Node {
			nodeType: NODE_PATH,
			token: this.currentToken,
		}

		this.advance()

		if currentNode == nil {
			pathNode = node
		} else {
			currentNode.next = node
		}

		currentNode = node

		if this.currentToken.tokenType != TOKEN_DOT {
			break
		}

		this.advance()
	}

	return nil, pathNode
}

func (this *Parser) parseModule() (error, *Node) {
	err := this.eat(TOKEN_MODULE)
	if err != nil {
		return err, nil
	}

	err, pathNode := this.parsePath()
	if err != nil {
		return err, nil
	}

	moduleNode := &Node {
		nodeType: NODE_MODULE,
		left: pathNode,
	}

	return nil, moduleNode
}

func (this *Parser) parseImport() (error, *Node) {
	err := this.eat(TOKEN_IMPORT)
	if err != nil {
		return err, nil
	}

	err, pathNode := this.parsePath()
	if err != nil {
		return err, nil
	}

	importNode := &Node {
		nodeType: NODE_IMPORT,
		left: pathNode,
	}

	return nil, importNode
}

func (this *Parser) parseImports() (error, *Node) {
	var importsNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType == TOKEN_IMPORT; {
		err, node := this.parseImport()
		if err != nil {
			return err, nil
		}

		if currentNode == nil {
			importsNode = node
		} else {
			currentNode.next = node
		}

		currentNode = node
	}

	return nil, importsNode
}

func (this *Parser) Parse() (error, *Node) {
	err := this.advance()
	if err != nil {
		return err, nil
	}

	err, moduleNode := this.parseModule()
	if err != nil {
		return err, nil
	}

	err, importsNode := this.parseImports()
	if err != nil {
		return err, nil
	}

	var statementsNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_EOF; {
		err, node := this.parseRootStatement()
		if err != nil {
			return err, nil
		}

		if currentNode == nil {
			statementsNode = node
		} else {
			currentNode.next = node
		}

		currentNode = node
	}

	programMetadataNode := &Node {
		nodeType: NODE_METADATA,
		left: moduleNode,
		right: importsNode,
	}

	root := &Node{
		nodeType: NODE_PROGRAM,
		left: programMetadataNode,
		right: statementsNode,
	}

	return nil, root
}
