package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSymbolInspect(t *testing.T) {
	tests := map[string]struct {
		sym  Symbol
		want string
	}{
		"only letters": {
			sym:  SymbolTable.Add("foo"),
			want: `:foo`,
		},
		"digit as the first char": {
			sym:  SymbolTable.Add("1foo"),
			want: `:"1foo"`,
		},
		"with underscores": {
			sym:  SymbolTable.Add("foo_bar"),
			want: `:foo_bar`,
		},
		"with an initial letter and digits": {
			sym:  SymbolTable.Add("foo1"),
			want: `:foo1`,
		},
		"with one byte escapes": {
			sym:  SymbolTable.Add("foo\nbar\t\r\v\f\a\b"),
			want: `:"foo\nbar\t\r\v\f\a\b"`,
		},
		"with non-ascii bytes": {
			sym:  SymbolTable.Add("foo\x02bar"),
			want: `:"foo\x02bar"`,
		},
		"with non-graphic unicode chars": {
			sym:  SymbolTable.Add("foo\U0010FFFFbar"),
			want: `:"foo\U0010FFFFbar"`,
		},
		"with emojis": {
			sym:  SymbolTable.Add("foo🐧bar"),
			want: `:"foo🐧bar"`,
		},
		"with spaces": {
			sym:  SymbolTable.Add("foo bar"),
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
