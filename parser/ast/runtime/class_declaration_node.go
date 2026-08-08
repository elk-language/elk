package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initClassDeclarationNode() {
	c := &value.ClassDeclarationNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			constant := args[1].MustReference().(ast.ExpressionNode)

			var body []ast.StatementNode
			if !args[2].IsUndefined() {
				bodyTuple := args[2].AsReference().(value.ArrayTuple)
				body = value.TransformArrayTupleIntoNativeArrayTuple(bodyTuple, func(v value.Value) ast.StatementNode {
					return v.AsReference().(ast.StatementNode)
				}).ToSlice()
			}

			var typeParams []ast.TypeParameterNode
			if !args[3].IsUndefined() {
				typeParamTuple := args[3].AsReference().(value.ArrayTuple)
				typeParams = value.TransformArrayTupleIntoNativeArrayTuple(typeParamTuple, func(v value.Value) ast.TypeParameterNode {
					return v.AsReference().(ast.TypeParameterNode)
				}).ToSlice()
			}

			var abstract bool
			if !args[4].IsUndefined() {
				abstract = value.Truthy(args[4])
			}
			var sealed bool
			if !args[5].IsUndefined() {
				sealed = value.Truthy(args[5])
			}
			var primitive bool
			if !args[6].IsUndefined() {
				primitive = value.Truthy(args[6])
			}
			var noInit bool
			if !args[7].IsUndefined() {
				noInit = value.Truthy(args[7])
			}
			var immutable bool
			if !args[8].IsUndefined() {
				immutable = value.Truthy(args[8])
			}

			var superclass ast.ExpressionNode
			if !args[9].IsUndefined() {
				superclass = args[9].MustReference().(ast.ExpressionNode)
			}
			var docComment string
			if !args[10].IsUndefined() {
				docComment = (string)(args[10].MustReference().(value.String))
			}
			var argLoc *position.Location
			if args[11].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[11].Pointer())
			}

			self := ast.NewClassDeclarationNode(
				argLoc,
				docComment,
				abstract,
				sealed,
				primitive,
				noInit,
				immutable,
				constant,
				typeParams,
				superclass,
				body,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(11),
	)

	vm.Def(
		c,
		"is_abstract",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.BoolVal(self.Abstract)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_sealed",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.BoolVal(self.Sealed)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_primitive",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.BoolVal(self.Primitive)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_no_init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.BoolVal(self.NoInit)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_immutable",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.BoolVal(self.Immutable)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"constant",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
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
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			entries := value.CastNativeArrayTuplePtr(&self.TypeParameters)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"superclass",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			if self.Superclass == nil {
				return value.Nil, value.Undefined
			}
			return value.Ref(self.Superclass), value.Undefined
		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			entries := value.CastNativeArrayTuplePtr(&self.Body)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"bytecode",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.Ref(self.Bytecode)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
