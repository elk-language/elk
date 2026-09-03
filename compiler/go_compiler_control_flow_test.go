package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

func TestGoGoExpression(t *testing.T) {
	tests := goTestTable{
		"with a single expression": {
			input: "go println('foo')",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("<closure>")
var sym3 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 *vm.NativeClosure
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t1 = vm.NewNativeClosure(
		func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) { // name: fn_cl0, sig: Std::Thread, loc: <main>:1:1
			var callFrame *vm.CallFrame
			_ = callFrame
			var t1 value.Value
			_ = t1
			var t2 []value.Value
			_ = t2
			var err value.Value
			_ = err

			callFrame = thread.AddNativeCallFrame(sym2, sym1, 1)
			defer thread.PopNativeCallFrame()
			t2 = value.ResizeNativeArgs(t2, 3)
			t2[0] = (value.KernelModule).ToValue()
			t2[1] = (value.String("foo")).ToValue()
			t1, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
			if err.IsNotUndefined() {
				thread.CaptureStackTrace()
				return value.Undefined, err
			}
			return t1, value.Undefined
		},
		0,
		position.NewLocation("<main>", position.NewSpan(position.New(0, 1, 1), position.New(0, 1, 1))),
	)
	thread.GoNative(t1)
}
`,
		},
		"with captured thread": {
			input: "a := go println('foo')",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("<closure>")
var sym3 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Thread
	_ = l0
	var t1 *vm.NativeClosure
	_ = t1
	var t2 *vm.Thread
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t1 = vm.NewNativeClosure(
		func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) { // name: fn_cl0, sig: Std::Thread, loc: <main>:1:6
			var callFrame *vm.CallFrame
			_ = callFrame
			var t1 value.Value
			_ = t1
			var t2 []value.Value
			_ = t2
			var err value.Value
			_ = err

			callFrame = thread.AddNativeCallFrame(sym2, sym1, 1)
			defer thread.PopNativeCallFrame()
			t2 = value.ResizeNativeArgs(t2, 3)
			t2[0] = (value.KernelModule).ToValue()
			t2[1] = (value.String("foo")).ToValue()
			t1, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
			if err.IsNotUndefined() {
				thread.CaptureStackTrace()
				return value.Undefined, err
			}
			return t1, value.Undefined
		},
		0,
		position.NewLocation("<main>", position.NewSpan(position.New(5, 1, 6), position.New(5, 1, 6))),
	)
	t2 = thread.GoNative(t1)
	l0 = (t2).ToValue()
}
`,
		},
		"with outer variables": {
			input: `
				a := 5
				go
					println("foo")
					println(a)
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("<closure>")
var sym3 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var t1 *vm.NativeClosure
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	{
		l0 := l0 // close: var a: Std::Int
		t1 = vm.NewNativeClosure(
			func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) { // name: fn_cl0, sig: Std::Thread, loc: <main>:3:5
				var callFrame *vm.CallFrame
				_ = callFrame
				var t1 []value.Value
				_ = t1
				var err value.Value
				_ = err
				var t2 value.Value
				_ = t2

				callFrame = thread.AddNativeCallFrame(sym2, sym1, 3)
				defer thread.PopNativeCallFrame()
				t1 = value.ResizeNativeArgs(t1, 3)
				t1[0] = (value.KernelModule).ToValue()
				t1[1] = (value.String("foo")).ToValue()
				callFrame.SetNativeLineNumber(4)
				_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
				if err.IsNotUndefined() {
					thread.CaptureStackTrace()
					return value.Undefined, err
				}
				t1 = value.ResizeNativeArgs(t1, 3)
				t1[0] = (value.KernelModule).ToValue()
				t1[1] = l0
				callFrame.SetNativeLineNumber(5)
				t2, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
				if err.IsNotUndefined() {
					thread.CaptureStackTrace()
					return value.Undefined, err
				}
				return t2, value.Undefined
			},
			0,
			position.NewLocation("<main>", position.NewSpan(position.New(16, 3, 5), position.New(16, 3, 5))),
		)
	}
	thread.GoNative(t1)
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

func TestGoForInExpression(t *testing.T) {
	tests := goTestTable{
		"int literal": {
			input: `
				for i in 20
					println(i)
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var i: Std::Int
	_ = l0
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(0)
	for {
		if !(value.Bool((l0).LessThanSmallInt(value.SmallInt(20)))) {
			break
		}
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = (l0).ToValue()
		callFrame.SetNativeLineNumber(3)
		t1, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		l0++
	}
}
`,
		},
		"int value": {
			input: `
				j := 20
				for i in j
					println(i)
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(20)).ToValue()
	l1 = value.SmallInt(0).ToValue()
	for {
		if !(value.Bool(value.LessThanInts(l1, l0))) {
			break
		}
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = l1
		callFrame.SetNativeLineNumber(4)
		t1, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		l1 = value.IncrementInt(l1)
	}
}
`,
		},
		"range literal": {
			input: `
				for i in 5...20
					println(i)
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var i: 5
	_ = l0
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	for {
		if !(value.Bool(value.LessThanEqualInts(l0, (value.SmallInt(20)).ToValue()))) {
			break
		}
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = l0
		callFrame.SetNativeLineNumber(3)
		t1, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		l0 = value.IncrementInt(l0)
	}
}
`,
		},
		"range value": {
			input: `
				r := 5...20
				for i in r
					println(i)
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var range0 = value.NewClosedRange((value.SmallInt(5)).ToValue(), (value.SmallInt(20)).ToValue())
var sym2 = value.ToSymbol("end")
var fn_method0 vm.NativeFunction // Std::ClosedRange.:end
var sym3 = value.ToSymbol("start")
var fn_method1 vm.NativeFunction // Std::ClosedRange.:start
var sym4 = value.ToSymbol("println@1")
var fn_method2 vm.NativeFunction // Std::Kernel::println@1

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
	var l0 *value.ClosedRange // var r: Std::ClosedRange[Std::Int]
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var t3 value.Value
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ClosedRangeClass).LookupMethod(sym2))
	fn_method1 = vm.MethodToFunc((value.ClosedRangeClass).LookupMethod(sym3))
	fn_method2 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym4))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = range0
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (l0).ToValue()
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method0(thread, t2) // receiver: Std::ClosedRange[Std::Int], name: end
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (l0).ToValue()
	t3, err = fn_method1(thread, t2) // receiver: Std::ClosedRange[Std::Int], name: start
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	l1 = t3
	for {
		if !(value.Bool(value.LessThanEqualInts(l1, t1))) {
			break
		}
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = l1
		callFrame.SetNativeLineNumber(4)
		t1, err = fn_method2(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		l1 = value.IncrementInt(l1)
	}
}
`,
		},
		"iterate": {
			input: `
				for i in [1, 2, 3]
					println(i)
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var l0 value.Value // var i: Std::Int
	_ = l0
	var t2 []value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for t1, err = range vm.Iterate(thread, (value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue())).ToValue()) {
		l0 = t1
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = l0
		callFrame.SetNativeLineNumber(3)
		t1, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
	}
}
`,
		},
		"with a pattern": {
			input: `
				for %[a, b] in %[%[1, 2], %[3, 4], %[5, 6]]
					println(a + b)
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var arrtuple0 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue())
var arrtuple1 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue())
var arrtuple2 = value.NewArrayTupleOfValueWithElements(0, (value.SmallInt(5)).ToValue(), (value.SmallInt(6)).ToValue())
var arrtuple3 = value.NewNativeArrayTupleWithElements[*value.ArrayTupleOfValue](0, arrtuple0, arrtuple1, arrtuple2)
var sym2 = value.ToSymbol("length")
var fn_method0 vm.NativeFunction // Std::ArrayTuple.:length
var sym3 = value.ToSymbol("println@1")
var fn_method1 vm.NativeFunction // Std::Kernel::println@1

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
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var t2 value.Bool
	_ = t2
	var t3 []value.Value
	_ = t3
	var t4 value.SmallInt
	_ = t4
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.ArrayTupleClass).LookupMethod(sym2))
	fn_method1 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for t1, err = range vm.Iterate(thread, (arrtuple3).ToValue()) {
		t2 = value.True
		if !(value.Bool(value.IsA(t1, value.TupleMixin))) {
			t2 = value.False
			goto lbl1
		}
		t3 = value.ResizeNativeArgs(t3, 2)
		t3[0] = t1
		callFrame.SetNativeLineNumber(2)
		t1, err = fn_method0(thread, t3) // receiver: Std::ArrayTuple[Std::Int], name: length
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		t4 = value.SmallInt((t1).AsInt())
		if t4 != 2 {
			t2 = value.False
			goto lbl1
		}
		t1, err = ((t1).AsReference().(value.ArrayTuple)).SubscriptInt(int(value.SmallInt(0)))
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		l0 = t1
		t1, err = ((t1).AsReference().(value.ArrayTuple)).SubscriptInt(int(value.SmallInt(1)))
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		l1 = t1
	lbl1:
		if !(t2) {
			thread.CaptureStackTrace()
			thread.Panic((value.NewPatternNotMatchedInForInLoopError()).ToValue())
		}
		t3 = value.ResizeNativeArgs(t3, 3)
		t3[0] = (value.KernelModule).ToValue()
		t3[1] = value.AddInts(l0, l1)
		callFrame.SetNativeLineNumber(3)
		t1, err = fn_method1(thread, t3) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
	}
}
`,
		},
		"with break": {
			input: `
				for i in [1, 2, 3, 4, 5]
					println(i)
					break if i > 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var l0 value.Value // var i: Std::Int
	_ = l0
	var t2 []value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
loop0:
	for t1, err = range vm.Iterate(thread, (value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue())).ToValue()) {
		l0 = t1
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = l0
		callFrame.SetNativeLineNumber(3)
		_, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(2)).ToValue())) {
			break loop0
		} else {
			t1 = value.Nil
		}
	}
}
`,
		},
		"with break with value": {
			input: `
				for i in [1, 2, 3, 4, 5]
					println(i)
					break :foo if i > 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
