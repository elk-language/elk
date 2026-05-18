package macros

import (
	"slices"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func initWaitGroup(env *types.GlobalEnvironment) {
	astModule := env.NamesToNamespace(symbol.Std, symbol.Elk, symbol.AST)
	exprNode := astModule.MustSubtype(symbol.ExpressionNode)
	goNode := astModule.MustSubtype(value.ToSymbol("GoExpressionNode"))
	waitGroupClass := env.NamesToNamespace(symbol.Std, symbol.Sync, symbol.WaitGroup).Singleton()

	types.DefMacro(
		waitGroupClass,
		`Expands to an expression with never type, an endless loop.
Useful in header files for function parameter default value.

Example:

	using Std::Sync::WaitGroup
	using Std::Sync::WaitGroup::spawn!

	wg := WaitGroup()

	spawn! wg, go
		sleep 5.seconds
	  println "thread 1"
	end

	spawn! wg, go
		sleep 2.seconds
	  println "thread 2"
	end

	wg.wait
`,
		"spawn!",
		[]*types.Parameter{
			types.NewParameter(
				value.ToSymbol("wait_group"),
				exprNode,
				types.NormalParameterKind,
				false,
			),
			types.NewParameter(
				value.ToSymbol("go_expr"),
				goNode,
				types.NormalParameterKind,
				false,
			),
		},
		exprNode,
		func(v *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
			waitGroupExpr := args[1].AsReference().(ast.ExpressionNode)
			goExpr := args[2].AsReference().(*ast.GoExpressionNode)

			result := ast.NewDoExpressionNode(
				position.ZeroLocation,
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						position.ZeroLocation,
						ast.NewMethodCallNode(
							position.ZeroLocation,
							waitGroupExpr,
							token.New(position.ZeroLocation, token.DOT),
							ast.NewPublicIdentifierNode(position.ZeroLocation, "start"),
							nil,
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						position.ZeroLocation,
						ast.NewGoExpressionNode(
							position.ZeroLocation,
							append(
								slices.Clone(goExpr.Body),
								ast.NewExpressionStatementNode(
									position.ZeroLocation,
									ast.NewMethodCallNode(
										position.ZeroLocation,
										waitGroupExpr,
										token.New(position.ZeroLocation, token.DOT),
										ast.NewPublicIdentifierNode(position.ZeroLocation, "end"),
										nil,
										nil,
									),
								),
							),
						),
					),
				},
				nil,
				nil,
			)

			return ast.Unhygienic(result).ToValue(), value.Undefined
		},
	)
}
