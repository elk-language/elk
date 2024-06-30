package error

import (
	"testing"

	"github.com/elk-language/elk/position"
	"github.com/google/go-cmp/cmp"
)

// Create a new location in tests
var L = position.NewLocation

// Create a new position in tests
var P = position.New

func TestErrorString(t *testing.T) {
	err := NewError(
		L("/opt/elk", P(0, 2, 1), P(5, 2, 1)),
		"foo bar",
		FAILURE,
	)

	diff := cmp.Diff(err.String(), "/opt/elk:2:1: foo bar")
	if diff != "" {
		t.Fatalf(diff)
	}
}

func TestErrorListAdd(t *testing.T) {
	got := ErrorList{
		NewError(
			L("/opt/elk", P(0, 1, 1), P(5, 2, 1)),
			"foo bar",
			FAILURE,
		),
	}

	got.Add("sick style dude!", L("/opt/elk", P(6, 2, 2), P(10, 2, 6)), FAILURE)

	want := ErrorList{
		NewError(
			L("/opt/elk", P(0, 1, 1), P(5, 2, 1)),
			"foo bar",
			FAILURE,
		),
		NewError(
			L("/opt/elk", P(6, 2, 2), P(10, 2, 6)),
			"sick style dude!",
			FAILURE,
		),
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Fatalf(diff)
	}
}

func TestErrorListError(t *testing.T) {
	err := ErrorList{
		NewError(
			L("/some/path", P(5, 2, 1), P(5, 2, 1)),
			"foo bar",
			FAILURE,
		),
		NewError(
			L("<main>", P(20, 4, 5), P(25, 4, 10)),
			"sick style dude!",
			FAILURE,
		),
	}

	got := err.Error()
	want := "/some/path:2:1: foo bar\n<main>:4:5: sick style dude!\n"

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(diff)
	}
}

func TestErrorListJoin(t *testing.T) {
	tests := map[string]struct {
		left  ErrorList
		right ErrorList
		want  ErrorList
	}{
		"return left when right is nil": {
			left: ErrorList{
				NewError(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAILURE,
				),
				NewError(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAILURE,
				),
			},
			right: nil,
			want: ErrorList{
				NewError(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAILURE,
				),
				NewError(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAILURE,
				),
			},
		},
		"return right when left is nil": {
			left: nil,
			right: ErrorList{
				NewError(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAILURE,
				),
				NewError(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAILURE,
				),
			},
			want: ErrorList{
				NewError(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAILURE,
				),
				NewError(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAILURE,
				),
			},
		},
		"return joined list": {
			left: ErrorList{
				NewError(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAILURE,
				),
				NewError(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAILURE,
				),
			},
			right: ErrorList{
				NewError(
					L("/foo/bar", P(50, 10, 2), P(51, 10, 3)),
					"baz",
					FAILURE,
				),
			},
			want: ErrorList{
				NewError(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAILURE,
				),
				NewError(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAILURE,
				),
				NewError(
					L("/foo/bar", P(50, 10, 2), P(51, 10, 3)),
					"baz",
					FAILURE,
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Join(tc.right)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
