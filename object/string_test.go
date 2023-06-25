package object

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStringByteLength(t *testing.T) {
	tests := map[string]struct {
		str  String
		want int
	}{
		"only ascii": {
			str:  "foo123",
			want: 6,
		},
		"unicode": {
			str:  "ślązak",
			want: 8,
		},
		"grapheme clusters": {
			str:  "👩‍🚀🇵🇱",
			want: 19,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.ByteLength()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringCharLength(t *testing.T) {
	tests := map[string]struct {
		str  String
		want int
	}{
		"only ascii": {
			str:  "foo123",
			want: 6,
		},
		"unicode": {
			str:  "ślązak",
			want: 6,
		},
		"grapheme clusters": {
			str:  "👩‍🚀🇵🇱",
			want: 5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.CharLength()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringGraphemeLength(t *testing.T) {
	tests := map[string]struct {
		str  String
		want int
	}{
		"only ascii": {
			str:  "foo123",
			want: 6,
		},
		"unicode": {
			str:  "ślązak",
			want: 6,
		},
		"graphemes clusters": {
			str:  "👩‍🚀🇵🇱",
			want: 2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.GraphemeLength()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringReverseBytes(t *testing.T) {
	tests := map[string]struct {
		str  String
		want String
	}{
		"only ascii": {
			str:  "foo123",
			want: "321oof",
		},
		"unicode": {
			str:  "ślązak",
			want: "kaz\x85\xc4l\x9b\xc5",
		},
		"grapheme clusters": {
			str:  "👩‍🚀🇵🇱",
			want: "\xb1\x87\x9f\U000351df\xf0\x80\x9a\x9f\xf0\x8d\x80⩑\x9f\xf0",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.ReverseBytes()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringReverseChars(t *testing.T) {
	tests := map[string]struct {
		str  String
		want String
	}{
		"only ascii": {
			str:  "foo123",
			want: "321oof",
		},
		"unicode": {
			str:  "ślązak",
			want: "kaząlś",
		},
		"grapheme clusters": {
			str:  "👩‍🚀🇵🇱",
			want: "🇱🇵🚀\u200d👩",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.ReverseChars()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestStringReverseGraphemes(t *testing.T) {
	tests := map[string]struct {
		str  String
		want String
	}{
		"only ascii": {
			str:  "foo123",
			want: "321oof",
		},
		"unicode": {
			str:  "ślązak",
			want: "kaząlś",
		},
		"grapheme clusters": {
			str:  "👩‍🚀🇵🇱",
			want: "🇵🇱👩‍🚀",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.str.ReverseGraphemes()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
