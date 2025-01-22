package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::ROMutex
func initROMutex() {
	// Instance methods
	c := &value.ROMutexClass.MethodContainer
	Def(
		c,
		"#init",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.ROMutex)(args[0].Pointer())
			arg := args[1]
			if arg.IsUndefined() {
				self.RWMutex = value.NewRWMutex()
			} else {
				self.RWMutex = (*value.RWMutex)(arg.Pointer())
			}

			return args[0], value.Undefined
		},
		DefWithParameters(1),
		DefWithOptionalParameters(1),
	)
	Def(
		c,
		"lock",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.ROMutex)(args[0].Pointer())
			self.Lock()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"unlock",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.ROMutex)(args[0].Pointer())
			if err := self.Unlock(); !err.IsUndefined() {
				return value.Undefined, err
			}

			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"rwmutex",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.ROMutex)(args[0].Pointer())
			return value.Ref(self.RWMutex), value.Undefined
		},
	)
	Def(
		c,
		"inspect",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := (*value.ROMutex)(args[0].Pointer())
			return value.Ref(value.String(self.Inspect())), value.Undefined
		},
	)

}
