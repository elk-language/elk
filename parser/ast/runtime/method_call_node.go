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
			argReceiver := args[0].MustReference().(ast.ExpressionNode)
			argOp := args[1].MustReference().(*token.Token)
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

			var argSpan *position.Span
			if args[5].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[5].Pointer())
			}
			self := ast.NewMethodCallNode(
				argSpan,
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
			result := value.Ref(value.String(self.MethodName))
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
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.MethodCallNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
