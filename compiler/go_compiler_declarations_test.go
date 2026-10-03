package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

func TestGoSingletonBlock(t *testing.T) {
	tests := goTestTable{
		"define in top-level": {
			input: `
				singleton
					def foo then :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 2, 5), P(44, 4, 7)), "singleton definitions cannot appear in this context"),
				diagnostic.NewFailure(L(P(20, 3, 6), P(36, 3, 22)), "method definitions cannot appear in this context"),
			},
		},
		"define in a method": {
			input: `
				def baz
					singleton
						def foo then :bar
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(18, 3, 6), P(59, 5, 8)), "singleton definitions cannot appear in this context"),
				diagnostic.NewFailure(L(P(34, 4, 7), P(50, 4, 23)), "method definitions cannot appear in this context"),
			},
		},
		"define in a setter": {
			input: `
				def baz=(arg)
					singleton
						def foo then :bar
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(14, 2, 14), P(16, 2, 16)), "cannot declare parameter `arg` without a type"),
				diagnostic.NewFailure(L(P(24, 3, 6), P(65, 5, 8)), "singleton definitions cannot appear in this context"),
				diagnostic.NewFailure(L(P(40, 4, 7), P(56, 4, 23)), "method definitions cannot appear in this context"),
			},
		},
		"define in a class": {
			input: `
				class Baz
					singleton
						def foo then :bar
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

var sym4 = value.ToSymbol("main")

var const0 *value.Class // Baz
var sym0 = value.ToSymbol("Baz")

var sym1 = value.ToSymbol("Baz::foo")
var sym2 = value.ToSymbol("<main>")
var sym3 = value.ToSymbol("bar")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Baz::foo, loc: <main>:4:7
	var callFrame *vm.CallFrame
	_ = callFrame

	return (sym3).ToValue(), value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0                 // Baz
	class = class.SingletonClass() // &Baz
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},
		"define in a mixin": {
			input: `
				mixin Baz
					singleton
						def foo then :bar
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

var sym4 = value.ToSymbol("main")

var const0 *value.Mixin // Baz
var sym0 = value.ToSymbol("Baz")

var sym1 = value.ToSymbol("Baz::foo")
var sym2 = value.ToSymbol("<main>")
var sym3 = value.ToSymbol("bar")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Baz::foo, loc: <main>:4:7
	var callFrame *vm.CallFrame
	_ = callFrame

	return (sym3).ToValue(), value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
	const0 = value.NewMixin()
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0                 // Baz
	class = class.SingletonClass() // &Baz
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},
		"define in an interface": {
			input: `
				interface Baz
					singleton
						def foo then :bar
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

var sym4 = value.ToSymbol("main")

var const0 *value.Interface // Baz
var sym0 = value.ToSymbol("Baz")

var sym1 = value.ToSymbol("Baz::foo")
var sym2 = value.ToSymbol("<main>")
var sym3 = value.ToSymbol("bar")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Baz::foo, loc: <main>:4:7
	var callFrame *vm.CallFrame
	_ = callFrame

	return (sym3).ToValue(), value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
	const0 = value.NewInterface()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass() // Baz
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},
		"define in a module": {
			input: `
				module Baz
					singleton
						def foo then :bar
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(21, 3, 6), P(62, 5, 8)), "singleton definitions cannot appear in this context"),
				diagnostic.NewFailure(L(P(37, 4, 7), P(53, 4, 23)), "method definitions cannot appear in this context"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			goCompilerTest(tc, t)
		})
	}
}

func TestGoGetter(t *testing.T) {
	tests := goTestTable{
		"define single getter": {
			input: `
				class Foo
					getter foo
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

var sym2 = value.ToSymbol("main")
var sym3 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("foo")

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

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
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
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.DefineGetter(class, sym1, 0)
}
`,
		},
		"define three getters": {
			input: `
				class Foo
					getter foo: Foo?, bar: Int?, baz: String?
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

var sym4 = value.ToSymbol("main")
var sym5 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("bar")
var sym2 = value.ToSymbol("baz")
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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym5, 1)
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
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("bar"): 0, symbol.ToSymbol("baz"): 1, symbol.ToSymbol("foo"): 2}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.DefineGetter(class, sym1, 0)
	vm.DefineGetter(class, sym2, 1)
	vm.DefineGetter(class, sym3, 2)
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

