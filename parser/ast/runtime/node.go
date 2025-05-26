package runtime

import (
	"fmt"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func initNode() {
	// Std::Elk::AST::Node
	c := &value.NodeMixin.MethodContainer
	vm.Def(
		c,
		"to_ast_node",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			return args[0], value.Undefined
		},
	)
	vm.Def(
		c,
		"traverse",
		func(v *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(ast.Node)
			fn := args[1]

			switch f := fn.SafeAsReference().(type) {
			case *vm.Closure:
				for node := range ast.Iter(self) {
					ok, err := v.CallClosure(f, value.Ref(node))
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					if value.Falsy(ok) {
						return value.False, value.Undefined
					}
				}
			default:
				for node := range ast.Iter(self) {
					ok, err := v.CallMethodByName(symbol.L_call, fn, value.Ref(node))
					if !err.IsUndefined() {
						return value.Undefined, err
					}
					if value.Falsy(ok) {
						return value.False, value.Undefined
					}
				}
			}

			return value.True, value.Undefined
		},
		vm.DefWithParameters(1),
	)
	vm.Def(
		c,
		"iter",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(ast.Node)
			iterator := ast.NewNodeIterator(self)
			return value.Ref(iterator), value.Undefined
		},
	)

	// &Std::Elk::AST::Node
	c = &value.NodeMixin.SingletonClass().MethodContainer
	vm.Def(
		c,
		"expr",
		func(vm *vm.VM, args []value.Value) (value.Value, value.Value) {
			node := args[1].AsReference()
			var expr ast.ExpressionNode

			switch node := node.(type) {
			case ast.StatementNode:
				expr = ast.NewDoExpressionNode(node.Location(), []ast.StatementNode{node}, nil, nil)
			case *value.ArrayTuple:
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

	// Std::Kernel
	c = &value.KernelModule.SingletonClass().MethodContainer
	vm.Def(
		c,
		"#splice",
		func(vm *vm.VM, args []value.Value) (value.Value, value.Value) {
			baseNode := args[1].AsReference().(ast.Node)

			var replacementNodes value.ArrayTuple
			if !args[2].IsUndefined() {
				replacementNodes = *(*value.ArrayTuple)(args[2].Pointer())
			}

			r := ds.MapSlice(
				replacementNodes,
				func(v value.Value) ast.Node {
					return v.AsReference().(ast.Node)
				},
			)
			result := ast.Splice(baseNode, nil, &r)

			return value.Ref(result), value.Undefined
		},
		vm.DefWithParameters(2),
	)
}
