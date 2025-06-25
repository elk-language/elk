package runtime

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initUInt8LiteralNode() {
	c := &value.UInt8LiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argValue := (string)(args[1].MustReference().(value.String))
			_, err := value.ParseBigIntWithErr(argValue, 0, value.UInt8LiteralNodeFormatErrorClass)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			var argLoc *position.Location
			if args[2].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[2].Pointer())
			}
			self := ast.NewUInt8LiteralNode(
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
			self := args[0].MustReference().(*ast.UInt8LiteralNode)
			result := value.Ref(value.String(self.Value))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UInt8LiteralNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UInt8LiteralNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UInt8LiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

	vm.Def(
		c,
		"to_uint8",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.UInt8LiteralNode)
			result, err := value.StrictParseUintWithErr(self.Value, 0, 8, value.UInt8LiteralNodeFormatErrorClass)
			if !err.IsUndefined() {
				return value.Undefined, err
			}

			return value.UInt8(result).ToValue(), value.Undefined
		},
	)

	c = &value.UInt8Class.MethodContainer
	vm.Def(
		c,
		"to_ast_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsUInt8()
			node := ast.NewUInt8LiteralNode(position.ZeroLocation, string(self.ToString()))
			return value.Ref(node), value.Undefined
		},
	)
	vm.Alias(c, "to_ast_expr_node", "to_ast_node")
	vm.Alias(c, "to_ast_pattern_node", "to_ast_node")
	vm.Alias(c, "to_ast_type_node", "to_ast_node")
}
