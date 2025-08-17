package checker_test

import (
	"testing"
)

func TestAny(t *testing.T) {
	tests := testTable{
		"any is assignable to any": {
			input: `
				var foo: any = 5
				var bar: any = foo
			`,
		},
		"union types are assignable to any": {
			input: `
				var foo: String | Int = 5
				var bar: any = foo
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
