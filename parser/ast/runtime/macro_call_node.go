package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initMacroCallNode() {
	c := &value.MacroCallNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argReceiver := args[1].MustReference().(ast.ExpressionNode)
			argMethodName := (string)(args[2].MustReference().(value.String))

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
			self := ast.NewMacroCallNode(
				argLoc,
				argReceiver,
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
		"receiver",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MacroCallNode)
			result := value.Ref(self.Receiver)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"macro_name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MacroCallNode)
			result := value.Ref(value.String(self.MacroName))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"positional_arguments",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MacroCallNode)

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
			self := args[0].MustReference().(*ast.MacroCallNode)

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
			self := args[0].MustReference().(*ast.MacroCallNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MacroCallNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MacroCallNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
