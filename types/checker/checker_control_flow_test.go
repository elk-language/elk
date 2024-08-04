package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestBreakExpression(t *testing.T) {
	tests := testTable{
		"the return type is never": {
			input: `
				loop
					a := break
					a = 4
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 4, 10), P(35, 4, 10)), "type `4` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(31, 4, 6), P(35, 4, 10)), "unreachable code"),
			},
		},
		"outside of a loop": {
			input: `break`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(4, 1, 5)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"outside of a loop with a nonexistent label": {
			input: `break$foo`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(8, 1, 9)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"with a nonexistent label": {
			input: `
				loop
					break$foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(15, 3, 6), P(23, 3, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"with a displaced label": {
			input: `
				$foo: loop; end
				loop
					break$foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 4, 6), P(43, 4, 14)), "label $foo does not exist or is not attached to an enclosing loop"),
				error.NewWarning(L("<main>", P(25, 3, 5), P(51, 5, 7)), "unreachable code"),
			},
		},
		"with a valid label": {
			input: `
				$foo: loop
					loop
						break$foo
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

func TestContinueExpression(t *testing.T) {
	tests := testTable{
		"the return type is never": {
			input: `
				loop
					a := continue
					a = 4
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(38, 4, 10), P(38, 4, 10)), "type `4` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(34, 4, 6), P(38, 4, 10)), "unreachable code"),
			},
		},
		"outside of a loop": {
			input: `continue`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(7, 1, 8)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"outside of a loop with a nonexistent label": {
			input: `continue$foo`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(0, 1, 1), P(11, 1, 12)), "cannot jump with `break` or `continue` outside of a loop"),
			},
		},
		"with a nonexistent label": {
			input: `
				loop
					continue$foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(15, 3, 6), P(26, 3, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
			},
		},
		"with a displaced label": {
			input: `
				$foo: loop; end
				loop
					continue$foo
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(35, 4, 6), P(46, 4, 17)), "label $foo does not exist or is not attached to an enclosing loop"),
				error.NewWarning(L("<main>", P(25, 3, 5), P(54, 5, 7)), "unreachable code"),
			},
		},
		"with a valid label": {
			input: `
				$foo: loop
					loop
						continue$foo
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

func TestReturnExpression(t *testing.T) {
	tests := testTable{
		"the return type is never": {
			input: `
				a := return
				a = 4
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(25, 3, 9), P(25, 3, 9)), "type `4` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(21, 3, 5), P(25, 3, 9)), "unreachable code"),
			},
		},
		"warn about values returned in the top level": {
			input: `
				return 4
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(12, 2, 12), P(12, 2, 12)), "values returned in void context will be ignored"),
			},
		},
		"warn about values returned in void methods": {
			input: `
				def foo
					return 4
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(25, 3, 13), P(25, 3, 13)), "values returned in void context will be ignored"),
			},
		},
		"accept matching return type": {
			input: `
				def foo: String
					return "foo"
				end
			`,
		},
		"invalid return type": {
			input: `
				def foo: String
					return 2
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(26, 3, 6), P(33, 3, 13)), "type `2` cannot be assigned to type `Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestDoExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				a := 5
				do
					var b: Int = a
				end
			`,
		},
		"returns the last expression": {
			input: `
				a := 2
				var b: Int = do
					"foo" + "bar"
					a + 2
				end
			`,
		},
		"returns nil when empty": {
			input: `
				var b: nil = do; end
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestNumericForExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				fornum ;;
					var b: Int? = a
				end
			`,
		},
		"use variables defined in the header": {
			input: `
				fornum i := 0; i < 8; i++
					var b: Int = i
				end
			`,
		},
		"typecheck the header and body": {
			input: `
				fornum a; b; c
					d
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(12, 2, 12), P(12, 2, 12)), "undefined local `a`"),
				error.NewFailure(L("<main>", P(15, 2, 15), P(15, 2, 15)), "undefined local `b`"),
				error.NewFailure(L("<main>", P(18, 2, 18), P(18, 2, 18)), "undefined local `c`"),
				error.NewFailure(L("<main>", P(25, 3, 6), P(25, 3, 6)), "undefined local `d`"),
			},
		},
		"returns never if condition is truthy": {
			input: `
				a := 2
				b := fornum ;true;
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(29, 3, 18), P(32, 3, 21)), "this condition will always have the same result since type `true` is truthy"),
				error.NewFailure(L("<main>", P(81, 7, 9), P(81, 7, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(77, 7, 5), P(81, 7, 9)), "unreachable code"),
			},
		},
		"returns never when there is no condition": {
			input: `
				a := 2
				b := fornum ;;
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(77, 7, 9), P(77, 7, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(73, 7, 5), P(77, 7, 9)), "unreachable code"),
			},
		},
		"returns nil if condition is falsy": {
			input: `
				a := 2
				var b: nil = fornum ;false;
					"foo" + "bar"
					a + 2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(37, 3, 26), P(41, 3, 30)), "this loop will never execute since type `false` is falsy"),
				error.NewWarning(L("<main>", P(49, 4, 6), P(62, 4, 19)), "unreachable code"),
			},
		},
		"returns a nilable body type if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = fornum ;a;
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(72, 5, 7)), "type `Std::String?` cannot be assigned to type `8`"),
			},
		},
		"cannot use void in the condition": {
			input: `
				def foo; end
				fornum foo(); foo(); foo()
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 19), P(40, 3, 23)), "cannot use type `void` as a value in this context"),
			},
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				fornum ;a;
					var b: true = a
				end
			`,
		},

		"returns nil with a naked break if condition is truthy": {
			input: `
				a := 2
				b := fornum ;true;
					c := false
					if c
						break
					end
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(29, 3, 18), P(32, 3, 21)), "this condition will always have the same result since type `true` is truthy"),
				error.NewFailure(L("<main>", P(109, 10, 9), P(109, 10, 9)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"returns the value given to break if condition is truthy": {
			input: `
				a := 2
				b := fornum ;true;
					c := false
					if c
						break "foo" + "bar"
					end
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(29, 3, 18), P(32, 3, 21)), "this condition will always have the same result since type `true` is truthy"),
				error.NewFailure(L("<main>", P(123, 10, 9), P(123, 10, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},
		"returns nil with a break if condition is falsy": {
			input: `
				a := 2
				var b: nil = fornum ;false;
					c := false
					if c
						break "foo" + "bar"
					end
					a + 2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(37, 3, 26), P(41, 3, 30)), "this loop will never execute since type `false` is falsy"),
				error.NewWarning(L("<main>", P(49, 4, 6), P(59, 4, 16)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = fornum ;a;
					c := false
					if c
						break
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(119, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = fornum ;a;
					c := false
					if c
						break 2.5
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(123, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"break from a nested labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: fornum ;a;
					var b: Int? = 9
					fornum ;b;
						break$foo 2.5
						"foo" + "bar"
					end
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(116, 7, 7), P(128, 7, 19)), "unreachable code"),
				error.NewFailure(L("<main>", P(36, 3, 16), P(145, 9, 7)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},

		"returns never with a naked continue if condition is truthy": {
			input: `
				a := 2
				b := fornum ;true;
					c := false
					if c
						continue
					end
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(29, 3, 18), P(32, 3, 21)), "this condition will always have the same result since type `true` is truthy"),
				error.NewFailure(L("<main>", P(112, 10, 9), P(112, 10, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(108, 10, 5), P(112, 10, 9)), "unreachable code"),
			},
		},
		"returns never with continue if condition is truthy": {
			input: `
				a := 2
				b := fornum ;true;
					continue "foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(29, 3, 18), P(32, 3, 21)), "this condition will always have the same result since type `true` is truthy"),
				error.NewWarning(L("<main>", P(68, 5, 6), P(72, 5, 10)), "unreachable code"),
				error.NewFailure(L("<main>", P(90, 7, 9), P(90, 7, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(86, 7, 5), P(90, 7, 9)), "unreachable code"),
			},
		},
		"returns nil with a continue if condition is falsy": {
			input: `
				a := 2
				var b: nil = fornum ;false;
					c := false
					if c
						continue "foo" + "bar"
					end
					a + 2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(37, 3, 26), P(41, 3, 30)), "this loop will never execute since type `false` is falsy"),
				error.NewWarning(L("<main>", P(49, 4, 6), P(59, 4, 16)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = fornum ;a;
					var c: Int? = 1
					if c
						continue
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(127, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = fornum ;a;
					c := false
					if c
						continue 2.5
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(126, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"continue a parent labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: fornum ;a;
					var b: Int? = 9
					fornum ;b;
						continue$foo 2.5
						"foo" + "bar"
					end
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(119, 7, 7), P(131, 7, 19)), "unreachable code"),
				error.NewFailure(L("<main>", P(36, 3, 16), P(148, 9, 7)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestLoopExpression(t *testing.T) {
	tests := testTable{
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				loop
					var b: Int? = a
				end
			`,
		},
		"returns never": {
			input: `
				a := 2
				b := loop
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(72, 7, 9), P(72, 7, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(68, 7, 5), P(72, 7, 9)), "unreachable code"),
			},
		},
		"returns nil when a naked break is present": {
			input: `
				var a: Int? = 2
				b := loop
					if a
						break
					end
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(82, 8, 9), P(82, 8, 9)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"returns the value given to break": {
			input: `
				var a: Int? = 2
				var b = loop
					if a
						break ""
					end
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(88, 8, 9), P(88, 8, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},
		"returns the union of values given to break": {
			input: `
				var a: Int? = 2
				var b = loop
					if a
						break ""
					else
						break 2.5
					end
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(114, 10, 9), P(114, 10, 9)), "type `3` cannot be assigned to type `\"\" | 2.5`"),
			},
		},
		"break nested labeled loop": {
			input: `
				var a: Int? = 2
				var b = $foo: loop
					loop
						break$foo ""
					end
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(98, 8, 9), P(98, 8, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},

		"returns never when a naked continue is present": {
			input: `
				var a: Int? = 2
				b := loop
					if a
						continue
					end
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(85, 8, 9), P(85, 8, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(81, 8, 5), P(85, 8, 9)), "unreachable code"),
			},
		},
		"does not return the value given to continue": {
			input: `
				var a: Int? = 2
				var b = loop
					if a
						continue ""
					end
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(91, 8, 9), P(91, 8, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(87, 8, 5), P(91, 8, 9)), "unreachable code"),
			},
		},
		"continue in nested labeled loop": {
			input: `
				var a: Int? = 2
				var b = $foo: loop
					loop
						continue$foo ""
					end
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(101, 8, 9), P(101, 8, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(97, 8, 5), P(101, 8, 9)), "unreachable code"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestWhileExpression(t *testing.T) {
	tests := testTable{
		"cannot use void in the condition": {
			input: `
				def foo; end
				while foo()
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 3, 11), P(32, 3, 15)), "cannot use type `void` as a value in this context"),
			},
		},
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				while a
					var b: Int? = a
				end
			`,
		},
		"returns never if condition is truthy": {
			input: `
				a := 2
				b := while true
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 3, 16), P(30, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
				error.NewFailure(L("<main>", P(78, 7, 9), P(78, 7, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(74, 7, 5), P(78, 7, 9)), "unreachable code"),
			},
		},
		"returns nil if condition is falsy": {
			input: `
				a := 2
				var b: nil = while false
					"foo" + "bar"
					a + 2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(35, 3, 24), P(39, 3, 28)), "this loop will never execute since type `false` is falsy"),
				error.NewWarning(L("<main>", P(46, 4, 6), P(59, 4, 19)), "unreachable code"),
			},
		},
		"returns a nilable body type if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = while a
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(69, 5, 7)), "type `Std::String?` cannot be assigned to type `8`"),
			},
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				while a
					var b: true = a
				end
			`,
		},

		"returns nil with a naked break if condition is truthy": {
			input: `
				a := 2
				b := while true
					c := false
					if c
						break
					end
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 3, 16), P(30, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
				error.NewFailure(L("<main>", P(106, 10, 9), P(106, 10, 9)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"returns the value given to break if condition is truthy": {
			input: `
				a := 2
				b := while true
					c := false
					if c
						break "foo" + "bar"
					end
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 3, 16), P(30, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
				error.NewFailure(L("<main>", P(120, 10, 9), P(120, 10, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},
		"returns nil with a break if condition is falsy": {
			input: `
				a := 2
				var b: nil = while false
					c := false
					if c
						break "foo" + "bar"
					end
					a + 2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(35, 3, 24), P(39, 3, 28)), "this loop will never execute since type `false` is falsy"),
				error.NewWarning(L("<main>", P(46, 4, 6), P(56, 4, 16)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = while a
					c := false
					if c
						break
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(116, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = while a
					c := false
					if c
						break 2.5
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(120, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"break from a nested labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: while a
					var b: Int? = 9
					while b
						break$foo 2.5
						"foo" + "bar"
					end
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(110, 7, 7), P(122, 7, 19)), "unreachable code"),
				error.NewFailure(L("<main>", P(36, 3, 16), P(139, 9, 7)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},

		"returns never with a naked continue if condition is truthy": {
			input: `
				a := 2
				b := while true
					continue
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 3, 16), P(30, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
				error.NewWarning(L("<main>", P(51, 5, 6), P(55, 5, 10)), "unreachable code"),
				error.NewFailure(L("<main>", P(73, 7, 9), P(73, 7, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(69, 7, 5), P(73, 7, 9)), "unreachable code"),
			},
		},
		"returns never with continue if condition is truthy": {
			input: `
				a := 2
				b := while true
					continue "foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 3, 16), P(30, 3, 19)), "this condition will always have the same result since type `true` is truthy"),
				error.NewWarning(L("<main>", P(65, 5, 6), P(69, 5, 10)), "unreachable code"),
				error.NewFailure(L("<main>", P(87, 7, 9), P(87, 7, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(83, 7, 5), P(87, 7, 9)), "unreachable code"),
			},
		},
		"returns nil with a continue if condition is falsy": {
			input: `
				a := 2
				var b: nil = while false
					continue "foo" + "bar"
					a + 2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(35, 3, 24), P(39, 3, 28)), "this loop will never execute since type `false` is falsy"),
				error.NewWarning(L("<main>", P(46, 4, 6), P(68, 4, 28)), "unreachable code"),
				error.NewWarning(L("<main>", P(74, 5, 6), P(78, 5, 10)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = while a
					c := false
					if c
						continue
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(119, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = while a
					c := false
					if c
						continue 2.5
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(123, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"continue a parent labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: while a
					var b: Int? = 9
					while b
						continue$foo 2.5
						"foo" + "bar"
					end
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(113, 7, 7), P(125, 7, 19)), "unreachable code"),
				error.NewFailure(L("<main>", P(36, 3, 16), P(142, 9, 7)), "type `nil | 2.5` cannot be assigned to type `8`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestUntilExpression(t *testing.T) {
	tests := testTable{
		"cannot use void in the condition": {
			input: `
				def foo; end
				until foo()
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 3, 11), P(32, 3, 15)), "cannot use type `void` as a value in this context"),
			},
		},
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				until a
					var b: Int? = a
				end
			`,
		},
		"returns never if condition is falsy": {
			input: `
				a := 2
				b := until false
					"foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 3, 16), P(31, 3, 20)), "this condition will always have the same result since type `false` is falsy"),
				error.NewFailure(L("<main>", P(79, 7, 9), P(79, 7, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(75, 7, 5), P(79, 7, 9)), "unreachable code"),
			},
		},
		"returns nil if condition is truthy": {
			input: `
				a := 2
				var b: nil = until true
					"foo" + "bar"
					a + 2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(35, 3, 24), P(38, 3, 27)), "this loop will never execute since type `true` is truthy"),
				error.NewWarning(L("<main>", P(45, 4, 6), P(58, 4, 19)), "unreachable code"),
			},
		},
		"returns a nilable body type if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = until a
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(69, 5, 7)), "type `Std::String?` cannot be assigned to type `8`"),
			},
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				until a
					var b: false = a
				end
			`,
		},

		"returns nil with a naked break if condition is falsy": {
			input: `
				a := 2
				b := until false
					break
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 3, 16), P(31, 3, 20)), "this condition will always have the same result since type `false` is falsy"),
				error.NewWarning(L("<main>", P(49, 5, 6), P(53, 5, 10)), "unreachable code"),
				error.NewFailure(L("<main>", P(71, 7, 9), P(71, 7, 9)), "type `3` cannot be assigned to type `nil`"),
			},
		},
		"returns the value given to break if condition is falsy": {
			input: `
				a := 2
				b := until false
					break "foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 3, 16), P(31, 3, 20)), "this condition will always have the same result since type `false` is falsy"),
				error.NewWarning(L("<main>", P(63, 5, 6), P(67, 5, 10)), "unreachable code"),
				error.NewFailure(L("<main>", P(85, 7, 9), P(85, 7, 9)), "type `3` cannot be assigned to type `Std::String`"),
			},
		},
		"returns nil with a break if condition is truthy": {
			input: `
				a := 2
				var b: nil = until true
					break "foo" + "bar"
					a + 2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(35, 3, 24), P(38, 3, 27)), "this loop will never execute since type `true` is truthy"),
				error.NewWarning(L("<main>", P(45, 4, 6), P(64, 4, 25)), "unreachable code"),
				error.NewWarning(L("<main>", P(70, 5, 6), P(74, 5, 10)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = until a
					var c: Int? = 3
					if c
						break
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(121, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to break if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = until a
					var c: Int? = 3
					if c
						break 2.5
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(125, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"break from a nested labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: until a
					var b: Int? = 9
					until b
						var c: Int? = 3
						if c
							break$foo 2.5
						end
						"foo" + "bar"
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(183, 12, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},

		"returns never with a naked continue if condition is falsy": {
			input: `
				a := 2
				b := until false
					continue
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 3, 16), P(31, 3, 20)), "this condition will always have the same result since type `false` is falsy"),
				error.NewWarning(L("<main>", P(52, 5, 6), P(56, 5, 10)), "unreachable code"),
				error.NewFailure(L("<main>", P(74, 7, 9), P(74, 7, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(70, 7, 5), P(74, 7, 9)), "unreachable code"),
			},
		},
		"returns never with continue if condition is falsy": {
			input: `
				a := 2
				b := until false
					continue "foo" + "bar"
					a + 2
				end
				b = 3
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(27, 3, 16), P(31, 3, 20)), "this condition will always have the same result since type `false` is falsy"),
				error.NewWarning(L("<main>", P(66, 5, 6), P(70, 5, 10)), "unreachable code"),
				error.NewFailure(L("<main>", P(88, 7, 9), P(88, 7, 9)), "type `3` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(84, 7, 5), P(88, 7, 9)), "unreachable code"),
			},
		},
		"returns nil with a continue if condition is truthy": {
			input: `
				a := 2
				var b: nil = until true
					continue "foo" + "bar"
					a + 2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(35, 3, 24), P(38, 3, 27)), "this loop will never execute since type `true` is truthy"),
				error.NewWarning(L("<main>", P(45, 4, 6), P(67, 4, 28)), "unreachable code"),
				error.NewWarning(L("<main>", P(73, 5, 6), P(77, 5, 10)), "unreachable code"),
			},
		},
		"returns a nilable body type with naked continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = until a
					var c: Int? = 3
					if c
						continue
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(124, 9, 7)), "type `nil | Std::String` cannot be assigned to type `8`"),
			},
		},
		"returns a union of body type, nil and the value given to continue if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = until a
					var c: Int? = 3
					if c
						continue 2.5
					end
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(128, 9, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"continue a parent labeled loop": {
			input: `
				var a: Int? = 2
				var b: 8 = $foo: until a
					var b: Int? = 9
					until b
						var c: Int? = 3
						if c
							continue$foo 2.5
						end
						"foo" + "bar"
					end
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(186, 12, 7)), "type `2.5 | Std::String | nil` cannot be assigned to type `8`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestUnlessExpression(t *testing.T) {
	tests := testTable{
		"cannot use void in the condition": {
			input: `
				def foo; end
				unless foo()
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(29, 3, 12), P(33, 3, 16)), "cannot use type `void` as a value in this context"),
			},
		},
		"checks modifier version": {
			input: `
				def foo; end
				nil unless foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(33, 3, 16), P(37, 3, 20)), "cannot use type `void` as a value in this context"),
			},
		},
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				unless a
					var b: Int? = a
				else
					a
				end
			`,
		},
		"returns the last else expression if condition is truthy": {
			input: `
				a := 2
				var b: Float = unless true
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(38, 3, 27), P(41, 3, 30)), "this condition will always have the same result since type `true` is truthy"),
				error.NewWarning(L("<main>", P(48, 4, 6), P(61, 4, 19)), "unreachable code"),
			},
		},
		"returns the last then expression if condition is falsy": {
			input: `
				a := 2
				var b: Int = unless false
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(36, 3, 25), P(40, 3, 29)), "this condition will always have the same result since type `false` is falsy"),
				error.NewWarning(L("<main>", P(86, 7, 6), P(89, 7, 9)), "unreachable code"),
			},
		},
		"returns a union of both branches if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = unless a
					"foo" + "bar"
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(88, 7, 7)), "type `Std::String | 2.2` cannot be assigned to type `8`"),
			},
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				unless a
					var b: false = a
				else
					var b: true = a
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

