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
}
