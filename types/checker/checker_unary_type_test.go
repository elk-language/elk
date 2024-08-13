package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestSingletonType(t *testing.T) {
	tests := testTable{
		"singleton of class": {
			input: `
				var a: &String = String
			`,
		},
		"singleton of mixin": {
			input: `
				mixin Foo; end
				var a: &Foo = Foo
			`,
		},
		"singleton of interface": {
			input: `
				interface Foo; end
				var a: &Foo = Foo
			`,
		},
		"singleton of module": {
			input: `
				module Foo; end
				var a: &Foo = Foo
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(32, 3, 12), P(35, 3, 15)), "cannot get singleton class of `Foo`"),
			},
		},
		"singleton of literal": {
			input: `
				var a: &1
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 2, 12), P(13, 2, 13)), "cannot get singleton class of `1`"),
			},
		},

		"singleton of class self": {
			input: `
				class Foo
					def foo
						var a: &self
					end
				end
			`,
		},
		"singleton of mixin self": {
			input: `
				mixin Foo
					def foo
						var a: &self
					end
				end
			`,
		},
		"singleton of module self": {
			input: `
				module Foo
					def foo
						var a: &self
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(42, 4, 14), P(46, 4, 18)), "type `Foo` must be a class or mixin to be used with the singleton type"),
			},
		},

		"singleton of type parameter without bounds": {
			input: `
				class Foo[V]
					def foo
						var a: &V
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 4, 14), P(45, 4, 15)), "type parameter `V` must have an upper bound that is a class, mixin or interface to be used with the singleton type"),
			},
		},
		"singleton of type parameter with literal upper bound": {
			input: `
				class Foo[V < 1]
					def foo
						var a: &V
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 4, 14), P(49, 4, 15)), "type parameter `V` must have an upper bound that is a class, mixin or interface to be used with the singleton type"),
			},
		},
		"singleton of type parameter with union upper bound": {
			input: `
				class Foo[V < String | Int]
					def foo
						var a: &V
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(59, 4, 14), P(60, 4, 15)), "type parameter `V` must have an upper bound that is a class, mixin or interface to be used with the singleton type"),
			},
		},
		"singleton of type parameter with intersection upper bound": {
			input: `
				interface Bar
					def bar; end
				end

				class Foo[V < String & Bar]
					def foo
						var a: &V
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(104, 8, 14), P(105, 8, 15)), "type parameter `V` must have an upper bound that is a class, mixin or interface to be used with the singleton type"),
			},
		},
		"singleton of type parameter with class upper bound": {
			input: `
				class Foo[V < String]
					def foo
						var a: &V
					end
				end
			`,
		},
		"singleton of type parameter with mixin upper bound": {
			input: `
				mixin Bar; end

				class Foo[V < Bar]
					def foo
						var a: &V
					end
				end
			`,
		},
		"singleton of type parameter with interface upper bound": {
			input: `
				interface Bar; end

				class Foo[V < Bar]
					def foo
						var a: &V
					end
				end
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestInstanceOfType(t *testing.T) {
	tests := testTable{
		"instance of singleton class": {
			input: `
				var a: ^(&String) = "foo"
			`,
		},
		"singleton of class": {
			input: `
				class Foo; end
				var a: ^Foo = Foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(31, 3, 12), P(34, 3, 15)), "cannot get instance of `Foo`"),
			},
		},
		"instance of literal": {
			input: `
				var a: ^1
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 2, 12), P(13, 2, 13)), "cannot get instance of `1`"),
			},
		},

		"instance of class self": {
			input: `
				class Foo
					def foo
						var a: ^self
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(41, 4, 14), P(45, 4, 18)), "type `Foo` must be a singleton class to be used with the instance of type"),
			},
		},
		"instance of singleton self": {
			input: `
				class Foo
					singleton
						def foo
							var a: ^self
						end
					end
				end
			`,
		},

		"instance of type parameter without bounds": {
			input: `
				class Foo[V]
					def foo
						var a: ^V
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(44, 4, 14), P(45, 4, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"instance of type parameter with literal upper bound": {
			input: `
				class Foo[V < 1]
					def foo
						var a: ^V
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 4, 14), P(49, 4, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"instance of type parameter with union upper bound": {
			input: `
				class Foo[V < String | Int]
					def foo
						var a: ^V
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(59, 4, 14), P(60, 4, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"instance of type parameter with intersection upper bound": {
			input: `
				interface Bar
					def bar; end
				end

				class Foo[V < &String & Bar]
					def foo
						var a: ^V
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(105, 8, 14), P(106, 8, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"singleton of type parameter with class upper bound": {
			input: `
				class Foo[V < String]
					def foo
						var a: ^V
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(53, 4, 14), P(54, 4, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"singleton of type parameter with singleton upper bound": {
			input: `
				class Foo[V < &String]
					def foo
						var a: ^V
					end
				end
			`,
		},
		"singleton of type parameter with Class upper bound": {
			input: `
				class Foo[V < Class]
					def foo
						var a: ^V
					end
				end
			`,
		},
		"singleton of type parameter with Mixin upper bound": {
			input: `
				class Foo[V < Mixin]
					def foo
						var a: ^V
					end
				end
			`,
		},
		"singleton of type parameter with Interface upper bound": {
			input: `
				class Foo[V < Interface]
					def foo
						var a: ^V
					end
				end
			`,
		},
		"singleton of type parameter with Module upper bound": {
			input: `
				class Foo[V < Module]
					def foo
						var a: ^V
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(53, 4, 14), P(54, 4, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"return the instance of self in singleton context": {
			input: `
				class Foo
					singleton
						def create: ^self
							new
						end
					end
				end

				var a: Foo = Foo.create
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
