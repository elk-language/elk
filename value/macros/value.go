package macros

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

func initValue(env *types.GlobalEnvironment) {
	astModule := env.StdSubtypeModule(symbol.Elk).MustSubtype(symbol.AST).(*types.Module)
	exprNode := astModule.MustSubtype(symbol.ExpressionNode)
	valueClass := env.StdSubtypeClass(symbol.Value).Singleton()

	types.DefMacro(
		valueClass,
		`Expands to an expression with never type, an endless loop.
Useful in header files for function parameter default value.

Example:

	def foo(a: Bar = Value::never_value!); end`,
		"never_value!",
		nil,
		exprNode,
		func(v *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
			return ast.NewLoopExpressionNode(position.ZeroLocation, nil).ToValue(), value.Undefined
		},
	)
}
