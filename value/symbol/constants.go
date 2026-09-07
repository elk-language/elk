package symbol

import (
	"slices"
)

var (
	C_Root                 = ToSymbol("Root")
	C_Std                  = ToSymbol("Std")
	C_Object               = ToSymbol("Object")
	C_Class                = ToSymbol("Class")
	C_Mixin                = ToSymbol("Mixin")
	C_Module               = ToSymbol("Module")
	C_Interface            = ToSymbol("Interface")
	C_Value                = ToSymbol("Value")
	C_String               = ToSymbol("String")
	C_Symbol               = ToSymbol("Symbol")
	C_Char                 = ToSymbol("Char")
	C_Float                = ToSymbol("Float")
	C_Float64              = ToSymbol("Float64")
	C_Float32              = ToSymbol("Float32")
	C_BigFloat             = ToSymbol("BigFloat")
	C_Int                  = ToSymbol("Int")
	C_Int64                = ToSymbol("Int64")
	C_Int32                = ToSymbol("Int32")
	C_Int16                = ToSymbol("Int16")
	C_Int8                 = ToSymbol("Int8")
	C_UInt                 = ToSymbol("UInt")
	C_UInt64               = ToSymbol("UInt64")
	C_UInt32               = ToSymbol("UInt32")
	C_UInt16               = ToSymbol("UInt16")
	C_UInt8                = ToSymbol("UInt8")
	C_ImmutableBox         = ToSymbol("ImmutableBox")
	C_Box                  = ToSymbol("Box")
	C_Bool                 = ToSymbol("Bool")
	C_False                = ToSymbol("False")
	C_True                 = ToSymbol("True")
	C_Nil                  = ToSymbol("Nil")
	C_Method               = ToSymbol("Method")
	C_Regex                = ToSymbol("Regex")
	C_ArrayList            = ToSymbol("ArrayList")
	C_List                 = ToSymbol("List")
	C_ArrayTuple           = ToSymbol("ArrayTuple")
	C_Tuple                = ToSymbol("Tuple")
	C_HashMap              = ToSymbol("HashMap")
	C_Map                  = ToSymbol("Map")
	C_HashRecord           = ToSymbol("HashRecord")
	C_Record               = ToSymbol("Record")
	C_HashSet              = ToSymbol("HashSet")
	C_Set                  = ToSymbol("Set")
	C_Pair                 = ToSymbol("Pair")
	C_Result               = ToSymbol("Result")
	C_Convertible          = ToSymbol("Convertible")
	C_Inspectable          = ToSymbol("Inspectable")
	C_AnyInt               = ToSymbol("AnyInt")
	C_Kernel               = ToSymbol("Kernel")
	C_Range                = ToSymbol("Range")
	C_BeginlessClosedRange = ToSymbol("BeginlessClosedRange")
	C_BeginlessOpenRange   = ToSymbol("BeginlessOpenRange")
	C_ClosedRange          = ToSymbol("ClosedRange")
	C_EndlessClosedRange   = ToSymbol("EndlessClosedRange")
	C_EndlessOpenRange     = ToSymbol("EndlessOpenRange")
	C_LeftOpenRange        = ToSymbol("LeftOpenRange")
	C_OpenRange            = ToSymbol("OpenRange")
	C_RightOpenRange       = ToSymbol("RightOpenRange")
	C_Comparable           = ToSymbol("Comparable")
	C_Generator            = ToSymbol("Generator")
	C_Iterable             = ToSymbol("Iterable")
	C_PrimitiveIterable    = ToSymbol("PrimitiveIterable")
	C_Error                = ToSymbol("Error")
	C_Thread               = ToSymbol("Thread")
	C_ThreadPool           = ToSymbol("ThreadPool")
	C_Channel              = ToSymbol("Channel")
	C_ReadChannel          = ToSymbol("ReadChannel")
	C_WriteChannel         = ToSymbol("WriteChannel")
	C_StackTrace           = ToSymbol("StackTrace")
	C_CallFrame            = ToSymbol("CallFrame")
	C_Promise              = ToSymbol("Promise")
	C_Key                  = ToSymbol("Key")
	C_Elk                  = ToSymbol("Elk")
	C_Type                 = ToSymbol("Type")
	C_Checker              = ToSymbol("Checker")
	C_Macro                = ToSymbol("Macro")
	C_Token                = ToSymbol("Token")
	C_Node                 = ToSymbol("Node")
	C_ExpressionNode       = ToSymbol("ExpressionNode")
	C_ConstantNode         = ToSymbol("ConstantNode")
	C_ComplexConstantNode  = ToSymbol("ComplexConstantNode")
	C_PatternNode          = ToSymbol("PatternNode")
	C_LiteralPatternNode   = ToSymbol("LiteralPatternNode")
	C_TypeNode             = ToSymbol("TypeNode")
	C_IdentifierNode       = ToSymbol("IdentifierNode")
	C_InstanceVariableNode = ToSymbol("InstanceVariableNode")
	C_AST                  = ToSymbol("AST")
	C_Test                 = ToSymbol("Test")
	C_Time                 = ToSymbol("Time")
	C_Date                 = ToSymbol("Date")
	C_DateTime             = ToSymbol("DateTime")
	C_Span                 = ToSymbol("Span")
	C_Closure              = ToSymbol("Closure")
	C_Sync                 = ToSymbol("Sync")
	C_WaitGroup            = ToSymbol("WaitGroup")
)

