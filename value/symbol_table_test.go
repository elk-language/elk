package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSymbolTableGet(t *testing.T) {
	tests := map[string]struct {
		table *SymbolTableStruct
		get   string
		want  Symbol
		ok    bool
	}{
		"return nil when empty table": {
			table: NewSymbolTable(),
			get:   "foo",
			want:  -1,
			ok:    false,
		},
		"return nil when no such symbol": {
			table: NewSymbolTable(SymbolTableWithNameTable(map[string]Symbol{
				"bar": 1,
			})),
			get:  "foo",
			want: -1,
			ok:   false,
		},
		"return symbol when present": {
			table: NewSymbolTable(SymbolTableWithNameTable(map[string]Symbol{
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
		table      *SymbolTableStruct
		add        string
		want       Symbol
		tableAfter *SymbolTableStruct
	}{
		"add to an empty table": {
			table: NewSymbolTable(),
			add:   "foo",
			want:  0,
			tableAfter: NewSymbolTable(
				SymbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
					},
				),
				SymbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
		},
		"add to a populated table": {
			table: NewSymbolTable(
				SymbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
					},
				),
				SymbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
			add:  "bar",
			want: 1,
			tableAfter: NewSymbolTable(
				SymbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
						"bar": 1,
					},
				),
				SymbolTableWithIdTable(
					[]string{
						"foo",
						"bar",
					},
				),
			),
		},
		"add an already existing symbol": {
			table: NewSymbolTable(
				SymbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
					},
				),
				SymbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
			add:  "foo",
			want: 0,
			tableAfter: NewSymbolTable(
				SymbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
					},
				),
				SymbolTableWithIdTable(
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
				cmp.AllowUnexported(SymbolTableStruct{}),
				cmpopts.IgnoreFields(SymbolTableStruct{}, "mutex"),
			}
			if diff := cmp.Diff(tc.tableAfter, tc.table, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