func TestIfExpression(t *testing.T) {
	tests := testTable{
		"cannot use void in the condition": {
			input: `
				def foo; end
				if foo()
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(25, 3, 8), P(29, 3, 12)), "cannot use type `void` as a value in this context"),
			},
		},
		"checks modifier version": {
			input: `
				def foo; end
				nil if foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(29, 3, 12), P(33, 3, 16)), "cannot use type `void` as a value in this context"),
			},
		},
		"checks modifier version with else": {
			input: `
				a := 2
				var b: Int = (a + 2 if true else 2.2)
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(39, 3, 28), P(42, 3, 31)), "this condition will always have the same result since type `true` is truthy"),
				error.NewWarning(L("<main>", P(49, 3, 38), P(51, 3, 40)), "unreachable code"),
			},
		},
		"has access to outer variables": {
			input: `
				var a: Int? = 5
				if a
					var b: Int? = a
				else
					a
				end
			`,
		},
		"returns the last then expression if condition is truthy": {
			input: `
				a := 2
				var b: Int = if true
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(32, 3, 21), P(35, 3, 24)), "this condition will always have the same result since type `true` is truthy"),
				error.NewWarning(L("<main>", P(81, 7, 6), P(84, 7, 9)), "unreachable code"),
			},
		},
		"returns the last else expression if condition is truthy": {
			input: `
				a := 2
				var b: Float = if false
					"foo" + "bar"
					a + 2
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(34, 3, 23), P(38, 3, 27)), "this condition will always have the same result since type `false` is falsy"),
				error.NewWarning(L("<main>", P(45, 4, 6), P(58, 4, 19)), "unreachable code"),
			},
		},
		"returns a union of both branches if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = if a
					"foo" + "bar"
				else
					2.2
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(84, 7, 7)), "type `Std::String | 2.2` cannot be assigned to type `8`"),
			},
		},
		"returns a union of all branches if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: Float? = nil
				var c: 8 = if a
					"foo" + "bar"
				else if b
					2.5
				else
					%/foo/
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(60, 4, 16), P(134, 10, 7)), "type `Std::String | 2.5 | Std::Regex` cannot be assigned to type `8`"),
			},
		},
		"returns a union of then and nil if the condition is neither truthy nor falsy": {
			input: `
				var a: Int? = 2
				var b: 8 = if a
					"foo" + "bar"
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(36, 3, 16), P(66, 5, 7)), "type `Std::String | nil` cannot be assigned to type `8`"),
			},
		},
		"returns nil when empty": {
			input: `
				var a: Int? = nil
				var b: nil = if a; end
			`,
		},

		"narrow Bool variable type by using truthiness": {
			input: `
				var a = false
				if a
					var b: true = a
				else
					var b: false = a
				end
			`,
		},
		"narrow nilable variable type by using truthiness": {
			input: `
				var a: Int? = nil
				if a
					var b: Int = a
				else
					var b: nil = a
				end
			`,
		},
		"narrow named nilable variable type by using truthiness": {
			input: `
				typedef Foo = Int?
				var a: Foo = nil
				if a
					var b: Int = a
				else
					var b: nil = a
				end
			`,
		},
		"narrow nilable variable type by using negated truthiness": {
			input: `
				var a: Int? = nil
				if !a
					var b: nil = a
				else
					var b: Int = a
				end
			`,
		},
		"narrow named nilable variable type by using negated truthiness": {
			input: `
				typedef Foo = Int?
				var a: Foo = nil
				if !a
					var b: nil = a
				else
					var b: Int = a
				end
			`,
		},
		"narrow union type by using <:": {
			input: `
				var a: Int | String = "foo"
				if a <: Int
					var b: Int = a
				else
					var b: String = a
				end
			`,
		},
		"narrow variable type by using <:": {
			input: `
				var a: Int? = nil
				if a <: Int
					var b: Int = a
				else
					var b: nil = a
				end
			`,
		},
		"narrow variable type by using <<:": {
			input: `
				var a: Int? = nil
				if a <<: Int
					var b: Int = a
				else
					var b: nil = a
				end
			`,
		},
		"narrow a few variables with &&": {
			input: `
				var a: Int? = nil
				var b = false
				if a && b
					var c: Int = a
					var d: true = b
				end
			`,
		},
		"narrow with an impossible && branch": {
			input: `
				var a: Int? = 3
				if a && nil
					a = :foo
				else
					a = :bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(46, 4, 10), P(49, 4, 13)), "type `:foo` cannot be assigned to type `never`"),
				error.NewFailure(L("<main>", P(69, 6, 10), P(72, 6, 13)), "type `:bar` cannot be assigned to type `Std::Int?`"),
				error.NewWarning(L("<main>", P(28, 3, 8), P(35, 3, 15)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L("<main>", P(42, 4, 6), P(50, 4, 14)), "unreachable code"),
			},
		},
		"narrow a few variables with ||": {
			input: `
				var a: Int? = nil
				var b = false
				if a || b
				else
					var c: nil = a
					var d: false = b
				end
			`,
		},
		"narrow with an impossible || branch": {
			input: `
				var a: Int? = nil
				if a || 3
					a = :foo
				else
					a = :bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(46, 4, 10), P(49, 4, 13)), "type `:foo` cannot be assigned to type `Std::Int?`"),
				error.NewFailure(L("<main>", P(69, 6, 10), P(72, 6, 13)), "type `:bar` cannot be assigned to type `never`"),
				error.NewWarning(L("<main>", P(30, 3, 8), P(35, 3, 13)), "this condition will always have the same result since type `Std::Int` is truthy"),
				error.NewWarning(L("<main>", P(65, 6, 6), P(73, 6, 14)), "unreachable code"),
			},
		},

		"narrow with ===": {
			input: `
				var a: Int | Float = 1
				var b: Float? = .2
				if a === b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(75, 5, 10), P(78, 5, 13)), "type `:foo` cannot be assigned to type `Std::Float`"),
				error.NewFailure(L("<main>", P(89, 6, 10), P(92, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float`"),
				error.NewFailure(L("<main>", P(112, 8, 10), P(115, 8, 13)), "type `:baz` cannot be assigned to type `Std::Int | Std::Float`"),
				error.NewFailure(L("<main>", P(126, 9, 10), P(130, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float?`"),
			},
		},
		"narrow with an impossible ===": {
			input: `
				var a: Int = 1
				var b: Float = .2
				if a === b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(49, 4, 8), P(49, 4, 8)), "this strict equality check is impossible, `Std::Int` cannot ever be equal to `Std::Float`"),
				error.NewFailure(L("<main>", P(66, 5, 10), P(69, 5, 13)), "type `:foo` cannot be assigned to type `never`"),
				error.NewFailure(L("<main>", P(80, 6, 10), P(83, 6, 13)), "type `:bar` cannot be assigned to type `never`"),
				error.NewFailure(L("<main>", P(103, 8, 10), P(106, 8, 13)), "type `:baz` cannot be assigned to type `Std::Int`"),
				error.NewFailure(L("<main>", P(117, 9, 10), P(121, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float`"),
			},
		},
		"narrow with !==": {
			input: `
				var a: Int | Float = 1
				var b: Float? = .2
				if a !== b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(75, 5, 10), P(78, 5, 13)), "type `:foo` cannot be assigned to type `Std::Int | Std::Float`"),
				error.NewFailure(L("<main>", P(89, 6, 10), P(92, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float?`"),
				error.NewFailure(L("<main>", P(112, 8, 10), P(115, 8, 13)), "type `:baz` cannot be assigned to type `Std::Float`"),
				error.NewFailure(L("<main>", P(126, 9, 10), P(130, 9, 14)), "type `:fizz` cannot be assigned to type `Std::Float`"),
			},
		},
		"narrow with an impossible !==": {
			input: `
				var a: Int = 1
				var b: Float = .2
				if a !== b
					a = :foo
					b = :bar
				else
					a = :baz
					b = :fizz
				end
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(49, 4, 8), P(49, 4, 8)), "this strict equality check is impossible, `Std::Int` cannot ever be equal to `Std::Float`"),
				error.NewFailure(L("<main>", P(66, 5, 10), P(69, 5, 13)), "type `:foo` cannot be assigned to type `Std::Int`"),
				error.NewFailure(L("<main>", P(80, 6, 10), P(83, 6, 13)), "type `:bar` cannot be assigned to type `Std::Float`"),
				error.NewFailure(L("<main>", P(103, 8, 10), P(106, 8, 13)), "type `:baz` cannot be assigned to type `never`"),
				error.NewFailure(L("<main>", P(117, 9, 10), P(121, 9, 14)), "type `:fizz` cannot be assigned to type `never`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestLogicalAnd(t *testing.T) {
	tests := testTable{
		"cannot use void on the left hand side": {
			input: `
				def foo; end
				foo() && foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 3, 5), P(26, 3, 9)), "cannot use type `void` as a value in this context"),
			},
		},
		"returns the right type when the left type is truthy": {
			input: `
				var a = "foo"
				var b = 2
				var c: Int = a && b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(50, 4, 18), P(50, 4, 18)), "this condition will always have the same result since type `Std::String` is truthy"),
			},
		},
		"returns the left type when the left type is falsy": {
			input: `
				var a = nil
				var b = 2
				var c: nil = a && b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(48, 4, 18), P(48, 4, 18)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L("<main>", P(53, 4, 23), P(53, 4, 23)), "unreachable code"),
			},
		},
		"returns a union of both types with only nil when the left can be both truthy and falsy": {
			input: `
				var a: String? = "foo"
				var b = 2
				var c: 9 = a && b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 16), P(62, 4, 21)), "type `nil | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types with only false when the left can be both truthy and falsy": {
			input: `
				var a = false
				var b = 2
				var c: 9 = a && b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 4, 16), P(53, 4, 21)), "type `false | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without duplication": {
			input: `
				var a: false | nil | Int = nil
				var b: Float | Int | nil = 2.2
				var c: 9 = a && b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(86, 4, 16), P(91, 4, 21)), "type `false | nil | Std::Float | Std::Int` cannot be assigned to type `9`"),
			},
		},

		"narrow left variable to non falsy": {
			input: `
				var a: false | nil | Int = nil
				a && a + 2
			`,
		},
		"narrow a few variables to non falsy": {
			input: `
				var a: Int? = nil
				var b: Int? = nil
				a && b && a + b
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestLogicalOr(t *testing.T) {
	tests := testTable{
		"cannot use void on the left hand side": {
			input: `
				def foo; end
				foo() || foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 3, 5), P(26, 3, 9)), "cannot use type `void` as a value in this context"),
			},
		},
		"returns the left type when it is truthy": {
			input: `
				var a = "foo"
				var b = 2
				var c: String = a || b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(53, 4, 21), P(53, 4, 21)), "this condition will always have the same result since type `Std::String` is truthy"),
				error.NewWarning(L("<main>", P(58, 4, 26), P(58, 4, 26)), "unreachable code"),
			},
		},
		"returns the right type when the left type is falsy": {
			input: `
				var a = nil
				var b = 2
				var c: Int = a || b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(48, 4, 18), P(48, 4, 18)), "this condition will always have the same result since type `nil` is falsy"),
			},
		},
		"returns a union of both types without nil when the left can be both truthy and falsy": {
			input: `
				var a: String? = "foo"
				var b = 2
				var c: 9 = a || b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 16), P(62, 4, 21)), "type `Std::String | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without false when the left can be both truthy and falsy": {
			input: `
				var a = false
				var b = 2
				var c: 9 = a || b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(48, 4, 16), P(53, 4, 21)), "type `true | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without duplication": {
			input: `
				var a: String | Int | nil = nil
				var b: Float | Int | nil = 2.2
				var c: 9 = a || b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(87, 4, 16), P(92, 4, 21)), "type `Std::String | Std::Int | Std::Float | nil` cannot be assigned to type `9`"),
			},
		},

		"narrow left variable to falsy": {
			input: `
				var a: false | nil | Int = nil
				a || var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(56, 3, 21), P(56, 3, 21)), "type `false | nil` cannot be assigned to type `9`"),
			},
		},
		"narrow a few variables to non falsy": {
			input: `
				var a: Int? = nil
				var b: Int? = nil
				a || b || var c: 9 = a && b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(70, 4, 26), P(70, 4, 26)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L("<main>", P(75, 4, 31), P(75, 4, 31)), "unreachable code"),
				error.NewFailure(L("<main>", P(70, 4, 26), P(75, 4, 31)), "type `nil` cannot be assigned to type `9`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestNilCoalescing(t *testing.T) {
	tests := testTable{
		"cannot use void on the left hand side": {
			input: `
				def foo; end
				foo() ?? foo()
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(22, 3, 5), P(26, 3, 9)), "cannot use type `void` as a value in this context"),
			},
		},
		"returns the left type when it is not nilable": {
			input: `
				var a = "foo"
				var b = 2
				var c: String = a ?? b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(53, 4, 21), P(53, 4, 21)), "this condition will always have the same result since type `Std::String` can never be nil"),
				error.NewWarning(L("<main>", P(58, 4, 26), P(58, 4, 26)), "unreachable code"),
			},
		},
		"returns the right type when the left type is nil": {
			input: `
				var a = nil
				var b = 2
				var c: Int = a ?? b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(48, 4, 18), P(48, 4, 18)), "this condition will always have the same result"),
			},
		},
		"returns a union of both types without nil when the left can be both nil and not nil": {
			input: `
				var a: String? = "foo"
				var b = 2
				var c: 9 = a ?? b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(57, 4, 16), P(62, 4, 21)), "type `Std::String | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"returns a union of both types without duplication": {
			input: `
				var a: String | Int | nil = nil
				var b: Float | Int | nil = 2.2
				var c: 9 = a ?? b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(87, 4, 16), P(92, 4, 21)), "type `Std::String | Std::Int | Std::Float | nil` cannot be assigned to type `9`"),
			},
		},

		"narrow left variable to non nilable": {
			input: `
				var a: false | nil | Int = nil
				a ?? var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(56, 3, 21), P(56, 3, 21)), "type `nil` cannot be assigned to type `9`"),
			},
		},
		"narrow a few variables to non nilable": {
			input: `
				var a: Int? = nil
				var b: Int? = nil
				a ?? b ?? var c: 9 = a && b
			`,
			err: error.ErrorList{
				error.NewWarning(L("<main>", P(70, 4, 26), P(70, 4, 26)), "this condition will always have the same result since type `nil` is falsy"),
				error.NewWarning(L("<main>", P(75, 4, 31), P(75, 4, 31)), "unreachable code"),
				error.NewFailure(L("<main>", P(70, 4, 26), P(75, 4, 31)), "type `nil` cannot be assigned to type `9`"),
			},
		},
		"narrow nested ||": {
			input: `
				var a: bool? = false
				var b: bool? = false
				(a || b) ?? do
					a = :foo
					b = :bar
				end
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(79, 5, 10), P(82, 5, 13)), "type `:foo` cannot be assigned to type `nil`"),
				error.NewFailure(L("<main>", P(93, 6, 10), P(96, 6, 13)), "type `:bar` cannot be assigned to type `nil`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
