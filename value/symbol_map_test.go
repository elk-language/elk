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
	}{
		"return nil when the map is empty": {
			symbolMap: make(value.SymbolMap),
			get:       1,
			want:      value.Nil,
		},
		"return nil when no such symbol": {
			symbolMap: value.SymbolMap{
				1: value.SmallInt(5).ToValue(),
			},
			get:  20,
			want: value.Nil,
		},
		"return the value when the key is present": {
			symbolMap: value.SymbolMap{
				1: value.SmallInt(5).ToValue(),
			},
			get:  1,
			want: value.SmallInt(5).ToValue(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.symbolMap.Get(tc.get)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
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
		symbolTableAfter *value.SymbolTableStruct
	}{
		"return nil when the map is empty": {
			symbolTable:      value.NewSymbolTable(),
			symbolMap:        make(value.SymbolMap),
			get:              "foo",
			want:             value.Nil,
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
				1: value.SmallInt(5).ToValue(),
			},
			get:  "foo",
			want: value.Nil,
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
				0: value.SmallInt(5).ToValue(),
			},
			get:  "foo",
			want: value.SmallInt(5).ToValue(),
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
			got := tc.symbolMap.GetString(tc.get)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
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
			value:     value.SmallInt(5).ToValue(),
			want: value.SymbolMap{
				1: value.SmallInt(5).ToValue(),
			},
		},
		"add to a populated map": {
			symbolMap: value.SymbolMap{
				1: value.SmallInt(5).ToValue(),
			},
			key:   20,
			value: value.Ref(value.RootModule),
			want: value.SymbolMap{
				1:  value.SmallInt(5).ToValue(),
				20: value.Ref(value.RootModule),
			},
		},
		"overwrite an already existing value": {
			symbolMap: value.SymbolMap{
				1:  value.SmallInt(5).ToValue(),
				20: value.Ref(value.RootModule),
			},
			key:   20,
			value: value.SmallInt(-2).ToValue(),
			want: value.SymbolMap{
				1:  value.SmallInt(5).ToValue(),
				20: value.SmallInt(-2).ToValue(),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.symbolMap.Set(tc.key, tc.value)
			got := tc.symbolMap
			opts := comparer.Options()
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
			value:     value.SmallInt(5).ToValue(),
			want: value.SymbolMap{
				1: value.SmallInt(5).ToValue(),
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
				0: value.SmallInt(5).ToValue(),
			},
			key:   "foo",
			value: value.Ref(value.RootModule),
			want: value.SymbolMap{
				0: value.SmallInt(5).ToValue(),
				1: value.Ref(value.RootModule),
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
				0: value.SmallInt(5).ToValue(),
			},
			key:   "bar",
			value: value.Ref(value.RootModule),
			want: value.SymbolMap{
				0: value.SmallInt(5).ToValue(),
				1: value.Ref(value.RootModule),
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
				0: value.SmallInt(5).ToValue(),
				1: value.Ref(value.RootModule),
			},
			key:   "bar",
			value: value.SmallInt(-2).ToValue(),
			want: value.SymbolMap{
				0: value.SmallInt(5).ToValue(),
				1: value.SmallInt(-2).ToValue(),
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
			opts := comparer.Options()
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
				value.ToSymbol("foo"): value.Int64(5).ToValue(),
			},
			want: "{foo: 5i64}",
		},
		"multiple entries": {
			symbolMap: value.SymbolMap{
				value.ToSymbol("foo"): value.Ref(value.String("baz")),
				value.ToSymbol("bar"): value.Ref(value.FloatClass),
			},
			want:    `{foo: "baz", bar: class Std::Float < Std::Object}`,
			wantAlt: `{bar: class Std::Float < Std::Object, foo: "baz"}`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.symbolMap.Inspect()
			opts := comparer.Options()
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				if tc.wantAlt == "" {
					t.Fatal(diff)
				}
				if diff := cmp.Diff(tc.wantAlt, got, opts...); diff != "" {
					t.Fatal(diff)
				}
			}
		})
	}
}
