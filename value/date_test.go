package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestDateFormat(t *testing.T) {
	tests := map[string]struct {
		date   value.Date
		format string
		want   string
		err    value.Value
	}{
		"simple date with zero padding": {
			date:   value.MakeDate(2023, 5, 3),
			format: "%Y-%m-%d",
			want:   `2023-05-03`,
		},
		"simple date with space padding": {
			date:   value.MakeDate(520, 5, 3),
			format: "%_Y %_m %_d",
			want:   ` 520  5  3`,
		},
		"literal values": {
			date:   value.MakeDate(2023, 5, 3),
			format: "%Y %% %m %t %d %n",
			want:   "2023 % 05 \t 03 \n",
		},
		"full year week based no padding": {
			date:   value.MakeDate(512, 5, 3),
			format: "%-G",
			want:   "512",
		},
		"full year week based space padding": {
			date:   value.MakeDate(512, 5, 3),
			format: "%_G",
			want:   " 512",
		},
		"full year week based zero padding": {
			date:   value.MakeDate(512, 5, 3),
			format: "%G",
			want:   "0512",
		},
		"full year no padding": {
			date:   value.MakeDate(512, 5, 3),
			format: "%-Y",
			want:   "512",
		},
		"full year space padding": {
			date:   value.MakeDate(512, 5, 3),
			format: "%_Y",
			want:   " 512",
		},
		"full year zero padding": {
			date:   value.MakeDate(512, 5, 3),
			format: "%Y",
			want:   "0512",
		},
		"century no padding": {
			date:   value.MakeDate(512, 5, 3),
			format: "%-C",
			want:   "5",
		},
		"century space padding": {
			date:   value.MakeDate(512, 5, 3),
			format: "%_C",
			want:   " 5",
		},
		"century zero padding": {
			date:   value.MakeDate(512, 5, 3),
			format: "%C",
			want:   "05",
		},
		"year last two week based no padding": {
			date:   value.MakeDate(2003, 5, 3),
			format: "%-g",
			want:   "3",
		},
		"year last two week based space padding": {
			date:   value.MakeDate(2003, 5, 3),
			format: "%_g",
			want:   " 3",
		},
		"year last two week based zero padding": {
			date:   value.MakeDate(2003, 5, 3),
			format: "%g",
			want:   "03",
		},
		"year last two no padding": {
			date:   value.MakeDate(2003, 5, 3),
			format: "%-y",
			want:   "3",
		},
		"year last two space padding": {
			date:   value.MakeDate(2003, 5, 3),
			format: "%_y",
			want:   " 3",
		},
		"year last two zero padding": {
			date:   value.MakeDate(2003, 5, 3),
			format: "%y",
			want:   "03",
		},
		"month no padding": {
			date:   value.MakeDate(2003, 5, 3),
			format: "%-m",
			want:   "5",
		},
		"month space padding": {
			date:   value.MakeDate(2003, 5, 3),
			format: "%_m",
			want:   " 5",
		},
		"month zero padding": {
			date:   value.MakeDate(2003, 5, 3),
			format: "%m",
			want:   "05",
		},
		"month full name": {
			date:   value.MakeDate(2003, 6, 3),
			format: "%B",
			want:   "June",
		},
		"month full name uppercase": {
			date:   value.MakeDate(2003, 6, 3),
			format: "%^B",
			want:   "JUNE",
		},
		"month abbreviated name": {
			date:   value.MakeDate(2003, 6, 3),
			format: "%b",
			want:   "Jun",
		},
		"month abbreviated name uppercase": {
			date:   value.MakeDate(2003, 6, 3),
			format: "%^b",
			want:   "JUN",
		},
		"day of month no padding": {
			date:   value.MakeDate(2003, 6, 3),
			format: "%-d",
			want:   "3",
		},
		"day of month space padding": {
			date:   value.MakeDate(2003, 6, 3),
			format: "%_d",
			want:   " 3",
		},
		"day of month zero padding": {
			date:   value.MakeDate(2003, 6, 3),
			format: "%d",
			want:   "03",
		},
		"day of year no padding": {
			date:   value.MakeDate(2003, 2, 3),
			format: "%-j",
			want:   "34",
		},
		"day of year space padding": {
			date:   value.MakeDate(2003, 2, 3),
			format: "%_j",
			want:   " 34",
		},
		"day of year zero padding": {
			date:   value.MakeDate(2003, 2, 3),
			format: "%j",
			want:   "034",
		},
		"weekday full name": {
			date:   value.MakeDate(2023, 11, 28),
			format: "%A",
			want:   "Tuesday",
		},
		"weekday full name uppercase": {
			date:   value.MakeDate(2023, 11, 28),
			format: "%^A",
			want:   "TUESDAY",
		},
		"weekday abbreviated name": {
			date:   value.MakeDate(2023, 11, 28),
			format: "%a",
			want:   "Tue",
		},
		"weekday abbreviated name uppercase": {
			date:   value.MakeDate(2023, 11, 28),
			format: "%^a",
			want:   "TUE",
		},
		"ISO week no padding": {
			date:   value.MakeDate(2023, 2, 3),
			format: "%-V",
			want:   "5",
		},
		"ISO week space padding": {
			date:   value.MakeDate(2023, 2, 3),
			format: "%_V",
			want:   " 5",
		},
		"ISO week zero padding": {
			date:   value.MakeDate(2023, 2, 3),
			format: "%V",
			want:   "05",
		},
		"week monday no padding": {
			date:   value.MakeDate(2023, 2, 3),
			format: "%-W",
			want:   "5",
		},
		"week monday space padding": {
			date:   value.MakeDate(2023, 2, 3),
			format: "%_W",
			want:   " 5",
		},
		"week monday zero padding": {
			date:   value.MakeDate(2023, 2, 3),
			format: "%W",
			want:   "05",
		},
		"week sunday no padding": {
			date:   value.MakeDate(2023, 2, 3),
			format: "%-U",
			want:   "5",
		},
		"week sunday space padding": {
			date:   value.MakeDate(2023, 2, 3),
			format: "%_U",
			want:   " 5",
		},
		"week sunday zero padding": {
			date:   value.MakeDate(2023, 2, 3),
			format: "%U",
			want:   "05",
		},
		"date": {
			date:   value.MakeDate(2023, 2, 3),
			format: "%D",
			want:   "02/03/23",
		},
		"ISO8601 date": {
			date:   value.MakeDate(2023, 2, 3),
			format: "%F",
			want:   "2023-02-03",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.date.Format(tc.format)
			opts := comparer.Options()
			if diff := cmp.Diff(tc.err, err, opts...); diff != "" {
				t.Fatal(diff)
			}
			if !tc.err.IsUndefined() {
				return
			}

			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
