package compiler_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

func TestGoInstanceVariables(t *testing.T) {
	tests := goTestTable{
		"initialise when declared": {
			input: `
				class Foo
					var @a: Int = 3
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(34, 3, 20), P(34, 3, 20)),
					"instance variables cannot be initialised when declared",
				),
			},
		},
		"declare in the top level": {
			input: "var @a: Float",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(0, 1, 1), P(12, 1, 13)),
					"cannot declare instance variable `@a` in this context",
				),
			},
		},
		"declare in a method": {
			input: "def foo; var @a: Float; end",
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(P(9, 1, 10), P(21, 1, 22)),
					"instance variable definitions cannot appear in this context",
				),
			},
		},
		"declare in a class": {
			input: "class Foo; var @a: Float?; end",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}
`,
		},
		"declare in a mixin": {
			input: "mixin Foo; var @a: Float; end",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Mixin // Foo
var sym0 = value.ToSymbol("Foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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
	const0 = value.NewMixin()
	value.AddConstant(parentNamespace, sym0, namespace)

}
`,
		},
		"declare in a module": {
			input: "module Foo; var @a: Float?; end",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}
`,
		},
		"short set add": {
			input: `
				class Foo
					var @foo: Int
					init
						@foo = 3
					end

					def incr_foo
						@foo += 2
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:#init")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:#init, loc: <main>:4:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 4)
	defer thread.PopNativeCallFrame()
	t1 = (value.SmallInt(3)).ToValue()
	value.SetInstanceVariable(self, 0, t1)
	return self, value.Undefined

}

var sym3 = value.ToSymbol("Foo.:incr_foo")

func fn_method1(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:incr_foo, loc: <main>:8:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym3, sym2, 8)
	defer thread.PopNativeCallFrame()
	t1 = value.AddInts(value.GetInstanceVariable(self, 0), (value.SmallInt(2)).ToValue())
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "#init", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
	vm.Def(&class.MethodContainer, "incr_foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method1(thread, args[0])
		return result, err
	})
}
`,
		},
		"set instance variable in a class instance method": {
			input: `
				class Foo
					var @foo: Int?

					def foo
						@foo = 2
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:5:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 5)
	defer thread.PopNativeCallFrame()
	t1 = (value.SmallInt(2)).ToValue()
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},
		"set instance variable from parent in a class instance method": {
			input: `
				class Bar
					var @foo: Int?
				end

				class Foo < Bar
					def foo
						@foo = 2
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")

var const0 *value.Class // Bar
var sym0 = value.ToSymbol("Bar")
var const1 *value.Class // Foo
var sym1 = value.ToSymbol("Foo")

var sym2 = value.ToSymbol("Foo.:foo")
var sym3 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:7:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym2, sym3, 7)
	defer thread.PopNativeCallFrame()
	t1 = (value.SmallInt(2)).ToValue()
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym3, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
	class = const1
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const1 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},
		"set instance variable from parent and self in a class instance method": {
			input: `
				class Bar
					var @bar: Int?
				end

				class Foo < Bar
					var @foo: Int?

					def foo
						@foo = 2
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")

var const0 *value.Class // Bar
var sym0 = value.ToSymbol("Bar")
var const1 *value.Class // Foo
var sym1 = value.ToSymbol("Foo")

var sym2 = value.ToSymbol("Foo.:foo")
var sym3 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:9:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym2, sym3, 9)
	defer thread.PopNativeCallFrame()
	t1 = (value.SmallInt(2)).ToValue()
	value.SetInstanceVariable(self, 1, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym3, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("bar"): 0}
	class = const1
	class.IvarIndices = value.IvarIndices{value.ToSymbol("bar"): 0, value.ToSymbol("foo"): 1}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const1 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},
		"set instance variable from mixin in a class instance method": {
			input: `
				mixin Bar
					var @foo: Int?
				end

				class Foo
					include Bar

					def foo
				  	@foo = 2
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")

var const0 *value.Mixin // Bar
var sym0 = value.ToSymbol("Bar")
var const1 *value.Class // Foo
var sym1 = value.ToSymbol("Foo")

var sym2 = value.ToSymbol("Foo.:foo")
var sym3 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:9:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym2, sym3, 9)
	defer thread.PopNativeCallFrame()
	t1 = (value.SmallInt(2)).ToValue()
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym3, 1)
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

	parentNamespace = (value.RootModule).ToValue()
	const1 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const1)
	value.AddConstant(parentNamespace, sym1, namespace)

	class = const1
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
	class = const1
	mixin = const0
	class.IncludeMixin(mixin)
}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const1
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const1 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"set instance variable in a mixin instance method": {
			input: `
				mixin Foo
					var @foo: Int?

					def foo
				  	@foo = 2
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")

var const0 *value.Mixin // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:foo")
var sym2 = value.ToSymbol("<main>")
var sym3 = value.ToSymbol("foo")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:5:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 5)
	defer thread.PopNativeCallFrame()
	t1 = (value.SmallInt(2)).ToValue()
	err = value.SetInstanceVariableByName(self, sym3, t1)
	if err.IsNotUndefined() {
		panic(err)
	}
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"set instance variable in a module method": {
			input: `
				module Foo
					var @foo: Int?

					def foo
				  	@foo = 2
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo::foo, loc: <main>:5:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 5)
	defer thread.PopNativeCallFrame()
	t1 = (value.SmallInt(2)).ToValue()
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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
	const0 = value.NewModule()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},
		"set instance variable in a class method": {
			input: `
				class Foo
					singleton
						var @foo: Int?

						def foo
							@foo = 2
						end
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo::foo, loc: <main>:6:7
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 6)
	defer thread.PopNativeCallFrame()
	t1 = (value.SmallInt(2)).ToValue()
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0                 // Foo
	class = class.SingletonClass() // &Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"set instance variable in a mixin singleton method": {
			input: `
				mixin Foo
					singleton
						var @foo: Int?

						def foo
							@foo = 2
						end
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var const0 *value.Mixin // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo::foo, loc: <main>:6:7
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 6)
	defer thread.PopNativeCallFrame()
	t1 = (value.SmallInt(2)).ToValue()
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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
	const0 = value.NewMixin()
	value.AddConstant(parentNamespace, sym0, namespace)

}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0                 // Foo
	class = class.SingletonClass() // &Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"set instance variable in an interface singleton method": {
			input: `
				interface Foo
					singleton
						var @foo: Int?

						def foo
							@foo = 2
						end
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var const0 *value.Interface // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo::foo, loc: <main>:6:7
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 6)
	defer thread.PopNativeCallFrame()
	t1 = (value.SmallInt(2)).ToValue()
	value.SetInstanceVariable(self, 0, t1)
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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
	const0 = value.NewInterface()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"set instance variable in a class": {
			input: `
				class Foo
					singleton
						var @a: Int?
					end
					@a = 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym3 = value.ToSymbol("<class: Foo>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame

	self = (const0).ToValue()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 2)
	defer thread.PopNativeCallFrame()
	value.SetInstanceVariable(self, 0, (value.SmallInt(2)).ToValue())
}
`,
		},

		"set instance variable in a mixin": {
			input: `
				mixin Foo
					singleton
						var @a: Int?
					end
					@a = 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Mixin // Foo