// lowercase symbols
var (
	L_to_string                 = ToSymbol("to_string")
	L_at                        = ToSymbol("at")
	L_call                      = ToSymbol("call")
	L_self                      = ToSymbol("self")
	L_contains                  = ToSymbol("contains")
	L_length                    = ToSymbol("length")
	L_hash                      = ToSymbol("hash")
	L_iter                      = ToSymbol("iter")
	L_inspect                   = ToSymbol("inspect")
	L_next                      = ToSymbol("next")
	L_stop_iteration            = ToSymbol("stop_iteration")
	L_channel_closed            = ToSymbol("channel_closed")
	L_colorize                  = ToSymbol("colorize")
	L_to_ast_expr_node          = ToSymbol("to_ast_expr_node")
	L_to_ast_const_node         = ToSymbol("to_ast_const_node")
	L_to_ast_complex_const_node = ToSymbol("to_ast_complex_const_node")
	L_to_ast_pattern_node       = ToSymbol("to_ast_pattern_node")
	L_to_ast_pattern_expr_node  = ToSymbol("to_ast_pattern_expr_node")
	L_to_ast_type_node          = ToSymbol("to_ast_type_node")
	L_to_ast_ident_node         = ToSymbol("to_ast_ident_node")
	L_to_ast_ivar_node          = ToSymbol("to_ast_ivar_node")
	L_message                   = ToSymbol("message")
	L_matches                   = ToSymbol("matches")
	L_remove                    = ToSymbol("remove")
	L_push                      = ToSymbol("push")
	L_view                      = ToSymbol("view")
	L_slice                     = ToSymbol("slice")
	L_diagnostics               = ToSymbol("diagnostics")
	L_source_map                = ToSymbol("source_map")
)

// special symbols
var (
	S_empty                    = ToSymbol("")                          // empty symbol
	S_splice                   = ToSymbol("#splice")                   // #splice
	S_init                     = ToSymbol("#init")                     // #init
	S_contains                 = ToSymbol("#contains")                 // #contains
	S_BuiltinAddable           = ToSymbol("#BuiltinAddable")           // #BuiltinAddable
	S_BuiltinSubtractable      = ToSymbol("#BuiltinSubtractable")      // #BuiltinSubtractable
	S_BuiltinMultipliable      = ToSymbol("#BuiltinMultipliable")      // #BuiltinMultipliable
	S_BuiltinDividable         = ToSymbol("#BuiltinDividable")         // #BuiltinDividable
	S_BuiltinNumeric           = ToSymbol("#BuiltinNumeric")           // #BuiltinNumeric
	S_BuiltinInt               = ToSymbol("#BuiltinInt")               // #BuiltinInt
	S_BuiltinLogicBitshiftable = ToSymbol("#BuiltinLogicBitshiftable") // #BuiltinLogicBitshiftable
	S_BuiltinEquatable         = ToSymbol("#BuiltinEquatable")         // #BuiltinEquatable
	S_BuiltinIterable          = ToSymbol("#BuiltinIterable")          // #BuiltinIterable
	S_BuiltinIterator          = ToSymbol("#BuiltinIterator")          // #BuiltinIterator
	S_BuiltinIncrementable     = ToSymbol("#BuiltinIncrementable")     // #BuiltinIncrementable
	S_BuiltinSubscriptable     = ToSymbol("#BuiltinSubscriptable")     // #BuiltinSubscriptable
)

