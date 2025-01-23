package vm

import (
	"os"
	"runtime/pprof"

	"github.com/elk-language/elk/value"
)

// ::Std::Debug
func initDebug() {
	c := &value.DebugModule.SingletonClass().MethodContainer
	Def(
		c,
		"start_cpu_profile",
		func(v *VM, args []value.Value) (value.Value, value.Value) {
			filePath := args[1].MustReference().(value.String)
			f, err := os.Create(string(filePath))
			if err != nil {
				return value.Undefined, value.Ref(value.NewError(value.FileSystemErrorClass, err.Error()))
			}
			pprof.StartCPUProfile(f)
			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"stop_cpu_profile",
		func(v *VM, args []value.Value) (value.Value, value.Value) {
			pprof.StopCPUProfile()
			f, err := os.Create("heap.prof")
			if err != nil {
				return value.Undefined, value.Ref(value.NewError(value.FileSystemErrorClass, err.Error()))
			}
			pprof.WriteHeapProfile(f)
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"inspect_value_stack",
		func(v *VM, args []value.Value) (value.Value, value.Value) {
			v.InspectValueStack()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"inspect_call_stack",
		func(v *VM, args []value.Value) (value.Value, value.Value) {
			v.InspectCallStack()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"stack_trace",
		func(v *VM, args []value.Value) (value.Value, value.Value) {
			return value.Ref(v.BuildStackTrace()), value.Undefined
		},
	)
}
