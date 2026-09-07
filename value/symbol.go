package value

import (
	"slices"
	"unsafe"

	"github.com/elk-language/elk/value/symbol"
)

// Interned string
type Symbol struct {
	Id symbol.Symbol
}

func SortSymbolKeys[V any](m map[Symbol]V) []Symbol {
	keys := make([]Symbol, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b Symbol) int {
		aString := a.String()
		bString := b.String()
		if aString < bString {
			return -1
		}
		if aString > bString {
			return 1
		}
		return 0
	})
	return keys
}

// Convert a string to a Symbol
func ToSymbol[T ~string](str T) Symbol {
	return S(symbol.ToSymbol(str))
}

// Convert a simple symbol to an elk symbol value
func S(sym symbol.Symbol) Symbol {
	return Symbol{Id: sym}
}

var SymbolClass *Class // ::Std::Symbol

func (s Symbol) ToValue() Value {
	return Value{
		flag: SYMBOL_FLAG,
		data: *(*uintptr)(unsafe.Pointer(&s)),
	}
}

func (Symbol) Class() *Class {
	return SymbolClass
}

func (Symbol) DirectClass() *Class {
	return SymbolClass
}

func (Symbol) SingletonClass() *Class {
	return nil
}

func (s Symbol) String() string {
	return s.Id.String()
}

func (s Symbol) ToString() String {
	return String(s.String())
}

func (s Symbol) InspectContent() string {
	return s.Id.InspectContent()
}

func (s Symbol) Inspect() string {
	return s.Id.Inspect()
}

func (s Symbol) Error() string {
	return s.Inspect()
}

func (s Symbol) InstanceVariables() *InstanceVariables {
	return nil
}

// Check whether s is equal to other
func (s Symbol) EqualVal(other Value) Value {
	if other.IsInlineSymbol() {
		return BoolVal(s.Id.EqualSymbol(other.AsInlineSymbol().Id))
	}

	return False.ToValue()
}

// Check whether s is equal to other
func (s Symbol) Equal(other Value) bool {
	if other.IsInlineSymbol() {
		return s.Id.EqualSymbol(other.AsInlineSymbol().Id)
	}

	return false
}

func (s Symbol) LaxEqual(other Value) bool {
	return s.Equal(other)
}

// Check whether s is equal to other
func (s Symbol) StrictEqualVal(other Value) Value {
	return s.EqualVal(other)
}

// Check whether s is equal to other
func (s Symbol) LaxEqualVal(other Value) Value {
	return s.EqualVal(other)
}

func (s Symbol) Hash() UInt64 {
	return UInt64(s.Id.Hash())
}

func initSymbol() {
	SymbolClass = NewClass()
	StdModule.AddConstantString("Symbol", Ref(SymbolClass))
	RegisterNativeClass("Std::Symbol", "value.SymbolClass")
}
