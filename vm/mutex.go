package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Mutex
func initMutex() {
	// Instance methods
	c := &value.MutexClass.MethodContainer
	Def(
		c,
		"lock",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Mutex)(args[0].Pointer())
			self.Lock()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"unlock",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Mutex)(args[0].Pointer())
			if err := self.Unlock(); !err.IsUndefined() {
				return value.Undefined, err
			}

			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.Mutex)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

}
