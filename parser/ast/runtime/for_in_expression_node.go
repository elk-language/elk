package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initForInExpressionNode() {
	c := &value.ForInExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argPattern := args[1].MustReference().(ast.PatternNode)
			argInExpression := args[2].MustReference().(ast.ExpressionNode)

			argThenTuple := args[3].MustReference().(*value.ArrayTuple)
			argThen := make([]ast.StatementNode, argThenTuple.Length())
			for i, el := range *argThenTuple {
				argThen[i] = el.MustReference().(ast.StatementNode)
			}

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.DefaultLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewForInExpressionNode(
				argLoc,
				argPattern,
				argInExpression,
				argThen,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"pattern",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ForInExpressionNode)
			result := value.Ref(self.Pattern)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"in_expression",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ForInExpressionNode)
			result := value.Ref(self.InExpression)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"then_body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ForInExpressionNode)

			collection := self.ThenBody
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
			self := args[0].MustReference().(*ast.ForInExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ForInExpressionNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ForInExpressionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
