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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argConstant := args[1].MustReference().(ast.ExpressionNode)

			var argAbstract bool
			if !args[2].IsUndefined() {
				argAbstract = value.Truthy(args[2])
			}

			var argTypeParameters []ast.TypeParameterNode
			if !args[3].IsUndefined() {
				argTypeParametersTuple := args[3].AsReference().(value.ArrayTuple)
				argTypeParameters = value.TransformArrayTupleIntoNativeArrayTuple(argTypeParametersTuple, func(v value.Value) ast.TypeParameterNode {
					return v.AsReference().(ast.TypeParameterNode)
				}).ToSlice()
			}

			var argBody []ast.StatementNode
			if !args[4].IsUndefined() {
				argBodyTuple := args[4].AsReference().(value.ArrayTuple)
				argBody = value.TransformArrayTupleIntoNativeArrayTuple(argBodyTuple, func(v value.Value) ast.StatementNode {
					return v.AsReference().(ast.StatementNode)
				}).ToSlice()
			}

			var argDocComment string
			if !args[5].IsUndefined() {
				argDocComment = string(args[5].MustReference().(value.String))
			}

			var argLoc *position.Location
			if args[6].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[6].Pointer())
			}
			self := ast.NewMixinDeclarationNode(
				argLoc,
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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			result := value.Ref(value.String(self.DocComment()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_abstract",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			result := value.BoolVal(self.Abstract)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"constant",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			entries := value.CastNativeArrayTuplePtr(&self.TypeParameters)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			entries := value.CastNativeArrayTuplePtr(&self.Body)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MixinDeclarationNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
