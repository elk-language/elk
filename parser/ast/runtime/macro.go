package runtime

import (
	"fmt"

	"github.com/elk-language/elk/ds"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
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
			typechecker := checker.NewMacroChecker()
			compiler := typechecker.CheckMacroExpression(node)
			dl := &typechecker.Errors.DiagnosticList
			if dl.IsFailure() {
				err := value.NewObject(
					value.ObjectWithClass(value.ElkTypeCheckerErrorClass),
					value.ObjectWithInstanceVariablesByName(value.SymbolMap{
						symbol.L_message:     value.String("macro eval checker error").ToValue(),
						symbol.L_diagnostics: (*value.DiagnosticList)(dl).ToValue(),
					}),
				).ToValue()
				return value.Undefined, err
			}

			promise := vm.NewBytecodePromise(thread.ThreadPool(), compiler.Bytecode(), value.GlobalObject.ToValue())

			result, _, err := promise.AwaitSync()
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			return result.ToValue(), value.Undefined
		},
		vm.DefWithParameters(1),
	)
}
