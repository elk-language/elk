package vm_test

import (
	"testing"

	"github.com/elk-language/elk/value"
)

func TestVMSource_Closure(t *testing.T) {
	tests := sourceTestTable{
		"variable from outer context": {
			source: `
				x := "outside"
				func := -> println(x)
				func.()
			`,
			wantStdout:   "outside\n",
			wantStackTop: value.Nil,
		},
		"upvalue through 2 frames": {
			source: `
				func1 := ->
					x := "outside"
					func2 := ->
						func3 := -> println(x)
						func3.()
					end
					func2.()
				end
				func1.()
			`,
			wantStdout:   "outside\n",
			wantStackTop: value.Nil,
		},
		"closed upvalue": {
			source: `
				outer := ->
					x := "outside"
					::Std::Pair(-> println(x), |v| -> x = v)
				end
				var %[show, set] = outer.()
				show.()
				set.("new value")
				show.()
			`,
			wantStdout:   "outside\nnew value\n",
			wantStackTop: value.Nil,
		},
		"upvalue in for in loop": {
			source: `
				closures := []
				for i in 1...2
					closures << -> println i
				end
				closures[0].()
				closures[1].()
			`,
			wantStdout:   "1\n2\n",
			wantStackTop: value.Nil,
		},
		"upvalue in modifier for in loop": {
			source: `
				closures := []
				(closures << -> println i) for i in 1...2
				closures[0].()
				closures[1].()
			`,
			wantStdout:   "1\n2\n",
			wantStackTop: value.Nil,
		},
		"upvalue in fornum loop": {
			source: `
				closures := []
				fornum i := 1; i <= 2; i++
					closures << -> println i
				end
				closures[0].()
				closures[1].()
			`,
			wantStdout:   "3\n3\n",
			wantStackTop: value.Nil,
		},
		"upvalue in body of fornum loop": {
			source: `
				closures := []
				fornum i := 1; i <= 2; i++
					j := i
					closures << -> println j
				end
				closures[0].()
				closures[1].()
			`,
			wantStdout:   "1\n2\n",
			wantStackTop: value.Nil,
		},
		"upvalue in while loop": {
			source: `
				closures := []
				i := 1
				while i <= 2
					j := i
					closures << -> println j
					i++
				end
				closures[0].()
				closures[1].()
			`,
			wantStdout:   "1\n2\n",
			wantStackTop: value.Nil,
		},
		"upvalue in do while loop": {
			source: `
				closures := []
				i := 1
				do
					j := i
					closures << -> println j
					i++
				end while i <= 2
				closures[0].()
				closures[1].()
			`,
			wantStdout:   "1\n2\n",
			wantStackTop: value.Nil,
		},
		"upvalue in loop": {
			source: `
				closures := []
				i := 1
				loop
					break if i > 2
					j := i
					closures << -> println j
					i++
				end
				closures[0].()
				closures[1].()
			`,
			wantStdout:   "1\n2\n",
			wantStackTop: value.Nil,
		},
		"upvalue in until loop": {
			source: `
				closures := []
				i := 1
				until i > 2
					j := i
					closures << -> println j
					i++
				end
				closures[0].()
				closures[1].()
			`,
			wantStdout:   "1\n2\n",
			wantStackTop: value.Nil,
		},
		"upvalue in do until loop": {
			source: `
				closures := []
				i := 1
				do
					j := i
					closures << -> println j
					i++
				end until i > 2
				closures[0].()
				closures[1].()
			`,
			wantStdout:   "1\n2\n",
			wantStackTop: value.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
