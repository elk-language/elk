package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initConstantDeclarationNode() {
	c := &value.ConstantDeclarationNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			constant := args[0].MustReference().(ast.ExpressionNode)

			var typeNode ast.TypeNode
			if !args[1].IsUndefined() {
				typeNode = args[1].MustReference().(ast.TypeNode)
			}

			var init ast.ExpressionNode
			if !args[2].IsUndefined() {
				init = args[2].MustReference().(ast.ExpressionNode)
			}

			var docComment string
			if !args[3].IsUndefined() {
				docComment = string(args[3].MustReference().(value.String))
			}

			var argSpan *position.Span
			if args[4].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[4].Pointer())
			}
			self := ast.NewConstantDeclarationNode(
				argSpan,
				docComment,
				constant,
				typeNode,
				init,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"doc_comment",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstantDeclarationNode)
			result := value.Ref(value.String(self.DocComment()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"constant",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstantDeclarationNode)
			result := value.Ref(self.Constant)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstantDeclarationNode)
			if self.TypeNode == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.TypeNode), value.Undefined

		},
	)

	vm.Def(
		c,
		"initialiser",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstantDeclarationNode)
			if self.Initialiser == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.Initialiser), value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstantDeclarationNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
