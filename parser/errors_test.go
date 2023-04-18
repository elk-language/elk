package parser

import (
	"testing"

	"github.com/elk-language/elk/lexer"
	"github.com/google/go-cmp/cmp"
)

func TestErrorString(t *testing.T) {
	err := &Error{
		Message: "foo bar",
		Position: lexer.Position{
			StartByte:  3,
			ByteLength: 10,
			Line:       2,
			Column:     1,
		},
	}

	diff := cmp.Diff(err.String(), "2:1: foo bar")
	if diff != "" {
		t.Fatalf(diff)
	}
}

func TestErrorListAdd(t *testing.T) {
	got := ErrorList{
		&Error{
			Message: "foo bar",
			Position: lexer.Position{
				StartByte:  3,
				ByteLength: 10,
				Line:       2,
				Column:     1,
			},
		},
	}

	got.Add("sick style dude!", lexer.Position{
		StartByte:  25,
		ByteLength: 3,
		Line:       4,
		Column:     5,
	})

	want := ErrorList{
		&Error{
			Message: "foo bar",
			Position: lexer.Position{
				StartByte:  3,
				ByteLength: 10,
				Line:       2,
				Column:     1,
			},
		},
		&Error{
			Message: "sick style dude!",
			Position: lexer.Position{
				StartByte:  25,
				ByteLength: 3,
				Line:       4,
				Column:     5,
			},
		},
	}

	diff := cmp.Diff(got, want)
	if diff != "" {
		t.Fatalf(diff)
	}
}

func TestErrorListError(t *testing.T) {
	err := ErrorList{
		&Error{
			Message: "foo bar",
			Position: lexer.Position{
				StartByte:  3,
				ByteLength: 10,
				Line:       2,
				Column:     1,
			},
		},
		&Error{
			Message: "sick style dude!",
			Position: lexer.Position{
				StartByte:  25,
				ByteLength: 3,
				Line:       4,
				Column:     5,
			},
		},
	}

	got := err.Error()
	want := "2:1: foo bar\n4:5: sick style dude!\n"

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(diff)
	}
}
