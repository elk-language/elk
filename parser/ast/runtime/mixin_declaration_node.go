package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initMixinDeclarationNode() {
	c := &value.MixinDeclarationNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argConstant := args[1].MustReference().(ast.ExpressionNode)

			var argAbstract bool
			if !args[2].IsUndefined() {
				argAbstract = value.Truthy(args[2])
			}

			var argTypeParameters []ast.TypeParameterNode
			if !args[3].IsUndefined() {
				argTypeParametersTuple := args[3].MustReference().(*value.ArrayTuple)
				argTypeParameters = make([]ast.TypeParameterNode, argTypeParametersTuple.Length())
				for i, el := range *argTypeParametersTuple {
					argTypeParameters[i] = el.MustReference().(ast.TypeParameterNode)
				}
			}

			var argBody []ast.StatementNode
			if !args[4].IsUndefined() {
				argBodyTuple := args[4].MustReference().(*value.ArrayTuple)
				argBody = make([]ast.StatementNode, argBodyTuple.Length())
				for i, el := range *argBodyTuple {
					argBody[i] = el.MustReference().(ast.StatementNode)
				}
			}

			var argDocComment string
			if !args[5].IsUndefined() {
				argDocComment = string(args[5].MustReference().(value.String))
			}

			var argSpan *position.Span
			if args[6].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[6].Pointer())
			}
			self := ast.NewMixinDeclarationNode(
				argSpan,
				argDocComment,
				argAbstract,
				argConstant,
				argTypeParameters,
				argBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(6),
	)

	vm.Def(
		c,
		"doc_comment",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			result := value.Ref(value.String(self.DocComment()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_abstract",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			result := value.ToElkBool(self.Abstract)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"constant",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			if self.Constant == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.Constant), value.Undefined

		},
	)

	vm.Def(
		c,
		"type_parameters",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)

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
			self := args[0].MustReference().(*ast.MixinDeclarationNode)

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
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
