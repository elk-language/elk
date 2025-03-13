package ast

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
			docComment := (string)(args[0].MustReference().(value.String))
			abstract := value.Truthy(args[0])
			sealed := value.Truthy(args[1])
			primitive := value.Truthy(args[2])
			noInit := value.Truthy(args[3])
			constant := args[4].MustReference().(ast.ExpressionNode)

			typeParamTuple := args[5].MustReference().(*value.ArrayTuple)
			typeParams := make([]ast.TypeParameterNode, typeParamTuple.Length())
			for _, el := range *typeParamTuple {
				typeParams = append(typeParams, el.MustReference().(ast.TypeParameterNode))
			}
			superclass := args[6].MustReference().(ast.ExpressionNode)

			bodyTuple := args[7].MustReference().(*value.ArrayTuple)
			body := make([]ast.StatementNode, bodyTuple.Length())
			for _, el := range *bodyTuple {
				body = append(body, el.MustReference().(ast.StatementNode))
			}

			var argSpan *position.Span
			if args[8].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[9].Pointer())
			}
			self := ast.NewClassDeclarationNode(
				argSpan,
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
		"abstract",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.ToElkBool(self.Abstract)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"sealed",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.ToElkBool(self.Sealed)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"primitive",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClassDeclarationNode)
			result := value.ToElkBool(self.Primitive)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"no_init",
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
			result := value.Ref(self.Constant)
			return result, value.Undefined

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
			result := value.Ref(self.Superclass)
			return result, value.Undefined

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

}