var fn_method0 vm.NativeFunction // Std::Kernel::println@1
var sym3 = value.ToSymbol("foo")

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
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var l0 value.Value // var i: Std::Int
	_ = l0
	var t2 []value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
loop0:
	for t1, err = range vm.Iterate(thread, (value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue())).ToValue()) {
		l0 = t1
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = l0
		callFrame.SetNativeLineNumber(3)
		_, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(2)).ToValue())) {
			break loop0
		} else {
			t1 = value.Nil
		}
	}
}
`,
		},
		"with break with action": {
			input: `
				a := 5
				for i in [1, 2, 3, 4, 5]
					println(i)
					break a = 20 if i > 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var l1 value.Value // var i: Std::Int
	_ = l1
	var t2 []value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
loop0:
	for t1, err = range vm.Iterate(thread, (value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue())).ToValue()) {
		l1 = t1
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = l1
		callFrame.SetNativeLineNumber(4)
		_, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		if value.Bool(value.GreaterThanInts(l1, (value.SmallInt(2)).ToValue())) {
			l0 = (value.SmallInt(20)).ToValue()
			break loop0
		} else {
			t1 = value.Nil
		}
	}
}
`,
		},
		"with labeled break": {
			input: `
				$foo: for i in [1, 2, 3, 4, 5]
					println(i)
					break[foo] if i > 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var l0 value.Value // var i: Std::Int
	_ = l0
	var t2 []value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
loop0:
	for t1, err = range vm.Iterate(thread, (value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue())).ToValue()) {
		l0 = t1
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = l0
		callFrame.SetNativeLineNumber(3)
		_, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(2)).ToValue())) {
			break loop0
		} else {
			t1 = value.Nil
		}
	}
}
`,
		},
		"with continue": {
			input: `
				for i in [1, 2, 3, 4, 5]
					continue if i > 2
					println(i)
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var l0 value.Value // var i: Std::Int
	_ = l0
	var t2 []value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
loop0:
	for t1, err = range vm.Iterate(thread, (value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue())).ToValue()) {
		l0 = t1
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(2)).ToValue())) {
			continue loop0
		}
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = l0
		callFrame.SetNativeLineNumber(4)
		t1, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
	}
}
`,
		},
		"with labeled continue": {
			input: `
				$foo: for i in [1, 2, 3, 4, 5]
					continue[foo] if i > 2
					println(i)
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var l0 value.Value // var i: Std::Int
	_ = l0
	var t2 []value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
loop0:
	for t1, err = range vm.Iterate(thread, (value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue())).ToValue()) {
		l0 = t1
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(2)).ToValue())) {
			continue loop0
		}
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = l0
		callFrame.SetNativeLineNumber(4)
		t1, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
	}
}
`,
		},
		"nested with break": {
			input: `
				for c in ['a', 'b', 'c', 'd']
					for i in [1, 2, 3, 4, 5]
						break if i > 2
						println(c, i)
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

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println")
var fn_method0 vm.NativeFunction // Std::Kernel::println

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
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var l0 value.Value // var c: Std::String
	_ = l0
	var t2 value.Value
	_ = t2
	var l1 value.Value // var i: Std::Int
	_ = l1
	var t3 []value.Value
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for t1, err = range vm.Iterate(thread, (value.NewNativeArrayListWithElements[value.String](0, value.String("a"), value.String("b"), value.String("c"), value.String("d"))).ToValue()) {
		l0 = t1
		t1 = value.Nil
	loop1:
		for t2, err = range vm.Iterate(thread, (value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue())).ToValue()) {
			l1 = t2
			if value.Bool(value.GreaterThanInts(l1, (value.SmallInt(2)).ToValue())) {
				t1 = value.Nil
				break loop1
			}
			t3 = value.ResizeNativeArgs(t3, 3)
			t3[0] = (value.KernelModule).ToValue()
			t3[1] = (value.NewArrayTupleOfValueWithElementsAndTotalCapacity(2, l0, l1)).ToValue()
			callFrame.SetNativeLineNumber(5)
			t2, err = fn_method0(thread, t3) // receiver: Std::Kernel, name: println
			if err.IsNotUndefined() {
				thread.CaptureStackTrace()
				thread.Panic(err)
			}
		}
	}
}
`,
		},
		"nested with labeled break": {
			input: `
				$foo: for c in ['a', 'b', 'c', 'd']
					for i in [1, 2, 3, 4, 5]
						break[foo] if i > 2
						println(c, i)
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

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println")
var fn_method0 vm.NativeFunction // Std::Kernel::println

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
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var l0 value.Value // var c: Std::String
	_ = l0
	var t2 value.Value
	_ = t2
	var l1 value.Value // var i: Std::Int
	_ = l1
	var t3 []value.Value
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
loop0:
	for t1, err = range vm.Iterate(thread, (value.NewNativeArrayListWithElements[value.String](0, value.String("a"), value.String("b"), value.String("c"), value.String("d"))).ToValue()) {
		l0 = t1
		t1 = value.Nil
		for t2, err = range vm.Iterate(thread, (value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue())).ToValue()) {
			l1 = t2
			if value.Bool(value.GreaterThanInts(l1, (value.SmallInt(2)).ToValue())) {
				break loop0
			}
			t3 = value.ResizeNativeArgs(t3, 3)
			t3[0] = (value.KernelModule).ToValue()
			t3[1] = (value.NewArrayTupleOfValueWithElementsAndTotalCapacity(2, l0, l1)).ToValue()
			callFrame.SetNativeLineNumber(5)
			t2, err = fn_method0(thread, t3) // receiver: Std::Kernel, name: println
			if err.IsNotUndefined() {
				thread.CaptureStackTrace()
				thread.Panic(err)
			}
		}
	}
}
`,
		},
		"nested with continue": {
			input: `
				for c in ['a', 'b', 'c', 'd']
					for i in [1, 2, 3, 4, 5]
						continue if i > 2
						println(c, i)
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

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println")
var fn_method0 vm.NativeFunction // Std::Kernel::println

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
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var l0 value.Value // var c: Std::String
	_ = l0
	var t2 value.Value
	_ = t2
	var l1 value.Value // var i: Std::Int
	_ = l1
	var t3 []value.Value
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for t1, err = range vm.Iterate(thread, (value.NewNativeArrayListWithElements[value.String](0, value.String("a"), value.String("b"), value.String("c"), value.String("d"))).ToValue()) {
		l0 = t1
		t1 = value.Nil
	loop1:
		for t2, err = range vm.Iterate(thread, (value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue())).ToValue()) {
			l1 = t2
			if value.Bool(value.GreaterThanInts(l1, (value.SmallInt(2)).ToValue())) {
				t1 = value.Nil
				continue loop1
			}
			t3 = value.ResizeNativeArgs(t3, 3)
			t3[0] = (value.KernelModule).ToValue()
			t3[1] = (value.NewArrayTupleOfValueWithElementsAndTotalCapacity(2, l0, l1)).ToValue()
			callFrame.SetNativeLineNumber(5)
			t2, err = fn_method0(thread, t3) // receiver: Std::Kernel, name: println
			if err.IsNotUndefined() {
				thread.CaptureStackTrace()
				thread.Panic(err)
			}
		}
	}
}
`,
		},
		"nested with labeled continue": {
			input: `
				$foo: for c in ['a', 'b', 'c', 'd']
					for i in [1, 2, 3, 4, 5]
						continue[foo] if i > 2
						println(c, i)
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

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println")
var fn_method0 vm.NativeFunction // Std::Kernel::println

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
	var t1 value.Value
	_ = t1
	var err value.Value
	_ = err
	var l0 value.Value // var c: Std::String
	_ = l0
	var t2 value.Value
	_ = t2
	var l1 value.Value // var i: Std::Int
	_ = l1
	var t3 []value.Value
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
loop0:
	for t1, err = range vm.Iterate(thread, (value.NewNativeArrayListWithElements[value.String](0, value.String("a"), value.String("b"), value.String("c"), value.String("d"))).ToValue()) {
		l0 = t1
		t1 = value.Nil
		for t2, err = range vm.Iterate(thread, (value.NewArrayListOfValueWithElements(0, (value.SmallInt(1)).ToValue(), (value.SmallInt(2)).ToValue(), (value.SmallInt(3)).ToValue(), (value.SmallInt(4)).ToValue(), (value.SmallInt(5)).ToValue())).ToValue()) {
			l1 = t2
			if value.Bool(value.GreaterThanInts(l1, (value.SmallInt(2)).ToValue())) {
				continue loop0
			}
			t3 = value.ResizeNativeArgs(t3, 3)
			t3[0] = (value.KernelModule).ToValue()
			t3[1] = (value.NewArrayTupleOfValueWithElementsAndTotalCapacity(2, l0, l1)).ToValue()
			callFrame.SetNativeLineNumber(5)
			t2, err = fn_method0(thread, t3) // receiver: Std::Kernel, name: println
			if err.IsNotUndefined() {
				thread.CaptureStackTrace()
				thread.Panic(err)
			}
		}
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

func TestGoReturnExpression(t *testing.T) {
	tests := goTestTable{
		"return a value in main": {
			input: "return 5.to_string",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("to_string")
var fn_method0 vm.NativeFunction // Std::Int.:to_string

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
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc((value.IntClass).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t2 = value.ResizeNativeArgs(t2, 2)
	t2[0] = (value.SmallInt(5)).ToValue()
	t1, err = fn_method0(thread, t2) // receiver: Std::Int, name: to_string
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	return
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(7, 1, 8), P(17, 1, 18)), "values returned in void context will be ignored"),
			},
		},
		// TODO: return in methods, in namespace bodies
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

// func TestBytecodeAwaitExpression(t *testing.T) {
// 	tests := bytecodeTestTable{
// 		"await in a synchronous context": {
// 			input: "await timeout(2.seconds)",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.GET_CONST8), 0,
// 					byte(bytecode.INT_2),
// 					byte(bytecode.CALL_METHOD8), 1,
// 					byte(bytecode.UNDEFINED),
// 					byte(bytecode.CALL_METHOD8), 2,
// 					byte(bytecode.AWAIT_SYNC),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(23, 1, 24)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 10),
// 				},
// 				[]value.Value{
// 					value.ToSymbol("Std::Kernel").ToValue(),
// 					value.Ref(vm.NewCallSiteInfo(value.ToSymbol("seconds"), 0)),
// 					value.Ref(vm.NewCallSiteInfo(value.ToSymbol("timeout"), 2)),
// 				},
// 			),
// 		},
// 		"await in an asynchronous context": {
// 			input: `
// 				async def foo
// 					await timeout(2.seconds)
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(56, 4, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 					bytecode.NewLineInfo(4, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						methodDefinitionsSymbol,
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.GET_SINGLETON),
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(56, 4, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 7),
// 							bytecode.NewLineInfo(4, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Std::Kernel").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("foo"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.PROMISE),
// 									byte(bytecode.RETURN),
// 									byte(bytecode.INT_2),
// 									byte(bytecode.CALL_METHOD8), 0,
// 									byte(bytecode.UNDEFINED),
// 									byte(bytecode.CALL_SELF8), 1,
// 									byte(bytecode.AWAIT),
// 									byte(bytecode.AWAIT_RESULT),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(5, 2, 5), P(55, 4, 7)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(2, 2),
// 									bytecode.NewLineInfo(4, 1),
// 									bytecode.NewLineInfo(3, 8),
// 									bytecode.NewLineInfo(4, 1),
// 								},
// 								1,
// 								1,
// 								[]value.Value{
// 									value.Ref(vm.NewCallSiteInfo(value.ToSymbol("seconds"), 0)),
// 									value.Ref(vm.NewCallSiteInfo(value.ToSymbol("timeout"), 2)),
// 								},
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"await_sync in an asynchronous context": {
// 			input: `
// 				async def foo
// 					await_sync timeout(2.seconds)
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(61, 4, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 					bytecode.NewLineInfo(4, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						methodDefinitionsSymbol,
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.GET_SINGLETON),
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_METHOD),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(61, 4, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 7),
// 							bytecode.NewLineInfo(4, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Std::Kernel").ToValue(),
// 							value.Ref(vm.NewBytecodeFunction(
// 								value.ToSymbol("foo"),
// 								[]byte{
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.PROMISE),
// 									byte(bytecode.RETURN),
// 									byte(bytecode.INT_2),
// 									byte(bytecode.CALL_METHOD8), 0,
// 									byte(bytecode.UNDEFINED),
// 									byte(bytecode.CALL_SELF8), 1,
// 									byte(bytecode.AWAIT_SYNC),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(5, 2, 5), P(60, 4, 7)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(2, 2),
// 									bytecode.NewLineInfo(4, 1),
// 									bytecode.NewLineInfo(3, 7),
// 									bytecode.NewLineInfo(4, 1),
// 								},
// 								1,
// 								1,
// 								[]value.Value{
// 									value.Ref(vm.NewCallSiteInfo(value.ToSymbol("seconds"), 0)),
// 									value.Ref(vm.NewCallSiteInfo(value.ToSymbol("timeout"), 2)),
// 								},
// 							)),
// 							value.ToSymbol("foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			bytecodeCompilerTest(tc, t)
// 		})
// 	}
// }

// func TestBytecodeModifierForIn(t *testing.T) {
// 	tests := bytecodeTestTable{
// 		"iterate": {
// 			input: `println(i) for i in [1, 2, 3]`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 2,
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.COPY),
// 					byte(bytecode.GET_ITERATOR),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.FOR_IN_BUILTIN), 0, 10,
// 					byte(bytecode.SET_LOCAL_2),
// 					byte(bytecode.GET_CONST8), 1,
// 					byte(bytecode.GET_LOCAL_2),
// 					byte(bytecode.CALL_METHOD8), 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.LOOP), 0, 14,
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(28, 1, 29)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 22),
// 				},
// 				[]value.Value{
// 					value.Ref(&value.ArrayList{
// 						value.SmallInt(1).ToValue(),
// 						value.SmallInt(2).ToValue(),
// 						value.SmallInt(3).ToValue(),
// 					}),
// 					value.ToSymbol("Std::Kernel").ToValue(),
// 					value.Ref(vm.NewCallSiteInfo(
// 						value.ToSymbol("println@1"),
// 						1,
// 					)),
// 				},
// 			),
// 		},
// 		"with a pattern": {
// 			input: `println(a + b) for %[a, b] in %[%[1, 2], %[3, 4], %[5, 6]]`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.PREP_LOCALS8), 3,
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.GET_ITERATOR),
// 					byte(bytecode.SET_LOCAL_1),
// 					byte(bytecode.GET_LOCAL_1),
// 					byte(bytecode.FOR_IN_BUILTIN), 0, 56,
// 					byte(bytecode.DUP),
// 					byte(bytecode.LOAD_VALUE_1),
// 					byte(bytecode.IS_A),
// 					byte(bytecode.JUMP_UNLESS_NP), 0, 33,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.CALL_METHOD8), 2,
// 					byte(bytecode.INT_2),
// 					byte(bytecode.EQUAL_INT),
// 					byte(bytecode.JUMP_UNLESS_NP), 0, 24,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.INT_0),
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.SET_LOCAL_2),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS_NP), 0, 13,
// 					byte(bytecode.POP),
// 					byte(bytecode.DUP),
// 					byte(bytecode.INT_1),
// 					byte(bytecode.SUBSCRIPT),
// 					byte(bytecode.DUP),
// 					byte(bytecode.SET_LOCAL_3),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.POP_SKIP_ONE),
// 					byte(bytecode.JUMP_UNLESS_NP), 0, 2,
// 					byte(bytecode.POP),
// 					byte(bytecode.TRUE),
// 					byte(bytecode.JUMP_IF), 0, 2,
// 					byte(bytecode.LOAD_VALUE_3),
// 					byte(bytecode.THROW),
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_CONST8), 4,
// 					byte(bytecode.GET_LOCAL_2),
// 					byte(bytecode.GET_LOCAL_3),
// 					byte(bytecode.ADD_INT),
// 					byte(bytecode.CALL_METHOD8), 5,
// 					byte(bytecode.POP),
// 					byte(bytecode.LOOP), 0, 60,
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(57, 1, 58)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 67),
// 				},
// 				[]value.Value{
// 					value.Ref(&value.ArrayTuple{
// 						value.Ref(&value.ArrayTuple{
// 							value.SmallInt(1).ToValue(),
// 							value.SmallInt(2).ToValue(),
// 						}),
// 						value.Ref(&value.ArrayTuple{
// 							value.SmallInt(3).ToValue(),
// 							value.SmallInt(4).ToValue(),
// 						}),
// 						value.Ref(&value.ArrayTuple{
// 							value.SmallInt(5).ToValue(),
// 							value.SmallInt(6).ToValue(),
// 						}),
// 					}),
// 					value.Ref(value.TupleMixin),
// 					value.Ref(vm.NewCallSiteInfo(value.ToSymbol("length"), 0)),
// 					value.Ref(value.NewError(
// 						value.PatternNotMatchedErrorClass,
// 						"assigned value does not match the pattern defined in for in loop",
// 					)),
// 					value.ToSymbol("Std::Kernel").ToValue(),
// 					value.Ref(vm.NewCallSiteInfo(
// 						value.ToSymbol("println@1"),
// 						1,
// 					)),
// 				},
// 			),
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			bytecodeCompilerTest(tc, t)
// 		})
// 	}
// }

