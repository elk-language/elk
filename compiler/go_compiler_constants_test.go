package compiler_test

import (
	"testing"
)

func TestGoGetConstant(t *testing.T) {
	tests := goTestTable{
		"infer namespace in method argument": {
			input: `
				class Color
					const RED: Color = new(255u8, 0u8, 0u8)
					const GREEN: Color = new(0u8, 255u8, 0u8)
					const BLUE: Color = new(0u8, 0u8, 255u8)

					getter red: UInt8
					getter green: UInt8
					getter blue: UInt8
					init(@red, @green, @blue); end
				end

				def baz(col: Color)
					println "col: #col"
				end
				baz(.::RED)
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym1 = value.ToSymbol("#init")
var cc_main_1 = &vm.CallCache{}
var sym2 = value.ToSymbol("RED")
var const1 value.Value // constant: Color::RED, loc: <main>:3:6
var cc_main_2 = &vm.CallCache{}
var sym3 = value.ToSymbol("GREEN")
var const2 value.Value // constant: Color::GREEN, loc: <main>:4:6
var cc_main_3 = &vm.CallCache{}
var sym4 = value.ToSymbol("BLUE")
var const3 value.Value // constant: Color::BLUE, loc: <main>:5:6
var sym13 = value.ToSymbol("main")
var sym15 = value.ToSymbol("baz")

var const0 *value.Class // Color
var sym0 = value.ToSymbol("Color")

var sym5 = value.ToSymbol("Color.:#init")
var sym6 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt8, l1 value.UInt8, l2 value.UInt8) (result value.Value, err value.Value) { // method: Color.:#init, loc: <main>:10:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 2, (l0).ToValue())
	value.SetInstanceVariable(self, 1, (l1).ToValue())
	value.SetInstanceVariable(self, 0, (l2).ToValue())
	return self, value.Undefined

}

var sym10 = value.ToSymbol("blue")
var sym11 = value.ToSymbol("green")
var sym12 = value.ToSymbol("red")
var sym7 = value.ToSymbol("Std::Kernel::baz")
var sym8 = value.ToSymbol("inspect")
var fn_method2 vm.NativeFunction // Std::Value.:inspect
var sym9 = value.ToSymbol("println@1")
var fn_method3 vm.NativeFunction // Std::Kernel::println@1
func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Std::Kernel::baz, loc: <main>:13:5
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var t3 value.Value
	_ = t3

	callFrame = thread.AddNativeCallFrame(sym7, sym6, 13)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = l0
	callFrame.SetNativeLineNumber(14)
	t1, err = fn_method2(thread, t2) // receiver: Color, name: inspect
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = self
	t2[1] = (value.String("col: ") + (t1).AsString()).ToValue()
	t3, err = fn_method3(thread, t2) // receiver: Std::Kernel, name: println@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	return t3, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var namespace value.Value
	_ = namespace
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(255)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_1, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym2, t1)
	const1 = t1

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(255)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_2, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym3, t1)
	const2 = t1

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(255)).ToValue()
	callFrame.SetNativeLineNumber(5)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_3, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym4, t1)
	const3 = t1

	methodDefinitions()
	fn_method2 = vm.MethodToFunc((value.ValueClass).LookupMethod(sym8))
	fn_method3 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym9))

	callFrame = thread.AddNativeCallFrame(sym13, sym6, 1)
	defer thread.PopNativeCallFrame()
	_, err = fn_method1(thread, (value.KernelModule).ToValue(), const1) // receiver: Std::Kernel, name: baz
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}

func initGlobalEnv() {
	var parentNamespace value.Value
	_ = parentNamespace
	var namespace value.Value
	_ = namespace
	var class *value.Class
	_ = class
	var superclass *value.Class
	_ = superclass
	var mixin *value.Mixin
	_ = mixin

	parentNamespace = (value.RootModule).ToValue()
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("blue"): 0, symbol.ToSymbol("green"): 1, symbol.ToSymbol("red"): 2}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Color
	vm.Def(class, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt8(), (args[2]).AsUInt8(), (args[3]).AsUInt8())
		return result, err
	}, vm.DefWithParameters(3))
	vm.DefineGetter(class, sym10, 0)
	vm.DefineGetter(class, sym11, 1)
	vm.DefineGetter(class, sym12, 2)
	class = (value.KernelModule).SingletonClass() // Std::Kernel
	vm.Def(class, "baz", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},
		"infer namespace in local declaration": {
			input: `
				class Color
					const RED: Color = new(255u8, 0u8, 0u8)
					const GREEN: Color = new(0u8, 255u8, 0u8)
					const BLUE: Color = new(0u8, 0u8, 255u8)

					getter red: UInt8
					getter green: UInt8
					getter blue: UInt8
					init(@red, @green, @blue); end
				end

				var baz: Color = .::RED
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym1 = value.ToSymbol("#init")
var cc_main_1 = &vm.CallCache{}
var sym2 = value.ToSymbol("RED")
var const1 value.Value // constant: Color::RED, loc: <main>:3:6
var cc_main_2 = &vm.CallCache{}
var sym3 = value.ToSymbol("GREEN")
var const2 value.Value // constant: Color::GREEN, loc: <main>:4:6
var cc_main_3 = &vm.CallCache{}
var sym4 = value.ToSymbol("BLUE")
var const3 value.Value // constant: Color::BLUE, loc: <main>:5:6
var sym10 = value.ToSymbol("main")

var const0 *value.Class // Color
var sym0 = value.ToSymbol("Color")

var sym5 = value.ToSymbol("Color.:#init")
var sym6 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt8, l1 value.UInt8, l2 value.UInt8) (result value.Value, err value.Value) { // method: Color.:#init, loc: <main>:10:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 2, (l0).ToValue())
	value.SetInstanceVariable(self, 1, (l1).ToValue())
	value.SetInstanceVariable(self, 0, (l2).ToValue())
	return self, value.Undefined

}

var sym7 = value.ToSymbol("blue")
var sym8 = value.ToSymbol("green")
var sym9 = value.ToSymbol("red")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var namespace value.Value
	_ = namespace
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var baz: Color
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(255)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_1, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym2, t1)
	const1 = t1

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(255)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_2, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym3, t1)
	const2 = t1

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(255)).ToValue()
	callFrame.SetNativeLineNumber(5)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_3, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym4, t1)
	const3 = t1

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym10, sym6, 1)
	defer thread.PopNativeCallFrame()
	l0 = const1
}

