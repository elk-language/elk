package vm_test

import (
	"testing"

	"github.com/elk-language/elk/value"
)

func TestVMSource_Go(t *testing.T) {
	tests := sourceTestTable{
		"handle an error thrown in a separate coroutine": {
			source: `
				go throw unchecked 5
				sleep 0.5.seconds
			`,
			wantStderr:   "Stack trace (the most recent call is last)\n 0: sourceName:2, in `<closure>`\nError! Uncaught thrown value: 5\n\n",
			wantStackTop: value.Nil,
		},
		"wait for a thread with a WaitGroup": {
			source: `
				def print_numbers(name: String)
					for i in 1...5
						println "$name: $i"
					end
				end


				wg := WaitGroup(1)
				go
					print_numbers("go1")
					wg.end
				end

				wg.wait
				print_numbers("main1")
			`,
			wantStdout:   "go1: 1\ngo1: 2\ngo1: 3\ngo1: 4\ngo1: 5\nmain1: 1\nmain1: 2\nmain1: 3\nmain1: 4\nmain1: 5\n",
			wantStackTop: value.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
