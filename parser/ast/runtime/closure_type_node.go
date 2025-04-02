package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initClosureTypeNode() {
	c := &value.ClosureTypeNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			var argParams []ast.ParameterNode
			if !args[0].IsUndefined() {
				argParamsTuple := args[0].MustReference().(*value.ArrayTuple)
				argParams := make([]ast.ParameterNode, argParamsTuple.Length())
				for i, el := range *argParamsTuple {
					argParams[i] = el.MustReference().(ast.ParameterNode)
				}
			}

			var argReturnType ast.TypeNode
			if !args[1].IsUndefined() {
				argReturnType = args[1].MustReference().(ast.TypeNode)
			}

			var argThrowType ast.TypeNode
			if !args[2].IsUndefined() {
				argThrowType = args[2].MustReference().(ast.TypeNode)
			}

			var argSpan *position.Span
			if args[3].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[3].Pointer())
			}
			self := ast.NewClosureTypeNode(
				argSpan,
				argParams,
				argReturnType,
				argThrowType,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"parameters",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureTypeNode)

			collection := self.Parameters
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
		"return_type",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureTypeNode)
			if self.ReturnType == nil {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.ReturnType), value.Undefined
		},
	)

	vm.Def(
		c,
		"throw_type",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureTypeNode)
			if !args[0].IsUndefined() {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.ThrowType), value.Undefined
		},
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureTypeNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined
		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureTypeNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ClosureTypeNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
