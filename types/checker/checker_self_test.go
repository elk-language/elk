package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestSelfType(t *testing.T) {
	tests := testTable{
		"use self in a variable declaration in top level": {
			input: `
				var a: self
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 2, 12), P(15, 2, 15)), "type `self` can appear only in method throw types, method return types and method bodies"),
			},
		},
		"use self in a variable declaration in a class body": {
			input: `
				class A
					var a: self
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(25, 3, 13), P(28, 3, 16)), "type `self` can appear only in method throw types, method return types and method bodies"),
			},
		},
		"use self in a variable declaration in a module body": {
			input: `
				module A
					var a: self
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(26, 3, 13), P(29, 3, 16)), "type `self` can appear only in method throw types, method return types and method bodies"),
			},
		},
		"use self in a variable declaration in a mixin body": {
			input: `
				mixin A
					var a: self
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(25, 3, 13), P(28, 3, 16)), "type `self` can appear only in method throw types, method return types and method bodies"),
			},
		},
		"use self in a variable declaration in a method": {
			input: `
				class A
					def foo
						var a: self
					end
				end
			`,
		},
		"use self in a method return type": {
			input: `
				class A
					def foo: self
						self
					end
				end
			`,
		},
		"use self in a method throw type": {
			input: `
				class A
					def foo! self; end
				end
			`,
		},
		"use self in a method param type": {
			input: `
				class A
					def foo(a: self); end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(29, 3, 17), P(32, 3, 20)), "type `self` can appear only in method throw types, method return types and method bodies"),
			},
		},
		"assign self to self": {
			input: `
				class A
					def foo
						var a: self = self
					end
				end
			`,
		},
		"assign self to class": {
			input: `
				class A
					def foo
						var a: A = self
					end
				end
			`,
		},
		"assign class instance to self": {
			input: `
				class A
					def foo
						var b: self = A()
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(46, 4, 21), P(48, 4, 23)), "type `A` cannot be assigned to type `self`"),
			},
		},
		"do not replace self in method calls on self": {
			input: `
				class Foo
					def foo
						var a: &self = self.class
					end
				end
			`,
		},
		"assign &self to singleton class of parent": {
			input: `
				class Foo
					def foo
						var a: &Object = self.class
					end
				end
			`,
		},
		"assign &self to singleton class of child": {
			input: `
				class Bar < Foo; end
				class Foo
					def foo
						var a: &Bar = self.class
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(73, 5, 21), P(82, 5, 30)), "type `&self` cannot be assigned to type `&Bar`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
