package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initReceiverlessMacroCallNode() {
	c := &value.ReceiverlessMacroCallNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argMethodName := args[1].MustReference().(ast.IdentifierNode)

			var argKind ast.MacroKind
			if !args[2].IsUndefined() {
				argKind = ast.MacroKind(args[2].AsUInt8())
			}

			var argPositionalArguments []ast.ExpressionNode
			if !args[3].IsUndefined() {
				argPositionalArgumentsTuple := args[3].MustReference().(*value.ArrayTuple)
				argPositionalArguments = make([]ast.ExpressionNode, argPositionalArgumentsTuple.Length())
				for i, el := range *argPositionalArgumentsTuple {
					argPositionalArguments[i] = el.MustReference().(ast.ExpressionNode)
				}
			}

			var argNamedArguments []ast.NamedArgumentNode
			if !args[4].IsUndefined() {
				argNamedArgumentsTuple := args[4].MustReference().(*value.ArrayTuple)
				argNamedArguments = make([]ast.NamedArgumentNode, argNamedArgumentsTuple.Length())
				for i, el := range *argNamedArgumentsTuple {
					argNamedArguments[i] = el.MustReference().(ast.NamedArgumentNode)
				}
			}

			var argLoc *position.Location
			if args[5].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[5].Pointer())
			}
			self := ast.NewReceiverlessMacroCallNode(
				argLoc,
				argKind,
				argMethodName,
				argPositionalArguments,
				argNamedArguments,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(5),
	)

	vm.Def(
		c,
		"macro_name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMacroCallNode)
			result := value.Ref(self.MacroName)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"positional_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMacroCallNode)

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
		"kind",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMacroCallNode)
			result := value.UInt8(self.Kind)
			return result.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"named_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMacroCallNode)

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
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMacroCallNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMacroCallNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ReceiverlessMacroCallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
