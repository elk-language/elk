package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initConstructorCallNode() {
	c := &value.ConstructorCallNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argClassNode := args[1].MustReference().(ast.ComplexConstantNode)

			var argPosArgs []ast.ExpressionNode
			if !args[2].IsUndefined() {
				argPosArgsTuple := args[2].MustReference().(*value.ArrayTuple)
				argPosArgs = make([]ast.ExpressionNode, argPosArgsTuple.Length())
				for i, el := range *argPosArgsTuple {
					argPosArgs[i] = el.MustReference().(ast.ExpressionNode)
				}
			}

			var argNamedArgs []ast.NamedArgumentNode
			if !args[3].IsUndefined() {
				argNamedArgsTuple := args[3].MustReference().(*value.ArrayTuple)
				argNamedArgs = make([]ast.NamedArgumentNode, argNamedArgsTuple.Length())
				for i, el := range *argNamedArgsTuple {
					argNamedArgs[i] = el.MustReference().(ast.NamedArgumentNode)
				}
			}

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewConstructorCallNode(
				argLoc,
				argClassNode,
				argPosArgs,
				argNamedArgs,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"class_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstructorCallNode)
			result := value.Ref(self.ClassNode)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"positional_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstructorCallNode)

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
			self := args[0].MustReference().(*ast.ConstructorCallNode)

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
			self := args[0].MustReference().(*ast.ConstructorCallNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstructorCallNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ConstructorCallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
