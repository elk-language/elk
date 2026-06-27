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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argConstant := args[1].MustReference().(ast.ExpressionNode)

			var argTypeParameters []ast.TypeParameterNode
			if !args[2].IsUndefined() {
				argTypeParametersTuple := args[2].AsReference().(value.ArrayTuple)
				argTypeParameters = value.TransformArrayTupleIntoNativeArrayTuple(argTypeParametersTuple, func(v value.Value) ast.TypeParameterNode {
					return v.AsReference().(ast.TypeParameterNode)
				}).ToSlice()
			}

			var argBody []ast.StructBodyStatementNode
			if !args[3].IsUndefined() {
				argBodyTuple := args[3].AsReference().(value.ArrayTuple)
				argBody = value.TransformArrayTupleIntoNativeArrayTuple(argBodyTuple, func(v value.Value) ast.StructBodyStatementNode {
					return v.AsReference().(ast.StructBodyStatementNode)
				}).ToSlice()
			}

			var argDocComment string
			if !args[4].IsUndefined() {
				argDocComment = string(args[4].MustReference().(value.String))
			}

			var argLoc *position.Location
			if args[5].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[5].Pointer())
			}
			self := ast.NewStructDeclarationNode(
				argLoc,
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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)
			result := value.Ref(value.String(self.DocComment()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"constant",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
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
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)
			entries := value.CastNativeArrayTuplePtr(&self.TypeParameters)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)
			entries := value.CastNativeArrayTuplePtr(&self.Body)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.StructDeclarationNode)
			result := value.Ref(value.String(self.String()))
			return result, value.Undefined
		},
	)

}
