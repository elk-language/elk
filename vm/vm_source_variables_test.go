package vm_test

import (
	"testing"

	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/value"
)

func TestVMSource_Variables(t *testing.T) {
	tests := sourceTestTable{
		"define and initialise a variable": {
			source:       "var a = 'foo'",
			wantStackTop: value.Ref(value.String("foo")),
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
			wantStackTop: value.SmallInt(18).ToValue(),
		},
		"define and set a variable": {
			source: `
				var a = 'foo'
				a = a + ' bar'
				a
			`,
			wantStackTop: value.Ref(value.String("foo bar")),
		},
		"define variables with a pattern": {
			source: `
				var [1, a] = [1, 25]
				a
			`,
			wantStackTop: value.SmallInt(25).ToValue(),
		},
		"override variables with a pattern": {
			source: `
				var a = 5
				var b = -7
				var [b, a] = [a, b]
				[a, b]
			`,
			wantStackTop: value.Ref(&value.ArrayList{
				value.SmallInt(-7).ToValue(),
				value.SmallInt(5).ToValue(),
			}),
		},
		"define variables with a pattern that does not match": {
			source: `
				var [1, 2, a] = [1, 25]
				a
			`,
			wantStackTop: value.SmallInt(25).ToValue(),
			wantRuntimeErr: value.Ref(value.NewError(
				value.PatternNotMatchedErrorClass,
				"assigned value does not match the pattern defined in variable declaration",
			)),
		},
		"try to read an uninitialised variable": {
			source: `
				var a: String
				a
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(23, 3, 5), P(23, 3, 5)), "cannot access uninitialised local `a`"),
			},
		},
		"try to read a nonexistent variable": {
			source: `
				a
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(5, 2, 5), P(5, 2, 5)), "undefined local `a`"),
			},
		},
		"increment": {
			source: `
				a := 1
				a++
			`,
			wantStackTop: value.SmallInt(2).ToValue(),
		},
		"decrement": {
			source: `
				a := 1
				a--
			`,
			wantStackTop: value.SmallInt(0).ToValue(),
		},
		"set add": {
			source: `
				a := 1
				a += 2
			`,
			wantStackTop: value.SmallInt(3).ToValue(),
		},
		"set subtract": {
			source: `
				a := 1
				a -= 2
			`,
			wantStackTop: value.SmallInt(-1).ToValue(),
		},
		"set multiply": {
			source: `
				a := 2
				a *= 3
			`,
			wantStackTop: value.SmallInt(6).ToValue(),
		},
		"set divide": {
			source: `
				a := 12
				a /= 3
			`,
			wantStackTop: value.SmallInt(4).ToValue(),
		},
		"set exponentiate": {
			source: `
				a := 12
				a **= 2
			`,
			wantStackTop: value.SmallInt(144).ToValue(),
		},
		"set modulo": {
			source: `
				a := 14
				a %= 3
			`,
			wantStackTop: value.SmallInt(2).ToValue(),
		},
		"set left bitshift": {
			source: `
				a := 14
				a <<= 3
			`,
			wantStackTop: value.SmallInt(112).ToValue(),
		},
		"set logic left bitshift": {
			source: `
				a := 14i8
				a <<<= 3
			`,
			wantStackTop: value.Int8(112).ToValue(),
		},
		"set right bitshift": {
			source: `
				a := 14
				a >>= 2
			`,
			wantStackTop: value.SmallInt(3).ToValue(),
		},
		"set logic right bitshift": {
			source: `
				a := 14i8
				a >>>= 2
			`,
			wantStackTop: value.Int8(3).ToValue(),
		},
		"set bitwise and": {
			source: `
				a := 14
				a &= 5
			`,
			wantStackTop: value.SmallInt(4).ToValue(),
		},
		"set bitwise or": {
			source: `
				a := 14
				a |= 5
			`,
			wantStackTop: value.SmallInt(15).ToValue(),
		},
		"set bitwise xor": {
			source: `
				a := 14
				a ^= 5
			`,
			wantStackTop: value.SmallInt(11).ToValue(),
		},
		"set logic or false": {
			source: `
				var a: Int | bool = false
				a ||= 5
			`,
			wantStackTop: value.SmallInt(5).ToValue(),
		},
		"set logic or nil": {
			source: `
				var a: Int? = nil
				a ||= 5
			`,
			wantStackTop: value.SmallInt(5).ToValue(),
		},
		"set logic or truthy": {
			source: `
				a := 1
				a ||= 5
			`,
			wantStackTop: value.SmallInt(1).ToValue(),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(16, 3, 5), P(16, 3, 5)), "this condition will always have the same result since type `Std::Int` is truthy"),
				error.NewWarning(L(P(22, 3, 11), P(22, 3, 11)), "unreachable code"),
			},
		},
		"set logic and nil": {
			source: `
				var a: Int? = nil
				a &&= 5
			`,
			wantStackTop: value.Nil,
		},
		"set logic and false": {
			source: `
				var a: Int | bool = false
				a &&= 5
			`,
			wantStackTop: value.False,
		},
		"set logic and truthy": {
			source: `
				a := 2
				a &&= 5
			`,
			wantStackTop: value.SmallInt(5).ToValue(),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(16, 3, 5), P(16, 3, 5)), "this condition will always have the same result since type `Std::Int` is truthy"),
			},
		},
		"set nil coalesce nil": {
			source: `
				var a: Int? = nil
				a ??= 5
			`,
			wantStackTop: value.SmallInt(5).ToValue(),
		},
		"set nil coalesce false": {
			source: `
				var a: Int | bool = false
				a ??= 5
			`,
			wantStackTop: value.False,
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(35, 3, 5), P(35, 3, 5)), "this condition will always have the same result since type `Std::Int | bool` can never be nil"),
				error.NewWarning(L(P(41, 3, 11), P(41, 3, 11)), "unreachable code"),
			},
		},
		"set nil coalesce truthy": {
			source: `
				a := 1
				a ??= 5
			`,
			wantStackTop: value.SmallInt(1).ToValue(),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(16, 3, 5), P(16, 3, 5)), "this condition will always have the same result since type `Std::Int` can never be nil"),
				error.NewWarning(L(P(22, 3, 11), P(22, 3, 11)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Values(t *testing.T) {
	tests := sourceTestTable{
		"define and initialise": {
			source:       "val a = 'foo'",
			wantStackTop: value.Ref(value.String("foo")),
		},
		"shadow": {
			source: `
				val a = 10
				val b = do
					var a = 5
					a + 3
				end
				a + b
			`,
			wantStackTop: value.SmallInt(18).ToValue(),
		},
		"define and set": {
			source: `
				val a = 'foo'
				a = a + ' bar'
				a
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(23, 3, 5), P(23, 3, 5)), "local value `a` cannot be reassigned"),
				error.NewFailure(L(P(27, 3, 9), P(36, 3, 18)), "type `Std::String` cannot be assigned to type `\"foo\"`"),
			},
		},
		"define variables with a pattern": {
			source: `
				val [1, a] = [1, 25]
				a
			`,
			wantStackTop: value.SmallInt(25).ToValue(),
		},
		"override variables with a pattern": {
			source: `
				val a = 5
				val b = -7
				val [b, a] = [a, b]
				[a, b]
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(39, 4, 10), P(39, 4, 10)), "local value `b` cannot be reassigned"),
				error.NewFailure(L(P(42, 4, 13), P(42, 4, 13)), "local value `a` cannot be reassigned"),
			},
		},
		"define with a pattern that does not match": {
			source: `
				val [1, 2, a] = [1, 25]
				a
			`,
			wantStackTop: value.SmallInt(25).ToValue(),
			wantRuntimeErr: value.Ref(value.NewError(
				value.PatternNotMatchedErrorClass,
				"assigned value does not match the pattern defined in value declaration",
			)),
		},
		"try to read uninitialised": {
			source: `
				val a: Int
				a
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(20, 3, 5), P(20, 3, 5)), "cannot access uninitialised local `a`"),
			},
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
				 	setter bar: String?

					def bar: String? then @bar
				end

				f := ::Foo()
				f.bar = "bar value"
				f.bar
			`,
			wantStackTop: value.Ref(value.String("bar value")),
		},
		"set an instance variable of an instance": {
			source: `
				class Foo
				 	getter bar: String?

					def bar=(arg: String?) then @bar = arg
				end

				f := ::Foo()
				f.bar = "bar value"
				f.bar
			`,
			wantStackTop: value.Ref(value.String("bar value")),
		},
		"set an instance variable of a class": {
			source: `
				class Foo
				  singleton
				 		getter bar: String?

						def bar=(arg: String?) then @bar = arg
					end
				end

				::Foo.bar = "bar value"
				::Foo.bar
			`,
			wantStackTop: value.Ref(value.String("bar value")),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