var (
	OpIncrement            = ToSymbol("++")  // `++`
	OpDecrement            = ToSymbol("--")  // `--`
	OpSubscriptSet         = ToSymbol("[]=") // `[]=`
	OpSubscript            = ToSymbol("[]")  // `[]`
	OpPop                  = ToSymbol("<<@") // `<<@`
	OpNegate               = ToSymbol("-@")  // `-@`
	OpUnaryPlus            = ToSymbol("+@")  // `+@`
	OpBitwiseNot           = ToSymbol("~")   // `~`
	OpAnd                  = ToSymbol("&")   // `&`
	OpAndNot               = ToSymbol("&~")  // `&~`
	OpOr                   = ToSymbol("|")   // `|`
	OpXor                  = ToSymbol("^")   // `^`
	OpSpaceship            = ToSymbol("<=>") // `<=>`
	OpModulo               = ToSymbol("%")   // `%`
	OpEqual                = ToSymbol("==")  // `==`
	OpNotEqual             = ToSymbol("!=")  // `!=`
	OpLaxEqual             = ToSymbol("=~")  // `=~`
	OpLaxNotEqual          = ToSymbol("!~")  // `!~`
	OpStrictEqual          = ToSymbol("===") // `===`
	OpStrictNotEqual       = ToSymbol("!==") // `!==`
	OpGreaterThan          = ToSymbol(">")   // `>`
	OpGreaterThanEqual     = ToSymbol(">=")  // `>=`
	OpLessThan             = ToSymbol("<")   // `<`
	OpLessThanEqual        = ToSymbol("<=")  // `<=`
	OpLeftBitshift         = ToSymbol("<<")  // `<<`
	OpLogicalLeftBitshift  = ToSymbol("<<<") // `<<<`
	OpRightBitshift        = ToSymbol(">>")  // `>>`
	OpLogicalRightBitshift = ToSymbol(">>>") // `>>>`
	OpAdd                  = ToSymbol("+")   // `+`
	OpSubtract             = ToSymbol("-")   // `-`
	OpMultiply             = ToSymbol("*")   // `*`
	OpDivide               = ToSymbol("/")   // `/`
	OpExponentiate         = ToSymbol("**")  // `**`
)

func IsLaxEqualityOperator(methodName Symbol) bool {
	switch methodName {
	case OpLaxEqual, OpLaxNotEqual:
		return true
	default:
		return false
	}
}

func IsEqualityOperator(methodName Symbol) bool {
	switch methodName {
	case OpEqual, OpNotEqual,
		OpLaxEqual, OpLaxNotEqual,
		OpStrictEqual, OpStrictNotEqual:
		return true
	default:
		return false
	}
}

func IsRelationalOperator(methodName Symbol) bool {
	switch methodName {
	case OpGreaterThan, OpGreaterThanEqual,
		OpLessThan, OpLessThanEqual:
		return true
	default:
		return false
	}
}

func RequiresNoParameters(methodName Symbol) bool {
	switch methodName {
	case OpIncrement, OpDecrement, OpNegate, OpUnaryPlus, OpBitwiseNot:
		return true
	default:
		return false
	}
}

func RequiresOneParameter(methodName Symbol) bool {
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

func SortKeys[V any](m map[Symbol]V) []Symbol {
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
