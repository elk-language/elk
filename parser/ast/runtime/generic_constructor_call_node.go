package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initGenericConstructorCallNode() {
	c := &value.GenericConstructorCallNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argClass := args[1].MustReference().(ast.ComplexConstantNode)

			argTypArgsTuple := args[2].MustReference().(*value.ArrayTuple)
			argTypeArgs := make([]ast.TypeNode, argTypArgsTuple.Length())
			for i, el := range *argTypArgsTuple {
				argTypeArgs[i] = el.MustReference().(ast.TypeNode)
			}
			var argPosArgs []ast.ExpressionNode
			if !args[3].IsUndefined() {
				argPosArgsTuple := args[3].MustReference().(*value.ArrayTuple)
				argPosArgs = make([]ast.ExpressionNode, argPosArgsTuple.Length())
				for i, el := range *argPosArgsTuple {
					argPosArgs[i] = el.MustReference().(ast.ExpressionNode)
				}
			}

			var argNamedArgs []ast.NamedArgumentNode
			if !args[4].IsUndefined() {
				argNamedArgsTuple := args[4].MustReference().(*value.ArrayTuple)
				argNamedArgs = make([]ast.NamedArgumentNode, argNamedArgsTuple.Length())
				for i, el := range *argNamedArgsTuple {
					argNamedArgs[i] = el.MustReference().(ast.NamedArgumentNode)
				}
			}

			var argSpan *position.Span
			if args[5].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[5].Pointer())
			}
			self := ast.NewGenericConstructorCallNode(
				argSpan,
				argClass,
				argTypeArgs,
				argPosArgs,
				argNamedArgs,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"class_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)
			result := value.Ref(self.ClassNode)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"type_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)

			collection := self.TypeArguments
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
		"positional_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)

			collection := self.PositionalArguments
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
		"named_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)

			collection := self.NamedArguments
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
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.GenericConstructorCallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
