package position

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestErrorString(t *testing.T) {
	err := NewError(
		NewLocation("/opt/elk", 0, 0, 2, 1),
		"foo bar",
	)

	diff := cmp.Diff(err.String(), "/opt/elk:2:1: foo bar")
	if diff != "" {
		t.Fatalf(diff)
	}
}

func TestErrorListAdd(t *testing.T) {
	got := ErrorList{
		NewError(
			NewLocation("", 0, 0, 2, 1),
			"foo bar",
		),
	}

	got.Add("sick style dude!", NewLocation("", 0, 0, 4, 5))

	want := ErrorList{
		NewError(
			NewLocation("", 0, 0, 2, 1),
			"foo bar",
		),
		NewError(
			NewLocation("", 0, 0, 4, 5),
			"sick style dude!",
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
			NewLocation("/some/path", 0, 0, 2, 1),
			"foo bar",
		),
		NewError(
			NewLocation("main", 0, 0, 4, 5),
			"sick style dude!",
		),
	}

	got := err.Error()
	want := "/some/path:2:1: foo bar\nmain:4:5: sick style dude!\n"

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
					NewLocation("/some/path", 0, 0, 2, 1),
					"foo bar",
				),
				NewError(
					NewLocation("main", 0, 0, 4, 5),
					"sick style dude!",
				),
			},
			right: nil,
			want: ErrorList{
				NewError(
					NewLocation("/some/path", 0, 0, 2, 1),
					"foo bar",
				),
				NewError(
					NewLocation("main", 0, 0, 4, 5),
					"sick style dude!",
				),
			},
		},
		"return right when left is nil": {
			left: nil,
			right: ErrorList{
				NewError(
					NewLocation("/some/path", 0, 0, 2, 1),
					"foo bar",
				),
				NewError(
					NewLocation("main", 0, 0, 4, 5),
					"sick style dude!",
				),
			},
			want: ErrorList{
				NewError(
					NewLocation("/some/path", 0, 0, 2, 1),
					"foo bar",
				),
				NewError(
					NewLocation("main", 0, 0, 4, 5),
					"sick style dude!",
				),
			},
		},
		"return joined list": {
			left: ErrorList{
				NewError(
					NewLocation("/some/path", 0, 0, 2, 1),
					"foo bar",
				),
				NewError(
					NewLocation("main", 0, 0, 4, 5),
					"sick style dude!",
				),
			},
			right: ErrorList{
				NewError(
					NewLocation("/foo/bar", 5, 8, 3, 2),
					"some new error",
				),
			},
			want: ErrorList{
				NewError(
					NewLocation("/some/path", 0, 0, 2, 1),
					"foo bar",
				),
				NewError(
					NewLocation("main", 0, 0, 4, 5),
					"sick style dude!",
				),
				NewError(
					NewLocation("/foo/bar", 5, 8, 3, 2),
					"some new error",
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
