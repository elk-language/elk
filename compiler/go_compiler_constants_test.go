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
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym9 = value.ToSymbol("main")   // main
var sym11 = value.ToSymbol("baz")   // baz
var sym12 = value.ToSymbol("#init") // #init
var cc_main_1 = &vm.CallCache{}
var sym13 = value.ToSymbol("RED") // RED
var const1 value.Value            // const Color::RED: Color, loc: <main>:3:6
var cc_main_2 = &vm.CallCache{}
var sym14 = value.ToSymbol("GREEN") // GREEN
var const2 value.Value              // const Color::GREEN: Color, loc: <main>:4:6
var cc_main_3 = &vm.CallCache{}
var sym15 = value.ToSymbol("BLUE") // BLUE
var const3 value.Value             // const Color::BLUE: Color, loc: <main>:5:6

var const0 *value.Class            // Color
var sym0 = value.ToSymbol("Color") // Color

var sym1 = value.ToSymbol("Color.:#init") // Color.:#init
var sym2 = value.ToSymbol("<main>")       // <main>
func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt8, l1 value.UInt8, l2 value.UInt8) (result value.Value, err value.Value) { // method: Color.:#init, loc: <main>:10:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 2, (l0).ToValue())
	value.SetInstanceVariable(self, 1, (l1).ToValue())
	value.SetInstanceVariable(self, 0, (l2).ToValue())
	return self, value.Undefined

}

