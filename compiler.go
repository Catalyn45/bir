package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Compiler struct {
	asts               []*Node
	irModule           *ir.Module
	symbolTables       Stack[*SymbolTable]
	currentFunction    *ir.Func
	currentStruct      types.Type
	blocks             Stack[*ir.Block]
	currentInstance    value.Value
}

func newCompiler(asts []*Node) *Compiler {
	m := ir.NewModule()
	return &Compiler{
		asts:               asts,
		irModule:           m,
		symbolTables:       Stack[*SymbolTable]{},
		blocks:             Stack[*ir.Block]{},
	}
}

func (this *Compiler) searchSymbol(symbolName string) (error, *Symbol) {
	var foundSymbol *Symbol = nil

	this.symbolTables.foreach(func(item *SymbolTable) (stop bool) {
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

func (this *Compiler) convertType(birType string) (error, types.Type) {
	switch birType {
	case "int":
		return nil, types.I64
	case "float":
		return nil, types.Double
	case "bool":
		return nil, types.I8
	case "void":
		return nil, types.Void
	}

	err, symbol := this.searchSymbol(birType)
	if err != nil {
		return err, nil
	}

	return nil, symbol.valueType
}

func (this *Compiler) walkBinaryExpression(node *Node) (error, value.Value) {
	if node.nodeType != NODE_BINARY_EXPRESSION {
		return fmt.Errorf("Not binary expression"), nil
	}

	block := this.blocks.peek()

	err, leftValue := this.walkExpression(node.left)
	if err != nil {
		return err, nil
	}

	err, rightValue := this.walkExpression(node.right)
	if err != nil {
		return err, nil
	}

	if node.token.tokenType == TOKEN_PLUS {
		if leftValue.Type() == types.Float {
			return nil, block.NewFAdd(leftValue, rightValue)
		} else {
			return nil, block.NewAdd(leftValue, rightValue)
		}
	}

	if node.token.tokenType == TOKEN_MINUS {
		if leftValue.Type() == types.Float {
			return nil, block.NewFSub(leftValue, rightValue)
		} else {
			return nil, block.NewSub(leftValue, rightValue)
		}
	}

	if node.token.tokenType == TOKEN_DIVIDE {
		if leftValue.Type() == types.Float {
			return nil, block.NewFDiv(leftValue, rightValue)
		} else {
			return nil, block.NewSDiv(leftValue, rightValue)
		}
	}

	if node.token.tokenType == TOKEN_MULTIPLY {
		if leftValue.Type() == types.Float {
			return nil, block.NewFMul(leftValue, rightValue)
		} else {
			return nil, block.NewMul(leftValue, rightValue)
		}
	}

	if node.token.tokenType == TOKEN_AND {
		return nil, block.NewAnd(leftValue, rightValue)
	}

	if node.token.tokenType == TOKEN_OR {
		return nil, block.NewOr(leftValue, rightValue)
	}

	if node.token.tokenType == TOKEN_EQUAL {
		if leftValue.Type() == types.Float {
			return nil, block.NewFCmp(enum.FPredOEQ, leftValue, rightValue)
		} else {
			return nil, block.NewICmp(enum.IPredEQ, leftValue, rightValue)
		}
	}

	if node.token.tokenType == TOKEN_DIFFERENT {
		if leftValue.Type() == types.Float {
			return nil, block.NewFCmp(enum.FPredONE, leftValue, rightValue)
		} else {
			return nil, block.NewICmp(enum.IPredNE, leftValue, rightValue)
		}
	}

	if node.token.tokenType == TOKEN_GREATER {
		if leftValue.Type() == types.Float {
			return nil, block.NewFCmp(enum.FPredOGT, leftValue, rightValue)
		} else {
			return nil, block.NewICmp(enum.IPredSGT, leftValue, rightValue)
		}
	}

	if node.token.tokenType == TOKEN_GREATER_EQUAL {
		if leftValue.Type() == types.Float {
			return nil, block.NewFCmp(enum.FPredOGE, leftValue, rightValue)
		} else {
			return nil, block.NewICmp(enum.IPredSGE, leftValue, rightValue)
		}
	}

	if node.token.tokenType == TOKEN_LESS {
		if leftValue.Type() == types.Float {
			return nil, block.NewFCmp(enum.FPredOLT, leftValue, rightValue)
		} else {
			return nil, block.NewICmp(enum.IPredSLT, leftValue, rightValue)
		}
	}

	if node.token.tokenType == TOKEN_LESS_EQUAL {
		if leftValue.Type() == types.Float {
			return nil, block.NewFCmp(enum.FPredOLE, leftValue, rightValue)
		} else {
			return nil, block.NewICmp(enum.IPredSLE, leftValue, rightValue)
		}
	}

	return fmt.Errorf("invalid operation"), nil
}

func (this *Compiler) walkLvalue(node *Node) (error, value.Value) {
	if node.nodeType == NODE_VARIABLE {
		err, symbol := this.searchSymbol(node.token.tokenValue)
		if err != nil {
			return err, nil
		}

		if symbol.value == nil {
			block := this.blocks.peek()

			allocated := block.NewAlloca(symbol.valueType)

			loaded := block.NewLoad(allocated.ElemType, allocated)

			return nil, loaded
		}

		return nil, symbol.value
	}

	if node.nodeType == NODE_MEMBER_ACCESS {
		err, value := this.walkExpression(node.left)
		if err != nil {
			return err, nil
		}

		if _, ok := value.Type().(*types.StructType); ok {
			err, symbol := this.searchSymbol(value.Type().Name())
			if err != nil {
				return err, nil
			}

			this.symbolTables.push(symbol.node.symbolTable)

			err, symbol = this.searchSymbol(node.token.tokenValue)
			if err != nil {
				return err, nil
			}

			if symbol.node.nodeType == NODE_FUNCTION_DECLARATION {
				this.currentInstance = value
			}

			return nil, symbol.value
		}
		
		return nil, value
	}

	return fmt.Errorf("can't eval lvalue expression"), nil
}

func (this *Compiler) walkExpression(node *Node) (error, value.Value) {
	block := this.blocks.peek()

	if node.nodeType == NODE_INT {
		intValue, err := strconv.ParseInt(node.token.tokenValue, 10, 64)
		if err != nil {
			return err, nil
		}

		return nil, constant.NewInt(types.I64, intValue)
	}

	if node.nodeType == NODE_FLOAT {
		floatValue, err := strconv.ParseFloat(node.token.tokenValue, 64)
		if err != nil {
			return err, nil
		}

		return nil, constant.NewFloat(types.Float, floatValue)
	}

	if node.nodeType == NODE_BOOL {
		if node.token.tokenValue == "true" {
			return nil, constant.NewBool(true)
		} else {
			return nil, constant.NewBool(false)
		}
	}

	if node.nodeType == NODE_BINARY_EXPRESSION {
		return this.walkBinaryExpression(node)
	}

	if node.nodeType == NODE_CALL {
		err, funcValue := this.walkLvalue(node.left)
		if err != nil {
			return err, nil
		}

		if _, ok := funcValue.Type().(*types.StructType); ok {
			return nil, funcValue
		}

		var arguments []value.Value

		if this.currentInstance != nil {
			arguments = append(arguments, this.currentInstance)
			this.currentInstance = nil
			this.symbolTables.pop()
		}

		for argument := node.right.right; argument != nil; argument = argument.next {
			err, argumentValue := this.walkExpression(argument)
			if err != nil {
				return err, nil
			}

			arguments = append(arguments, argumentValue)
		}

		return nil, block.NewCall(funcValue, arguments...)
	}

	if node.nodeType == NODE_VARIABLE_DECLARATION {
		err, irType := this.convertType(node.symbol.simbolType.name)
		if err != nil {
			return err, nil
		}

		block := this.blocks.peek()
		allocationValue := block.NewAlloca(irType)

		node.symbol.value = allocationValue

		if node.right != nil {
			err, initValue := this.walkExpression(node.right)
			if err != nil {
				return err, nil
			}

			block.NewStore(initValue, allocationValue)
		}

		return nil, block.NewLoad(allocationValue.ElemType, allocationValue)
	}

	err, expressionValue := this.walkLvalue(node)
	if err != nil {
		return err, nil
	}

	if pointerType, ok := expressionValue.Type().(*types.PointerType); ok {
		expressionValue = block.NewLoad(pointerType.ElemType, expressionValue)
	}

	return nil, expressionValue
}

func (this *Compiler) walk(node *Node) error {
	for node != nil {
		if node.nodeType == NODE_RETURN {
			err, returnValue := this.walkExpression(node.left)
			if err != nil {
				return err
			}

			block := this.blocks.pop()

			block.NewRet(returnValue)
		} else if node.nodeType == NODE_ASSIGNMENT {
			err, assignmentSource := this.walkLvalue(node.left)
			if err != nil {
				return err
			}

			err, assignmentValue := this.walkExpression(node.right)
			if err != nil {
				return err
			}

			block := this.blocks.peek()

			block.NewStore(assignmentValue, assignmentSource)
		} else if node.nodeType == NODE_IF {
			err, ifExpression := this.walkExpression(node.left)
			if err != nil {
				return err
			}

			thenBlock := this.currentFunction.NewBlock("")
			elseBlock := this.currentFunction.NewBlock("")
			exitBlock := this.currentFunction.NewBlock("")

			this.symbolTables.push(node.symbolTable)

			this.symbolTables.push(node.right.left.symbolTable)
			this.blocks.push(thenBlock)

			err = this.walk(node.right.left)
			if err != nil {
				return err
			}

			if thenBlock.Term == nil {
				this.blocks.pop()
				thenBlock.NewBr(exitBlock)
			}
			this.symbolTables.pop()

			this.blocks.push(elseBlock)

			if node.right.right != nil {
				this.symbolTables.push(node.right.right.symbolTable)

				err = this.walk(node.right.right)
				if err != nil {
					return err
				}

				this.symbolTables.pop()
			}

			if elseBlock.Term == nil {
				this.blocks.pop()
				elseBlock.NewBr(exitBlock)
			}

			block := this.blocks.pop()
			block.NewCondBr(ifExpression, thenBlock, elseBlock)

			this.symbolTables.pop()

			this.blocks.push(exitBlock)
		} else {
			err, _ := this.walkExpression(node)
			if err != nil {
				return err
			}
		}

		node = node.next
	}

	return nil
}

func (this *Compiler) walkRoot(node *Node) error {
	for node != nil {
		if node.nodeType == NODE_STRUCT {
			this.symbolTables.push(node.symbolTable)

			var fieldTypes []types.Type
			for field := node.right; field != nil; field = field.next {
				err, convertedType := this.convertType(field.symbol.simbolType.name)
				if err != nil {
					return err
				}

				fieldTypes = append(fieldTypes, convertedType)
			}

			structType := this.irModule.NewTypeDef(node.token.tokenValue, types.NewStruct(fieldTypes...))

			node.symbol.valueType = structType

			this.symbolTables.pop()
		} else if node.nodeType == NODE_IMPLEMENT {
			strcutName := node.token.tokenValue

			err, found := this.searchSymbol(strcutName)
			if err != nil {
				return err
			}

			this.currentStruct = found.valueType

			err = this.walkRoot(node.right)
			if err != nil {
				return err
			}

			this.currentStruct = nil
		} else if node.nodeType == NODE_FUNCTION {
			this.symbolTables.push(node.left.symbolTable)

			symbol := node.left.symbol
			birSignature := symbol.simbolType.signature

			err, returnType := this.convertType(birSignature.returnType)
			if err != nil {
				return err
			}

			var signature []*ir.Param

			for _, birparam := range birSignature.parameters {
				err, paramType := this.convertType(birparam.paramType)
				if err != nil {
					return err
				}

				parameter := ir.NewParam(birparam.name, paramType)

				birparam.node.symbol.value = parameter

				signature = append(signature, parameter)
			}

			functionPrefix := ""
			if this.currentStruct != nil {
				selfParameter := ir.NewParam("this", this.currentStruct)
				signature = append([]*ir.Param{selfParameter}, signature...)
				functionPrefix = this.currentStruct.Name() + "_"
			}

			function := this.irModule.NewFunc(
				functionPrefix + symbol.name,
				returnType,
				signature...,
			)

			node.left.symbol.value = function

			this.currentFunction = function

			block := function.NewBlock("")

			this.blocks.push(block)
			err = this.walk(node.right)
			if err != nil {
				return err
			}

			this.currentFunction = nil

			this.symbolTables.pop()
		}

		node = node.next
	}

	return nil
}

func (this *Compiler) Compile() error {
	for _, ast := range this.asts {
		this.symbolTables.push(ast.symbolTable)

		err := this.walkRoot(ast.right)
		if err != nil {
			return err
		}

		this.symbolTables.pop()
	}

	program := this.irModule.String()

	fmt.Println(program)

	err := os.WriteFile("main.ll", []byte(program), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("clang", "main.ll", "-o", "output.exe")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
