package timescanner_test

import (
	"testing"

	"github.com/elk-language/elk/value/timescanner"
	"github.com/google/go-cmp/cmp"
)

type tokenValue struct {
	Token timescanner.Token
	Value string
}

func T(token timescanner.Token) tokenValue {
	return tokenValue{
		Token: token,
	}
}

func V(token timescanner.Token, value string) tokenValue {
	return tokenValue{
		Token: token,
		Value: value,
	}
}

// Represents a single test case.
type testCase struct {
	input string
	want  []tokenValue
}

// Type of the test table.
type testTable map[string]testCase

// Function which powers all timescanner tests.
// Inspects if the produced stream of tokens
// matches the expected one.
func tokenTest(tc testCase, t *testing.T) {
	t.Helper()
	scanner := timescanner.New(tc.input)
	var got []tokenValue
	for {
		tok, value := scanner.Next()
		if tok == timescanner.END_OF_FILE {
			break
		}
		got = append(got, V(tok, value))
	}
	diff := cmp.Diff(tc.want, got)
	if diff != "" {
		t.Fatal(diff)
	}
}

func TestTimescanner(t *testing.T) {
	tests := testTable{
		"empty format string": {
			input: "",
			want:  nil,
		},
		"no format options": {
			input: "some cool text",
			want: []tokenValue{
				V(timescanner.TEXT, "some cool text"),
			},
		},
		"only format directives": {
			input: "%G%V%u%H%M%S%z",
			want: []tokenValue{
				T(timescanner.FULL_YEAR_WEEK_BASED_ZERO_PADDED),
				T(timescanner.WEEK_OF_WEEK_BASED_YEAR_ZERO_PADDED),
				T(timescanner.DAY_OF_WEEK_NUMBER),
				T(timescanner.HOUR_OF_DAY_ZERO_PADDED),
				T(timescanner.MINUTE_OF_HOUR_ZERO_PADDED),
				T(timescanner.SECOND_OF_MINUTE_ZERO_PADDED),
				T(timescanner.TIMEZONE_OFFSET),
			},
		},
		"date with text": {
			input: "%Y-%m-%d",
			want: []tokenValue{
				T(timescanner.FULL_YEAR_ZERO_PADDED),
				V(timescanner.TEXT, "-"),
				T(timescanner.MONTH_ZERO_PADDED),
				V(timescanner.TEXT, "-"),
				T(timescanner.DAY_OF_MONTH_ZERO_PADDED),
			},
		},
		"fake format option": {
			input: "%Y-%m-%!-%d",
			want: []tokenValue{
				T(timescanner.FULL_YEAR_ZERO_PADDED),
				V(timescanner.TEXT, "-"),
				T(timescanner.MONTH_ZERO_PADDED),
				V(timescanner.TEXT, "-"),
				V(timescanner.INVALID_FORMAT_DIRECTIVE, "%!"),
				V(timescanner.TEXT, "-"),
				T(timescanner.DAY_OF_MONTH_ZERO_PADDED),
			},
		},
		"modifiers": {
			input: "%-Y %_m %d %^B",
			want: []tokenValue{
				T(timescanner.FULL_YEAR),
				V(timescanner.TEXT, " "),
				T(timescanner.MONTH_SPACE_PADDED),
				V(timescanner.TEXT, " "),
				T(timescanner.DAY_OF_MONTH_ZERO_PADDED),
				V(timescanner.TEXT, " "),
				T(timescanner.MONTH_FULL_NAME_UPPERCASE),
			},
		},
		"add heavy text": {
			input: "Today is %G,%n of %V through %u%H,%tin %M:%S:%z",
			want: []tokenValue{
				V(timescanner.TEXT, "Today is "),
				T(timescanner.FULL_YEAR_WEEK_BASED_ZERO_PADDED),
				V(timescanner.TEXT, ","),
				T(timescanner.NEWLINE),
				V(timescanner.TEXT, " of "),
				T(timescanner.WEEK_OF_WEEK_BASED_YEAR_ZERO_PADDED),
				V(timescanner.TEXT, " through "),
				T(timescanner.DAY_OF_WEEK_NUMBER),
				T(timescanner.HOUR_OF_DAY_ZERO_PADDED),
				V(timescanner.TEXT, ","),
				T(timescanner.TAB),
				V(timescanner.TEXT, "in "),
				T(timescanner.MINUTE_OF_HOUR_ZERO_PADDED),
				V(timescanner.TEXT, ":"),
				T(timescanner.SECOND_OF_MINUTE_ZERO_PADDED),
				V(timescanner.TEXT, ":"),
				T(timescanner.TIMEZONE_OFFSET),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
