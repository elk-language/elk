package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initModuleDeclarationNode() {
	c := &value.ModuleDeclarationNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argDocComment := string(args[0].MustReference().(value.String))
			argConstant := args[1].MustReference().(ast.ExpressionNode)

			argBodyTuple := args[2].MustReference().(*value.ArrayTuple)
			argBody := make([]ast.StatementNode, argBodyTuple.Length())
			for i, el := range *argBodyTuple {
				argBody[i] = el.MustReference().(ast.StatementNode)
			}

			var argSpan *position.Span
			if args[3].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[3].Pointer())
			}
			self := ast.NewModuleDeclarationNode(
				argSpan,
				argDocComment,
				argConstant,
				argBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"constant",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModuleDeclarationNode)
			result := value.Ref(self.Constant)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModuleDeclarationNode)

			collection := self.Body
			arrayTuple := value.NewArrayTupleWithLength(len(collection))
			for i, el := range collection {
				arrayTuple.SetAt(i, value.Ref(el))
			}
			result := value.Ref(arrayTuple)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"bytecode",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModuleDeclarationNode)
			result := value.Ref(self.Bytecode)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModuleDeclarationNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
