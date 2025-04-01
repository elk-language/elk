package runtime

import (
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Std::Elk::Parser
func initParser() {
	// Singleton methods
	c := &value.ElkParserClass.SingletonClass().MethodContainer
	vm.Def(
		c,
		"parse",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			source := string(args[1].AsReference().(value.String))

			var sourceName string
			if args[2].IsUndefined() {
				sourceName = "<main>"
			} else {
				sourceName = string(args[2].AsReference().(value.String))
			}

			ast, diag := parser.Parse(sourceName, source)

			result := &parser.Result{
				AST:         ast,
				Diagnostics: diag,
			}
			return value.Ref(result), value.Undefined
		},
		vm.DefWithParameters(2),
	)
}
