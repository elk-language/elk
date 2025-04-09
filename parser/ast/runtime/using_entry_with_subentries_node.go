package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initUsingEntryWithSubentriesNode() {
	c := &value.UsingEntryWithSubentriesNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argNamespace := args[1].MustReference().(ast.UsingEntryNode)

			argSubentriesTuple := args[2].MustReference().(*value.ArrayTuple)
			argSubentries := make([]ast.UsingSubentryNode, argSubentriesTuple.Length())
			for i, el := range *argSubentriesTuple {
				argSubentries[i] = el.MustReference().(ast.UsingSubentryNode)
			}

			var argLoc *position.Location
			if args[3].IsUndefined() {
				argLoc = position.DefaultLocation
			} else {
				argLoc = (*position.Location)(args[3].Pointer())
			}
			self := ast.NewUsingEntryWithSubentriesNode(
				argLoc,
				argNamespace,
				argSubentries,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"namespace",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UsingEntryWithSubentriesNode)
			result := value.Ref(self.Namespace)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"subentries",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UsingEntryWithSubentriesNode)

			collection := self.Subentries
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
			self := args[0].MustReference().(*ast.UsingEntryWithSubentriesNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UsingEntryWithSubentriesNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UsingEntryWithSubentriesNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

}