func TestGoIfExpression(t *testing.T) {
	tests := goTestTable{
		"resolve static condition with empty then and else": {
			input: `a := if false; end`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: nil
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Nil
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(8, 1, 9), P(12, 1, 13)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"empty then and else": {
			input: "a := true; b := if a; end",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Bool // var a: bool
	_ = l0
	var l1 value.Value // var b: nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.True
	if l0 {
		t1 = value.Nil
	} else {
		t1 = value.Nil
	}
	l1 = t1
}
`,
		},
		"resolve static condition with then branch": {
			input: `
				a := if true
					10
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(10)).ToValue()
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(13, 2, 13), P(16, 2, 16)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"resolve static condition with then branch to nil": {
			input: `
				a := if false
					10
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: nil
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Nil
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(13, 2, 13), P(17, 2, 17)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewWarning(L(P(24, 3, 6), P(26, 3, 8)), "unreachable code"),
			},
		},
		"resolve static condition with then and else branches": {
			input: `
				a := if false
					10
				else
					5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(13, 2, 13), P(17, 2, 17)), "this condition will always have the same result since type `false` is falsy"),
				diagnostic.NewWarning(L(P(24, 3, 6), P(26, 3, 8)), "unreachable code"),
			},
		},
		"with then branch": {
			input: `
				var a: Int? = 5
				if a
					a = a * 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	if value.Truthy(l0) {
		l0 = value.MultiplyInts(l0, (value.SmallInt(2)).ToValue())
	}
}
`,
		},
		"with then branch and value": {
			input: `
				var a: Int? = 5
				result := if a
					a = a * 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.Value // var result: Std::Int | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	if value.Truthy(l0) {
		l0 = value.MultiplyInts(l0, (value.SmallInt(2)).ToValue())
		t1 = l0
	} else {
		t1 = value.Nil
	}
	l1 = t1
}
`,
		},
		"with then and else branches": {
			input: `
				var a: Int? = 5
				if a
					a = a * 2
				else
					a = 30
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	if value.Truthy(l0) {
		l0 = value.MultiplyInts(l0, (value.SmallInt(2)).ToValue())
	} else {
		l0 = (value.SmallInt(30)).ToValue()
	}
}
`,
		},
		"with then and else branches and value": {
			input: `
				var a: Int? = 5
				result := if a
					a = a * 2
				else
					a = 30
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.Value // var result: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	if value.Truthy(l0) {
		l0 = value.MultiplyInts(l0, (value.SmallInt(2)).ToValue())
		t1 = l0
	} else {
		l0 = (value.SmallInt(30)).ToValue()
		t1 = l0
	}
	l1 = t1
}
`,
		},
		"with native bool condition": {
			input: `
				a := 5
				b := false
				if b
					a = a * 2
				else
					a = 30
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	l1 = value.False
	if l1 {
		l0 = value.MultiplyInts(l0, (value.SmallInt(2)).ToValue())
	} else {
		l0 = (value.SmallInt(30)).ToValue()
	}
}
`,
		},
		"with native type narrowing": {
			input: `
				var a: Float? = 10.5
				if a
					a = a * 2.0
				else
					a = 69.420
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Float?
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.Float(10.5)).ToValue()
	if value.Truthy(l0) {
		l0 = (((l0).AsFloat()).MultiplyFloat(value.Float(2))).ToValue()
	} else {
		l0 = (value.Float(69.42)).ToValue()
	}
}
`,
		},
		"is an expression": {
			input: `
				var a: Int? = 5
				b := if a
					"foo"
				else
					5
				end
				b
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.Value // var b: Std::String | Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	if value.Truthy(l0) {
		t1 = (value.String("foo")).ToValue()
	} else {
		t1 = (value.SmallInt(5)).ToValue()
	}
	l1 = t1
}
`,
		},
		"can be chained": {
			input: `
				a := 5.6
				var result: String
				if a <= 2.0
					result = "foo"
				else if a <= 12.9
					result = "bar"
				else
				  result = "baz"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l1 value.String // var result: Std::String
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(5.6)
	if value.Bool((l0).LessThanEqualFloat(value.Float(2))) {
		l1 = value.String("foo")
	} else {
		if value.Bool((l0).LessThanEqualFloat(value.Float(12.9))) {
			l1 = value.String("bar")
		} else {
			l1 = value.String("baz")
		}
	}
}
`,
		},
		"can be chained with value": {
			input: `
				a := 5.6
				result := if a <= 2.0
					"foo"
				else if a <= 12.9
					"bar"
				else
				  "baz"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l1 value.String // var result: Std::String
	_ = l1
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(5.6)
	if value.Bool((l0).LessThanEqualFloat(value.Float(2))) {
		t1 = (value.String("foo")).ToValue()
	} else {
		if value.Bool((l0).LessThanEqualFloat(value.Float(12.9))) {
			t2 = (value.String("bar")).ToValue()
		} else {
			t2 = (value.String("baz")).ToValue()
		}
		t1 = t2
	}
	l1 = (t1).AsString()
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

func TestGoUnlessExpression(t *testing.T) {
	tests := goTestTable{
		"resolve static condition with empty then and else": {
			input: "a := unless true; end",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: nil
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Nil
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(12, 1, 13), P(15, 1, 16)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"empty then and else": {
			input: "a := true; b := unless a; end",
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Bool // var a: bool
	_ = l0
	var l1 value.Value // var b: nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.True
	if !(l0) {
		t1 = value.Nil
	} else {
		t1 = value.Nil
	}
	l1 = t1
}
`,
		},
		"resolve static condition with then branch": {
			input: `
				a := unless false
					10
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(10)).ToValue()
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(17, 2, 17), P(21, 2, 21)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"resolve static condition with then branch to nil": {
			input: `
				a := unless true
					10
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: nil
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Nil
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(17, 2, 17), P(20, 2, 20)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewWarning(L(P(27, 3, 6), P(29, 3, 8)), "unreachable code"),
			},
		},
		"resolve static condition with then and else branches": {
			input: `
				a := unless true
					10
				else
					5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(17, 2, 17), P(20, 2, 20)), "this condition will always have the same result since type `true` is truthy"),
				diagnostic.NewWarning(L(P(27, 3, 6), P(29, 3, 8)), "unreachable code"),
			},
		},
		"with then branch": {
			input: `
				var a: Int? = 5
				unless a
					a = 30
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	if value.Falsy(l0) {
		l0 = (value.SmallInt(30)).ToValue()
	}
}
`,
		},
		"with then branch and value": {
			input: `
				var a: Int? = 5
				result := unless a
					a = 30
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.Value // var result: Std::Int | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	if value.Falsy(l0) {
		l0 = (value.SmallInt(30)).ToValue()
		t1 = l0
	} else {
		t1 = value.Nil
	}
	l1 = t1
}
`,
		},
		"with then and else branches": {
			input: `
				var a: Int? = 5
				unless a
					a = 30
				else
					a = a * 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	if value.Falsy(l0) {
		l0 = (value.SmallInt(30)).ToValue()
	} else {
		l0 = value.MultiplyInts(l0, (value.SmallInt(2)).ToValue())
	}
}
`,
		},
		"with then and else branches and value": {
			input: `
				var a: Int? = 5
				result := unless a
					a = 30
				else
					a = a * 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.Value // var result: Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	if value.Falsy(l0) {
		l0 = (value.SmallInt(30)).ToValue()
		t1 = l0
	} else {
		l0 = value.MultiplyInts(l0, (value.SmallInt(2)).ToValue())
		t1 = l0
	}
	l1 = t1
}
`,
		},
		"with native bool condition": {
			input: `
				a := 5
				b := false
				unless b
					a = 30
				else
					a = a * 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	l1 = value.False
	if !(l1) {
		l0 = (value.SmallInt(30)).ToValue()
	} else {
		l0 = value.MultiplyInts(l0, (value.SmallInt(2)).ToValue())
	}
}
`,
		},
		"is an expression": {
			input: `
				var a: Int? = 5
				b := unless a
					"foo"
				else
					5
				end
				b
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int?
	_ = l0
	var l1 value.Value // var b: Std::String | Std::Int
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	if value.Falsy(l0) {
		t1 = (value.String("foo")).ToValue()
	} else {
		t1 = (value.SmallInt(5)).ToValue()
	}
	l1 = t1
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

