package vm_test

import (
	"testing"

	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"
)

func TestVMSource_Locals(t *testing.T) {
	tests := sourceTestTable{
		"define and initialise a variable": {
			source:       "var a = 'foo'",
			wantStackTop: value.String("foo"),
		},
		"shadow a variable": {
			source: `
				var a = 10
				var b = do
					var a = 5
					a + 3
				end
				a + b
			`,
			wantStackTop: value.SmallInt(18),
		},
		"define and set a variable": {
			source: `
				var a = 'foo'
				a = a + ' bar'
				a
			`,
			wantStackTop: value.String("foo bar"),
		},
		"try to read an uninitialised variable": {
			source: `
				var a
				a
			`,
			wantCompileErr: errors.ErrorList{
				errors.NewError(L(P(15, 3, 5), P(15, 3, 5)), "cannot access an uninitialised local: a"),
			},
		},
		"try to read a nonexistent variable": {
			source: `
				a
			`,
			wantCompileErr: errors.ErrorList{
				errors.NewError(L(P(5, 2, 5), P(5, 2, 5)), "undeclared variable: a"),
			},
		},
		"increment": {
			source: `
				a := 1
				a++
			`,
			wantStackTop: value.SmallInt(2),
		},
		"decrement": {
			source: `
				a := 1
				a--
			`,
			wantStackTop: value.SmallInt(0),
		},
		"set add": {
			source: `
				a := 1
				a += 2
			`,
			wantStackTop: value.SmallInt(3),
		},
		"set subtract": {
			source: `
				a := 1
				a -= 2
			`,
			wantStackTop: value.SmallInt(-1),
		},
		"set multiply": {
			source: `
				a := 2
				a *= 3
			`,
			wantStackTop: value.SmallInt(6),
		},
		"set divide": {
			source: `
				a := 12
				a /= 3
			`,
			wantStackTop: value.SmallInt(4),
		},
		"set exponentiate": {
			source: `
				a := 12
				a **= 2
			`,
			wantStackTop: value.SmallInt(144),
		},
		"set modulo": {
			source: `
				a := 14
				a %= 3
			`,
			wantStackTop: value.SmallInt(2),
		},
		"set left bitshift": {
			source: `
				a := 14
				a <<= 3
			`,
			wantStackTop: value.SmallInt(112),
		},
		"set logic left bitshift": {
			source: `
				a := 14i8
				a <<<= 3
			`,
			wantStackTop: value.Int8(112),
		},
		"set right bitshift": {
			source: `
				a := 14
				a >>= 2
			`,
			wantStackTop: value.SmallInt(3),
		},
		"set logic right bitshift": {
			source: `
				a := 14i8
				a >>>= 2
			`,
			wantStackTop: value.Int8(3),
		},
		"set bitwise and": {
			source: `
				a := 14
				a &= 5
			`,
			wantStackTop: value.SmallInt(4),
		},
		"set bitwise or": {
			source: `
				a := 14
				a |= 5
			`,
			wantStackTop: value.SmallInt(15),
		},
		"set bitwise xor": {
			source: `
				a := 14
				a ^= 5
			`,
			wantStackTop: value.SmallInt(11),
		},
		"set logic or false": {
			source: `
				a := false
				a ||= 5
			`,
			wantStackTop: value.SmallInt(5),
		},
		"set logic or nil": {
			source: `
				a := nil
				a ||= 5
			`,
			wantStackTop: value.SmallInt(5),
		},
		"set logic or truthy": {
			source: `
				a := 1
				a ||= 5
			`,
			wantStackTop: value.SmallInt(1),
		},
		"set logic and nil": {
			source: `
				a := nil
				a &&= 5
			`,
			wantStackTop: value.Nil,
		},
		"set logic and false": {
			source: `
				a := false
				a &&= 5
			`,
			wantStackTop: value.False,
		},
		"set logic and truthy": {
			source: `
				a := 2
				a &&= 5
			`,
			wantStackTop: value.SmallInt(5),
		},
		"set nil coalesce nil": {
			source: `
				a := nil
				a ??= 5
			`,
			wantStackTop: value.SmallInt(5),
		},
		"set nil coalesce false": {
			source: `
				a := false
				a ??= 5
			`,
			wantStackTop: value.False,
		},
		"set nil coalesce truthy": {
			source: `
				a := 1
				a ??= 5
			`,
			wantStackTop: value.SmallInt(1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_InstanceVariables(t *testing.T) {
	tests := sourceTestTable{
		"read an instance variable of an instance": {
			source: `
				class Foo
				 	setter bar

					def bar then @bar
				end

				f := ::Foo()
				f.bar = "bar value"
				f.bar
			`,
			wantStackTop: value.String("bar value"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
		},
		"read an instance variable of a primitive": {
			source: `
				def foo then @foo
				self.foo()
			`,
			wantRuntimeErr: value.NewError(
				value.PrimitiveValueErrorClass,
				"cannot access instance variables of a primitive value `<GlobalObject>`",
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"set an instance variable of a primitive": {
			source: `
				def foo=(arg) then @foo = arg
				self.foo = 2
			`,
			wantRuntimeErr: value.NewError(
				value.PrimitiveValueErrorClass,
				"cannot set instance variables of a primitive value `<GlobalObject>`",
			),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"set an instance variable of an instance": {
			source: `
				class Foo
				 	getter bar

					def bar=(arg) then @bar = arg
				end

				f := ::Foo()
				f.bar = "bar value"
				f.bar
			`,
			wantStackTop: value.String("bar value"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
		},
		"set an instance variable of a class": {
			source: `
				class Foo
				  singleton
				 		getter bar

						def bar=(arg) then @bar = arg
					end
				end

				::Foo.bar = "bar value"
				::Foo.bar
			`,
			wantStackTop: value.String("bar value"),
			teardown: func() {
				value.RootModule.Constants.DeleteString("Foo")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
