package main

import "fmt"

const (
	NODE_PROGRAM              = iota
	NODE_EXPRESSION           = iota
	NODE_ASSIGNMENT           = iota
	NODE_PATH                 = iota
	NODE_DECLARATION          = iota
	NODE_PARAMETER            = iota
	NODE_STATEMENT            = iota
	NODE_RETURN               = iota
	NODE_IF                   = iota
	NODE_VARIABLE             = iota
	NODE_INT                  = iota
	NODE_FLOAT                = iota
	NODE_STRING               = iota
	NODE_WHILE                = iota
	NODE_FOR                  = iota
	NODE_MEMBER               = iota
	NODE_MODULE               = iota
	NODE_CONST                = iota
	NODE_ITERATION            = iota
	NODE_STRUCT               = iota
	NODE_FUNCTION             = iota
	NODE_CONSTRUCTOR          = iota
	NODE_IMPLEMENT            = iota
	NODE_INTERFACE            = iota
	NODE_IMPORT               = iota
	NODE_BOOL_TYPE            = iota
	NODE_INT_TYPE             = iota
	NODE_FLOAT_TYPE           = iota
	NODE_STRING_TYPE          = iota
	NODE_CUSTOM_TYPE          = iota
	NODE_BINARY_EXPRESSION    = iota
	NODE_NOT                  = iota
	NODE_CALL                 = iota
	NODE_MEMBER_ACCESS        = iota
	NODE_INDEX                = iota
	NODE_WITH                 = iota
	NODE_LINK                 = iota
)

