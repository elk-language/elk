package checker

import (
	"fmt"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/macros"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func NewGlobalEnvironment() *types.GlobalEnvironment {
	env := macros.NewGlobalEnvironment()

	initMacro(env)

	return env
}

// Std::Macro
func initMacro(env *types.GlobalEnvironment) {
	astType := env.StdSubtypeModule(symbol.Elk).MustSubtype(symbol.AST).(*types.Module)
	exprType := astType.MustSubtype(symbol.ExpressionNode)
	macroType := env.StdSubtypeModule(symbol.Macro)

	types.DefMacro(
		macroType,
		`Evaluates the if condition and returns the "then" body nodes if it's result is truthy.
Otherwise returns the "else" body nodes if they're present or a "nil" node.

Example:

	compile_if!(if Env.get('OS') == "unix"
		def some_unix_method
		end
	end)
`,
		"compile_if!",
		[]*types.Parameter{
			types.NewParameter(
				value.ToSymbol("if_node"),
				types.NewUnion(
					astType.MustSubtypeString("IfExpressionNode"),
					astType.MustSubtypeString("UnlessExpressionNode"),
				),
				types.NormalParameterKind,
				false,
			),
		},
		exprType,
		func(v *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
			node := args[1].AsReference().(ast.ExpressionNode)

			var conditionNode ast.ExpressionNode
			var thenBody []ast.StatementNode
			var elseBody []ast.StatementNode
			switch node := node.(type) {
			case *ast.IfExpressionNode:
				conditionNode = node.Condition
				thenBody = node.ThenBody
				elseBody = node.ElseBody
			case *ast.UnlessExpressionNode:
				conditionNode = node.Condition
				elseBody = node.ThenBody
				thenBody = node.ElseBody
			default:
				panic(fmt.Sprintf("invalid node: %T", node))
			}

			conditionResult, err := EvalNode(v, conditionNode)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			if value.Truthy(conditionResult) {
				return ast.NewDoExpressionNode(node.Location(), thenBody, nil, nil).ToValue(), value.Undefined
			}

			return ast.NewDoExpressionNode(node.Location(), elseBody, nil, nil).ToValue(), value.Undefined
		},
	)
}