func TestGoBreak(t *testing.T) {
	tests := goTestTable{
		"in top level": {
			input: "break",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(4, 1, 5)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"in top level with a label": {
			input: "break[foo]",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(8, 1, 9)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"nonexistent label": {
			input: `
				loop
					break[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(15, 3, 6), P(23, 3, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to an expression": {
			input: `
				loop
					$foo: 1 + 2
					break[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(32, 4, 6), P(40, 4, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to a different loop": {
			input: `
				$foo: loop
					println("foo")
				end

				loop
					break[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(59, 7, 6), P(67, 7, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
				diagnostic.NewWarning(L(P(49, 6, 5), P(76, 8, 7)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoContinue(t *testing.T) {
	tests := goTestTable{
		"in top level": {
			input: "continue",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(7, 1, 8)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"in top level with a label": {
			input: "continue[foo]",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(11, 1, 12)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"nonexistent label": {
			input: `
				loop
					continue[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(15, 3, 6), P(26, 3, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to an expression": {
			input: `
				loop
					$foo: 1 + 2
					continue[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(32, 4, 6), P(43, 4, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"label attached to a different loop": {
			input: `
				$foo: loop
					println("foo")
				end

				loop
					continue[foo]
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(59, 7, 6), P(70, 7, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
				diagnostic.NewWarning(L(P(49, 6, 5), P(79, 8, 7)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoLoopExpression(t *testing.T) {
	tests := goTestTable{
		"empty body": {
			input: `
				loop
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for {
	}
}
`,
		},
		"with a body": {
			input: `
				a := 0
				loop
					a = a + 1
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
	}
}
`,
		},
		"with continue": {
			input: `
				a := 0
				loop
					a = a + 1
					continue
					println("foo")
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		continue loop0
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(55, 6, 6), P(68, 6, 19)), "unreachable code"),
			},
		},
		"with labeled continue": {
			input: `
				a := 0
				$foo: loop
					a = a + 1
					continue[foo]
					println("foo")
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		continue loop0
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(66, 6, 6), P(79, 6, 19)), "unreachable code"),
			},
		},
		"continue in a nested loop": {
			input: `
			 	j := 0
				loop
					j += 1
					i := 0
					loop
						continue if i >= 5
						i += 1
					end
					continue if j >= 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
	loop1:
		for {
			if value.Bool(value.GreaterThanEqualInts(l1, (value.SmallInt(5)).ToValue())) {
				continue loop1
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			continue loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(108, 10, 6), P(125, 10, 23)), "unreachable code"),
			},
		},

		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: loop
					j += 1
					i := 0
					loop
						continue[foo] if i + j > 8
						i += 1
					end
					continue if j >= 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				continue loop0
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			continue loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(122, 10, 6), P(139, 10, 23)), "unreachable code"),
			},
		},
		"labeled continue with value in a nested loop": {
			input: `
			 	j := 0
				result := $foo: loop
					j += 1
					i := 0
					loop
						continue[foo](:bar) if i + j > 8
						i += 1
					end
					continue if j >= 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("bar")

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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var result: never
	_ = l1
	var t1 value.Value
	_ = t1
	var l2 value.Value // var i: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l2 = (value.SmallInt(0)).ToValue()
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l2, l0), (value.SmallInt(8)).ToValue())) {
				t1 = (sym2).ToValue()
				continue loop0
			}
			l2 = value.AddInts(l2, (value.SmallInt(1)).ToValue())
		}
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			t1 = value.Nil
			continue loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(138, 10, 6), P(155, 10, 23)), "unreachable code"),
			},
		},
		"with break": {
			input: `
				a := 0
				loop
					a = a + 1
					break if a > 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
		},
		"with a labeled break": {
			input: `
				a := 0
				$foo: loop
					a = a + 1
					break[foo] if a > 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
		},
		"break in a nested loop": {
			input: `
			 	j := 0
				loop
					j += 1
					i := 0
					loop
						break if i >= 5
						i += 1
					end
					break if j >= 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
	loop1:
		for {
			if value.Bool(value.GreaterThanEqualInts(l1, (value.SmallInt(5)).ToValue())) {
				break loop1
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$outer: loop
					j += 1
					i := 0
					loop
						break[outer] if i >= 5
						i += 1
					end
					break if j >= 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
		for {
			if value.Bool(value.GreaterThanEqualInts(l1, (value.SmallInt(5)).ToValue())) {
				break loop0
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(120, 10, 6), P(134, 10, 20)), "unreachable code"),
			},
		},
		"labeled break with value in a nested loop": {
			input: `
			 	j := 0
				result := ($outer: loop
					j += 1
					i := 0
					loop
						break[outer](:bar) if i >= 5
						i += 1
					end
					break if j >= 5
				end)
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("bar")

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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var result: Std::Symbol | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var l2 value.Value // var i: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l2 = (value.SmallInt(0)).ToValue()
		for {
			if value.Bool(value.GreaterThanEqualInts(l2, (value.SmallInt(5)).ToValue())) {
				t1 = (sym2).ToValue()
				break loop0
			}
			l2 = value.AddInts(l2, (value.SmallInt(1)).ToValue())
		}
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			t1 = value.Nil
			break loop0
		}
	}
	l1 = t1
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(137, 10, 6), P(151, 10, 20)), "unreachable code"),
			},
		},
		"break with value": {
			input: `
				a := 0
				result := loop
					a = a + 1
					break true if a > 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Bool // var result: bool
	_ = l1
	var t1 value.Bool
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(5)).ToValue())) {
			t1 = value.True
			break loop0
		}
	}
	l1 = t1
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

func TestGoLogicalOrOperator(t *testing.T) {
	tests := goTestTable{
		"simple static left truthy": {
			input: `
				a := "foo"
				b := a || true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.String // var a: Std::String
	_ = l0
	var l1 value.String // var b: Std::String
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.String("foo")
	l1 = l0
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(25, 3, 10), P(25, 3, 10)), "this condition will always have the same result since type `Std::String` is truthy"),
				diagnostic.NewWarning(L(P(30, 3, 15), P(33, 3, 18)), "unreachable code"),
			},
		},
		"simple static left falsy": {
			input: `
				a := nil
				b := a || true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: nil
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Nil
	l1 = value.True
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(23, 3, 10), P(23, 3, 10)), "this condition will always have the same result since type `nil` is falsy"),
			},
		},
		"simple dynamic": {
			input: `
				var a: String? = "foo"
				b := a || true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::String?
	_ = l0
	var l1 value.Value // var b: Std::String | bool
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.String("foo")).ToValue()
	t1 = l0
	if value.Falsy(t1) {
		t1 = (value.True).ToValue()
	}
	l1 = t1
}
`,
		},
		"nested static": {
			input: `
				a := "foo"
				b := a || true || 3
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.String // var a: Std::String
	_ = l0
	var l1 value.String // var b: Std::String
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.String("foo")
	l1 = l0
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(25, 3, 10), P(25, 3, 10)), "this condition will always have the same result since type `Std::String` is truthy"),
				diagnostic.NewWarning(L(P(30, 3, 15), P(33, 3, 18)), "unreachable code"),
				diagnostic.NewWarning(L(P(25, 3, 10), P(33, 3, 18)), "this condition will always have the same result since type `Std::String` is truthy"),
				diagnostic.NewWarning(L(P(38, 3, 23), P(38, 3, 23)), "unreachable code"),
			},
		},
		"nested dynamic": {
			input: `
				var a: String? = "foo"
				var b: Int? = 5
				c := a || b || 3
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::String?
	_ = l0
	var l1 value.Value // var b: Std::Int?
	_ = l1
	var l2 value.Value // var c: Std::Int | Std::String
	_ = l2
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.String("foo")).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	t1 = l0
	if value.Falsy(t1) {
		t1 = l1
	}
	t2 = t1
	if value.Falsy(t2) {
		t2 = (value.SmallInt(3)).ToValue()
	}
	l2 = t2
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

func TestGoLogicalAndOperator(t *testing.T) {
	tests := goTestTable{
		"simple static left truthy": {
			input: `
				a := "foo"
				b := a && true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.String // var a: Std::String
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.String("foo")
	l1 = value.True
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(25, 3, 10), P(25, 3, 10)), "this condition will always have the same result since type `Std::String` is truthy"),
			},
		},
		"simple static left falsy": {
			input: `
				a := nil
				b := a && true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: nil
	_ = l0
	var l1 value.Value // var b: nil
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Nil
	l1 = l0
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(23, 3, 10), P(23, 3, 10)), "this condition will always have the same result since type `nil` is falsy"),
				diagnostic.NewWarning(L(P(28, 3, 15), P(31, 3, 18)), "unreachable code"),
			},
		},
		"simple dynamic": {
			input: `
				var a: String? = "foo"
				b := a && true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::String?
	_ = l0
	var l1 value.Value // var b: nil | bool
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.String("foo")).ToValue()
	t1 = l0
	if value.Truthy(t1) {
		t1 = (value.True).ToValue()
	}
	l1 = t1
}
`,
		},
		"nested": {
			input: `
				var a: String? = "foo"
				var b: Int? = 5
				c := a && b && 3
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::String?
	_ = l0
	var l1 value.Value // var b: Std::Int?
	_ = l1
	var l2 value.Value // var c: nil | Std::Int
	_ = l2
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.String("foo")).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	t1 = l0
	if value.Truthy(t1) {
		t1 = l1
	}
	t2 = t1
	if value.Truthy(t2) {
		t2 = (value.SmallInt(3)).ToValue()
	}
	l2 = t2
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