var sym6 = value.ToSymbol("blue")             // blue
var sym7 = value.ToSymbol("green")            // green
var sym8 = value.ToSymbol("red")              // red
var sym3 = value.ToSymbol("Std::Kernel::baz") // Std::Kernel::baz
var sym4 = value.ToSymbol("inspect")          // inspect
var fn_method2 vm.NativeFunction              // Std::Value.:inspect
var sym5 = value.ToSymbol("println@1")        // println@1
var fn_method3 vm.NativeFunction              // Std::Kernel::println@1
func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Std::Kernel::baz, loc: <main>:13:5
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var t3 value.Value
	_ = t3

	callFrame = thread.AddNativeCallFrame(sym2, sym3, 13)
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var err value.Value
	_ = err
	var namespace value.Value
	_ = namespace
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	fn_method2 = vm.MethodToFunc((value.ValueClass).LookupMethod(sym4))
	fn_method3 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym5))

	namespace = (const0).ToValue()
	self = namespace
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(255)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = thread.CallMethodByNameWithCache(sym12, &cc_main_1, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const1 = t1
	value.AddConstant(namespace, sym13, const1)

	namespace = (const0).ToValue()
	self = namespace
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(255)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = thread.CallMethodByNameWithCache(sym12, &cc_main_2, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const2 = t1
	value.AddConstant(namespace, sym14, const2)

	namespace = (const0).ToValue()
	self = namespace
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(255)).ToValue()
	callFrame.SetNativeLineNumber(5)
	t1, err = thread.CallMethodByNameWithCache(sym12, &cc_main_3, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const3 = t1
	value.AddConstant(namespace, sym15, const3)

	callFrame = thread.AddNativeCallFrame(sym2, sym9, 1)
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
	vm.DefineGetter(class, sym6, 0)
	vm.DefineGetter(class, sym7, 1)
	vm.DefineGetter(class, sym8, 2)
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
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym6 = value.ToSymbol("main")  // main
var sym8 = value.ToSymbol("#init") // #init
var cc_main_1 = &vm.CallCache{}
var sym9 = value.ToSymbol("RED") // RED
var const1 value.Value           // const Color::RED: Color, loc: <main>:3:6
var cc_main_2 = &vm.CallCache{}
var sym10 = value.ToSymbol("GREEN") // GREEN
var const2 value.Value              // const Color::GREEN: Color, loc: <main>:4:6
var cc_main_3 = &vm.CallCache{}
var sym11 = value.ToSymbol("BLUE") // BLUE
var const3 value.Value             // const Color::BLUE: Color, loc: <main>:5:6

var const0 *value.Class            // Color
var sym0 = value.ToSymbol("Color") // Color

var sym1 = value.ToSymbol("Color.:#init") // Color.:#init
var sym2 = value.ToSymbol("<main>")       // <main>
func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt8, l1 value.UInt8, l2 value.UInt8) (result value.Value, err value.Value) { // method: Color.:#init, loc: <main>:10:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 2, (l0).ToValue())
	value.SetInstanceVariable(self, 1, (l1).ToValue())
	value.SetInstanceVariable(self, 0, (l2).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("blue")  // blue
var sym4 = value.ToSymbol("green") // green
var sym5 = value.ToSymbol("red")   // red

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
	var l0 value.Value // var baz: Color
	_ = l0
	var namespace value.Value
	_ = namespace
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()

	namespace = (const0).ToValue()
	self = namespace
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(255)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = thread.CallMethodByNameWithCache(sym8, &cc_main_1, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const1 = t1
	value.AddConstant(namespace, sym9, const1)

	namespace = (const0).ToValue()
	self = namespace
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(255)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = thread.CallMethodByNameWithCache(sym8, &cc_main_2, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const2 = t1
	value.AddConstant(namespace, sym10, const2)

	namespace = (const0).ToValue()
	self = namespace
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(255)).ToValue()
	callFrame.SetNativeLineNumber(5)
	t1, err = thread.CallMethodByNameWithCache(sym8, &cc_main_3, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const3 = t1
	value.AddConstant(namespace, sym11, const3)

	callFrame = thread.AddNativeCallFrame(sym2, sym6, 1)
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
	vm.DefineGetter(class, sym3, 0)
	vm.DefineGetter(class, sym4, 1)
	vm.DefineGetter(class, sym5, 2)
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
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym6 = value.ToSymbol("main")  // main
var sym8 = value.ToSymbol("==")    // ==
var fn_method1 vm.NativeFunction   // Std::Value.:==
var sym9 = value.ToSymbol("#init") // #init
var cc_main_1 = &vm.CallCache{}
var sym10 = value.ToSymbol("RED") // RED
var const1 value.Value            // const Color::RED: Color, loc: <main>:3:6
var cc_main_2 = &vm.CallCache{}
var sym11 = value.ToSymbol("GREEN") // GREEN
var const2 value.Value              // const Color::GREEN: Color, loc: <main>:4:6
var cc_main_3 = &vm.CallCache{}
var sym12 = value.ToSymbol("BLUE") // BLUE
var const3 value.Value             // const Color::BLUE: Color, loc: <main>:5:6

var const0 *value.Class            // Color
var sym0 = value.ToSymbol("Color") // Color

var sym1 = value.ToSymbol("Color.:#init") // Color.:#init
var sym2 = value.ToSymbol("<main>")       // <main>
func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt8, l1 value.UInt8, l2 value.UInt8) (result value.Value, err value.Value) { // method: Color.:#init, loc: <main>:10:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 2, (l0).ToValue())
	value.SetInstanceVariable(self, 1, (l1).ToValue())
	value.SetInstanceVariable(self, 0, (l2).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("blue")  // blue
var sym4 = value.ToSymbol("green") // green
var sym5 = value.ToSymbol("red")   // red

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
	var l0 value.Value // var baz: Color
	_ = l0
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var namespace value.Value
	_ = namespace
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	fn_method1 = vm.MethodToFunc((value.ValueClass).LookupMethod(sym8))

	namespace = (const0).ToValue()
	self = namespace
	t1 = value.ResizeNativeArgs(t1, 5)
	t1[0] = const0.CreateInstance()
	t1[1] = (value.UInt8(255)).ToValue()
	t1[2] = (value.UInt8(0)).ToValue()
	t1[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t2, err = thread.CallMethodByNameWithCache(sym9, &cc_main_1, t1...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const1 = t2
	value.AddConstant(namespace, sym10, const1)

	namespace = (const0).ToValue()
	self = namespace
	t1 = value.ResizeNativeArgs(t1, 5)
	t1[0] = const0.CreateInstance()
	t1[1] = (value.UInt8(0)).ToValue()
	t1[2] = (value.UInt8(255)).ToValue()
	t1[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t2, err = thread.CallMethodByNameWithCache(sym9, &cc_main_2, t1...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const2 = t2
	value.AddConstant(namespace, sym11, const2)

	namespace = (const0).ToValue()
	self = namespace
	t1 = value.ResizeNativeArgs(t1, 5)
	t1[0] = const0.CreateInstance()
	t1[1] = (value.UInt8(0)).ToValue()
	t1[2] = (value.UInt8(0)).ToValue()
	t1[3] = (value.UInt8(255)).ToValue()
	callFrame.SetNativeLineNumber(5)
	t2, err = thread.CallMethodByNameWithCache(sym9, &cc_main_3, t1...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const3 = t2
	value.AddConstant(namespace, sym12, const3)

	callFrame = thread.AddNativeCallFrame(sym2, sym6, 1)
	defer thread.PopNativeCallFrame()
	l0 = const1
	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = l0
	t1[1] = const3
	callFrame.SetNativeLineNumber(14)
	_, err = fn_method1(thread, t1) // receiver: Color, name: ==
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
	vm.DefineGetter(class, sym3, 0)
	vm.DefineGetter(class, sym4, 1)
	vm.DefineGetter(class, sym5, 2)
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
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym6 = value.ToSymbol("main")      // main
var sym8 = value.ToSymbol("==")        // ==
var fn_method1 vm.NativeFunction       // Std::Value.:==
var sym9 = value.ToSymbol("println@1") // println@1
var fn_method2 vm.NativeFunction       // Std::Kernel::println@1
var sym10 = value.ToSymbol("#init")    // #init
var cc_main_1 = &vm.CallCache{}
var sym11 = value.ToSymbol("RED") // RED
var const1 value.Value            // const Color::RED: Color, loc: <main>:3:6
var cc_main_2 = &vm.CallCache{}
var sym12 = value.ToSymbol("GREEN") // GREEN
var const2 value.Value              // const Color::GREEN: Color, loc: <main>:4:6
var cc_main_3 = &vm.CallCache{}
var sym13 = value.ToSymbol("BLUE") // BLUE
var const3 value.Value             // const Color::BLUE: Color, loc: <main>:5:6

var const0 *value.Class            // Color
var sym0 = value.ToSymbol("Color") // Color

var sym1 = value.ToSymbol("Color.:#init") // Color.:#init
var sym2 = value.ToSymbol("<main>")       // <main>
func fn_method0(thread *vm.Thread, self value.Value, l0 value.UInt8, l1 value.UInt8, l2 value.UInt8) (result value.Value, err value.Value) { // method: Color.:#init, loc: <main>:10:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 2, (l0).ToValue())
	value.SetInstanceVariable(self, 1, (l1).ToValue())
	value.SetInstanceVariable(self, 0, (l2).ToValue())
	return self, value.Undefined

}

var sym3 = value.ToSymbol("blue")  // blue
var sym4 = value.ToSymbol("green") // green
var sym5 = value.ToSymbol("red")   // red

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
	var l0 value.Value // var baz: Color
	_ = l0
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var namespace value.Value
	_ = namespace
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	fn_method1 = vm.MethodToFunc((value.ValueClass).LookupMethod(sym8))
	fn_method2 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym9))

	namespace = (const0).ToValue()
	self = namespace
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(255)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = thread.CallMethodByNameWithCache(sym10, &cc_main_1, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const1 = t1
	value.AddConstant(namespace, sym11, const1)

	namespace = (const0).ToValue()
	self = namespace
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(255)).ToValue()
	t2[3] = (value.UInt8(0)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = thread.CallMethodByNameWithCache(sym10, &cc_main_2, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const2 = t1
	value.AddConstant(namespace, sym12, const2)

	namespace = (const0).ToValue()
	self = namespace
	t2 = value.ResizeNativeArgs(t2, 5)
	t2[0] = const0.CreateInstance()
	t2[1] = (value.UInt8(0)).ToValue()
	t2[2] = (value.UInt8(0)).ToValue()
	t2[3] = (value.UInt8(255)).ToValue()
	callFrame.SetNativeLineNumber(5)
	t1, err = thread.CallMethodByNameWithCache(sym10, &cc_main_3, t2...) // receiver: %self, name: #init
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const3 = t1
	value.AddConstant(namespace, sym13, const3)

	callFrame = thread.AddNativeCallFrame(sym2, sym6, 1)
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
	vm.DefineGetter(class, sym3, 0)
	vm.DefineGetter(class, sym4, 1)
	vm.DefineGetter(class, sym5, 2)
}
`,
		},
		"builtin absolute path ::Std": {
			input: "a := ::Std",
			want: `package main

import (
	"github.com/elk-language/elk"
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")   // main
var sym1 = value.ToSymbol("<main>") // <main>

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
	callFrame = thread.AddNativeCallFrame(sym1, sym0, 1)
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
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym1 = value.ToSymbol("main")   // main
var sym2 = value.ToSymbol("<main>") // <main>

var const0 *value.Module         // Foo
var sym0 = value.ToSymbol("Foo") // Foo

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
	callFrame = thread.AddNativeCallFrame(sym2, sym1, 1)
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
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")   // main
var sym1 = value.ToSymbol("<main>") // <main>

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
	callFrame = thread.AddNativeCallFrame(sym1, sym0, 1)
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
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym2 = value.ToSymbol("main")   // main
var sym3 = value.ToSymbol("<main>") // <main>
var sym6 = value.ToSymbol("BAZ")    // BAZ
var const2 value.SmallInt           // const Foo::Bar::BAZ: 5, loc: <main>:4:7

var const0 *value.Module         // Foo
var sym0 = value.ToSymbol("Foo") // Foo
var const1 *value.Module         // Foo::Bar
var sym1 = value.ToSymbol("Bar") // Bar

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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var namespace value.Value
	_ = namespace
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	namespace = (const1).ToValue()
	self = namespace
	const2 = value.SmallInt(5)
	value.AddConstant(namespace, sym6, (const2).ToValue())

	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
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
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")   // main
var sym1 = value.ToSymbol("<main>") // <main>

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
	callFrame = thread.AddNativeCallFrame(sym1, sym0, 1)
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
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym2 = value.ToSymbol("main")   // main
var sym3 = value.ToSymbol("<main>") // <main>
var sym7 = value.ToSymbol("BAR")    // BAR
var const2 value.SmallInt           // const Foo::BAR: 3, loc: <main>:3:6

var const0 *value.Module         // Foo
var sym0 = value.ToSymbol("Foo") // Foo
var const1 *value.Module         // Foo::Baz
var sym1 = value.ToSymbol("Baz") // Baz

var sym4 = value.ToSymbol("<module: Foo>") // <module: Foo>

var sym5 = value.ToSymbol("<module: Foo::Baz>") // <module: Foo::Baz>
var sym6 = value.ToSymbol("println@1")          // println@1
var fn_method0 vm.NativeFunction                // Std::Kernel::println@1

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
	var namespace value.Value
	_ = namespace
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym6))

	namespace = (const0).ToValue()
	self = namespace
	const2 = value.SmallInt(3)
	value.AddConstant(namespace, sym7, (const2).ToValue())

	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
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
	callFrame = thread.AddNativeCallFrame(sym3, sym4, 2)
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
	callFrame = thread.AddNativeCallFrame(sym3, sym5, 4)
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
		"define constant that call a singleton method": {
			input: `
				module Foo
					def lol: Int then 3

					const LOL: Int = lol()
				end
`,
			want: `package main

import (
	"github.com/elk-language/elk"
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym3 = value.ToSymbol("main") // main
var sym5 = value.ToSymbol("lol")  // lol
var cc_main_1 = &vm.CallCache{}
var sym6 = value.ToSymbol("LOL") // LOL
var const1 value.Value           // const Foo::LOL: Std::Int, loc: <main>:5:6

var const0 *value.Module         // Foo
var sym0 = value.ToSymbol("Foo") // Foo

var sym1 = value.ToSymbol("Foo::lol") // Foo::lol
var sym2 = value.ToSymbol("<main>")   // <main>
func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo::lol, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return (value.SmallInt(3)).ToValue(), value.Undefined

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

	var callFrame *vm.CallFrame
	_ = callFrame
	var namespace value.Value
	_ = namespace
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()

	namespace = (const0).ToValue()
	self = namespace
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = self
	callFrame.SetNativeLineNumber(5)
	t1, err = thread.CallMethodByNameWithCache(sym5, &cc_main_1, t2...) // receiver: self, name: lol
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	const1 = t1
	value.AddConstant(namespace, sym6, const1)

	callFrame = thread.AddNativeCallFrame(sym2, sym3, 1)
	defer thread.PopNativeCallFrame()
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

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass() // Foo
	vm.Def(class, "lol", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},
		"relative path Foo": {
			input: "const Foo = 3",
			want: `package main

import (
	"github.com/elk-language/elk"
	_ "github.com/elk-language/elk/ini"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")   // main
var sym1 = value.ToSymbol("<main>") // <main>
var sym2 = value.ToSymbol("Foo")    // Foo
var const0 value.SmallInt           // const Foo: 3, loc: <main>:1:1

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
	var namespace value.Value
	_ = namespace
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	namespace = (value.RootModule).ToValue()
	self = namespace
	const0 = value.SmallInt(3)
	value.AddConstant(namespace, sym2, (const0).ToValue())

	callFrame = thread.AddNativeCallFrame(sym1, sym0, 1)
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
