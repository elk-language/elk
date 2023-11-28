package vm

import (
	"github.com/elk-language/elk/value"
)

func init() {
	DefineMethodWithOptions(
		value.TimeClass.SingletonClass().Methods,
		"now",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return value.TimeNow(), nil
		},
	)
	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"format",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			arg := args[1]
			switch a := arg.(type) {
			case value.String:
				result, err := self.Format(string(a))
				if err != nil {
					return nil, err
				}
				return value.String(result), nil
			default:
				return nil, value.Errorf(
					value.ArgumentErrorClass,
					"expected a format string, got: %s",
					arg.Inspect(),
				)
			}
		},
		NativeMethodWithStringParameters("format_string"),
	)
	value.TimeClass.DefineAliasString("strftime", "format")

	DefineMethodWithOptions(
		value.TimeClass.Methods,
		"zone",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(value.Time)
			return self.Zone(), nil
		},
	)
	value.TimeClass.DefineAliasString("timezone", "zone")
}
