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
	asts []*Node
	irModule *ir.Module
	symbolTables Stack[*SymbolTable]
	currentFunction *ir.Func
	blocks Stack[*ir.Block]
}

func newCompiler(asts []*Node) *Compiler {
	m := ir.NewModule()
	return &Compiler {
		asts: asts,
		irModule: m,
		symbolTables: Stack[*SymbolTable]{},
		blocks: Stack[*ir.Block]{},
	}
}

func (this *Compiler) searchSymbol(symbolName string) (error, *Symbol) {
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

func (this *Compiler) converType(birType string) (error, types.Type) {
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

	return fmt.Errorf("Invalid type, can't convert"), nil
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

	if pointerType, ok := leftValue.Type().(*types.PointerType); ok {
		leftValue = block.NewLoad(pointerType.ElemType, leftValue)
	}

	err, rightValue := this.walkExpression(node.right)
	if err != nil {
		return err, nil
	}

	if pointerType, ok := rightValue.Type().(*types.PointerType); ok {
		rightValue = block.NewLoad(pointerType.ElemType, rightValue)
	}

	if node.token.tokenType == TOKEN_PLUS {
		return nil, block.NewAdd(leftValue, rightValue)
	}

	if node.token.tokenType == TOKEN_MINUS {
		return nil, block.NewSub(leftValue, rightValue)
	}

	if node.token.tokenType == TOKEN_DIVIDE {
		if leftValue.Type() == types.Float {
			return nil, block.NewFDiv(leftValue, rightValue)
		} else {
			return nil, block.NewSDiv(leftValue, rightValue)
		}
	}

	if node.token.tokenType == TOKEN_MULTIPLY {
		return nil, block.NewMul(leftValue, rightValue)
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

	if node.nodeType == NODE_VARIABLE {
		err, symbol := this.searchSymbol(node.token.tokenValue)
		if err != nil {
			return err, nil
		}

		return nil, symbol.value
	}

	if node.nodeType == NODE_BINARY_EXPRESSION {
		return this.walkBinaryExpression(node)
	}

	if node.nodeType == NODE_CALL {
		err, funcValue := this.walkExpression(node.left)
		if err != nil {
			return err, nil
		}

		var arguments []value.Value
		for argument := node.right.right; argument != nil; argument = argument.next {
			err, argumentValue := this.walkExpression(argument)
			if err != nil {
				return err, nil
			}

			if pointerType, ok := argumentValue.Type().(*types.PointerType); ok {
				argumentValue = block.NewLoad(pointerType.ElemType, argumentValue)
			}

			arguments = append(arguments, argumentValue)
		}

		return nil, block.NewCall(funcValue, arguments...)
	}

	if node.nodeType == NODE_VARIABLE_DECLARATION {
		err, irType := this.converType(node.symbol.simbolType.name)
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

	return fmt.Errorf("can't eval expression"), nil
}

func (this *Compiler) walk(node *Node) error {
	for node != nil {
		if node.nodeType == NODE_RETURN {
			err, returnValue := this.walkExpression(node.left)
			if err != nil {
				return err
			}

			block := this.blocks.pop()

			// deref pointer in case of pointer type
			// TODO don't do that if struct
			if pointerType, ok := returnValue.Type().(*types.PointerType); ok {
				returnValue = block.NewLoad(pointerType.ElemType, returnValue)
			}

			block.NewRet(returnValue)
		} else if node.nodeType == NODE_ASSIGNMENT {
			err, assignmentSource := this.walkExpression(node.left)
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
		if node.nodeType == NODE_FUNCTION {
			this.symbolTables.push(node.left.symbolTable)

			symbol := node.left.symbol
			birSignature := symbol.simbolType.signature

			err, returnType := this.converType(birSignature.returnType)
			if err != nil {
				return err
			}

			var signature []*ir.Param
			for _, birparam := range birSignature.parameters {
				err, paramType := this.converType(birparam.paramType)
				if err != nil {
					return err
				}

				parameter := ir.NewParam(birparam.name, paramType)

				birparam.node.symbol.value = parameter

				signature = append(signature, parameter)
			}

			function := this.irModule.NewFunc(
				symbol.name,
				returnType,
				signature...
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
