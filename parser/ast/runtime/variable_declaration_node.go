package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initVariableDeclarationNode() {
	c := &value.VariableDeclarationNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argName := (string)(args[0].MustReference().(value.String))

			var argTypeNode ast.TypeNode
			if !args[1].IsUndefined() {
				argTypeNode = args[1].MustReference().(ast.TypeNode)
			}

			var argInitialiser ast.ExpressionNode
			if !args[2].IsUndefined() {
				argInitialiser = args[2].MustReference().(ast.ExpressionNode)
			}

			var argDocComment string
			if !args[3].IsUndefined() {
				argDocComment = string(args[3].MustReference().(value.String))
			}

			var argSpan *position.Span
			if args[4].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[4].Pointer())
			}
			self := ast.NewVariableDeclarationNode(
				argSpan,
				argDocComment,
				argName,
				argTypeNode,
				argInitialiser,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"doc_comment",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariableDeclarationNode)
			result := value.Ref(value.String(self.DocComment()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariableDeclarationNode)
			result := value.Ref(value.String(self.Name))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariableDeclarationNode)
			if self.TypeNode == nil {
				return value.Nil, value.Undefined
			}
			result := value.Ref(self.TypeNode)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"initialiser",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariableDeclarationNode)
			if self.Initialiser == nil {
				return value.Nil, value.Undefined
			}
			result := value.Ref(self.Initialiser)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.VariableDeclarationNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
