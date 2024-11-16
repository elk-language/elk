package vm_test

import (
	"testing"

	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/value"
)

func TestVMSource_Must(t *testing.T) {
	tests := sourceTestTable{
		"must with value": {
			source: `
				println "1"
				a := must 5
				println a
			`,
			wantStdout:   "1\n5\n",
			wantStackTop: value.NilType{},
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(26, 3, 10), P(31, 3, 15)), "unnecessary `must`, type `5` is not nilable"),
			},
		},
		"must with nil": {
			source: `
				println "1"
				var a: Int? = nil
				b := must a
				println b
			`,
			wantStdout:     "1\n",
			wantRuntimeErr: value.NewUnexpectedNilError(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_As(t *testing.T) {
	tests := sourceTestTable{
		"valid downcast": {
			source: `
				println "1"
				var a: Int | Float = 5
				b := a as ::Std::Int
				println b
			`,
			wantStdout:   "1\n5\n",
			wantStackTop: value.NilType{},
		},
		"valid upcast": {
			source: `
				println "1"
				a := 5
				b := a as ::Std::Value
				println b.inspect
			`,
			wantStdout:   "1\n5\n",
			wantStackTop: value.NilType{},
		},
		"invalid cast": {
			source: `
				println "1"
				var a: Int | Float = 5
				b := a as ::Std::Float
				println b.inspect
			`,
			wantStdout:     "1\n",
			wantRuntimeErr: value.NewError(value.TypeErrorClass, "failed type cast, `5` is not an instance of `Std::Float`"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_ThrowCatch(t *testing.T) {
	tests := sourceTestTable{
		"throw": {
			source: `
				println "1"
				throw unchecked :foo
				println "2"
			`,
			wantStdout:     "1\n",
			wantRuntimeErr: value.ToSymbol("foo"),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(46, 4, 5), P(56, 4, 15)), "unreachable code"),
			},
		},
		"throw and catch": {
			source: `
				println "1"
				a := do
					println "2"
					throw :foo
					println "3"
					1
				catch :foo
					println "4"
					2
				end
				println "5"
				a
			`,
			wantStdout:   "1\n2\n4\n5\n",
			wantStackTop: value.SmallInt(2),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(67, 6, 6), P(77, 6, 16)), "unreachable code"),
			},
		},
		"throw and catch in second branch": {
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
				catch :foo
					println "5"
					3
				end
				println "6"
				a
			`,
			wantStdout:   "1\n2\n5\n6\n",
			wantStackTop: value.SmallInt(3),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(67, 6, 6), P(77, 6, 16)), "unreachable code"),
			},
		},
		"throw and catch with pattern matching": {
			source: `
				a := do
					println "1"
					throw "foo"
					println "2"
					1
				catch "bar" || "baz"
					println "3"
					2
				catch ::Std::String(length: > 2) as str
					println "4, str: ${str}"
					3
				end
				println "5"
				a
			`,
			wantStdout:   "1\n4, str: foo\n5\n",
			wantStackTop: value.SmallInt(3),
		},
		"throw unchecked and do not catch": {
			source: `
				println "1"
				a := do
					println "2"
					throw unchecked :foo
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
			wantStdout:     "1\n2\n",
			wantStackTop:   value.SmallInt(3),
			wantRuntimeErr: value.ToSymbol("foo"),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(77, 6, 6), P(87, 6, 16)), "unreachable code"),
			},
		},
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
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(51, 5, 6), P(60, 5, 15)), "thrown value of type `:foo` must be caught"),
				error.NewWarning(L(P(67, 6, 6), P(77, 6, 16)), "unreachable code"),
			},
		},
		"throw and catch in parent": {
			source: `
				a := do
					println "1"
					do
						println "2"
						throw :foo
						println "3"
					catch :bar
						println "4"
						1
					end
					println "5"
					2
				catch :foo
					println "6"
					3
				end
				println "7"
				a
			`,
			wantStdout:   "1\n2\n6\n7\n",
			wantStackTop: value.SmallInt(3),
		},
		"throw in catch and catch in parent": {
			source: `
				a := do
					do
						println "1"
						throw :foo
						println "2"
					catch :foo
						do
							println "3"
							throw :bar
							println "4"
						catch :biz
							println "5"
						end
						println "6"
					end
					1
				catch :bar
					println "7"
					2
				end
				println "8"
				a
			`,
			wantStdout:   "1\n3\n7\n8\n",
			wantStackTop: value.SmallInt(2),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(143, 11, 8), P(153, 11, 18)), "unreachable code"),
				error.NewWarning(L(P(62, 6, 7), P(72, 6, 17)), "unreachable code"),
			},
		},
		"finally without throw": {
			source: `
				a := do
					println "1"
					1
				finally
					println "2"
					2
				end
				println "3"
				a
			`,
			wantStdout:   "1\n2\n3\n",
			wantStackTop: value.SmallInt(1),
		},
		"throw, catch and execute finally": {
			source: `
				a := do
					println "1"
					throw :foo
					println "2"
					2
				catch :foo
					println "3"
					3
				finally
					println "4"
					4
				end
				println "5"
				a
			`,
			wantStdout:   "1\n3\n4\n5\n",
			wantStackTop: value.SmallInt(3),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(51, 5, 6), P(61, 5, 16)), "unreachable code"),
			},
		},
		"throw, execute finally and rethrow": {
			source: `
				a := do
					println "1"
					throw :foo
					println "2"
					1
				finally
					println "3"
					2
				end
				println "4"
				a
			`,
			wantStdout:     "1\n3\n",
			wantRuntimeErr: value.ToSymbol("foo"),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(79, 7, 7), P(89, 7, 17)), "unreachable code"),
			},
		},
		"throw, execute finally, throw and catch in finally, rethrow": {
			source: `
				a := do
					println "1"
					throw unchecked :foo
					println "2"
					1
				finally
					do
						println "3"
						throw :bar
						println "4"
						2
					catch :bar
						println "5"
						3
					end
				end
				println "6"
				a
			`,
			wantStdout:     "1\n3\n5\n",
			wantRuntimeErr: value.ToSymbol("foo"),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(141, 11, 7), P(151, 11, 17)), "unreachable code"),
				error.NewWarning(L(P(61, 5, 6), P(71, 5, 16)), "unreachable code"),
				error.NewWarning(L(P(224, 18, 5), P(234, 18, 15)), "unreachable code"),
			},
		},
		"execute finally, throw and catch in finally": {
			source: `
				a := do
					println "1"
					1
				finally
					do
						println "2"
						throw :bar
						println "3"
						2
					catch :bar
						println "4"
						3
					end
				end
				println "5"
				a
			`,
			wantStdout:   "1\n2\n4\n5\n",
			wantStackTop: value.SmallInt(1),
		},
		"throw, execute finally, throw and catch in parent": {
			source: `
				a := do
					println "1"
					do
						println "2"
						throw unchecked :foo
						println "3"
						1
					finally
						println "4"
						throw :bar
						println "5"
						2
					end
				catch :bar
					println "6"
					3
				end
				println "7"
				a
			`,
			wantStdout:   "1\n2\n4\n6\n7\n",
			wantStackTop: value.SmallInt(3),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(163, 12, 7), P(173, 12, 17)), "unreachable code"),
				error.NewWarning(L(P(89, 7, 7), P(99, 7, 17)), "unreachable code"),
			},
		},
		"execute finally, throw and catch in parent": {
			source: `
				a := do
					println "1"
					do
						println "2"
						1
					finally
						println "3"
						throw :bar
						println "4"
						2
					end
				catch :bar
					println "5"
					3
				end
				println "6"
				a
			`,
			wantStdout:   "1\n2\n3\n5\n6\n",
			wantStackTop: value.SmallInt(3),
		},
		"throw, execute finally, rethrow and catch in parent": {
			source: `
				a := do
					println "1"
					do
						println "2"
						throw :foo
						println "3"
						1
					finally
						println "4"
						2
					end
				catch :foo
					println "5"
					3
				end
				println "6"
				a
			`,
			wantStdout:   "1\n2\n4\n5\n6\n",
			wantStackTop: value.SmallInt(3),
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
			wantStackTop: value.SmallInt(3),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(57, 5, 6), P(67, 5, 16)), "unreachable code"),
			},
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
			},
		},
		"throw in nested method and catch in parent context": {
			source: `
				def foo
					println "1"
					throw :foo
					println "2"
					1
				end

				def bar
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
			wantStackTop: value.SmallInt(5),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("bar"))
			},
		},
		"throw, catch and break in loop": {
			source: `
			  i := 0
				loop
					println i
					do
						throw :foo if i == 5
					catch :foo
						break i
					end
					i++
				end
			`,
			wantStdout:   "0\n1\n2\n3\n4\n5\n",
			wantStackTop: value.SmallInt(5),
		},
		"execute finally before return": {
			source: `
				def foo
					println "1"
					do
						println "2"
						return println("3") ?? 1
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
			wantStackTop: value.SmallInt(1),
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
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
			wantStackTop: value.SmallInt(1),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(130, 11, 15), P(134, 11, 19)), "unreachable code"),
				error.NewWarning(L(P(241, 20, 6), P(251, 20, 16)), "unreachable code"),
			},
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
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
				self.foo = 25
			`,
			wantStdout:   "1\n2\n3\n4\n",
			wantStackTop: value.SmallInt(25),
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(132, 10, 14), P(136, 10, 18)), " values returned in void context will be ignored"),
				error.NewWarning(L(P(191, 15, 16), P(201, 15, 16)), "unreachable code"),
			},
			teardown: func() {
				delete(value.GlobalObjectSingletonClass.Methods, value.ToSymbol("foo"))
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
			wantCompileErr: error.ErrorList{
				error.NewWarning(L(P(121, 11, 14), P(125, 11, 18)), "values returned in void context will be ignored"),
				error.NewWarning(L(P(180, 16, 6), P(190, 16, 16)), "unreachable code"),
			},
		},
		"execute finally before break": {
			source: `
				def foo: Int
					println("2")
					1
				end
				a := loop
					do
						println "1"
						break foo()
						println "3"
						2
					finally
						println "4"
						3
					end
					println "5"
					4
				end
				println "6"
				a
			`,
			wantStdout:   "1\n2\n4\n6\n",
			wantStackTop: value.SmallInt(1),
		},
		"execute nested finally before break": {
			source: `
				def foo: Int
					println("2")
					1
				end
				var a
				do
					a = loop
						do
							do
								println "1"
								break foo()
								println "3"
								2
							finally
								println "4"
								3
							end
						finally
							println "5"
							4
						end
						println "6"
						5
					end
					println "7"
					6
				finally
					println "8"
					7
				end
				println "9"
				a
			`,
			wantStdout:   "1\n2\n4\n5\n7\n8\n9\n",
			wantStackTop: value.SmallInt(1),
		},
		"execute finally before continue": {
			source: `
				def foo: Int
					println("2")
					1
				end
				a := (do
					do
						println "1"
						continue foo()
						println "3"
						2
					finally
						println "4"
						3
					end
					println "5"
					4
				end while false)
				println "6"
				a
			`,
			wantStdout:   "1\n2\n4\n6\n",
			wantStackTop: value.SmallInt(1),
		},
		"execute nested finally before continue": {
			source: `
				def foo: Int
					println("2")
					1
				end
				var a
				do
					a = (do
						do
							do
								println "1"
								continue foo()
								println "3"
								2
							finally
								println "4"
								3
							end
						finally
							println "5"
							4
						end
						println "6"
						5
					end while false)
					println "7"
					6
				finally
					println "8"
					7
				end
				println "9"
				a
			`,
			wantStdout:   "1\n2\n4\n5\n7\n8\n9\n",
			wantStackTop: value.SmallInt(1),
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
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `iter` is not available to value of class `Std::Object`: Std::Object{}",
			),
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
			wantRuntimeErr: value.NewError(
				value.NoMethodErrorClass,
				"method `next` is not available to value of class `InvalidIterator`: InvalidIterator{}",
			),
		},
		"loop over a list": {
			source: `
				for i in [1, 2, 3, :foo, 'bar']
					print(i.inspect, " ")
				end
			`,
			wantStackTop: value.Nil,
			wantStdout:   `1 2 3 :foo "bar" `,
		},
		"loop over a string": {
			source: `
				for i in "Poznań jest √🔥"
					print(i.inspect, ", ")
				end
			`,
			wantStackTop: value.Nil,
			wantStdout:   "`P`, `o`, `z`, `n`, `a`, `ń`, ` `, `j`, `e`, `s`, `t`, ` `, `√`, `🔥`, ",
		},
		"loop over a hashmap with a pattern": {
			source: `
				h := {
					"foo" => 21,
					"elo" => 54,
					"grim" => -8,
				}
				sum := 0
				for ::Std::Pair(key: _, value) in h
					sum += value
				end
				println(sum.inspect)
			`,
			wantStackTop: value.Nil,
			wantStdout:   "67\n",
		},
		"loop over a string byte iterator": {
			source: `
				for i in "Poznań jest √🔥".byte_iter
					print(i.inspect, ", ")
				end
			`,
			wantStackTop: value.Nil,
			wantStdout:   "80u8, 111u8, 122u8, 110u8, 97u8, 197u8, 132u8, 32u8, 106u8, 101u8, 115u8, 116u8, 32u8, 226u8, 136u8, 154u8, 240u8, 159u8, 148u8, 165u8, ",
		},
		"loop over a arrayTuple": {
			source: `
				for i in %[1, 2, 3, :foo, 'bar']
					print(i.inspect, " ")
				end
			`,
			wantStackTop: value.Nil,
			wantStdout:   `1 2 3 :foo "bar" `,
		},
		"with break": {
			source: `
				for i in [1, 2, 3, 4, 5]
					break if i > 3
					print(i.inspect, " ")
				end
			`,
			wantStackTop: value.Nil,
			wantStdout:   `1 2 3 `,
		},
		"with break with value": {
			source: `
				for i in [1, 2, 3, 4, 5]
					break i if i > 3
					print(i.inspect, " ")
				end
			`,
			wantStackTop: value.SmallInt(4),
			wantStdout:   `1 2 3 `,
		},
		"nested": {
			source: `
				for i in [1, 2, 3]
					for j in [8, 9, 10]
						print(i.inspect, ":", j.inspect, " ")
					end
				end
			`,
			wantStackTop: value.Nil,
			wantStdout:   `1:8 1:9 1:10 2:8 2:9 2:10 3:8 3:9 3:10 `,
		},
		"nested with break": {
			source: `
				for i in [1, 2, 3]
					for j in [8, 9, 10]
						break if j == 9
						print(i.inspect, ":", j.inspect, " ")
					end
				end
			`,
			wantStackTop: value.Nil,
			wantStdout:   `1:8 2:8 3:8 `,
		},
		"nested with labeled break": {
			source: `
				$outer: for i in [1, 2, 3]
					for j in [8, 9, 10]
						break$outer if j == 10
						print(i.inspect, ":", j.inspect, " ")
					end
				end
			`,
			wantStackTop: value.Nil,
			wantStdout:   `1:8 1:9 `,
		},
		"nested with labeled break with value": {
			source: `
				$outer: for i in [1, 2, 3]
					for j in [8, 9, 10]
						break$outer j if j == 10
						print(i.inspect, ":", j.inspect, " ")
					end
				end
			`,
			wantStackTop: value.SmallInt(10),
			wantStdout:   `1:8 1:9 `,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_NumericFor(t *testing.T) {
	tests := sourceTestTable{
		"calculate the sum of consecutive natural numbers": {
			source: `
				a := 0
				fornum i := 1; i < 6; i += 1
					a += i
				end
				a
			`,
			wantStackTop: value.SmallInt(15),
		},
		"create a repeated string": {
			source: `
				a := ""
				fornum i := 20; i > 0; i -= 2
					a += "-"
				end
				a
			`,
			wantStackTop: value.String("----------"),
		},
		"calculate the factorial of 10": {
			source: `
				a := 1
				fornum i := 2; i <= 10; i += 1
					a *= i
				end
				a
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"return the value of the last iteration": {
			source: `
				a := 1
				fornum i := 2; i <= 10; i += 1
					a *= i
				end
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"return nil when no iterations": {
			source: `
				a := 1
				fornum i := 20; i <= 10; i += 1
					a *= i
				end
			`,
			wantStackTop: value.Nil,
		},
		"return nil after break": {
			source: `
				a := 1
				fornum i := 2; i <= 10; i += 1
					a *= i
					break if a > 200
				end
			`,
			wantStackTop: value.Nil,
		},
		"return a value using break": {
			source: `
				a := 1
				fornum i := 2; i <= 10; i += 1
					a *= i
					break a if a > 200
				end
			`,
			wantStackTop: value.SmallInt(720),
		},
		"nested with continue": {
			source: `
				fornum j := 1; j <= 5; j += 1
					fornum i := 1; i <= 5; i += 1
						continue if i + j > 5
						println j.to_string + ":" + i.to_string
					end
				end
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n2:1\n2:2\n2:3\n3:1\n3:2\n4:1\n",
			wantStackTop: value.Nil,
		},
		"nested with a labeled continue": {
			source: `
				$foo: fornum j := 1; j <= 5; j += 1
					fornum i := 1; i <= 5; i += 1
						continue$foo if i % 2 == 0 || j % 2 == 0
						println j.to_string + ":" + i.to_string
					end
				end
			`,
			wantStdout:   "1:1\n3:1\n5:1\n",
			wantStackTop: value.Nil,
		},
		"nested with break": {
			source: `
				fornum j := 1;; j += 1
					fornum i := 1;; i += 1
						println j.to_string + ":" + i.to_string
						break if i >= 5
					end
					break if j >= 5
				end
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n1:5\n2:1\n2:2\n2:3\n2:4\n2:5\n3:1\n3:2\n3:3\n3:4\n3:5\n4:1\n4:2\n4:3\n4:4\n4:5\n5:1\n5:2\n5:3\n5:4\n5:5\n",
			wantStackTop: value.Nil,
		},
		"nested with a labeled break": {
			source: `
				$foo: fornum j := 1;; j += 1
					fornum i := 1;; i += 1
						println j.to_string + ":" + i.to_string
						break$foo if i >= 5
					end
					break if j >= 5
				end
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n1:5\n",
			wantStackTop: value.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_While(t *testing.T) {
	tests := sourceTestTable{
		"calculate the sum of consecutive natural numbers": {
			source: `
				a := 0
				i := 0
				while i < 6
					a += i
					i += 1
				end
				a
			`,
			wantStackTop: value.SmallInt(15),
		},
		"return nil with break": {
			source: `
				a := 0
				i := 0
				while true
					a += i
					i += 1
					break if i >= 6
				end
			`,
			wantStackTop: value.Nil,
		},
		"with break": {
			source: `
				a := 0
				i := 0
				while true
					a += i
					i += 1
					break if i >= 6
				end
				a
			`,
			wantStackTop: value.SmallInt(15),
		},
		"nested with break": {
			source: `
				j := 0
				while true
					j += 1
					i := 0
					while true
						break if i >= 5
						i += 1
						println j.to_string + ":" + i.to_string
					end
					break if j >= 5
				end
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n1:5\n2:1\n2:2\n2:3\n2:4\n2:5\n3:1\n3:2\n3:3\n3:4\n3:5\n4:1\n4:2\n4:3\n4:4\n4:5\n5:1\n5:2\n5:3\n5:4\n5:5\n",
			wantStackTop: value.Nil,
		},
		"nested with a labeled break": {
			source: `
				j := 0
				$foo: while true
					j += 1
					i := 0
					while true
						break$foo if i >= 5
						i += 1
						println j.to_string + ":" + i.to_string
					end
					break if j >= 5
				end
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n1:5\n",
			wantStackTop: value.Nil,
		},
		"continue": {
			source: `
				i := 0
				while i < 2
					i += 1
					println "before"
					continue println "during"
					println "after"
				end
			`,
			wantStdout:   "before\nduring\nbefore\nduring\n",
			wantStackTop: value.Nil,
		},
		"nested with continue": {
			source: `
				j := 0
				while j < 5
					j += 1
					i := 0
					while i < 5
						i += 1
						continue if i + j > 5
						println j.to_string + ":" + i.to_string
					end
				end
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n2:1\n2:2\n2:3\n3:1\n3:2\n4:1\n",
			wantStackTop: value.Nil,
		},
		"nested with a labeled continue": {
			source: `
				j := 0
				$foo: while j < 5
					j += 1
					i := 0
					while i < 5
						i += 1
						continue$foo if i % 2 == 0 || j % 2 == 0
						println j.to_string + ":" + i.to_string
					end
				end
			`,
			wantStdout:   "1:1\n3:1\n5:1\n",
			wantStackTop: value.Nil,
		},
		"return a value with break": {
			source: `
				a := 0
				i := 0
				while true
					a += i
					i += 1
					break a if i >= 6
				end
			`,
			wantStackTop: value.SmallInt(15),
		},
		"create a repeated string": {
			source: `
				a := ""
				i := 20
				while i > 0
				  a += "-"
					i -= 2
				end
				a
			`,
			wantStackTop: value.String("----------"),
		},
		"calculate the factorial of 10": {
			source: `
				a := 1
				i := 2
				while i <= 10
					a *= i
					i += 1
				end
				a
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"return the value of the last iteration": {
			source: `
				a := 1
				i := 2
				while i <= 10
					a *= i
					i += 1
					a
				end
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"return nil when no iterations": {
			source: `
				a := 1
				i := 20
				while i <= 10
				  a *= i
					i += 1
				end
			`,
			wantStackTop: value.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_ModifierWhile(t *testing.T) {
	tests := sourceTestTable{
		"calculate the sum of consecutive natural numbers": {
			source: `
				a := 0
				i := 0
				do
					a += i
					i += 1
				end while i < 6
				a
			`,
			wantStackTop: value.SmallInt(15),
		},
		"return nil with break": {
			source: `
				a := 0
				i := 0
				do
					a += i
					i += 1
					break if i >= 6
				end while true
			`,
			wantStackTop: value.Nil,
		},
		"with break": {
			source: `
				a := 0
				i := 0
				do
					a += i
					i += 1
					break if i >= 6
				end while true
				a
			`,
			wantStackTop: value.SmallInt(15),
		},
		"nested with break": {
			source: `
				j := 0
				do
					j += 1
					i := 0
					do
						break if i >= 5
						i += 1
						println j.to_string + ":" + i.to_string
					end while true
					break if j >= 5
				end while true
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n1:5\n2:1\n2:2\n2:3\n2:4\n2:5\n3:1\n3:2\n3:3\n3:4\n3:5\n4:1\n4:2\n4:3\n4:4\n4:5\n5:1\n5:2\n5:3\n5:4\n5:5\n",
			wantStackTop: value.Nil,
		},
		"nested with a labeled break": {
			source: `
				j := 0
				$foo: do
					j += 1
					i := 0
					do
						break$foo if i >= 5
						i += 1
						println j.to_string + ":" + i.to_string
					end while true
					break if j >= 5
				end while true
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n1:5\n",
			wantStackTop: value.Nil,
		},
		"continue": {
			source: `
				i := 0
				do
					i += 1
					println "before"
					continue println "during"
					println "after"
				end while i < 2
			`,
			wantStdout:   "before\nduring\nbefore\nduring\n",
			wantStackTop: value.Nil,
		},
		"nested with continue": {
			source: `
				j := 0
				do
					j += 1
					i := 0
					do
						i += 1
						continue if i + j > 5
						println j.to_string + ":" + i.to_string
					end while i < 5
				end while j < 5
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n2:1\n2:2\n2:3\n3:1\n3:2\n4:1\n",
			wantStackTop: value.Nil,
		},
		"nested with a labeled continue": {
			source: `
				j := 0
				$foo: do
					j += 1
					i := 0
					do
						i += 1
						continue$foo if i % 2 == 0 || j % 2 == 0
						println j.to_string + ":" + i.to_string
					end while i < 5
				end while j < 5
			`,
			wantStdout:   "1:1\n3:1\n5:1\n",
			wantStackTop: value.Nil,
		},
		"return a value with break": {
			source: `
				a := 0
				i := 0
				do
					a += i
					i += 1
					break a if i >= 6
				end while true
			`,
			wantStackTop: value.SmallInt(15),
		},
		"create a repeated string": {
			source: `
				a := ""
				i := 20
				do
				  a += "-"
					i -= 2
				end while i > 0
				a
			`,
			wantStackTop: value.String("----------"),
		},
		"calculate the factorial of 10": {
			source: `
				a := 1
				i := 2
				do
					a *= i
					i += 1
				end while i <= 10
				a
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"return the value of the last iteration": {
			source: `
				a := 1
				i := 2
				do
					a *= i
					i += 1
					a
				end while i <= 10
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"always does at least one iteration": {
			source: `
				a := 1
				i := 20
				do
				  a *= i
					i += 1
				end while i <= 10
			`,
			wantStackTop: value.SmallInt(21),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_Until(t *testing.T) {
	tests := sourceTestTable{
		"calculate the sum of consecutive natural numbers": {
			source: `
				a := 0
				i := 0
				until i >= 6
					a += i
					i += 1
				end
				a
			`,
			wantStackTop: value.SmallInt(15),
		},
		"return nil with break": {
			source: `
				a := 0
				i := 0
				until false
					a += i
					i += 1
					break if i >= 6
				end
			`,
			wantStackTop: value.Nil,
		},
		"with break": {
			source: `
				a := 0
				i := 0
				until false
					a += i
					i += 1
					break if i >= 6
				end
				a
			`,
			wantStackTop: value.SmallInt(15),
		},
		"nested with break": {
			source: `
				j := 0
				until false
					j += 1
					i := 0
					until false
						break if i >= 5
						i += 1
						println j.to_string + ":" + i.to_string
					end
					break if j >= 5
				end
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n1:5\n2:1\n2:2\n2:3\n2:4\n2:5\n3:1\n3:2\n3:3\n3:4\n3:5\n4:1\n4:2\n4:3\n4:4\n4:5\n5:1\n5:2\n5:3\n5:4\n5:5\n",
			wantStackTop: value.Nil,
		},
		"nested with a labeled break": {
			source: `
				j := 0
				$foo: until false
					j += 1
					i := 0
					until false
						break$foo if i >= 5
						i += 1
						println j.to_string + ":" + i.to_string
					end
					break if j >= 5
				end
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n1:5\n",
			wantStackTop: value.Nil,
		},
		"continue": {
			source: `
				i := 0
				until i >= 2
					i += 1
					println "before"
					continue println "during"
					println "after"
				end
			`,
			wantStdout:   "before\nduring\nbefore\nduring\n",
			wantStackTop: value.Nil,
		},
		"nested with continue": {
			source: `
				j := 0
				until j >= 5
					j += 1
					i := 0
					until i >= 5
						i += 1
						continue if i + j > 5
						println j.to_string + ":" + i.to_string
					end
				end
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n2:1\n2:2\n2:3\n3:1\n3:2\n4:1\n",
			wantStackTop: value.Nil,
		},
		"nested with a labeled continue": {
			source: `
				j := 0
				$foo: until j >= 5
					j += 1
					i := 0
					until i >= 5
						i += 1
						continue$foo if i % 2 == 0 || j % 2 == 0
						println j.to_string + ":" + i.to_string
					end
				end
			`,
			wantStdout:   "1:1\n3:1\n5:1\n",
			wantStackTop: value.Nil,
		},
		"return a value with break": {
			source: `
				a := 0
				i := 0
				until false
					a += i
					i += 1
					break a if i >= 6
				end
			`,
			wantStackTop: value.SmallInt(15),
		},
		"create a repeated string": {
			source: `
				a := ""
				i := 20
				until i <= 0
				  a += "-"
					i -= 2
				end
				a
			`,
			wantStackTop: value.String("----------"),
		},
		"calculate the factorial of 10": {
			source: `
				a := 1
				i := 2
				until i > 10
					a *= i
					i += 1
				end
				a
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"return the value of the last iteration": {
			source: `
				a := 1
				i := 2
				until i > 10
					a *= i
					i += 1
					a
				end
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"return nil when no iterations": {
			source: `
				a := 1
				i := 20
				until i > 10
				  a *= i
					i += 1
				end
			`,
			wantStackTop: value.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_ModifierUntil(t *testing.T) {
	tests := sourceTestTable{
		"calculate the sum of consecutive natural numbers": {
			source: `
				a := 0
				i := 0
				do
					a += i
					i += 1
				end until i >= 6
				a
			`,
			wantStackTop: value.SmallInt(15),
		},
		"return nil with break": {
			source: `
				a := 0
				i := 0
				do
					a += i
					i += 1
					break if i >= 6
				end until false
			`,
			wantStackTop: value.Nil,
		},
		"with break": {
			source: `
				a := 0
				i := 0
				do
					a += i
					i += 1
					break if i >= 6
				end until false
				a
			`,
			wantStackTop: value.SmallInt(15),
		},
		"nested with break": {
			source: `
				j := 0
				do
					j += 1
					i := 0
					do
						break if i >= 5
						i += 1
						println j.to_string + ":" + i.to_string
					end until false
					break if j >= 5
				end until false
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n1:5\n2:1\n2:2\n2:3\n2:4\n2:5\n3:1\n3:2\n3:3\n3:4\n3:5\n4:1\n4:2\n4:3\n4:4\n4:5\n5:1\n5:2\n5:3\n5:4\n5:5\n",
			wantStackTop: value.Nil,
		},
		"nested with a labeled break": {
			source: `
				j := 0
				$foo: do
					j += 1
					i := 0
					do
						break$foo if i >= 5
						i += 1
						println j.to_string + ":" + i.to_string
					end until false
					break if j >= 5
				end until false
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n1:5\n",
			wantStackTop: value.Nil,
		},
		"continue": {
			source: `
				i := 0
				do
					i += 1
					println "before"
					continue println "during"
					println "after"
				end until i >= 2
			`,
			wantStdout:   "before\nduring\nbefore\nduring\n",
			wantStackTop: value.Nil,
		},
		"nested with continue": {
			source: `
				j := 0
				do
					j += 1
					i := 0
					do
						i += 1
						continue if i + j > 5
						println j.to_string + ":" + i.to_string
					end until i >= 5
				end until j >= 5
			`,
			wantStdout:   "1:1\n1:2\n1:3\n1:4\n2:1\n2:2\n2:3\n3:1\n3:2\n4:1\n",
			wantStackTop: value.Nil,
		},
		"nested with a labeled continue": {
			source: `
				j := 0
				$foo: do
					j += 1
					i := 0
					do
						i += 1
						continue$foo if i % 2 == 0 || j % 2 == 0
						println j.to_string + ":" + i.to_string
					end until i >= 5
				end until j >= 5
			`,
			wantStdout:   "1:1\n3:1\n5:1\n",
			wantStackTop: value.Nil,
		},
		"return a value with break": {
			source: `
				a := 0
				i := 0
				do
					a += i
					i += 1
					break a if i >= 6
				end until false
			`,
			wantStackTop: value.SmallInt(15),
		},
		"create a repeated string": {
			source: `
				a := ""
				i := 20
				do
				  a += "-"
					i -= 2
				end until i <= 0
				a
			`,
			wantStackTop: value.String("----------"),
		},
		"calculate the factorial of 10": {
			source: `
				a := 1
				i := 2
				do
					a *= i
					i += 1
				end until i > 10
				a
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"return the value of the last iteration": {
			source: `
				a := 1
				i := 2
				do
					a *= i
					i += 1
					a
				end until i > 10
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"always does at least one iteration": {
			source: `
				a := 1
				i := 20
				do
				  a *= i
					i += 1
				end until i > 10
			`,
			wantStackTop: value.SmallInt(21),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_IfExpressions(t *testing.T) {
	tests := sourceTestTable{
		"return nil when condition is truthy and then is empty": {
			source:       "if true; end",
			wantStackTop: value.Nil,
		},
		"return nil when condition is falsy and then is empty": {
			source:       "if false; end",
			wantStackTop: value.Nil,
		},
		"execute the then branch": {
			source: `
				a := 5
				if a
					a = a + 2
				end
			`,
			wantStackTop: value.SmallInt(7),
		},
		"execute the empty else branch": {
			source: `
				a := 5
				if false
					a = a * 2
				end
			`,
			wantStackTop: value.Nil,
		},
		"execute the then branch instead of else": {
			source: `
				a := 5
				if a
					a = a + 2
				else
					a = 30
				end
			`,
			wantStackTop: value.SmallInt(7),
		},
		"execute the else branch instead of then": {
			source: `
				a := 5
				if nil
					a = a + 2
				else
					a = 30
				end
			`,
			wantStackTop: value.SmallInt(30),
		},
		"is an expression": {
			source: `
				a := 5
				b := if a
					"foo"
				else
					5
				end
				b
			`,
			wantStackTop: value.String("foo"),
		},
		"modifier binds more strongly than assignment": {
			source: `
				a := 5
				b := "foo" if a else 5
				b
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(43, 4, 5), P(43, 4, 5)), "undeclared variable: b"),
			},
		},
		"modifier returns the left side if the condition is satisfied": {
			source: `
				a := 5
				"foo" if a else 5
			`,
			wantStackTop: value.String("foo"),
		},
		"modifier returns the right side if the condition is not satisfied": {
			source: `
				a := nil
				"foo" if a else 5
			`,
			wantStackTop: value.SmallInt(5),
		},
		"modifier returns nil when condition is not satisfied": {
			source: `
				a := nil
				"foo" if a
			`,
			wantStackTop: value.Nil,
		},
		"can access variables defined in the condition": {
			source: `
				a + " bar" if a := "foo"
			`,
			wantStackTop: value.String("foo bar"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_UnlessExpressions(t *testing.T) {
	tests := sourceTestTable{
		"return nil when condition is falsy and then is empty": {
			source:       "unless false; end",
			wantStackTop: value.Nil,
		},
		"return nil when condition is truthy and then is empty": {
			source:       "unless true; end",
			wantStackTop: value.Nil,
		},
		"execute the then branch": {
			source: `
				a := nil
				unless a
					a = 7
				end
			`,
			wantStackTop: value.SmallInt(7),
		},
		"execute the empty else branch": {
			source: `
				a := 5
				unless true
					a = a * 2
				end
			`,
			wantStackTop: value.Nil,
		},
		"execute the then branch instead of else": {
			source: `
				a := false
				unless a
					a = 10
				else
					a = a + 2
				end
			`,
			wantStackTop: value.SmallInt(10),
		},
		"execute the else branch instead of then": {
			source: `
				a := 5
				unless a
					a = 30
				else
					a = a + 2
				end
			`,
			wantStackTop: value.SmallInt(7),
		},
		"is an expression": {
			source: `
				a := 5
				b := unless a
					"foo"
				else
					5
				end
				b
			`,
			wantStackTop: value.SmallInt(5),
		},
		"modifier binds more strongly than assignment": {
			source: `
				a := 5
				b := "foo" unless a
				b
			`,
			wantCompileErr: error.ErrorList{
				error.NewFailure(L(P(40, 4, 5), P(40, 4, 5)), "undeclared variable: b"),
			},
		},
		"modifier returns the left side if the condition is satisfied": {
			source: `
				a := nil
				"foo" unless a
			`,
			wantStackTop: value.String("foo"),
		},
		"modifier returns nil if the condition is not satisfied": {
			source: `
				a := 5
				"foo" unless a
			`,
			wantStackTop: value.Nil,
		},
		"can access variables defined in the condition": {
			source: `
				a unless a := false
			`,
			wantStackTop: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LogicalOrOperator(t *testing.T) {
	tests := sourceTestTable{
		"return right operand if left is nil": {
			source:       "nil || 4",
			wantStackTop: value.SmallInt(4),
		},
		"return right operand (nil) if left is nil": {
			source:       "nil || nil",
			wantStackTop: value.Nil,
		},
		"return right operand (false) if left is nil": {
			source:       "nil || false",
			wantStackTop: value.False,
		},
		"return right operand if left is false": {
			source:       "false || 'foo'",
			wantStackTop: value.String("foo"),
		},
		"return left operand if it's truthy": {
			source:       "3 || 'foo'",
			wantStackTop: value.SmallInt(3),
		},
		"return right nested operand if left are falsy": {
			source:       "false || nil || 4",
			wantStackTop: value.SmallInt(4),
		},
		"return middle nested operand if left is falsy": {
			source:       "false || 2 || 5",
			wantStackTop: value.SmallInt(2),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_LogicalAndOperator(t *testing.T) {
	tests := sourceTestTable{
		"return left operand if left is nil": {
			source:       "nil && 4",
			wantStackTop: value.Nil,
		},
		"return left operand if left is false": {
			source:       "false && 'foo'",
			wantStackTop: value.False,
		},
		"return right operand if left is truthy": {
			source:       "3 && 'foo'",
			wantStackTop: value.String("foo"),
		},
		"return right operand (false) if left is truthy": {
			source:       "3 && false",
			wantStackTop: value.False,
		},
		"return right nested operand if left are truthy": {
			source:       "4 && 'bar' && 16",
			wantStackTop: value.SmallInt(16),
		},
		"return middle nested operand if left is truthy": {
			source:       "4 && nil && 5",
			wantStackTop: value.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}

func TestVMSource_NilCoalescingOperator(t *testing.T) {
	tests := sourceTestTable{
		"return right operand if left is nil": {
			source:       "nil ?? 4",
			wantStackTop: value.SmallInt(4),
		},
		"return right operand (nil) if left is nil": {
			source:       "nil ?? nil",
			wantStackTop: value.Nil,
		},
		"return right operand (false) if left is nil": {
			source:       "nil ?? false",
			wantStackTop: value.False,
		},
		"return left operand if left is false": {
			source:       "false ?? 'foo'",
			wantStackTop: value.False,
		},
		"return left operand if it's not nil": {
			source:       "3 ?? 'foo'",
			wantStackTop: value.SmallInt(3),
		},
		"return right nested operand if left are nil": {
			source:       "nil ?? nil ?? 4",
			wantStackTop: value.SmallInt(4),
		},
		"return middle nested operand if left is nil": {
			source:       "nil ?? false ?? 5",
			wantStackTop: value.False,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
