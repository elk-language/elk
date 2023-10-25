package value

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSimpleSymbolMapGet(t *testing.T) {
	tests := map[string]struct {
		symbolMap SimpleSymbolMap
		get       Symbol
		want      Value
		ok        bool
	}{
		"return nil when the map is empty": {
			symbolMap: make(SimpleSymbolMap),
			get:       1,
			want:      nil,
			ok:        false,
		},
		"return nil when no such symbol": {
			symbolMap: SimpleSymbolMap{
				1: SmallInt(5),
			},
			get:  20,
			want: nil,
			ok:   false,
		},
		"return the value when the key is present": {
			symbolMap: SimpleSymbolMap{
				1: SmallInt(5),
			},
			get:  1,
			want: SmallInt(5),
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

func TestSimpleSymbolMapGetString(t *testing.T) {
	tests := map[string]struct {
		symbolTable      *symbolTableStruct
		symbolMap        SimpleSymbolMap
		get              string
		want             Value
		ok               bool
		symbolTableAfter *symbolTableStruct
	}{
		"return nil when the map is empty": {
			symbolTable:      newSymbolTable(),
			symbolMap:        make(SimpleSymbolMap),
			get:              "foo",
			want:             nil,
			ok:               false,
			symbolTableAfter: newSymbolTable(),
		},
		"return nil when no such symbol": {
			symbolTable: newSymbolTable(
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
			symbolMap: SimpleSymbolMap{
				1: SmallInt(5),
			},
			get:  "foo",
			want: nil,
			symbolTableAfter: newSymbolTable(
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
		"return the value when the key is present": {
			symbolTable: newSymbolTable(
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
			symbolMap: SimpleSymbolMap{
				0: SmallInt(5),
			},
			get:  "foo",
			want: SmallInt(5),
			ok:   true,
			symbolTableAfter: newSymbolTable(
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
			originalSymbolTable := SymbolTable
			SymbolTable = tc.symbolTable
			got, ok := tc.symbolMap.GetString(tc.get)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			if diff := cmp.Diff(tc.ok, ok); diff != "" {
				t.Fatalf(diff)
			}
			opts := []cmp.Option{
				cmp.AllowUnexported(symbolTableStruct{}),
				cmpopts.IgnoreFields(symbolTableStruct{}, "mutex"),
			}
			if diff := cmp.Diff(tc.symbolTableAfter, SymbolTable, opts...); diff != "" {
				t.Fatalf(diff)
			}
			SymbolTable = originalSymbolTable
		})
	}
}

func TestSimpleSymbolMapSet(t *testing.T) {
	tests := map[string]struct {
		symbolMap SimpleSymbolMap
		key       Symbol
		value     Value
		want      SimpleSymbolMap
	}{
		"add to an empty map": {
			symbolMap: SimpleSymbolMap{},
			key:       1,
			value:     SmallInt(5),
			want: SimpleSymbolMap{
				1: SmallInt(5),
			},
		},
		"add to a populated map": {
			symbolMap: SimpleSymbolMap{
				1: SmallInt(5),
			},
			key:   20,
			value: RootModule,
			want: SimpleSymbolMap{
				1:  SmallInt(5),
				20: RootModule,
			},
		},
		"overwrite an already existing value": {
			symbolMap: SimpleSymbolMap{
				1:  SmallInt(5),
				20: RootModule,
			},
			key:   20,
			value: SmallInt(-2),
			want: SimpleSymbolMap{
				1:  SmallInt(5),
				20: SmallInt(-2),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.symbolMap.Set(tc.key, tc.value)
			got := tc.symbolMap
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				BigFloatComparer,
				FloatComparer,
				Float32Comparer,
				Float64Comparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSimpleSymbolMapSetString(t *testing.T) {
	tests := map[string]struct {
		symbolTable      *symbolTableStruct
		symbolMap        SimpleSymbolMap
		key              string
		value            Value
		want             SimpleSymbolMap
		symbolTableAfter *symbolTableStruct
	}{
		"add to an empty map": {
			symbolTable: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"bar": 0,
						"foo": 1,
					},
				),
				symbolTableWithIdTable(
					[]string{
						0: "bar",
						1: "foo",
					},
				),
			),
			symbolMap: SimpleSymbolMap{},
			key:       "foo",
			value:     SmallInt(5),
			want: SimpleSymbolMap{
				1: SmallInt(5),
			},
			symbolTableAfter: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"bar": 0,
						"foo": 1,
					},
				),
				symbolTableWithIdTable(
					[]string{
						0: "bar",
						1: "foo",
					},
				),
			),
		},
		"add to a populated map": {
			symbolTable: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"bar": 0,
						"foo": 1,
					},
				),
				symbolTableWithIdTable(
					[]string{
						0: "bar",
						1: "foo",
					},
				),
			),
			symbolMap: SimpleSymbolMap{
				0: SmallInt(5),
			},
			key:   "foo",
			value: RootModule,
			want: SimpleSymbolMap{
				0: SmallInt(5),
				1: RootModule,
			},
			symbolTableAfter: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"bar": 0,
						"foo": 1,
					},
				),
				symbolTableWithIdTable(
					[]string{
						0: "bar",
						1: "foo",
					},
				),
			),
		},
		"add a new symbol": {
			symbolTable: newSymbolTable(
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
			symbolMap: SimpleSymbolMap{
				0: SmallInt(5),
			},
			key:   "bar",
			value: RootModule,
			want: SimpleSymbolMap{
				0: SmallInt(5),
				1: RootModule,
			},
			symbolTableAfter: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
						"bar": 1,
					},
				),
				symbolTableWithIdTable(
					[]string{
						0: "foo",
						1: "bar",
					},
				),
			),
		},
		"overwrite an already existing value": {
			symbolTable: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
						"bar": 1,
					},
				),
				symbolTableWithIdTable(
					[]string{
						0: "foo",
						1: "bar",
					},
				),
			),
			symbolMap: SimpleSymbolMap{
				0: SmallInt(5),
				1: RootModule,
			},
			key:   "bar",
			value: SmallInt(-2),
			want: SimpleSymbolMap{
				0: SmallInt(5),
				1: SmallInt(-2),
			},
			symbolTableAfter: newSymbolTable(
				symbolTableWithNameTable(
					map[string]Symbol{
						"foo": 0,
						"bar": 1,
					},
				),
				symbolTableWithIdTable(
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
			originalSymbolTable := SymbolTable
			SymbolTable = tc.symbolTable
			tc.symbolMap.SetString(tc.key, tc.value)
			got := tc.symbolMap
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
				BigFloatComparer,
				FloatComparer,
				Float32Comparer,
				Float64Comparer,
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Logf("got: %s, want: %s", got.Inspect(), tc.want.Inspect())
				t.Fatalf(diff)
			}
			opts = []cmp.Option{
				cmp.AllowUnexported(symbolTableStruct{}),
				cmpopts.IgnoreFields(symbolTableStruct{}, "mutex"),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
			}
			if diff := cmp.Diff(tc.symbolTableAfter, SymbolTable, opts...); diff != "" {
				t.Fatalf(diff)
			}
			SymbolTable = originalSymbolTable
		})
	}
}

func TestSimpleSymbolMapInspect(t *testing.T) {
	tests := map[string]struct {
		symbolMap SimpleSymbolMap
		want      string
		// ordering in maps is unpredictable
		// so this field can be used to provide
		// a second acceptable result
		wantAlt string
	}{
		"empty map": {
			symbolMap: SimpleSymbolMap{},
			want:      "{}",
		},
		"single entry": {
			symbolMap: SimpleSymbolMap{
				SymbolTable.Add("foo"): Int64(5),
			},
			want: "{ foo: 5i64 }",
		},
		"multiple entries": {
			symbolMap: SimpleSymbolMap{
				SymbolTable.Add("foo"): String("baz"),
				SymbolTable.Add("bar"): FloatClass,
			},
			want:    `{ foo: "baz", bar: class Std::Float < Std::Numeric }`,
			wantAlt: `{ bar: class Std::Float < Std::Numeric, foo: "baz" }`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.symbolMap.Inspect()
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
			}
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
