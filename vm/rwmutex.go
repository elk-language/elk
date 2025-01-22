package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::RWMutex
func initRWMutex() {
	// Instance methods
	c := &value.RWMutexClass.MethodContainer
	Def(
		c,
		"lock",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.RWMutex)(args[0].Pointer())
			self.Lock()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"read_lock",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.RWMutex)(args[0].Pointer())
			self.ReadLock()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"unlock",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.RWMutex)(args[0].Pointer())
			if err := self.Unlock(); !err.IsUndefined() {
				return value.Undefined, err
			}

			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"read_unlock",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.RWMutex)(args[0].Pointer())
			if err := self.ReadUnlock(); !err.IsUndefined() {
				return value.Undefined, err
			}

			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.RWMutex)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

}
