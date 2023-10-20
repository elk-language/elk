package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSymbolTableGet(t *testing.T) {
	tests := map[string]struct {
		table *symbolTableStruct
		get   string
		want  Symbol
		ok    bool
	}{
		"return nil when empty table": {
			table: newSymbolTable(),
			get:   "foo",
			want:  -1,
			ok:    false,
		},
		"return nil when no such symbol": {
			table: newSymbolTable(symbolTableWithNameTable(map[string]Symbol{
				"bar": 1,
			})),
			get:  "foo",
			want: -1,
			ok:   false,
		},
		"return symbol when present": {
			table: newSymbolTable(symbolTableWithNameTable(map[string]Symbol{
				"foo": 1,
			})),
			get:  "foo",
			want: 1,
			ok:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.table.Get(tc.get)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSymbolTableAdd(t *testing.T) {
	tests := map[string]struct {
		table      *symbolTableStruct
		add        string
		want       Symbol
		tableAfter *symbolTableStruct
	}{
		"add to an empty table": {
			table: newSymbolTable(),
			add:   "foo",
			want:  0,
			tableAfter: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
					},
				),
				symbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
		},
		"add to a populated table": {
			table: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
					},
				),
				symbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
			add:  "bar",
			want: 1,
			tableAfter: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
						"bar": 1,
					},
				),
				symbolTableWithIdTable(
					[]string{
						"foo",
						"bar",
					},
				),
			),
		},
		"add an already existing symbol": {
			table: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
					},
				),
				symbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
			add:  "foo",
			want: 0,
			tableAfter: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
					},
				),
				symbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.table.Add(tc.add)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			opts := []cmp.Option{
				cmp.AllowUnexported(symbolTableStruct{}),
				cmpopts.IgnoreFields(symbolTableStruct{}, "mutex"),
			}
			if diff := cmp.Diff(tc.tableAfter, tc.table, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