var nodeStrings = []string{
	"NODE_PROGRAM",
	"NODE_EXPRESSION",
	"NODE_ASSIGNMENT",
	"NODE_PATH",
	"NODE_DECLARATION",
	"NODE_PARAMETER",
	"NODE_STATEMENT",
	"NODE_RETURN",
	"NODE_IF",
	"NODE_VARIABLE",
	"NODE_INT",
	"NODE_FLOAT",
	"NODE_STRING",
	"NODE_WHILE",
	"NODE_FOR",
	"NODE_MEMBER",
	"NODE_MODULE",
	"NODE_CONST",
	"NODE_ITERATION",
	"NODE_STRUCT",
	"NODE_FUNCTION",
	"NODE_CONSTRUCTOR",
	"NODE_IMPLEMENT",
	"NODE_INTERFACE",
	"NODE_IMPORT",
	"NODE_BOOL_TYPE",
	"NODE_INT_TYPE",
	"NODE_FLOAT_TYPE",
	"NODE_STRING_TYPE",
	"NODE_CUSTOM_TYPE",
	"NODE_BINARY_EXPRESSION",
	"NODE_NOT",
	"NODE_CALL",
	"NODE_MEMBER_ACCESS",
	"NODE_INDEX",
	"NODE_WITH",
	"NODE_LINK",
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

func contains(arr []int, target int) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

func clone(arr []int) []int {
	// Create a new slice of the same length
	clone := make([]int, len(arr))

	// Copy the original slice to the new one
	copy(clone, arr)

	return clone
}

const (
	emptySymbol         = " "
	barSymbol           = "│"
	normalNodeSymbol    = "└─> "
	continousNodeSymbol = "├─> "
	nextNodeSymbol      = " ─> "
)

func (this *Node) Dump(indent int, bars *[]int, nodeSymbol string) {
	for i := 0; i < indent-4; i++ {
		if contains(*bars, i) {
			fmt.Print(barSymbol)
		} else {
			fmt.Print(emptySymbol)
		}
	}

	fmt.Print(nodeSymbol)

	fmt.Println(this.ToString())

	clonedBars := clone(*bars)

	leftNodeCharacter := normalNodeSymbol
	if this.right != nil {
		clonedBars = append(clonedBars, indent)
		leftNodeCharacter = continousNodeSymbol
	}

	if this.left != nil {
		this.left.Dump(indent+4, &clonedBars, leftNodeCharacter)
	}

	if this.right != nil {
		this.right.Dump(indent+4, bars, normalNodeSymbol)
	}

	if this.next != nil {
		this.next.Dump(indent, bars, nextNodeSymbol)
	}
}

type Parser struct {
	lexer        *Lexer
	currentToken *Token

	asAllowed bool
}

func newParser(lexer *Lexer) *Parser {
	return &Parser{
		lexer:        lexer,
		currentToken: nil,
		asAllowed: false,
	}
}

func (this *Parser) invalidTokenError(expectedTokenType int) error {
	err := fmt.Errorf("Invalid token: %s, expected: %s, line: %d, column: %d", this.currentToken.toString(), tokenTypesString[expectedTokenType], this.currentToken.line, this.currentToken.column)

	panic(err)
}

func (this *Parser) unexpectedTokenError() error {
	err := fmt.Errorf("Unexpected token: %s, line: %d, column: %d", this.currentToken.toString(), this.currentToken.line, this.currentToken.column)
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

func (this *Parser) parseLiteral() (error, *Node) {
	var literalNode *Node
	if this.currentToken.tokenType == TOKEN_INT_LITERAL {
		literalNode = &Node{
			nodeType: NODE_INT,
			token:    this.currentToken,
		}
	} else if this.currentToken.tokenType == TOKEN_FLOAT_LITERAL {
		literalNode = &Node{
			nodeType: NODE_FLOAT,
			token:    this.currentToken,
		}
	} else if this.currentToken.tokenType == TOKEN_STRING_LITERAL {
		literalNode = &Node{
			nodeType: NODE_STRING,
			token:    this.currentToken,
		}
	} else {
		return this.unexpectedTokenError(), nil
	}

	this.advance()

	return nil, literalNode
}

func (this *Parser) parsePrimary() (error, *Node) {
	if this.currentToken.tokenType == TOKEN_OPEN_PARANTHESIS {
		this.advance()

		err, primaryNode := this.parseExpression()
		if err != nil {
			return err, nil
		}

		err = this.eat(TOKEN_CLOSED_PARANTHESIS)
		if err != nil {
			return err, nil
		}

		return nil, primaryNode
	} else if this.currentToken.tokenType == TOKEN_IDENTIFIER {
		primaryNode := &Node{
			nodeType: NODE_VARIABLE,
			token:    this.currentToken,
		}

		this.advance()

		return nil, primaryNode
	}

	return this.parseLiteral()
}

func (this *Parser) parseArguments() (error, *Node) {
	err := this.eat(TOKEN_OPEN_PARANTHESIS)
	if err != nil {
		return err, nil
	}

	var argumentsNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_PARANTHESIS; {
		err, node := this.parseExpression()
		if err != nil {
			return err, nil
		}

		if currentNode == nil {
			argumentsNode = node
		} else {
			currentNode.next = node
		}

		currentNode = node

		if this.currentToken.tokenType == TOKEN_COMMA {
			this.advance()
		}
	}
	this.advance()

	return nil, argumentsNode
}

func (this *Parser) parsePostfix() (error, *Node) {
	err, left := this.parsePrimary()
	if err != nil {
		return err, nil
	}

	for {
		var templateNode *Node = nil
		if this.currentToken.tokenType == TOKEN_DOUBLE_COLONS {
			err, templateNode = this.parseTemplateSpecification()
			if err != nil {
				return err, nil
			}

			err = this.expectToken(TOKEN_OPEN_PARANTHESIS)
			if err != nil {
				return err, nil
			}
		}

		if this.currentToken.tokenType == TOKEN_OPEN_PARANTHESIS {
			err, arguments := this.parseArguments()
			if err != nil {
				return err, nil
			}

			left = &Node{
				nodeType: NODE_CALL,
				left:     left,
				right:    &Node{
					nodeType: NODE_LINK,
					left: templateNode,
					right: arguments,
				},
			}
		} else if this.asAllowed && this.currentToken.tokenType == TOKEN_AS {
			this.advance()

			err, identifier := this.parseIdentifier(false)
			if err != nil {
				return err, nil
			}

			identifier.right = left
			left = identifier
		} else if this.currentToken.tokenType == TOKEN_DOT {
			this.advance()

			err := this.expectToken(TOKEN_IDENTIFIER)
			if err != nil {
				return err, nil
			}

			left = &Node{
				nodeType: NODE_MEMBER_ACCESS,
				token:    this.currentToken,
				left:     left,
			}

			this.advance()
		} else if this.currentToken.tokenType == TOKEN_OPEN_SQUARE {
			this.advance()

			err, expression := this.parseExpression()
			if err != nil {
				return err, nil
			}

			left = &Node{
				nodeType: NODE_INDEX,
				left:     left,
				right:    expression,
			}

			err = this.eat(TOKEN_CLOSED_SQUARE)
			if err != nil {
				return err, nil
			}
		} else {
			break
		}
	}

	return nil, left
}

func (this *Parser) parseUnary() (error, *Node) {
	if this.currentToken.tokenType == TOKEN_NOT {
		this.advance()

		err, expression := this.parseUnary()
		if err != nil {
			return err, nil
		}

		return nil, &Node{
			nodeType: NODE_NOT,
			left:     expression,
		}
	}

	return this.parsePostfix()
}

func (this *Parser) parseMultiplicative() (error, *Node) {
	err, left := this.parseUnary()
	if err != nil {
		return err, nil
	}

	for currentToken := this.currentToken; currentToken.tokenType == TOKEN_MULTIPLY || currentToken.tokenType == TOKEN_DIVIDE; currentToken = this.currentToken {
		this.advance()

		err, right := this.parseUnary()
		if err != nil {
			return err, nil
		}

		left = &Node{
			nodeType: NODE_BINARY_EXPRESSION,
			token:    currentToken,
			left:     left,
			right:    right,
		}
	}

	return nil, left
}

func (this *Parser) parseAdditive() (error, *Node) {
	err, left := this.parseMultiplicative()
	if err != nil {
		return err, nil
	}

	for currentToken := this.currentToken; currentToken.tokenType == TOKEN_PLUS || currentToken.tokenType == TOKEN_MINUS; currentToken = this.currentToken {
		this.advance()

		err, right := this.parseMultiplicative()
		if err != nil {
			return err, nil
		}

		left = &Node{
			nodeType: NODE_BINARY_EXPRESSION,
			token:    currentToken,
			left:     left,
			right:    right,
		}
	}

	return nil, left
}

func (this *Parser) parseRelational() (error, *Node) {
	err, left := this.parseAdditive()
	if err != nil {
		return err, nil
	}

	for currentToken := this.currentToken; currentToken.tokenType == TOKEN_GREATER ||
		currentToken.tokenType == TOKEN_GREATER_EQUAL ||
		currentToken.tokenType == TOKEN_LESS ||
		currentToken.tokenType == TOKEN_LESS_EQUAL; currentToken = this.currentToken {
		this.advance()

		err, right := this.parseAdditive()
		if err != nil {
			return err, nil
		}

		left = &Node{
			nodeType: NODE_BINARY_EXPRESSION,
			token:    currentToken,
			left:     left,
			right:    right,
		}
	}

	return nil, left
}

func (this *Parser) parseEquality() (error, *Node) {
	err, left := this.parseRelational()
	if err != nil {
		return err, nil
	}

	for currentToken := this.currentToken; currentToken.tokenType == TOKEN_EQUAL ||
		currentToken.tokenType == TOKEN_DIFFERENT; currentToken = this.currentToken {
		this.advance()

		err, right := this.parseRelational()
		if err != nil {
			return err, nil
		}

		left = &Node{
			nodeType: NODE_BINARY_EXPRESSION,
			token:    currentToken,
			left:     left,
			right:    right,
		}
	}

	return nil, left
}

func (this *Parser) parseAnd() (error, *Node) {
	err, left := this.parseEquality()
	if err != nil {
		return err, nil
	}

	for currentToken := this.currentToken; currentToken.tokenType == TOKEN_AND; currentToken = this.currentToken {
		this.advance()

		err, right := this.parseEquality()
		if err != nil {
			return err, nil
		}

		left = &Node{
			nodeType: NODE_BINARY_EXPRESSION,
			token:    currentToken,
			left:     left,
			right:    right,
		}
	}

	return nil, left
}

func (this *Parser) parseOr() (error, *Node) {
	err, left := this.parseAnd()
	if err != nil {
		return err, nil
	}

	for currentToken := this.currentToken; currentToken.tokenType == TOKEN_OR; currentToken = this.currentToken {
		this.advance()

		err, right := this.parseAnd()
		if err != nil {
			return err, nil
		}

		left = &Node{
			nodeType: NODE_BINARY_EXPRESSION,
			token:    currentToken,
			left:     left,
			right:    right,
		}
	}

	return nil, left
}

func (this *Parser) parseExpression() (error, *Node) {
	return this.parseOr()
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

	assignmentNode := &Node{
		nodeType: NODE_ASSIGNMENT,
		left:     expression,
		right:    rightExpression,
	}

	return nil, assignmentNode
}

func (this *Parser) parseTemplate() (error, *Node) {
	err := this.eat(TOKEN_LESS)
	if err != nil {
		return err, nil
	}

	var templateNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_GREATER; {
		err, node := this.parseType()
		if err != nil {
			return err, nil
		}

		if currentNode == nil {
			templateNode = node
		} else {
			currentNode.next = node
		}

		currentNode = node

		if this.currentToken.tokenType == TOKEN_COMMA {
			this.advance()
		}
	}

	this.advance()

	return nil, templateNode
}

func (this *Parser) parseTemplateSpecification() (error, *Node) {
	err := this.eat(TOKEN_DOUBLE_COLONS)
	if err != nil {
		return err, nil
	}

	return this.parseTemplate()
}

func (this *Parser) parseType() (error, *Node) {
	var node *Node
	if this.currentToken.tokenType == TOKEN_BOOL {
		node = &Node{nodeType: NODE_BOOL_TYPE}
	} else if this.currentToken.tokenType == TOKEN_INT {
		node = &Node{nodeType: NODE_INT_TYPE}
	} else if this.currentToken.tokenType == TOKEN_FLOAT {
		node = &Node{nodeType: NODE_FLOAT_TYPE}
	} else if this.currentToken.tokenType == TOKEN_STRING {
		node = &Node{nodeType: NODE_STRING_TYPE}
	} else if this.currentToken.tokenType == TOKEN_IDENTIFIER {
		node = &Node{nodeType: NODE_CUSTOM_TYPE, token: this.currentToken}
	} else {
		return this.unexpectedTokenError(), nil
	}

	this.advance()

	if this.currentToken.tokenType == TOKEN_LESS {
		err, templateNode := this.parseTemplate()
		if err != nil {
			return err, nil
		}

		node.left = templateNode
	}

	return nil, node
}

func (this *Parser) parseTypeSpecification() (error, *Node) {
	err := this.eat(TOKEN_COLONS)
	if err != nil {
		return err, nil
	}

	return this.parseType()
}

func (this *Parser) parseIdentifier(typeRequired bool) (error, *Node) {
	err := this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	variableNode := &Node{
		nodeType: NODE_VARIABLE,
		token:    this.currentToken,
	}

	this.advance()

	var typeNode *Node = nil
	if this.currentToken.tokenType == TOKEN_COLONS || typeRequired {
		err, typeNode = this.parseTypeSpecification()
		if err != nil {
			return err, nil
		}
	}

	variableNode.left = typeNode

	return nil, variableNode
}

func (this *Parser) parseVariableDeclaration() (error, *Node) {
	err := this.eat(TOKEN_VAR)
	if err != nil {
		return err, nil
	}

	err, variableNode := this.parseIdentifier(false)
	if err != nil {
		return err, nil
	}

	var expressionNode *Node = nil
	if this.currentToken.tokenType == TOKEN_ASSIGN {
		this.advance()

		var err error
		err, expressionNode = this.parseExpression()
		if err != nil {
			return err, nil
		}
	}

	if variableNode.left == nil && expressionNode == nil {
		return fmt.Errorf("Variable needs to be either typed or initialized"), nil
	}

	variableNode.right = expressionNode

	return nil, variableNode
}

func (this *Parser) parseConstant() (error, *Node) {
	err := this.eat(TOKEN_CONST)
	if err != nil {
		return err, nil
	}

	err, constNode := this.parseIdentifier(false)
	if err != nil {
		return err, nil
	}

	err = this.eat(TOKEN_ASSIGN)
	if err != nil {
		return err, nil
	}

	err, literalNode := this.parseLiteral()
	if err != nil {
		return err, nil
	}

	constNode.right = literalNode

	return nil, constNode
}

func (this *Parser) parseIf() (error, *Node) {
	err := this.eat(TOKEN_IF)
	if err != nil {
		return err, nil
	}

	this.asAllowed = true
	err, expressionNode := this.parseExpression()
	if err != nil {
		return err, nil
	}
	this.asAllowed = false

	err, statementsNode := this.parseStatementsBlock()
	if err != nil {
		return err, nil
	}

	branchNode := Node{
		nodeType: NODE_LINK,
		left:     statementsNode,
	}

	if this.currentToken.tokenType == TOKEN_ELSE {
		this.advance()

		var elseNode *Node
		var err error

		if this.currentToken.tokenType == TOKEN_IF {
			err, elseNode = this.parseIf()
		} else {
			err, elseNode = this.parseStatementsBlock()
		}

		if err != nil {
			return err, nil
		}

		branchNode.right = elseNode
	}

	ifNode := &Node{
		nodeType: NODE_IF,
		left:     expressionNode,
		right:    &branchNode,
	}

	return nil, ifNode
}

func (this *Parser) parseWhile() (error, *Node) {
	err := this.eat(TOKEN_WHILE)
	if err != nil {
		return err, nil
	}

	this.asAllowed = true
	err, expressionNode := this.parseExpression()
	if err != nil {
		return err, nil
	}
	this.asAllowed = false

	err, statementsNode := this.parseStatementsBlock()
	if err != nil {
		return err, nil
	}

	branchNode := Node{
		nodeType: NODE_LINK,
		left:     statementsNode,
	}

	if this.currentToken.tokenType == TOKEN_ELSE {
		this.advance()

		err, elseNode := this.parseStatementsBlock()
		if err != nil {
			return err, nil
		}

		branchNode.right = elseNode
	}

	whileNode := &Node{
		nodeType: NODE_IF,
		left:     expressionNode,
		right:    &branchNode,
	}

	return nil, whileNode
}

func (this *Parser) parseFor() (error, *Node) {
	err := this.eat(TOKEN_FOR)
	if err != nil {
		return err, nil
	}

	err, iteratorNode := this.parseIdentifier(false)
	if err != nil {
		return err, nil
	}

	err = this.eat(TOKEN_IN)
	if err != nil {
		return err, nil
	}

	err, iterableNode := this.parseExpression()
	if err != nil {
		return err, nil
	}

	err, statementsNode := this.parseStatementsBlock()
	if err != nil {
		return err, nil
	}

	iterationNode := &Node{
		nodeType: NODE_ITERATION,
		left:     iteratorNode,
		right:    iterableNode,
	}

	forNode := &Node{
		nodeType: NODE_FOR,
		left:     iterationNode,
		right:    statementsNode,
	}

	return nil, forNode
}

func (this *Parser) parseReturn() (error, *Node) {
	err := this.eat(TOKEN_RETURN)
	if err != nil {
		return err, nil
	}

	returnNode := &Node{
		nodeType: NODE_RETURN,
	}

	if this.currentToken.tokenType == TOKEN_CLOSED_BRACKET {
		return nil, returnNode
	}

	err, expressionNode := this.parseExpression()
	if err != nil {
		return err, nil
	}

	returnNode.left = expressionNode

	return nil, returnNode
}

func (this *Parser) parseWith() (error, *Node) {
	err := this.eat(TOKEN_WITH)
	if err != nil {
		return err, nil
	}

	this.asAllowed = true
	err, expression := this.parseExpression()
	if err != nil {
		return err, nil
	}
	this.asAllowed = false

	var statementsNode *Node = nil
	if this.currentToken.tokenType == TOKEN_OPEN_BRACKET {
		err, statementsNode = this.parseStatementsBlock()
		if err != nil {
			return err, nil
		}
	}

	return nil, &Node{
		nodeType: NODE_WITH,
		left: &Node{
			nodeType: NODE_LINK,
			right:    expression,
		},
		right: statementsNode,
	}
}

func (this *Parser) parseStatement() (error, *Node) {
	if this.currentToken.tokenType == TOKEN_CONST {
		return this.parseConstant()
	}

	if this.currentToken.tokenType == TOKEN_VAR {
		return this.parseVariableDeclaration()
	}

	if this.currentToken.tokenType == TOKEN_IF {
		return this.parseIf()
	}

	if this.currentToken.tokenType == TOKEN_WHILE {
		return this.parseWhile()
	}

	if this.currentToken.tokenType == TOKEN_WITH {
		return this.parseWith()
	}

	if this.currentToken.tokenType == TOKEN_FOR {
		return this.parseFor()
	}

	if this.currentToken.tokenType == TOKEN_RETURN {
		return this.parseReturn()
	}

	return this.parseExpressionStatement()
}

func (this *Parser) parseStatementsBlock() (error, *Node) {
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

func (this *Parser) parseStructBlock() (error, *Node) {
	err := this.eat(TOKEN_OPEN_BRACKET)
	if err != nil {
		return err, nil
	}

	var membersNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_BRACKET; {
		err, node := this.parseIdentifier(true)
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
	err := this.eat(TOKEN_STRUCT)
	if err != nil {
		return err, nil
	}

	err = this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	structNode := &Node{
		nodeType: NODE_STRUCT,
		token:    this.currentToken,
	}

	this.advance()

	var templateNode *Node = nil
	if this.currentToken.tokenType == TOKEN_DOUBLE_COLONS {
		err, templateNode = this.parseTemplateSpecification()
		if err != nil {
			return err, nil
		}
	}

	err, membersNode := this.parseStructBlock()
	if err != nil {
		return err, nil
	}

	structNode.left = templateNode
	structNode.right = membersNode

	return nil, structNode
}

func (this *Parser) parseFunctionParameters() (error, *Node) {
	err := this.eat(TOKEN_OPEN_PARANTHESIS)
	if err != nil {
		return err, nil
	}

	var parametersNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_PARANTHESIS; {
		err, node := this.parseIdentifier(true)
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

func (this *Parser) parseFunctionDeclaration(isConstructor bool) (error, *Node) {
	if !isConstructor {
		err := this.eat(TOKEN_FUNCTION)
		if err != nil {
			return err, nil
		}
	}

	err := this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	functionDeclarationNode := &Node{
		nodeType: NODE_LINK,
		token:    this.currentToken,
	}

	this.advance()

	var templateNode *Node = nil
	if this.currentToken.tokenType == TOKEN_DOUBLE_COLONS {
		var err error
		err, templateNode = this.parseTemplateSpecification()
		if err != nil {
			return err, nil
		}
	}

	err, parametersNode := this.parseFunctionParameters()
	if err != nil {
		return err, nil
	}

	var functionTypeNode *Node = nil
	if !isConstructor && this.currentToken.tokenType == TOKEN_COLONS {
		err, functionTypeNode = this.parseTypeSpecification()
		if err != nil {
			return err, nil
		}
	}

	functionDeclarationNode.left = &Node {
		nodeType: NODE_LINK,
		left: functionTypeNode,
		right: templateNode,
	}

	functionDeclarationNode.right = parametersNode

	return nil, functionDeclarationNode
}

func (this *Parser) parseInterfaceBlock() (error, *Node) {
	err := this.eat(TOKEN_OPEN_BRACKET)
	if err != nil {
		return err, nil
	}

	var functionDeclarationsNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_BRACKET; {
		err, node := this.parseFunctionDeclaration(false)
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

func (this *Parser) parseFunction(isConstructor bool) (error, *Node) {
	err, functionDeclarationNode := this.parseFunctionDeclaration(isConstructor)
	if err != nil {
		return err, nil
	}

	err, statementsNode := this.parseStatementsBlock()
	if err != nil {
		return err, nil
	}

	nodeType := NODE_FUNCTION
	if isConstructor {
		nodeType = NODE_CONSTRUCTOR
	}

	functionNode := &Node{
		nodeType: nodeType,
		left:     functionDeclarationNode,
		right:    statementsNode,
	}

	return nil, functionNode
}

func (this *Parser) parseImplementBlock() (error, *Node) {
	err := this.eat(TOKEN_OPEN_BRACKET)
	if err != nil {
		return err, nil
	}

	var functionsNode *Node = nil

	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_BRACKET; {
		var isConstructor = false
		if this.currentToken.tokenType == TOKEN_IDENTIFIER && this.currentToken.tokenValue == "init" {
			isConstructor = true
		}

		err, node := this.parseFunction(isConstructor)
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
	err := this.eat(TOKEN_INTERFACE)
	if err != nil {
		return err, nil
	}

	err = this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	interfaceNode := &Node{
		nodeType: NODE_INTERFACE,
		token:    this.currentToken,
	}

	this.advance()

	var templateNode *Node = nil
	if this.currentToken.tokenType == TOKEN_DOUBLE_COLONS {
		err, templateNode = this.parseTemplateSpecification()
		if err != nil {
			return err, nil
		}
	}

	err, functionDeclarationsNode := this.parseInterfaceBlock()
	if err != nil {
		return err, nil
	}

	interfaceNode.left = templateNode
	interfaceNode.right = functionDeclarationsNode

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

	implementNode := &Node{
		nodeType: NODE_IMPLEMENT,
		token:    this.currentToken,
	}

	this.advance()

	var templateNode *Node = nil
	if this.currentToken.tokenType == TOKEN_DOUBLE_COLONS {
		err, templateNode = this.parseTemplateSpecification()
		if err != nil {
			return err, nil
		}
	}

	err, functionDefinitionsNode := this.parseImplementBlock()
	if err != nil {
		return err, nil
	}

	implementNode.left = templateNode
	implementNode.right = functionDefinitionsNode

	return nil, implementNode
}

func (this *Parser) parseExportDeclaration() (error, *Node) {
	if this.currentToken.tokenType == TOKEN_STRUCT {
		return this.parseStruct()
	}

	if this.currentToken.tokenType == TOKEN_INTERFACE {
		return this.parseInterface()
	}

	if this.currentToken.tokenType == TOKEN_IMPLEMENT {
		return this.parseImplement()
	}

	if this.currentToken.tokenType == TOKEN_FUNCTION {
		return this.parseFunction(false)
	}

	if this.currentToken.tokenType == TOKEN_CONST {
		return this.parseConstant()
	}

	return this.unexpectedTokenError(), nil
}

func (this *Parser) parseExportBlock() (error, *Node) {
	var declarationsNode *Node = nil
	for currentNode := (*Node)(nil); this.currentToken.tokenType != TOKEN_CLOSED_BRACKET; {
		err, node := this.parseExportDeclaration()
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
	err := this.eat(TOKEN_EXPORT)
	if err != nil {
		return err, nil
	}

	if this.currentToken.tokenType == TOKEN_OPEN_BRACKET {
		return this.parseExportBlock()
	}

	return this.parseExportDeclaration()
}

func (this *Parser) parseRootStatement() (error, *Node) {
	if this.currentToken.tokenType == TOKEN_STRUCT {
		return this.parseStruct()
	}

	if this.currentToken.tokenType == TOKEN_INTERFACE {
		return this.parseInterface()
	}

	if this.currentToken.tokenType == TOKEN_IMPLEMENT {
		return this.parseImplement()
	}

	if this.currentToken.tokenType == TOKEN_FUNCTION {
		return this.parseFunction(false)
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
	err := this.expectToken(TOKEN_IDENTIFIER)
	if err != nil {
		return err, nil
	}

	pathNode := &Node{
		nodeType: NODE_PATH,
		token:    this.currentToken,
	}

	this.advance()

	for this.currentToken.tokenType == TOKEN_DOT {
		this.advance()

		err := this.expectToken(TOKEN_IDENTIFIER)
		if err != nil {
			return err, nil
		}

		pathNode = &Node{
			nodeType: NODE_PATH,
			token:    this.currentToken,
			left:     pathNode,
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

	moduleNode := &Node{
		nodeType: NODE_MODULE,
		left:     pathNode,
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

	importNode := &Node{
		nodeType: NODE_IMPORT,
		left:     pathNode,
	}

	if this.currentToken.tokenType == TOKEN_AS {
		this.advance()

		err = this.expectToken(TOKEN_IDENTIFIER)
		if err != nil {
			return err, nil
		}

		importNode.token = this.currentToken

		this.advance()
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

	programMetadataNode := &Node{
		nodeType: NODE_LINK,
		left:     moduleNode,
		right:    importsNode,
	}

	root := &Node{
		nodeType: NODE_PROGRAM,
		left:     programMetadataNode,
		right:    statementsNode,
	}

	return nil, root
}
