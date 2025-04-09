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
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			constant := args[1].MustReference().(ast.ExpressionNode)

			var body []ast.StatementNode
			if !args[2].IsUndefined() {
				bodyTuple := args[2].MustReference().(*value.ArrayTuple)
				body = make([]ast.StatementNode, bodyTuple.Length())
				for _, el := range *bodyTuple {
					body = append(body, el.MustReference().(ast.StatementNode))
				}
			}

			var typeParams []ast.TypeParameterNode
			if !args[3].IsUndefined() {
				typeParamTuple := args[3].MustReference().(*value.ArrayTuple)
				typeParams = make([]ast.TypeParameterNode, typeParamTuple.Length())
				for _, el := range *typeParamTuple {
					typeParams = append(typeParams, el.MustReference().(ast.TypeParameterNode))
				}
			}

			var abstract bool
			if !args[4].IsUndefined() {
				abstract = value.Truthy(args[3])
			}
			var sealed bool
			if !args[5].IsUndefined() {
				sealed = value.Truthy(args[4])
			}
			var primitive bool
			if !args[6].IsUndefined() {
				primitive = value.Truthy(args[5])
			}
			var noInit bool
			if !args[7].IsUndefined() {
				noInit = value.Truthy(args[6])
			}

			var superclass ast.ExpressionNode
			if !args[8].IsUndefined() {
				superclass = args[8].MustReference().(ast.ExpressionNode)
			}
			var docComment string
			if !args[9].IsUndefined() {
				docComment = (string)(args[9].MustReference().(value.String))
			}

			var argLoc *position.Location
			if args[10].IsUndefined() {
				argLoc = position.DefaultLocation
			} else {
				argLoc = (*position.Location)(args[10].Pointer())
			}
			self := ast.NewClassDeclarationNode(
				argLoc,
				docComment,
				abstract,
				sealed,
				primitive,
				noInit,
				constant,
				typeParams,
				superclass,
				body,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(10),
	)

	vm.Def(
		c,
		"is_abstract",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.ToElkBool(self.Abstract)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_sealed",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.ToElkBool(self.Sealed)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_primitive",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.ToElkBool(self.Primitive)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"is_no_init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.ToElkBool(self.NoInit)
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"constant",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
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
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)

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
		"superclass",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
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
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)

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
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.Ref(self.Bytecode)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