func TestGoNilCoalescingOperator(t *testing.T) {
	tests := goTestTable{
		"simple static left non-nilable": {
			input: `
				a := "foo"
				b := a ?? true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.String // var a: Std::String
	_ = l0
	var l1 value.String // var b: Std::String
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.String("foo")
	l1 = l0
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(25, 3, 10), P(25, 3, 10)), "this condition will always have the same result since type `Std::String` can never be nil"),
				diagnostic.NewWarning(L(P(30, 3, 15), P(33, 3, 18)), "unreachable code"),
			},
		},
		"simple static left nil": {
			input: `
				a := nil
				b := a ?? true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: nil
	_ = l0
	var l1 value.Bool // var b: bool
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Nil
	l1 = value.True
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(23, 3, 10), P(23, 3, 10)), "this condition will always have the same result"),
			},
		},
		"nested static": {
			input: `
				a := "foo"
				b := a ?? true ?? 3
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.String // var a: Std::String
	_ = l0
	var l1 value.String // var b: Std::String
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.String("foo")
	l1 = l0
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(25, 3, 10), P(25, 3, 10)), "this condition will always have the same result since type `Std::String` can never be nil"),
				diagnostic.NewWarning(L(P(30, 3, 15), P(33, 3, 18)), "unreachable code"),
				diagnostic.NewWarning(L(P(25, 3, 10), P(33, 3, 18)), "this condition will always have the same result since type `Std::String` can never be nil"),
				diagnostic.NewWarning(L(P(38, 3, 23), P(38, 3, 23)), "unreachable code"),
			},
		},
		"nested dynamic": {
			input: `
				var a: String? = "foo"
				var b: Int? = 5
				c := a ?? b ?? 35
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::String?
	_ = l0
	var l1 value.Value // var b: Std::Int?
	_ = l1
	var l2 value.Value // var c: Std::Int | Std::String
	_ = l2
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.String("foo")).ToValue()
	l1 = (value.SmallInt(5)).ToValue()
	t1 = l0
	if value.IsNil(t1) {
		t1 = l1
	}
	t2 = t1
	if value.IsNil(t2) {
		t2 = (value.SmallInt(35)).ToValue()
	}
	l2 = t2
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

