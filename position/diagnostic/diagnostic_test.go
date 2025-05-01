package diagnostic

import (
	"testing"

	"github.com/elk-language/elk/position"
	"github.com/google/go-cmp/cmp"
)

// Create a new location in tests
var L = position.NewLocation
var LP = position.NewLocationWithParent

// Create a new position in tests
var P = position.New

// Create a new span in tests
var S = position.NewSpan

func TestDiagnosticString(t *testing.T) {
	err := NewDiagnostic(
		L("/opt/elk", S(P(0, 2, 1), P(5, 2, 1))),
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
			L("/opt/elk", S(P(0, 1, 1), P(5, 2, 1))),
			"foo bar",
			FAIL,
		),
	}

	got.Add("sick style dude!", L("/opt/elk", S(P(6, 2, 2), P(10, 2, 6))), FAIL)

	want := DiagnosticList{
		NewDiagnostic(
			L("/opt/elk", S(P(0, 1, 1), P(5, 2, 1))),
			"foo bar",
			FAIL,
		),
		NewDiagnostic(
			L("/opt/elk", S(P(6, 2, 2), P(10, 2, 6))),
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
			L("/some/path", S(P(5, 2, 1), P(5, 2, 1))),
			"foo bar",
			FAIL,
		),
		NewDiagnostic(
			L("<main>", S(P(20, 4, 5), P(25, 4, 10))),
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

func TestDiagnosticList_HumanStringWithSource(t *testing.T) {
	list := DiagnosticList{
		NewDiagnostic(
			L("fixtures/if.elk", S(P(3, 1, 4), P(7, 1, 8))),
			"undefined variable 'style'",
			FAIL,
		),
		NewDiagnostic(
			LP(
				"fixtures/if.elk",
				S(P(9, 1, 10), P(10, 1, 11)),
				L("fixtures/if.elk", S(P(0, 1, 1), P(1, 1, 2))),
			),
			"invalid operator '=='",
			FAIL,
		),
	}

	got, err := list.HumanString(false, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	want := `[FAIL] undefined variable 'style'

  fixtures/if.elk:1:4
    1 | if style == FOO
           └───┤
               └ Here

[FAIL] invalid operator '=='

  fixtures/if.elk:1:10
    1 | if style == FOO
                 └┤
                  └ Here

  fixtures/if.elk:1:1
    1 | if style == FOO
        └┤
         └ Here

`

	if diff := cmp.Diff(want, got); diff != "" {
		t.Log(got)
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
					L("/some/path", S(P(0, 1, 1), P(5, 2, 1))),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", S(P(18, 2, 3), P(20, 4, 5))),
					"sick style dude!",
					FAIL,
				),
			},
			right: nil,
			want: DiagnosticList{
				NewDiagnostic(
					L("/some/path", S(P(0, 1, 1), P(5, 2, 1))),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", S(P(18, 2, 3), P(20, 4, 5))),
					"sick style dude!",
					FAIL,
				),
			},
		},
		"return right when left is nil": {
			left: nil,
			right: DiagnosticList{
				NewDiagnostic(
					L("/some/path", S(P(0, 1, 1), P(5, 2, 1))),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", S(P(18, 2, 3), P(20, 4, 5))),
					"sick style dude!",
					FAIL,
				),
			},
			want: DiagnosticList{
				NewDiagnostic(
					L("/some/path", S(P(0, 1, 1), P(5, 2, 1))),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", S(P(18, 2, 3), P(20, 4, 5))),
					"sick style dude!",
					FAIL,
				),
			},
		},
		"return joined list": {
			left: DiagnosticList{
				NewDiagnostic(
					L("/some/path", S(P(0, 1, 1), P(5, 2, 1))),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", S(P(18, 2, 3), P(20, 4, 5))),
					"sick style dude!",
					FAIL,
				),
			},
			right: DiagnosticList{
				NewDiagnostic(
					L("/foo/bar", S(P(50, 10, 2), P(51, 10, 3))),
					"baz",
					FAIL,
				),
			},
			want: DiagnosticList{
				NewDiagnostic(
					L("/some/path", S(P(0, 1, 1), P(5, 2, 1))),
					"foo bar",
					FAIL,
				),
				NewDiagnostic(
					L("<main>", S(P(18, 2, 3), P(20, 4, 5))),
					"sick style dude!",
					FAIL,
				),
				NewDiagnostic(
					L("/foo/bar", S(P(50, 10, 2), P(51, 10, 3))),
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
