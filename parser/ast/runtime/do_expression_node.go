package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initDoExpressionNode() {
	c := &value.DoExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {

			arg0Tuple := args[0].MustReference().(*value.ArrayTuple)
			arg0 := make([]ast.StatementNode, arg0Tuple.Length())
			for i, el := range *arg0Tuple {
				arg0[i] = el.MustReference().(ast.StatementNode)
			}

			arg1Tuple := args[1].MustReference().(*value.ArrayTuple)
			arg1 := make([]*ast.CatchNode, arg1Tuple.Length())
			for i, el := range *arg1Tuple {
				arg1[i] = el.MustReference().(*ast.CatchNode)
			}

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
			self := ast.NewDoExpressionNode(
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
		"body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)

			collection := self.Body
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
		"catches",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)

			collection := self.Catches
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
		"finally_body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)

			collection := self.Finally
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
			self := args[0].MustReference().(*ast.DoExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
