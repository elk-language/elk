package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initCallableTypeNode() {
	c := &value.CallableTypeNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			var argParams []ast.ParameterNode
			if !args[1].IsUndefined() {
				argParamsTuple := args[1].MustReference().(*value.ArrayTuple)
				argParams := make([]ast.ParameterNode, argParamsTuple.Length())
				for i, el := range *argParamsTuple {
					argParams[i] = el.MustReference().(ast.ParameterNode)
				}
			}

			var argReturnType ast.TypeNode
			if !args[2].IsUndefined() {
				argReturnType = args[2].MustReference().(ast.TypeNode)
			}

			var argThrowType ast.TypeNode
			if !args[3].IsUndefined() {
				argThrowType = args[3].MustReference().(ast.TypeNode)
			}

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}

			var argClosure bool
			if !args[5].IsUndefined() {
				argClosure = value.Truthy(args[5])
			}
			self := ast.NewCallableTypeNode(
				argLoc,
				argParams,
				argReturnType,
				argThrowType,
				argClosure,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"parameters",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallableTypeNode)

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
			self := args[0].MustReference().(*ast.CallableTypeNode)
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
			self := args[0].MustReference().(*ast.CallableTypeNode)
			if !args[0].IsUndefined() {
				return value.Nil, value.Undefined
			}

			return value.Ref(self.ThrowType), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallableTypeNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined
		},
	)
	vm.Def(
		c,
		"is_closure",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallableTypeNode)
			result := value.ToElkBool(self.IsClosure)
			return result, value.Undefined
		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallableTypeNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CallableTypeNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
