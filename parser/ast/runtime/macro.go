package runtime

import (
	"fmt"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/macros"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func initMacro() {
	// Std::Macro
	c := value.MacroModule.SingletonClass()

	vm.Def(
		c,
		"expr",
		func(vm *vm.Thread, args []value.Value) (value.Value, value.Value) {
			node := args[1].AsReference()
			var expr ast.ExpressionNode

			switch node := node.(type) {
			case ast.StatementNode:
				expr = ast.NewDoExpressionNode(node.Location(), []ast.StatementNode{node}, nil, nil)
			case value.ArrayTuple:
				body := ds.MapSeq2(
					node.Elements(),
					func(_ int, v value.Value) ast.StatementNode {
						return v.MustReference().(ast.StatementNode)
					},
				)
				expr = ast.NewDoExpressionNode(position.ZeroLocation, body, nil, nil)
			default:
				panic(fmt.Sprintf("invalid node argument to expr: %T", node))
			}

			return value.Ref(expr), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"unhygienic",
		func(vm *vm.Thread, args []value.Value) (value.Value, value.Value) {
			node := args[1].AsReference().(ast.Node)
			result := ast.NewUnhygienicNode(
				node.Location(),
				node,
			)
			return value.Ref(result), value.Undefined
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"eval_node",
		func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
			node := args[1].AsReference().(ast.ExpressionNode)
			return checker.EvalNode(thread, node)
		},
		vm.DefWithParameters(1),
	)

	vm.Def(
		c,
		"append_expression",
		func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
			context := args[1].AsReference().(*macros.Context)
			node := args[2].AsReference().(ast.ExpressionNode)
			kindArg := args[3].AsSymbol()

			var kind types.MacroAppendExpressionKind
			switch kindArg.Id {
			case symbol.L_end_of_file:
				kind = types.MacroAppendExpressionEndOfFileKind
			case symbol.L_end_of_namespace:
				kind = types.MacroAppendExpressionEndOfNamespaceKind
			default:
				return value.Undefined, value.NewError(
					value.ArgumentErrorClass,
					fmt.Sprintf(
						"invalid `kind`, expected `:end_of_file` or `:end_of_namespace`, got: `%s`",
						kindArg.Inspect(),
					),
				).ToValue()
			}

			var key string
			if args[4].IsNotUndefined() {
				key = string(args[4].AsString())
			}
			context.TypeChecker.RegisterMacroAppendExpression(nil, key, kind, node)
			return value.Nil, value.Undefined
		},
		vm.DefWithParameters(4),
	)
}
