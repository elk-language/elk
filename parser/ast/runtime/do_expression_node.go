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
			var argBody []ast.StatementNode
			if !args[0].IsUndefined() {
				argBodyTuple := args[0].MustReference().(*value.ArrayTuple)
				argBody = make([]ast.StatementNode, argBodyTuple.Length())
				for i, el := range *argBodyTuple {
					argBody[i] = el.MustReference().(ast.StatementNode)
				}
			}

			var argCatches []*ast.CatchNode
			if !args[1].IsUndefined() {
				argCatchesTuple := args[1].MustReference().(*value.ArrayTuple)
				argCatches = make([]*ast.CatchNode, argCatchesTuple.Length())
				for i, el := range *argCatchesTuple {
					argCatches[i] = el.MustReference().(*ast.CatchNode)
				}
			}

			var argFinally []ast.StatementNode
			if !args[2].IsUndefined() {
				argFinallyTuple := args[2].MustReference().(*value.ArrayTuple)
				argFinally = make([]ast.StatementNode, argFinallyTuple.Length())
				for i, el := range *argFinallyTuple {
					argFinally[i] = el.MustReference().(ast.StatementNode)
				}
			}

			var argSpan *position.Span
			if args[3].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[3].Pointer())
			}
			self := ast.NewDoExpressionNode(
				argSpan,
				argBody,
				argCatches,
				argFinally,
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

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.DoExpressionNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
