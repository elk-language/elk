package checker

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
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
		`Evaluates the condition and returns the given "then_node" if it's result is truthy.
Otherwise returns the "else_node" if it is present or a "nil" node.

Example:

	macro_if! Env.get('OS') == "unix", do
		def some_unix_method
		end
	end
`,
		"macro_if!",
		[]*types.Parameter{
			types.NewParameter(
				value.ToSymbol("condition"),
				exprType,
				types.NormalParameterKind,
				false,
			),
			types.NewParameter(
				value.ToSymbol("then_node"),
				exprType,
				types.NormalParameterKind,
				false,
			),
			types.NewParameter(
				value.ToSymbol("else_node"),
				types.NewNilable(exprType),
				types.DefaultValueParameterKind,
				false,
			),
		},
		exprType,
		func(v *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
			conditionNode := args[1].AsReference().(ast.ExpressionNode)
			thenNode := args[2]

			var elseNode value.Value
			if args[3].IsNotUndefined() {
				elseNode = args[3]
			} else {
				elseNode = ast.NewNilLiteralNode(position.ZeroLocation).ToValue()
			}

			conditionResult, err := EvalNode(v, conditionNode)
			if err.IsNotUndefined() {
				return value.Undefined, err
			}

			if value.Truthy(conditionResult) {
				return thenNode, value.Undefined
			}

			return elseNode, value.Undefined
		},
	)
}