func initGlobalEnv() {
	var parentNamespace value.Value
	_ = parentNamespace
	var namespace value.Value
	_ = namespace
	var class *value.Class
	_ = class
	var superclass *value.Class
	_ = superclass
	var mixin *value.Mixin
	_ = mixin

	parentNamespace = (value.RootModule).ToValue()
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("blue"): 0, symbol.ToSymbol("green"): 1, symbol.ToSymbol("red"): 2}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Color
	vm.Def(class, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt8(), (args[2]).AsUInt8(), (args[3]).AsUInt8())
		return result, err
	}, vm.DefWithParameters(3))
	vm.DefineGetter(class, sym7, 0)
	vm.DefineGetter(class, sym8, 1)
	vm.DefineGetter(class, sym9, 2)
}
`,
		},
		"infer namespace in binary operator": {
			input: `
				class Color
					const RED: Color = new(255u8, 0u8, 0u8)
					const GREEN: Color = new(0u8, 255u8, 0u8)
					const BLUE: Color = new(0u8, 0u8, 255u8)

					getter red: UInt8
					getter green: UInt8
					getter blue: UInt8
					init(@red, @green, @blue); end
				end

				baz := Color::RED
				baz == .::BLUE
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym1 = value.ToSymbol("#init")
var cc_main_1 = &vm.CallCache{}
var sym2 = value.ToSymbol("RED")
var const1 value.Value // constant: Color::RED, loc: <main>:3:6
var cc_main_2 = &vm.CallCache{}
var sym3 = value.ToSymbol("GREEN")
var const2 value.Value // constant: Color::GREEN, loc: <main>:4:6
var cc_main_3 = &vm.CallCache{}
var sym4 = value.ToSymbol("BLUE")
var const3 value.Value // constant: Color::BLUE, loc: <main>:5:6
var sym10 = value.ToSymbol("main")
var sym12 = value.ToSymbol("==")
var fn_method1 vm.NativeFunction // Std::Value.:==

