package symbol

import "github.com/elk-language/elk/value"

var (
	Std      = value.ToSymbol("Std")
	String   = value.ToSymbol("String")
	Symbol   = value.ToSymbol("Symbol")
	Char     = value.ToSymbol("Char")
	Float    = value.ToSymbol("Float")
	Float64  = value.ToSymbol("Float64")
	Float32  = value.ToSymbol("Float32")
	BigFloat = value.ToSymbol("BigFloat")
	Int      = value.ToSymbol("Int")
	Int64    = value.ToSymbol("Int64")
	Int32    = value.ToSymbol("Int32")
	Int16    = value.ToSymbol("Int16")
	Int8     = value.ToSymbol("Int8")
	UInt64   = value.ToSymbol("UInt64")
	UInt32   = value.ToSymbol("UInt32")
	UInt16   = value.ToSymbol("UInt16")
	UInt8    = value.ToSymbol("UInt8")
	False    = value.ToSymbol("False")
	True     = value.ToSymbol("True")
	Nil      = value.ToSymbol("Nil")
	Method   = value.ToSymbol("Method")
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
