package vm_test

import (
	"testing"

	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/value"
)

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