var const0 *value.Class // Color
var sym0 = value.ToSymbol("Color")

var sym5 = value.ToSymbol("Color.:#init")
var sym6 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt8, l1 value.UInt8, l2 value.UInt8) (result value.Value, err value.Value) { // method: Color.:#init, loc: <main>:10:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 2, (l0).ToValue())
	value.SetInstanceVariable(self, 1, (l1).ToValue())
	value.SetInstanceVariable(self, 0, (l2).ToValue())
	return self, value.Undefined

}

var sym7 = value.ToSymbol("blue")
var sym8 = value.ToSymbol("green")
var sym9 = value.ToSymbol("red")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var namespace value.Value
	_ = namespace
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var baz: Color
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(255)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_1, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym2, t1)
	const1 = t1

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(255)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_2, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym3, t1)
	const2 = t1

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(255)).ToValue()
	callFrame.SetNativeLineNumber(5)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_3, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym4, t1)
	const3 = t1

	methodDefinitions()
	fn_method1 = vm.MethodToFunc((value.ValueClass).LookupMethod(sym12))

	callFrame = thread.AddNativeCallFrame(sym10, sym6, 1)
	defer thread.PopNativeCallFrame()
	l0 = const1
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = const3
	callFrame.SetNativeLineNumber(14)
	_, err = fn_method1(thread, t2) // receiver: Color, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}

func initGlobalEnv() {
	var parentNamespace value.Value
	_ = parentNamespace
	var namespace value.Value
	_ = namespace
	var class *value.Class
	_ = class
	var superclass *value.Class
	_ = superclass
	var mixin *value.Mixin
	_ = mixin

	parentNamespace = (value.RootModule).ToValue()
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("blue"): 0, symbol.ToSymbol("green"): 1, symbol.ToSymbol("red"): 2}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Color
	vm.Def(class, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt8(), (args[2]).AsUInt8(), (args[3]).AsUInt8())
		return result, err
	}, vm.DefWithParameters(3))
	vm.DefineGetter(class, sym7, 0)
	vm.DefineGetter(class, sym8, 1)
	vm.DefineGetter(class, sym9, 2)
}
`,
		},
		"infer namespace in switch pattern": {
			input: `
				class Color
					const RED: Color = new(255u8, 0u8, 0u8)
					const GREEN: Color = new(0u8, 255u8, 0u8)
					const BLUE: Color = new(0u8, 0u8, 255u8)

					getter red: UInt8
					getter green: UInt8
					getter blue: UInt8
					init(@red, @green, @blue); end
				end

				baz := Color::RED
				switch baz
				case .::GREEN
					println("green")
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym1 = value.ToSymbol("#init")
var cc_main_1 = &vm.CallCache{}
var sym2 = value.ToSymbol("RED")
var const1 value.Value // constant: Color::RED, loc: <main>:3:6
var cc_main_2 = &vm.CallCache{}
var sym3 = value.ToSymbol("GREEN")
var const2 value.Value // constant: Color::GREEN, loc: <main>:4:6
var cc_main_3 = &vm.CallCache{}
var sym4 = value.ToSymbol("BLUE")
var const3 value.Value // constant: Color::BLUE, loc: <main>:5:6
var sym10 = value.ToSymbol("main")
var sym12 = value.ToSymbol("==")
var fn_method1 vm.NativeFunction // Std::Value.:==
var sym13 = value.ToSymbol("println@1")
var fn_method2 vm.NativeFunction // Std::Kernel::println@1

