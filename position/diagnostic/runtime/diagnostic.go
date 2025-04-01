package runtime

import (
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/lexer/colorizer"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Std::Diagnostic
func initDiagnostic() {
	c := &value.DiagnosticClass.MethodContainer

	vm.Def(
		c,
		"#init",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Diagnostic)(args[0].Pointer())

			message := string(args[1].AsReference().(value.String))
			location := (*position.Location)(args[2].Pointer())

			var severity diagnostic.Severity
			if !args[3].IsUndefined() {
				severity = diagnostic.Severity(args[3].AsUInt8())
			}

			self.Location = location
			self.Message = message
			self.Severity = severity
			return value.Ref(self), value.Undefined
		},
		vm.DefWithParameters(3),
	)

	vm.Def(
		c,
		"message",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Diagnostic)(args[0].Pointer())
			return value.Ref(value.String(self.Message)), value.Undefined
		},
	)

	vm.Def(
		c,
		"location",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Diagnostic)(args[0].Pointer())
			return value.Ref((*value.Location)(self.Location)), value.Undefined
		},
	)

	vm.Def(
		c,
		"severity",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Diagnostic)(args[0].Pointer())
			return (value.UInt8)(self.Severity).ToValue(), value.Undefined
		},
	)

	vm.Def(
		c,
		"severity_name",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Diagnostic)(args[0].Pointer())
			return value.Ref(value.String(self.Severity.String())), value.Undefined
		},
	)

	vm.Def(
		c,
		"to_string",
		func(_ *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Diagnostic)(args[0].Pointer())
			return value.Ref(value.String(self.String())), value.Undefined
		},
	)

	vm.Def(
		c,
		"to_human_string",
		func(v *vm.VM, args []value.Value) (value.Value, value.Value) {
			self := (*diagnostic.Diagnostic)(args[0].Pointer())

			style := true
			if !args[1].IsUndefined() {
				style = value.Truthy(args[1])
			}

			var colorizer colorizer.Colorizer
			if args[2].IsUndefined() {
				colorizer = lexer.Colorizer{}
			} else if !args[2].IsNil() {
				colorizer = vm.MakeColorizer(v, args[2])
			}

			var result string
			var err error
			if args[3].IsUndefined() {
				result, err = self.HumanString(style, colorizer)
			} else {
				result, err = self.HumanStringWithSource(
					string(args[3].AsReference().(value.String)),
					style,
					colorizer,
				)
			}

			if err != nil {
				return value.Undefined, value.Ref(value.NewError(value.ColorizerErrorClass, err.Error()))
			}
			return value.Ref(value.String(result)), value.Undefined
		},
		vm.DefWithParameters(3),
	)
}
