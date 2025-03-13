package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initSetterDeclarationNode() {
	c := &value.SetterDeclarationNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {

			argDocComment := string(args[0].MustReference().(value.String))
			argEntriesTuple := args[1].MustReference().(*value.ArrayTuple)
			argEntries := make([]ast.ParameterNode, argEntriesTuple.Length())
			for i, el := range *argEntriesTuple {
				argEntries[i] = el.MustReference().(ast.ParameterNode)
			}

			var argSpan *position.Span
			if args[2].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[2].Pointer())
			}
			self := ast.NewSetterDeclarationNode(
				argSpan,
				argDocComment,
				argEntries,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"span",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SetterDeclarationNode)
			result := value.Ref((value.String)(self.DocComment()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"entries",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.SetterDeclarationNode)

			collection := self.Entries
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
			self := args[0].MustReference().(*ast.SetterDeclarationNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
