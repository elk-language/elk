package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSymbolMapGet(t *testing.T) {
	tests := map[string]struct {
		symbolMap value.SymbolMap
		get       value.Symbol
		want      value.Value
		ok        bool
	}{
		"return nil when the map is empty": {
			symbolMap: make(value.SymbolMap),
			get:       1,
			want:      nil,
			ok:        false,
		},
		"return nil when no such symbol": {
			symbolMap: value.SymbolMap{
				1: value.SmallInt(5),
			},
			get:  20,
			want: nil,
			ok:   false,
		},
		"return the value when the key is present": {
			symbolMap: value.SymbolMap{
				1: value.SmallInt(5),
			},
			get:  1,
			want: value.SmallInt(5),
			ok:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := tc.symbolMap.Get(tc.get)
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

func TestSymbolMapGetString(t *testing.T) {
	tests := map[string]struct {
		symbolTable      *value.SymbolTableStruct
		symbolMap        value.SymbolMap
		get              string
		want             value.Value
		ok               bool
		symbolTableAfter *value.SymbolTableStruct
	}{
		"return nil when the map is empty": {
			symbolTable:      value.NewSymbolTable(),
			symbolMap:        make(value.SymbolMap),
			get:              "foo",
			want:             nil,
			ok:               false,
			symbolTableAfter: value.NewSymbolTable(),
		},
		"return nil when no such symbol": {
			symbolTable: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"foo": 0,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
			symbolMap: value.SymbolMap{
				1: value.SmallInt(5),
			},
			get:  "foo",
			want: nil,
			symbolTableAfter: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"foo": 0,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
		},
		"return the value when the key is present": {
			symbolTable: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"foo": 0,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
			symbolMap: value.SymbolMap{
				0: value.SmallInt(5),
			},
			get:  "foo",
			want: value.SmallInt(5),
			ok:   true,
			symbolTableAfter: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"foo": 0,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			originalSymbolTable := value.SymbolTable
			value.SymbolTable = tc.symbolTable
			got, ok := tc.symbolMap.GetString(tc.get)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatalf(diff)
			}
			opts := []cmp.Option{
				cmp.AllowUnexported(value.SymbolTableStruct{}),
				cmpopts.IgnoreFields(value.SymbolTableStruct{}, "mutex"),
			}
			if diff := cmp.Diff(tc.symbolTableAfter, value.SymbolTable, opts...); diff != "" {
				t.Fatalf(diff)
			}
			value.SymbolTable = originalSymbolTable
		})
	}
}

func TestSymbolMapSet(t *testing.T) {
	tests := map[string]struct {
		symbolMap value.SymbolMap
		key       value.Symbol
		value     value.Value
		want      value.SymbolMap
	}{
		"add to an empty map": {
			symbolMap: value.SymbolMap{},
			key:       1,
			value:     value.SmallInt(5),
			want: value.SymbolMap{
				1: value.SmallInt(5),
			},
		},
		"add to a populated map": {
			symbolMap: value.SymbolMap{
				1: value.SmallInt(5),
			},
			key:   20,
			value: value.RootModule,
			want: value.SymbolMap{
				1:  value.SmallInt(5),
				20: value.RootModule,
			},
		},
		"overwrite an already existing value": {
			symbolMap: value.SymbolMap{
				1:  value.SmallInt(5),
				20: value.RootModule,
			},
			key:   20,
			value: value.SmallInt(-2),
			want: value.SymbolMap{
				1:  value.SmallInt(5),
				20: value.SmallInt(-2),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.symbolMap.Set(tc.key, tc.value)
			got := tc.symbolMap
			opts := comparer.Comparer
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSymbolMapSetString(t *testing.T) {
	tests := map[string]struct {
		symbolTable      *value.SymbolTableStruct
		symbolMap        value.SymbolMap
		key              string
		value            value.Value
		want             value.SymbolMap
		symbolTableAfter *value.SymbolTableStruct
	}{
		"add to an empty map": {
			symbolTable: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"bar": 0,
						"foo": 1,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "bar",
						1: "foo",
					},
				),
			),
			symbolMap: value.SymbolMap{},
			key:       "foo",
			value:     value.SmallInt(5),
			want: value.SymbolMap{
				1: value.SmallInt(5),
			},
			symbolTableAfter: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"bar": 0,
						"foo": 1,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "bar",
						1: "foo",
					},
				),
			),
		},
		"add to a populated map": {
			symbolTable: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"bar": 0,
						"foo": 1,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "bar",
						1: "foo",
					},
				),
			),
			symbolMap: value.SymbolMap{
				0: value.SmallInt(5),
			},
			key:   "foo",
			value: value.RootModule,
			want: value.SymbolMap{
				0: value.SmallInt(5),
				1: value.RootModule,
			},
			symbolTableAfter: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"bar": 0,
						"foo": 1,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "bar",
						1: "foo",
					},
				),
			),
		},
		"add a new symbol": {
			symbolTable: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"foo": 0,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "foo",
					},
				),
			),
			symbolMap: value.SymbolMap{
				0: value.SmallInt(5),
			},
			key:   "bar",
			value: value.RootModule,
			want: value.SymbolMap{
				0: value.SmallInt(5),
				1: value.RootModule,
			},
			symbolTableAfter: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"foo": 0,
						"bar": 1,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "foo",
						1: "bar",
					},
				),
			),
		},
		"overwrite an already existing value": {
			symbolTable: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"foo": 0,
						"bar": 1,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "foo",
						1: "bar",
					},
				),
			),
			symbolMap: value.SymbolMap{
				0: value.SmallInt(5),
				1: value.RootModule,
			},
			key:   "bar",
			value: value.SmallInt(-2),
			want: value.SymbolMap{
				0: value.SmallInt(5),
				1: value.SmallInt(-2),
			},
			symbolTableAfter: value.NewSymbolTable(
				value.SymbolTableWithNameTable(
					map[string]value.Symbol{
						"foo": 0,
						"bar": 1,
					},
				),
				value.SymbolTableWithIdTable(
					[]string{
						0: "foo",
						1: "bar",
					},
				),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			originalSymbolTable := value.SymbolTable
			value.SymbolTable = tc.symbolTable
			tc.symbolMap.SetString(tc.key, tc.value)
			got := tc.symbolMap
			opts := comparer.Comparer
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.symbolTableAfter, value.SymbolTable, opts...); diff != "" {
				t.Fatalf(diff)
			}
			value.SymbolTable = originalSymbolTable
		})
	}
}

func TestSymbolMapInspect(t *testing.T) {
	tests := map[string]struct {
		symbolMap value.SymbolMap
		want      string
		// ordering in maps is unpredictable
		// so this field can be used to provide
		// a second acceptable result
		wantAlt string
	}{
		"empty map": {
			symbolMap: value.SymbolMap{},
			want:      "{}",
		},
		"single entry": {
			symbolMap: value.SymbolMap{
				value.ToSymbol("foo"): value.Int64(5),
			},
			want: "{ foo: 5i64 }",
		},
		"multiple entries": {
			symbolMap: value.SymbolMap{
				value.ToSymbol("foo"): value.String("baz"),
				value.ToSymbol("bar"): value.FloatClass,
			},
			want:    `{ foo: "baz", bar: class Std::Float < Std::Numeric }`,
			wantAlt: `{ bar: class Std::Float < Std::Numeric, foo: "baz" }`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.symbolMap.Inspect()
			opts := comparer.Comparer
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				if tc.wantAlt == "" {
					t.Fatalf(diff)
				}
				if diff := cmp.Diff(tc.wantAlt, got, opts...); diff != "" {
					t.Fatalf(diff)
				}
			}
		})
	}
}
