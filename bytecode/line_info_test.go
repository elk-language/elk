package bytecode

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLineInfoListLast(t *testing.T) {
	tests := map[string]struct {
		in   LineInfoList
		want *LineInfo
	}{
		"return the last element when there are a few": {
			in:   LineInfoList{NewLineInfo(1, 3), NewLineInfo(2, 1)},
			want: NewLineInfo(2, 1),
		},
		"return nil when the list is empty": {
			in:   LineInfoList{},
			want: nil,
		},
		"return nil when the list is nil": {
			in:   nil,
			want: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.in.Last()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestLineInfoListGetLineNumber(t *testing.T) {
	tests := map[string]struct {
		list  LineInfoList
		index int
		want  int
	}{
		"return the first element": {
			list:  LineInfoList{NewLineInfo(1, 3), NewLineInfo(2, 1)},
			index: 2,
			want:  1,
		},
		"return the second element": {
			list:  LineInfoList{NewLineInfo(1, 3), NewLineInfo(2, 1), NewLineInfo(3, 5)},
			index: 3,
			want:  2,
		},
		"return -1 when no such index": {
			list:  LineInfoList{NewLineInfo(1, 3), NewLineInfo(2, 1), NewLineInfo(3, 5)},
			index: 30,
			want:  -1,
		},
		"return -1 when the list is empty": {
			list:  LineInfoList{},
			index: 6,
			want:  -1,
		},
		"return nil when the list is nil": {
			list:  nil,
			index: 6,
			want:  -1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.list.GetLineNumber(tc.index)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestLineInfoListAddLineNumber(t *testing.T) {
	tests := map[string]struct {
		list LineInfoList
		line int
		want LineInfoList
	}{
		"increment the instruction count when the line is the same": {
			list: LineInfoList{NewLineInfo(1, 3), NewLineInfo(2, 1)},
			line: 2,
			want: LineInfoList{NewLineInfo(1, 3), NewLineInfo(2, 2)},
		},
		"add a new line info when the line is different": {
			list: LineInfoList{NewLineInfo(1, 3), NewLineInfo(2, 1)},
			line: 3,
			want: LineInfoList{NewLineInfo(1, 3), NewLineInfo(2, 1), NewLineInfo(3, 1)},
		},
		"add a new line info when the list is empty": {
			list: LineInfoList{},
			line: 3,
			want: LineInfoList{NewLineInfo(3, 1)},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.list.AddLineNumber(tc.line, 1)
			got := tc.list
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
