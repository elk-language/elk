package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Result
func initResult() {
	// Instance methods
	c := &value.ResultClass.MethodContainer
	Def(
		c,
		"value",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsResult()
			return self.Value(), value.Undefined
		},
	)
	Def(
		c,
		"err",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsResult()
			return self.Err(), value.Undefined
		},
	)
	Def(
		c,
		"ok",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsResult()
			return value.BoolVal(self.Ok()), value.Undefined
		},
	)

	// Singleton methods
	c = &value.ResultClass.SingletonClass().MethodContainer
	Def(
		c,
		"ok",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			val := args[1]
			return value.MakeOkResult(val).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"err",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			val := args[1]
			return value.MakeErrResult(val).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
}
