package vm_test

import (
	"testing"

	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"
)

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
				"method `iterator` is not available to value of class `Std::Object`: Std::Object{}",
			),
		},
		"loop over an invalid iterable": {
			source: `
				class InvalidIterator
					def iterator then self
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
		"loop over a string byte iterator": {
			source: `
				for i in "Poznań jest √🔥".byte_iterator
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
			wantCompileErr: errors.ErrorList{
				errors.NewError(L(P(43, 4, 5), P(43, 4, 5)), "undeclared variable: b"),
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
			wantCompileErr: errors.ErrorList{
				errors.NewError(L(P(40, 4, 5), P(40, 4, 5)), "undeclared variable: b"),
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

func TestVMSource_Switch(t *testing.T) {
	tests := sourceTestTable{
		"match no value": {
			source: `
				switch 20
		    case 0 then :a
				case 5 then :b
				case 10 then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.Nil,
		},
		"match with variable": {
			source: `
				switch 20
		    case 0 then :a
				case n then n + 2
				case 10 then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.SmallInt(22),
		},
		"match no value with else": {
			source: `
				switch 20
		    case 0 then :a
				case 5 then :b
				case 10 then :c
				case 15 then :d
				else :e
				end
			`,
			wantStackTop: value.ToSymbol("e"),
		},
		"match boolean": {
			source: `
				switch true
		    case 0 then :a
				case false then :b
				case true then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match nil": {
			source: `
				switch nil
		    case 0 then :a
				case nil then :b
				case true then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("b"),
		},
		"match string": {
			source: `
				switch "some string"
		    case 0 then :a
				case nil then :b
				case "some other string" then :b
				case "some string" then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match interpolated string": {
			source: `
				switch "some 7 string"
		    case 0 then :a
				case nil then :b
				case "some other string" then :b
				case "some ${5 + 2} string" then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match string with regex": {
			source: `
				switch "some string"
		    case 0 then :a
				case nil then :b
				case %/^some other string$/ then :b
				case %/some str\w+/ then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match string with interpolated regex": {
			source: `
				switch "some 7 string"
		    case 0 then :a
				case nil then :b
				case %/^some other string$/ then :b
				case %/some ${2 + 5} str\w+/ then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match symbol": {
			source: `
				switch :foo
		    case 0 then :a
				case :bar then :b
				case "foo" then :b
				case :foo then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match interpolated symbol": {
			source: `
				switch :foo7
		    case 0 then :a
				case :bar then :b
				case "foo7" then :b
				case :"foo${2 + 5}" then :c
				case 15 then :d
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match comparison with and": {
			source: `
				switch 10
		    case > 20 then :a
				case > 5 && < 8 then :b
				case > 9 && < 15 then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match or": {
			source: `
				switch "foo"
		    case "bar" || "foo" then :a
				case > 5 && < 8 then :b
				case > 9 && < 15 then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("a"),
		},
		"match non-list with list patterns": {
			source: `
				switch :foo
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},
		"match tuple with list patterns": {
			source: `
				switch %[1, 6, 9, 20]
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},
		"match list": {
			source: `
				switch [1, 6, 9, 20]
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match empty list": {
			source: `
				switch []
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match list with rest elements": {
			source: `
				switch [1, 6, 9, 20]
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, > 5, 2] then :c
				case [1, *a, > 15] then a
				case 15 then :e
				end
			`,
			wantStackTop: &value.ArrayList{
				value.SmallInt(6),
				value.SmallInt(9),
			},
		},
		"match list with unnamed rest elements": {
			source: `
				switch [1, 6, 9, 20]
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, > 5, 2] then :c
				case [1, *, < 15] then :d
				case [1, *, > 15] then :e
				case 15 then :f
				end
			`,
			wantStackTop: value.ToSymbol("e"),
		},
		"match nested lists": {
			source: `
				switch [1, 6, [17, 43, [71, 28]], 20]
		    case < 9 then :a
				case [1, 6, 10] then :b
				case [< 2, 6, [17, 43, [42, 28]], 20] then :c
				case [1, 6, [17, > 40, [71, 28]], > 15] then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},

		"match non-list with tuple patterns": {
			source: `
				switch :foo
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},
		"match list with tuple patterns": {
			source: `
				switch [1, 6, 9, 20]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match tuple": {
			source: `
				switch %[1, 6, 9, 20]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, > 5, 20] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match empty tuple": {
			source: `
				switch %[]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[] then :c
				case == 10 then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("c"),
		},
		"match tuple with rest elements": {
			source: `
				switch %[1, 6, 9, 20]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, > 5, 2] then :c
				case %[1, *a, > 15] then a
				case 15 then :e
				end
			`,
			wantStackTop: &value.ArrayList{
				value.SmallInt(6),
				value.SmallInt(9),
			},
		},
		"match tuple with unnamed rest elements": {
			source: `
				switch %[1, 6, 9, 20]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, > 5, 2] then :c
				case %[1, *, < 15] then :d
				case %[1, *, > 15] then :e
				case 15 then :f
				end
			`,
			wantStackTop: value.ToSymbol("e"),
		},
		"match nested tuples": {
			source: `
				switch %[1, 6, %[17, 43, %[71, 28]], 20]
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %[< 2, 6, %[17, 43, %[42, 28]], 20] then :c
				case %[1, 6, %[17, > 40, %[71, 28]], > 15] then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},

		"match empty map": {
			source: `
				switch {}
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case { 1 => > 2, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case {} then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
		"match map": {
			source: `
				switch { 1 => 5.5, foo: "bar", "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case { 1 => > 2, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case { 1 => > 2, foo: "baz" || "bar", "baz" => < 13 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
		"match map with record": {
			source: `
				switch { 1 => 5.5, foo: "bar", "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %{ 1 => > 2, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case %{ 1 => > 2, foo: "baz" || "bar", "baz" => < 13 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
		"match nested maps": {
			source: `
				switch { 1 => 5.5, foo: ["bar", 5, 4, { elo: "mordo" }], "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case { 1 => > 2, foo: ["baz" || "bar", 5, 4, { eli: "mordo" }], "baz" => < 13 } then :c
				case { 1 => > 2, foo: ["baz" || "bar", 5, 4, { elo: %/^mord\w+$/ }], "baz" => < 13 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},

		"match empty record": {
			source: `
				switch %{}
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %{ 1 => > 2, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case %{} then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
		"match record": {
			source: `
				switch %{ 1 => 5.5, foo: "bar", "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %{ 1 => > 2, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case %{ 1 => > 2, foo: "baz" || "bar", "baz" => < 13 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
		"match record with map": {
			source: `
				switch %{ 1 => 5.5, foo: "bar", "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case { 1 => > 2, foo: "baz" || "bar", "baz" => 12.2 } then :c
				case { 1 => > 2, foo: "baz" || "bar", "baz" => < 13 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.Nil,
		},
		"match nested records": {
			source: `
				switch %{ 1 => 5.5, foo: ["bar", 5, 4, %{ elo: "mordo" }], "baz" => 12.5 }
		    case < 9 then :a
				case %[1, 6, 10] then :b
				case %{ 1 => > 2, foo: ["baz" || "bar", 5, 4, %{ eli: "mordo" }], "baz" => < 13 } then :c
				case %{ 1 => > 2, foo: ["baz" || "bar", 5, 4, %{ elo: %/^mord\w+$/ }], "baz" => < 13 } then :d
				case 15 then :e
				end
			`,
			wantStackTop: value.ToSymbol("d"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