var sym0 = value.ToSymbol("Foo")

var sym3 = value.ToSymbol("<mixin: Foo>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)
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
	const0 = value.NewMixin()
	value.AddConstant(parentNamespace, sym0, namespace)

}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame

	self = (const0).ToValue()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 2)
	defer thread.PopNativeCallFrame()
	value.SetInstanceVariable(self, 0, (value.SmallInt(2)).ToValue())
}
`,
		},

		"set instance variable in an interface": {
			input: `
				interface Foo
					singleton
						var @a: Int?
					end
					b := @a = 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Interface // Foo
var sym0 = value.ToSymbol("Foo")

var sym3 = value.ToSymbol("<interface: Foo>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)
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
	const0 = value.NewInterface()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var b: Std::Int
	_ = l0
	var t1 value.Value
	_ = t1

	self = (const0).ToValue()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 2)
	defer thread.PopNativeCallFrame()
	t1 = (value.SmallInt(2)).ToValue()
	value.SetInstanceVariable(self, 0, t1)
	l0 = t1
}
`,
		},

		"set instance variable in a module": {
			input: `
				module Foo
					var @a: Int?
					@a = 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym3 = value.ToSymbol("<module: Foo>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame

	self = value.Ref((const0).SingletonClass())
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 2)
	defer thread.PopNativeCallFrame()
	value.SetInstanceVariable(self, 0, (value.SmallInt(2)).ToValue())
}
`,
		},

		// TODO: fix
		// 		"get box of instance variable in a class instance method": {
		// 			input: `
		// 				class Foo
		// 					var @foo: Int?

		// 					def foo: ^Int?
		// 				  	&@foo
		// 					end
		// 				end
		// 			`,
		// 			want: vm.NewBytecodeFunctionNoParams(
		// 				mainSymbol,
		// 				[]byte{
		// 					byte(bytecode.LOAD_VALUE_0),
		// 					byte(bytecode.EXEC),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.LOAD_VALUE_1),
		// 					byte(bytecode.EXEC),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.LOAD_VALUE_2),
		// 					byte(bytecode.EXEC),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.NIL),
		// 					byte(bytecode.RETURN),
		// 				},
		// 				L(P(0, 1, 1), P(85, 8, 8)),
		// 				bytecode.LineInfoList{
		// 					bytecode.NewLineInfo(1, 9),
		// 					bytecode.NewLineInfo(8, 2),
		// 				},
		// 				[]value.Value{
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						namespaceDefinitionsSymbol,
		// 						[]byte{
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.LOAD_VALUE_1),
		// 							byte(bytecode.DEF_NAMESPACE), 1,
		// 							byte(bytecode.GET_CONST8), 1,
		// 							byte(bytecode.GET_CONST8), 2,
		// 							byte(bytecode.SET_SUPERCLASS),
		// 							byte(bytecode.NIL),
		// 							byte(bytecode.RETURN),
		// 						},
		// 						L(P(0, 1, 1), P(85, 8, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(1, 10),
		// 							bytecode.NewLineInfo(8, 2),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Root").ToValue(),
		// 							value.ToSymbol("Foo").ToValue(),
		// 							value.ToSymbol("Std::Object").ToValue(),
		// 						},
		// 					)),
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						value.ToSymbol("<ivarIndices>"),
		// 						[]byte{
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.LOAD_VALUE_1),
		// 							byte(bytecode.DEF_IVARS),
		// 							byte(bytecode.NIL),
		// 							byte(bytecode.RETURN),
		// 						},
		// 						L(P(0, 1, 1), P(85, 8, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(1, 4),
		// 							bytecode.NewLineInfo(8, 2),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Foo").ToValue(),
		// 							value.Ref(&value.IvarIndices{
		// 								value.ToSymbol("foo"): 0,
		// 							}),
		// 						},
		// 					)),
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						value.ToSymbol("<methodDefinitions>"),
		// 						[]byte{
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.LOAD_VALUE_1),
		// 							byte(bytecode.LOAD_VALUE_2),
		// 							byte(bytecode.DEF_METHOD),
		// 							byte(bytecode.POP),
		// 							byte(bytecode.NIL),
		// 							byte(bytecode.RETURN),
		// 						},
		// 						L(P(0, 1, 1), P(85, 8, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(1, 6),
		// 							bytecode.NewLineInfo(8, 2),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Foo").ToValue(),
		// 							value.Ref(vm.NewBytecodeFunctionNoParams(
		// 								value.ToSymbol("Foo.:foo"),
		// 								[]byte{
		// 									byte(bytecode.INT_0),
		// 									byte(bytecode.CALL_METHOD8), 0,
		// 									byte(bytecode.RETURN),
		// 								},
		// 								L(P(41, 5, 6), P(76, 7, 8)),
		// 								bytecode.LineInfoList{
		// 									bytecode.NewLineInfo(6, 3),
		// 									bytecode.NewLineInfo(7, 1),
		// 								},
		// 								[]value.Value{
		// 									value.Ref(&value.CallSiteInfo{
		// 										Name:          value.ToSymbol("#box_of_ivar_index"),
		// 										ArgumentCount: 1,
		// 									}),
		// 								},
		// 							)),
		// 							value.ToSymbol("foo").ToValue(),
		// 						},
		// 					)),
		// 				},
		// 			),
		// 		},

		// TODO: fix
		// 		"get box of instance variable from parent in a class instance method": {
		// 			input: `
		// 				class Bar
		// 					var @foo: Int?
		// 				end

		// 				class Foo < Bar
		// 					def foo: ^Int?
		// 				  	&@foo
		// 					end
		// 				end
		// 			`,
		// 			want: vm.NewBytecodeFunctionNoParams(
		// 				mainSymbol,
		// 				[]byte{
		// 					byte(bytecode.LOAD_VALUE_0),
		// 					byte(bytecode.EXEC),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.LOAD_VALUE_1),
		// 					byte(bytecode.EXEC),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.LOAD_VALUE_2),
		// 					byte(bytecode.EXEC),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.NIL),
		// 					byte(bytecode.RETURN),
		// 				},
		// 				L(P(0, 1, 1), P(113, 10, 8)),
		// 				bytecode.LineInfoList{
		// 					bytecode.NewLineInfo(1, 9),
		// 					bytecode.NewLineInfo(10, 2),
		// 				},
		// 				[]value.Value{
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						namespaceDefinitionsSymbol,
		// 						[]byte{
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.LOAD_VALUE_1),
		// 							byte(bytecode.DEF_NAMESPACE), 1,
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.LOAD_VALUE_2),
		// 							byte(bytecode.DEF_NAMESPACE), 1,
		// 							byte(bytecode.GET_CONST8), 1,
		// 							byte(bytecode.GET_CONST8), 3,
		// 							byte(bytecode.SET_SUPERCLASS),
		// 							byte(bytecode.GET_CONST8), 2,
		// 							byte(bytecode.GET_CONST8), 1,
		// 							byte(bytecode.SET_SUPERCLASS),
		// 							byte(bytecode.NIL),
		// 							byte(bytecode.RETURN),
		// 						},
		// 						L(P(0, 1, 1), P(113, 10, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(1, 20),
		// 							bytecode.NewLineInfo(10, 2),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Root").ToValue(),
		// 							value.ToSymbol("Bar").ToValue(),
		// 							value.ToSymbol("Foo").ToValue(),
		// 							value.ToSymbol("Std::Object").ToValue(),
		// 						},
		// 					)),
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						value.ToSymbol("<ivarIndices>"),
		// 						[]byte{
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.LOAD_VALUE_1),
		// 							byte(bytecode.DEF_IVARS),
		// 							byte(bytecode.GET_CONST8), 2,
		// 							byte(bytecode.LOAD_VALUE_3),
		// 							byte(bytecode.DEF_IVARS),
		// 							byte(bytecode.NIL),
		// 							byte(bytecode.RETURN),
		// 						},
		// 						L(P(0, 1, 1), P(113, 10, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(1, 8),
		// 							bytecode.NewLineInfo(10, 2),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Bar").ToValue(),
		// 							value.Ref(&value.IvarIndices{
		// 								value.ToSymbol("foo"): 0,
		// 							}),
		// 							value.ToSymbol("Foo").ToValue(),
		// 							value.Ref(&value.IvarIndices{
		// 								value.ToSymbol("foo"): 0,
		// 							}),
		// 						},
		// 					)),
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						value.ToSymbol("<methodDefinitions>"),
		// 						[]byte{
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.LOAD_VALUE_1),
		// 							byte(bytecode.LOAD_VALUE_2),
		// 							byte(bytecode.DEF_METHOD),
		// 							byte(bytecode.POP),
		// 							byte(bytecode.NIL),
		// 							byte(bytecode.RETURN),
		// 						},
		// 						L(P(0, 1, 1), P(113, 10, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(1, 6),
		// 							bytecode.NewLineInfo(10, 2),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Foo").ToValue(),
		// 							value.Ref(vm.NewBytecodeFunctionNoParams(
		// 								value.ToSymbol("Foo.:foo"),
		// 								[]byte{
		// 									byte(bytecode.INT_0),
		// 									byte(bytecode.CALL_METHOD8), 0,
		// 									byte(bytecode.RETURN),
		// 								},
		// 								L(P(69, 7, 6), P(104, 9, 8)),
		// 								bytecode.LineInfoList{
		// 									bytecode.NewLineInfo(8, 3),
		// 									bytecode.NewLineInfo(9, 1),
		// 								},
		// 								[]value.Value{
		// 									value.Ref(&value.CallSiteInfo{
		// 										Name:          value.ToSymbol("#box_of_ivar_index"),
		// 										ArgumentCount: 1,
		// 									}),
		// 								},
		// 							)),
		// 							value.ToSymbol("foo").ToValue(),
		// 						},
		// 					)),
		// 				},
		// 			),
		// 		},

		// TODO: fix
		// 		"get box of instance variable from mixin in a class instance method": {
		// 			input: `
		// 				mixin Bar
		// 					var @foo: Int?
		// 				end

		// 				class Foo
		// 					include Bar

		// 					def foo: ^Int?
		// 				  	&@foo
		// 					end
		// 				end
		// 			`,
		// 			want: vm.NewBytecodeFunctionNoParams(
		// 				mainSymbol,
		// 				[]byte{
		// 					byte(bytecode.LOAD_VALUE_0),
		// 					byte(bytecode.EXEC),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.LOAD_VALUE_1),
		// 					byte(bytecode.EXEC),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.LOAD_VALUE_2),
		// 					byte(bytecode.EXEC),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.NIL),
		// 					byte(bytecode.RETURN),
		// 				},
		// 				L(P(0, 1, 1), P(125, 12, 8)),
		// 				bytecode.LineInfoList{
		// 					bytecode.NewLineInfo(1, 9),
		// 					bytecode.NewLineInfo(12, 2),
		// 				},
		// 				[]value.Value{
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						namespaceDefinitionsSymbol,
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
		// 						L(P(0, 1, 1), P(125, 12, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(1, 20),
		// 							bytecode.NewLineInfo(12, 2),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Root").ToValue(),
		// 							value.ToSymbol("Bar").ToValue(),
		// 							value.ToSymbol("Foo").ToValue(),
		// 							value.ToSymbol("Std::Object").ToValue(),
		// 						},
		// 					)),
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						value.ToSymbol("<ivarIndices>"),
		// 						[]byte{
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.LOAD_VALUE_1),
		// 							byte(bytecode.DEF_IVARS),
		// 							byte(bytecode.NIL),
		// 							byte(bytecode.RETURN),
		// 						},
		// 						L(P(0, 1, 1), P(125, 12, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(1, 4),
		// 							bytecode.NewLineInfo(12, 2),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Foo").ToValue(),
		// 							value.Ref(&value.IvarIndices{
		// 								value.ToSymbol("foo"): 0,
		// 							}),
		// 						},
		// 					)),
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						value.ToSymbol("<methodDefinitions>"),
		// 						[]byte{
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.LOAD_VALUE_1),
		// 							byte(bytecode.LOAD_VALUE_2),
		// 							byte(bytecode.DEF_METHOD),
		// 							byte(bytecode.POP),
		// 							byte(bytecode.NIL),
		// 							byte(bytecode.RETURN),
		// 						},
		// 						L(P(0, 1, 1), P(125, 12, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(1, 6),
		// 							bytecode.NewLineInfo(12, 2),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Foo").ToValue(),
		// 							value.Ref(vm.NewBytecodeFunctionNoParams(
		// 								value.ToSymbol("Foo.:foo"),
		// 								[]byte{
		// 									byte(bytecode.INT_0),
		// 									byte(bytecode.CALL_METHOD8), 0,
		// 									byte(bytecode.RETURN),
		// 								},
		// 								L(P(81, 9, 6), P(116, 11, 8)),
		// 								bytecode.LineInfoList{
		// 									bytecode.NewLineInfo(10, 3),
		// 									bytecode.NewLineInfo(11, 1),
		// 								},
		// 								[]value.Value{
		// 									value.Ref(&value.CallSiteInfo{
		// 										Name:          value.ToSymbol("#box_of_ivar_index"),
		// 										ArgumentCount: 1,
		// 									}),
		// 								},
		// 							)),
		// 							value.ToSymbol("foo").ToValue(),
		// 						},
		// 					)),
		// 				},
		// 			),
		// 		},

		// TODO: fix
		// 		"get box of instance variable in a mixin instance method": {
		// 			input: `
		// 				mixin Foo
		// 					var @foo: Int?

		// 					def foo: ^Int?
		// 				  	&@foo
		// 					end
		// 				end
		// 			`,
		// 			want: vm.NewBytecodeFunctionNoParams(
		// 				mainSymbol,
		// 				[]byte{
		// 					byte(bytecode.LOAD_VALUE_0),
		// 					byte(bytecode.EXEC),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.LOAD_VALUE_1),
		// 					byte(bytecode.EXEC),
		// 					byte(bytecode.POP),
		// 					byte(bytecode.NIL),
		// 					byte(bytecode.RETURN),
		// 				},
		// 				L(P(0, 1, 1), P(85, 8, 8)),
		// 				bytecode.LineInfoList{
		// 					bytecode.NewLineInfo(1, 6),
		// 					bytecode.NewLineInfo(8, 2),
		// 				},
		// 				[]value.Value{
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						namespaceDefinitionsSymbol,
		// 						[]byte{
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.LOAD_VALUE_1),
		// 							byte(bytecode.DEF_NAMESPACE), 2,
		// 							byte(bytecode.NIL),
		// 							byte(bytecode.RETURN),
		// 						},
		// 						L(P(0, 1, 1), P(85, 8, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(1, 5),
		// 							bytecode.NewLineInfo(8, 2),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Root").ToValue(),
		// 							value.ToSymbol("Foo").ToValue(),
		// 						},
		// 					)),
		// 					value.Ref(vm.NewBytecodeFunctionNoParams(
		// 						value.ToSymbol("<methodDefinitions>"),
		// 						[]byte{
		// 							byte(bytecode.GET_CONST8), 0,
		// 							byte(bytecode.LOAD_VALUE_1),
		// 							byte(bytecode.LOAD_VALUE_2),
		// 							byte(bytecode.DEF_METHOD),
		// 							byte(bytecode.POP),
		// 							byte(bytecode.NIL),
		// 							byte(bytecode.RETURN),
		// 						},
		// 						L(P(0, 1, 1), P(85, 8, 8)),
		// 						bytecode.LineInfoList{
		// 							bytecode.NewLineInfo(1, 6),
		// 							bytecode.NewLineInfo(8, 2),
		// 						},
		// 						[]value.Value{
		// 							value.ToSymbol("Foo").ToValue(),
		// 							value.Ref(vm.NewBytecodeFunctionNoParams(
		// 								value.ToSymbol("Foo.:foo"),
		// 								[]byte{
		// 									byte(bytecode.LOAD_VALUE_0),
		// 									byte(bytecode.CALL_METHOD8), 1,
		// 									byte(bytecode.RETURN),
		// 								},
		// 								L(P(41, 5, 6), P(76, 7, 8)),
		// 								bytecode.LineInfoList{
		// 									bytecode.NewLineInfo(6, 3),
		// 									bytecode.NewLineInfo(7, 1),
		// 								},
		// 								[]value.Value{
		// 									value.ToSymbol("foo").ToValue(),
		// 									value.Ref(&value.CallSiteInfo{
		// 										Name:          value.ToSymbol("#box_of_ivar_name"),
		// 										ArgumentCount: 1,
		// 									}),
		// 								},
		// 							)),
		// 							value.ToSymbol("foo").ToValue(),
		// 						},
		// 					)),
		// 				},
		// 			),
		// 		},

		"get instance variable in a class instance method": {
			input: `
				class Foo
					var @foo: Int?

					def foo
				  	@foo
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:5:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 5)
	defer thread.PopNativeCallFrame()
	return value.GetInstanceVariable(self, 0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"get instance variable from parent in a class instance method": {
			input: `
				class Bar
					var @foo: Int?
				end

				class Foo < Bar
					def foo
				  	@foo
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")

var const0 *value.Class // Bar
var sym0 = value.ToSymbol("Bar")
var const1 *value.Class // Foo
var sym1 = value.ToSymbol("Foo")

var sym2 = value.ToSymbol("Foo.:foo")
var sym3 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:7:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym2, sym3, 7)
	defer thread.PopNativeCallFrame()
	return value.GetInstanceVariable(self, 0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym3, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
	class = const1
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const1 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"get instance variable from parent and self in a class instance method": {
			input: `
				class Bar
					var @bar: Int?
				end

				class Foo < Bar
					var @foo: Int?

					def foo
				  	@bar.must + @foo.must
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")

var const0 *value.Class // Bar
var sym0 = value.ToSymbol("Bar")
var const1 *value.Class // Foo
var sym1 = value.ToSymbol("Foo")

var sym2 = value.ToSymbol("Foo.:foo")
var sym3 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:9:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym2, sym3, 9)
	defer thread.PopNativeCallFrame()
	t1 = value.GetInstanceVariable(self, 0)
	callFrame.SetNativeLineNumber(10)
	err = value.Must(t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	t1 = value.GetInstanceVariable(self, 1)
	err = value.Must(t1)
	if err.IsNotUndefined() {
		thread.CaptureStackTrace()
		return result, err
	}
	return value.AddInts(t1, t1), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym3, 1)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const0
	class.IvarIndices = value.IvarIndices{value.ToSymbol("bar"): 0}
	class = const1
	class.IvarIndices = value.IvarIndices{value.ToSymbol("bar"): 0, value.ToSymbol("foo"): 1}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const1 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"get instance variable from mixin in a class instance method": {
			input: `
				mixin Bar
					var @foo: Int?
				end

				class Foo
					include Bar

					def foo
				  	@foo
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")

var const0 *value.Mixin // Bar
var sym0 = value.ToSymbol("Bar")
var const1 *value.Class // Foo
var sym1 = value.ToSymbol("Foo")

var sym2 = value.ToSymbol("Foo.:foo")
var sym3 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:9:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym2, sym3, 9)
	defer thread.PopNativeCallFrame()
	return value.GetInstanceVariable(self, 0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)

	methodDefinitions()
	callFrame = thread.AddNativeCallFrame(sym4, sym3, 1)
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

	parentNamespace = (value.RootModule).ToValue()
	const1 = value.NewClassWithOptions(value.ClassWithSuperclass(nil))
	namespace = value.Ref(const1)
	value.AddConstant(parentNamespace, sym1, namespace)

	class = const1
	superclass = value.ObjectClass
	class.SetSuperclass(superclass)
	class = const1
	mixin = const0
	class.IncludeMixin(mixin)
}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = const1
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const1 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"get instance variable in a mixin instance method": {
			input: `
				mixin Foo
					var @foo: Int?

					def foo
				  	@foo
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym4 = value.ToSymbol("main")

var const0 *value.Mixin // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo.:foo")
var sym2 = value.ToSymbol("<main>")
var sym3 = value.ToSymbol("foo")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo.:foo, loc: <main>:5:6
	var callFrame *vm.CallFrame
	_ = callFrame
	var t1 value.Value
	_ = t1

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 5)
	defer thread.PopNativeCallFrame()
	t1, err = value.GetInstanceVariableByName(self, sym3)
	if err.IsNotUndefined() {
		panic(err)
	}
	return t1, value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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

	class = const0 // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"get instance variable in a module method": {
			input: `
				module Foo
					var @foo: Int?

					def foo
				  	@foo
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo::foo, loc: <main>:5:6
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 5)
	defer thread.PopNativeCallFrame()
	return value.GetInstanceVariable(self, 0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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
	const0 = value.NewModule()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"get instance variable in a class method": {
			input: `
				class Foo
					singleton
						var @foo: Int?

						def foo
				  		@foo
						end
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo::foo, loc: <main>:6:7
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 6)
	defer thread.PopNativeCallFrame()
	return value.GetInstanceVariable(self, 0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0                 // Foo
	class = class.SingletonClass() // &Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"get instance variable in a mixin singleton method": {
			input: `
				mixin Foo
					singleton
						var @foo: Int?

						def foo
				  		@foo
						end
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var const0 *value.Mixin // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo::foo, loc: <main>:6:7
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 6)
	defer thread.PopNativeCallFrame()
	return value.GetInstanceVariable(self, 0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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
	const0 = value.NewMixin()
	value.AddConstant(parentNamespace, sym0, namespace)

}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = const0                 // Foo
	class = class.SingletonClass() // &Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"get instance variable in an interface singleton method": {
			input: `
				interface Foo
					singleton
						var @foo: Int?

						def foo
				  		@foo
						end
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym3 = value.ToSymbol("main")

var const0 *value.Interface // Foo
var sym0 = value.ToSymbol("Foo")

var sym1 = value.ToSymbol("Foo::foo")
var sym2 = value.ToSymbol("<main>")

func fn_method0(thread *vm.Thread, self value.Value) (result value.Value, err value.Value) { // method: Foo::foo, loc: <main>:6:7
	var callFrame *vm.CallFrame
	_ = callFrame

	callFrame = thread.AddNativeCallFrame(sym1, sym2, 6)
	defer thread.PopNativeCallFrame()
	return value.GetInstanceVariable(self, 0), value.Undefined

}

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
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
	const0 = value.NewInterface()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("foo"): 0}
}

func methodDefinitions() {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass() // Foo
	vm.Def(&class.MethodContainer, "foo", func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) {
		result, err := fn_method0(thread, args[0])
		return result, err
	})
}
`,
		},

		"get instance variable in a class": {
			input: `
				class Foo
					singleton
						var @a: Int?
					end
					b := @a
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Class // Foo
var sym0 = value.ToSymbol("Foo")

var sym3 = value.ToSymbol("<class: Foo>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var b: Std::Int?
	_ = l0

	self = (const0).ToValue()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 2)
	defer thread.PopNativeCallFrame()
	l0 = value.GetInstanceVariable(self, 0)
}
`,
		},

		"get instance variable in a mixin": {
			input: `
				mixin Foo
					singleton
						var @a: Int?
					end
					b := @a
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Mixin // Foo
var sym0 = value.ToSymbol("Foo")

var sym3 = value.ToSymbol("<mixin: Foo>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)
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
	const0 = value.NewMixin()
	value.AddConstant(parentNamespace, sym0, namespace)

}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var b: Std::Int?
	_ = l0

	self = (const0).ToValue()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 2)
	defer thread.PopNativeCallFrame()
	l0 = value.GetInstanceVariable(self, 0)
}
`,
		},

		"get instance variable in an interface": {
			input: `
				interface Foo
					singleton
						var @a: Int?
					end
					b := @a
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Interface // Foo
var sym0 = value.ToSymbol("Foo")

var sym3 = value.ToSymbol("<interface: Foo>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)
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
	const0 = value.NewInterface()
	namespace = value.Ref(const0)
	value.AddConstant(parentNamespace, sym0, namespace)

}
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var b: Std::Int?
	_ = l0

	self = (const0).ToValue()
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 2)
	defer thread.PopNativeCallFrame()
	l0 = value.GetInstanceVariable(self, 0)
}
`,
		},

		"get instance variable in a module": {
			input: `
				module Foo
					var @a: Int?
					b := @a
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym1 = value.ToSymbol("main")
var sym2 = value.ToSymbol("<main>")

var const0 *value.Module // Foo
var sym0 = value.ToSymbol("Foo")

var sym3 = value.ToSymbol("<module: Foo>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)

	initGlobalEnv()

	ivarIndices(thread)
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
func ivarIndices(thread *vm.Thread) {
	var class *value.Class
	_ = class

	class = (const0).SingletonClass()
	class.IvarIndices = value.IvarIndices{value.ToSymbol("a"): 0}
}

func fn_ns_expr0(thread *vm.Thread) { // namespace: Foo, loc: <main>:2:5
	var self value.Value
	_ = self
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var b: Std::Int?
	_ = l0

	self = value.Ref((const0).SingletonClass())
	callFrame = thread.AddNativeCallFrame(sym3, sym2, 2)
	defer thread.PopNativeCallFrame()
	l0 = value.GetInstanceVariable(self, 0)
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

func TestGoLocalVariables(t *testing.T) {
	tests := goTestTable{
		"declare": {
			input: "var a: Int",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
}
`,
		},
		"declare and initialise": {
			input: "var a = 3",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(3)).ToValue()
}
`,
		},
		// "declare with a pattern": {
		// 	input: "var [1, a] = [1, 2]",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.LOAD_VALUE_0),
		// 			byte(bytecode.COPY),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.IS_A),
		// 			byte(bytecode.JUMP_UNLESS_NP), 0, 33,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.CALL_METHOD8), 2,
		// 			byte(bytecode.INT_2),
		// 			byte(bytecode.EQUAL_INT),
		// 			byte(bytecode.JUMP_UNLESS_NP), 0, 24,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.INT_0),
		// 			byte(bytecode.SUBSCRIPT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.EQUAL),
		// 			byte(bytecode.POP_SKIP_ONE),
		// 			byte(bytecode.JUMP_UNLESS_NP), 0, 13,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SUBSCRIPT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.TRUE),
		// 			byte(bytecode.POP_SKIP_ONE),
		// 			byte(bytecode.JUMP_UNLESS_NP), 0, 2,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.TRUE),
		// 			byte(bytecode.JUMP_IF), 0, 2,
		// 			byte(bytecode.LOAD_VALUE_3),
		// 			byte(bytecode.THROW),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(18, 1, 19)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 49),
		// 		},
		// 		[]value.Value{
		// 			value.Ref(&value.ArrayListOfValue{value.SmallInt(1).ToValue(), value.SmallInt(2).ToValue()}),
		// 			value.Ref(value.ListMixin),
		// 			value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
		// 			value.Ref(value.NewError(value.PatternNotMatchedErrorClass, "assigned value does not match the pattern defined in variable declaration")),
		// 		},
		// 	),
		// },

		"read undeclared": {
			input: "a",
			want:  ``,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(0, 1, 1)), "undefined local `a`"),
			},
		},
		"assign undeclared": {
			input: "a = 3",
			want: `
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(0, 1, 1)), "undefined local `a`"),
			},
		},
		"assign uninitialised": {
			input: `
				var a: String
				a = 'foo'
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.String // var a: Std::String
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.String("foo")
}
`,
		},
		"assign initialised": {
			input: `
				var a = 'foo'
				a = 'bar'
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.String // var a: Std::String
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.String("foo")
	l0 = value.String("bar")
}
`,
		},
		"short add set": {
			input: `
				var a = 'foo'
				a += 'bar'
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.String // var a: Std::String
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.String("foo")
	l0 = l0.ConcatString(value.String("bar"))
}
`,
		},
		"read uninitialised": {
			input: `
				var a: Int
				a + 2
			`,
			want: ``,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"read initialised": {
			input: `
				var a = 5
				b := a + 2
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	l1 = value.AddInts(l0, (value.SmallInt(2)).ToValue())
}
`,
		},
		"read initialised in child scope": {
			input: `
				var a = 5
				b := do
					a + 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	l1 = value.AddInts(l0, (value.SmallInt(2)).ToValue())
}
`,
		},
		"shadow in child scope": {
			input: `
				var a = 5
				b := 2 + do
					var a = 10
					a + 12
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var a: Std::Int
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	l2 = (value.SmallInt(10)).ToValue()
	l1 = (value.SmallInt(2)).AddInt(value.AddInts(l2, (value.SmallInt(12)).ToValue()))
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

func TestGoUpvalue(t *testing.T) {
	tests := goTestTable{
		"create a pointer": {
			input: `
				a := 5
				b := &a
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Box // var b: Std::Box[Std::Int]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	l1 = (*value.BoxOfValue)(&l0)
}
`,
		},
		"create a pointer to an immutable local": {
			input: `
				val a = 5
				b := &a
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 5
	_ = l0
	var l1 value.ImmutableBox // var b: Std::ImmutableBox[5]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(5)
	l1 = value.NewImmutableNativeBox(&l0)
}
`,
		},
		"create a pointer to a native value": {
			input: `
				a := 5.2
				b := &a
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: Std::Float
	_ = l0
	var l1 value.Box // var b: Std::Box[Std::Float]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(5.2)
	l1 = value.NewNativeBox(&l0)
}
`,
		},
		"create a pointer to an immutable native local": {
			input: `
				val a = 5.5
				b := &a
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Float // var a: 5.5
	_ = l0
	var l1 value.ImmutableBox // var b: Std::ImmutableBox[5.5]
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.Float(5.5)
	l1 = value.NewImmutableNativeBox(&l0)
}
`,
		},
		"in child scope": {
			input: `
				upvalue := 5
				fn := -> println upvalue
			`,
			want: `package main

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("<closure>")
var sym3 = value.ToSymbol("println@1")
var fn_method0 vm.NativeFunction // Std::Kernel::println@1

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var upvalue: Std::Int
	_ = l0
	var l1 vm.Closure // var fn: %||: void
	_ = l1
	var t1 *vm.NativeClosure
	_ = t1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	fn_method0 = vm.MethodToFunc(((value.KernelModule).SingletonClass()).LookupMethod(sym3))

	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(5)).ToValue()
	t1 = vm.NewNativeClosure(
		func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) { // name: fn_cl0, sig: %||: void, loc: <main>:3:11
			var callFrame *vm.CallFrame
			_ = callFrame
			var t1 value.Value
			_ = t1
			var t2 []value.Value
			_ = t2
			var err value.Value
			_ = err

			callFrame = thread.AddNativeCallFrame(sym2, sym1, 3)
			defer thread.PopNativeCallFrame()
			t2 = value.ResizeNativeArgs(t2, 3)
			t2[0] = (value.KernelModule).ToValue()
			t2[1] = l0
			t1, err = fn_method0(thread, t2) // receiver: Std::Kernel, name: println@1
			if err.IsNotUndefined() {
				thread.CaptureStackTrace()
				return value.Undefined, err
			}
			return t1, value.Undefined
		},
		0,
		position.NewLocation("<main>", position.NewSpan(position.New(28, 3, 11), position.New(28, 3, 11))),
	)
	l1 = t1
}
`,
		},
		"upvalues from multiple scopes": {
			input: `
				a := 2
				b := 10

				if b == 40
					g := 5
					e := 9
					f := 20
					h := -> g + f
					if a == 20
						a := 10
						b := 50
						c := 30
						d := -> a + b + c
					end
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")
var sym2 = value.ToSymbol("<closure>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.Value // var g: Std::Int
	_ = l2
	var l3 value.Value // var e: Std::Int
	_ = l3
	var l4 value.Value // var f: Std::Int
	_ = l4
	var l5 vm.Closure // var h: %||: Std::Int
	_ = l5
	var t1 *vm.NativeClosure
	_ = t1
	var l6 value.Value // var a: Std::Int
	_ = l6
	var l7 value.Value // var b: Std::Int
	_ = l7
	var l8 value.Value // var c: Std::Int
	_ = l8
	var l9 vm.Closure // var d: %||: Std::Int
	_ = l9
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = (value.SmallInt(2)).ToValue()
	l1 = (value.SmallInt(10)).ToValue()
	if value.Bool(value.EqualInts(l1, (value.SmallInt(40)).ToValue())) {
		l2 = (value.SmallInt(5)).ToValue()
		l3 = (value.SmallInt(9)).ToValue()
		l4 = (value.SmallInt(20)).ToValue()
		t1 = vm.NewNativeClosure(
			func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) { // name: fn_cl0, sig: %||: Std::Int, loc: <main>:9:11
				var callFrame *vm.CallFrame
				_ = callFrame

				callFrame = thread.AddNativeCallFrame(sym2, sym1, 9)
				defer thread.PopNativeCallFrame()
				return value.AddInts(l2, l4), value.Undefined
			},
			0,
			position.NewLocation("<main>", position.NewSpan(position.New(87, 9, 11), position.New(87, 9, 11))),
		)
		l5 = t1
		if value.Bool(value.EqualInts(l0, (value.SmallInt(20)).ToValue())) {
			l6 = (value.SmallInt(10)).ToValue()
			l7 = (value.SmallInt(50)).ToValue()
			l8 = (value.SmallInt(30)).ToValue()
			t1 = vm.NewNativeClosure(
				func(thread *vm.Thread, args []value.Value) (value.Value, value.Value) { // name: fn_cl1, sig: %||: Std::Int, loc: <main>:14:12
					var callFrame *vm.CallFrame
					_ = callFrame

					callFrame = thread.AddNativeCallFrame(sym2, sym1, 14)
					defer thread.PopNativeCallFrame()
					return value.AddInts(value.AddInts(l6, l7), l8), value.Undefined
				},
				0,
				position.NewLocation("<main>", position.NewSpan(position.New(165, 14, 12), position.New(165, 14, 12))),
			)
			l9 = t1
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

func TestGoLocalValues(t *testing.T) {
	tests := goTestTable{
		"declare without initialisation": {
			input: "val a: Int",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.Value // var a: Std::Int
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(0, 1, 1), P(9, 1, 10)), "a value must be initialised on declaration `a`"),
			},
		},
		"declare and initialise": {
			input: "val a = 3",
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 3
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(3)
}
`,
		},
		// "declare with a pattern": {
		// 	input: "val [1, a] = [1, 2]",
		// 	want: vm.NewBytecodeFunctionNoParams(
		// 		mainSymbol,
		// 		[]byte{
		// 			byte(bytecode.PREP_LOCALS8), 1,
		// 			byte(bytecode.LOAD_VALUE_0),
		// 			byte(bytecode.COPY),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.LOAD_VALUE_1),
		// 			byte(bytecode.IS_A),
		// 			byte(bytecode.JUMP_UNLESS_NP), 0, 33,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.CALL_METHOD8), 2,
		// 			byte(bytecode.INT_2),
		// 			byte(bytecode.EQUAL_INT),
		// 			byte(bytecode.JUMP_UNLESS_NP), 0, 24,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.INT_0),
		// 			byte(bytecode.SUBSCRIPT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.EQUAL),
		// 			byte(bytecode.POP_SKIP_ONE),
		// 			byte(bytecode.JUMP_UNLESS_NP), 0, 13,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.INT_1),
		// 			byte(bytecode.SUBSCRIPT),
		// 			byte(bytecode.DUP),
		// 			byte(bytecode.SET_LOCAL_1),
		// 			byte(bytecode.TRUE),
		// 			byte(bytecode.POP_SKIP_ONE),
		// 			byte(bytecode.JUMP_UNLESS_NP), 0, 2,
		// 			byte(bytecode.POP),
		// 			byte(bytecode.TRUE),
		// 			byte(bytecode.JUMP_IF), 0, 2,
		// 			byte(bytecode.LOAD_VALUE_3),
		// 			byte(bytecode.THROW),
		// 			byte(bytecode.RETURN),
		// 		},
		// 		L(P(0, 1, 1), P(18, 1, 19)),
		// 		bytecode.LineInfoList{
		// 			bytecode.NewLineInfo(1, 49),
		// 		},
		// 		[]value.Value{
		// 			value.Ref(&value.ArrayListOfValue{value.SmallInt(1).ToValue(), value.SmallInt(2).ToValue()}),
		// 			value.Ref(value.ListMixin),
		// 			value.Ref(value.NewCallSiteInfo(value.ToSymbol("length"), 0)),
		// 			value.Ref(value.NewError(value.PatternNotMatchedErrorClass, "assigned value does not match the pattern defined in value declaration")),
		// 		},
		// 	),
		// },
		"assign uninitialised": {
			input: `
				val a: String
				a = 'foo'
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.String // var a: Std::String
	_ = l0
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.String("foo")
}
`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 2, 5), P(17, 2, 17)), "a value must be initialised on declaration `a`"),
			},
		},
		"assign initialised": {
			input: `
				val a = 'foo'
				a = 'bar'
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(23, 3, 5), P(23, 3, 5)), "local value `a` cannot be reassigned"),
				diagnostic.NewFailure(L(P(27, 3, 9), P(31, 3, 13)), "type `\"bar\"` cannot be assigned to type `\"foo\"`"),
			},
		},
		"read uninitialised": {
			input: `
				val a: Int
				a + 2
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(5, 2, 5), P(14, 2, 14)), "a value must be initialised on declaration `a`"),
				diagnostic.NewFailure(L(P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"read initialised": {
			input: `
				val a = 5
				b := a + 2
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 5
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(5)
	l1 = (l0).AddSmallInt(value.SmallInt(2))
}
`,
		},
		"read initialised in child scope": {
			input: `
				val a = 5
				b := do
					a + 2
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 5
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(5)
	l1 = (l0).AddSmallInt(value.SmallInt(2))
}
`,
		},
		"shadow in child scope": {
			input: `
				val a = 5
				b := 2 + do
					val a = 10
					a + 12
				end
			`,
			want: `package main

import (
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

var _ = symbol.Value
var _ = vm.New
var _ = value.Truthy

var sym0 = value.ToSymbol("main")
var sym1 = value.ToSymbol("<main>")

func main() { // loc: <main>
	thread := vm.New()
	_ = thread
	var callFrame *vm.CallFrame
	_ = callFrame
	var l0 value.SmallInt // var a: 5
	_ = l0
	var l1 value.Value // var b: Std::Int
	_ = l1
	var l2 value.SmallInt // var a: 10
	_ = l2
	var self value.Value
	_ = self

	self = value.Ref(value.GlobalObject)
	callFrame = thread.AddNativeCallFrame(sym0, sym1, 1)
	defer thread.PopNativeCallFrame()
	l0 = value.SmallInt(5)
	l2 = value.SmallInt(10)
	l1 = (value.SmallInt(2)).AddInt((l2).AddSmallInt(value.SmallInt(12)))
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