func TestGoSetter(t *testing.T) {
	tests := goTestTable{
		"define single setter": {
			input: `
				class Foo
					setter foo: String?
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

var sym2 = value.ToSymbol("main")
var sym3 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("foo")

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

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
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
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.DefineSetter(class, sym1, 0)
}
`,
		},
		"define three setters": {
			input: `
				class Foo
					setter foo: Foo?, baz: String?
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

var sym3 = value.ToSymbol("main")
var sym4 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("baz")
var sym2 = value.ToSymbol("foo")

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

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym4, 1)
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
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("baz"): 0, symbol.ToSymbol("foo"): 1}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.DefineSetter(class, sym1, 0)
	vm.DefineSetter(class, sym2, 1)
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

func TestGoAttr(t *testing.T) {
	tests := goTestTable{
		"define single attr": {
			input: `
				class Foo
					attr foo: String?
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

var sym2 = value.ToSymbol("main")
var sym3 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("foo")

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

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
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
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.DefineGetter(class, sym1, 0)
	vm.DefineSetter(class, sym1, 0)
}
`,
		},
		"define three attrs": {
			input: `
				class Foo
					attr foo: Foo?, baz: String?
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

var sym3 = value.ToSymbol("main")
var sym4 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("baz")
var sym2 = value.ToSymbol("foo")

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

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym4, 1)
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
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("baz"): 0, symbol.ToSymbol("foo"): 1}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.DefineGetter(class, sym1, 0)
	vm.DefineSetter(class, sym1, 0)
	vm.DefineGetter(class, sym2, 1)
	vm.DefineSetter(class, sym2, 1)
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

func TestGoAlias(t *testing.T) {
	tests := goTestTable{
		"define single alias": {
			input: `
				class Foo
					def bar; end
					alias foo bar
				end
				f := Foo()
				f.foo
				f.bar
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

var sym3 = value.ToSymbol("main")
var sym5 = value.ToSymbol("foo")
var sym6 = value.ToSymbol("bar")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:bar")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:bar, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

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
	var l0 value.Value // var f: Foo
	_ = l0
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = const0.CreateInstance()
	_, err = fn_method0(thread, l0) // receiver: Foo, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method0(thread, l0) // receiver: Foo, name: bar
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

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(class, "bar", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},
		"define three aliases": {
			input: `
				class Foo
					def bar; end
					def delete; end
					def plus; end
					alias foo bar, remove delete, add plus
				end
				f := Foo()
				f.bar
				f.foo
				f.delete
				f.remove
				f.plus
				f.add
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

var sym5 = value.ToSymbol("main")
var sym7 = value.ToSymbol("bar")
var sym8 = value.ToSymbol("foo")
var sym9 = value.ToSymbol("delete")
var sym10 = value.ToSymbol("remove")
var sym11 = value.ToSymbol("plus")
var sym12 = value.ToSymbol("add")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym4 = value.ToSymbol("Foo.:plus")

func fn_method2(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:plus, loc: <main>:5:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Nil, value.Undefined

}

var sym1 = value.ToSymbol("Foo.:bar")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:bar, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	return value.Nil, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:delete")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:delete, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame

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
	var l0 value.Value // var f: Foo
	_ = l0
	var err value.Value
	_ = err
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym5, sym2, 1)
	defer thread.PopNativeCallFrame()
	l0 = const0.CreateInstance()
	_, err = fn_method0(thread, l0) // receiver: Foo, name: bar
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method0(thread, l0) // receiver: Foo, name: foo
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0) // receiver: Foo, name: delete
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method1(thread, l0) // receiver: Foo, name: remove
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method2(thread, l0) // receiver: Foo, name: plus
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		thread.Panic(err)
	}
	_, err = fn_method2(thread, l0) // receiver: Foo, name: add
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

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(class, "add", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method2(thread, args[0])
		return result, err
	})
	vm.Def(class, "bar", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
	vm.Def(class, "delete", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
	vm.Def(class, "plus", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method2(thread, args[0])
		return result, err
	})
	vm.Def(class, "remove", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
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

func TestGoDefClass(t *testing.T) {
	tests := goTestTable{
		"class with a relative name without a body": {
			input: "class Foo; end",
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

var const0 *value.Class // Foo
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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
`,
		},
		"named class inside of a method": {
			input: `
				def foo
				  class ::Bar; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(19, 3, 7), P(34, 3, 22)), "class definitions cannot appear in this context"),
			},
		},
		"class with an absolute parent": {
			input: `
				class Bar; end
				class Foo < ::Bar; end
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

var sym2 = value.ToSymbol("main")
var sym3 = value.ToSymbol("<main>")

var const0 *value.Class // Bar
var sym0 = value.ToSymbol("Bar")
var const1 *value.Class // Foo
var sym1 = value.ToSymbol("Foo")

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

	initGlobalEnv()
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	parentNamespace = (value.RootModule).ToValue()
	const1 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const1)
	value.AddConstant(parentNamespace, sym1, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
	class = const1
	superclass = const0
	class.SetSuperclass(superclass)
}
`,
		},
		"class with an absolute nested parent": {
			input: `
				module Baz
					class Bar; end
				end
				class Foo < Baz::Bar; end
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

var sym3 = value.ToSymbol("main")
var sym4 = value.ToSymbol("<main>")

var const0 *value.Module // Baz
var sym0 = value.ToSymbol("Baz")
var const1 *value.Class // Baz::Bar
var sym1 = value.ToSymbol("Bar")
var const2 *value.Class // Foo
var sym2 = value.ToSymbol("Foo")

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

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym3, sym4, 1)
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

	parentNamespace = (const0).ToValue()
	const1 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const1)
	value.AddConstant(parentNamespace, sym1, namespace)

	parentNamespace = (value.RootModule).ToValue()
	const2 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const2)
	value.AddConstant(parentNamespace, sym2, namespace)

	class = const1
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
	class = const2
	superclass = const1
	class.SetSuperclass(superclass)
}
`,
		},
		"class with an absolute name without a body": {
			input: "class ::Foo; end",
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

var const0 *value.Class // Foo
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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
`,
		},
		"class with an absolute nested name without a body": {
			input: "class ::Std::Int::Foo; end",
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

var const0 *value.Class // Std::Int::Foo
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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
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

	parentNamespace = (value.IntClass).ToValue()
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}
`,
		},
		"class with a body": {
			input: `
				class Foo
					a := 1
					a += 2
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

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym3 = value.ToSymbol("<class: Foo>")

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

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0

	self = (const0).ToValue()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 2)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.AddInts(l0, (value.SmallInt(2)).ToValue())
}
`,
		},
		"nested classes": {
			input: `
				class Foo
					class Bar
						a := 1
						a += 2
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

var sym2 = value.ToSymbol("main")
var sym3 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")
var const1 *value.Class // Foo::Bar
var sym1 = value.ToSymbol("Bar")

var sym4 = value.ToSymbol("<class: Foo>")

var sym5 = value.ToSymbol("<class: Foo::Bar>")

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

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym2, sym3, 1)
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	parentNamespace = (const0).ToValue()
	const1 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const1)
	value.AddConstant(parentNamespace, sym1, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
	class = const1
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame

	self = (const0).ToValue()
	callFrame = thread.AddNativeCallFrame(sym4, sym3, 2)
	defer thread.PopNativeCallFrame()
	fn_ns_expr1(thread)
}

func fn_ns_expr1(thread *vm.Thread) { // namespace: Foo::Bar, loc: <main>:3:6
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0

	self = (const1).ToValue()
	callFrame = thread.AddNativeCallFrame(sym5, sym3, 3)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.AddInts(l0, (value.SmallInt(2)).ToValue())
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

func TestGoDefModule(t *testing.T) {
	tests := goTestTable{
		"module with a relative name without a body": {
			input: "module Foo; end",
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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
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
`,
		},
		"named module inside of a method": {
			input: `
				def foo
					module Bar; end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(18, 3, 6), P(32, 3, 20)), "module definitions cannot appear in this context"),
			},
		},
		"module with an absolute name without a body": {
			input: "module ::Foo; end",
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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
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
`,
		},
		"module with an absolute nested name without a body": {
			input: "module ::Std::Int::Foo; end",
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

var const0 *value.Module // Std::Int::Foo
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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
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

	parentNamespace = (value.IntClass).ToValue()
	const0 = value.NewModule()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}
`,
		},
		"module with a body": {
			input: `
				module Foo
					a := 1
					a += 2
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

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym3 = value.ToSymbol("<module: Foo>")

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

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym1, sym2, 1)
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

}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0

	self = value.Ref((const0).SingletonClass())
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 2)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.AddInts(l0, (value.SmallInt(2)).ToValue())
}
`,
		},
		"nested modules": {
			input: `
				module Foo
					module Bar
						a := 1
						a += 2
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