func TestGoNumericFor(t *testing.T) {
	tests := goTestTable{
		"for without initialiser, condition, increment and body but with value": {
			input: `
				result := fornum ;;
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var result: never
	_ = l0
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for {
	}
}
`,
		},
		"for without initialiser, condition, increment and body": {
			input: `
				fornum ;;
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for {
	}
}
`,
		},
		"for without initialiser, condition and increment": {
			input: `
				a := 0
				fornum ;;
					a += 1
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
	}
}
`,
		},
		"for without initialiser, condition and increment but with value": {
			input: `
				a := 0
				result := fornum ;;
					a += 1
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var result: never
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
	}
}
`,
		},
		"for with break": {
			input: `
				a := 0
				fornum ;;
					a += 1
					break if a > 10
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(10)).ToValue())) {
			break loop0
		}
	}
}
`,
		},
		"for with break and result value": {
			input: `
				a := 0
				result := fornum ;;
					a += 1
					break if a > 10
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var result: nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(10)).ToValue())) {
			t1 = value.Nil
			break loop0
		}
	}
	l1 = t1
}
`,
		},
		"for with labeled break": {
			input: `
				a := 0
				$foo: fornum ;;
					a += 1
					break[foo] if a > 10
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(10)).ToValue())) {
			break loop0
		}
	}
}
`,
		},
		"nested for with continue": {
			input: `
				fornum a := 0;;
					fornum ;; a += 1
						continue if a > 10
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

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
	loop1:
		for {
			if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(10)).ToValue())) {
				continue loop1
			}
			l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},
		"nested for with a labeled continue": {
			input: `
				$foo: fornum a := 0;;
					fornum ;; a += 1
						continue[foo] if a > 10
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

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		for {
			if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(10)).ToValue())) {
				continue loop0
			}
			l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},
		"nested for with break": {
			input: `
				fornum a := 0;;
					fornum ;; a += 1
						break if a > 10
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

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
	loop1:
		for {
			if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(10)).ToValue())) {
				break loop1
			}
			l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},
		"nested for with a labeled break": {
			input: `
				$foo: fornum a := 0;;
					fornum ;; a += 1
						break[foo] if a > 10
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

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		for {
			if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(10)).ToValue())) {
				break loop0
			}
			l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},

		"for with break with value": {
			input: `
				fornum a := 0;;
					a += 1
					break 5 if a > 10
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.GreaterThanInts(l0, (value.SmallInt(10)).ToValue())) {
			break loop0
		}
	}
}
`,
		},
		"for with initialiser, without condition and increment": {
			input: `
				fornum a := 0;;
					a += 1
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
	}
}
`,
		},
		"for with initialiser, condition, without increment": {
			input: `
				fornum a := 0; a < 5;
					a += 1
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
	}
}
`,
		},
		"for with initialiser, condition, without increment but with result value": {
			input: `
				result := fornum a := 0; a < 5;
					a += 1
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var result: Std::Int?
	_ = l0
	var t1 value.Value
	_ = t1
	var l1 value.Value // var a: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t1 = value.Nil
	l1 = (value.SmallInt(0)).ToValue()
	for {
		if !(value.Bool(value.LessThanInts(l1, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		t1 = l1
	}
	l0 = t1
}
`,
		},
		"for with initialiser, condition and increment": {
			input: `
				a := 0
				fornum i := 0; i < 5; i += 1
					a += i
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	l1 = (value.SmallInt(0)).ToValue()
	for {
		if !(value.Bool(value.LessThanInts(l1, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l0 = value.AddInts(l0, l1)
		l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
	}
}
`,
		},
		"for with initialiser, condition, increment and result value": {
			input: `
				a := 0
				result := fornum i := 0; i < 5; i += 1
					a += i
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var result: Std::Int?
	_ = l1
	var t1 value.Value
	_ = t1
	var l2 value.Value // var i: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
	l2 = (value.SmallInt(0)).ToValue()
	for {
		if !(value.Bool(value.LessThanInts(l2, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l0 = value.AddInts(l0, l2)
		t1 = l0
		l2 = value.AddInts(l2, (value.SmallInt(1)).ToValue())
	}
	l1 = t1
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

func TestGoModifierWhile(t *testing.T) {
	tests := goTestTable{
		"single line": {
			input: `
			  i := 0
				i += 1 while i < 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
	}
}
`,
		},
		"multiline": {
			input: `
			  i := 0
				do
					i += 1
				end while i < 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
	}
}
`,
		},
		"with break": {
			input: `
			  i := 0
				do
					i += 1
					break if i < 5
				end while true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(66, 6, 15), P(69, 6, 18)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"with break with result value": {
			input: `
			  i := 0
				result := (do
					i += 1
					break if i < 5
				end while true)
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var l1 value.Value // var result: nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			t1 = value.Nil
			break loop0
		}
	}
	l1 = t1
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(77, 6, 15), P(80, 6, 18)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"with labeled break": {
			input: `
				i := 0
				$foo: do
					i += 1
					break[foo] if i < 5
				end while true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(76, 6, 15), P(79, 6, 18)), "this condition will always have the same result since type `true` is truthy"),
			},
		},

		"with break with value": {
			input: `
				i := 0
				result := (do
					i += 1
					break true if i < 5
				end while true)
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var l1 value.Bool // var result: bool
	_ = l1
	var t1 value.Bool
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			t1 = value.True
			break loop0
		}
	}
	l1 = t1
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(81, 6, 15), P(84, 6, 18)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"continue in a nested loop": {
			input: `
			 	j := 0
				do
					j += 1
					i := 0
					do
						continue if i + j > 8
						i += 1
					end while i < 5
				end while j < 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
	loop1:
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				continue loop1
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
			if !(value.Bool(value.LessThanInts(l1, (value.SmallInt(5)).ToValue()))) {
				break
			}
		}
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
	}
}
`,
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						continue[foo] if i + j > 8
						i += 1
					end while i < 5
				end while j < 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				continue loop0
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
			if !(value.Bool(value.LessThanInts(l1, (value.SmallInt(5)).ToValue()))) {
				break
			}
		}
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
	}
}
`,
		},
		"break in a nested loop": {
			input: `
			 	j := 0
				do
					j += 1
					i := 0
					do
						break if i + j > 8
						i += 1
					end while i < 5
				end while j < 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
	loop1:
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				break loop1
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
			if !(value.Bool(value.LessThanInts(l1, (value.SmallInt(5)).ToValue()))) {
				break
			}
		}
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
	}
}
`,
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						break[foo] if i + j > 8
						i += 1
					end while i < 5
				end while j < 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				break loop0
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
			if !(value.Bool(value.LessThanInts(l1, (value.SmallInt(5)).ToValue()))) {
				break
			}
		}
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
	}
}
`,
		},
		"labeled break with value in a nested loop": {
			input: `
			 	j := 0
				result := ($foo: do
					j += 1
					i := 0
					do
						break[foo](:bar) if i + j > 8
						i += 1
					end while i < 5
				end while j < 5)
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("bar")

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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var result: Std::Symbol | Std::Int | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var l2 value.Value // var i: Std::Int
	_ = l2
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l2 = (value.SmallInt(0)).ToValue()
		t2 = value.Nil
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l2, l0), (value.SmallInt(8)).ToValue())) {
				t1 = (sym2).ToValue()
				break loop0
			}
			l2 = value.AddInts(l2, (value.SmallInt(1)).ToValue())
			t2 = l2
			if !(value.Bool(value.LessThanInts(l2, (value.SmallInt(5)).ToValue()))) {
				break
			}
		}
		t1 = t2
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
	}
	l1 = t1
}
`,
		},
		"static infinite": {
			input: `
				do
					println("foo")
				end while true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for {
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("foo")).ToValue()
		callFrame.SetNativeLineNumber(3)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(42, 4, 15), P(45, 4, 18)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"static one iteration": {
			input: `
				do
					println("foo")
				end while false
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for {
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("foo")).ToValue()
		callFrame.SetNativeLineNumber(3)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		break
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(42, 4, 15), P(46, 4, 19)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoWhile(t *testing.T) {
	tests := goTestTable{
		"with a body": {
			input: `
			  i := 0
				while i < 5
					i += 1
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
	}
}
`,
		},
		"with a body and result value": {
			input: `
			  i := 0
				result := while i < 5
					i += 1
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var l1 value.Value // var result: Std::Int?
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
	for {
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		t1 = l0
	}
	l1 = t1
}
`,
		},
		"with break": {
			input: `
			  i := 0
				while true
					i += 1
					break if i < 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(23, 3, 11), P(26, 3, 14)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"with break and result value": {
			input: `
			  i := 0
				result := while true
					i += 1
					break if i < 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var l1 value.Value // var result: nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			t1 = value.Nil
			break loop0
		}
	}
	l1 = t1
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(33, 3, 21), P(36, 3, 24)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"with labeled break": {
			input: `
			  i := 0
				$foo: while true
					i += 1
					break[foo] if i < 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(29, 3, 17), P(32, 3, 20)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"with complex labeled break": {
			input: `
			  i := 0
				$'foo bar': while true
					i += 1
					break[$'foo bar'] if i < 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(35, 3, 23), P(38, 3, 26)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"with break with value": {
			input: `
			  i := 0
				r := while true
					i += 1
					break true if i < 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var l1 value.Bool // var r: bool
	_ = l1
	var t1 value.Bool
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			t1 = value.True
			break loop0
		}
	}
	l1 = t1
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(28, 3, 16), P(31, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
			},
		},

		"continue in a nested loop": {
			input: `
			 	j := 0
				while j < 5
					j += 1
					i := 0
					while i < 5
						continue if i + j > 8
						i += 1
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

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
	loop1:
		for {
			if !(value.Bool(value.LessThanInts(l1, (value.SmallInt(5)).ToValue()))) {
				break
			}
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				continue loop1
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: while j < 5
					j += 1
					i := 0
					while i < 5
						continue[foo] if i + j > 8
						i += 1
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

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
		for {
			if !(value.Bool(value.LessThanInts(l1, (value.SmallInt(5)).ToValue()))) {
				break
			}
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				continue loop0
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},
		"labeled continue in a nested loop with value": {
			input: `
			 	j := 0
				result := $foo: while j < 5
					j += 1
					i := 0
					while i < 5
						continue[foo](:bump) if i + j > 8
						i += 1
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

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("bump")

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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var result: Std::Symbol | Std::Int | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var l2 value.Value // var i: Std::Int
	_ = l2
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l2 = (value.SmallInt(0)).ToValue()
		t2 = value.Nil
		for {
			if !(value.Bool(value.LessThanInts(l2, (value.SmallInt(5)).ToValue()))) {
				break
			}
			if value.Bool(value.GreaterThanInts(value.AddInts(l2, l0), (value.SmallInt(8)).ToValue())) {
				t1 = (sym2).ToValue()
				continue loop0
			}
			l2 = value.AddInts(l2, (value.SmallInt(1)).ToValue())
			t2 = l2
		}
		t1 = t2
	}
	l1 = t1
}
`,
		},
		"break in a nested loop": {
			input: `
			 	j := 0
				while j < 5
					j += 1
					i := 0
					while i < 5
						break if i + j > 8
						i += 1
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

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
	loop1:
		for {
			if !(value.Bool(value.LessThanInts(l1, (value.SmallInt(5)).ToValue()))) {
				break
			}
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				break loop1
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},

		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: while j < 5
					j += 1
					i := 0
					while i < 5
						break[foo] if i + j > 8
						i += 1
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

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
		for {
			if !(value.Bool(value.LessThanInts(l1, (value.SmallInt(5)).ToValue()))) {
				break
			}
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				break loop0
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},
		"labeled break in a nested loop with a value": {
			input: `
			 	j := 0
				result := $foo: while j < 5
					j += 1
					i := 0
					while i < 5
						break[foo](:bump) if i + j > 8
						i += 1
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

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("bump")

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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var result: Std::Symbol | Std::Int | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var l2 value.Value // var i: Std::Int
	_ = l2
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l2 = (value.SmallInt(0)).ToValue()
		t2 = value.Nil
		for {
			if !(value.Bool(value.LessThanInts(l2, (value.SmallInt(5)).ToValue()))) {
				break
			}
			if value.Bool(value.GreaterThanInts(value.AddInts(l2, l0), (value.SmallInt(8)).ToValue())) {
				t1 = (sym2).ToValue()
				break loop0
			}
			l2 = value.AddInts(l2, (value.SmallInt(1)).ToValue())
			t2 = l2
		}
		t1 = t2
	}
	l1 = t1
}
`,
		},
		"without a body": {
			input: `
				i := 0
				while i < 5; end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		if !(value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue()))) {
			break
		}
	}
}
`,
		},
		"static infinite": {
			input: `
				while true
					println("foo")
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for {
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("foo")).ToValue()
		callFrame.SetNativeLineNumber(3)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(11, 2, 11), P(14, 2, 14)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
		"static impossible": {
			input: `
				while false
					println("foo")
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(11, 2, 11), P(15, 2, 15)), "this loop will never execute since type `false` is falsy"),
				diagnostic.NewWarning(L(P(22, 3, 6), P(36, 3, 20)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoModifierUntil(t *testing.T) {
	tests := goTestTable{
		"single line": {
			input: `
			  i := 0
				i += 1 until i >= 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
	}
}
`,
		},
		"single line with result value": {
			input: `
			  i := 0
				result := (i += 1 until i >= 5)
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var l1 value.Value // var result: Std::Int?
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		t1 = l0
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
	}
	l1 = t1
}
`,
		},
		"multiline": {
			input: `
			  i := 0
				do
					i += 1
				end until i >= 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
	}
}
`,
		},
		"with break": {
			input: `
			  i := 0
				do
					i += 1
					break if i < 5
				end until false
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(66, 6, 15), P(70, 6, 19)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"with break with result value": {
			input: `
			  i := 0
				result := (do
					i += 1
					break if i < 5
				end until false)
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var l1 value.Value // var result: nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			t1 = value.Nil
			break loop0
		}
	}
	l1 = t1
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(77, 6, 15), P(81, 6, 19)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"with labeled break": {
			input: `
			  i := 0
				$foo: do
					i += 1
					break[foo] if i < 5
				end until false
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(77, 6, 15), P(81, 6, 19)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"with break with value": {
			input: `
			  i := 0
				result := (do
					i += 1
					break true if i < 5
				end until false)
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var l1 value.Bool // var result: bool
	_ = l1
	var t1 value.Bool
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			t1 = value.True
			break loop0
		}
	}
	l1 = t1
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(82, 6, 15), P(86, 6, 19)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"continue in a nested loop": {
			input: `
			 	j := 0
				do
					j += 1
					i := 0
					do
						continue if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
	loop1:
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				continue loop1
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
			if value.Bool(value.GreaterThanEqualInts(l1, (value.SmallInt(5)).ToValue())) {
				break
			}
		}
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
	}
}
`,
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						continue[foo] if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				continue loop0
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
			if value.Bool(value.GreaterThanEqualInts(l1, (value.SmallInt(5)).ToValue())) {
				break
			}
		}
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
	}
}
`,
		},
		"labeled continue with value in a nested loop": {
			input: `
			 	j := 0
				result := ($foo: do
					j += 1
					i := 0
					do
						continue[foo](:bar) if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5)
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("bar")

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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var result: Std::Symbol | Std::Int | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var l2 value.Value // var i: Std::Int
	_ = l2
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l2 = (value.SmallInt(0)).ToValue()
		t2 = value.Nil
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l2, l0), (value.SmallInt(8)).ToValue())) {
				t1 = (sym2).ToValue()
				continue loop0
			}
			l2 = value.AddInts(l2, (value.SmallInt(1)).ToValue())
			t2 = l2
			if value.Bool(value.GreaterThanEqualInts(l2, (value.SmallInt(5)).ToValue())) {
				break
			}
		}
		t1 = t2
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
	}
	l1 = t1
}
`,
		},
		"break in a nested loop": {
			input: `
			 	j := 0
				do
					j += 1
					i := 0
					do
						break if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
	loop1:
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				break loop1
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
			if value.Bool(value.GreaterThanEqualInts(l1, (value.SmallInt(5)).ToValue())) {
				break
			}
		}
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
	}
}
`,
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: do
					j += 1
					i := 0
					do
						break[foo] if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				break loop0
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
			if value.Bool(value.GreaterThanEqualInts(l1, (value.SmallInt(5)).ToValue())) {
				break
			}
		}
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
	}
}
`,
		},
		"labeled break with value in a nested loop": {
			input: `
			 	j := 0
				result := ($foo: do
					j += 1
					i := 0
					do
						break[foo](:bar) if i + j > 8
						i += 1
					end until i >= 5
				end until j >= 5)
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("bar")

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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var result: Std::Symbol | Std::Int | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var l2 value.Value // var i: Std::Int
	_ = l2
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l2 = (value.SmallInt(0)).ToValue()
		t2 = value.Nil
		for {
			if value.Bool(value.GreaterThanInts(value.AddInts(l2, l0), (value.SmallInt(8)).ToValue())) {
				t1 = (sym2).ToValue()
				break loop0
			}
			l2 = value.AddInts(l2, (value.SmallInt(1)).ToValue())
			t2 = l2
			if value.Bool(value.GreaterThanEqualInts(l2, (value.SmallInt(5)).ToValue())) {
				break
			}
		}
		t1 = t2
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
	}
	l1 = t1
}
`,
		},
		"static infinite": {
			input: `
				do
					println("foo")
				end until false
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for {
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("foo")).ToValue()
		callFrame.SetNativeLineNumber(3)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(42, 4, 15), P(46, 4, 19)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"static one iteration": {
			input: `
				do
					println("foo")
				end until true
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for {
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("foo")).ToValue()
		callFrame.SetNativeLineNumber(3)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		break
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(42, 4, 15), P(45, 4, 18)), "this condition will always have the same result since type `true` is truthy"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoUntil(t *testing.T) {
	tests := goTestTable{
		"with a body": {
			input: `
			  i := 0
				until i >= 5
					i += 1
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
	}
}
`,
		},
		"with break": {
			input: `
			  i := 0
				until false
					i += 1
					break if i < 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(23, 3, 11), P(27, 3, 15)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"with break and result value": {
			input: `
			  i := 0
				result := until false
					i += 1
					break if i < 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var l1 value.Value // var result: nil
	_ = l1
	var t1 value.Value
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			t1 = value.Nil
			break loop0
		}
	}
	l1 = t1
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(33, 3, 21), P(37, 3, 25)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"with labeled break": {
			input: `
			  i := 0
				$foo: until false
					i += 1
					break[foo] if i < 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(29, 3, 17), P(33, 3, 21)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"with break with value": {
			input: `
			  i := 0
				until false
					i += 1
					break true if i < 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			break loop0
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(23, 3, 11), P(27, 3, 15)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"with break with value and result value": {
			input: `
			  i := 0
				result := until false
					i += 1
					break true if i < 5
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var l1 value.Bool // var result: bool
	_ = l1
	var t1 value.Bool
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		if value.Bool(value.LessThanInts(l0, (value.SmallInt(5)).ToValue())) {
			t1 = value.True
			break loop0
		}
	}
	l1 = t1
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(33, 3, 21), P(37, 3, 25)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"continue in a nested loop": {
			input: `
			 	j := 0
				until j >= 5
					j += 1
					i := 0
					until i >= 5
						continue if i + j > 8
						i += 1
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

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
	loop1:
		for {
			if value.Bool(value.GreaterThanEqualInts(l1, (value.SmallInt(5)).ToValue())) {
				break
			}
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				continue loop1
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},
		"labeled continue in a nested loop": {
			input: `
			 	j := 0
				$foo: until j >= 5
					j += 1
					i := 0
					until i >= 5
						continue[foo] if i + j > 8
						i += 1
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

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
		for {
			if value.Bool(value.GreaterThanEqualInts(l1, (value.SmallInt(5)).ToValue())) {
				break
			}
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				continue loop0
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},
		"break in a nested loop": {
			input: `
			 	j := 0
				until j >= 5
					j += 1
					i := 0
					until i >= 5
						break if i + j > 8
						i += 1
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

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
	loop1:
		for {
			if value.Bool(value.GreaterThanEqualInts(l1, (value.SmallInt(5)).ToValue())) {
				break
			}
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				break loop1
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},
		"labeled break in a nested loop": {
			input: `
			 	j := 0
				$foo: until j >= 5
					j += 1
					i := 0
					until i >= 5
						break[foo] if i + j > 8
						i += 1
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

var _ = symbol.Value
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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var i: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
loop0:
	for {
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l1 = (value.SmallInt(0)).ToValue()
		for {
			if value.Bool(value.GreaterThanEqualInts(l1, (value.SmallInt(5)).ToValue())) {
				break
			}
			if value.Bool(value.GreaterThanInts(value.AddInts(l1, l0), (value.SmallInt(8)).ToValue())) {
				break loop0
			}
			l1 = value.AddInts(l1, (value.SmallInt(1)).ToValue())
		}
	}
}
`,
		},
		"labeled break in a nested loop with value": {
			input: `
			 	j := 0
				result := $foo: until j >= 5
					j += 1
					i := 0
					until i >= 5
						break[foo](:bump) if i + j > 8
						i += 1
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

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("bump")

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
	var l0 value.Value // var j: Std::Int
	_ = l0
	var l1 value.Value // var result: Std::Symbol | Std::Int | nil
	_ = l1
	var t1 value.Value
	_ = t1
	var l2 value.Value // var i: Std::Int
	_ = l2
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	t1 = value.Nil
loop0:
	for {
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
		l0 = value.AddInts(l0, (value.SmallInt(1)).ToValue())
		l2 = (value.SmallInt(0)).ToValue()
		t2 = value.Nil
		for {
			if value.Bool(value.GreaterThanEqualInts(l2, (value.SmallInt(5)).ToValue())) {
				break
			}
			if value.Bool(value.GreaterThanInts(value.AddInts(l2, l0), (value.SmallInt(8)).ToValue())) {
				t1 = (sym2).ToValue()
				break loop0
			}
			l2 = value.AddInts(l2, (value.SmallInt(1)).ToValue())
			t2 = l2
		}
		t1 = t2
	}
	l1 = t1
}
`,
		},
		"without a body": {
			input: `
				i := 0
				until i >= 5; end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var l0 value.Value // var i: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(0)).ToValue()
	for {
		if value.Bool(value.GreaterThanEqualInts(l0, (value.SmallInt(5)).ToValue())) {
			break
		}
	}
}
`,
		},
		"static infinite": {
			input: `
				until false
					println("foo")
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	for {
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("foo")).ToValue()
		callFrame.SetNativeLineNumber(3)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
	}
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(11, 2, 11), P(15, 2, 15)), "this condition will always have the same result since type `false` is falsy"),
			},
		},
		"static impossible": {
			input: `
				until true
					println("foo")
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(11, 2, 11), P(14, 2, 14)), "this loop will never execute since type `true` is truthy"),
				diagnostic.NewWarning(L(P(21, 3, 6), P(35, 3, 20)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoCatch(t *testing.T) {
	tests := goTestTable{
		"simple catch": {
			input: `
				result := do
					throw :foo
				catch String() as str
					str
				catch :foo
					3
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("foo")
var cc_main_1 = &vm.CallCache{}

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
	var l0 value.Value // var result: Std::String | Std::Int
	_ = l0
	var t1 value.Value
	_ = t1
	var t2 value.Value
	_ = t2
	var l1 value.Value // var str: any
	_ = l1
	var t3 []value.Value
	_ = t3
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	thread.CaptureStackTrace()
	err = (sym2).ToValue()
	goto lbl1
	t1 = value.Nil
	goto lbl3
lbl1:
	t2 = err
	l1 = t2
	if value.Bool(value.IsA(t2, value.StringClass)) {
		t1 = l1
		goto lbl3
	}
	t3 = value.ResizeNativeArgs(t3, 3)
	t3[0] = t2
	t3[1] = (sym2).ToValue()
	callFrame.SetNativeLineNumber(6)
	t2, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t3...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t2) {
		t1 = (value.SmallInt(3)).ToValue()
		goto lbl3
	}
	thread.Panic(t2)
lbl3:
	l0 = t1
}
`,
		},
		"catch with stack trace": {
			input: `
				do
					throw :foo
				catch String() as str, stack_trace
					puts str
				catch :foo
					puts "symbol!"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("foo")
var sym3 = value.ToSymbol("puts@1")
var fn_method0 vm.NativeFunction // Std::Kernel::puts@1
var cc_main_1 = &vm.CallCache{}

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
	var t1 value.Value
	_ = t1
	var l0 value.Value // var stack_trace: Std::StackTrace
	_ = l0
	var l1 value.Value // var str: any
	_ = l1
	var t2 []value.Value
	_ = t2
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	thread.CaptureStackTrace()
	err = (sym2).ToValue()
	goto lbl1
	goto lbl3
lbl1:
	t1 = err
	l0 = thread.ErrStackTrace().ToValue()
	l1 = t1
	if value.Bool(value.IsA(t1, value.StringClass)) {
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = l1
		callFrame.SetNativeLineNumber(5)
		_, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: puts@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		goto lbl3
	}
	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = t1
	t2[1] = (sym2).ToValue()
	callFrame.SetNativeLineNumber(6)
	t1, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t2...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t1) {
		t2 = value.ResizeNativeArgs(t2, 3)
		t2[0] = (value.KernelModule).ToValue()
		t2[1] = (value.String("symbol!")).ToValue()
		callFrame.SetNativeLineNumber(7)
		_, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: puts@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		goto lbl3
	}
	thread.Panic(t1)
lbl3:
}
`,
		},
		"finally": {
			input: `
				do
					println("foo")
				finally
					println("bar")
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = (value.KernelModule).ToValue()
	t1[1] = (value.String("foo")).ToValue()
	callFrame.SetNativeLineNumber(3)
	_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		goto lbl1
	}
	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = (value.KernelModule).ToValue()
	t1[1] = (value.String("bar")).ToValue()
	callFrame.SetNativeLineNumber(5)
	_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	goto lbl3
lbl1:
	t2 = err
	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = (value.KernelModule).ToValue()
	t1[1] = (value.String("bar")).ToValue()
	_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	thread.Panic(t2)
lbl3:
}
`,
		},
		"finally between continue and loop": {
			input: `
				a := true
				while a
					do
						do
							println("foo")
							continue
						finally
							println("finally 1")
						end
					finally
						println("finally 2")
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

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Bool // var a: bool
	_ = l0
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.True
loop0:
	for {
		if !(l0) {
			break
		}
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("foo")).ToValue()
		callFrame.SetNativeLineNumber(6)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			goto lbl2
		}
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 1")).ToValue()
		callFrame.SetNativeLineNumber(9)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			goto lbl1
		}
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 2")).ToValue()
		callFrame.SetNativeLineNumber(12)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		continue loop0
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 1")).ToValue()
		callFrame.SetNativeLineNumber(9)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			goto lbl1
		}
		goto lbl4
	lbl2:
		t2 = err
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 1")).ToValue()
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			goto lbl1
		}
		err = t2
		goto lbl1
	lbl4:
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 2")).ToValue()
		callFrame.SetNativeLineNumber(12)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		goto lbl6
	lbl1:
		t2 = err
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 2")).ToValue()
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		thread.Panic(t2)
	lbl6:
	}
}
`,
		},
		"finally between break and loop": {
			input: `
				a := true
				while a
					do
						do
							println("foo")
							break
						finally
							println("finally 1")
						end
					finally
						println("finally 2")
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

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
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

	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Bool // var a: bool
	_ = l0
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var t2 value.Value
	_ = t2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.True
loop0:
	for {
		if !(l0) {
			break
		}
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("foo")).ToValue()
		callFrame.SetNativeLineNumber(6)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			goto lbl2
		}
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 1")).ToValue()
		callFrame.SetNativeLineNumber(9)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			goto lbl1
		}
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 2")).ToValue()
		callFrame.SetNativeLineNumber(12)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		break loop0
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 1")).ToValue()
		callFrame.SetNativeLineNumber(9)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			goto lbl1
		}
		goto lbl4
	lbl2:
		t2 = err
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 1")).ToValue()
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			goto lbl1
		}
		err = t2
		goto lbl1
	lbl4:
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 2")).ToValue()
		callFrame.SetNativeLineNumber(12)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		goto lbl6
	lbl1:
		t2 = err
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("finally 2")).ToValue()
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		thread.Panic(t2)
	lbl6:
	}
}
`,
		},
		"catch and finally": {
			input: `
				def foo! :foo
					println "foo"
					throw :foo
				end

				do
					foo()
				catch :foo
					println "bar"
				finally
					println "baz"
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var cc_main_1 = &vm.CallCache{}

var sym0 = value.ToSymbol("Std::Kernel::foo")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
var fn_method1 vm.NativeFunction // Std::Kernel::println@1
var sym3 = value.ToSymbol("foo")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Std::Kernel::foo, loc: <main>:2:5
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 []value.Value
	_ = t1

	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = self
	t1[1] = (value.String("foo")).ToValue()
	callFrame.SetNativeLineNumber(3)
	_, err = fn_method1(thread, t1) // receiver: Std::Kernel, name: println@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	thread.CaptureStackTrace()
	return result, (sym3).ToValue()
	return value.Nil, value.Undefined

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
	var t1 []value.Value
	_ = t1
	var t2 value.Value
	_ = t2
	var t3 value.Value
	_ = t3
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	methodDefinitions()
	fn_method1 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym4, sym1, 1)
	defer thread.PopNativeCallFrame()
	_, err = fn_method0(thread, (value.KernelModule).ToValue()) // receiver: Std::Kernel, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		goto lbl1
	}
	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = (value.KernelModule).ToValue()
	t1[1] = (value.String("baz")).ToValue()
	callFrame.SetNativeLineNumber(12)
	_, err = fn_method1(thread, t1) // receiver: Std::Kernel, name: println@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	goto lbl3
lbl1:
	t2 = err
	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = t2
	t1[1] = (sym3).ToValue()
	callFrame.SetNativeLineNumber(9)
	t3, err = thread.CallMethodByNameWithCache(symbol.OpEqual, &cc_main_1, t1...) // receiver: any, name: ==
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	if value.ToBool(t3) {
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("bar")).ToValue()
		callFrame.SetNativeLineNumber(10)
		_, err = fn_method1(thread, t1) // receiver: Std::Kernel, name: println@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
		goto lbl2
	}
	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = (value.KernelModule).ToValue()
	t1[1] = (value.String("baz")).ToValue()
	callFrame.SetNativeLineNumber(12)
	_, err = fn_method1(thread, t1) // receiver: Std::Kernel, name: println@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	thread.Panic(t2)
lbl2:
	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = (value.KernelModule).ToValue()
	t1[1] = (value.String("baz")).ToValue()
	_, err = fn_method1(thread, t1) // receiver: Std::Kernel, name: println@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
lbl3:
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (value.KernelModule).SingletonClass() // Std::Kernel
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
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

func TestGoDefer(t *testing.T) {
	tests := goTestTable{
		"multiple defer": {
			input: `
				puts "1. open file"
				defer puts "2. close file"

				puts "3. open TCP socket"
				defer puts "4. close TCP socket"
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("puts@1")
var fn_method0 vm.NativeFunction // Std::Kernel::puts@1

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
	var t1 []value.Value
	_ = t1
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = (value.KernelModule).ToValue()
	t1[1] = (value.String("1. open file")).ToValue()
	callFrame.SetNativeLineNumber(2)
	_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: puts@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	defer func() {
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("2. close file")).ToValue()
		callFrame.SetNativeLineNumber(3)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: puts@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
	}()
	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = (value.KernelModule).ToValue()
	t1[1] = (value.String("3. open TCP socket")).ToValue()
	callFrame.SetNativeLineNumber(5)
	_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: puts@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	defer func() {
		t1 = value.ResizeNativeArgs(t1, 3)
		t1[0] = (value.KernelModule).ToValue()
		t1[1] = (value.String("4. close TCP socket")).ToValue()
		callFrame.SetNativeLineNumber(6)
		_, err = fn_method0(thread, t1) // receiver: Std::Kernel, name: puts@1
		if err.IsNotUndefined() {
			thread.CaptureStackTrace()
			thread.Panic(err)
		}
	}()
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

// func TestBytecodeThrow(t *testing.T) {
// 	tests := bytecodeTestTable{
// 		"with a value": {
// 			input: `throw unchecked :foo`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.THROW),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(19, 1, 20)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 				},
// 				[]value.Value{
// 					value.ToSymbol("foo").ToValue(),
// 				},
// 			),
// 		},
// 		"without a value": {
// 			input: `throw`,
// 			err: diagnostic.DiagnosticList{
// 				diagnostic.NewFailure(L(P(0, 1, 1), P(4, 1, 5)), "thrown value of type `Std::Error` must be caught"),
// 			},
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			bytecodeCompilerTest(tc, t)
// 		})
// 	}
// }
