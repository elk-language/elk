package value_test

import (
	"testing"

	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

func TestTimeFormat(t *testing.T) {
	tests := map[string]struct {
		time   value.Time
		format string
		want   string
		err    *value.Error
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := tc.time.Format(tc.format)
			if diff := cmp.Diff(tc.err, err); diff != "" {
				t.Fatalf(diff)
			}
			if tc.err != nil {
				return
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
