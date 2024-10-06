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

	if this.currentToken.tokenType == TOKEN_TYPE_STRUCT {
		return this.parseStruct()
	}

	if this.currentToken.tokenType == TOKEN_TYPE_INTERFACE {
		return this.parseInterface()
	}

	if this.currentToken.tokenType == TOKEN_IMPLEMENT {
		return this.parseImplement()
	}

	if this.currentToken.tokenType == TOKEN_RETURN {
		return this.parseReturn()
	}

	return this.parseExpressionStatement()
}

func (this *Parser) parseRootStatement() (error, *Node) {
	if this.currentToken.tokenType == TOKEN_IMPORT {
		return this.parseImport()
	}

	if this.currentToken.tokenType == TOKEN_FUNCTION {
		return this.parseFunction()
	}

	if this.currentToken.tokenType == TOKEN_CONST {
		return this.parseConstant()
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
