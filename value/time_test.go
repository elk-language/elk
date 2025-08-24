package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestTimeFormat(t *testing.T) {
	tests := map[string]struct {
		time   value.Time
		format string
		want   string
		err    value.Value
	}{
		"hour no padding": {
			time:   value.MakeTime(8, 9, 0, 0),
			format: "%-H",
			want:   "8",
		},
		"hour space padding": {
			time:   value.MakeTime(8, 9, 0, 0),
			format: "%_H",
			want:   " 8",
		},
		"hour zero padding": {
			time:   value.MakeTime(8, 9, 0, 0),
			format: "%H",
			want:   "08",
		},
		"hour 12 no padding": {
			time:   value.MakeTime(16, 9, 0, 0),
			format: "%-I",
			want:   "4",
		},
		"hour 12 space padding": {
			time:   value.MakeTime(16, 9, 0, 0),
			format: "%_I",
			want:   " 4",
		},
		"hour 12 zero padding": {
			time:   value.MakeTime(16, 9, 0, 0),
			format: "%I",
			want:   "04",
		},
		"meridiem lowercase PM": {
			time:   value.MakeTime(16, 9, 0, 0),
			format: "%I %P",
			want:   "04 pm",
		},
		"meridiem lowercase 12 PM": {
			time:   value.MakeTime(12, 9, 0, 0),
			format: "%I %P",
			want:   "12 pm",
		},
		"meridiem lowercase 12 AM": {
			time:   value.MakeTime(0, 9, 0, 0),
			format: "%I %P",
			want:   "12 am",
		},
		"meridiem lowercase 6 AM": {
			time:   value.MakeTime(6, 9, 0, 0),
			format: "%I %P",
			want:   "06 am",
		},
		"meridiem uppercase PM": {
			time:   value.MakeTime(16, 9, 0, 0),
			format: "%I %p %^P",
			want:   "04 PM PM",
		},
		"meridiem uppercase 12 PM": {
			time:   value.MakeTime(12, 9, 0, 0),
			format: "%I %p %^P",
			want:   "12 PM PM",
		},
		"meridiem uppercase 12 AM": {
			time:   value.MakeTime(0, 9, 0, 0),
			format: "%I %p %^P",
			want:   "12 AM AM",
		},
		"meridiem uppercase 6 AM": {
			time:   value.MakeTime(6, 9, 0, 0),
			format: "%I %p %^P",
			want:   "06 AM AM",
		},
		"minute no padding": {
			time:   value.MakeTime(16, 9, 0, 0),
			format: "%-M",
			want:   "9",
		},
		"minute space padding": {
			time:   value.MakeTime(16, 9, 0, 0),
			format: "%_M",
			want:   " 9",
		},
		"minute zero padding": {
			time:   value.MakeTime(16, 9, 0, 0),
			format: "%M",
			want:   "09",
		},
		"second no padding": {
			time:   value.MakeTime(16, 9, 7, 0),
			format: "%-S",
			want:   "7",
		},
		"second space padding": {
			time:   value.MakeTime(16, 9, 7, 0),
			format: "%_S",
			want:   " 7",
		},
		"second zero padding": {
			time:   value.MakeTime(16, 9, 7, 0),
			format: "%S",
			want:   "07",
		},
		"millisecond no padding": {
			time:   value.MakeTime(16, 9, 7, 62_000_000),
			format: "%-L:%-3N",
			want:   "62:62",
		},
		"millisecond space padding": {
			time:   value.MakeTime(16, 9, 7, 62_000_000),
			format: "%_L:%_3N",
			want:   " 62: 62",
		},
		"millisecond zero padding": {
			time:   value.MakeTime(16, 9, 7, 62_000_000),
			format: "%L:%3N",
			want:   "062:062",
		},
		"time 12": {
			time:   value.MakeTime(16, 9, 7, 0),
			format: "%r",
			want:   "04:09:07 PM",
		},
		"time 24": {
			time:   value.MakeTime(16, 9, 7, 0),
			format: "%R",
			want:   "16:09",
		},
		"time 24 seconds": {
			time:   value.MakeTime(16, 9, 7, 0),
			format: "%T",
			want:   "16:09:07",
		},
		"microsecond no padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_000),
			format: "%-6N",
			want:   "32765",
		},
		"microsecond space padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_000),
			format: "%_6N",
			want:   " 32765",
		},
		"microsecond zero padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_000),
			format: "%6N",
			want:   "032765",
		},
		"nanosecond no padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%-9N",
			want:   "32765198",
		},
		"nanosecond space padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%_9N",
			want:   " 32765198",
		},
		"nanosecond zero padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%9N",
			want:   "032765198",
		},
		"picosecond no padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%-12N",
			want:   "32765198000",
		},
		"picosecond space padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%_12N",
			want:   " 32765198000",
		},
		"picosecond zero padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%12N",
			want:   "032765198000",
		},
		"femtosecond no padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%-15N",
			want:   "32765198000000",
		},
		"femtosecond space padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%_15N",
			want:   " 32765198000000",
		},
		"femtosecond zero padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%15N",
			want:   "032765198000000",
		},
		"attosecond no padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%-18N",
			want:   "32765198000000000",
		},
		"attosecond space padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%_18N",
			want:   " 32765198000000000",
		},
		"attosecond zero padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%18N",
			want:   "032765198000000000",
		},
		"zeptosecond no padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%-21N",
			want:   "32765198000000000000",
		},
		"zeptosecond space padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%_21N",
			want:   " 32765198000000000000",
		},
		"zeptosecond zero padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%21N",
			want:   "032765198000000000000",
		},
		"yoctosecond no padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%-24N",
			want:   "32765198000000000000000",
		},
		"yoctosecond space padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%_24N",
			want:   " 32765198000000000000000",
		},
		"yoctosecond zero padding": {
			time:   value.MakeTime(16, 9, 7, 32_765_198),
			format: "%24N",
			want:   "032765198000000000000000",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.time.Format(tc.format)
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
