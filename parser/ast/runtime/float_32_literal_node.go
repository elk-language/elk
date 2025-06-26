package runtime

import (
	"strconv"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initFloat32LiteralNode() {
	c := &value.Float32LiteralNodeClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			argValue := (string)(args[1].MustReference().(value.String))
			_, err := strconv.ParseFloat(argValue, 32)
			if err != nil {
				return value.Undefined,
					value.Ref(
						value.NewError(
							value.Float32LiteralNodeFormatErrorClass,
							err.Error(),
						),
					)
			}

			var argLoc *position.Location
			if args[2].IsUndefined() {
				argLoc = position.ZeroLocation
			} else {
				argLoc = (*position.Location)(args[2].Pointer())
			}
			self := ast.NewFloat32LiteralNode(
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
			self := args[0].MustReference().(*ast.Float32LiteralNode)
			result := value.Ref(value.String(self.Value))
			return result, value.Undefined

		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.Float32LiteralNode)
			result := value.Ref((*value.Location)(self.Location()))
			return result, value.Undefined

		},
	)
	vm.Def(
		c,
		"==",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.Float32LiteralNode)
			other := args[1]
			return value.ToElkBool(self.Equal(other)), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.Float32LiteralNode)
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

	vm.Def(
		c,
		"to_float32",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].MustReference().(*ast.Float32LiteralNode)
			result, err := strconv.ParseFloat(self.Value, 32)
			if err != nil {
				return value.Undefined,
					value.Ref(
						value.NewError(
							value.FloatLiteralNodeFormatErrorClass,
							err.Error(),
						),
					)
			}

			return value.Float32(result).ToValue(), value.Undefined
		},
	)

	c = &value.Float32Class.MethodContainer
	vm.Def(
		c,
		"to_ast_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsFloat32()
			node := ast.NewFloat32LiteralNode(position.ZeroLocation, string(self.ToString()))
			return value.Ref(node), value.Undefined
		},
	)
	vm.Alias(c, "to_ast_expr_node", "to_ast_node")
	vm.Alias(c, "to_ast_pattern_node", "to_ast_node")
	vm.Alias(c, "to_ast_type_node", "to_ast_node")
}
