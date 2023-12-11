package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestSymbolInspect(t *testing.T) {
	tests := map[string]struct {
		sym  value.Symbol
		want string
	}{
		"only letters": {
			sym:  value.ToSymbol("foo"),
			want: `:foo`,
		},
		"digit as the first char": {
			sym:  value.ToSymbol("1foo"),
			want: `:"1foo"`,
		},
		"with underscores": {
			sym:  value.ToSymbol("foo_bar"),
			want: `:foo_bar`,
		},
		"with an initial letter and digits": {
			sym:  value.ToSymbol("foo1"),
			want: `:foo1`,
		},
		"with one byte escapes": {
			sym:  value.ToSymbol("foo\nbar\t\r\v\f\a\b"),
			want: `:"foo\nbar\t\r\v\f\a\b"`,
		},
		"with non-ascii bytes": {
			sym:  value.ToSymbol("foo\x02bar"),
			want: `:"foo\x02bar"`,
		},
		"with non-graphic unicode chars": {
			sym:  value.ToSymbol("foo\U0010FFFFbar"),
			want: `:"foo\U0010FFFFbar"`,
		},
		"with emojis": {
			sym:  value.ToSymbol("foo🐧bar"),
			want: `:"foo🐧bar"`,
		},
		"with spaces": {
			sym:  value.ToSymbol("foo bar"),
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