var sym2 = value.ToSymbol("main")
var sym3 = value.ToSymbol("<main>")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")
var const1 *value.Module // Foo::Bar
var sym1 = value.ToSymbol("Bar")

var sym4 = value.ToSymbol("<module: Foo>")

var sym5 = value.ToSymbol("<module: Foo::Bar>")

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

	initGlobalEnv()
	callFrame = thread.AddNativeCallFrame(sym2, sym3, 1)
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
	callFrame = thread.AddNativeCallFrame(sym4, sym3, 2)
	defer thread.PopNativeCallFrame()
	fn_ns_expr1(thread)
}

func fn_ns_expr1(thread *vm.Thread) { // namespace: Foo::Bar, loc: <main>:3:6
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0

	self = value.Ref((const1).SingletonClass())
	callFrame = thread.AddNativeCallFrame(sym5, sym3, 3)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(1)).ToValue()
	l0 = value.AddInts(l0, (value.SmallInt(2)).ToValue())
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

func TestGoDefMethod(t *testing.T) {
	tests := goTestTable{
		"define method in top level": {
			input: `
				def foo then :bar
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

var sym3 = value.ToSymbol("main")

var sym0 = value.ToSymbol("Std::Kernel::foo")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("bar")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Std::Kernel::foo, loc: <main>:2:5
	var callFrame *vm.CallFrame
	_ = callFrame

	return (sym2).ToValue(), value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym1, 1)
	defer thread.PopNativeCallFrame()
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (value.KernelModule).SingletonClass() // Std::Kernel
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},
		"define a setter": {
			input: `
				def foo=(a: Int)
					println(a + 2)
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

var sym3 = value.ToSymbol("main")

var sym0 = value.ToSymbol("Std::Kernel::foo=")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
var fn_method1 vm.NativeFunction // Std::Kernel::println@1
func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Std::Kernel::foo=, loc: <main>:2:5
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1
	var t2 []value.Value
	_ = t2

	t2 = value.ResizeNativeArgs(t2, 3)
	t2[0] = self
	t2[1] = value.AddInts(l0, (value.SmallInt(2)).ToValue())
	callFrame.SetNativeLineNumber(3)
	t1, err = fn_method1(thread, t2) // receiver: Std::Kernel, name: println@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	return l0, value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	methodDefinitions()
	fn_method1 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym3, sym1, 1)
	defer thread.PopNativeCallFrame()
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (value.KernelModule).SingletonClass() // Std::Kernel
	vm.Def(class, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
		},
		"define a setter with return": {
			input: `
				def foo=(a: Int)
					println(a + 2)
					return "siema"
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

var sym3 = value.ToSymbol("main")

var sym0 = value.ToSymbol("Std::Kernel::foo=")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("println@1")
var fn_method1 vm.NativeFunction // Std::Kernel::println@1
func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value) (result value.Value, err value.Value) { // method: Std::Kernel::foo=, loc: <main>:2:5
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 []value.Value
	_ = t1

	t1 = value.ResizeNativeArgs(t1, 3)
	t1[0] = self
	t1[1] = value.AddInts(l0, (value.SmallInt(2)).ToValue())
	callFrame.SetNativeLineNumber(3)
	_, err = fn_method1(thread, t1) // receiver: Std::Kernel, name: println@1
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	return l0, value.Undefined
	return l0, value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	methodDefinitions()
	fn_method1 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym2))

	callFrame = thread.AddNativeCallFrame(sym3, sym1, 1)
	defer thread.PopNativeCallFrame()
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (value.KernelModule).SingletonClass() // Std::Kernel
	vm.Def(class, "foo=", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1])
		return result, err
	}, vm.DefWithParameters(1))
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(54, 4, 13), P(60, 4, 19)), "values returned in void context will be ignored"),
			},
		},
		"define generator": {
			input: `
				def *foo(a: Int, b: Int = 2): Int ! String
					yield a + b
					throw "lol"
					return 10
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position"
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

var fn_method0 *vm.BytecodeFunction // method: Std::Kernel::foo, loc: <main>:2:5

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

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (value.KernelModule).SingletonClass() // Std::Kernel
	fn_method0 = vm.NewBytecodeFunctionWithOptions(
		vm.BytecodeFunctionWithInstructions([]byte{
			byte(bytecode.GET_LOCAL_2),
			byte(bytecode.JUMP_UNLESS_UNDEF), 0x00, 0x02,
			byte(bytecode.INT_2),
			byte(bytecode.SET_LOCAL_2),
			byte(bytecode.GENERATOR),
			byte(bytecode.RETURN),
			byte(bytecode.GET_LOCAL_1),
			byte(bytecode.GET_LOCAL_2),
			byte(bytecode.ADD_INT),
			byte(bytecode.YIELD),
			byte(bytecode.LOAD_VALUE_0),
			byte(bytecode.THROW),
			byte(bytecode.POP),
			byte(bytecode.LOAD_INT_8), 0x0a,
			byte(bytecode.YIELD),
			byte(bytecode.STOP_ITERATION),
			byte(bytecode.YIELD),
			byte(bytecode.STOP_ITERATION),
			byte(bytecode.LOOP), 0x00, 0x04,
		}),
		vm.BytecodeFunctionWithLocation(position.NewLocation("<main>", position.NewSpan(position.New(5, 2, 5), position.New(103, 6, 7)))),
		vm.BytecodeFunctionWithUpvalueCount(0),
		vm.BytecodeFunctionWithStringName("Std::Kernel::foo"),
		vm.BytecodeFunctionWithParameters(2),
		vm.BytecodeFunctionWithOptionalParameters(1),
		vm.BytecodeFunctionWithCatchEntriesVar(
			vm.NewCatchEntry(-1, -1, 8, false),
		),
		vm.BytecodeFunctionWithLineInfoListVar(
			bytecode.NewLineInfo(2, 7),
			bytecode.NewLineInfo(6, 1),
			bytecode.NewLineInfo(3, 4),
			bytecode.NewLineInfo(4, 3),
			bytecode.NewLineInfo(5, 4),
			bytecode.NewLineInfo(6, 2),
			bytecode.NewLineInfo(2, 3),
		),
		vm.BytecodeFunctionWithValuesVar(
			(value.String("lol")).ToValue(),
		),
	)
	vm.DefBytecode(class, "foo", fn_method0)
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(87, 5, 6), P(95, 5, 14)), "unreachable code"),
			},
		},
		"define an async method": {
			input: `
				async def foo(a: Int, b: Int = 2): Int ! String
					await timeout(5.seconds)
					return a + b
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk"
	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.C_Value
var _ = vm.New
var _ = value.Truthy

func init() { elk.InitNative() }

var sym4 = value.ToSymbol("main")
var sym5 = value.ToSymbol("<main>")

var sym0 = value.ToSymbol("Std::Int")
var sym1 = value.ToSymbol("seconds")
var sym2 = value.ToSymbol("Std::Kernel")
var sym3 = value.ToSymbol("timeout")
var fn_method0 *vm.BytecodeFunction // method: Std::Kernel::foo, loc: <main>:2:5

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

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym5, 1)
	defer thread.PopNativeCallFrame()
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (value.KernelModule).SingletonClass() // Std::Kernel
	fn_method0 = vm.NewBytecodeFunctionWithOptions(
		vm.BytecodeFunctionWithInstructions([]byte{
			byte(bytecode.GET_LOCAL_2),
			byte(bytecode.JUMP_UNLESS_UNDEF), 0x00, 0x02,
			byte(bytecode.INT_2),
			byte(bytecode.SET_LOCAL_2),
			byte(bytecode.GET_LOCAL_3),
			byte(bytecode.PROMISE),
			byte(bytecode.RETURN),
			byte(bytecode.SELF),
			byte(bytecode.INT_5),
			byte(bytecode.CALL_METHOD_NT8), 0x00,
			byte(bytecode.UNDEFINED),
			byte(bytecode.CALL_METHOD_NT8), 0x01,
			byte(bytecode.AWAIT),
			byte(bytecode.AWAIT_RESULT),
			byte(bytecode.POP),
			byte(bytecode.GET_LOCAL_1),
			byte(bytecode.GET_LOCAL_2),
			byte(bytecode.ADD_INT),
			byte(bytecode.RETURN),
		}),
		vm.BytecodeFunctionWithLocation(position.NewLocation("<main>", position.NewSpan(position.New(5, 2, 5), position.New(107, 5, 7)))),
		vm.BytecodeFunctionWithUpvalueCount(0),
		vm.BytecodeFunctionWithStringName("Std::Kernel::foo"),
		vm.BytecodeFunctionWithParameters(3),
		vm.BytecodeFunctionWithOptionalParameters(2),
		vm.BytecodeFunctionWithLineInfoListVar(
			bytecode.NewLineInfo(2, 8),
			bytecode.NewLineInfo(5, 1),
			bytecode.NewLineInfo(3, 10),
			bytecode.NewLineInfo(4, 4),
		),
		vm.BytecodeFunctionWithValuesVar(
			(vm.NewNativeCallSiteInfo((value.GetClass(sym0)).GetMethod(sym1).(*vm.NativeMethod), 0)).ToValue(),
			(vm.NewNativeCallSiteInfo((value.GetSingletonClass(sym2)).GetMethod(sym3).(*vm.NativeMethod), 2)).ToValue(),
		),
	)
	vm.DefBytecode(class, "foo", fn_method0)
}
`,
		},
		"define method with required parameters in top level": {
			input: `
				def foo(a: Int, b: Int)
					c := 5
					a + b + c
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

var sym2 = value.ToSymbol("main")

var sym0 = value.ToSymbol("Std::Kernel::foo")
var sym1 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value, l1 value.Value) (result value.Value, err value.Value) { // method: Std::Kernel::foo, loc: <main>:2:5
	var callFrame *vm.CallFrame
	_ = callFrame
	var l2 value.Value // var c: Std::Int
	_ = l2

	l2 = (value.SmallInt(5)).ToValue()
	return value.AddInts(value.AddInts(l0, l1), l2), value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym2, sym1, 1)
	defer thread.PopNativeCallFrame()
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (value.KernelModule).SingletonClass() // Std::Kernel
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1], args[2])
		return result, err
	}, vm.DefWithParameters(2))
}
`,
		},
		"define method with ivar parameters": {
			input: `
				class Bar
					init(@a: Int, @b: Int); end

					def foo(@a, @b)
						c := 5
						a + b + c
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

var sym4 = value.ToSymbol("main")

var const0 *value.Class // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value, l1 value.Value) (result value.Value, err value.Value) { // method: Bar.:#init, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame

	value.SetInstanceVariable(self, 0, l0)
	value.SetInstanceVariable(self, 1, l1)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Bar.:foo")

func fn_method1(thread *vm.Thread, self value.Value, l0 value.Value, l1 value.Value) (result value.Value, err value.Value) { // method: Bar.:foo, loc: <main>:5:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var l2 value.Value // var c: Std::Int
	_ = l2

	value.SetInstanceVariable(self, 0, l0)
	value.SetInstanceVariable(self, 1, l1)
	l2 = (value.SmallInt(5)).ToValue()
	return value.AddInts(value.AddInts(l0, l1), l2), value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym2, 1)
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
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("a"): 0, symbol.ToSymbol("b"): 1}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Bar
	vm.Def(class, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1], args[2])
		return result, err
	}, vm.DefWithParameters(2))
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0], args[1], args[2])
		return result, err
	}, vm.DefWithParameters(2))
}
`,
		},
		"define method with default ivar parameters": {
			input: `
				class Bar
					var @a: Int?
					var @b: Int?
					def foo(@a: Int = 5, @b: Int = 21)
						c := 5
						a + b + c
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

var sym3 = value.ToSymbol("main")

var const0 *value.Class // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar.:foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, arg_l0 value.Value, arg_l1 value.Value) (result value.Value, err value.Value) { // method: Bar.:foo, loc: <main>:5:6
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var callFrame *vm.CallFrame
	_ = callFrame
	var l2 value.Value // var c: Std::Int
	_ = l2

	if (arg_l0).IsUndefined() {
		l0 = (value.SmallInt(5)).ToValue()
	} else {
		l0 = arg_l0
	}
	value.SetInstanceVariable(self, 0, l0)
	if (arg_l1).IsUndefined() {
		l1 = (value.SmallInt(21)).ToValue()
	} else {
		l1 = arg_l1
	}
	value.SetInstanceVariable(self, 1, l1)
	l2 = (value.SmallInt(5)).ToValue()
	return value.AddInts(value.AddInts(l0, l1), l2), value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
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
	class.IvarIndices = value.IvarIndices{symbol.ToSymbol("a"): 0, symbol.ToSymbol("b"): 1}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Bar
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1], args[2])
		return result, err
	}, vm.DefWithParameters(2))
}
`,
		},
		"define method with optional parameters in top level": {
			input: `
				def foo(a: Int, b: Float = 5.2, c: Int = 10)
					d := 5
					a + b + c + d
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

var sym2 = value.ToSymbol("main")

var sym0 = value.ToSymbol("Std::Kernel::foo")
var sym1 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value, arg_l1 value.Value, arg_l2 value.Value) (result value.Value, err value.Value) { // method: Std::Kernel::foo, loc: <main>:2:5
	var l1 value.Float // var b: Std::Float
	_ = l1
	var l2 value.Value // var c: Std::Int
	_ = l2
	var callFrame *vm.CallFrame
	_ = callFrame
	var l3 value.Value // var d: Std::Int
	_ = l3
	var t1 value.Value
	_ = t1

	if (arg_l1).IsUndefined() {
		l1 = value.Float(5.2)
	} else {
		l1 = (arg_l1).AsFloat()
	}
	if (arg_l2).IsUndefined() {
		l2 = (value.SmallInt(10)).ToValue()
	} else {
		l2 = arg_l2
	}
	l3 = (value.SmallInt(5)).ToValue()
	callFrame.SetNativeLineNumber(4)
	t1, err = value.AddInt(l0, (l1).ToValue())
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	return ((((t1).AsFloat()).AddInt(l2)).AddInt(l3)).ToValue(), value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym2, sym1, 1)
	defer thread.PopNativeCallFrame()
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (value.KernelModule).SingletonClass() // Std::Kernel
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1], args[2], args[3])
		return result, err
	}, vm.DefWithParameters(3))
}
`,
		},
		"define method with required parameters in a class": {
			input: `
				class Bar
					def foo(a: Int, b: Int)
						c := 5
						a + b + c
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

var sym3 = value.ToSymbol("main")

var const0 *value.Class // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar.:foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value, l1 value.Value) (result value.Value, err value.Value) { // method: Bar.:foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var l2 value.Value // var c: Std::Int
	_ = l2

	l2 = (value.SmallInt(5)).ToValue()
	return value.AddInts(value.AddInts(l0, l1), l2), value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Bar
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1], args[2])
		return result, err
	}, vm.DefWithParameters(2))
}
`,
		},
		"define method with required parameters in a module": {
			input: `
				module Bar
					def foo(a: Int, b: Int)
						c := 5
						a + b + c
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

var sym3 = value.ToSymbol("main")

var const0 *value.Module // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value, l1 value.Value) (result value.Value, err value.Value) { // method: Bar::foo, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var l2 value.Value // var c: Std::Int
	_ = l2

	l2 = (value.SmallInt(5)).ToValue()
	return value.AddInts(value.AddInts(l0, l1), l2), value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
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

	class = (const0).SingletonClass() // Bar
	vm.Def(class, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1], args[2])
		return result, err
	}, vm.DefWithParameters(2))
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

func TestGoDefInit(t *testing.T) {
	tests := goTestTable{
		"define init in top level": {
			input: `
				init then :bar
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 2, 5), P(18, 2, 18)), "init definitions cannot appear outside of classes"),
			},
		},
		"define init in a module": {
			input: `
				module Foo
					init then :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(21, 3, 6), P(34, 3, 19)), "init definitions cannot appear outside of classes"),
			},
		},
		"define init in a mixin": {
			input: `
				mixin Foo
					init then :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(20, 3, 6), P(33, 3, 19)), "init definitions cannot appear outside of classes"),
			},
		},
		"define init in an interface": {
			input: `
				interface Foo
					init then :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(24, 3, 6), P(37, 3, 19)), "init definitions cannot appear outside of classes"),
				diagnostic.NewFailure(L(P(24, 3, 6), P(37, 3, 19)), "method `#init` cannot have a body because it is abstract"),
			},
		},
		"define init in a method": {
			input: `
				def foo
					init then :bar
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(18, 3, 6), P(31, 3, 19)), "method definitions cannot appear in this context"),
			},
		},
		"define init in init": {
			input: `
				class Foo
				  init
					  init then :bar
				  end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(33, 4, 8), P(46, 4, 21)), "method definitions cannot appear in this context"),
			},
		},
		"define with required parameters in a class": {
			input: `
				class Bar
					init(a: Int, b: Int)
						c := 5
						a + b + c
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

var sym3 = value.ToSymbol("main")

var const0 *value.Class // Bar
var sym0 = value.ToSymbol("Bar")

var sym1 = value.ToSymbol("Bar.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value, l0 value.Value, l1 value.Value) (result value.Value, err value.Value) { // method: Bar.:#init, loc: <main>:3:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var l2 value.Value // var c: Std::Int
	_ = l2

	l2 = (value.SmallInt(5)).ToValue()
	return self, value.Undefined

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
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 1)
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
	const0 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

	class = const0
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Bar
	vm.Def(class, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0], args[1], args[2])
		return result, err
	}, vm.DefWithParameters(2))
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

// func TestBytecodeDefMixin(t *testing.T) {
// 	tests := bytecodeTestTable{
// 		"mixin with a relative name without a body": {
// 			input: "mixin Foo; end",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				nil,
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(13, 1, 14)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 5),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						nil,
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 2,
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(13, 1, 14)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 7),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"mixin with an absolute name without a body": {
// 			input: "mixin ::Foo; end",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				nil,
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(15, 1, 16)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 5),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						nil,
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 2,
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(15, 1, 16)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 7),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"named mixin inside of a method": {
// 			input: `
// 				def foo
// 					mixin Bar; end
// 				end
// 			`,
// 			err: diagnostic.DiagnosticList{
// 				diagnostic.NewFailure(L(P(18, 3, 6), P(31, 3, 19)), "mixin definitions cannot appear in this context"),
// 			},
// 		},
// 		"mixin with an absolute nested name without a body": {
// 			input: "mixin ::Std::Int::Foo; end",
// 			want: vm.NewBytecodeFunctionNoParams(
// 				nil,
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(25, 1, 26)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 5),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						nil,
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 2,
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(25, 1, 26)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 7),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Std::Int").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"mixin with a body": {
// 			input: `
// 				mixin Foo
// 					a := 1
// 					a + 2
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				nil,
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_CONST8), 1,
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.INIT_NAMESPACE),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(45, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 					bytecode.NewLineInfo(2, 4),
// 					bytecode.NewLineInfo(5, 1),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						nil,
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 2,
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(45, 5, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 5),
// 							bytecode.NewLineInfo(5, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 						},
// 					)),
// 					value.ToSymbol("Foo").ToValue(),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						nil,
// 						value.ToSymbol("<mixin: Foo>"),
// 						[]byte{
// 							byte(bytecode.PREP_LOCALS8), 1,
// 							byte(bytecode.INT_1),
// 							byte(bytecode.SET_LOCAL_1),
// 							byte(bytecode.GET_LOCAL_1),
// 							byte(bytecode.INT_2),
// 							byte(bytecode.ADD_INT),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(5, 2, 5), P(44, 5, 7)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(3, 4),
// 							bytecode.NewLineInfo(4, 3),
// 							bytecode.NewLineInfo(5, 3),
// 						},
// 						nil,
// 					)),
// 				},
// 			),
// 		},
// 		"nested mixins": {
// 			input: `
// 				mixin Foo
// 					mixin Bar
// 						a := 1
// 						a + 2
// 					end
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				nil,
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.GET_CONST8), 1,
// 					byte(bytecode.LOAD_VALUE_2),
// 					byte(bytecode.INIT_NAMESPACE),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(71, 7, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 					bytecode.NewLineInfo(2, 4),
// 					bytecode.NewLineInfo(7, 1),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						nil,
// 						namespaceDefinitionsSymbol,
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 2,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_NAMESPACE), 2,
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(71, 7, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 10),
// 							bytecode.NewLineInfo(7, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Bar").ToValue(),
// 						},
// 					)),
// 					value.ToSymbol("Foo").ToValue(),
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						nil,
// 						value.ToSymbol("<mixin: Foo>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.INIT_NAMESPACE),
// 							byte(bytecode.POP),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(5, 2, 5), P(70, 7, 7)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(3, 4),
// 							bytecode.NewLineInfo(7, 3),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Foo::Bar").ToValue(),
// 							value.Ref(vm.NewBytecodeFunctionNoParams(
// 								nil,
// 								value.ToSymbol("<mixin: Foo::Bar>"),
// 								[]byte{
// 									byte(bytecode.PREP_LOCALS8), 1,
// 									byte(bytecode.INT_1),
// 									byte(bytecode.SET_LOCAL_1),
// 									byte(bytecode.GET_LOCAL_1),
// 									byte(bytecode.INT_2),
// 									byte(bytecode.ADD_INT),
// 									byte(bytecode.POP),
// 									byte(bytecode.NIL),
// 									byte(bytecode.RETURN),
// 								},
// 								L(P(20, 3, 6), P(62, 6, 8)),
// 								bytecode.LineInfoList{
// 									bytecode.NewLineInfo(4, 4),
// 									bytecode.NewLineInfo(5, 3),
// 									bytecode.NewLineInfo(6, 3),
// 								},
// 								nil,
// 							)),
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

