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
		want  *Symbol
	}{
		"return nil when empty table": {
			table: newSymbolTable(),
			get:   "foo",
			want:  nil,
		},
		"return nil when no such symbol": {
			table: newSymbolTable(symbolTableWithNameTable(map[string]*Symbol{
				"bar": newSymbol("bar", 1),
			})),
			get:  "foo",
			want: nil,
		},
		"return symbol when present": {
			table: newSymbolTable(symbolTableWithNameTable(map[string]*Symbol{
				"foo": newSymbol("foo", 1),
			})),
			get:  "foo",
			want: newSymbol("foo", 1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.table.Get(tc.get)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSymbolTableAdd(t *testing.T) {
	tests := map[string]struct {
		table      *symbolTableStruct
		add        string
		want       *Symbol
		tableAfter *symbolTableStruct
	}{
		"add to an empty table": {
			table: newSymbolTable(),
			add:   "foo",
			want:  newSymbol("foo", 0),
			tableAfter: newSymbolTable(
				symbolTableWithNameTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 0),
					},
				),
				symbolTableWithIdTable(
					[]*Symbol{
						newSymbol("foo", 0),
					},
				),
			),
		},
		"add to a populated table": {
			table: newSymbolTable(
				symbolTableWithNameTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 0),
					},
				),
				symbolTableWithIdTable(
					[]*Symbol{
						newSymbol("foo", 0),
					},
				),
			),
			add:  "bar",
			want: newSymbol("bar", 1),
			tableAfter: newSymbolTable(
				symbolTableWithNameTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 0),
						"bar": newSymbol("bar", 1),
					},
				),
				symbolTableWithIdTable(
					[]*Symbol{
						newSymbol("foo", 0),
						newSymbol("bar", 1),
					},
				),
			),
		},
		"add an already existing symbol": {
			table: newSymbolTable(
				symbolTableWithNameTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 0),
					},
				),
				symbolTableWithIdTable(
					[]*Symbol{
						newSymbol("foo", 0),
					},
				),
			),
			add:  "foo",
			want: newSymbol("foo", 0),
			tableAfter: newSymbolTable(
				symbolTableWithNameTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 0),
					},
				),
				symbolTableWithIdTable(
					[]*Symbol{
						newSymbol("foo", 0),
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
