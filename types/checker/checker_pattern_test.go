package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
)

func TestPatterns(t *testing.T) {
	tests := testTable{
		"public identifier pattern": {
			input: `
				var [a] = [1]
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 16), P(34, 3, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"redeclare public identifier pattern": {
			input: `
				var a: Int
				var [a] = [1]
				var b: 9 = a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(49, 4, 16), P(49, 4, 16)), "type `Std::Int` cannot be assigned to type `9`"),
			},
		},
		"redeclare public identifier with a different type": {
			input: `
				var a: String
				var [a] = [1]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(28, 3, 10), P(28, 3, 10)), "type `Std::Int` cannot be assigned to type `Std::String`"),
			},
		},
		"private identifier pattern": {
			input: `
				var [_a] = ["", 8]
				var b: 9 = _a
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(39, 3, 16), P(40, 3, 17)), "type `Std::String | Std::Int` cannot be assigned to type `9`"),
			},
		},
		"list pattern with rest": {
			input: `
				var [a, *b] = ["", 8, 1]
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(45, 3, 16), P(45, 3, 16)), "type `Std::ArrayList[Std::String | Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"list pattern with invalid value": {
			input: `
				var [a] = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(11, 2, 11)), "type `8` cannot be matched against a list pattern"),
			},
		},
		"tuple pattern with rest": {
			input: `
				var %[a, *b] = ["", 8, 1]
				var c: 9 = b
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(46, 3, 16), P(46, 3, 16)), "type `Std::ArrayList[Std::String | Std::Int]` cannot be assigned to type `9`"),
			},
		},
		"tuple pattern with invalid value": {
			input: `
				var %[a] = 8
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(9, 2, 9), P(12, 2, 12)), "type `8` cannot be matched against a tuple pattern"),
			},
		},
		"pattern with as and new variable": {
			input: `
				var %[a as b] = [8]
				a = b
			`,
		},
		"pattern with as and existing variable": {
			input: `
				var b: String
				var %[a as b] = [8]
			`,
			err: error.ErrorList{
				error.NewFailure(L("<main>", P(34, 3, 16), P(34, 3, 16)), "type `Std::Int` cannot be assigned to type `Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
