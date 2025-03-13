package ast

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
			arg0 := args[0].MustReference().(ast.PatternNode)
			arg1 := args[1].MustReference().(ast.ExpressionNode)

			arg2Tuple := args[2].MustReference().(*value.ArrayTuple)
			arg2 := make([]ast.StatementNode, arg2Tuple.Length())
			for i, el := range *arg2Tuple {
				arg2[i] = el.MustReference().(ast.StatementNode)
			}

			var argSpan *position.Span
			if args[3].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[3].Pointer())
			}
			self := ast.NewForInExpressionNode(
				argSpan,
				arg0,
				arg1,
				arg2,
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

}
