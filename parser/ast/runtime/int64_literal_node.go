package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initInt64LiteralNode() {
	c := &value.Int64LiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argValue := (string)(args[1].MustReference().(value.String))
			_, err := value.ParseBigIntWithErr(argValue, 0, value.Int64LiteralNodeFormatErrorClass)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			var argLoc *position.Location
			if args[2].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[2].Pointer())
			}
			self := ast.NewInt64LiteralNode(
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
			self := args[0].MustReference().(*ast.Int64LiteralNode)
			result := value.Ref(value.String(self.Value))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.Int64LiteralNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.Int64LiteralNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.Int64LiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

	vm.Def(
		c,
		"to_int64",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.Int64LiteralNode)
			result, err := value.StrictParseIntWithErr(self.Value, 0, 64, value.Int64LiteralNodeFormatErrorClass)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			return value.Int64(result).ToValue(), value.Undefined
		},
	)

	c = &value.Int64Class.MethodContainer
	vm.Def(
		c,
		"to_ast_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsInt64()
			node := ast.NewInt64LiteralNode(position.ZeroLocation, string(self.ToString()))
			return value.Ref(node), value.Undefined
		},
	)
	vm.Alias(c, "to_ast_expr_node", "to_ast_node")
	vm.Alias(c, "to_ast_pattern_node", "to_ast_node")
	vm.Alias(c, "to_ast_pattern_expr_node", "to_ast_node")
	vm.Alias(c, "to_ast_type_node", "to_ast_node")
}
