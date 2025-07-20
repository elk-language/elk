package ds

import (
	"iter"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReverseSeq(t *testing.T) {
	tests := map[string]struct {
		seq  iter.Seq[string]
		want []string
	}{
		"reverse elements": {
			seq: func(yield func(string) bool) {
				yield("foo")
				yield("bar")
				yield("baz")
			},
			want: []string{
				"baz",
				"bar",
				"foo",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got []string
			for v := range ReverseSeq(tc.seq) {
				got = append(got, v)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestReverseSeq2(t *testing.T) {
	tests := map[string]struct {
		seq  iter.Seq2[string, int]
		want []Pair[string, int]
	}{
		"reverse elements": {
			seq: func(yield func(string, int) bool) {
				yield("foo", 5)
				yield("bar", -13)
				yield("baz", 92)
			},
			want: []Pair[string, int]{
				MakePair("baz", 92),
				MakePair("bar", -13),
				MakePair("foo", 5),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got []Pair[string, int]
			for k, v := range ReverseSeq2(tc.seq) {
				got = append(got, MakePair(k, v))
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
