package object

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSymbolMapGet(t *testing.T) {
	tests := map[string]struct {
		symbolMap SymbolMap
		get       *Symbol
		want      Object
	}{
		"return nil when the map is empty": {
			symbolMap: make(SymbolMap),
			get:       newSymbol("foo", 1),
			want:      nil,
		},
		"return nil when no such symbol": {
			symbolMap: SymbolMap{
				1: SmallInt(5),
			},
			get:  newSymbol("foo", 20),
			want: nil,
		},
		"return the value when the key is present": {
			symbolMap: SymbolMap{
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

func TestSymbolMapGetId(t *testing.T) {
	tests := map[string]struct {
		symbolMap SymbolMap
		get       SymbolId
		want      Object
	}{
		"return nil when the map is empty": {
			symbolMap: make(SymbolMap),
			get:       1,
			want:      nil,
		},
		"return nil when no such symbol": {
			symbolMap: SymbolMap{
				1: SmallInt(5),
			},
			get:  20,
			want: nil,
		},
		"return the value when the key is present": {
			symbolMap: SymbolMap{
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

func TestSymbolMapGetString(t *testing.T) {
	tests := map[string]struct {
		symbolTable      *symbolTableStruct
		symbolMap        SymbolMap
		get              string
		want             Object
		symbolTableAfter *symbolTableStruct
	}{
		"return nil when the map is empty": {
			symbolTable:      newSymbolTable(),
			symbolMap:        make(SymbolMap),
			get:              "foo",
			want:             nil,
			symbolTableAfter: newSymbolTable(),
		},
		"return nil when no such symbol": {
			symbolTable: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 20),
					},
				),
				symbolTableWithLastId(20),
			),
			symbolMap: SymbolMap{
				1: SmallInt(5),
			},
			get:  "foo",
			want: nil,
			symbolTableAfter: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 20),
					},
				),
				symbolTableWithLastId(20),
			),
		},
		"return the value when the key is present": {
			symbolTable: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 1),
					},
				),
				symbolTableWithLastId(1),
			),
			symbolMap: SymbolMap{
				1: SmallInt(5),
			},
			get:  "foo",
			want: SmallInt(5),
			symbolTableAfter: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 1),
					},
				),
				symbolTableWithLastId(1),
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

func TestSymbolMapSet(t *testing.T) {
	tests := map[string]struct {
		symbolMap SymbolMap
		key       *Symbol
		value     Object
		want      SymbolMap
	}{
		"add to an empty map": {
			symbolMap: SymbolMap{},
			key:       newSymbol("foo", 1),
			value:     SmallInt(5),
			want: SymbolMap{
				1: SmallInt(5),
			},
		},
		"add to a populated map": {
			symbolMap: SymbolMap{
				1: SmallInt(5),
			},
			key:   newSymbol("foo", 20),
			value: RootModule,
			want: SymbolMap{
				1:  SmallInt(5),
				20: RootModule,
			},
		},
		"overwrite an already existing value": {
			symbolMap: SymbolMap{
				1:  SmallInt(5),
				20: RootModule,
			},
			key:   newSymbol("foo", 20),
			value: SmallInt(-2),
			want: SymbolMap{
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
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSymbolMapSetId(t *testing.T) {
	tests := map[string]struct {
		symbolMap SymbolMap
		key       SymbolId
		value     Object
		want      SymbolMap
	}{
		"add to an empty map": {
			symbolMap: SymbolMap{},
			key:       1,
			value:     SmallInt(5),
			want: SymbolMap{
				1: SmallInt(5),
			},
		},
		"add to a populated map": {
			symbolMap: SymbolMap{
				1: SmallInt(5),
			},
			key:   20,
			value: RootModule,
			want: SymbolMap{
				1:  SmallInt(5),
				20: RootModule,
			},
		},
		"overwrite an already existing value": {
			symbolMap: SymbolMap{
				1:  SmallInt(5),
				20: RootModule,
			},
			key:   20,
			value: SmallInt(-2),
			want: SymbolMap{
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
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}

func TestSymbolMapSetString(t *testing.T) {
	tests := map[string]struct {
		symbolTable      *symbolTableStruct
		symbolMap        SymbolMap
		key              string
		value            Object
		want             SymbolMap
		symbolTableAfter *symbolTableStruct
	}{
		"add to an empty map": {
			symbolTable: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 1),
					},
				),
				symbolTableWithLastId(20),
			),
			symbolMap: SymbolMap{},
			key:       "foo",
			value:     SmallInt(5),
			want: SymbolMap{
				1: SmallInt(5),
			},
			symbolTableAfter: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 1),
					},
				),
				symbolTableWithLastId(20),
			),
		},
		"add to a populated map": {
			symbolTable: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 20),
					},
				),
				symbolTableWithLastId(20),
			),
			symbolMap: SymbolMap{
				1: SmallInt(5),
			},
			key:   "foo",
			value: RootModule,
			want: SymbolMap{
				1:  SmallInt(5),
				20: RootModule,
			},
			symbolTableAfter: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 20),
					},
				),
				symbolTableWithLastId(20),
			),
		},
		"add a new symbol": {
			symbolTable: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 20),
					},
				),
				symbolTableWithLastId(20),
			),
			symbolMap: SymbolMap{
				1: SmallInt(5),
			},
			key:   "bar",
			value: RootModule,
			want: SymbolMap{
				1:  SmallInt(5),
				21: RootModule,
			},
			symbolTableAfter: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 20),
						"bar": newSymbol("bar", 21),
					},
				),
				symbolTableWithLastId(21),
			),
		},
		"overwrite an already existing value": {
			symbolTable: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 20),
					},
				),
				symbolTableWithLastId(20),
			),
			symbolMap: SymbolMap{
				1:  SmallInt(5),
				20: RootModule,
			},
			key:   "foo",
			value: SmallInt(-2),
			want: SymbolMap{
				1:  SmallInt(5),
				20: SmallInt(-2),
			},
			symbolTableAfter: newSymbolTable(
				symbolTableWithTable(
					map[string]*Symbol{
						"foo": newSymbol("foo", 20),
					},
				),
				symbolTableWithLastId(20),
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
			}
			if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
				t.Fatalf(diff)
			}
			opts = []cmp.Option{
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
