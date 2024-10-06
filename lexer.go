package main

import (
	"fmt"
	"unicode"
)

const(
	TOKEN_TRUE = iota
	TOKEN_FALSE = iota
	TOKEN_STRING_LITERAL = iota
	TOKEN_INT_LITERAL = iota
	TOKEN_FLOAT_LITERAL = iota

	TOKEN_FUNCTION = iota
	TOKEN_RETURN = iota
	TOKEN_IMPORT = iota
	TOKEN_WITH = iota
	TOKEN_CONST = iota
	TOKEN_EXPORT = iota
	TOKEN_PUBLIC = iota
	TOKEN_PRIVATE = iota

	TOKEN_ASSIGN = iota
	TOKEN_EQUAL = iota
	TOKEN_DIFFERENT = iota
	TOKEN_AND = iota
	TOKEN_OR = iota
	TOKEN_NOT = iota

	TOKEN_PLUS = iota
	TOKEN_MINUS = iota
	TOKEN_DIVIDE = iota
	TOKEN_MULTIPLY = iota

	TOKEN_ADD_ASSIGN = iota
	TOKEN_SUBSTRACT_ASSIGN = iota
	TOKEN_DIVIDE_ASSIGN = iota
	TOKEN_MULTIPLY_ASSIGN = iota

	TOKEN_OPEN_PARANTHESIS = iota
	TOKEN_CLOSED_PARANTHESIS = iota
	TOKEN_OPEN_BRACKET = iota
	TOKEN_CLOSED_BRACKET = iota
	TOKEN_OPEN_SQUARE = iota
	TOKEN_CLOSED_SQUARE = iota

	TOKEN_TYPE_INT = iota
	TOKEN_TYPE_FLOAT = iota
	TOKEN_TYPE_STRING = iota
	TOKEN_TYPE_BOOL = iota

	TOKEN_TYPE_VAR = iota
	TOKEN_TYPE_STRUCT = iota
	TOKEN_TYPE_INTERFACE = iota
	TOKEN_IMPLEMENT = iota

	TOKEN_IF = iota
	TOKEN_ELSE = iota
	TOKEN_WHILE = iota
	TOKEN_FOR = iota

	TOKEN_COMMA = iota
	TOKEN_DOT = iota
	TOKEN_COLONS = iota

	TOKEN_IDENTIFIER = iota

	TOKEN_EOF = iota
)

var tokenTypesString = []string{
	"TOKEN_TRUE",
	"TOKEN_FALSE",
	"TOKEN_STRING_LITERAL",
	"TOKEN_INT_LITERAL",
	"TOKEN_FLOAT_LITERAL",

	"TOKEN_FUNCTION",
	"TOKEN_RETURN",
	"TOKEN_IMPORT",
	"TOKEN_WITH",
	"TOKEN_CONST",
	"TOKEN_EXPORT",
	"TOKEN_PUBLIC",
	"TOKEN_PRIVATE",

	"TOKEN_ASSIGN",
	"TOKEN_EQUAL",
	"TOKEN_DIFFERENT",
	"TOKEN_AND",
	"TOKEN_OR",
	"TOKEN_NOT",

	"TOKEN_PLUS",
	"TOKEN_MINUS",
	"TOKEN_DIVIDE",
	"TOKEN_MULTIPLY",

	"TOKEN_ADD_ASSIGN",
	"TOKEN_SUBSTRACT_ASSIGN",
	"TOKEN_DIVIDE_ASSIGN",
	"TOKEN_MULTIPLY_ASSIGN",

	"TOKEN_OPEN_PARANTHESIS",
	"TOKEN_CLOSED_PARANTHESIS",
	"TOKEN_OPEN_BRACKET",
	"TOKEN_CLOSED_BRACKET",
	"TOKEN_OPEN_SQUARE",
	"TOKEN_CLOSED_SQUARE",

	"TOKEN_TYPE_INT",
	"TOKEN_TYPE_FLOAT",
	"TOKEN_TYPE_STRING",
	"TOKEN_TYPE_BOOL",

	"TOKEN_TYPE_VAR",
	"TOKEN_TYPE_STRUCT",
	"TOKEN_TYPE_INTERFACE",
	"TOKEN_IMPLEMENT",

	"TOKEN_IF",
	"TOKEN_ELSE",
	"TOKEN_WHILE",
	"TOKEN_FOR",

	"TOKEN_COMMA",
	"TOKEN_DOT",
	"TOKEN_COLONS",

	"TOKEN_IDENTIFIER",

	"TOKEN_EOF",
}

type Token struct {
	tokenType int
	tokenValue string
}

func (this *Token) toString() string {
	return fmt.Sprintf(
		"{ tokenType: %s, tokenValue: %s}",
		tokenTypesString[this.tokenType],
		this.tokenValue,
	)
}

type Lexer struct {
	text string

	currentPosition int
	currentLine int
	currentColumn int
}

func newLexer(text string) *Lexer {
	return &Lexer{
		text: text,
		currentPosition: 0,
		currentLine: 1,
		currentColumn: 0,
	}
}