var const0 *value.Class // Color
var sym0 = value.ToSymbol("Color")

var sym5 = value.ToSymbol("Color.:#init")
var sym6 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt8, l1 value.UInt8, l2 value.UInt8) (result value.Value, err value.Value) { // method: Color.:#init, loc: <main>:10:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 2, (l0).ToValue())
	value.SetInstanceVariable(self, 1, (l1).ToValue())
	value.SetInstanceVariable(self, 0, (l2).ToValue())
	return self, value.Undefined

}

var sym7 = value.ToSymbol("blue")
var sym8 = value.ToSymbol("green")
var sym9 = value.ToSymbol("red")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var namespace value.Value
	_ = namespace
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var baz: Color
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(255)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_1, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym2, t1)
	const1 = t1

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(255)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_2, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym3, t1)
	const2 = t1

	namespace = (const0).ToValue()
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(255)).ToValue()
	callFrame.SetNativeLineNumber(5)
	t1, err = thread.CallMethodByNameWithCache(sym1, &cc_main_3, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	value.AddConstant(namespace, sym4, t1)
	const3 = t1

	methodDefinitions()
	fn_method1 = vm.MethodToFunc((value.ValueClass).LookupMethod(sym12))
	fn_method2 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym13))

	callFrame = thread.AddNativeCallFrame(sym10, sym6, 1)
	defer thread.PopNativeCallFrame()
	l0 = const1
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = l0
	t2[1] = const2
	callFrame.SetNativeLineNumber(15)
	t1, err = fn_method1(thread, t2) // receiver: Color, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t1) {
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = (value.String("green")).ToValue()
		callFrame.SetNativeLineNumber(16)
		_, err = fn_method2(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		goto lbl1
	}
lbl1:
}

func initGlobalEnv() {
	var parentNamespace value.Value
	_ = parentNamespace
	var namespace value.Value
	_ = namespace
	var class *value.Class
	_ = class
	var superclass *value.Class
	_ = superclass
	var mixin *value.Mixin
	_ = mixin

	parentNamespace = (value.RootModule).ToValue()
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("blue"): 0, symbol.ToSymbol("green"): 1, symbol.ToSymbol("red"): 2}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Color
	vm.Def(class, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], (args[1]).AsUInt8(), (args[2]).AsUInt8(), (args[3]).AsUInt8())
		return result, err
	}, vm.DefWithParameters(3))
	vm.DefineGetter(class, sym7, 0)
	vm.DefineGetter(class, sym8, 1)
	vm.DefineGetter(class, sym9, 2)
}
`,
		},
		"builtin absolute path ::Std": {
			input: "a := ::Std",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.StdModule).ToValue()
}
`,
		},
		"custom absolute path": {
			input: `
				module Foo; end
				a := ::Foo
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Foo
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = (const0).ToValue()
}

func initGlobalEnv() {
	var parentNamespace value.Value
	_ = parentNamespace
	var namespace value.Value
	_ = namespace
	var class *value.Class
	_ = class
	var superclass *value.Class
	_ = superclass
	var mixin *value.Mixin
	_ = mixin

	parentNamespace = (value.RootModule).ToValue()
	const0 = value.NewModule()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}
`,
		},
		"builtin absolute nested path ::Std::Float::INF": {
			input: "a := ::Std::Float::INF",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.FloatInf()
}
`,
		},
		"custom absolute nested path ::Std::Float::INF": {
			input: `
				module Foo
					module Bar
						const BAZ = 5
					end
				end
				a := ::Foo::Bar::BAZ
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym2 = value.ToSymbol("BAZ")
var const2 value.SmallInt // constant: Foo::Bar::BAZ, loc: <main>:4:7
var sym3 = value.ToSymbol("main")
var sym4 = value.ToSymbol("<main>")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")
var const1 *value.Module // Foo::Bar
var sym1 = value.ToSymbol("Bar")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var namespace value.Value
	_ = namespace
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	namespace = (const1).ToValue()
	value.AddConstant(namespace, sym2, (value.SmallInt(5)).ToValue())
	const2 = value.SmallInt(5)
	callFrame = thread.AddNativeCallFrame(sym3, sym4, 1)
	defer thread.PopNativeCallFrame()
	l0 = (const2).ToValue()
}

