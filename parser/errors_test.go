package parser

import (
	"testing"

	"github.com/elk-language/elk/lexer"
	"github.com/google/go-cmp/cmp"
)

func TestErrorString(t *testing.T) {
	err := NewError(
		lexer.NewPosition(0, 0, 2, 1),
		"foo bar",
	)

	diff := cmp.Diff(err.String(), "2:1: foo bar")
	if diff != "" {
		t.Fatalf(diff)
	}
}

func TestErrorStringWithPath(t *testing.T) {
	err := NewError(
		lexer.NewPosition(0, 0, 2, 1),
		"foo bar",
	)

	diff := cmp.Diff(err.StringWithPath("/some/path"), "/some/path:2:1: foo bar")
	if diff != "" {
		t.Fatalf(diff)
	}
}

func TestErrorListAdd(t *testing.T) {
	got := ErrorList{
		NewError(
			lexer.NewPosition(0, 0, 2, 1),
			"foo bar",
		),
	}

	got.Add("sick style dude!", lexer.NewPosition(0, 0, 4, 5))

	want := ErrorList{
		NewError(
			lexer.NewPosition(0, 0, 2, 1),
			"foo bar",
		),
		NewError(
			lexer.NewPosition(0, 0, 4, 5),
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
			lexer.NewPosition(0, 0, 2, 1),
			"foo bar",
		),
		NewError(
			lexer.NewPosition(0, 0, 4, 5),
			"sick style dude!",
		),
	}

	got := err.Error()
	want := "2:1: foo bar\n4:5: sick style dude!\n"

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(diff)
	}
}
