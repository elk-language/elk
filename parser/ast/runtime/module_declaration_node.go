package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initModuleDeclarationNode() {
	c := &value.ModuleDeclarationNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			argConstant := args[1].MustReference().(ast.ExpressionNode)

			var argBody []ast.StatementNode
			if !args[2].IsUndefined() {
				argBodyTuple := args[2].AsReference().(value.ArrayTuple)
				argBody = value.TransformArrayTupleIntoNativeArrayTuple(argBodyTuple, func(v value.Value) ast.StatementNode {
					return v.AsReference().(ast.StatementNode)
				}).ToSlice()
			}

			var argDocComment string
			if !args[3].IsUndefined() {
				argDocComment = string(args[3].MustReference().(value.String))
			}

			var argLoc *position.Location
			if args[4].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[4].Pointer())
			}
			self := ast.NewModuleDeclarationNode(
				argLoc,
				argDocComment,
				argConstant,
				argBody,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(4),
	)

	vm.Def(
		c,
		"constant",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModuleDeclarationNode)
			if self.Constant == nil {
				return value.Nil, value.Undefined
			}

			result := value.Ref(self.Constant)
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"body",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModuleDeclarationNode)
			entries := value.CastNativeArrayTuplePtr(&self.Body)
			return entries.ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"doc_comment",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModuleDeclarationNode)
			result := value.Ref((value.String)(self.DocComment()))
			return result, value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModuleDeclarationNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModuleDeclarationNode)
			other := args[1]
			return value.BoolVal(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.ModuleDeclarationNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)
}
