package vm

import (
	"sync"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Std::Once
func initOnce() {
	// Singleton methods
	c := &value.OnceClass.SingletonClass().MethodContainer
	Def(
		c,
		"memo",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return OnceMemo(args[1]).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"fn",
		func(vm *Thread, args []value.Value) (value.Value, value.Value) {
			return OnceFn(args[1]).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)

	// Instance methods
	c = &value.OnceClass.MethodContainer
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
	once.Native().Do(func() {
		_, err = vm.CallCallable(fn)
	})

	return err
}

func OnceMemo(fn value.Value) *NativeClosure {
	var once sync.Once
	var memoResult value.Value
	var memoErr value.Value

	return NewNativeClosure(
		func(vm *Thread, args []value.Value) (returnVal value.Value, err value.Value) {
			once.Do(func() {
				memoResult, memoErr = vm.CallCallable(fn)
			})

			return memoResult, memoErr
		},
		0,
		position.ZeroLocation,
	)
}

func OnceFn(fn value.Value) *NativeClosure {
	var once sync.Once

	return NewNativeClosure(
		func(vm *Thread, args []value.Value) (returnVal value.Value, err value.Value) {
			once.Do(func() {
				vm.CallCallable(fn)
			})

			return value.Nil, value.Undefined
		},
		0,
		position.ZeroLocation,
	)
}
