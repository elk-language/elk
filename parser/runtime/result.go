package runtime

import (
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Std::Elk::Parser::Result
func initResult() {
	// Instance methods
	c := &value.ElkParserResultClass.MethodContainer
	vm.Def(
		c,
		"ast",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*parser.Result)(args[0].Pointer())
			return value.Ref(self.AST), value.Undefined
		},
	)
	vm.Def(
		c,
		"diagnostics",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*parser.Result)(args[0].Pointer())
			return value.Ref((*value.DiagnosticList)(&self.Diagnostics)), value.Undefined
		},
	)
}
