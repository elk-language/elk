package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initCharLiteralNode() {
	c := &value.CharLiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argValue := (rune)(args[1].AsChar())

			var argLoc *position.Location
			if args[2].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[2].Pointer())
			}
			self := ast.NewCharLiteralNode(
				argLoc,
				argValue,
			)
			return value.Ref(self), value.Undefined

		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"value",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CharLiteralNode)
			result := value.Char(self.Value).ToValue()
			return result, value.Undefined

		},
	)
	vm.Alias(c, "to_char", "value")

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CharLiteralNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CharLiteralNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.CharLiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

	c = &value.CharClass.MethodContainer
	vm.Def(
		c,
		"to_ast_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsChar()
			node := ast.NewCharLiteralNode(position.ZeroLocation, self.Rune())
			return value.Ref(node), value.Undefined
		},
	)
	vm.Alias(c, "to_ast_expr_node", "to_ast_node")
	vm.Alias(c, "to_ast_pattern_node", "to_ast_node")
	vm.Alias(c, "to_ast_pattern_expr_node", "to_ast_node")
	vm.Alias(c, "to_ast_type_node", "to_ast_node")

	vm.Def(
		c,
		"to_ast_ident_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsChar()
			node := ast.NewPublicIdentifierNode(position.ZeroLocation, string(self.Rune()))
			return value.Ref(node), value.Undefined
		},
	)
	vm.Def(
		c,
		"to_ast_const_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsChar()
			node := ast.NewPublicConstantNode(position.ZeroLocation, string(self.Rune()))
			return value.Ref(node), value.Undefined
		},
	)
	vm.Def(
		c,
		"to_ast_ivar_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsChar()
			node := ast.NewPublicInstanceVariableNode(position.ZeroLocation, string(self.Rune()))
			return value.Ref(node), value.Undefined
		},
	)
}
