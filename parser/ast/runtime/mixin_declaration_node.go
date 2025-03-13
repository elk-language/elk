package ast

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
			argDocComment := string(args[0].MustReference().(value.String))
			argAbstract := value.Truthy(args[0])
			argConstant := args[1].MustReference().(ast.ExpressionNode)

			argTypeParametersTuple := args[2].MustReference().(*value.ArrayTuple)
			argTypeParameters := make([]ast.TypeParameterNode, argTypeParametersTuple.Length())
			for i, el := range *argTypeParametersTuple {
				argTypeParameters[i] = el.MustReference().(ast.TypeParameterNode)
			}

			argBodyTuple := args[3].MustReference().(*value.ArrayTuple)
			argBody := make([]ast.StatementNode, argBodyTuple.Length())
			for i, el := range *argBodyTuple {
				argBody[i] = el.MustReference().(ast.StatementNode)
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
		"abstract",
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
			result := value.Ref(self.Constant)
			return result, value.Undefined

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
		"includes_and_implements",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)

			collection := self.IncludesAndImplements
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
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			result := value.Ref(self.Bytecode)
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

}
