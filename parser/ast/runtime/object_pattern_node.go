package ast

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initObjectPatternNode() {
	c := &value.ObjectPatternNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argObjectType := args[0].MustReference().(ast.ComplexConstantNode)

			argAttributesTuple := args[1].MustReference().(*value.ArrayTuple)
			argAttributes := make([]ast.PatternNode, argAttributesTuple.Length())
			for i, el := range *argAttributesTuple {
				argAttributes[i] = el.MustReference().(ast.PatternNode)
			}

			var argSpan *position.Span
			if args[2].IsUndefined() {
				argSpan = position.DefaultSpan
			} else {
				argSpan = (*position.Span)(args[2].Pointer())
			}
			self := ast.NewObjectPatternNode(
				argSpan,
				argObjectType,
				argAttributes,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"object_type",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ObjectPatternNode)
			result := value.Ref(self.ObjectType)
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"attributes",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ObjectPatternNode)

			collection := self.Attributes
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
			self := args[0].MustReference().(*ast.ObjectPatternNode)
			result := value.Ref((*value.Span)(self.Span()))
			return result, value.Undefined

		},
	)

}
