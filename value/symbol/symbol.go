package symbol

import (
	"slices"

	"github.com/elk-language/elk/value"
)

var (
	Std               = value.ToSymbol("Std")
	Object            = value.ToSymbol("Object")
	Class             = value.ToSymbol("Class")
	Mixin             = value.ToSymbol("Mixin")
	Module            = value.ToSymbol("Module")
	Interface         = value.ToSymbol("Interface")
	Value             = value.ToSymbol("Value")
	String            = value.ToSymbol("String")
	Symbol            = value.ToSymbol("Symbol")
	Char              = value.ToSymbol("Char")
	Float             = value.ToSymbol("Float")
	Float64           = value.ToSymbol("Float64")
	Float32           = value.ToSymbol("Float32")
	BigFloat          = value.ToSymbol("BigFloat")
	Int               = value.ToSymbol("Int")
	Int64             = value.ToSymbol("Int64")
	Int32             = value.ToSymbol("Int32")
	Int16             = value.ToSymbol("Int16")
	Int8              = value.ToSymbol("Int8")
	UInt64            = value.ToSymbol("UInt64")
	UInt32            = value.ToSymbol("UInt32")
	UInt16            = value.ToSymbol("UInt16")
	UInt8             = value.ToSymbol("UInt8")
	Bool              = value.ToSymbol("Bool")
	False             = value.ToSymbol("False")
	True              = value.ToSymbol("True")
	Nil               = value.ToSymbol("Nil")
	Method            = value.ToSymbol("Method")
	Regex             = value.ToSymbol("Regex")
	ArrayList         = value.ToSymbol("ArrayList")
	ArrayTuple        = value.ToSymbol("ArrayTuple")
	HashMap           = value.ToSymbol("HashMap")
	HashRecord        = value.ToSymbol("HashRecord")
	HashSet           = value.ToSymbol("HashSet")
	Pair              = value.ToSymbol("Pair")
	StringConvertible = value.ToSymbol("StringConvertible")
	Inspectable       = value.ToSymbol("Inspectable")
	AnyInt            = value.ToSymbol("AnyInt")
)

var (
	Empty  = value.ToSymbol("")
	M_init = value.ToSymbol("#init")
	M_call = value.ToSymbol("call")
	M_self = value.ToSymbol("self")
)

var (
	OpIncrement            = value.ToSymbol("++")  // `++`
	OpDecrement            = value.ToSymbol("--")  // `--`
	OpSubscriptSet         = value.ToSymbol("[]=") // `[]=`
	OpSubscript            = value.ToSymbol("[]")  // `[]`
	OpNegate               = value.ToSymbol("-@")  // `-@`
	OpUnaryPlus            = value.ToSymbol("+@")  // `+@`
	OpBitwiseNot           = value.ToSymbol("~")   // `~`
	OpAnd                  = value.ToSymbol("&")   // `&`
	OpAndNot               = value.ToSymbol("&~")  // `&~`
	OpOr                   = value.ToSymbol("|")   // `|`
	OpXor                  = value.ToSymbol("^")   // `^`
	OpSpaceship            = value.ToSymbol("<=>") // `<=>`
	OpModulo               = value.ToSymbol("%")   // `%`
	OpEqual                = value.ToSymbol("==")  // `==`
	OpLaxEqual             = value.ToSymbol("=~")  // `=~`
	OpStrictEqual          = value.ToSymbol("===") // `===`
	OpGreaterThan          = value.ToSymbol(">")   // `>`
	OpGreaterThanEqual     = value.ToSymbol(">=")  // `>=`
	OpLessThan             = value.ToSymbol("<")   // `<`
	OpLessThanEqual        = value.ToSymbol("<=")  // `<=`
	OpLeftBitshift         = value.ToSymbol("<<")  // `<<`
	OpLogicalLeftBitshift  = value.ToSymbol("<<<") // `<<<`
	OpRightBitshift        = value.ToSymbol(">>")  // `>>`
	OpLogicalRightBitshift = value.ToSymbol(">>>") // `>>>`
	OpAdd                  = value.ToSymbol("+")   // `+`
	OpSubtract             = value.ToSymbol("-")   // `-`
	OpMultiply             = value.ToSymbol("*")   // `*`
	OpDivide               = value.ToSymbol("/")   // `/`
	OpExponentiate         = value.ToSymbol("**")  // `**`
)

func SortKeys[V any](m map[value.Symbol]V) []value.Symbol {
	keys := make([]value.Symbol, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.SortFunc(keys, func(a, b value.Symbol) int {
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
