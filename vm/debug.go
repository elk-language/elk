package vm

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/elk-language/elk/value"
)

// ::Std::Debug
func initDebug() {
	c := value.DebugModule.SingletonClass()
	Def(
		c,
		"start_cpu_profile",
		func(v *Thread, args []value.Value) (value.Value, value.Value) {
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
		func(v *Thread, args []value.Value) (value.Value, value.Value) {
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
		func(v *Thread, args []value.Value) (value.Value, value.Value) {
			v.InspectValueStack()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"inspect_call_stack",
		func(v *Thread, args []value.Value) (value.Value, value.Value) {
			v.InspectCallStack()
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"inspect_inheritance",
		func(v *Thread, args []value.Value) (value.Value, value.Value) {
			klass := args[1].AsReference().(*value.Class)
			return value.Ref(value.String(klass.InspectInheritance())), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"stack_trace",
		func(v *Thread, args []value.Value) (value.Value, value.Value) {
			return value.Ref(v.BuildStackTrace()), value.Undefined
		},
	)
	Def(
		c,
		"inspect_native_type",
		func(v *Thread, args []value.Value) (value.Value, value.Value) {
			val := args[1].ToInterface()
			return value.String(fmt.Sprintf("%T", val)).ToValue(), value.Undefined
		},
		DefWithParameters(1),
	)
	Def(
		c,
		"inspect_current_bytecode",
		func(v *Thread, args []value.Value) (value.Value, value.Value) {
			if v.bytecode == nil {
				return value.Undefined, value.NewError(
					value.TypeErrorClass,
					"Cannot inspect bytecode of a native method",
				).ToValue()
			}

			v.bytecode.Disassemble(os.Stdout)
			return value.Nil, value.Undefined
		},
	)
	Def(
		c,
		"inspect_bytecode",
		func(v *Thread, args []value.Value) (value.Value, value.Value) {
			method := args[1]
			switch method := method.AsReference().(type) {
			case *BytecodeFunction:
				method.Disassemble(v.Stdout)
			case *BytecodeClosure:
				method.Bytecode.Disassemble(v.Stdout)
			default:
				return value.Undefined, value.NewError(
					value.TypeErrorClass,
					"Given method has no bytecode to inspect",
				).ToValue()
			}

			return value.Nil, value.Undefined
		},
		DefWithParameters(1),
	)
}
