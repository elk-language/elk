package vm

import (
	"runtime"

	"github.com/elk-language/elk/value"
)

// ::Std::Runtime
func initRuntime() {
	c := &value.RuntimeModule.SingletonClass().MethodContainer
	Def(
		c,
		"gc",
		func(v *VM, args []value.Value) (value.Value, value.Value) {
			runtime.GC()
			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)

}
