package object

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSymbolInspect(t *testing.T) {
	tests := map[string]struct {
		sym  *Symbol
		want string
	}{
		"only letters": {
			sym:  newSymbol("foo", 1),
			want: `:foo`,
		},
		"digit as the first char": {
			sym:  newSymbol("1foo", 1),
			want: `:"1foo"`,
		},
		"with underscores": {
			sym:  newSymbol("foo_bar", 1),
			want: `:foo_bar`,
		},
		"with an initial letter and digits": {
			sym:  newSymbol("foo1", 1),
			want: `:foo1`,
		},
		"with one byte escapes": {
			sym:  newSymbol("foo\nbar\t\r\v\f\a\b", 1),
			want: `:"foo\nbar\t\r\v\f\a\b"`,
		},
		"with non-ascii bytes": {
			sym:  newSymbol("foo\x02bar", 1),
			want: `:"foo\x02bar"`,
		},
		"with non-graphic unicode chars": {
			sym:  newSymbol("foo\U0010FFFFbar", 1),
			want: `:"foo\U0010FFFFbar"`,
		},
		"with emojis": {
			sym:  newSymbol("foo🐧bar", 1),
			want: `:"foo🐧bar"`,
		},
		"with spaces": {
			sym:  newSymbol("foo bar", 1),
			want: `:"foo bar"`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.sym.Inspect()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
