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
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Result)(args[0].Pointer())
			return self.Value(), value.Undefined
		},
	)
	Def(
		c,
		"err",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Result)(args[0].Pointer())
			return self.Err(), value.Undefined
		},
	)
	Def(
		c,
		"ok",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Result)(args[0].Pointer())
			return value.ToElkBool(self.Ok()), value.Undefined
		},
	)

	// Singleton methods
	c = &value.ResultClass.SingletonClass().MethodContainer
	Def(
		c,
		"ok",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			val := args[1]
			return value.Ref(value.NewOkResult(val)), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"err",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			val := args[1]
			return value.Ref(value.NewErrResult(val)), value.Undefined
		},
		DefWithParameters(1),
	)
}
