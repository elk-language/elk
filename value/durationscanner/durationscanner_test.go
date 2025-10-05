package durationscanner_test

import (
	"testing"

	"github.com/elk-language/elk/value/durationscanner"
	"github.com/google/go-cmp/cmp"
)

type tokenValue struct {
	Token durationscanner.Token
	Value string
}

func T(token durationscanner.Token) tokenValue {
	return tokenValue{
		Token: token,
	}
}

func V(token durationscanner.Token, value string) tokenValue {
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
	scanner := durationscanner.New(tc.input)
	var got []tokenValue
	for {
		tok, value := scanner.Next()
		if tok == durationscanner.END_OF_FILE {
			break
		}
		got = append(got, V(tok, value))
	}
	diff := cmp.Diff(tc.want, got)
	if diff != "" {
		t.Fatal(diff)
	}
}

func TestDurationscanner(t *testing.T) {
	tests := testTable{
		"empty string": {
			input: "",
			want:  nil,
		},
		"only whitespace": {
			input: "   \n \t\r    ",
			want:  nil,
		},

		"date": {
			input: "200Y4M21D",
			want: []tokenValue{
				V(durationscanner.YEARS, "200"),
				V(durationscanner.MONTHS, "4"),
				V(durationscanner.DAYS, "21"),
			},
		},
		"date with spacing": {
			input: "200Y 4M 21D",
			want: []tokenValue{
				V(durationscanner.YEARS, "200"),
				V(durationscanner.MONTHS, "4"),
				V(durationscanner.DAYS, "21"),
			},
		},

		"time": {
			input: "23h40m15s300ms400us50ns",
			want: []tokenValue{

				V(durationscanner.HOURS, "23"),
				V(durationscanner.MINUTES, "40"),
				V(durationscanner.SECONDS, "15"),
				V(durationscanner.MILLISECONDS, "300"),
				V(durationscanner.MICROSECONDS, "400"),
				V(durationscanner.NANOSECONDS, "50"),
			},
		},
		"time with spacing": {
			input: "23h 40m 15s 300ms 400us 50ns",
			want: []tokenValue{
				V(durationscanner.HOURS, "23"),
				V(durationscanner.MINUTES, "40"),
				V(durationscanner.SECONDS, "15"),
				V(durationscanner.MILLISECONDS, "300"),
				V(durationscanner.MICROSECONDS, "400"),
				V(durationscanner.NANOSECONDS, "50"),
			},
		},

		"datetime": {
			input: "200Y4M21D23h40m15s300ms400us50ns",
			want: []tokenValue{
				V(durationscanner.YEARS, "200"),
				V(durationscanner.MONTHS, "4"),
				V(durationscanner.DAYS, "21"),
				V(durationscanner.HOURS, "23"),
				V(durationscanner.MINUTES, "40"),
				V(durationscanner.SECONDS, "15"),
				V(durationscanner.MILLISECONDS, "300"),
				V(durationscanner.MICROSECONDS, "400"),
				V(durationscanner.NANOSECONDS, "50"),
			},
		},
		"datetime with spacing": {
			input: "200Y 4M 21D 23h 40m 15s 300ms 400us 50ns",
			want: []tokenValue{
				V(durationscanner.YEARS, "200"),
				V(durationscanner.MONTHS, "4"),
				V(durationscanner.DAYS, "21"),
				V(durationscanner.HOURS, "23"),
				V(durationscanner.MINUTES, "40"),
				V(durationscanner.SECONDS, "15"),
				V(durationscanner.MILLISECONDS, "300"),
				V(durationscanner.MICROSECONDS, "400"),
				V(durationscanner.NANOSECONDS, "50"),
			},
		},

		"integer years": {
			input: "200Y",
			want: []tokenValue{
				V(durationscanner.YEARS, "200"),
			},
		},
		"integer years with underscores": {
			input: "2_000Y",
			want: []tokenValue{
				V(durationscanner.YEARS, "2000"),
			},
		},
		"float years": {
			input: "22.910Y",
			want: []tokenValue{
				V(durationscanner.YEARS, "22.910"),
			},
		},
		"float years leading dot": {
			input: ".910Y",
			want: []tokenValue{
				V(durationscanner.YEARS, "0.910"),
			},
		},
		"float years with underscores": {
			input: "2_2.9_10Y",
			want: []tokenValue{
				V(durationscanner.YEARS, "22.910"),
			},
		},

		"integer months": {
			input: "200M",
			want: []tokenValue{
				V(durationscanner.MONTHS, "200"),
			},
		},
		"integer months with underscores": {
			input: "2_000M",
			want: []tokenValue{
				V(durationscanner.MONTHS, "2000"),
			},
		},
		"float months": {
			input: "22.910M",
			want: []tokenValue{
				V(durationscanner.MONTHS, "22.910"),
			},
		},
		"float months leading dot": {
			input: ".910M",
			want: []tokenValue{
				V(durationscanner.MONTHS, "0.910"),
			},
		},
		"float months with underscores": {
			input: "2_2.9_10M",
			want: []tokenValue{
				V(durationscanner.MONTHS, "22.910"),
			},
		},

		"integer days": {
			input: "200D",
			want: []tokenValue{
				V(durationscanner.DAYS, "200"),
			},
		},
		"integer days with underscores": {
			input: "2_000D",
			want: []tokenValue{
				V(durationscanner.DAYS, "2000"),
			},
		},
		"float days": {
			input: "22.910D",
			want: []tokenValue{
				V(durationscanner.DAYS, "22.910"),
			},
		},
		"float days leading dot": {
			input: ".910D",
			want: []tokenValue{
				V(durationscanner.DAYS, "0.910"),
			},
		},
		"float days with underscores": {
			input: "2_2.9_10D",
			want: []tokenValue{
				V(durationscanner.DAYS, "22.910"),
			},
		},

		"integer hours": {
			input: "200h",
			want: []tokenValue{
				V(durationscanner.HOURS, "200"),
			},
		},
		"integer hours with underscores": {
			input: "2_000h",
			want: []tokenValue{
				V(durationscanner.HOURS, "2000"),
			},
		},
		"float hours": {
			input: "22.910h",
			want: []tokenValue{
				V(durationscanner.HOURS, "22.910"),
			},
		},
		"float hours leading dot": {
			input: ".910h",
			want: []tokenValue{
				V(durationscanner.HOURS, "0.910"),
			},
		},
		"float hours with underscores": {
			input: "2_2.9_10h",
			want: []tokenValue{
				V(durationscanner.HOURS, "22.910"),
			},
		},

		"integer minutes": {
			input: "200m",
			want: []tokenValue{
				V(durationscanner.MINUTES, "200"),
			},
		},
		"integer minutes with underscores": {
			input: "2_000m",
			want: []tokenValue{
				V(durationscanner.MINUTES, "2000"),
			},
		},
		"float minutes": {
			input: "22.910m",
			want: []tokenValue{
				V(durationscanner.MINUTES, "22.910"),
			},
		},
		"float minutes leading dot": {
			input: ".910m",
			want: []tokenValue{
				V(durationscanner.MINUTES, "0.910"),
			},
		},
		"float minutes with underscores": {
			input: "2_2.9_10m",
			want: []tokenValue{
				V(durationscanner.MINUTES, "22.910"),
			},
		},

		"integer seconds": {
			input: "200s",
			want: []tokenValue{
				V(durationscanner.SECONDS, "200"),
			},
		},
		"integer seconds with underscores": {
			input: "2_000s",
			want: []tokenValue{
				V(durationscanner.SECONDS, "2000"),
			},
		},
		"float seconds": {
			input: "22.910s",
			want: []tokenValue{
				V(durationscanner.SECONDS, "22.910"),
			},
		},
		"float seconds leading dot": {
			input: ".910s",
			want: []tokenValue{
				V(durationscanner.SECONDS, "0.910"),
			},
		},
		"float seconds with underscores": {
			input: "2_2.9_10s",
			want: []tokenValue{
				V(durationscanner.SECONDS, "22.910"),
			},
		},

		"integer milliseconds with bad second char": {
			input: "200ml",
			want: []tokenValue{
				V(durationscanner.MINUTES, "200"),
				V(durationscanner.ERROR, "unexpected char 'l', expected a digit"),
			},
		},
		"integer milliseconds": {
			input: "200ms",
			want: []tokenValue{
				V(durationscanner.MILLISECONDS, "200"),
			},
		},
		"integer milliseconds with underscores": {
			input: "2_000ms",
			want: []tokenValue{
				V(durationscanner.MILLISECONDS, "2000"),
			},
		},
		"float milliseconds": {
			input: "22.910ms",
			want: []tokenValue{
				V(durationscanner.MILLISECONDS, "22.910"),
			},
		},
		"float milliseconds leading dot": {
			input: ".910ms",
			want: []tokenValue{
				V(durationscanner.MILLISECONDS, "0.910"),
			},
		},
		"float milliseconds with underscores": {
			input: "2_2.9_10ms",
			want: []tokenValue{
				V(durationscanner.MILLISECONDS, "22.910"),
			},
		},

		"integer microseconds with bad second char": {
			input: "200ul",
			want: []tokenValue{
				V(durationscanner.ERROR, "unexpected char 'l', expected 's'"),
			},
		},
		"integer microseconds": {
			input: "200us",
			want: []tokenValue{
				V(durationscanner.MICROSECONDS, "200"),
			},
		},
		"integer microseconds with greek": {
			input: "600µs",
			want: []tokenValue{
				V(durationscanner.MICROSECONDS, "600"),
			},
		},
		"integer microseconds with underscores": {
			input: "2_000us",
			want: []tokenValue{
				V(durationscanner.MICROSECONDS, "2000"),
			},
		},
		"float microseconds": {
			input: "22.910us",
			want: []tokenValue{
				V(durationscanner.MICROSECONDS, "22.910"),
			},
		},
		"float microseconds leading dot": {
			input: ".910us",
			want: []tokenValue{
				V(durationscanner.MICROSECONDS, "0.910"),
			},
		},
		"float microseconds with underscores": {
			input: "2_2.9_10us",
			want: []tokenValue{
				V(durationscanner.MICROSECONDS, "22.910"),
			},
		},

		"integer nanoseconds with bad second char": {
			input: "200nl",
			want: []tokenValue{
				V(durationscanner.ERROR, "unexpected char 'l', expected 's'"),
			},
		},
		"integer nanoseconds": {
			input: "200ns",
			want: []tokenValue{
				V(durationscanner.NANOSECONDS, "200"),
			},
		},
		"integer nanoseconds with underscores": {
			input: "2_000ns",
			want: []tokenValue{
				V(durationscanner.NANOSECONDS, "2000"),
			},
		},
		"float nanoseconds": {
			input: "22.910ns",
			want: []tokenValue{
				V(durationscanner.NANOSECONDS, "22.910"),
			},
		},
		"float nanoseconds leading dot": {
			input: ".910ns",
			want: []tokenValue{
				V(durationscanner.NANOSECONDS, "0.910"),
			},
		},
		"float nanoseconds with underscores": {
			input: "2_2.9_10ns",
			want: []tokenValue{
				V(durationscanner.NANOSECONDS, "22.910"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokenTest(tc, t)
		})
	}
}
