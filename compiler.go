package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	// "github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Compiler struct {
	asts []*Node
	irModule *ir.Module
	symbolTables Stack[*SymbolTable]
	currentFunction *ir.Func
}

func newCompile(asts []*Node) *Compiler {
	m := ir.NewModule()
	return &Compiler {
		asts: asts,
		irModule: m,
		symbolTables: Stack[*SymbolTable]{},
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

func (this *Compiler) convertToValue(varType types.Type) types.Type {
	if varType == types.I64Ptr {
		return types.I64
	}

	return nil
}

func (this *Compiler) walkExpression(node *Node) (error, value.Value) {
	if node.nodeType == NODE_INT {
		intValue, err := strconv.Atoi(node.token.tokenValue)
		if err != nil {
			return err, nil
		}

		return nil, constant.NewInt(types.I64, int64(intValue))
	}

	if node.nodeType == NODE_FLOAT {

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

	return fmt.Errorf("can't eval expression"), nil
}

func (this *Compiler) walk(node *Node, block *ir.Block) error {
	for node != nil {
		if node.nodeType == NODE_VARIABLE_DECLARATION {
			err, irType := this.converType(node.symbol.simbolType.name)
			if err != nil {
				return err
			}

			node.symbol.value = block.NewAlloca(irType)
		}

		if node.nodeType == NODE_RETURN {
			err, returnValue := this.walkExpression(node.left)
			if err != nil {
				return err
			}

			// deref pointer in case of pointer type
			// TODO don't do that if struct
			if pointerType, ok := returnValue.Type().(*types.PointerType); ok {
				returnValue = block.NewLoad(pointerType.ElemType, returnValue)
			}

			retBlock := this.currentFunction.NewBlock("")
			
			retBlock.NewRet(returnValue)

			block.NewBr(retBlock)
		}

		if node.nodeType == NODE_ASSIGNMENT {
			err, assignmentSource := this.walkExpression(node.left)
			if err != nil {
				return err
			}

			err, assignmentValue := this.walkExpression(node.right)
			if err != nil {
				return err
			}

			block.NewStore(assignmentValue, assignmentSource)
		}

		if node.nodeType == NODE_IF {
			err, ifExpression := this.walkExpression(node.left)
			if err != nil {
				return err
			}

			// exitCtx := this.currentFunction.NewBlock("if.exit")

			outerCtx := this.currentFunction.NewBlock("")

			thenCtx := this.currentFunction.NewBlock("")
			err = this.walk(node.right.left, thenCtx)
			if err != nil {
				return err
			}

			thenCtx.NewRet(nil)

			// thenCtx.NewBr(exitCtx)
			// elseCtx := this.currentFunction.NewBlock("if.else")

			// err = this.walk(node.right.right, thenCtx)
			// if err != nil {
			// 	return err
			// }

			// elseCtx.NewBr(exitCtx)

			// exitCtx.NewRet(nil)


			tr := outerCtx.NewCondBr(ifExpression, thenCtx, nil)
			if tr != nil {
				print(tr)
			}
		}

		node = node.next
	}

	return nil
}

func (this *Compiler) walkRoot(node *Node) error {
	for node != nil {
		if node.nodeType == NODE_FUNCTION {
			this.symbolTables.push(node.symbolTable)

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

				signature = append(signature, ir.NewParam(birparam.name, paramType))
			}

			function := this.irModule.NewFunc(
				symbol.name,
				returnType,
				signature...
			)

			this.currentFunction = function

			block := function.NewBlock("")
			err = this.walk(node.right, block)
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
		err := this.walkRoot(ast.right)
		if err != nil {
			return err
		}
	}

	program := this.irModule.String()

	fmt.Println(program)

	err := os.WriteFile("main.ll", []byte(program), 0644)
	if err != nil {
		return err
	}

	return nil
}
