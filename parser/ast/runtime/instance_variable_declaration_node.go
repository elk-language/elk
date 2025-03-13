package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initInstanceVariableDeclarationNode() {
	c := &value.InstanceVariableDeclarationNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argDocComment := (string)(args[0].MustReference().(value.String))
			argName := (string)(args[1].MustReference().(value.String))
			argTypeNode := args[2].MustReference().(ast.TypeNode)

			var argSpan *position.Span
			if args[3].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[3].Pointer())
			}
			self := ast.NewInstanceVariableDeclarationNode(
				argSpan,
				argDocComment,
				argName,
				argTypeNode,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InstanceVariableDeclarationNode)
			result := value.Ref(value.String(self.Name))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InstanceVariableDeclarationNode)
			result := value.Ref(self.TypeNode)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.InstanceVariableDeclarationNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
