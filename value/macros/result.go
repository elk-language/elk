package macros

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func initResult(env *types.GlobalEnvironment) {
	astModule := env.StdSubtypeModule(symbol.Elk).MustSubtype(symbol.AST).(*types.Module)
	patternNode := astModule.MustSubtype(symbol.PatternNode)
	result := env.StdSubtypeClass(symbol.Result).Singleton()

	types.DefMacro(
		result,
		`Expands to a pattern that matches successful Result values

Example:

	switch divide(10, 2)
	case Result::ok!
		puts "ok: #value" #> value: 5
	case Result::err!
		puts "err: #err"
	end`,
		"ok!",
		nil,
		patternNode,
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			val := ast.NewObjectPatternNode(
				position.ZeroLocation,
				ast.NewConstantLookupNode(
					position.ZeroLocation,
					ast.NewConstantLookupNode(
						position.ZeroLocation,
						nil,
						ast.NewPublicConstantNode(position.ZeroLocation, "Std"),
					),
					ast.NewPublicConstantNode(position.ZeroLocation, "Result"),
				),
				[]ast.PatternNode{
					ast.NewSymbolKeyValuePatternNode(
						position.ZeroLocation,
						"value",
						ast.NewAsPatternNode(
							position.ZeroLocation,
							ast.NewUnaryExpressionNode(
								position.ZeroLocation,
								token.New(
									position.ZeroLocation,
									token.NOT_EQUAL,
								),
								ast.NewNilLiteralNode(position.ZeroLocation),
							),
							ast.NewPublicIdentifierNode(position.ZeroLocation, "value"),
						),
					),
					ast.NewSymbolKeyValuePatternNode(
						position.ZeroLocation,
						"err",
						ast.NewAsPatternNode(
							position.ZeroLocation,
							ast.NewNilLiteralNode(position.ZeroLocation),
							ast.NewPublicIdentifierNode(position.ZeroLocation, "err"),
						),
					),
				},
			)

			return value.Ref(val), value.Undefined
		},
	)

	types.DefMacro(
		result,
		`Expands to a pattern that matches failure Result values

Example:

	switch divide(10, 0)
	case Result::ok!
		puts "ok: #value"
	case Result::err!
		puts "err: #err" #> err: "division by zero"
	end`,
		"err!",
		nil,
		patternNode,
		func(v *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
			val := ast.NewObjectPatternNode(
				position.ZeroLocation,
				ast.NewConstantLookupNode(
					position.ZeroLocation,
					ast.NewConstantLookupNode(
						position.ZeroLocation,
						nil,
						ast.NewPublicConstantNode(position.ZeroLocation, "Std"),
					),
					ast.NewPublicConstantNode(position.ZeroLocation, "Result"),
				),
				[]ast.PatternNode{
					ast.NewSymbolKeyValuePatternNode(
						position.ZeroLocation,
						"err",
						ast.NewAsPatternNode(
							position.ZeroLocation,
							ast.NewUnaryExpressionNode(
								position.ZeroLocation,
								token.New(
									position.ZeroLocation,
									token.NOT_EQUAL,
								),
								ast.NewNilLiteralNode(position.ZeroLocation),
							),
							ast.NewPublicIdentifierNode(position.ZeroLocation, "err"),
						),
					),
					ast.NewSymbolKeyValuePatternNode(
						position.ZeroLocation,
						"value",
						ast.NewAsPatternNode(
							position.ZeroLocation,
							ast.NewNilLiteralNode(position.ZeroLocation),
							ast.NewPublicIdentifierNode(position.ZeroLocation, "value"),
						),
					),
				},
			)

			return value.Ref(val), value.Undefined
		},
	)
}
