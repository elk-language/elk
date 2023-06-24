package object

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSimpleSymbolMapGet(t *testing.T) {
	tests := map[string]struct {
		symbolMap SimpleSymbolMap
		get       *Symbol
		want      Value
	}{
		"return nil when the map is empty": {
			symbolMap: make(SimpleSymbolMap),
			get:       newSymbol("foo", 1),
			want:      nil,
		},
		"return nil when no such symbol": {
			symbolMap: SimpleSymbolMap{
				1: SmallInt(5),
			},
			get:  newSymbol("foo", 20),
			want: nil,
		},
		"return the value when the key is present": {
			symbolMap: SimpleSymbolMap{
				1: SmallInt(5),
			},
			get:  newSymbol("foo", 1),
			want: SmallInt(5),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.symbolMap.Get(tc.get)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSimpleSymbolMapGetId(t *testing.T) {
	tests := map[string]struct {
		symbolMap SimpleSymbolMap
		get       SymbolId
		want      Value
	}{
		"return nil when the map is empty": {
			symbolMap: make(SimpleSymbolMap),
			get:       1,
			want:      nil,
		},
		"return nil when no such symbol": {
			symbolMap: SimpleSymbolMap{
				1: SmallInt(5),
			},
			get:  20,
			want: nil,
		},
		"return the value when the key is present": {
			symbolMap: SimpleSymbolMap{
				1: SmallInt(5),
			},
			get:  1,
			want: SmallInt(5),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.symbolMap.GetId(tc.get)
			if diff := cmp.Diff(tc.want, got); diff != "" {
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
		symbolTableAfter *symbolTableStruct
	}{
		"return nil when the map is empty": {
			symbolTable:      newSymbolTable(),
			symbolMap:        make(SimpleSymbolMap),
			get:              "foo",
			want:             nil,
			symbolTableAfter: newSymbolTable(),
		},
		"return nil when no such symbol": {
			symbolTable: newSymbolTable(
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
			symbolMap: SimpleSymbolMap{
				1: SmallInt(5),
			},
			get:  "foo",
			want: nil,
			symbolTableAfter: newSymbolTable(
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
		"return the value when the key is present": {
			symbolTable: newSymbolTable(
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
			symbolMap: SimpleSymbolMap{
				0: SmallInt(5),
			},
			get:  "foo",
			want: SmallInt(5),
			symbolTableAfter: newSymbolTable(
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
			SymbolTable = tc.symbolTable
			got := tc.symbolMap.GetString(tc.get)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf(diff)
			}
			opts := []cmp.Option{
				cmp.AllowUnexported(symbolTableStruct{}),
				cmpopts.IgnoreFields(symbolTableStruct{}, "mutex"),
			}
			if diff := cmp.Diff(tc.symbolTableAfter, SymbolTable, opts...); diff != "" {
				t.Fatalf(diff)
			}
			SymbolTable = newSymbolTable()
		})
	}
}

func TestSimpleSymbolMapSet(t *testing.T) {
	tests := map[string]struct {
		symbolMap SimpleSymbolMap
		key       *Symbol
		value     Value
		want      SimpleSymbolMap
	}{
		"add to an empty map": {
			symbolMap: SimpleSymbolMap{},
			key:       newSymbol("foo", 1),
			value:     SmallInt(5),
			want: SimpleSymbolMap{
				1: SmallInt(5),
			},
		},
		"add to a populated map": {
			symbolMap: SimpleSymbolMap{
				1: SmallInt(5),
			},
			key:   newSymbol("foo", 20),
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
			key:   newSymbol("foo", 20),
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
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSimpleSymbolMapSetId(t *testing.T) {
	tests := map[string]struct {
		symbolMap SimpleSymbolMap
		key       SymbolId
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
			tc.symbolMap.SetId(tc.key, tc.value)
			got := tc.symbolMap
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
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
					map[string]*Symbol{
						"bar": newSymbol("bar", 0),
						"foo": newSymbol("foo", 1),
					},
				),
				symbolTableWithIdTable(
					[]*Symbol{
						newSymbol("bar", 0),
						newSymbol("foo", 1),
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
					map[string]*Symbol{
						"bar": newSymbol("bar", 0),
						"foo": newSymbol("foo", 1),
					},
				),
				symbolTableWithIdTable(
					[]*Symbol{
						newSymbol("bar", 0),
						newSymbol("foo", 1),
					},
				),
			),
		},
		"add to a populated map": {
			symbolTable: newSymbolTable(
				symbolTableWithNameTable(
					map[string]*Symbol{
						"bar": newSymbol("bar", 0),
						"foo": newSymbol("foo", 1),
					},
				),
				symbolTableWithIdTable(
					[]*Symbol{
						newSymbol("bar", 0),
						newSymbol("foo", 1),
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
					map[string]*Symbol{
						"bar": newSymbol("bar", 0),
						"foo": newSymbol("foo", 1),
					},
				),
				symbolTableWithIdTable(
					[]*Symbol{
						newSymbol("bar", 0),
						newSymbol("foo", 1),
					},
				),
			),
		},
		"add a new symbol": {
			symbolTable: newSymbolTable(
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
		"overwrite an already existing value": {
			symbolTable: newSymbolTable(
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			SymbolTable = tc.symbolTable
			tc.symbolMap.SetString(tc.key, tc.value)
			got := tc.symbolMap
			opts := []cmp.Option{
				cmpopts.IgnoreUnexported(Class{}, Module{}),
				cmpopts.IgnoreFields(Class{}, "ConstructorFunc"),
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
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
			SymbolTable = newSymbolTable()
		})
	}
}

func TestSimpleSymbolMapInspect(t *testing.T) {
	tests := map[string]struct {
		symbolMap SimpleSymbolMap
		want      string
	}{
		"empty map": {
			symbolMap: SimpleSymbolMap{},
			want:      "{}",
		},
		"single entry": {
			symbolMap: SimpleSymbolMap{
				SymbolTable.Add("foo").Id: Int64(5),
			},
			want: "{ foo: 5i64 }",
		},
		"multiple entries": {
			symbolMap: SimpleSymbolMap{
				SymbolTable.Add("foo").Id: String("baz"),
				SymbolTable.Add("bar").Id: FloatClass,
			},
			want: `{ foo: "baz", bar: class Std::Float < Std::Numeric }`,
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
				t.Fatalf(diff)
			}
		})
	}
}
