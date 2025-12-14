package vm

import (
	"github.com/elk-language/elk/value"
)

// Std::Once
func initOnce() {
	// Instance methods
	c := &value.OnceClass.MethodContainer
	Def(
		c,
		"call",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			self := (*value.Once)(args[0].Pointer())
			err := OnceDo(vm, self, args[1])
			if !err.IsUndefined() {
				return value.Undefined, err
			}
			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)

}

func OnceDo(vm *Thread, once *value.Once, fn value.Value) (err value.Value) {
	once.Native.Do(func() {
		_, err = vm.CallCallable(fn)
	})

	return err
}
