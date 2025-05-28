package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initMethodCallNode() {
	c := &value.MethodCallNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argReceiver := args[1].MustReference().(ast.ExpressionNode)
			argOp := args[2].MustReference().(*token.Token)
			argMethodName := args[3].MustReference().(ast.IdentifierNode)

			var argPositionalArguments []ast.ExpressionNode
			if !args[4].IsUndefined() {
				argPositionalArgumentsTuple := args[4].MustReference().(*value.ArrayTuple)
				argPositionalArguments = make([]ast.ExpressionNode, argPositionalArgumentsTuple.Length())
				for i, el := range *argPositionalArgumentsTuple {
					argPositionalArguments[i] = el.MustReference().(ast.ExpressionNode)
				}
			}

			var argNamedArguments []ast.NamedArgumentNode
			if !args[5].IsUndefined() {
				argNamedArgumentsTuple := args[5].MustReference().(*value.ArrayTuple)
				argNamedArguments = make([]ast.NamedArgumentNode, argNamedArgumentsTuple.Length())
				for i, el := range *argNamedArgumentsTuple {
					argNamedArguments[i] = el.MustReference().(ast.NamedArgumentNode)
				}
			}

			var argLoc *position.Location
			if args[6].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[6].Pointer())
			}
			self := ast.NewMethodCallNode(
				argLoc,
				argReceiver,
				argOp,
				argMethodName,
				argPositionalArguments,
				argNamedArguments,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(6),
	)

	vm.Def(
		c,
		"receiver",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodCallNode)
			result := value.Ref(self.Receiver)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"op",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodCallNode)
			result := value.Ref(self.Op)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"method_name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodCallNode)
			result := value.Ref(self.MethodName)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"positional_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodCallNode)

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
			self := args[0].MustReference().(*ast.MethodCallNode)

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
		"tail_call",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodCallNode)
			result := value.ToElkBool(self.TailCall)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodCallNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodCallNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodCallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