// func TestBytecodeInclude(t *testing.T) {
// 	tests := bytecodeTestTable{
// 		"include a global constant in a class": {
// 			input: `
// 				mixin Bar; end
// 				class Foo
// 					include ::Bar
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				nil,
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(60, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						nil,
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 2,
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.GET_CONST8), 3,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.INCLUDE),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(60, 5, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 20),
// 							bytecode.NewLineInfo(5, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Bar").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"include a global constant in a mixin": {
// 			input: `
// 				mixin Bar; end
// 				mixin Foo
// 					include ::Bar
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				nil,
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(60, 5, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 					bytecode.NewLineInfo(5, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						nil,
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 2,
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_NAMESPACE), 2,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.INCLUDE),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(60, 5, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 15),
// 							bytecode.NewLineInfo(5, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Bar").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"include two constants in a class": {
// 			input: `
// 				mixin Bar; end
// 				mixin Baz; end
// 				class Foo
// 					include ::Bar, ::Baz
// 				end
// 			`,
// 			want: vm.NewBytecodeFunctionNoParams(
// 				nil,
// 				mainSymbol,
// 				[]byte{
// 					byte(bytecode.LOAD_VALUE_0),
// 					byte(bytecode.EXEC),
// 					byte(bytecode.POP),
// 					byte(bytecode.NIL),
// 					byte(bytecode.RETURN),
// 				},
// 				L(P(0, 1, 1), P(86, 6, 8)),
// 				bytecode.LineInfoList{
// 					bytecode.NewLineInfo(1, 3),
// 					bytecode.NewLineInfo(6, 2),
// 				},
// 				[]value.Value{
// 					value.Ref(vm.NewBytecodeFunctionNoParams(
// 						nil,
// 						value.ToSymbol("<namespaceDefinitions>"),
// 						[]byte{
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_1),
// 							byte(bytecode.DEF_NAMESPACE), 2,
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_2),
// 							byte(bytecode.DEF_NAMESPACE), 2,
// 							byte(bytecode.GET_CONST8), 0,
// 							byte(bytecode.LOAD_VALUE_3),
// 							byte(bytecode.DEF_NAMESPACE), 1,
// 							byte(bytecode.GET_CONST8), 3,
// 							byte(bytecode.GET_CONST8), 4,
// 							byte(bytecode.SET_SUPERCLASS),
// 							byte(bytecode.GET_CONST8), 3,
// 							byte(bytecode.GET_CONST8), 1,
// 							byte(bytecode.INCLUDE),
// 							byte(bytecode.GET_CONST8), 3,
// 							byte(bytecode.GET_CONST8), 2,
// 							byte(bytecode.INCLUDE),
// 							byte(bytecode.NIL),
// 							byte(bytecode.RETURN),
// 						},
// 						L(P(0, 1, 1), P(86, 6, 8)),
// 						bytecode.LineInfoList{
// 							bytecode.NewLineInfo(1, 30),
// 							bytecode.NewLineInfo(6, 2),
// 						},
// 						[]value.Value{
// 							value.ToSymbol("Root").ToValue(),
// 							value.ToSymbol("Bar").ToValue(),
// 							value.ToSymbol("Baz").ToValue(),
// 							value.ToSymbol("Foo").ToValue(),
// 							value.ToSymbol("Std::Object").ToValue(),
// 						},
// 					)),
// 				},
// 			),
// 		},
// 		"include in top level": {
// 			input: `
// 				mixin Bar; end
// 				include ::Bar
// 			`,
// 			err: diagnostic.DiagnosticList{
// 				diagnostic.NewFailure(
// 					L(P(24, 3, 5), P(36, 3, 17)),
// 					"cannot include mixins in this context",
// 				),
// 			},
// 		},
// 		"include in a module": {
// 			input: `
// 				mixin Bar; end
// 				module Foo
// 					include ::Bar
// 				end
// 			`,
// 			err: diagnostic.DiagnosticList{
// 				diagnostic.NewFailure(
// 					L(P(40, 4, 6), P(52, 4, 18)),
// 					"cannot include mixins in this context",
// 				),
// 			},
// 		},
// 		"include in an interface": {
// 			input: `
// 				mixin Bar; end
// 				interface Foo
// 					include ::Bar
// 				end
// 			`,
// 			err: diagnostic.DiagnosticList{
// 				diagnostic.NewFailure(
// 					L(P(43, 4, 6), P(55, 4, 18)),
// 					"cannot include mixins in this context",
// 				),
// 			},
// 		},
// 		"include in a method": {
// 			input: `
// 				mixin Bar; end
// 				def foo
// 					include ::Bar
// 				end
// 			`,
// 			err: diagnostic.DiagnosticList{
// 				diagnostic.NewFailure(
// 					L(P(37, 4, 6), P(49, 4, 18)),
// 					"cannot include mixins in this context",
// 				),
// 			},
// 		},
// 	}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			bytecodeCompilerTest(tc, t)
// 		})
// 	}
// }
