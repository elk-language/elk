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
	Def(
		c,
		"unwrap",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			self := args[0].AsResult()
			if self.Ok() {
				return self.Value(), value.Undefined
			}
			return value.Undefined, self.Err()
		},
	)
	Alias(c, "or_throw", "unwrap")

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
	Def(
		c,
		"merge",
		func(_ *Thread, args []value.Value) (value.Value, value.Value) {
			return args[1], value.Undefined
		},
		DefWithParameters(1),
	)
}
