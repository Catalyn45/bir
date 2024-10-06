package main

import "fmt"

const (
	NODE_PROGRAM     = iota
	NODE_LINK        = iota
	NODE_ADD         = iota
	NODE_SUBSTRACT   = iota
	NODE_MULTIPLY    = iota
	NODE_DIVIDE      = iota
	NODE_EXPRESSION  = iota
	NODE_ASSIGNMENT  = iota
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
	NODE_INTERFACE_FUNCTION_DELCARATION      = iota
	NODE_IMPLEMENT_FUNCTION_DEFINITION = iota
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

type Node struct {
	nodeType int
	token    *Token
	left     *Node
	right    *Node
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

func (this *Parser) advance() error {
	err, token := this.lexer.next()
	if err != nil {
		return err
	}

	this.currentToken = token
	return nil
}

func (this *Parser) parseExpression() (error, *Node) {
	return nil, nil
}

func (this *Parser) parseExpressionStatement() (error, *Node) {
	err, expression := this.parseExpression()
	if err != nil {
		return err, nil
	}

	if this.currentToken.tokenType != TOKEN_EQUAL {
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
	if this.currentToken.tokenType != TOKEN_COLONS {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

	if this.currentToken.tokenType == TOKEN_TYPE_BOOL {
		return nil, &Node { nodeType: NODE_TYPE_BOOL }
	} else if this.currentToken.tokenType == TOKEN_TYPE_INT {
		return nil, &Node { nodeType: NODE_TYPE_INT }
	} else if this.currentToken.tokenType == TOKEN_TYPE_FLOAT {
		return nil, &Node { nodeType: NODE_TYPE_FLOAT }
	} else if this.currentToken.tokenType == TOKEN_TYPE_STRING {
		return nil, &Node { nodeType: NODE_TYPE_STRING }
	} else if this.currentToken.tokenType == TOKEN_IDENTIFIER {
		return nil, &Node { nodeType: NODE_TYPE_CUSTOM }
	}

	return fmt.Errorf("Invalid token"), nil
}

func (this *Parser) parseTypedIdentifier() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_IDENTIFIER {
		return fmt.Errorf("invalid token"), nil
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
	if this.currentToken.tokenType != TOKEN_TYPE_VAR {
		return fmt.Errorf("invalid token"), nil
	}
	this.advance()

	if this.currentToken.tokenType != TOKEN_IDENTIFIER {
		return fmt.Errorf("invalid token"), nil
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
		return fmt.Errorf("invalid token"), nil
	}

	variableNode.left = variableTypeNode
	variableNode.right = expressionNode

	return nil, variableNode
}

func (this *Parser) parseConstant() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_CONST {
		return fmt.Errorf("invalid token"), nil
	}
	this.advance()

	if this.currentToken.tokenType != TOKEN_IDENTIFIER {
		return fmt.Errorf("invalid token"), nil
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

	if this.currentToken.tokenType != TOKEN_EQUAL {
		return fmt.Errorf("invalid Token"), nil
	}
	this.advance()

	err, expressionNode := this.parseExpression()
	if err != nil {
		return err, nil
	}

	constNode.left = constTypeNode
	constNode.right = expressionNode

	return nil, constNode
}

func (this *Parser) parseIf() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_IF {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

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
	if this.currentToken.tokenType != TOKEN_WHILE {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

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
	if this.currentToken.tokenType != TOKEN_FOR {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

	if this.currentToken.tokenType != TOKEN_IDENTIFIER {
		return fmt.Errorf("invalid token"), nil
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

	if this.currentToken.tokenType != TOKEN_IN {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

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
	if this.currentToken.tokenType != TOKEN_RETURN {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

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
	if this.currentToken.tokenType != TOKEN_OPEN_BRACKET {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

	root := &Node{
		nodeType: NODE_STATEMENT,
	}
	currentNode := root

	for this.currentToken.tokenType != TOKEN_CLOSED_BRACKET {
		err, node := this.parseStatement()
		if err != nil {
			return err, nil
		}

		if currentNode.left == nil {
			currentNode.left = node
		} else {
			currentNode.right = &Node{
				nodeType: NODE_LINK,
				left:     node,
			}

			currentNode = node
		}
	}
	this.advance()

	return nil, root
}

func (this *Parser) parseMembers() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_OPEN_BRACKET {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

	root := &Node{
		nodeType: NODE_MEMBER,
	}
	currentNode := root

	for this.currentToken.tokenType != TOKEN_CLOSED_BRACKET {
		err, node := this.parseTypedIdentifier()
		if err != nil {
			return err, nil
		}

		if currentNode.left == nil {
			currentNode.left = node
		} else {
			currentNode.right = &Node{
				nodeType: NODE_LINK,
				left:     node,
			}

			currentNode = node
		}
	}
	this.advance()

	return nil, root
}

func (this *Parser) parseStruct() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_TYPE_STRUCT {
		return fmt.Errorf("invalid token"), nil
	}
	this.advance()

	if this.currentToken.tokenType != TOKEN_IDENTIFIER {
		return fmt.Errorf("invalid token"), nil
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
	if this.currentToken.tokenType != TOKEN_OPEN_PARANTHESIS {
		return fmt.Errorf("invalid token"), nil
	}
	this.advance()

	root := &Node{
		nodeType: NODE_PROGRAM,
	}
	currentNode := root

	for {
		err, node := this.parseTypedIdentifier()
		if err != nil {
			return err, nil
		}

		if currentNode.left == nil {
			currentNode.left = node
		} else {
			currentNode.right = &Node{
				nodeType: NODE_LINK,
				left:     node,
			}

			currentNode = node
		}

		if this.currentToken.tokenType != TOKEN_COMMA {
			break
		}
		this.advance()
	}

	if this.currentToken.tokenType != TOKEN_CLOSED_PARANTHESIS {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

	return nil, root
}

func (this *Parser) parseFunctionDeclaration() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_FUNCTION {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

	if this.currentToken.tokenType != TOKEN_IDENTIFIER {
		return fmt.Errorf("Invalid token"), nil
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

	return nil, parametersNode
}

func (this *Parser) parseInterfaceFunctionDeclarations() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_OPEN_BRACKET {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

	root := &Node{
		nodeType: NODE_INTERFACE_FUNCTION_DELCARATION,
	}
	currentNode := root

	for this.currentToken.tokenType != TOKEN_CLOSED_BRACKET {
		err, node := this.parseFunctionDeclaration()
		if err != nil {
			return err, nil
		}

		if currentNode.left == nil {
			currentNode.left = node
		} else {
			currentNode.right = &Node{
				nodeType: NODE_LINK,
				left:     node,
			}

			currentNode = node
		}
	}
	this.advance()

	return nil, root
}

func (this *Parser) parseFunction() (error, *Node) {
	err, functionDeclarationNode := this.parseFunctionDeclaration()
	if err != nil {
		return fmt.Errorf("invalid token"), nil
	}

	err, statementsNode := this.parseBlock()
	if err != nil {
		return fmt.Errorf("invalid token"), nil
	}

	functionNode := &Node {
		nodeType: NODE_FUNCTION,
		left: functionDeclarationNode,
		right: statementsNode,
	}

	return nil, functionNode
}

func (this *Parser) parseImplementFunctionDefinitions() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_OPEN_BRACKET {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

	root := &Node{
		nodeType: NODE_IMPLEMENT_FUNCTION_DEFINITION,
	}
	currentNode := root

	for this.currentToken.tokenType != TOKEN_CLOSED_BRACKET {
		err, node := this.parseFunction()
		if err != nil {
			return err, nil
		}

		if currentNode.left == nil {
			currentNode.left = node
		} else {
			currentNode.right = &Node{
				nodeType: NODE_LINK,
				left:     node,
			}

			currentNode = node
		}
	}

	this.advance()

	return nil, root
}

func (this *Parser) parseInterface() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_TYPE_INTERFACE {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

	if this.currentToken.tokenType != TOKEN_IDENTIFIER {
		return fmt.Errorf("Invalid token"), nil
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
	if this.currentToken.tokenType != TOKEN_IMPLEMENT {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

	if this.currentToken.tokenType != TOKEN_IDENTIFIER {
		return fmt.Errorf("Invalid token"), nil
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

	return fmt.Errorf("Invalid token"), nil
}

func (this *Parser) parseDeclarations() (error, *Node) {
	root := &Node {
		nodeType: NODE_STATEMENT,
	}

	currentNode := root

	for this.currentToken.tokenType != TOKEN_CLOSED_BRACKET {
		err, node := this.parseDeclaration()
		if err != nil {
			return err, nil
		}

		if currentNode.left == nil {
			currentNode.left = node
		} else {
			currentNode.right = &Node{
				nodeType: NODE_LINK,
				left:     node,
			}

			currentNode = node
		}
	}

	return nil, root
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

	return fmt.Errorf("invalid token"), nil
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

	return fmt.Errorf("Invalid token in root declaration"), nil
}

func (this *Parser) parseModule() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_MODULE {
		return fmt.Errorf("invalid token"), nil
	}
	this.advance()

	root := &Node {
		nodeType: NODE_MODULE,
	}

	currentNode := root

	for {
		if this.currentToken.tokenType != TOKEN_IDENTIFIER {
			return fmt.Errorf("invalid token"), nil
		}

		node := &Node {
			nodeType: NODE_MODULE,
			token: this.currentToken,
		}

		this.advance()

		if currentNode.left == nil {
			currentNode.left = node
		} else {
			currentNode.right = &Node{
				nodeType: NODE_LINK,
				left:     node,
			}

			currentNode = node
		}
		if this.currentToken.tokenType != TOKEN_DOT {
			break
		}
		this.advance()
	}

	return nil, root
}

func (this *Parser) parseImport() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_IMPORT {
		return fmt.Errorf("invalid token"), nil
	}
	this.advance()

	root := &Node {
		nodeType: NODE_IMPORT,
	}

	currentNode := root

	for {
		if this.currentToken.tokenType != TOKEN_IDENTIFIER {
			return fmt.Errorf("invalid token"), nil
		}

		node := &Node {
			nodeType: NODE_IMPORT,
			token: this.currentToken,
		}

		this.advance()

		if currentNode.left == nil {
			currentNode.left = node
		} else {
			currentNode.right = &Node{
				nodeType: NODE_LINK,
				left:     node,
			}

			currentNode = node
		}
		if this.currentToken.tokenType != TOKEN_DOT {
			break
		}
		this.advance()
	}

	return nil, root
}

func (this *Parser) parseImports() (error, *Node) {
	importsNode := &Node {
		nodeType: NODE_IMPORT,
	}

	currentNode := importsNode

	for this.currentToken.tokenType == TOKEN_IMPORT {
		err, node := this.parseImport()
		if err != nil {
			return fmt.Errorf("invalid token"), nil
		}

		if currentNode.left == nil {
			currentNode.left = node
		} else {
			currentNode.right = &Node{
				nodeType: NODE_LINK,
				left:     node,
			}

			currentNode = node
		}
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

	statementsNode := &Node {
		nodeType: NODE_STATEMENT,
	}

	currentNode := statementsNode

	for this.currentToken.tokenType != TOKEN_EOF {
		err, node := this.parseRootStatement()
		if err != nil {
			return err, nil
		}

		if currentNode.left == nil {
			currentNode.left = node
		} else {
			currentNode.right = &Node{
				nodeType: NODE_LINK,
				left:     node,
			}

			currentNode = node
		}
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
