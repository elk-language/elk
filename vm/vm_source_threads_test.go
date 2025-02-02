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
				using Std::Sync::WaitGroup

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
		"synchronise threads with a mutex": {
			source: `
				using Std::Sync::{Mutex, WaitGroup}

				class Counter
					attr n: Int, m: Mutex

					init
						@n = 0
						@m = Mutex()
					end

					def incr
						@m.lock
						@n++
						@m.unlock
					end
				end

				def work(c: Counter, wg: WaitGroup)
					for i in 1...50
						c.incr
					end

					wg.end
				end

				c := Counter()
				wg := WaitGroup(5)

				for i in 1...5
					go work(c, wg)
				end

				wg.wait
				println "counter: ${c.n}"
			`,
			wantStdout:   "counter: 250\n",
			wantStackTop: value.Nil,
		},
		"synchronise threads with a channel": {
			source: `
				using Std::Sync::WaitGroup
				ch := Channel::[Int]()
				wg := WaitGroup(2)
				go
					ch << 5
					ch << 10
					ch << 2
					ch.close
					wg.end
				end

				go
					for i in ch
						println i.inspect
					end
					wg.end
				end

				wg.wait
			`,
			wantStdout:   "5\n10\n2\n",
			wantStackTop: value.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			vmSourceTest(tc, t)
		})
	}
}
