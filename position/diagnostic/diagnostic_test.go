package diagnostic

import (
	"testing"

	"github.com/elk-language/elk/position"
	"github.com/google/go-cmp/cmp"
)

// Create a new location in tests
var L = position.NewLocation

// Create a new position in tests
var P = position.New

func TestDiagnosticString(t *testing.T) {
	err := NewDiagnostic(
		L("/opt/elk", P(0, 2, 1), P(5, 2, 1)),
		"foo bar",
		FAIL,
	)

	diff := cmp.Diff(err.String(), "/opt/elk:2:1: foo bar")
	if diff != "" {
		t.Fatal(diff)
	}
}

func TestDiagnosticListAdd(t *testing.T) {
	got := DiagnosticList{
		NewDiagnostic(
			L("/opt/elk", P(0, 1, 1), P(5, 2, 1)),
			"foo bar",
			FAIL,
		),
	}

	got.Add("sick style dude!", L("/opt/elk", P(6, 2, 2), P(10, 2, 6)), FAIL)

	want := DiagnosticList{
		NewDiagnostic(
			L("/opt/elk", P(0, 1, 1), P(5, 2, 1)),
			"foo bar",
			FAIL,
		),
		NewDiagnostic(
			L("/opt/elk", P(6, 2, 2), P(10, 2, 6)),
			"sick style dude!",
			FAIL,
		),
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Fatal(diff)
	}
}

func TestDiagnosticListError(t *testing.T) {
	err := DiagnosticList{
		NewDiagnostic(
			L("/some/path", P(5, 2, 1), P(5, 2, 1)),
			"foo bar",
			FAIL,
		),
		NewDiagnostic(
			L("<main>", P(20, 4, 5), P(25, 4, 10)),
			"sick style dude!",
			FAIL,
		),
	}

	got := err.Error()
	want := "/some/path:2:1: foo bar\n<main>:4:5: sick style dude!\n"

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatal(diff)
	}
}

func TestDiagnosticListJoin(t *testing.T) {
	tests := map[string]struct {
		left  DiagnosticList
		right DiagnosticList
		want  DiagnosticList
	}{
		"return left when right is nil": {
			left: DiagnosticList{
				NewDiagnostic(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAIL,
				),
			},
			right: nil,
			want: DiagnosticList{
				NewDiagnostic(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAIL,
				),
			},
		},
		"return right when left is nil": {
			left: nil,
			right: DiagnosticList{
				NewDiagnostic(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAIL,
				),
			},
			want: DiagnosticList{
				NewDiagnostic(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAIL,
				),
			},
		},
		"return joined list": {
			left: DiagnosticList{
				NewDiagnostic(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAIL,
				),
			},
			right: DiagnosticList{
				NewDiagnostic(
					L("/foo/bar", P(50, 10, 2), P(51, 10, 3)),
					"baz",
					FAIL,
				),
			},
			want: DiagnosticList{
				NewDiagnostic(
					L("/some/path", P(0, 1, 1), P(5, 2, 1)),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", P(18, 2, 3), P(20, 4, 5)),
					"sick style dude!",
					FAIL,
				),
				NewDiagnostic(
					L("/foo/bar", P(50, 10, 2), P(51, 10, 3)),
					"baz",
					FAIL,
				),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.left.Join(tc.right)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
