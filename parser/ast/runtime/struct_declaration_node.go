package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initStructDeclarationNode() {
	c := &value.StructDeclarationNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argConstant := args[0].MustReference().(ast.ExpressionNode)

			var argTypeParameters []ast.TypeParameterNode
			if !args[1].IsUndefined() {
				argTypeParametersTuple := args[1].MustReference().(*value.ArrayTuple)
				argTypeParameters = make([]ast.TypeParameterNode, argTypeParametersTuple.Length())
				for i, el := range *argTypeParametersTuple {
					argTypeParameters[i] = el.MustReference().(ast.TypeParameterNode)
				}
			}

			var argBody []ast.StructBodyStatementNode
			if !args[2].IsUndefined() {
				argBodyTuple := args[2].MustReference().(*value.ArrayTuple)
				argBody = make([]ast.StructBodyStatementNode, argBodyTuple.Length())
				for i, el := range *argBodyTuple {
					argBody[i] = el.MustReference().(ast.StructBodyStatementNode)
				}
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
			self := ast.NewStructDeclarationNode(
				argSpan,
				argDocComment,
				argConstant,
				argTypeParameters,
				argBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"doc_comment",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)
			result := value.Ref(value.String(self.DocComment()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"constant",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)
			if self.Constant == nil {
				return value.Nil, value.Undefined
			}
			result := value.Ref(self.Constant)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"type_parameters",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)

			collection := self.TypeParameters
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
		"body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)

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
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)
			result := value.Ref(value.String(self.String()))
			return result, value.Undefined
		},
	)

}
