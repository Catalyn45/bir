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

func (this *Parser) parseExpressionStatement() (error, *Node) {
	return nil, nil
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
		this.advance()

		var err error
		err, variableTypeNode = this.parseType()
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

func (this *Parser) parseIf() (error, *Node) {
	if this.currentToken.tokenType != TOKEN_IF {
		return fmt.Errorf("Invalid token"), nil
	}
	this.advance()

	err, expressionNode := this.parseExpression()
	if err != nil {
		return err, nil
	}

	err, statementsNode := this.parseStatements()
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
			err, elseNode = this.parseStatements()
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

	err, statementsNode := this.parseStatements()
	if err != nil {
		return err, nil
	}

	branchNode := Node {
		nodeType: NODE_BRANCH,
		left: statementsNode,
	}

	if this.currentToken.tokenType == TOKEN_ELSE {
		this.advance()

		err, elseNode := this.parseStatements()
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

func (this *Parser) parseStatements() (error, *Node) {
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
}

func (this *Parser) parseRootStatement() (error, *Node) {
	if this.currentToken.tokenType == TOKEN_IMPORT {
		return this.parseImport()
	}

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

func (this *Parser) parse() (error, *Node) {
	err := this.advance()
	if err != nil {
		return err, nil
	}

	root := &Node{
		nodeType: NODE_PROGRAM,
	}
	currentNode := root

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

	return nil, root
}