func (this *Lexer) parsePlus() (error, *Token) {
	currentCharacter := this.text[this.currentPosition]
	if currentCharacter != '+' {
		return fmt.Errorf("invalid character"), nil
	}

	err := this.advance()
	if err == nil && this.text[this.currentPosition] == '=' {
		return this.SimpleToken(TOKEN_ADD_ASSIGN)
	}

	return nil, &Token {tokenType: TOKEN_PLUS}
}

func (this *Lexer) parseMinus() (error, *Token) {
	currentCharacter := this.text[this.currentPosition]
	if currentCharacter != '-' {
		return fmt.Errorf("invalid character"), nil
	}

	err := this.advance()
	if err == nil && this.text[this.currentPosition] == '=' {
		return this.SimpleToken(TOKEN_SUBSTRACT_ASSIGN)
	}

	return nil, &Token {tokenType: TOKEN_MINUS}
}

func (this *Lexer) parseMultiply() (error, *Token) {
	currentCharacter := this.text[this.currentPosition]
	if currentCharacter != '*' {
		return fmt.Errorf("invalid character"), nil
	}

	err := this.advance()
	if err == nil && this.text[this.currentPosition] == '=' {
		return this.SimpleToken(TOKEN_MULTIPLY_ASSIGN)
	}

	return nil, &Token {tokenType: TOKEN_MULTIPLY}
}

func (this *Lexer) parseDivide() (error, *Token) {
	currentCharacter := this.text[this.currentPosition]
	if currentCharacter != '/' {
		return fmt.Errorf("invalid character"), nil
	}

	err := this.advance()
	if err == nil && this.text[this.currentPosition] == '=' {
		return this.SimpleToken(TOKEN_DIVIDE_ASSIGN)
	}

	return nil, &Token {tokenType: TOKEN_DIVIDE}
}

func (this *Lexer) parseEqual() (error, *Token) {
	currentCharacter := this.text[this.currentPosition]
	if currentCharacter != '=' {
		return fmt.Errorf("invalid character"), nil
	}

	err := this.advance()
	if err == nil && this.text[this.currentPosition] == '=' {
		return this.SimpleToken(TOKEN_EQUAL)
	}

	return nil, &Token {tokenType: TOKEN_ASSIGN}
}

func (this *Lexer) parseNot() (error, *Token) {
	currentCharacter := this.text[this.currentPosition]
	if currentCharacter != '!' {
		return fmt.Errorf("invalid character"), nil
	}

	err := this.advance()
	if err == nil && this.text[this.currentPosition] == '=' {
		return this.SimpleToken(TOKEN_DIFFERENT)
	}

	return nil, &Token {tokenType: TOKEN_NOT}
}

func (this *Lexer) getKeyword(text string) (error, *Token) {
	switch text {
	case "if":
		return nil, &Token{tokenType: TOKEN_IF}
	case "else":
		return nil, &Token{tokenType: TOKEN_ELSE}
	case "for":
		return nil, &Token{tokenType: TOKEN_FOR}
	case "while":
		return nil, &Token{tokenType: TOKEN_WHILE}
	case "int":
		return nil, &Token{tokenType: TOKEN_TYPE_INT}
	case "float":
		return nil, &Token{tokenType: TOKEN_TYPE_FLOAT}
	case "string":
		return nil, &Token{tokenType: TOKEN_TYPE_STRING}
	case "bool":
		return nil, &Token{tokenType: TOKEN_TYPE_BOOL}
	case "true":
		return nil, &Token{tokenType: TOKEN_TRUE}
	case "false":
		return nil, &Token{tokenType: TOKEN_FALSE}
	case "and":
		return nil, &Token{tokenType: TOKEN_AND}
	case "or":
		return nil, &Token{tokenType: TOKEN_OR}
	case "not":
		return nil, &Token{tokenType: TOKEN_NOT}
	case "var":
		return nil, &Token{tokenType: TOKEN_TYPE_VAR}
	case "struct":
		return nil, &Token{tokenType: TOKEN_TYPE_STRUCT}
	case "interface":
		return nil, &Token{tokenType: TOKEN_TYPE_INTERFACE}
	case "implement":
		return nil, &Token{tokenType: TOKEN_IMPLEMENT}
	case "function":
		return nil, &Token{tokenType: TOKEN_FUNCTION}
	case "return":
		return nil, &Token{tokenType: TOKEN_RETURN}
	case "import":
		return nil, &Token{tokenType: TOKEN_IMPORT}
	case "with":
		return nil, &Token{tokenType: TOKEN_WITH}
	case "const":
		return nil, &Token{tokenType: TOKEN_CONST}
	case "export":
		return nil, &Token{tokenType: TOKEN_EXPORT}
	case "public":
		return nil, &Token{tokenType: TOKEN_PUBLIC}
	case "private":
		return nil, &Token{tokenType: TOKEN_PRIVATE}
	}

	return fmt.Errorf("Invalid keyword"), nil
}

