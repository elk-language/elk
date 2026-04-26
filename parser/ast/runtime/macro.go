package runtime

import (
	"fmt"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func initMacro() {
	// Std::Macro
	c := &value.MacroModule.SingletonClass().MethodContainer

	vm.Def(
		c,
		"expr",
		func(vm *vm.Thread, args []value.Value) (value.Value, value.Value) {
			node := args[1].AsReference()
			var expr ast.ExpressionNode

			switch node := node.(type) {
			case ast.StatementNode:
				expr = ast.NewDoExpressionNode(node.Location(), []ast.StatementNode{node}, nil, nil)
			case *value.ArrayTupleOfValue:
				body := ds.MapSlice(
					*node,
					func(v value.Value) ast.StatementNode {
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
}
