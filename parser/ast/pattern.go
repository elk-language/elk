package ast

import "slices"

// All nodes that should be valid in pattern matching should
// implement this interface
type PatternNode interface {
	Node
	patternNode()
}

func (*InvalidNode) patternNode()                    {}
func (*AsPatternNode) patternNode()                  {}
func (*BinHashSetLiteralNode) patternNode()          {}
func (*BinArrayTupleLiteralNode) patternNode()       {}
func (*BinArrayListLiteralNode) patternNode()        {}
func (*HexHashSetLiteralNode) patternNode()          {}
func (*HexArrayTupleLiteralNode) patternNode()       {}
func (*HexArrayListLiteralNode) patternNode()        {}
func (*SymbolHashSetLiteralNode) patternNode()       {}
func (*SymbolArrayTupleLiteralNode) patternNode()    {}
func (*SymbolArrayListLiteralNode) patternNode()     {}
func (*WordHashSetLiteralNode) patternNode()         {}
func (*WordArrayTupleLiteralNode) patternNode()      {}
func (*WordArrayListLiteralNode) patternNode()       {}
func (*SymbolKeyValuePatternNode) patternNode()      {}
func (*KeyValuePatternNode) patternNode()            {}
func (*ObjectPatternNode) patternNode()              {}
func (*RecordPatternNode) patternNode()              {}
func (*MapPatternNode) patternNode()                 {}
func (*RestPatternNode) patternNode()                {}
func (*SetPatternNode) patternNode()                 {}
func (*ListPatternNode) patternNode()                {}
func (*TuplePatternNode) patternNode()               {}
func (*ConstantLookupNode) patternNode()             {}
func (*PublicConstantNode) patternNode()             {}
func (*PrivateConstantNode) patternNode()            {}
func (*GenericConstantNode) patternNode()            {}
func (*PublicIdentifierNode) patternNode()           {}
func (*PrivateIdentifierNode) patternNode()          {}
func (*RangeLiteralNode) patternNode()               {}
func (*BinaryPatternNode) patternNode()              {}
func (*UnaryExpressionNode) patternNode()            {}
func (*TrueLiteralNode) patternNode()                {}
func (*FalseLiteralNode) patternNode()               {}
func (*NilLiteralNode) patternNode()                 {}
func (*CharLiteralNode) patternNode()                {}
func (*RawCharLiteralNode) patternNode()             {}
func (*DoubleQuotedStringLiteralNode) patternNode()  {}
func (*InterpolatedStringLiteralNode) patternNode()  {}
func (*RawStringLiteralNode) patternNode()           {}
func (*SimpleSymbolLiteralNode) patternNode()        {}
func (*InterpolatedSymbolLiteralNode) patternNode()  {}
func (*IntLiteralNode) patternNode()                 {}
func (*Int64LiteralNode) patternNode()               {}
func (*UInt64LiteralNode) patternNode()              {}
func (*Int32LiteralNode) patternNode()               {}
func (*UInt32LiteralNode) patternNode()              {}
func (*Int16LiteralNode) patternNode()               {}
func (*UInt16LiteralNode) patternNode()              {}
func (*Int8LiteralNode) patternNode()                {}
func (*UInt8LiteralNode) patternNode()               {}
func (*FloatLiteralNode) patternNode()               {}
func (*Float32LiteralNode) patternNode()             {}
func (*Float64LiteralNode) patternNode()             {}
func (*BigFloatLiteralNode) patternNode()            {}
func (*UninterpolatedRegexLiteralNode) patternNode() {}
func (*InterpolatedRegexLiteralNode) patternNode()   {}
func (*UnquoteNode) patternNode()                    {}

func anyPatternDeclaresVariables(patterns []PatternNode) bool {
	return slices.ContainsFunc(patterns, PatternDeclaresVariables)
}

func PatternDeclaresVariables(pattern PatternNode) bool {
	switch pat := pattern.(type) {
	case *PublicIdentifierNode, *PrivateIdentifierNode, *AsPatternNode:
		return true
	case *BinaryPatternNode:
		return PatternDeclaresVariables(pat.Left) ||
			PatternDeclaresVariables(pat.Right)
	case *ObjectPatternNode:
		return anyPatternDeclaresVariables(pat.Attributes)
	case *SymbolKeyValuePatternNode:
		return PatternDeclaresVariables(pat.Value)
	case *KeyValuePatternNode:
		return PatternDeclaresVariables(pat.Value)
	case *MapPatternNode:
		return anyPatternDeclaresVariables(pat.Elements)
	case *RecordPatternNode:
		return anyPatternDeclaresVariables(pat.Elements)
	case *ListPatternNode:
		return anyPatternDeclaresVariables(pat.Elements)
	case *TuplePatternNode:
		return anyPatternDeclaresVariables(pat.Elements)
	case *RestPatternNode:
		switch pat.Identifier.(type) {
		case *PrivateIdentifierNode, *PublicIdentifierNode:
			return true
		}
		return false
	default:
		return false
	}
}

type PatternExpressionNode interface {
	Node
	ExpressionNode
	PatternNode
}