func (this *Lexer) parseIdentifier() (error, *Token) {
	currentCharacter := this.text[this.currentPosition]

	if !this.isIdentifierKeywordLetter(true, currentCharacter) {
		return fmt.Errorf("Invalid character"), nil
	}

	value := string(currentCharacter)
	for {
		err := this.advance()
		if err != nil {
			break
		}

		currentCharacter = this.text[this.currentPosition]

		if (!this.isIdentifierKeywordLetter(false, currentCharacter)) {
			break
		}
		
		value = value + string(currentCharacter)
	}

	err, keywordToken := this.getKeyword(value)
	if err == nil {
		return nil, keywordToken
	}

	return nil, &Token{
		tokenType: TOKEN_IDENTIFIER,
		tokenValue: value,
	}
}

func (this *Lexer) advance() error {
	this.currentPosition += 1
	if this.currentPosition >= len(this.text) {
		return fmt.Errorf("End of file reached")
	}

	if this.text[this.currentPosition] == '\n' {
		this.currentLine += 1
		this.currentColumn = 0
	} else {
		this.currentColumn += 1
	}

	return nil
}

func (this *Lexer) parseString() (error, *Token) {
	currentCharacter := this.text[this.currentPosition]
	if currentCharacter != '"' {
		return fmt.Errorf("Invalid character"), nil
	}

	value := ""
	for {
		err := this.advance()
		if err != nil {
			return fmt.Errorf("String opened but not closed"), nil
		}

		currentCharacter = this.text[this.currentPosition]
		if currentCharacter == '"' {
			break
		}

		value = value + string(currentCharacter)
	}

	this.advance()
	return nil, &Token{
		tokenType: TOKEN_STRING_LITERAL,
		tokenValue: value,
	}
}

func (this *Lexer) parseNumber() (error, *Token) {
	currentCharacter := this.text[this.currentPosition]
	if currentCharacter < '0' || currentCharacter > '9' {
		return fmt.Errorf("Invalid character"), nil
	}

	commaFound := false
	startWithZero := currentCharacter == '0'
	value := string(currentCharacter)

	for {
		err := this.advance()
		if err != nil {
			break
		}

		currentCharacter = this.text[this.currentPosition]
		if currentCharacter == '.' {
			if commaFound {
				return fmt.Errorf("Invalid number"), nil
			}

			commaFound = true
			this.advance()

			continue
		}

		if currentCharacter < '0' || currentCharacter > '9' {
			break
		}

		if startWithZero {
			return fmt.Errorf("Can't have multiple zeros at start of a number"), nil
		}

		value = value + string(currentCharacter)
	}

	if commaFound {
		return nil, &Token{
			tokenType: TOKEN_FLOAT_LITERAL,
			tokenValue: value,
		}
	}

	return nil, &Token{
		tokenType: TOKEN_INT_LITERAL,
		tokenValue: value,
	}
}

func (this *Lexer) isIdentifierKeywordLetter(firstCharacter bool, character byte) bool {
	if character == '_' {
		return true
	}

	if character >= 'a' && character <= 'z' {
		return true
	}

	if character >= 'A' && character <= 'Z' {
		return true
	}
	
	if firstCharacter {
		return false
	}

	return character >= '0' && character <= '9'
}

func (this *Lexer) SimpleToken(token_type int) (error, *Token) {
	this.advance()
	return nil, &Token{
		tokenType: token_type,
	}
}

func (this *Lexer) next() (error, *Token) {
	var currentCharacter byte
	for {
		if this.currentPosition >= len(this.text) {
			return nil, &Token{ tokenType: TOKEN_EOF }
		}

		currentCharacter = this.text[this.currentPosition]
		if !unicode.IsSpace(rune(currentCharacter)) {
			break
		}

		this.advance()
	}

	switch currentCharacter {
	case '+':
		return this.parsePlus()
	case '-':
		return this.parseMinus()
	case '*':
		return this.parseMultiply()
	case '/':
		return this.parseDivide()
	case '=':
		return this.parseEqual()
	case '!':
		return this.parseNot()
	case '(':
		return this.SimpleToken(TOKEN_OPEN_PARANTHESIS)
	case ')':
		return this.SimpleToken(TOKEN_CLOSED_PARANTHESIS)
	case '{':
		return this.SimpleToken(TOKEN_OPEN_BRACKET)
	case '}':
		return this.SimpleToken(TOKEN_CLOSED_BRACKET)
	case '[':
		return this.SimpleToken(TOKEN_OPEN_SQUARE)
	case ']':
		return this.SimpleToken(TOKEN_CLOSED_SQUARE)
	case ',':
		return this.SimpleToken(TOKEN_COMMA)
	case ':':
		return this.SimpleToken(TOKEN_COLONS)
	case '"':
		return this.parseString()
	}

	if currentCharacter >= '0' && currentCharacter <= '9' {
		return this.parseNumber()
	}

	if this.isIdentifierKeywordLetter(true, currentCharacter) {
		return this.parseIdentifier()
	}

	return fmt.Errorf("Invalid token"), nil
}