func initGlobalEnv() {
	var parentNamespace value.Value
	_ = parentNamespace
	var namespace value.Value
	_ = namespace
	var class *value.Class
	_ = class
	var superclass *value.Class
	_ = superclass
	var mixin *value.Mixin
	_ = mixin

	parentNamespace = (value.RootModule).ToValue()
	const0 = value.NewModule()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	parentNamespace = (const0).ToValue()
	const1 = value.NewModule()
	namespace = value.Ref(const1)
	value.AddConstant(parentNamespace, sym1, namespace)

}
`,
		},
		"relative path from using": {
			input: `
				using Std::Float::INF as I
				a := I
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: &Std::Float
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.FloatInf()).ToValue()
}
`,
		},
		"relative path": {
			input: `
				module Foo
					const BAR = 3
					module Baz
						println BAR
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym2 = value.ToSymbol("BAR")
var const2 value.SmallInt // constant: Foo::BAR, loc: <main>:3:6
var sym3 = value.ToSymbol("main")
var sym4 = value.ToSymbol("<main>")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")
var const1 *value.Module // Foo::Baz
var sym1 = value.ToSymbol("Baz")

var sym5 = value.ToSymbol("<module: Foo>")

var sym6 = value.ToSymbol("<module: Foo::Baz>")
var sym7 = value.ToSymbol("println@1")
var fn_method0 vm.NativeFunction // Std::Kernel::println@1

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var namespace value.Value
	_ = namespace
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	namespace = (const0).ToValue()
	value.AddConstant(namespace, sym2, (value.SmallInt(3)).ToValue())
	const2 = value.SmallInt(3)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym7))

	callFrame = thread.AddNativeCallFrame(sym3, sym4, 1)
	defer thread.PopNativeCallFrame()
	fn_ns_expr0(thread)
}

func initGlobalEnv() {
	var parentNamespace value.Value
	_ = parentNamespace
	var namespace value.Value
	_ = namespace
	var class *value.Class
	_ = class
	var superclass *value.Class
	_ = superclass
	var mixin *value.Mixin
	_ = mixin

	parentNamespace = (value.RootModule).ToValue()
	const0 = value.NewModule()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	parentNamespace = (const0).ToValue()
	const1 = value.NewModule()
	namespace = value.Ref(const1)
	value.AddConstant(parentNamespace, sym1, namespace)

}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame

	self = value.Ref((const0).SingletonClass())
	callFrame = thread.AddNativeCallFrame(sym5, sym4, 2)
	defer thread.PopNativeCallFrame()
	fn_ns_expr1(thread)
}

func fn_ns_expr1(thread *vm.Thread) { // namespace: Foo::Baz, loc: <main>:4:6
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err

	self = value.Ref((const1).SingletonClass())
	callFrame = thread.AddNativeCallFrame(sym6, sym4, 4)
	defer thread.PopNativeCallFrame()
	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = (value.KernelModule).ToValue()
	t1[1] = (const2).ToValue()
	callFrame.SetNativeLineNumber(5)
	_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
}
`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoDefConstant(t *testing.T) {
	tests := goTestTable{
		"relative path Foo": {
			input: "const Foo = 3",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("Foo")
var const0 value.SmallInt // constant: Foo, loc: <main>:1:1
var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread

	defer func() {
		switch r := recover().(type) {
		case value.Value:
			thread.Exit(r)
		case nil:
		default:
			panic(r)
		}
	}()

	var namespace value.Value
	_ = namespace
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	namespace = (value.RootModule).ToValue()
	value.AddConstant(namespace, sym0, (value.SmallInt(3)).ToValue())
	const0 = value.SmallInt(3)
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
	defer thread.PopNativeCallFrame()
}
`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}
