package ast

// All nodes that should be valid in type annotations should
// implement this interface
type TypeNode interface {
	Node
	typeNode()
}

func (*InvalidNode) typeNode()                   {}
func (*SelfLiteralNode) typeNode()               {}
func (*AnyTypeNode) typeNode()                   {}
func (*NeverTypeNode) typeNode()                 {}
func (*VoidTypeNode) typeNode()                  {}
func (*UnionTypeNode) typeNode()                 {}
func (*IntersectionTypeNode) typeNode()          {}
func (*BinaryTypeNode) typeNode()                {}
func (*NilableTypeNode) typeNode()               {}
func (*InstanceOfTypeNode) typeNode()            {}
func (*SingletonTypeNode) typeNode()             {}
func (*ClosureTypeNode) typeNode()               {}
func (*NotTypeNode) typeNode()                   {}
func (*UnaryTypeNode) typeNode()                 {}
func (*PublicConstantNode) typeNode()            {}
func (*PrivateConstantNode) typeNode()           {}
func (*ConstantLookupNode) typeNode()            {}
func (*GenericConstantNode) typeNode()           {}
func (*NilLiteralNode) typeNode()                {}
func (*BoolLiteralNode) typeNode()               {}
func (*TrueLiteralNode) typeNode()               {}
func (*FalseLiteralNode) typeNode()              {}
func (*CharLiteralNode) typeNode()               {}
func (*RawCharLiteralNode) typeNode()            {}
func (*RawStringLiteralNode) typeNode()          {}
func (*InterpolatedStringLiteralNode) typeNode() {}
func (*DoubleQuotedStringLiteralNode) typeNode() {}
func (*SimpleSymbolLiteralNode) typeNode()       {}
func (*InterpolatedSymbolLiteralNode) typeNode() {}
func (*IntLiteralNode) typeNode()                {}
func (*Int64LiteralNode) typeNode()              {}
func (*Int32LiteralNode) typeNode()              {}
func (*Int16LiteralNode) typeNode()              {}
func (*Int8LiteralNode) typeNode()               {}
func (*UInt64LiteralNode) typeNode()             {}
func (*UInt32LiteralNode) typeNode()             {}
func (*UInt16LiteralNode) typeNode()             {}
func (*UInt8LiteralNode) typeNode()              {}
func (*FloatLiteralNode) typeNode()              {}
func (*Float64LiteralNode) typeNode()            {}
func (*Float32LiteralNode) typeNode()            {}
func (*BigFloatLiteralNode) typeNode()           {}
func (*UnquoteNode) typeNode()                   {}
func (*MacroCallNode) typeNode()                 {}
func (*ReceiverlessMacroCallNode) typeNode()     {}
func (*ScopedMacroCallNode) typeNode()           {}
func (*MacroBoundaryNode) typeNode()             {}

type StringTypeNode interface {
	Node
	TypeNode
	StringLiteralNode
}

type StringOrSymbolTypeNode interface {
	Node
	TypeNode
	StringOrSymbolLiteralNode
}
