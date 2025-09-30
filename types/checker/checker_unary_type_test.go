package checker_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
)

func TestSingletonType(t *testing.T) {
	tests := testTable{
		"assign a singleton type of a class to Class": {
			input: `
				var a: Class = String
			`,
		},
		"assign a singleton type of a class to Mixin": {
			input: `
				var a: Mixin = String
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(20, 2, 20), P(25, 2, 25)), "type `&Std::String` cannot be assigned to type `Std::Mixin`"),
			},
		},
		"assign a singleton type of a mixin to Mixin": {
			input: `
				mixin Foo; end
				var a: Mixin = Foo
			`,
		},
		"assign a singleton type of an interface to Interface": {
			input: `
				interface Foo; end
				var a: Interface = Foo
			`,
		},
		"assign a module to Module": {
			input: `
				module Foo; end
				var a: Module = Foo
			`,
		},
		"assign class of child to singleton type of parent": {
			input: `
				class Foo; end
				class Bar < Foo; end
				var a: &Foo = Bar
			`,
		},

		"assign a class to a singleton type of the class": {
			input: `
				var a: &String = String
			`,
		},
		"assign a mixin to a singleton type of the mixin": {
			input: `
				mixin Foo; end
				var a: &Foo = Foo
			`,
		},
		"assign an interface to a singleton type of the interface": {
			input: `
				interface Foo; end
				var a: &Foo = Foo
			`,
		},
		"assign a module to a singleton type of the module": {
			input: `
				module Foo; end
				var a: &Foo = Foo
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(32, 3, 12), P(35, 3, 15)), "cannot get singleton class of `Foo`"),
			},
		},
		"singleton type of a literal": {
			input: `
				var a: &1
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(12, 2, 12), P(13, 2, 13)), "cannot get singleton class of `1`"),
			},
		},

		"can get the singleton of self in a class": {
			input: `
				class Foo
					def foo
						var a: &self
					end
				end
			`,
		},
		"can get the singleton of self in a mixin": {
			input: `
				mixin Foo
					def foo
						var a: &self
					end
				end
			`,
		},
		"can get the singleton of self in a module": {
			input: `
				module Foo
					def foo
						var a: &self
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(42, 4, 14), P(46, 4, 18)), "type `Foo` must be a class or mixin to be used with the singleton type"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(44, 4, 14), P(45, 4, 15)), "type parameter `V` must have an upper bound that is a class, mixin or interface to be used with the singleton type"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(48, 4, 14), P(49, 4, 15)), "type parameter `V` must have an upper bound that is a class, mixin or interface to be used with the singleton type"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(59, 4, 14), P(60, 4, 15)), "type parameter `V` must have an upper bound that is a class, mixin or interface to be used with the singleton type"),
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
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(104, 8, 14), P(105, 8, 15)), "type parameter `V` must have an upper bound that is a class, mixin or interface to be used with the singleton type"),
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
				var a: %(&String) = "foo"
			`,
		},
		"singleton of class": {
			input: `
				class Foo; end
				var a: %Foo = Foo()
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(31, 3, 12), P(34, 3, 15)), "cannot get instance of `Foo`"),
			},
		},
		"instance of literal": {
			input: `
				var a: %1
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(12, 2, 12), P(13, 2, 13)), "cannot get instance of `1`"),
			},
		},

		"instance of class self": {
			input: `
				class Foo
					def foo
						var a: %self
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(41, 4, 14), P(45, 4, 18)), "type `Foo` must be a singleton class to be used with the instance of type"),
			},
		},
		"instance of singleton self": {
			input: `
				class Foo
					singleton
						def foo
							var a: %self
						end
					end
				end
			`,
		},

		"instance of type parameter without bounds": {
			input: `
				class Foo[V]
					def foo
						var a: %V
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(44, 4, 14), P(45, 4, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"instance of type parameter with literal upper bound": {
			input: `
				class Foo[V < 1]
					def foo
						var a: %V
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(48, 4, 14), P(49, 4, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"instance of type parameter with union upper bound": {
			input: `
				class Foo[V < String | Int]
					def foo
						var a: %V
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(59, 4, 14), P(60, 4, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"instance of type parameter with intersection upper bound": {
			input: `
				interface Bar
					def bar; end
				end

				class Foo[V < &String & Bar]
					def foo
						var a: %V
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(105, 8, 14), P(106, 8, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"singleton of type parameter with class upper bound": {
			input: `
				class Foo[V < String]
					def foo
						var a: %V
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(53, 4, 14), P(54, 4, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"singleton of type parameter with singleton upper bound": {
			input: `
				class Foo[V < &String]
					def foo
						var a: %V
					end
				end
			`,
		},
		"singleton of type parameter with Class upper bound": {
			input: `
				class Foo[V < Class]
					def foo
						var a: %V
					end
				end
			`,
		},
		"singleton of type parameter with Mixin upper bound": {
			input: `
				class Foo[V < Mixin]
					def foo
						var a: %V
					end
				end
			`,
		},
		"singleton of type parameter with Interface upper bound": {
			input: `
				class Foo[V < Interface]
					def foo
						var a: %V
					end
				end
			`,
		},
		"singleton of type parameter with Module upper bound": {
			input: `
				class Foo[V < Module]
					def foo
						var a: %V
					end
				end
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(53, 4, 14), P(54, 4, 15)), "type parameter `V` must have an upper bound that is a singleton class to be used with the instance of type"),
			},
		},
		"return the instance of self in singleton context": {
			input: `
				class Foo
					singleton
						def create: %self
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

func TestUnaryMinusType(t *testing.T) {
	tests := testTable{
		"assign positive int to a negative int type": {
			input: `
				var a: -1 = 1
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(17, 2, 17), P(17, 2, 17)), "type `1` cannot be assigned to type `-1`"),
			},
		},
		"assign negative int to a positive int type": {
			input: `
				var a: 1 = -1
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(16, 2, 16), P(17, 2, 17)), "type `-1` cannot be assigned to type `1`"),
			},
		},
		"assign negative int to a negative int type": {
			input: `
				var a: -1 = -1
			`,
		},
		"apply to an invalid type": {
			input: `
				var a: -"c"
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(12, 2, 12), P(15, 2, 15)), "unary operator `-` cannot be used on type `\"c\"`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestBoxType(t *testing.T) {
	tests := testTable{
		"assign to box type": {
			input: `
				var a: Box[Int] = Box(1)
				var b: ^Int = a
			`,
		},
		"assign to immutable box type": {
			input: `
				var a: ImmutableBox[Int] = ImmutableBox(1)
				var b: *Int = a
			`,
		},
		"box type is distinct from its value": {
			input: `
				var a: ^Int = Box(1)
				var b: Int = a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(43, 3, 18), P(43, 3, 18)), "type `Std::Box[Std::Int]` cannot be assigned to type `Std::Int`"),
			},
		},
		"immutable box type is distinct from its value": {
			input: `
				var a: *Int = ImmutableBox(1)
				var b: Int = a
			`,
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L("<main>", P(52, 3, 18), P(52, 3, 18)), "type `Std::ImmutableBox[Std::Int]` cannot be assigned to type `Std::Int`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
