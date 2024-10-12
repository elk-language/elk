package symbol

import (
	"slices"

	"github.com/elk-language/elk/value"
)

var (
	Std                  = value.ToSymbol("Std")
	Object               = value.ToSymbol("Object")
	Class                = value.ToSymbol("Class")
	Mixin                = value.ToSymbol("Mixin")
	Module               = value.ToSymbol("Module")
	Interface            = value.ToSymbol("Interface")
	Value                = value.ToSymbol("Value")
	String               = value.ToSymbol("String")
	Symbol               = value.ToSymbol("Symbol")
	Char                 = value.ToSymbol("Char")
	Float                = value.ToSymbol("Float")
	Float64              = value.ToSymbol("Float64")
	Float32              = value.ToSymbol("Float32")
	BigFloat             = value.ToSymbol("BigFloat")
	Int                  = value.ToSymbol("Int")
	Int64                = value.ToSymbol("Int64")
	Int32                = value.ToSymbol("Int32")
	Int16                = value.ToSymbol("Int16")
	Int8                 = value.ToSymbol("Int8")
	UInt64               = value.ToSymbol("UInt64")
	UInt32               = value.ToSymbol("UInt32")
	UInt16               = value.ToSymbol("UInt16")
	UInt8                = value.ToSymbol("UInt8")
	Bool                 = value.ToSymbol("Bool")
	False                = value.ToSymbol("False")
	True                 = value.ToSymbol("True")
	Nil                  = value.ToSymbol("Nil")
	Method               = value.ToSymbol("Method")
	Regex                = value.ToSymbol("Regex")
	ArrayList            = value.ToSymbol("ArrayList")
	List                 = value.ToSymbol("List")
	ArrayTuple           = value.ToSymbol("ArrayTuple")
	Tuple                = value.ToSymbol("Tuple")
	HashMap              = value.ToSymbol("HashMap")
	Map                  = value.ToSymbol("Map")
	HashRecord           = value.ToSymbol("HashRecord")
	Record               = value.ToSymbol("Record")
	HashSet              = value.ToSymbol("HashSet")
	Set                  = value.ToSymbol("Set")
	Pair                 = value.ToSymbol("Pair")
	StringConvertible    = value.ToSymbol("StringConvertible")
	Inspectable          = value.ToSymbol("Inspectable")
	AnyInt               = value.ToSymbol("AnyInt")
	Kernel               = value.ToSymbol("Kernel")
	Range                = value.ToSymbol("Range")
	BeginlessClosedRange = value.ToSymbol("BeginlessClosedRange")
	BeginlessOpenRange   = value.ToSymbol("BeginlessOpenRange")
	ClosedRange          = value.ToSymbol("ClosedRange")
	EndlessClosedRange   = value.ToSymbol("EndlessClosedRange")
	EndlessOpenRange     = value.ToSymbol("EndlessOpenRange")
	LeftOpenRange        = value.ToSymbol("LeftOpenRange")
	OpenRange            = value.ToSymbol("OpenRange")
	RightOpenRange       = value.ToSymbol("RightOpenRange")
	Comparable           = value.ToSymbol("Comparable")
)

// lowercase symbols
var (
	L_call     = value.ToSymbol("call")
	L_self     = value.ToSymbol("self")
	L_contains = value.ToSymbol("contains")
	L_length   = value.ToSymbol("length")
)

// special symbols
var (
	S_empty    = value.ToSymbol("")
	S_init     = value.ToSymbol("#init")
	S_contains = value.ToSymbol("#contains")
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
	OpNotEqual             = value.ToSymbol("!=")  // `!=`
	OpLaxEqual             = value.ToSymbol("=~")  // `=~`
	OpLaxNotEqual          = value.ToSymbol("!~")  // `!~`
	OpStrictEqual          = value.ToSymbol("===") // `===`
	OpStrictNotEqual       = value.ToSymbol("!==") // `!==`
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

func IsEqualityOperator(methodName value.Symbol) bool {
	switch methodName {
	case OpEqual, OpNotEqual,
		OpLaxEqual, OpLaxNotEqual,
		OpStrictEqual, OpStrictNotEqual:
		return true
	default:
		return false
	}
}

func IsRelationalOperator(methodName value.Symbol) bool {
	switch methodName {
	case OpGreaterThan, OpGreaterThanEqual,
		OpLessThan, OpLessThanEqual:
		return true
	default:
		return false
	}
}

func RequiresNoParameters(methodName value.Symbol) bool {
	switch methodName {
	case OpIncrement, OpDecrement, OpNegate, OpUnaryPlus, OpBitwiseNot:
		return true
	default:
		return false
	}
}

func RequiresOneParameter(methodName value.Symbol) bool {
	switch methodName {
	case OpAdd, OpSubtract, OpMultiply,
		OpDivide, OpExponentiate, OpLogicalRightBitshift,
		OpLogicalLeftBitshift, OpRightBitshift, OpLeftBitshift,
		OpLessThan, OpLessThanEqual, OpGreaterThan, OpGreaterThanEqual,
		OpStrictEqual, OpStrictNotEqual, OpLaxEqual, OpLaxNotEqual,
		OpEqual, OpNotEqual, OpModulo, OpSpaceship, OpXor,
		OpOr, OpAnd, OpAndNot, OpSubscript:
		return true
	default:
		return false
	}
}

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
