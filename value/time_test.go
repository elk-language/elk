package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestTimeFormat(t *testing.T) {
	tests := map[string]struct {
		time   *value.Time
		format string
		want   string
		err    value.Value
	}{
		"simple date with zero padding": {
			time:   value.NewTime(2023, 5, 3, 0, 0, 0, 0, nil),
			format: "%Y-%m-%d",
			want:   `2023-05-03`,
		},
		"simple date with space padding": {
			time:   value.NewTime(520, 5, 3, 0, 0, 0, 0, nil),
			format: "%_Y %_m %_d",
			want:   ` 520  5  3`,
		},
		"literal values": {
			time:   value.NewTime(2023, 5, 3, 0, 0, 0, 0, nil),
			format: "%Y %% %m %t %d %n",
			want:   "2023 % 05 \t 03 \n",
		},
		"full year week based no padding": {
			time:   value.NewTime(512, 5, 3, 0, 0, 0, 0, nil),
			format: "%-G",
			want:   "512",
		},
		"full year week based space padding": {
			time:   value.NewTime(512, 5, 3, 0, 0, 0, 0, nil),
			format: "%_G",
			want:   " 512",
		},
		"full year week based zero padding": {
			time:   value.NewTime(512, 5, 3, 0, 0, 0, 0, nil),
			format: "%G",
			want:   "0512",
		},
		"full year no padding": {
			time:   value.NewTime(512, 5, 3, 0, 0, 0, 0, nil),
			format: "%-Y",
			want:   "512",
		},
		"full year space padding": {
			time:   value.NewTime(512, 5, 3, 0, 0, 0, 0, nil),
			format: "%_Y",
			want:   " 512",
		},
		"full year zero padding": {
			time:   value.NewTime(512, 5, 3, 0, 0, 0, 0, nil),
			format: "%Y",
			want:   "0512",
		},
		"century no padding": {
			time:   value.NewTime(512, 5, 3, 0, 0, 0, 0, nil),
			format: "%-C",
			want:   "5",
		},
		"century space padding": {
			time:   value.NewTime(512, 5, 3, 0, 0, 0, 0, nil),
			format: "%_C",
			want:   " 5",
		},
		"century zero padding": {
			time:   value.NewTime(512, 5, 3, 0, 0, 0, 0, nil),
			format: "%C",
			want:   "05",
		},
		"year last two week based no padding": {
			time:   value.NewTime(2003, 5, 3, 0, 0, 0, 0, nil),
			format: "%-g",
			want:   "3",
		},
		"year last two week based space padding": {
			time:   value.NewTime(2003, 5, 3, 0, 0, 0, 0, nil),
			format: "%_g",
			want:   " 3",
		},
		"year last two week based zero padding": {
			time:   value.NewTime(2003, 5, 3, 0, 0, 0, 0, nil),
			format: "%g",
			want:   "03",
		},
		"year last two no padding": {
			time:   value.NewTime(2003, 5, 3, 0, 0, 0, 0, nil),
			format: "%-y",
			want:   "3",
		},
		"year last two space padding": {
			time:   value.NewTime(2003, 5, 3, 0, 0, 0, 0, nil),
			format: "%_y",
			want:   " 3",
		},
		"year last two zero padding": {
			time:   value.NewTime(2003, 5, 3, 0, 0, 0, 0, nil),
			format: "%y",
			want:   "03",
		},
		"month no padding": {
			time:   value.NewTime(2003, 5, 3, 0, 0, 0, 0, nil),
			format: "%-m",
			want:   "5",
		},
		"month space padding": {
			time:   value.NewTime(2003, 5, 3, 0, 0, 0, 0, nil),
			format: "%_m",
			want:   " 5",
		},
		"month zero padding": {
			time:   value.NewTime(2003, 5, 3, 0, 0, 0, 0, nil),
			format: "%m",
			want:   "05",
		},
		"month full name": {
			time:   value.NewTime(2003, 6, 3, 0, 0, 0, 0, nil),
			format: "%B",
			want:   "June",
		},
		"month full name uppercase": {
			time:   value.NewTime(2003, 6, 3, 0, 0, 0, 0, nil),
			format: "%^B",
			want:   "JUNE",
		},
		"month abbreviated name": {
			time:   value.NewTime(2003, 6, 3, 0, 0, 0, 0, nil),
			format: "%b",
			want:   "Jun",
		},
		"month abbreviated name uppercase": {
			time:   value.NewTime(2003, 6, 3, 0, 0, 0, 0, nil),
			format: "%^b",
			want:   "JUN",
		},
		"day of month no padding": {
			time:   value.NewTime(2003, 6, 3, 0, 0, 0, 0, nil),
			format: "%-d",
			want:   "3",
		},
		"day of month space padding": {
			time:   value.NewTime(2003, 6, 3, 0, 0, 0, 0, nil),
			format: "%_d",
			want:   " 3",
		},
		"day of month zero padding": {
			time:   value.NewTime(2003, 6, 3, 0, 0, 0, 0, nil),
			format: "%d",
			want:   "03",
		},
		"day of year no padding": {
			time:   value.NewTime(2003, 2, 3, 0, 0, 0, 0, nil),
			format: "%-j",
			want:   "34",
		},
		"day of year space padding": {
			time:   value.NewTime(2003, 2, 3, 0, 0, 0, 0, nil),
			format: "%_j",
			want:   " 34",
		},
		"day of year zero padding": {
			time:   value.NewTime(2003, 2, 3, 0, 0, 0, 0, nil),
			format: "%j",
			want:   "034",
		},
		"hour no padding": {
			time:   value.NewTime(2003, 2, 3, 8, 9, 0, 0, nil),
			format: "%-H",
			want:   "8",
		},
		"hour space padding": {
			time:   value.NewTime(2003, 2, 3, 8, 9, 0, 0, nil),
			format: "%_H",
			want:   " 8",
		},
		"hour zero padding": {
			time:   value.NewTime(2003, 2, 3, 8, 9, 0, 0, nil),
			format: "%H",
			want:   "08",
		},
		"hour 12 no padding": {
			time:   value.NewTime(2003, 2, 3, 16, 9, 0, 0, nil),
			format: "%-I",
			want:   "4",
		},
		"hour 12 space padding": {
			time:   value.NewTime(2003, 2, 3, 16, 9, 0, 0, nil),
			format: "%_I",
			want:   " 4",
		},
		"hour 12 zero padding": {
			time:   value.NewTime(2003, 2, 3, 16, 9, 0, 0, nil),
			format: "%I",
			want:   "04",
		},
		"meridiem lowercase PM": {
			time:   value.NewTime(2003, 2, 3, 16, 9, 0, 0, nil),
			format: "%I %P",
			want:   "04 pm",
		},
		"meridiem lowercase 12 PM": {
			time:   value.NewTime(2003, 2, 3, 12, 9, 0, 0, nil),
			format: "%I %P",
			want:   "12 pm",
		},
		"meridiem lowercase 12 AM": {
			time:   value.NewTime(2003, 2, 3, 0, 9, 0, 0, nil),
			format: "%I %P",
			want:   "12 am",
		},
		"meridiem lowercase 6 AM": {
			time:   value.NewTime(2003, 2, 3, 6, 9, 0, 0, nil),
			format: "%I %P",
			want:   "06 am",
		},
		"meridiem uppercase PM": {
			time:   value.NewTime(2003, 2, 3, 16, 9, 0, 0, nil),
			format: "%I %p %^P",
			want:   "04 PM PM",
		},
		"meridiem uppercase 12 PM": {
			time:   value.NewTime(2003, 2, 3, 12, 9, 0, 0, nil),
			format: "%I %p %^P",
			want:   "12 PM PM",
		},
		"meridiem uppercase 12 AM": {
			time:   value.NewTime(2003, 2, 3, 0, 9, 0, 0, nil),
			format: "%I %p %^P",
			want:   "12 AM AM",
		},
		"meridiem uppercase 6 AM": {
			time:   value.NewTime(2003, 2, 3, 6, 9, 0, 0, nil),
			format: "%I %p %^P",
			want:   "06 AM AM",
		},
		"minute no padding": {
			time:   value.NewTime(2003, 2, 3, 16, 9, 0, 0, nil),
			format: "%-M",
			want:   "9",
		},
		"minute space padding": {
			time:   value.NewTime(2003, 2, 3, 16, 9, 0, 0, nil),
			format: "%_M",
			want:   " 9",
		},
		"minute zero padding": {
			time:   value.NewTime(2003, 2, 3, 16, 9, 0, 0, nil),
			format: "%M",
			want:   "09",
		},
		"second no padding": {
			time:   value.NewTime(2003, 2, 3, 16, 9, 7, 0, nil),
			format: "%-S",
			want:   "7",
		},
		"second space padding": {
			time:   value.NewTime(2003, 2, 3, 16, 9, 7, 0, nil),
			format: "%_S",
			want:   " 7",
		},
		"second zero padding": {
			time:   value.NewTime(2003, 2, 3, 16, 9, 7, 0, nil),
			format: "%S",
			want:   "07",
		},
		"millisecond no padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 62_000_000, nil),
			format: "%-L:%-3N",
			want:   "62:62",
		},
		"millisecond space padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 62_000_000, nil),
			format: "%_L:%_3N",
			want:   " 62: 62",
		},
		"millisecond zero padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 62_000_000, nil),
			format: "%L:%3N",
			want:   "062:062",
		},
		"timezone name": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%Z",
			want:   "CET",
		},
		"timezone offset": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%z",
			want:   "+0100",
		},
		"timezone offset with colon": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%:z",
			want:   "+01:00",
		},
		"weekday full name": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 0, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%A",
			want:   "Tuesday",
		},
		"weekday full name uppercase": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 0, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%^A",
			want:   "TUESDAY",
		},
		"weekday abbreviated name": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 0, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%a",
			want:   "Tue",
		},
		"weekday abbreviated name uppercase": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 0, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%^a",
			want:   "TUE",
		},
		"unix seconds": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 0, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%s",
			want:   "1701184147",
		},
		"unix milliseconds": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 232_000_000, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%Q:%3s",
			want:   "1701184147232:1701184147232",
		},
		"unix microseconds": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 232_765_000, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%6s",
			want:   "1701184147232765",
		},
		"unix nanoseconds": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 232_765_397, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%9s",
			want:   "1701184147232765397",
		},
		"unix picoseconds": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 232_765_397, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%12s",
			want:   "1701184147232765397000",
		},
		"unix femtoseconds": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 232_765_397, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%15s",
			want:   "1701184147232765397000000",
		},
		"unix attoseconds": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 232_765_397, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%18s",
			want:   "1701184147232765397000000000",
		},
		"unix zeptoseconds": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 232_765_397, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%21s",
			want:   "1701184147232765397000000000000",
		},
		"unix yoctoseconds": {
			time:   value.NewTime(2023, 11, 28, 16, 9, 7, 232_765_397, value.MustLoadTimezone("Europe/Warsaw")),
			format: "%24s",
			want:   "1701184147232765397000000000000000",
		},
		"ISO week no padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%-V",
			want:   "5",
		},
		"ISO week space padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%_V",
			want:   " 5",
		},
		"ISO week zero padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%V",
			want:   "05",
		},
		"week monday no padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%-W",
			want:   "5",
		},
		"week monday space padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%_W",
			want:   " 5",
		},
		"week monday zero padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%W",
			want:   "05",
		},
		"week sunday no padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%-U",
			want:   "5",
		},
		"week sunday space padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%_U",
			want:   " 5",
		},
		"week sunday zero padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%U",
			want:   "05",
		},
		"date and time": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%c",
			want:   "Fri Feb  3 16:09:07 2023",
		},
		"date and time uppercase": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%^c",
			want:   "FRI FEB  3 16:09:07 2023",
		},
		"date": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%D",
			want:   "02/03/23",
		},
		"ISO8601 date": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%F",
			want:   "2023-02-03",
		},
		"time 12": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%r",
			want:   "04:09:07 PM",
		},
		"time 24": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%R",
			want:   "16:09",
		},
		"time 24 seconds": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%T",
			want:   "16:09:07",
		},
		"date1": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%+",
			want:   "Fri Feb  3 16:09:07 UTC 2023",
		},
		"date1 uppercase": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 0, nil),
			format: "%^+",
			want:   "FRI FEB  3 16:09:07 UTC 2023",
		},
		"microsecond no padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_000, nil),
			format: "%-6N",
			want:   "32765",
		},
		"microsecond space padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_000, nil),
			format: "%_6N",
			want:   " 32765",
		},
		"microsecond zero padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_000, nil),
			format: "%6N",
			want:   "032765",
		},
		"nanosecond no padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%-9N",
			want:   "32765198",
		},
		"nanosecond space padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%_9N",
			want:   " 32765198",
		},
		"nanosecond zero padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%9N",
			want:   "032765198",
		},
		"picosecond no padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%-12N",
			want:   "32765198000",
		},
		"picosecond space padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%_12N",
			want:   " 32765198000",
		},
		"picosecond zero padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%12N",
			want:   "032765198000",
		},
		"femtosecond no padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%-15N",
			want:   "32765198000000",
		},
		"femtosecond space padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%_15N",
			want:   " 32765198000000",
		},
		"femtosecond zero padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%15N",
			want:   "032765198000000",
		},
		"attosecond no padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%-18N",
			want:   "32765198000000000",
		},
		"attosecond space padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%_18N",
			want:   " 32765198000000000",
		},
		"attosecond zero padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%18N",
			want:   "032765198000000000",
		},
		"zeptosecond no padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%-21N",
			want:   "32765198000000000000",
		},
		"zeptosecond space padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%_21N",
			want:   " 32765198000000000000",
		},
		"zeptosecond zero padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%21N",
			want:   "032765198000000000000",
		},
		"yoctosecond no padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%-24N",
			want:   "32765198000000000000000",
		},
		"yoctosecond space padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%_24N",
			want:   " 32765198000000000000000",
		},
		"yoctosecond zero padding": {
			time:   value.NewTime(2023, 2, 3, 16, 9, 7, 32_765_198, nil),
			format: "%24N",
			want:   "032765198000000000000000",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.time.Format(tc.format)
			if diff := cmp.Diff(tc.err, err); diff != "" {
				t.Fatalf(diff)
			}
			if !tc.err.IsNil() {
				return
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
