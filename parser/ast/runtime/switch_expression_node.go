package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initSwitchExpressionNode() {
	c := &value.SwitchExpressionNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argValue := args[1].MustReference().(ast.ExpressionNode)

			argCasesTuple := args[2].MustReference().(*value.ArrayTuple)
			argCases := make([]*ast.CaseNode, argCasesTuple.Length())
			for i, el := range *argCasesTuple {
				argCases[i] = el.MustReference().(*ast.CaseNode)
			}

			var argElseBody []ast.StatementNode
			if !args[3].IsUndefined() {
				argElseBodyTuple := args[3].MustReference().(*value.ArrayTuple)
				argElseBody = make([]ast.StatementNode, argElseBodyTuple.Length())
				for i, el := range *argElseBodyTuple {
					argElseBody[i] = el.MustReference().(ast.StatementNode)
				}
			}

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.DefaultLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewSwitchExpressionNode(
				argLoc,
				argValue,
				argCases,
				argElseBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"value",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchExpressionNode)
			result := value.Ref(self.Value)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"cases",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchExpressionNode)

			collection := self.Cases
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
		"else_body",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchExpressionNode)

			collection := self.ElseBody
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
			self := args[0].MustReference().(*ast.SwitchExpressionNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchExpressionNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SwitchExpressionNode)
			result := value.Ref(value.String(self.String()))
			return result, value.Undefined
		},
	)

}
