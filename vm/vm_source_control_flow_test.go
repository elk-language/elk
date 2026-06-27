package vm_test

import (
	"testing"

	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/value"
)

func TestVMSource_ThrowCatch(t *testing.T) {
	tests := sourceTestTable{
		"throw and do not catch": {
			source: `
				println "1"
				a := do
					println "2"
					throw :foo
					println "3"
					1
				catch :bar
					println "4"
					2
				catch :baz
					println "5"
					3
				end
				println "6"
				a
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(51, 5, 6), P(60, 5, 15)), "thrown value of type `:foo` must be caught"),
				diagnostic.NewWarning(L(P(67, 6, 6), P(77, 6, 16)), "unreachable code"),
			},
		},

		"throw in method and catch in parent context": {
			source: `
				def foo! :foo
					println "1"
					throw :foo
					println "2"
					1
				end

				a := do
					foo()
					println "3"
					2
				catch :foo
					println "4"
					3
				end
				println "5"
				a
			`,
			wantStdout:   "1\n4\n5\n",
			wantStackTop: value.SmallInt(3).ToValue(),
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(57, 5, 6), P(67, 5, 16)), "unreachable code"),
			},
		},
		"throw in nested method and catch in parent context": {
			source: `
				def foo! :foo
					println "1"
					throw :foo
					println "2"
					1
				end

				def bar! :foo
					do
						foo()
						println "3"
						2
					catch :bar
						println "4"
						3
					end
				end

				a := do
					bar()
					println "5"
					4
				catch :foo
					println "6"
					5
				end
				println "7"
				a
			`,
			wantStdout:   "1\n6\n7\n",
			wantStackTop: value.SmallInt(5).ToValue(),
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(57, 5, 6), P(67, 5, 16)), "unreachable code"),
			},
		},

		"execute finally before return": {
			source: `
				def bar: Int
					println("3")
					1
				end
				def foo
					println "1"
					do
						println "2"
						return bar()
					finally
						println "4"
						2
					end
					println "5"
					3
				end
				foo()
			`,
			wantStdout:   "1\n2\n3\n4\n",
			wantStackTop: value.SmallInt(1).ToValue(),
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(119, 10, 14), P(123, 10, 18)), "values returned in void context will be ignored"),
				diagnostic.NewWarning(L(P(178, 15, 6), P(188, 15, 16)), "unreachable code"),
			},
		},
		"execute nested finally before return": {
			source: `
				def bar: Int
					println("3")
					1
				end
				def foo
					println "1"
					do
						do
							println "2"
							return bar()
						finally
							println "4"
							2
						end
					finally
						println "5"
						3
					end
					println "6"
					4
				end
				foo()
			`,
			wantStdout:   "1\n2\n3\n4\n5\n",
			wantStackTop: value.SmallInt(1).ToValue(),
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(130, 11, 15), P(134, 11, 19)), "values returned in void context will be ignored"),
				diagnostic.NewWarning(L(P(241, 20, 6), P(251, 20, 16)), "unreachable code"),
			},
		},
		"execute finally before return in a setter method": {
			source: `
				def bar: Int
					println("3")
					1
				end
				def foo=(value: Int)
					println "1"
					do
						println "2"
						return bar()
					finally
						println "4"
						2
					end
					println "5"
					3
				end
				Kernel.foo = 25
			`,
			wantStdout:   "1\n2\n3\n4\n",
			wantStackTop: value.SmallInt(25).ToValue(),
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(132, 10, 14), P(136, 10, 18)), "values returned in void context will be ignored"),
				diagnostic.NewWarning(L(P(191, 15, 6), P(201, 15, 16)), "unreachable code"),
			},
		},
		"execute finally before return in a module": {
			source: `
				def foo: Int
					println("3")
					1
				end

				module D
					println "1"
					do
						println "2"
						return foo()
					finally
						println "4"
						2
					end
					println "5"
					3
				end
				nil
			`,
			wantStdout:   "1\n2\n3\n4\n",
			wantStackTop: value.Nil,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewWarning(L(P(121, 11, 14), P(125, 11, 18)), "values returned in void context will be ignored"),
				diagnostic.NewWarning(L(P(180, 16, 6), P(190, 16, 16)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_ForIn(t *testing.T) {
	tests := sourceTestTable{
		"loop over a non-iterable": {
			source: `
				for i in ::Std::Object()
					print(i.inspect, " ")
				end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(14, 2, 14), P(28, 2, 28)), "type `Std::Object` does not implement interface `Std::PrimitiveIterable[any, any]`:\n\n  - missing method `Std::PrimitiveIterable.:iter` with signature: `def iter(): Std::Iterator[any, any]`"),
				diagnostic.NewFailure(L(P(14, 2, 14), P(28, 2, 28)), "type `Std::Object` cannot be iterated over, it does not implement `Std::PrimitiveIterable[any, any]`"),
			},
		},
		"loop over an invalid iterable": {
			source: `
				class InvalidIterator
					def iter then self
				end

				for i in ::InvalidIterator()
					print(i.inspect, " ")
				end
			`,
			wantCompileErr: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(P(73, 6, 14), P(91, 6, 32)), "type `InvalidIterator` does not implement interface `Std::PrimitiveIterable[any, any]`:\n\n  - incorrect implementation of `Std::PrimitiveIterable.:iter`\n      is:        `def iter(): void`\n      should be: `def iter(): Std::Iterator[any, any]`"),
				diagnostic.NewFailure(L(P(73, 6, 14), P(91, 6, 32)), "type `InvalidIterator` cannot be iterated over, it does not implement `Std::PrimitiveIterable[any, any]`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
