package runtime

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

// Std::Elk::Type::Checker::Error
func initError() {
	// Instance methods
	c := &value.ElkTypeCheckerErrorClass.MethodContainer
	vm.Def(
		c,
		"#init",
		func(_ *vm.Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsReference().(*value.Object)
			message := args[1]
			diagnostics := args[2]

			self.SetMessage(message)
			self.SetInstanceVariable(symbol.L_diagnostics, diagnostics)

			return value.Ref(self), value.Undefined
		},
		vm.DefWithParameters(2),
	)

	vm.Def(
		c,
		"diagnostics",
		func(thread *vm.Thread, args []value.Value) (returnVal value.Value, err value.Value) {
			self := args[0].AsReference().(*value.Object)
			return self.GetInstanceVariable(symbol.L_diagnostics), value.Undefined
		},
	)
}
