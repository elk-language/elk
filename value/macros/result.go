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
				position.DefaultLocation,
				ast.NewConstantLookupNode(
					position.DefaultLocation,
					ast.NewConstantLookupNode(
						position.DefaultLocation,
						nil,
						ast.NewPublicConstantNode(position.DefaultLocation, "Std"),
					),
					ast.NewPublicConstantNode(position.DefaultLocation, "Result"),
				),
				[]ast.PatternNode{
					ast.NewSymbolKeyValuePatternNode(
						position.DefaultLocation,
						"value",
						ast.NewAsPatternNode(
							position.DefaultLocation,
							ast.NewUnaryExpressionNode(
								position.DefaultLocation,
								token.New(
									position.DefaultLocation,
									token.NOT_EQUAL,
								),
								ast.NewNilLiteralNode(position.DefaultLocation),
							),
							ast.NewPublicIdentifierNode(position.DefaultLocation, "value"),
						),
					),
					ast.NewSymbolKeyValuePatternNode(
						position.DefaultLocation,
						"err",
						ast.NewAsPatternNode(
							position.DefaultLocation,
							ast.NewNilLiteralNode(position.DefaultLocation),
							ast.NewPublicIdentifierNode(position.DefaultLocation, "err"),
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
				position.DefaultLocation,
				ast.NewConstantLookupNode(
					position.DefaultLocation,
					ast.NewConstantLookupNode(
						position.DefaultLocation,
						nil,
						ast.NewPublicConstantNode(position.DefaultLocation, "Std"),
					),
					ast.NewPublicConstantNode(position.DefaultLocation, "Result"),
				),
				[]ast.PatternNode{
					ast.NewSymbolKeyValuePatternNode(
						position.DefaultLocation,
						"err",
						ast.NewAsPatternNode(
							position.DefaultLocation,
							ast.NewUnaryExpressionNode(
								position.DefaultLocation,
								token.New(
									position.DefaultLocation,
									token.NOT_EQUAL,
								),
								ast.NewNilLiteralNode(position.DefaultLocation),
							),
							ast.NewPublicIdentifierNode(position.DefaultLocation, "err"),
						),
					),
					ast.NewSymbolKeyValuePatternNode(
						position.DefaultLocation,
						"value",
						ast.NewAsPatternNode(
							position.DefaultLocation,
							ast.NewNilLiteralNode(position.DefaultLocation),
							ast.NewPublicIdentifierNode(position.DefaultLocation, "value"),
						),
					),
				},
			)

			return value.Ref(val), value.Undefined
		},
	)
}
