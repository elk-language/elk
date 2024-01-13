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
					print(i.inspect, " ")
				end
			`,
			wantStackTop: value.Nil,
			wantStdout:   `c"P" c"o" c"z" c"n" c"a" c"ń" c" " c"j" c"e" c"s" c"t" c" " c"√" c"🔥" `,
		},
		"loop over a tuple": {
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
				for i := 1; i < 6; i += 1
					a += i
				end
				a
			`,
			wantStackTop: value.SmallInt(15),
		},
		"create a repeated string": {
			source: `
				a := ""
				for i := 20; i > 0; i -= 2
					a += "-"
				end
				a
			`,
			wantStackTop: value.String("----------"),
		},
		"calculate the factorial of 10": {
			source: `
				a := 1
				for i := 2; i <= 10; i += 1
					a *= i
				end
				a
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"return the value of the last iteration": {
			source: `
				a := 1
				for i := 2; i <= 10; i += 1
					a *= i
				end
			`,
			wantStackTop: value.SmallInt(3628800),
		},
		"return nil when no iterations": {
			source: `
				a := 1
				for i := 20; i <= 10; i += 1
					a *= i
				end
			`,
			wantStackTop: value.Nil,
		},
		"return nil after break": {
			source: `
				a := 1
				for i := 2; i <= 10; i += 1
					a *= i
					break if a > 200
				end
			`,
			wantStackTop: value.Nil,
		},
		"return a value using break": {
			source: `
				a := 1
				for i := 2; i <= 10; i += 1
					a *= i
					break a if a > 200
				end
			`,
			wantStackTop: value.SmallInt(720),
		},
		"nested with continue": {
			source: `
				for j := 1; j <= 5; j += 1
					for i := 1; i <= 5; i += 1
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
				$foo: for j := 1; j <= 5; j += 1
					for i := 1; i <= 5; i += 1
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
				for j := 1;; j += 1
					for i := 1;; i += 1
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
				$foo: for j := 1;; j += 1
					for i := 1;; i += 1
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
