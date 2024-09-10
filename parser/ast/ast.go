// Package ast defines types
// used by the Elk parser.
//
// All the nodes of the Abstract Syntax Tree
// constructed by the Elk parser are defined in this package.
package ast

import (
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/regex/flag"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value/symbol"
)

// Checks whether all expressions in the given list are static.
func isExpressionSliceStatic(elements []ExpressionNode) bool {
	for _, element := range elements {
		if !element.IsStatic() {
			return false
		}
	}
	return true
}

// Checks whether all expressions in the given list are static.
func areExpressionsStatic(elements ...ExpressionNode) bool {
	for _, element := range elements {
		if element != nil && !element.IsStatic() {
			return false
		}
	}
	return true
}

// Checks whether all nodes in the given list are static.
func areNodesStatic(elements ...Node) bool {
	for _, element := range elements {
		if element != nil && !element.IsStatic() {
			return false
		}
	}
	return true
}

// Turn an expression to a statement
func ExpressionToStatement(expr ExpressionNode) StatementNode {
	return NewExpressionStatementNode(expr.Span(), expr)
}

// Turn an expression to a collection of statements.
func ExpressionToStatements(expr ExpressionNode) []StatementNode {
	return []StatementNode{ExpressionToStatement(expr)}
}

// Every node type implements this interface.
type Node interface {
	position.SpanInterface
	IsStatic() bool // Value is known at compile-time
	Type(*types.GlobalEnvironment) types.Type
	SetType(types.Type)
	SkipTypechecking() bool
}

type DocCommentableNode interface {
	DocComment() string
	SetDocComment(string)
}

type DocCommentableNodeBase struct {
	comment string
}

func (d *DocCommentableNodeBase) DocComment() string {
	return d.comment
}

func (d *DocCommentableNodeBase) SetDocComment(comment string) {
	d.comment = comment
}

// Every node type implements this interface.
type TypedNode interface {
	Node
	SetType(types.Type)
}

// Base typed AST node.
type TypedNodeBase struct {
	span *position.Span
	typ  types.Type
}

func (t *TypedNodeBase) Type(*types.GlobalEnvironment) types.Type {
	return t.typ
}

func (t *TypedNodeBase) SkipTypechecking() bool {
	return t.typ != nil
}

func (t *TypedNodeBase) SetType(typ types.Type) {
	t.typ = typ
}

func (t *TypedNodeBase) Span() *position.Span {
	return t.span
}

func (t *TypedNodeBase) SetSpan(span *position.Span) {
	t.span = span
}

// Base AST node.
type NodeBase struct {
	span *position.Span
}

func (*NodeBase) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Void{}
}

func (*NodeBase) SetType(types.Type) {}

func (t *NodeBase) SkipTypechecking() bool {
	return false
}

func (n *NodeBase) Span() *position.Span {
	return n.span
}

func (n *NodeBase) SetSpan(span *position.Span) {
	n.span = span
}

// Check whether the node can be used as a left value
// in a variable/constant declaration.
func IsValidDeclarationTarget(node Node) bool {
	switch node.(type) {
	case *PrivateIdentifierNode, *PublicIdentifierNode:
		return true
	default:
		return false
	}
}

// Check whether the node can be used as a left value
// in an assignment expression.
func IsValidAssignmentTarget(node Node) bool {
	switch node.(type) {
	case *PrivateIdentifierNode, *PublicIdentifierNode,
		*AttributeAccessNode, *InstanceVariableNode, *SubscriptExpressionNode:
		return true
	default:
		return false
	}
}

// Check whether the node can be used as a range pattern element.
func IsValidRangePatternElement(node Node) bool {
	switch node.(type) {
	case *TrueLiteralNode, *FalseLiteralNode, *NilLiteralNode, *CharLiteralNode,
		*RawCharLiteralNode, *RawStringLiteralNode, *DoubleQuotedStringLiteralNode,
		*InterpolatedStringLiteralNode, *SimpleSymbolLiteralNode, *InterpolatedSymbolLiteralNode,
		*FloatLiteralNode, *Float64LiteralNode, *Float32LiteralNode, *BigFloatLiteralNode,
		*IntLiteralNode, *Int64LiteralNode, *UInt64LiteralNode, *Int32LiteralNode, *UInt32LiteralNode,
		*Int16LiteralNode, *UInt16LiteralNode, *Int8LiteralNode, *UInt8LiteralNode,
		*PublicConstantNode, *PrivateConstantNode, *ConstantLookupNode, *UnaryExpressionNode:
		return true
	default:
		return false
	}
}

// Check whether the node is a constant.
func IsConstant(node Node) bool {
	switch node.(type) {
	case *PrivateConstantNode, *PublicConstantNode:
		return true
	default:
		return false
	}
}

// Check whether the node is a complex constant.
func IsComplexConstant(node Node) bool {
	switch node.(type) {
	case *PrivateConstantNode, *PublicConstantNode, *ConstantLookupNode:
		return true
	default:
		return false
	}
}

// Check whether the node is a valid right operand of the pipe operator `|>`.
func IsValidPipeExpressionTarget(node Node) bool {
	switch node.(type) {
	case *MethodCallNode, *GenericMethodCallNode, *ReceiverlessMethodCallNode,
		*GenericReceiverlessMethodCallNode, *AttributeAccessNode, *ConstructorCallNode, *CallNode:
		return true
	default:
		return false
	}
}

// Represents a single statement, so for example
// a single valid "line" of Elk code.
// Usually its an expression optionally terminated with a newline ors semicolon.
type StatementNode interface {
	Node
	statementNode()
}

func (*InvalidNode) statementNode()             {}
func (*ExpressionStatementNode) statementNode() {}
func (*EmptyStatementNode) statementNode()      {}
func (*ImportStatementNode) statementNode()     {}

// Represents a single statement of a struct body
// optionally terminated with a newline or semicolon.
type StructBodyStatementNode interface {
	Node
	structBodyStatementNode()
}

func (*InvalidNode) structBodyStatementNode()            {}
func (*EmptyStatementNode) structBodyStatementNode()     {}
func (*ParameterStatementNode) structBodyStatementNode() {}

// All nodes that should be able to appear as
// elements of word collection literals should
// implement this interface.
type WordCollectionContentNode interface {
	Node
	ExpressionNode
	wordCollectionContentNode()
}

func (*InvalidNode) wordCollectionContentNode()          {}
func (*RawStringLiteralNode) wordCollectionContentNode() {}

// All nodes that should be able to appear as
// elements of symbol collection literals should
// implement this interface.
type SymbolCollectionContentNode interface {
	Node
	ExpressionNode
	symbolCollectionContentNode()
}

func (*InvalidNode) symbolCollectionContentNode()             {}
func (*SimpleSymbolLiteralNode) symbolCollectionContentNode() {}

// All nodes that should be able to appear as
// elements of Int collection literals should
// implement this interface.
type IntCollectionContentNode interface {
	Node
	ExpressionNode
	intCollectionContentNode()
}

func (*InvalidNode) intCollectionContentNode()    {}
func (*IntLiteralNode) intCollectionContentNode() {}

// All expression nodes implement this interface.
type ExpressionNode interface {
	Node
	expressionNode()
}

func (*InvalidNode) expressionNode()                       {}
func (*TypeExpressionNode) expressionNode()                {}
func (*InstanceVariableDeclarationNode) expressionNode()   {}
func (*VariablePatternDeclarationNode) expressionNode()    {}
func (*ValuePatternDeclarationNode) expressionNode()       {}
func (*PostfixExpressionNode) expressionNode()             {}
func (*ModifierNode) expressionNode()                      {}
func (*ModifierIfElseNode) expressionNode()                {}
func (*ModifierForInNode) expressionNode()                 {}
func (*AssignmentExpressionNode) expressionNode()          {}
func (*BinaryExpressionNode) expressionNode()              {}
func (*LogicalExpressionNode) expressionNode()             {}
func (*UnaryExpressionNode) expressionNode()               {}
func (*TrueLiteralNode) expressionNode()                   {}
func (*FalseLiteralNode) expressionNode()                  {}
func (*NilLiteralNode) expressionNode()                    {}
func (*UndefinedLiteralNode) expressionNode()              {}
func (*InstanceVariableNode) expressionNode()              {}
func (*SimpleSymbolLiteralNode) expressionNode()           {}
func (*InterpolatedSymbolLiteralNode) expressionNode()     {}
func (*IntLiteralNode) expressionNode()                    {}
func (*Int64LiteralNode) expressionNode()                  {}
func (*UInt64LiteralNode) expressionNode()                 {}
func (*Int32LiteralNode) expressionNode()                  {}
func (*UInt32LiteralNode) expressionNode()                 {}
func (*Int16LiteralNode) expressionNode()                  {}
func (*UInt16LiteralNode) expressionNode()                 {}
func (*Int8LiteralNode) expressionNode()                   {}
func (*UInt8LiteralNode) expressionNode()                  {}
func (*FloatLiteralNode) expressionNode()                  {}
func (*BigFloatLiteralNode) expressionNode()               {}
func (*Float64LiteralNode) expressionNode()                {}
func (*Float32LiteralNode) expressionNode()                {}
func (*UninterpolatedRegexLiteralNode) expressionNode()    {}
func (*InterpolatedRegexLiteralNode) expressionNode()      {}
func (*RawStringLiteralNode) expressionNode()              {}
func (*CharLiteralNode) expressionNode()                   {}
func (*RawCharLiteralNode) expressionNode()                {}
func (*DoubleQuotedStringLiteralNode) expressionNode()     {}
func (*InterpolatedStringLiteralNode) expressionNode()     {}
func (*VariableDeclarationNode) expressionNode()           {}
func (*ValueDeclarationNode) expressionNode()              {}
func (*PublicIdentifierNode) expressionNode()              {}
func (*PublicIdentifierAsNode) expressionNode()            {}
func (*PrivateIdentifierNode) expressionNode()             {}
func (*PublicConstantNode) expressionNode()                {}
func (*PublicConstantAsNode) expressionNode()              {}
func (*PrivateConstantNode) expressionNode()               {}
func (*SelfLiteralNode) expressionNode()                   {}
func (*DoExpressionNode) expressionNode()                  {}
func (*SingletonBlockExpressionNode) expressionNode()      {}
func (*SwitchExpressionNode) expressionNode()              {}
func (*IfExpressionNode) expressionNode()                  {}
func (*UnlessExpressionNode) expressionNode()              {}
func (*WhileExpressionNode) expressionNode()               {}
func (*UntilExpressionNode) expressionNode()               {}
func (*LoopExpressionNode) expressionNode()                {}
func (*NumericForExpressionNode) expressionNode()          {}
func (*ForInExpressionNode) expressionNode()               {}
func (*BreakExpressionNode) expressionNode()               {}
func (*LabeledExpressionNode) expressionNode()             {}
func (*ReturnExpressionNode) expressionNode()              {}
func (*ContinueExpressionNode) expressionNode()            {}
func (*ThrowExpressionNode) expressionNode()               {}
func (*ConstantDeclarationNode) expressionNode()           {}
func (*ConstantLookupNode) expressionNode()                {}
func (*MethodLookupNode) expressionNode()                  {}
func (*UsingAllEntryNode) expressionNode()                 {}
func (*UsingEntryWithSubentriesNode) expressionNode()      {}
func (*ClosureLiteralNode) expressionNode()                {}
func (*ClassDeclarationNode) expressionNode()              {}
func (*ModuleDeclarationNode) expressionNode()             {}
func (*MixinDeclarationNode) expressionNode()              {}
func (*InterfaceDeclarationNode) expressionNode()          {}
func (*StructDeclarationNode) expressionNode()             {}
func (*MethodDefinitionNode) expressionNode()              {}
func (*InitDefinitionNode) expressionNode()                {}
func (*MethodSignatureDefinitionNode) expressionNode()     {}
func (*GenericConstantNode) expressionNode()               {}
func (*GenericTypeDefinitionNode) expressionNode()         {}
func (*TypeDefinitionNode) expressionNode()                {}
func (*AliasDeclarationNode) expressionNode()              {}
func (*GetterDeclarationNode) expressionNode()             {}
func (*SetterDeclarationNode) expressionNode()             {}
func (*AttrDeclarationNode) expressionNode()               {}
func (*UsingExpressionNode) expressionNode()               {}
func (*IncludeExpressionNode) expressionNode()             {}
func (*ExtendExpressionNode) expressionNode()              {}
func (*EnhanceExpressionNode) expressionNode()             {}
func (*ImplementExpressionNode) expressionNode()           {}
func (*NewExpressionNode) expressionNode()                 {}
func (*GenericConstructorCallNode) expressionNode()        {}
func (*ConstructorCallNode) expressionNode()               {}
func (*SubscriptExpressionNode) expressionNode()           {}
func (*NilSafeSubscriptExpressionNode) expressionNode()    {}
func (*CallNode) expressionNode()                          {}
func (*GenericMethodCallNode) expressionNode()             {}
func (*MethodCallNode) expressionNode()                    {}
func (*GenericReceiverlessMethodCallNode) expressionNode() {}
func (*ReceiverlessMethodCallNode) expressionNode()        {}
func (*AttributeAccessNode) expressionNode()               {}
func (*KeyValueExpressionNode) expressionNode()            {}
func (*SymbolKeyValueExpressionNode) expressionNode()      {}
func (*ArrayListLiteralNode) expressionNode()              {}
func (*WordArrayListLiteralNode) expressionNode()          {}
func (*WordArrayTupleLiteralNode) expressionNode()         {}
func (*WordHashSetLiteralNode) expressionNode()            {}
func (*SymbolArrayListLiteralNode) expressionNode()        {}
func (*SymbolArrayTupleLiteralNode) expressionNode()       {}
func (*SymbolHashSetLiteralNode) expressionNode()          {}
func (*HexArrayListLiteralNode) expressionNode()           {}
func (*HexArrayTupleLiteralNode) expressionNode()          {}
func (*HexHashSetLiteralNode) expressionNode()             {}
func (*BinArrayListLiteralNode) expressionNode()           {}
func (*BinArrayTupleLiteralNode) expressionNode()          {}
func (*BinHashSetLiteralNode) expressionNode()             {}
func (*ArrayTupleLiteralNode) expressionNode()             {}
func (*HashSetLiteralNode) expressionNode()                {}
func (*HashMapLiteralNode) expressionNode()                {}
func (*HashRecordLiteralNode) expressionNode()             {}
func (*RangeLiteralNode) expressionNode()                  {}

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
func (*BinaryTypeExpressionNode) typeNode()      {}
func (*NilableTypeNode) typeNode()               {}
func (*InstanceOfTypeNode) typeNode()            {}
func (*SingletonTypeNode) typeNode()             {}
func (*ClosureTypeNode) typeNode()               {}
func (*NotTypeNode) typeNode()                   {}
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

func anyPatternDeclaresVariables(patterns []PatternNode) bool {
	for _, pat := range patterns {
		if PatternDeclaresVariables(pat) {
			return true
		}
	}
	return false
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

// All nodes that represent regexes should
// implement this interface.
type RegexLiteralNode interface {
	Node
	PatternExpressionNode
	regexLiteralNode()
}

func (*InvalidNode) regexLiteralNode()                    {}
func (*UninterpolatedRegexLiteralNode) regexLiteralNode() {}
func (*InterpolatedRegexLiteralNode) regexLiteralNode()   {}

type StringOrSymbolLiteralNode interface {
	Node
	PatternExpressionNode
	stringOrSymbolLiteralNode()
}

func (*InvalidNode) stringOrSymbolLiteralNode()                   {}
func (*InterpolatedSymbolLiteralNode) stringOrSymbolLiteralNode() {}
func (*SimpleSymbolLiteralNode) stringOrSymbolLiteralNode()       {}
func (*DoubleQuotedStringLiteralNode) stringOrSymbolLiteralNode() {}
func (*RawStringLiteralNode) stringOrSymbolLiteralNode()          {}
func (*InterpolatedStringLiteralNode) stringOrSymbolLiteralNode() {}

// All nodes that represent strings should
// implement this interface.
type StringLiteralNode interface {
	Node
	PatternExpressionNode
	StringOrSymbolLiteralNode
	stringLiteralNode()
}

func (*InvalidNode) stringLiteralNode()                   {}
func (*DoubleQuotedStringLiteralNode) stringLiteralNode() {}
func (*RawStringLiteralNode) stringLiteralNode()          {}
func (*InterpolatedStringLiteralNode) stringLiteralNode() {}

// All nodes that represent simple strings (without interpolation)
// should implement this interface.
type SimpleStringLiteralNode interface {
	Node
	ExpressionNode
	StringLiteralNode
	StringOrSymbolLiteralNode
	simpleStringLiteralNode()
}

func (*InvalidNode) simpleStringLiteralNode()                   {}
func (*DoubleQuotedStringLiteralNode) simpleStringLiteralNode() {}
func (*RawStringLiteralNode) simpleStringLiteralNode()          {}

// All nodes that should be valid in parameter declaration lists
// of methods or functions should implement this interface.
type ParameterNode interface {
	Node
	parameterNode()
	IsOptional() bool
}

func (*InvalidNode) parameterNode()            {}
func (*FormalParameterNode) parameterNode()    {}
func (*MethodParameterNode) parameterNode()    {}
func (*SignatureParameterNode) parameterNode() {}
func (*AttributeParameterNode) parameterNode() {}

// checks whether the given parameter is a positional rest parameter.
func IsPositionalRestParam(p ParameterNode) bool {
	switch param := p.(type) {
	case *MethodParameterNode:
		return param.Kind == PositionalRestParameterKind
	case *FormalParameterNode:
		return param.Kind == PositionalRestParameterKind
	case *SignatureParameterNode:
		return param.Kind == PositionalRestParameterKind
	default:
		return false
	}
}

// checks whether the given parameter is a named rest parameter.
func IsNamedRestParam(p ParameterNode) bool {
	switch param := p.(type) {
	case *MethodParameterNode:
		return param.Kind == NamedRestParameterKind
	case *FormalParameterNode:
		return param.Kind == NamedRestParameterKind
	case *SignatureParameterNode:
		return param.Kind == NamedRestParameterKind
	default:
		return false
	}
}

// Represents a type variable in generics like `class Foo[+V]; end`
type TypeParameterNode interface {
	Node
	typeVariableNode()
}

func (*InvalidNode) typeVariableNode()              {}
func (*VariantTypeParameterNode) typeVariableNode() {}

// All nodes that should be valid in constant lookups
// should implement this interface.
type ComplexConstantNode interface {
	TypedNode
	TypeNode
	ExpressionNode
	PatternNode
	PatternExpressionNode
	complexConstantNode()
}

func (*InvalidNode) complexConstantNode()         {}
func (*PublicConstantNode) complexConstantNode()  {}
func (*PrivateConstantNode) complexConstantNode() {}
func (*ConstantLookupNode) complexConstantNode()  {}
func (*GenericConstantNode) complexConstantNode() {}
func (*NilLiteralNode) complexConstantNode()      {}

// Represents all nodes that are valid in using declarations
type UsingEntryNode interface {
	TypedNode
	ExpressionNode
	usingEntryNode()
}

func (*InvalidNode) usingEntryNode()                  {}
func (*PublicConstantNode) usingEntryNode()           {}
func (*PrivateConstantNode) usingEntryNode()          {}
func (*ConstantLookupNode) usingEntryNode()           {}
func (*MethodLookupNode) usingEntryNode()             {}
func (*UsingAllEntryNode) usingEntryNode()            {}
func (*UsingEntryWithSubentriesNode) usingEntryNode() {}

type UsingSubentryNode interface {
	TypedNode
	ExpressionNode
	usingSubentryNode()
}

func (*InvalidNode) usingSubentryNode()            {}
func (*PublicConstantNode) usingSubentryNode()     {}
func (*PublicConstantAsNode) usingSubentryNode()   {}
func (*PublicIdentifierNode) usingSubentryNode()   {}
func (*PublicIdentifierAsNode) usingSubentryNode() {}

// All nodes that should be valid constants
// should implement this interface.
type ConstantNode interface {
	Node
	TypeNode
	ExpressionNode
	UsingEntryNode
	ComplexConstantNode
	constantNode()
}

func (*InvalidNode) constantNode()         {}
func (*PublicConstantNode) constantNode()  {}
func (*PrivateConstantNode) constantNode() {}

// All nodes that should be valid identifiers
// should implement this interface.
type IdentifierNode interface {
	Node
	PatternExpressionNode
	identifierNode()
}

func (*InvalidNode) identifierNode()           {}
func (*PublicIdentifierNode) identifierNode()  {}
func (*PrivateIdentifierNode) identifierNode() {}

// Nodes that implement this interface can appear
// inside of a String literal.
type StringLiteralContentNode interface {
	Node
	stringLiteralContentNode()
}

func (*InvalidNode) stringLiteralContentNode()                     {}
func (*StringInspectInterpolationNode) stringLiteralContentNode()  {}
func (*StringInterpolationNode) stringLiteralContentNode()         {}
func (*StringLiteralContentSectionNode) stringLiteralContentNode() {}

// Nodes that implement this interface can appear
// inside of a Regex literal.
type RegexLiteralContentNode interface {
	Node
	regexLiteralContentNode()
}

func (*InvalidNode) regexLiteralContentNode()                    {}
func (*RegexInterpolationNode) regexLiteralContentNode()         {}
func (*RegexLiteralContentSectionNode) regexLiteralContentNode() {}

// Nodes that implement this interface represent
// symbol literals.
type SymbolLiteralNode interface {
	Node
	ExpressionNode
	symbolLiteralNode()
}

func (*InvalidNode) symbolLiteralNode()                   {}
func (*SimpleSymbolLiteralNode) symbolLiteralNode()       {}
func (*InterpolatedSymbolLiteralNode) symbolLiteralNode() {}

// Nodes that implement this interface represent
// named arguments in method calls.
type NamedArgumentNode interface {
	Node
	namedArgumentNode()
}

func (*InvalidNode) namedArgumentNode()           {}
func (*NamedCallArgumentNode) namedArgumentNode() {}

type ProgramState uint8

const (
	UNCHECKED ProgramState = iota
	CHECKING_NAMESPACES
	CHECKED_NAMESPACES

	CHECKING_METHODS
	CHECKED_METHODS

	CHECKING_EXPRESSIONS
	CHECKED_EXPRESSIONS
)

// Represents a single Elk program (usually a single file).
type ProgramNode struct {
	NodeBase
	Body        []StatementNode
	ImportPaths []string
	State       ProgramState
}

func (*ProgramNode) IsStatic() bool {
	return false
}

// Create a new program node.
func NewProgramNode(span *position.Span, body []StatementNode) *ProgramNode {
	return &ProgramNode{
		NodeBase: NodeBase{span: span},
		Body:     body,
	}
}

// Represents an empty statement eg. a statement with only a semicolon or a newline.
type EmptyStatementNode struct {
	NodeBase
}

func (*EmptyStatementNode) IsStatic() bool {
	return false
}

// Create a new empty statement node.
func NewEmptyStatementNode(span *position.Span) *EmptyStatementNode {
	return &EmptyStatementNode{
		NodeBase: NodeBase{span: span},
	}
}

// Expression optionally terminated with a newline or a semicolon.
type ImportStatementNode struct {
	NodeBase
	Path    StringLiteralNode
	FsPaths []string // resolved file system paths
}

func (i *ImportStatementNode) IsStatic() bool {
	return false
}

// Create a new import statement node eg. `import "foo"`
func NewImportStatementNode(span *position.Span, path StringLiteralNode) *ImportStatementNode {
	return &ImportStatementNode{
		NodeBase: NodeBase{span: span},
		Path:     path,
	}
}

// Expression optionally terminated with a newline or a semicolon.
type ExpressionStatementNode struct {
	NodeBase
	Expression ExpressionNode
}

func (e *ExpressionStatementNode) IsStatic() bool {
	return e.Expression.IsStatic()
}

// Create a new expression statement node eg. `5 * 2\n`
func NewExpressionStatementNode(span *position.Span, expr ExpressionNode) *ExpressionStatementNode {
	return &ExpressionStatementNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Same as [NewExpressionStatementNode] but returns an interface
func NewExpressionStatementNodeI(span *position.Span, expr ExpressionNode) StatementNode {
	return &ExpressionStatementNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Formal parameter optionally terminated with a newline or a semicolon.
type ParameterStatementNode struct {
	NodeBase
	Parameter ParameterNode
}

func (*ParameterStatementNode) IsStatic() bool {
	return false
}

// Create a new formal parameter statement node eg. `foo: Bar\n`
func NewParameterStatementNode(span *position.Span, param ParameterNode) *ParameterStatementNode {
	return &ParameterStatementNode{
		NodeBase:  NodeBase{span: span},
		Parameter: param,
	}
}

// Same as [NewParameterStatementNode] but returns an interface
func NewParameterStatementNodeI(span *position.Span, param ParameterNode) StructBodyStatementNode {
	return &ParameterStatementNode{
		NodeBase:  NodeBase{span: span},
		Parameter: param,
	}
}

// Represents a variable declaration with patterns eg. `var [foo, { bar }] = baz()`
type VariablePatternDeclarationNode struct {
	NodeBase
	Pattern     PatternNode
	Initialiser ExpressionNode // value assigned to the variable
}

func (*VariablePatternDeclarationNode) IsStatic() bool {
	return false
}

// Create a new variable declaration node with patterns eg. `var [foo, { bar }] = baz()`
func NewVariablePatternDeclarationNode(span *position.Span, pattern PatternNode, init ExpressionNode) *VariablePatternDeclarationNode {
	return &VariablePatternDeclarationNode{
		NodeBase:    NodeBase{span: span},
		Pattern:     pattern,
		Initialiser: init,
	}
}

// Represents an instance variable declaration eg. `var @foo: String`
type InstanceVariableDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name     string   // name of the variable
	TypeNode TypeNode // type of the variable
}

func (*InstanceVariableDeclarationNode) IsStatic() bool {
	return false
}

// Create a new instance variable declaration node eg. `var @foo: String`
func NewInstanceVariableDeclarationNode(span *position.Span, docComment string, name string, typ TypeNode) *InstanceVariableDeclarationNode {
	return &InstanceVariableDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Name:     name,
		TypeNode: typ,
	}
}

// Represents a variable declaration eg. `var foo: String`
type VariableDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name        string         // name of the variable
	TypeNode    TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
}

func (*VariableDeclarationNode) IsStatic() bool {
	return false
}

// Create a new variable declaration node eg. `var foo: String`
func NewVariableDeclarationNode(span *position.Span, docComment string, name string, typ TypeNode, init ExpressionNode) *VariableDeclarationNode {
	return &VariableDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Name:        name,
		TypeNode:    typ,
		Initialiser: init,
	}
}

// Represents a value pattern declaration eg. `val [foo, { bar }] = baz()`
type ValuePatternDeclarationNode struct {
	NodeBase
	Pattern     PatternNode
	Initialiser ExpressionNode // value assigned to the value
}

func (*ValuePatternDeclarationNode) IsStatic() bool {
	return false
}

// Create a new value declaration node eg. `val foo: String`
func NewValuePatternDeclarationNode(span *position.Span, pattern PatternNode, init ExpressionNode) *ValuePatternDeclarationNode {
	return &ValuePatternDeclarationNode{
		NodeBase:    NodeBase{span: span},
		Pattern:     pattern,
		Initialiser: init,
	}
}

// Represents a value declaration eg. `val foo: String`
type ValueDeclarationNode struct {
	TypedNodeBase
	Name        string         // name of the value
	TypeNode    TypeNode       // type of the value
	Initialiser ExpressionNode // value assigned to the value
}

func (*ValueDeclarationNode) IsStatic() bool {
	return false
}

// Create a new value declaration node eg. `val foo: String`
func NewValueDeclarationNode(span *position.Span, name string, typ TypeNode, init ExpressionNode) *ValueDeclarationNode {
	return &ValueDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Name:          name,
		TypeNode:      typ,
		Initialiser:   init,
	}
}

// Assignment with the specified operator.
type AssignmentExpressionNode struct {
	TypedNodeBase
	Op    *token.Token   // operator
	Left  ExpressionNode // left hand side
	Right ExpressionNode // right hand side
}

func (*AssignmentExpressionNode) IsStatic() bool {
	return false
}

// Create a new assignment expression node eg. `foo = 3`
func NewAssignmentExpressionNode(span *position.Span, op *token.Token, left, right ExpressionNode) *AssignmentExpressionNode {
	return &AssignmentExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Expression of an operator with two operands eg. `2 + 5`, `foo > bar`
type BinaryExpressionNode struct {
	TypedNodeBase
	Op     *token.Token   // operator
	Left   ExpressionNode // left hand side
	Right  ExpressionNode // right hand side
	static bool
}

func (b *BinaryExpressionNode) IsStatic() bool {
	return b.static
}

// Create a new binary expression node.
func NewBinaryExpressionNode(span *position.Span, op *token.Token, left, right ExpressionNode) *BinaryExpressionNode {
	return &BinaryExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
		static:        areExpressionsStatic(left, right),
	}
}

// Same as [NewBinaryExpressionNode] but returns an interface
func NewBinaryExpressionNodeI(span *position.Span, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return NewBinaryExpressionNode(span, op, left, right)
}

// Expression of a logical operator with two operands eg. `foo && bar`
type LogicalExpressionNode struct {
	TypedNodeBase
	Op     *token.Token   // operator
	Left   ExpressionNode // left hand side
	Right  ExpressionNode // right hand side
	static bool
}

func (l *LogicalExpressionNode) IsStatic() bool {
	return l.static
}

// Create a new logical expression node.
func NewLogicalExpressionNode(span *position.Span, op *token.Token, left, right ExpressionNode) *LogicalExpressionNode {
	return &LogicalExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
		static:        areExpressionsStatic(left, right),
	}
}

// Same as [NewLogicalExpressionNode] but returns an interface
func NewLogicalExpressionNodeI(span *position.Span, op *token.Token, left, right ExpressionNode) ExpressionNode {
	return &LogicalExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Expression of an operator with one operand eg. `!foo`, `-bar`
type UnaryExpressionNode struct {
	TypedNodeBase
	Op    *token.Token   // operator
	Right ExpressionNode // right hand side
}

func (u *UnaryExpressionNode) IsStatic() bool {
	return u.Right.IsStatic()
}

// Create a new unary expression node.
func NewUnaryExpressionNode(span *position.Span, op *token.Token, right ExpressionNode) *UnaryExpressionNode {
	return &UnaryExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Right:         right,
	}
}

// Postfix expression eg. `foo++`, `bar--`
type PostfixExpressionNode struct {
	NodeBase
	Op         *token.Token // operator
	Expression ExpressionNode
}

func (i *PostfixExpressionNode) IsStatic() bool {
	return false
}

// Create a new postfix expression node.
func NewPostfixExpressionNode(span *position.Span, op *token.Token, expr ExpressionNode) *PostfixExpressionNode {
	return &PostfixExpressionNode{
		NodeBase:   NodeBase{span: span},
		Op:         op,
		Expression: expr,
	}
}

// `bool` literal.
type BoolLiteralNode struct {
	NodeBase
}

func (*BoolLiteralNode) IsStatic() bool {
	return true
}

func (*BoolLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Bool{}
}

// Create a new `bool` literal node.
func NewBoolLiteralNode(span *position.Span) *BoolLiteralNode {
	return &BoolLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// `true` literal.
type TrueLiteralNode struct {
	NodeBase
}

func (*TrueLiteralNode) IsStatic() bool {
	return true
}

func (*TrueLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.True{}
}

// Create a new `true` literal node.
func NewTrueLiteralNode(span *position.Span) *TrueLiteralNode {
	return &TrueLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// `self` literal.
type FalseLiteralNode struct {
	NodeBase
}

func (*FalseLiteralNode) IsStatic() bool {
	return true
}

func (*FalseLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.False{}
}

// Create a new `false` literal node.
func NewFalseLiteralNode(span *position.Span) *FalseLiteralNode {
	return &FalseLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// `self` literal.
type SelfLiteralNode struct {
	TypedNodeBase
}

func (*SelfLiteralNode) IsStatic() bool {
	return false
}

// Create a new `self` literal node.
func NewSelfLiteralNode(span *position.Span) *SelfLiteralNode {
	return &SelfLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
	}
}

// `nil` literal.
type NilLiteralNode struct {
	NodeBase
}

func (*NilLiteralNode) SetType(types.Type) {}

func (*NilLiteralNode) IsStatic() bool {
	return true
}

func (*NilLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Nil{}
}

// Create a new `nil` literal node.
func NewNilLiteralNode(span *position.Span) *NilLiteralNode {
	return &NilLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// `undefined` literal.
type UndefinedLiteralNode struct {
	NodeBase
}

func (*UndefinedLiteralNode) IsStatic() bool {
	return true
}

// Create a new `undefined` literal node.
func NewUndefinedLiteralNode(span *position.Span) *UndefinedLiteralNode {
	return &UndefinedLiteralNode{
		NodeBase: NodeBase{span: span},
	}
}

// `void` type.
type VoidTypeNode struct {
	NodeBase
}

func (*VoidTypeNode) IsStatic() bool {
	return true
}

// Create a new `void` type node.
func NewVoidTypeNode(span *position.Span) *VoidTypeNode {
	return &VoidTypeNode{
		NodeBase: NodeBase{span: span},
	}
}

// `never` type.
type NeverTypeNode struct {
	NodeBase
}

func (*NeverTypeNode) IsStatic() bool {
	return true
}

func (*NeverTypeNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Never{}
}

// Create a new `never` type node.
func NewNeverTypeNode(span *position.Span) *NeverTypeNode {
	return &NeverTypeNode{
		NodeBase: NodeBase{span: span},
	}
}

// `any` type.
type AnyTypeNode struct {
	NodeBase
}

func (*AnyTypeNode) IsStatic() bool {
	return true
}

func (*AnyTypeNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Any{}
}

// Create a new `any` type node.
func NewAnyTypeNode(span *position.Span) *AnyTypeNode {
	return &AnyTypeNode{
		NodeBase: NodeBase{span: span},
	}
}

// Raw string literal enclosed with single quotes eg. `'foo'`.
type RawStringLiteralNode struct {
	TypedNodeBase
	Value string // value of the string literal
}

func (*RawStringLiteralNode) IsStatic() bool {
	return true
}

// Create a new raw string literal node eg. `'foo'`.
func NewRawStringLiteralNode(span *position.Span, val string) *RawStringLiteralNode {
	return &RawStringLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Char literal eg. `c"a"`
type CharLiteralNode struct {
	TypedNodeBase
	Value rune // value of the string literal
}

func (*CharLiteralNode) IsStatic() bool {
	return true
}

// Create a new char literal node eg. `c"a"`
func NewCharLiteralNode(span *position.Span, val rune) *CharLiteralNode {
	return &CharLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Raw Char literal eg. `a`
type RawCharLiteralNode struct {
	TypedNodeBase
	Value rune // value of the char literal
}

func (*RawCharLiteralNode) IsStatic() bool {
	return true
}

// Create a new raw char literal node eg. r`a`
func NewRawCharLiteralNode(span *position.Span, val rune) *RawCharLiteralNode {
	return &RawCharLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Int literal eg. `5`, `125_355`, `0xff`
type IntLiteralNode struct {
	TypedNodeBase
	Value string
}

func (*IntLiteralNode) IsStatic() bool {
	return true
}

// Create a new int literal node eg. `5`, `125_355`, `0xff`
func NewIntLiteralNode(span *position.Span, val string) *IntLiteralNode {
	return &IntLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Int64 literal eg. `5i64`, `125_355i64`, `0xffi64`
type Int64LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*Int64LiteralNode) IsStatic() bool {
	return true
}

// Create a new Int64 literal node eg. `5i64`, `125_355i64`, `0xffi64`
func NewInt64LiteralNode(span *position.Span, val string) *Int64LiteralNode {
	return &Int64LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// UInt64 literal eg. `5u64`, `125_355u64`, `0xffu64`
type UInt64LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*UInt64LiteralNode) IsStatic() bool {
	return true
}

// Create a new UInt64 literal node eg. `5u64`, `125_355u64`, `0xffu64`
func NewUInt64LiteralNode(span *position.Span, val string) *UInt64LiteralNode {
	return &UInt64LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Int32 literal eg. `5i32`, `1_20i32`, `0xffi32`
type Int32LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*Int32LiteralNode) IsStatic() bool {
	return true
}

// Create a new Int32 literal node eg. `5i32`, `1_20i32`, `0xffi32`
func NewInt32LiteralNode(span *position.Span, val string) *Int32LiteralNode {
	return &Int32LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// UInt32 literal eg. `5u32`, `1_20u32`, `0xffu32`
type UInt32LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*UInt32LiteralNode) IsStatic() bool {
	return true
}

// Create a new UInt32 literal node eg. `5u32`, `1_20u32`, `0xffu32`
func NewUInt32LiteralNode(span *position.Span, val string) *UInt32LiteralNode {
	return &UInt32LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Int16 literal eg. `5i16`, `1_20i16`, `0xffi16`
type Int16LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*Int16LiteralNode) IsStatic() bool {
	return true
}

// Create a new Int16 literal node eg. `5i16`, `1_20i16`, `0xffi16`
func NewInt16LiteralNode(span *position.Span, val string) *Int16LiteralNode {
	return &Int16LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// UInt16 literal eg. `5u16`, `1_20u16`, `0xffu16`
type UInt16LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*UInt16LiteralNode) IsStatic() bool {
	return true
}

// Create a new UInt16 literal node eg. `5u16`, `1_20u16`, `0xffu16`
func NewUInt16LiteralNode(span *position.Span, val string) *UInt16LiteralNode {
	return &UInt16LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Int8 literal eg. `5i8`, `1_20i8`, `0xffi8`
type Int8LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*Int8LiteralNode) IsStatic() bool {
	return true
}

// Create a new Int8 literal node eg. `5i8`, `1_20i8`, `0xffi8`
func NewInt8LiteralNode(span *position.Span, val string) *Int8LiteralNode {
	return &Int8LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// UInt8 literal eg. `5u8`, `1_20u8`, `0xffu8`
type UInt8LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*UInt8LiteralNode) IsStatic() bool {
	return true
}

// Create a new UInt8 literal node eg. `5u8`, `1_20u8`, `0xffu8`
func NewUInt8LiteralNode(span *position.Span, val string) *UInt8LiteralNode {
	return &UInt8LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Float literal eg. `5.2`, `.5`, `45e20`
type FloatLiteralNode struct {
	TypedNodeBase
	Value string
}

func (*FloatLiteralNode) IsStatic() bool {
	return true
}

// Create a new float literal node eg. `5.2`, `.5`, `45e20`
func NewFloatLiteralNode(span *position.Span, val string) *FloatLiteralNode {
	return &FloatLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// BigFloat literal eg. `5.2bf`, `.5bf`, `45e20bf`
type BigFloatLiteralNode struct {
	TypedNodeBase
	Value string
}

func (*BigFloatLiteralNode) IsStatic() bool {
	return true
}

// Create a new BigFloat literal node eg. `5.2bf`, `.5bf`, `45e20bf`
func NewBigFloatLiteralNode(span *position.Span, val string) *BigFloatLiteralNode {
	return &BigFloatLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Float64 literal eg. `5.2f64`, `.5f64`, `45e20f64`
type Float64LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*Float64LiteralNode) IsStatic() bool {
	return true
}

// Create a new Float64 literal node eg. `5.2f64`, `.5f64`, `45e20f64`
func NewFloat64LiteralNode(span *position.Span, val string) *Float64LiteralNode {
	return &Float64LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Float32 literal eg. `5.2f32`, `.5f32`, `45e20f32`
type Float32LiteralNode struct {
	TypedNodeBase
	Value string
}

func (*Float32LiteralNode) IsStatic() bool {
	return true
}

// Create a new Float32 literal node eg. `5.2f32`, `.5f32`, `45e20f32`
func NewFloat32LiteralNode(span *position.Span, val string) *Float32LiteralNode {
	return &Float32LiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a syntax error.
type InvalidNode struct {
	NodeBase
	Token *token.Token
}

func (*InvalidNode) SetType(types.Type) {}

func (*InvalidNode) IsStatic() bool {
	return false
}

func (*InvalidNode) IsOptional() bool {
	return false
}

// Create a new invalid node.
func NewInvalidNode(span *position.Span, tok *token.Token) *InvalidNode {
	return &InvalidNode{
		NodeBase: NodeBase{span: span},
		Token:    tok,
	}
}

func NewInvalidExpressionNode(span *position.Span, tok *token.Token) ExpressionNode {
	return NewInvalidNode(span, tok)
}

func NewInvalidPatternNode(span *position.Span, tok *token.Token) PatternNode {
	return NewInvalidNode(span, tok)
}

func NewInvalidPatternExpressionNode(span *position.Span, tok *token.Token) PatternExpressionNode {
	return NewInvalidNode(span, tok)
}

// Represents a single section of characters of a string literal eg. `foo` in `"foo${bar}"`.
type StringLiteralContentSectionNode struct {
	NodeBase
	Value string
}

func (*StringLiteralContentSectionNode) IsStatic() bool {
	return true
}

// Create a new string literal content section node eg. `foo` in `"foo${bar}"`.
func NewStringLiteralContentSectionNode(span *position.Span, val string) *StringLiteralContentSectionNode {
	return &StringLiteralContentSectionNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a single inspect interpolated section of a string literal eg. `bar + 2` in `"foo#{bar + 2}"`
type StringInspectInterpolationNode struct {
	NodeBase
	Expression ExpressionNode
}

func (*StringInspectInterpolationNode) IsStatic() bool {
	return false
}

// Create a new string inspect interpolation node eg. `bar + 2` in `"foo#{bar + 2}"`
func NewStringInspectInterpolationNode(span *position.Span, expr ExpressionNode) *StringInspectInterpolationNode {
	return &StringInspectInterpolationNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Represents a single interpolated section of a string literal eg. `bar + 2` in `"foo${bar + 2}"`
type StringInterpolationNode struct {
	NodeBase
	Expression ExpressionNode
}

func (*StringInterpolationNode) IsStatic() bool {
	return false
}

// Create a new string interpolation node eg. `bar + 2` in `"foo${bar + 2}"`
func NewStringInterpolationNode(span *position.Span, expr ExpressionNode) *StringInterpolationNode {
	return &StringInterpolationNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Represents an interpolated string literal eg. `"foo ${bar} baz"`
type InterpolatedStringLiteralNode struct {
	NodeBase
	Content []StringLiteralContentNode
}

func (*InterpolatedStringLiteralNode) IsStatic() bool {
	return false
}

func (*InterpolatedStringLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return globalEnv.StdSubtype(symbol.String)
}

// Create a new interpolated string literal node eg. `"foo ${bar} baz"`
func NewInterpolatedStringLiteralNode(span *position.Span, cont []StringLiteralContentNode) *InterpolatedStringLiteralNode {
	return &InterpolatedStringLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  cont,
	}
}

// Represents a simple double quoted string literal eg. `"foo baz"`
type DoubleQuotedStringLiteralNode struct {
	TypedNodeBase
	Value string
}

func (*DoubleQuotedStringLiteralNode) IsStatic() bool {
	return true
}

// Create a new double quoted string literal node eg. `"foo baz"`
func NewDoubleQuotedStringLiteralNode(span *position.Span, val string) *DoubleQuotedStringLiteralNode {
	return &DoubleQuotedStringLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a public identifier eg. `foo`.
type PublicIdentifierNode struct {
	TypedNodeBase
	Value string
}

func (*PublicIdentifierNode) IsStatic() bool {
	return false
}

// Create a new public identifier node eg. `foo`.
func NewPublicIdentifierNode(span *position.Span, val string) *PublicIdentifierNode {
	return &PublicIdentifierNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a type expression `type String?`
type TypeExpressionNode struct {
	NodeBase
	TypeNode TypeNode
}

func (*TypeExpressionNode) IsStatic() bool {
	return false
}

// Create a new type expression `type String?`
func NewTypeExpressionNode(span *position.Span, typeNode TypeNode) *TypeExpressionNode {
	return &TypeExpressionNode{
		NodeBase: NodeBase{span: span},
		TypeNode: typeNode,
	}
}

// Represents an uninterpolated regex literal eg. `%/foo/`
type UninterpolatedRegexLiteralNode struct {
	NodeBase
	Content string
	Flags   bitfield.BitField8
}

func (*UninterpolatedRegexLiteralNode) Type(env *types.GlobalEnvironment) types.Type {
	return env.StdSubtype(symbol.Regex)
}

func (*UninterpolatedRegexLiteralNode) IsStatic() bool {
	return true
}

func (r *UninterpolatedRegexLiteralNode) IsCaseInsensitive() bool {
	return r.Flags.HasFlag(flag.CaseInsensitiveFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetCaseInsensitive() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.CaseInsensitiveFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsMultiline() bool {
	return r.Flags.HasFlag(flag.MultilineFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetMultiline() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.MultilineFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsDotAll() bool {
	return r.Flags.HasFlag(flag.DotAllFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetDotAll() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.DotAllFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsUngreedy() bool {
	return r.Flags.HasFlag(flag.UngreedyFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetUngreedy() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.UngreedyFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsASCII() bool {
	return r.Flags.HasFlag(flag.ASCIIFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetASCII() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ASCIIFlag)
	return r
}

func (r *UninterpolatedRegexLiteralNode) IsExtended() bool {
	return r.Flags.HasFlag(flag.ExtendedFlag)
}

func (r *UninterpolatedRegexLiteralNode) SetExtended() *UninterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ExtendedFlag)
	return r
}

// Create a new uninterpolated regex literal node eg. `%/foo/`.
func NewUninterpolatedRegexLiteralNode(span *position.Span, content string, flags bitfield.BitField8) *UninterpolatedRegexLiteralNode {
	return &UninterpolatedRegexLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  content,
		Flags:    flags,
	}
}

// Represents a single section of characters of a regex literal eg. `foo` in `%/foo${bar}/`.
type RegexLiteralContentSectionNode struct {
	NodeBase
	Value string
}

func (*RegexLiteralContentSectionNode) IsStatic() bool {
	return true
}

// Create a new regex literal content section node eg. `foo` in `%/foo${bar}/`.
func NewRegexLiteralContentSectionNode(span *position.Span, val string) *RegexLiteralContentSectionNode {
	return &RegexLiteralContentSectionNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a single interpolated section of a regex literal eg. `bar + 2` in `%/foo${bar + 2}/`
type RegexInterpolationNode struct {
	NodeBase
	Expression ExpressionNode
}

func (*RegexInterpolationNode) IsStatic() bool {
	return false
}

// Create a new regex interpolation node eg. `bar + 2` in `%/foo${bar + 2}/`
func NewRegexInterpolationNode(span *position.Span, expr ExpressionNode) *RegexInterpolationNode {
	return &RegexInterpolationNode{
		NodeBase:   NodeBase{span: span},
		Expression: expr,
	}
}

// Represents an Interpolated regex literal eg. `%/foo${1 + 2}bar/`
type InterpolatedRegexLiteralNode struct {
	NodeBase
	Content []RegexLiteralContentNode
	Flags   bitfield.BitField8
}

func (*InterpolatedRegexLiteralNode) Type(env *types.GlobalEnvironment) types.Type {
	return env.StdSubtype(symbol.Regex)
}

func (*InterpolatedRegexLiteralNode) IsStatic() bool {
	return false
}

func (r *InterpolatedRegexLiteralNode) IsCaseInsensitive() bool {
	return r.Flags.HasFlag(flag.CaseInsensitiveFlag)
}

func (r *InterpolatedRegexLiteralNode) SetCaseInsensitive() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.CaseInsensitiveFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsMultiline() bool {
	return r.Flags.HasFlag(flag.MultilineFlag)
}

func (r *InterpolatedRegexLiteralNode) SetMultiline() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.MultilineFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsDotAll() bool {
	return r.Flags.HasFlag(flag.DotAllFlag)
}

func (r *InterpolatedRegexLiteralNode) SetDotAll() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.DotAllFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsUngreedy() bool {
	return r.Flags.HasFlag(flag.UngreedyFlag)
}

func (r *InterpolatedRegexLiteralNode) SetUngreedy() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.UngreedyFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsASCII() bool {
	return r.Flags.HasFlag(flag.ASCIIFlag)
}

func (r *InterpolatedRegexLiteralNode) SetASCII() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ASCIIFlag)
	return r
}

func (r *InterpolatedRegexLiteralNode) IsExtended() bool {
	return r.Flags.HasFlag(flag.ExtendedFlag)
}

func (r *InterpolatedRegexLiteralNode) SetExtended() *InterpolatedRegexLiteralNode {
	r.Flags.SetFlag(flag.ExtendedFlag)
	return r
}

// Create a new interpolated regex literal node eg. `%/foo${1 + 3}bar/`.
func NewInterpolatedRegexLiteralNode(span *position.Span, content []RegexLiteralContentNode, flags bitfield.BitField8) *InterpolatedRegexLiteralNode {
	return &InterpolatedRegexLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  content,
		Flags:    flags,
	}
}

// Represents an instance variable eg. `@foo`
type InstanceVariableNode struct {
	TypedNodeBase
	Value string
}

func (*InstanceVariableNode) IsStatic() bool {
	return false
}

// Create an instance variable node eg. `@foo`.
func NewInstanceVariableNode(span *position.Span, val string) *InstanceVariableNode {
	return &InstanceVariableNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a private identifier eg. `_foo`
type PrivateIdentifierNode struct {
	TypedNodeBase
	Value string
}

func (*PrivateIdentifierNode) IsStatic() bool {
	return false
}

// Create a new private identifier node eg. `_foo`.
func NewPrivateIdentifierNode(span *position.Span, val string) *PrivateIdentifierNode {
	return &PrivateIdentifierNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a public constant eg. `Foo`.
type PublicConstantNode struct {
	TypedNodeBase
	Value string
}

func (*PublicConstantNode) IsStatic() bool {
	return false
}

// Create a new public constant node eg. `Foo`.
func NewPublicConstantNode(span *position.Span, val string) *PublicConstantNode {
	return &PublicConstantNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a private constant eg. `_Foo`
type PrivateConstantNode struct {
	TypedNodeBase
	Value string
}

func (*PrivateConstantNode) IsStatic() bool {
	return false
}

// Create a new private constant node eg. `_Foo`.
func NewPrivateConstantNode(span *position.Span, val string) *PrivateConstantNode {
	return &PrivateConstantNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Value:         val,
	}
}

// Represents a `catch` eg.
//
//	catch SomeError(message)
//		print("awesome!")
//	end
type CatchNode struct {
	NodeBase
	Pattern PatternNode
	Body    []StatementNode // do expression body
}

func (*CatchNode) IsStatic() bool {
	return false
}

// Create a new `catch` node eg.
//
//	catch SomeError(message)
//		print("awesome!")
//	end
func NewCatchNode(span *position.Span, pattern PatternNode, body []StatementNode) *CatchNode {
	return &CatchNode{
		NodeBase: NodeBase{span: span},
		Pattern:  pattern,
		Body:     body,
	}
}

// Represents a `do` expression eg.
//
//	do
//		print("awesome!")
//	end
type DoExpressionNode struct {
	TypedNodeBase
	Body    []StatementNode // do expression body
	Catches []*CatchNode
	Finally []StatementNode
}

func (*DoExpressionNode) IsStatic() bool {
	return false
}

// Create a new `do` expression node eg.
//
//	do
//		print("awesome!")
//	end
func NewDoExpressionNode(span *position.Span, body []StatementNode, catches []*CatchNode, finally []StatementNode) *DoExpressionNode {
	return &DoExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Body:          body,
		Catches:       catches,
		Finally:       finally,
	}
}

// Represents a `singleton` block expression eg.
//
//	singleton
//		def hello then println("awesome!")
//	end
type SingletonBlockExpressionNode struct {
	TypedNodeBase
	Body []StatementNode // do expression body
}

func (*SingletonBlockExpressionNode) SkipTypechecking() bool {
	return false
}

func (*SingletonBlockExpressionNode) IsStatic() bool {
	return false
}

// Create a new `singleton` block expression node eg.
//
//	singleton
//		def hello then println("awesome!")
//	end
func NewSingletonBlockExpressionNode(span *position.Span, body []StatementNode) *SingletonBlockExpressionNode {
	return &SingletonBlockExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Body:          body,
	}
}

// Represents an `if`, `unless`, `while` or `until` modifier expression eg. `return true if foo`.
type ModifierNode struct {
	TypedNodeBase
	Modifier *token.Token   // modifier token
	Left     ExpressionNode // left hand side
	Right    ExpressionNode // right hand side
}

func (*ModifierNode) IsStatic() bool {
	return false
}

// Create a new modifier node eg. `return true if foo`.
func NewModifierNode(span *position.Span, mod *token.Token, left, right ExpressionNode) *ModifierNode {
	return &ModifierNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Modifier:      mod,
		Left:          left,
		Right:         right,
	}
}

// Represents an `if .. else` modifier expression eg. `foo = 1 if bar else foo = 2`
type ModifierIfElseNode struct {
	TypedNodeBase
	ThenExpression ExpressionNode // then expression body
	Condition      ExpressionNode // if condition
	ElseExpression ExpressionNode // else expression body
}

func (*ModifierIfElseNode) IsStatic() bool {
	return false
}

// Create a new modifier `if` .. `else` node eg. `foo = 1 if bar else foo = 2“.
func NewModifierIfElseNode(span *position.Span, then, cond, els ExpressionNode) *ModifierIfElseNode {
	return &ModifierIfElseNode{
		TypedNodeBase:  TypedNodeBase{span: span},
		ThenExpression: then,
		Condition:      cond,
		ElseExpression: els,
	}
}

// Represents an `for .. in` modifier expression eg. `println(i) for i in 10..30`
type ModifierForInNode struct {
	NodeBase
	ThenExpression ExpressionNode // then expression body
	Parameter      PatternNode
	InExpression   ExpressionNode // expression that will be iterated through
}

func (*ModifierForInNode) IsStatic() bool {
	return false
}

// Create a new modifier `for` .. `in` node eg. `println(i) for i in 10..30`
func NewModifierForInNode(span *position.Span, then ExpressionNode, param PatternNode, in ExpressionNode) *ModifierForInNode {
	return &ModifierForInNode{
		NodeBase:       NodeBase{span: span},
		ThenExpression: then,
		Parameter:      param,
		InExpression:   in,
	}
}

// Pattern with two operands eg. `> 10 && < 50`
type BinaryPatternNode struct {
	NodeBase
	Op    *token.Token // operator
	Left  PatternNode  // left hand side
	Right PatternNode  // right hand side
}

func (*BinaryPatternNode) IsStatic() bool {
	return false
}

// Create a new binary pattern node eg. `> 10 && < 50`
func NewBinaryPatternNode(span *position.Span, op *token.Token, left, right PatternNode) *BinaryPatternNode {
	return &BinaryPatternNode{
		NodeBase: NodeBase{span: span},
		Op:       op,
		Left:     left,
		Right:    right,
	}
}

// Same as [NewBinaryPatternNode] but returns an interface
func NewBinaryPatternNodeI(span *position.Span, op *token.Token, left, right PatternNode) PatternNode {
	return NewBinaryPatternNode(span, op, left, right)
}

// Represents an as pattern eg. `> 5 && < 20 as foo`
type AsPatternNode struct {
	NodeBase
	Pattern PatternNode
	Name    IdentifierNode
}

func (*AsPatternNode) IsStatic() bool {
	return false
}

// Create an Object pattern node eg. `Foo(foo: 5, bar: a, c)`
func NewAsPatternNode(span *position.Span, pattern PatternNode, name IdentifierNode) *AsPatternNode {
	return &AsPatternNode{
		NodeBase: NodeBase{span: span},
		Pattern:  pattern,
		Name:     name,
	}
}

// Represents a symbol value pattern eg. `foo: bar`
type SymbolKeyValuePatternNode struct {
	NodeBase
	Key   string
	Value PatternNode
}

func (s *SymbolKeyValuePatternNode) IsStatic() bool {
	return false
}

// Create a symbol key value node eg. `foo: bar`
func NewSymbolKeyValuePatternNode(span *position.Span, key string, val PatternNode) *SymbolKeyValuePatternNode {
	return &SymbolKeyValuePatternNode{
		NodeBase: NodeBase{span: span},
		Key:      key,
		Value:    val,
	}
}

// Represents a key value pattern eg. `foo => bar`
type KeyValuePatternNode struct {
	NodeBase
	Key   PatternNode
	Value PatternNode
}

func (k *KeyValuePatternNode) IsStatic() bool {
	return false
}

// Create a key value pattern node eg. `foo => bar`
func NewKeyValuePatternNode(span *position.Span, key, val PatternNode) *KeyValuePatternNode {
	return &KeyValuePatternNode{
		NodeBase: NodeBase{span: span},
		Key:      key,
		Value:    val,
	}
}

// Represents an Object pattern eg. `Foo(foo: 5, bar: a, c)`
type ObjectPatternNode struct {
	NodeBase
	Class      ComplexConstantNode
	Attributes []PatternNode
}

func (m *ObjectPatternNode) IsStatic() bool {
	return false
}

// Create an Object pattern node eg. `Foo(foo: 5, bar: a, c)`
func NewObjectPatternNode(span *position.Span, class ComplexConstantNode, attrs []PatternNode) *ObjectPatternNode {
	return &ObjectPatternNode{
		NodeBase:   NodeBase{span: span},
		Class:      class,
		Attributes: attrs,
	}
}

// Represents a Record pattern eg. `%{ foo: 5, bar: a, 5 => >= 10 }`
type RecordPatternNode struct {
	NodeBase
	Elements []PatternNode
}

func (m *RecordPatternNode) IsStatic() bool {
	return false
}

// Create a Record pattern node eg. `%{ foo: 5, bar: a, 5 => >= 10 }`
func NewRecordPatternNode(span *position.Span, elements []PatternNode) *RecordPatternNode {
	return &RecordPatternNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Same as [NewRecordPatternNode] but returns an interface
func NewRecordPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewRecordPatternNode(span, elements)
}

// Represents a Map pattern eg. `{ foo: 5, bar: a, 5 => >= 10 }`
type MapPatternNode struct {
	NodeBase
	Elements []PatternNode
}

func (m *MapPatternNode) IsStatic() bool {
	return false
}

// Create a Map pattern node eg. `{ foo: 5, bar: a, 5 => >= 10 }`
func NewMapPatternNode(span *position.Span, elements []PatternNode) *MapPatternNode {
	return &MapPatternNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Same as [NewMapPatternNode] but returns an interface
func NewMapPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewMapPatternNode(span, elements)
}

// Represents a Set pattern eg. `^[1, "foo"]`
type SetPatternNode struct {
	NodeBase
	Elements []PatternNode
}

func (s *SetPatternNode) IsStatic() bool {
	return false
}

// Create a Set pattern node eg. `^[1, "foo"]`
func NewSetPatternNode(span *position.Span, elements []PatternNode) *SetPatternNode {
	return &SetPatternNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Same as [NewSetPatternNode] but returns an interface
func NewSetPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewSetPatternNode(span, elements)
}

// Represents a rest element in a list pattern eg. `*a`
type RestPatternNode struct {
	NodeBase
	Identifier IdentifierNode
}

func (r *RestPatternNode) IsStatic() bool {
	return false
}

// Create a rest pattern node eg. `*a`
func NewRestPatternNode(span *position.Span, ident IdentifierNode) *RestPatternNode {
	return &RestPatternNode{
		NodeBase:   NodeBase{span: span},
		Identifier: ident,
	}
}

// Represents a Tuple pattern eg. `%[1, a, >= 10]`
type TuplePatternNode struct {
	NodeBase
	Elements []PatternNode
}

func (l *TuplePatternNode) IsStatic() bool {
	return false
}

// Create a Tuple pattern node eg. `%[1, a, >= 10]`
func NewTuplePatternNode(span *position.Span, elements []PatternNode) *TuplePatternNode {
	return &TuplePatternNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Same as [NewTuplePatternNode] but returns an interface
func NewTuplePatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewTuplePatternNode(span, elements)
}

// Represents a List pattern eg. `[1, a, >= 10]`
type ListPatternNode struct {
	NodeBase
	Elements []PatternNode
}

func (l *ListPatternNode) IsStatic() bool {
	return false
}

// Create a List pattern node eg. `[1, a, >= 10]`
func NewListPatternNode(span *position.Span, elements []PatternNode) *ListPatternNode {
	return &ListPatternNode{
		NodeBase: NodeBase{span: span},
		Elements: elements,
	}
}

// Same as [NewListPatternNode] but returns an interface
func NewListPatternNodeI(span *position.Span, elements []PatternNode) PatternNode {
	return NewListPatternNode(span, elements)
}

// Represents a `case` node eg. `case 3 then println("eureka!")`
type CaseNode struct {
	NodeBase
	Pattern PatternNode
	Body    []StatementNode
}

func (*CaseNode) IsStatic() bool {
	return false
}

// Create a new `case` node
func NewCaseNode(span *position.Span, pattern PatternNode, body []StatementNode) *CaseNode {
	return &CaseNode{
		NodeBase: NodeBase{span: span},
		Pattern:  pattern,
		Body:     body,
	}
}

// Represents a `switch` expression eg.
//
//	switch a
//	case 3
//	  println("eureka!")
//	case nil
//	  println("boo")
//	else
//	  println("nothing")
//	end
type SwitchExpressionNode struct {
	NodeBase
	Value    ExpressionNode
	Cases    []*CaseNode
	ElseBody []StatementNode
}

func (*SwitchExpressionNode) IsStatic() bool {
	return false
}

// Create a new `switch` expression node
func NewSwitchExpressionNode(span *position.Span, val ExpressionNode, cases []*CaseNode, els []StatementNode) *SwitchExpressionNode {
	return &SwitchExpressionNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
		Cases:    cases,
		ElseBody: els,
	}
}

// Represents an `if` expression eg. `if foo then println("bar")`
type IfExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // if condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

func (*IfExpressionNode) IsStatic() bool {
	return false
}

// Create a new `if` expression node eg. `if foo then println("bar")`
func NewIfExpressionNode(span *position.Span, cond ExpressionNode, then, els []StatementNode) *IfExpressionNode {
	return &IfExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		ThenBody:      then,
		Condition:     cond,
		ElseBody:      els,
	}
}

// Represents an `unless` expression eg. `unless foo then println("bar")`
type UnlessExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // unless condition
	ThenBody  []StatementNode // then expression body
	ElseBody  []StatementNode // else expression body
}

func (*UnlessExpressionNode) IsStatic() bool {
	return false
}

// Create a new `unless` expression node eg. `unless foo then println("bar")`
func NewUnlessExpressionNode(span *position.Span, cond ExpressionNode, then, els []StatementNode) *UnlessExpressionNode {
	return &UnlessExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		ThenBody:      then,
		Condition:     cond,
		ElseBody:      els,
	}
}

// Represents a `while` expression eg. `while i < 5 then i += 5`
type WhileExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // while condition
	ThenBody  []StatementNode // then expression body
}

func (*WhileExpressionNode) IsStatic() bool {
	return false
}

// Create a new `while` expression node eg. `while i < 5 then i += 5`
func NewWhileExpressionNode(span *position.Span, cond ExpressionNode, then []StatementNode) *WhileExpressionNode {
	return &WhileExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Condition:     cond,
		ThenBody:      then,
	}
}

// Represents a `until` expression eg. `until i >= 5 then i += 5`
type UntilExpressionNode struct {
	TypedNodeBase
	Condition ExpressionNode  // until condition
	ThenBody  []StatementNode // then expression body
}

func (*UntilExpressionNode) IsStatic() bool {
	return false
}

// Create a new `until` expression node eg. `until i >= 5 then i += 5`
func NewUntilExpressionNode(span *position.Span, cond ExpressionNode, then []StatementNode) *UntilExpressionNode {
	return &UntilExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Condition:     cond,
		ThenBody:      then,
	}
}

// Represents a `loop` expression.
type LoopExpressionNode struct {
	TypedNodeBase
	ThenBody []StatementNode // then expression body
}

func (*LoopExpressionNode) IsStatic() bool {
	return false
}

// Create a new `loop` expression node eg. `loop println('elk is awesome')`
func NewLoopExpressionNode(span *position.Span, then []StatementNode) *LoopExpressionNode {
	return &LoopExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		ThenBody:      then,
	}
}

// Represents a numeric `for` expression eg. `fornum i := 0; i < 10; i += 1 then println(i)`
type NumericForExpressionNode struct {
	TypedNodeBase
	Initialiser ExpressionNode  // i := 0
	Condition   ExpressionNode  // i < 10
	Increment   ExpressionNode  // i += 1
	ThenBody    []StatementNode // then expression body
}

func (*NumericForExpressionNode) IsStatic() bool {
	return false
}

// Create a new numeric `fornum` expression eg. `for i := 0; i < 10; i += 1 then println(i)`
func NewNumericForExpressionNode(span *position.Span, init, cond, incr ExpressionNode, then []StatementNode) *NumericForExpressionNode {
	return &NumericForExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Initialiser:   init,
		Condition:     cond,
		Increment:     incr,
		ThenBody:      then,
	}
}

// Represents a `for in` expression eg. `for i in 5..15 then println(i)`
type ForInExpressionNode struct {
	NodeBase
	Parameter    PatternNode     // parameter
	InExpression ExpressionNode  // expression that will be iterated through
	ThenBody     []StatementNode // then expression body
}

func (*ForInExpressionNode) IsStatic() bool {
	return false
}

// Create a new `for in` expression node eg. `for i in 5..15 then println(i)`
func NewForInExpressionNode(span *position.Span, param PatternNode, inExpr ExpressionNode, then []StatementNode) *ForInExpressionNode {
	return &ForInExpressionNode{
		NodeBase:     NodeBase{span: span},
		Parameter:    param,
		InExpression: inExpr,
		ThenBody:     then,
	}
}

// Represents a labeled expression eg. `$foo: 1 + 2`
type LabeledExpressionNode struct {
	NodeBase
	Label      string
	Expression ExpressionNode
}

func (l *LabeledExpressionNode) Type(env *types.GlobalEnvironment) types.Type {
	return l.Expression.Type(env)
}

func (l *LabeledExpressionNode) IsStatic() bool {
	return l.Expression.IsStatic()
}

// Create a new labeled expression node eg. `$foo: 1 + 2`
func NewLabeledExpressionNode(span *position.Span, label string, expr ExpressionNode) *LabeledExpressionNode {
	return &LabeledExpressionNode{
		NodeBase:   NodeBase{span: span},
		Label:      label,
		Expression: expr,
	}
}

// Represents a `break` expression eg. `break`, `break false`
type BreakExpressionNode struct {
	NodeBase
	Label string
	Value ExpressionNode
}

func (*BreakExpressionNode) IsStatic() bool {
	return false
}

func (*BreakExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

// Create a new `break` expression node eg. `break`
func NewBreakExpressionNode(span *position.Span, label string, val ExpressionNode) *BreakExpressionNode {
	return &BreakExpressionNode{
		NodeBase: NodeBase{span: span},
		Label:    label,
		Value:    val,
	}
}

// Represents a `return` expression eg. `return`, `return true`
type ReturnExpressionNode struct {
	NodeBase
	Value ExpressionNode
}

func (*ReturnExpressionNode) IsStatic() bool {
	return false
}

func (*ReturnExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

// Create a new `return` expression node eg. `return`, `return true`
func NewReturnExpressionNode(span *position.Span, val ExpressionNode) *ReturnExpressionNode {
	return &ReturnExpressionNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a `continue` expression eg. `continue`, `continue "foo"`
type ContinueExpressionNode struct {
	NodeBase
	Label string
	Value ExpressionNode
}

func (*ContinueExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

func (*ContinueExpressionNode) IsStatic() bool {
	return false
}

// Create a new `continue` expression node eg. `continue`, `continue "foo"`
func NewContinueExpressionNode(span *position.Span, label string, val ExpressionNode) *ContinueExpressionNode {
	return &ContinueExpressionNode{
		NodeBase: NodeBase{span: span},
		Label:    label,
		Value:    val,
	}
}

// Represents a `throw` expression eg. `throw ArgumentError.new("foo")`
type ThrowExpressionNode struct {
	NodeBase
	Value ExpressionNode
}

func (*ThrowExpressionNode) Type(*types.GlobalEnvironment) types.Type {
	return types.Never{}
}

func (*ThrowExpressionNode) IsStatic() bool {
	return false
}

// Create a new `throw` expression node eg. `throw ArgumentError.new("foo")`
func NewThrowExpressionNode(span *position.Span, val ExpressionNode) *ThrowExpressionNode {
	return &ThrowExpressionNode{
		NodeBase: NodeBase{span: span},
		Value:    val,
	}
}

// Represents a constant declaration eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
type ConstantDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant    ExpressionNode // name of the constant
	TypeNode    TypeNode       // type of the constant
	Initialiser ExpressionNode // value assigned to the constant
}

func (*ConstantDeclarationNode) IsStatic() bool {
	return false
}

// Create a new constant declaration node eg. `const Foo: ArrayList[String] = ["foo", "bar"]`
func NewConstantDeclarationNode(span *position.Span, docComment string, constant ExpressionNode, typ TypeNode, init ExpressionNode) *ConstantDeclarationNode {
	return &ConstantDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:    constant,
		TypeNode:    typ,
		Initialiser: init,
	}
}

// Type expression of an operator with two operands eg. `String | Int`
type BinaryTypeExpressionNode struct {
	TypedNodeBase
	Op    *token.Token // operator
	Left  TypeNode     // left hand side
	Right TypeNode     // right hand side
}

func (*BinaryTypeExpressionNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String | Int`
func NewBinaryTypeExpressionNode(span *position.Span, op *token.Token, left, right TypeNode) *BinaryTypeExpressionNode {
	return &BinaryTypeExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Same as [NewBinaryTypeExpressionNode] but returns an interface
func NewBinaryTypeExpressionNodeI(span *position.Span, op *token.Token, left, right TypeNode) TypeNode {
	return &BinaryTypeExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Op:            op,
		Left:          left,
		Right:         right,
	}
}

// Union type eg. `String & Int & Float`
type IntersectionTypeNode struct {
	TypedNodeBase
	Elements []TypeNode
}

func (*IntersectionTypeNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String & Int`
func NewIntersectionTypeNode(span *position.Span, elements []TypeNode) *IntersectionTypeNode {
	return &IntersectionTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Union type eg. `String | Int | Float`
type UnionTypeNode struct {
	TypedNodeBase
	Elements []TypeNode
}

func (*UnionTypeNode) IsStatic() bool {
	return false
}

// Create a new binary type expression node eg. `String | Int`
func NewUnionTypeNode(span *position.Span, elements []TypeNode) *UnionTypeNode {
	return &UnionTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Represents an optional or nilable type eg. `String?`
type NilableTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (*NilableTypeNode) IsStatic() bool {
	return false
}

// Create a new nilable type node eg. `String?`
func NewNilableTypeNode(span *position.Span, typ TypeNode) *NilableTypeNode {
	return &NilableTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		TypeNode:      typ,
	}
}

// Represents an instance type eg. `^self`
type InstanceOfTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (*InstanceOfTypeNode) IsStatic() bool {
	return false
}

// Create a new instance of type node eg. `^self`
func NewInstanceOfTypeNode(span *position.Span, typ TypeNode) *InstanceOfTypeNode {
	return &InstanceOfTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		TypeNode:      typ,
	}
}

// Represents a singleton type eg. `&String`
type SingletonTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (*SingletonTypeNode) IsStatic() bool {
	return false
}

// Create a new singleton type node eg. `&String`
func NewSingletonTypeNode(span *position.Span, typ TypeNode) *SingletonTypeNode {
	return &SingletonTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		TypeNode:      typ,
	}
}

// Represents a not type eg. `~String`
type NotTypeNode struct {
	TypedNodeBase
	TypeNode TypeNode // right hand side
}

func (*NotTypeNode) IsStatic() bool {
	return false
}

// Create a new not type node eg. `~String`
func NewNotTypeNode(span *position.Span, typ TypeNode) *NotTypeNode {
	return &NotTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		TypeNode:      typ,
	}
}

// Represents a constant lookup expressions eg. `Foo::Bar`
type ConstantLookupNode struct {
	TypedNodeBase
	Left  ExpressionNode      // left hand side
	Right ComplexConstantNode // right hand side
}

func (*ConstantLookupNode) IsStatic() bool {
	return false
}

// Create a new constant lookup expression node eg. `Foo::Bar`
func NewConstantLookupNode(span *position.Span, left ExpressionNode, right ComplexConstantNode) *ConstantLookupNode {
	return &ConstantLookupNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Left:          left,
		Right:         right,
	}
}

// Represents a method lookup expression eg. `Foo::bar`, `a::c`
type MethodLookupNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Name     string
}

func (*MethodLookupNode) IsStatic() bool {
	return false
}

// Create a new method lookup expression node eg. `Foo::bar`, `a::c`
func NewMethodLookupNode(span *position.Span, receiver ExpressionNode, name string) *MethodLookupNode {
	return &MethodLookupNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Receiver:      receiver,
		Name:          name,
	}
}

// Represents a using all entry node eg. `Foo::*`, `A::B::C::*`
type UsingAllEntryNode struct {
	TypedNodeBase
	Namespace ExpressionNode
}

func (*UsingAllEntryNode) IsStatic() bool {
	return false
}

// Create a new using all entry node eg. `Foo::*`, `A::B::C::*`
func NewUsingAllEntryNode(span *position.Span, namespace UsingEntryNode) *UsingAllEntryNode {
	return &UsingAllEntryNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Namespace:     namespace,
	}
}

// Represents a using entry node with subentries eg. `Foo::{Bar, baz}`, `A::B::C::{lol, foo as epic, Gro as Moe}`
type UsingEntryWithSubentriesNode struct {
	NodeBase
	Namespace  ExpressionNode
	Subentries []UsingSubentryNode
}

func (*UsingEntryWithSubentriesNode) IsStatic() bool {
	return false
}

// Create a new using all entry node eg. `Foo::*`, `A::B::C::*`
func NewUsingEntryWithSubentriesNode(span *position.Span, namespace UsingEntryNode, subentries []UsingSubentryNode) *UsingEntryWithSubentriesNode {
	return &UsingEntryWithSubentriesNode{
		NodeBase:   NodeBase{span: span},
		Namespace:  namespace,
		Subentries: subentries,
	}
}

// Represents an identifier with as in using declarations
// eg. `foo as bar`.
type PublicIdentifierAsNode struct {
	NodeBase
	TargetName string
	AsName     string
}

func (*PublicIdentifierAsNode) IsStatic() bool {
	return false
}

// Create a new identifier with as eg. `foo as bar`.
func NewPublicIdentifierAsNode(span *position.Span, target, as string) *PublicIdentifierAsNode {
	return &PublicIdentifierAsNode{
		NodeBase:   NodeBase{span: span},
		TargetName: target,
		AsName:     as,
	}
}

// Represents a constant with as in using declarations
// eg. `Foo as Bar`.
type PublicConstantAsNode struct {
	NodeBase
	TargetName string
	AsName     string
}

func (*PublicConstantAsNode) IsStatic() bool {
	return false
}

// Create a new identifier with as eg. `Foo as Bar`.
func NewPublicConstantAsNode(span *position.Span, target, as string) *PublicConstantAsNode {
	return &PublicConstantAsNode{
		NodeBase:   NodeBase{span: span},
		TargetName: target,
		AsName:     as,
	}
}

// Indicates whether the parameter is a rest param
type ParameterKind uint8

const (
	NormalParameterKind ParameterKind = iota
	PositionalRestParameterKind
	NamedRestParameterKind
)

// Represents a formal parameter in function or struct declarations eg. `foo: String = 'bar'`
type FormalParameterNode struct {
	NodeBase
	Name        string         // name of the variable
	TypeNode    TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
	Kind        ParameterKind
}

func (*FormalParameterNode) IsStatic() bool {
	return false
}

func (f *FormalParameterNode) IsOptional() bool {
	return f.Initialiser != nil
}

// Create a new formal parameter node eg. `foo: String = 'bar'`
func NewFormalParameterNode(span *position.Span, name string, typ TypeNode, init ExpressionNode, kind ParameterKind) *FormalParameterNode {
	return &FormalParameterNode{
		NodeBase:    NodeBase{span: span},
		Name:        name,
		TypeNode:    typ,
		Initialiser: init,
		Kind:        kind,
	}
}

// Represents a formal parameter in method declarations eg. `foo: String = 'bar'`
type MethodParameterNode struct {
	NodeBase
	Name                string         // name of the variable
	TypeNode            TypeNode       // type of the variable
	Initialiser         ExpressionNode // value assigned to the variable
	SetInstanceVariable bool           // whether an instance variable with this name gets automatically assigned
	Kind                ParameterKind
}

func (*MethodParameterNode) IsStatic() bool {
	return false
}

func (f *MethodParameterNode) IsOptional() bool {
	return f.Initialiser != nil
}

// Create a new formal parameter node eg. `foo: String = 'bar'`
func NewMethodParameterNode(span *position.Span, name string, setIvar bool, typ TypeNode, init ExpressionNode, kind ParameterKind) *MethodParameterNode {
	return &MethodParameterNode{
		NodeBase:            NodeBase{span: span},
		SetInstanceVariable: setIvar,
		Name:                name,
		TypeNode:            typ,
		Initialiser:         init,
		Kind:                kind,
	}
}

// Represents a signature parameter in method and function signatures eg. `foo?: String`
type SignatureParameterNode struct {
	NodeBase
	Name     string   // name of the variable
	TypeNode TypeNode // type of the variable
	Optional bool     // whether this parameter is optional
	Kind     ParameterKind
}

func (*SignatureParameterNode) IsStatic() bool {
	return false
}

func (f *SignatureParameterNode) IsOptional() bool {
	return f.Optional
}

// Create a new signature parameter node eg. `foo?: String`
func NewSignatureParameterNode(span *position.Span, name string, typ TypeNode, opt bool, kind ParameterKind) *SignatureParameterNode {
	return &SignatureParameterNode{
		NodeBase: NodeBase{span: span},
		Name:     name,
		TypeNode: typ,
		Optional: opt,
		Kind:     kind,
	}
}

// Represents an attribute declaration in getters, setters and accessors eg. `foo: String`
type AttributeParameterNode struct {
	NodeBase
	Name        string         // name of the variable
	TypeNode    TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
}

func (*AttributeParameterNode) IsStatic() bool {
	return false
}

func (a *AttributeParameterNode) IsOptional() bool {
	return a.Initialiser != nil
}

// Create a new attribute declaration in getters, setters and accessors eg. `foo: String`
func NewAttributeParameterNode(span *position.Span, name string, typ TypeNode, init ExpressionNode) *AttributeParameterNode {
	return &AttributeParameterNode{
		NodeBase:    NodeBase{span: span},
		Name:        name,
		TypeNode:    typ,
		Initialiser: init,
	}
}

// Represents a closure type eg. `|i: Int|: String`
type ClosureTypeNode struct {
	TypedNodeBase
	Parameters []ParameterNode // formal parameters of the closure separated by semicolons
	ReturnType TypeNode
	ThrowType  TypeNode
}

func (*ClosureTypeNode) IsStatic() bool {
	return false
}

// Create a new closure type node eg. `|i: Int|: String`
func NewClosureTypeNode(span *position.Span, params []ParameterNode, retType TypeNode, throwType TypeNode) *ClosureTypeNode {
	return &ClosureTypeNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Parameters:    params,
		ReturnType:    retType,
		ThrowType:     throwType,
	}
}

// Represents a closure eg. `|i| -> println(i)`
type ClosureLiteralNode struct {
	TypedNodeBase
	Parameters []ParameterNode // formal parameters of the closure separated by semicolons
	ReturnType TypeNode
	ThrowType  TypeNode
	Body       []StatementNode // body of the closure
}

func (*ClosureLiteralNode) IsStatic() bool {
	return false
}

// Create a new closure expression node eg. `|i| -> println(i)`
func NewClosureLiteralNode(span *position.Span, params []ParameterNode, retType TypeNode, throwType TypeNode, body []StatementNode) *ClosureLiteralNode {
	return &ClosureLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Parameters:    params,
		ReturnType:    retType,
		ThrowType:     throwType,
		Body:          body,
	}
}

// Represents a class declaration eg. `class Foo; end`
type ClassDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Abstract       bool
	Sealed         bool
	Primitive      bool
	Constant       ExpressionNode      // The constant that will hold the class value
	TypeParameters []TypeParameterNode // Generic type variable definitions
	Superclass     ExpressionNode      // the super/parent class of this class
	Body           []StatementNode     // body of the class
}

func (*ClassDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*ClassDeclarationNode) IsStatic() bool {
	return false
}

// Create a new class declaration node eg. `class Foo; end`
func NewClassDeclarationNode(
	span *position.Span,
	docComment string,
	abstract bool,
	sealed bool,
	primitive bool,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	superclass ExpressionNode,
	body []StatementNode,
) *ClassDeclarationNode {

	return &ClassDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Abstract:       abstract,
		Sealed:         sealed,
		Primitive:      primitive,
		Constant:       constant,
		TypeParameters: typeParams,
		Superclass:     superclass,
		Body:           body,
	}
}

// Represents a module declaration eg. `module Foo; end`
type ModuleDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant ExpressionNode  // The constant that will hold the module value
	Body     []StatementNode // body of the module
}

func (*ModuleDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*ModuleDeclarationNode) IsStatic() bool {
	return false
}

// Create a new module declaration node eg. `module Foo; end`
func NewModuleDeclarationNode(
	span *position.Span,
	docComment string,
	constant ExpressionNode,
	body []StatementNode,
) *ModuleDeclarationNode {

	return &ModuleDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant: constant,
		Body:     body,
	}
}

// Represents a mixin declaration eg. `mixin Foo; end`
type MixinDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Abstract              bool
	Constant              ExpressionNode      // The constant that will hold the mixin value
	TypeParameters        []TypeParameterNode // Generic type variable definitions
	Body                  []StatementNode     // body of the mixin
	IncludesAndImplements []ExpressionNode
}

func (*MixinDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*MixinDeclarationNode) IsStatic() bool {
	return false
}

// Create a new mixin declaration node eg. `mixin Foo; end`
func NewMixinDeclarationNode(
	span *position.Span,
	docComment string,
	abstract bool,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	body []StatementNode,
) *MixinDeclarationNode {

	return &MixinDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Abstract:       abstract,
		Constant:       constant,
		TypeParameters: typeParams,
		Body:           body,
	}
}

// Represents an interface declaration eg. `interface Foo; end`
type InterfaceDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant       ExpressionNode      // The constant that will hold the interface value
	TypeParameters []TypeParameterNode // Generic type variable definitions
	Body           []StatementNode     // body of the interface
	Implements     []*ImplementExpressionNode
}

func (*InterfaceDeclarationNode) SkipTypechecking() bool {
	return false
}

func (*InterfaceDeclarationNode) IsStatic() bool {
	return false
}

// Create a new interface declaration node eg. `interface Foo; end`
func NewInterfaceDeclarationNode(
	span *position.Span,
	docComment string,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	body []StatementNode,
) *InterfaceDeclarationNode {

	return &InterfaceDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:       constant,
		TypeParameters: typeParams,
		Body:           body,
	}
}

// Represents a struct declaration eg. `struct Foo; end`
type StructDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant       ExpressionNode            // The constant that will hold the struct value
	TypeParameters []TypeParameterNode       // Generic type variable definitions
	Body           []StructBodyStatementNode // body of the struct
}

func (*StructDeclarationNode) IsStatic() bool {
	return false
}

// Create a new struct declaration node eg. `struct Foo; end`
func NewStructDeclarationNode(
	span *position.Span,
	docComment string,
	constant ExpressionNode,
	typeParams []TypeParameterNode,
	body []StructBodyStatementNode,
) *StructDeclarationNode {

	return &StructDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:       constant,
		TypeParameters: typeParams,
		Body:           body,
	}
}

// Represents the variance of a type parameter.
type Variance uint8

const (
	INVARIANT Variance = iota
	COVARIANT
	CONTRAVARIANT
)

// Represents a type parameter eg. `+V`
type VariantTypeParameterNode struct {
	TypedNodeBase
	Variance   Variance // Variance level of this type parameter
	Name       string   // Name of the type parameter eg. `T`
	UpperBound TypeNode
	LowerBound TypeNode
}

func (*VariantTypeParameterNode) IsStatic() bool {
	return false
}

// Create a new type variable node eg. `+V`
func NewVariantTypeParameterNode(span *position.Span, variance Variance, name string, lower, upper TypeNode) *VariantTypeParameterNode {
	return &VariantTypeParameterNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Variance:      variance,
		Name:          name,
		LowerBound:    lower,
		UpperBound:    upper,
	}
}

// Represents a symbol literal with simple content eg. `:foo`, `:'foo bar`, `:"lol"`
type SimpleSymbolLiteralNode struct {
	TypedNodeBase
	Content string
}

func (*SimpleSymbolLiteralNode) IsStatic() bool {
	return true
}

// Create a simple symbol literal node eg. `:foo`, `:'foo bar`, `:"lol"`
func NewSimpleSymbolLiteralNode(span *position.Span, cont string) *SimpleSymbolLiteralNode {
	return &SimpleSymbolLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Content:       cont,
	}
}

// Represents an interpolated symbol eg. `:"foo ${bar + 2}"`
type InterpolatedSymbolLiteralNode struct {
	NodeBase
	Content *InterpolatedStringLiteralNode
}

func (*InterpolatedSymbolLiteralNode) IsStatic() bool {
	return false
}

func (*InterpolatedSymbolLiteralNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return globalEnv.StdSubtype(symbol.Symbol)
}

// Create an interpolated symbol literal node eg. `:"foo ${bar + 2}"`
func NewInterpolatedSymbolLiteralNode(span *position.Span, cont *InterpolatedStringLiteralNode) *InterpolatedSymbolLiteralNode {
	return &InterpolatedSymbolLiteralNode{
		NodeBase: NodeBase{span: span},
		Content:  cont,
	}
}

// Represents a method definition eg. `def foo: String then 'hello world'`
type MethodDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name           string
	TypeParameters []TypeParameterNode
	Parameters     []ParameterNode // formal parameters
	ReturnType     TypeNode
	ThrowType      TypeNode
	Body           []StatementNode // body of the method
	Sealed         bool
	Abstract       bool
}

func (*MethodDefinitionNode) IsStatic() bool {
	return false
}

// Whether the method is a setter.
func (m *MethodDefinitionNode) IsSetter() bool {
	if len(m.Name) > 0 {
		firstChar, _ := utf8.DecodeRuneInString(m.Name)
		lastChar := m.Name[len(m.Name)-1]
		if (unicode.IsLetter(firstChar) || firstChar == '_') && lastChar == '=' {
			return true
		}
	}

	return false
}

// Create a method definition node eg. `def foo: String then 'hello world'`
func NewMethodDefinitionNode(
	span *position.Span,
	docComment string,
	abstract bool,
	sealed bool,
	name string,
	typeParams []TypeParameterNode,
	params []ParameterNode,
	returnType,
	throwType TypeNode,
	body []StatementNode,
) *MethodDefinitionNode {
	return &MethodDefinitionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Abstract:       abstract,
		Sealed:         sealed,
		Name:           name,
		TypeParameters: typeParams,
		Parameters:     params,
		ReturnType:     returnType,
		ThrowType:      throwType,
		Body:           body,
	}
}

// Represents a constructor definition eg. `init then 'hello world'`
type InitDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Parameters []ParameterNode // formal parameters
	ThrowType  TypeNode
	Body       []StatementNode // body of the method
}

func (*InitDefinitionNode) IsStatic() bool {
	return false
}

// Create a constructor definition node eg. `init then 'hello world'`
func NewInitDefinitionNode(span *position.Span, params []ParameterNode, throwType TypeNode, body []StatementNode) *InitDefinitionNode {
	return &InitDefinitionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Parameters:    params,
		ThrowType:     throwType,
		Body:          body,
	}
}

// Represents a method signature definition eg. `sig to_string(val: Int): String`
type MethodSignatureDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Name           string
	TypeParameters []TypeParameterNode
	Parameters     []ParameterNode // formal parameters
	ReturnType     TypeNode
	ThrowType      TypeNode
}

func (*MethodSignatureDefinitionNode) IsStatic() bool {
	return false
}

// Create a method signature node eg. `sig to_string(val: Int): String`
func NewMethodSignatureDefinitionNode(span *position.Span, docComment, name string, typeParams []TypeParameterNode, params []ParameterNode, returnType, throwType TypeNode) *MethodSignatureDefinitionNode {
	return &MethodSignatureDefinitionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Name:           name,
		TypeParameters: typeParams,
		Parameters:     params,
		ReturnType:     returnType,
		ThrowType:      throwType,
	}
}

// Represents a generic constant in type annotations eg. `ArrayList[String]`
type GenericConstantNode struct {
	TypedNodeBase
	Constant      ComplexConstantNode
	TypeArguments []TypeNode
}

func (*GenericConstantNode) IsStatic() bool {
	return true
}

// Create a generic constant node eg. `ArrayList[String]`
func NewGenericConstantNode(span *position.Span, constant ComplexConstantNode, args []TypeNode) *GenericConstantNode {
	return &GenericConstantNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Constant:      constant,
		TypeArguments: args,
	}
}

// Represents a new generic type definition eg. `typedef Nilable[T] = T | nil`
type GenericTypeDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	TypeParameters []TypeParameterNode // Generic type variable definitions
	Constant       ComplexConstantNode // new name of the type
	TypeNode       TypeNode            // the type
}

func (*GenericTypeDefinitionNode) IsStatic() bool {
	return false
}

// Create a generic type definition node eg. `typedef Nilable[T] = T | nil`
func NewGenericTypeDefinitionNode(span *position.Span, docComment string, constant ComplexConstantNode, typeVars []TypeParameterNode, typ TypeNode) *GenericTypeDefinitionNode {
	return &GenericTypeDefinitionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant:       constant,
		TypeParameters: typeVars,
		TypeNode:       typ,
	}
}

// Represents a new type definition eg. `typedef StringList = ArrayList[String]`
type TypeDefinitionNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Constant ComplexConstantNode // new name of the type
	TypeNode TypeNode            // the type
}

func (*TypeDefinitionNode) IsStatic() bool {
	return false
}

// Create a type definition node eg. `typedef StringList = ArrayList[String]`
func NewTypeDefinitionNode(span *position.Span, docComment string, constant ComplexConstantNode, typ TypeNode) *TypeDefinitionNode {
	return &TypeDefinitionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Constant: constant,
		TypeNode: typ,
	}
}

// A single alias entry eg. `new_name old_name`
type AliasDeclarationEntry struct {
	NodeBase
	NewName string
	OldName string
}

func (*AliasDeclarationEntry) IsStatic() bool {
	return false
}

// Create an alias alias entry eg. `new_name old_name`
func NewAliasDeclarationEntry(span *position.Span, newName, oldName string) *AliasDeclarationEntry {
	return &AliasDeclarationEntry{
		NodeBase: NodeBase{span: span},
		NewName:  newName,
		OldName:  oldName,
	}
}

// Represents a new alias declaration eg. `alias push append, add plus`
type AliasDeclarationNode struct {
	TypedNodeBase
	Entries []*AliasDeclarationEntry
}

func (*AliasDeclarationNode) IsStatic() bool {
	return false
}

// Create an alias declaration node eg. `alias push append, add plus`
func NewAliasDeclarationNode(span *position.Span, entries []*AliasDeclarationEntry) *AliasDeclarationNode {
	return &AliasDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Entries:       entries,
	}
}

// Represents a new getter declaration eg. `getter foo: String`
type GetterDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Entries []ParameterNode
}

func (*GetterDeclarationNode) IsStatic() bool {
	return false
}

// Create a getter declaration node eg. `getter foo: String`
func NewGetterDeclarationNode(span *position.Span, docComment string, entries []ParameterNode) *GetterDeclarationNode {
	return &GetterDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Entries: entries,
	}
}

// Represents a new setter declaration eg. `setter foo: String`
type SetterDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Entries []ParameterNode
}

func (*SetterDeclarationNode) IsStatic() bool {
	return false
}

// Create a setter declaration node eg. `setter foo: String`
func NewSetterDeclarationNode(span *position.Span, docComment string, entries []ParameterNode) *SetterDeclarationNode {
	return &SetterDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Entries: entries,
	}
}

// Represents a new setter declaration eg. `attr foo: String`
type AttrDeclarationNode struct {
	TypedNodeBase
	DocCommentableNodeBase
	Entries []ParameterNode
}

func (*AttrDeclarationNode) IsStatic() bool {
	return false
}

// Create an attribute declaration node eg. `attr foo: String`
func NewAttrDeclarationNode(span *position.Span, docComment string, entries []ParameterNode) *AttrDeclarationNode {
	return &AttrDeclarationNode{
		TypedNodeBase: TypedNodeBase{span: span},
		DocCommentableNodeBase: DocCommentableNodeBase{
			comment: docComment,
		},
		Entries: entries,
	}
}

// Represents a using expression eg. `using Foo`
type UsingExpressionNode struct {
	TypedNodeBase
	Entries []UsingEntryNode
}

func (*UsingExpressionNode) SkipTypechecking() bool {
	return false
}

func (*UsingExpressionNode) IsStatic() bool {
	return false
}

// Create a using expression node eg. `using Foo`
func NewUsingExpressionNode(span *position.Span, consts []UsingEntryNode) *UsingExpressionNode {
	return &UsingExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Entries:       consts,
	}
}

// Represents an include expression eg. `include Enumerable[V]`
type IncludeExpressionNode struct {
	TypedNodeBase
	Constants []ComplexConstantNode
}

func (*IncludeExpressionNode) SkipTypechecking() bool {
	return false
}

func (*IncludeExpressionNode) IsStatic() bool {
	return false
}

// Create an include expression node eg. `include Enumerable[V]`
func NewIncludeExpressionNode(span *position.Span, consts []ComplexConstantNode) *IncludeExpressionNode {
	return &IncludeExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Constants:     consts,
	}
}

// Represents an extend expression eg. `extend Enumerable[V]`
type ExtendExpressionNode struct {
	NodeBase
	Constants []ComplexConstantNode
}

func (*ExtendExpressionNode) IsStatic() bool {
	return false
}

// Create an extend expression node eg. `extend Enumerable[V]`
func NewExtendExpressionNode(span *position.Span, consts []ComplexConstantNode) *ExtendExpressionNode {
	return &ExtendExpressionNode{
		NodeBase:  NodeBase{span: span},
		Constants: consts,
	}
}

// Represents an enhance expression eg. `enhance Enumerable[V]`
type EnhanceExpressionNode struct {
	NodeBase
	Constants []ComplexConstantNode
}

func (*EnhanceExpressionNode) IsStatic() bool {
	return false
}

// Create an enhance expression node eg. `enhance Enumerable[V]`
func NewEnhanceExpressionNode(span *position.Span, consts []ComplexConstantNode) *EnhanceExpressionNode {
	return &EnhanceExpressionNode{
		NodeBase:  NodeBase{span: span},
		Constants: consts,
	}
}

// Represents an enhance expression eg. `implement Enumerable[V]`
type ImplementExpressionNode struct {
	TypedNodeBase
	Constants []ComplexConstantNode
}

func (*ImplementExpressionNode) SkipTypechecking() bool {
	return false
}

func (*ImplementExpressionNode) IsStatic() bool {
	return false
}

// Create an enhance expression node eg. `implement Enumerable[V]`
func NewImplementExpressionNode(span *position.Span, consts []ComplexConstantNode) *ImplementExpressionNode {
	return &ImplementExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Constants:     consts,
	}
}

// Represents a named argument in a function call eg. `foo: 123`
type NamedCallArgumentNode struct {
	NodeBase
	Name  string
	Value ExpressionNode
}

func (*NamedCallArgumentNode) IsStatic() bool {
	return false
}

// Create a named argument node eg. `foo: 123`
func NewNamedCallArgumentNode(span *position.Span, name string, val ExpressionNode) *NamedCallArgumentNode {
	return &NamedCallArgumentNode{
		NodeBase: NodeBase{span: span},
		Name:     name,
		Value:    val,
	}
}

// Represents a new expression eg. `new(123)`
type NewExpressionNode struct {
	TypedNodeBase
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*NewExpressionNode) IsStatic() bool {
	return false
}

// Create a new expression node eg. `new(123)`
func NewNewExpressionNode(span *position.Span, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *NewExpressionNode {
	return &NewExpressionNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a constructor call eg. `String(123)`
type ConstructorCallNode struct {
	TypedNodeBase
	Class               ComplexConstantNode // class that is being instantiated
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*ConstructorCallNode) IsStatic() bool {
	return false
}

// Create a constructor call node eg. `String(123)`
func NewConstructorCallNode(span *position.Span, class ComplexConstantNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ConstructorCallNode {
	return &ConstructorCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		Class:               class,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a constructor call eg. `ArrayList::[Int](1, 2, 3)`
type GenericConstructorCallNode struct {
	TypedNodeBase
	Class               ComplexConstantNode // class that is being instantiated
	TypeArguments       []TypeNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*GenericConstructorCallNode) IsStatic() bool {
	return false
}

// Create a constructor call node eg. `ArrayList::[Int](1, 2, 3)`
func NewGenericConstructorCallNode(span *position.Span, class ComplexConstantNode, typeArgs []TypeNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *GenericConstructorCallNode {
	return &GenericConstructorCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		Class:               class,
		TypeArguments:       typeArgs,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents attribute access eg. `foo.bar`
type AttributeAccessNode struct {
	TypedNodeBase
	Receiver      ExpressionNode
	AttributeName string
}

func (*AttributeAccessNode) IsStatic() bool {
	return false
}

// Create an attribute access node eg. `foo.bar`
func NewAttributeAccessNode(span *position.Span, recv ExpressionNode, attrName string) *AttributeAccessNode {
	return &AttributeAccessNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Receiver:      recv,
		AttributeName: attrName,
	}
}

// Represents subscript access eg. `arr[5]`
type SubscriptExpressionNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Key      ExpressionNode
	static   bool
}

func (s *SubscriptExpressionNode) IsStatic() bool {
	return s.static
}

// Create a subscript expression node eg. `arr[5]`
func NewSubscriptExpressionNode(span *position.Span, recv, key ExpressionNode) *SubscriptExpressionNode {
	return &SubscriptExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Receiver:      recv,
		Key:           key,
		static:        recv.IsStatic() && key.IsStatic(),
	}
}

// Represents nil-safe subscript access eg. `arr?[5]`
type NilSafeSubscriptExpressionNode struct {
	TypedNodeBase
	Receiver ExpressionNode
	Key      ExpressionNode
	static   bool
}

func (s *NilSafeSubscriptExpressionNode) IsStatic() bool {
	return s.static
}

// Create a nil-safe subscript expression node eg. `arr?[5]`
func NewNilSafeSubscriptExpressionNode(span *position.Span, recv, key ExpressionNode) *NilSafeSubscriptExpressionNode {
	return &NilSafeSubscriptExpressionNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Receiver:      recv,
		Key:           key,
		static:        recv.IsStatic() && key.IsStatic(),
	}
}

// Represents a method call eg. `'123'.()`
type CallNode struct {
	NodeBase
	Receiver            ExpressionNode
	NilSafe             bool
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*CallNode) IsStatic() bool {
	return false
}

// Create a method call node eg. `'123'.to_int()`
func NewCallNode(span *position.Span, recv ExpressionNode, nilSafe bool, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *CallNode {
	return &CallNode{
		NodeBase:            NodeBase{span: span},
		Receiver:            recv,
		NilSafe:             nilSafe,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a method call eg. `foo.bar::[String](a)`
type GenericMethodCallNode struct {
	TypedNodeBase
	Receiver            ExpressionNode
	Op                  *token.Token
	MethodName          string
	TypeArguments       []TypeNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*GenericMethodCallNode) IsStatic() bool {
	return false
}

// Create a method call node eg. `foo.bar::[String](a)`
func NewGenericMethodCallNode(span *position.Span, recv ExpressionNode, op *token.Token, methodName string, typeArgs []TypeNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *GenericMethodCallNode {
	return &GenericMethodCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		Receiver:            recv,
		Op:                  op,
		MethodName:          methodName,
		TypeArguments:       typeArgs,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a method call eg. `'123'.to_int()`
type MethodCallNode struct {
	TypedNodeBase
	Receiver            ExpressionNode
	Op                  *token.Token
	MethodName          string
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*MethodCallNode) IsStatic() bool {
	return false
}

// Create a method call node eg. `'123'.to_int()`
func NewMethodCallNode(span *position.Span, recv ExpressionNode, op *token.Token, methodName string, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *MethodCallNode {
	return &MethodCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		Receiver:            recv,
		Op:                  op,
		MethodName:          methodName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a function-like call eg. `to_string(123)`
type ReceiverlessMethodCallNode struct {
	TypedNodeBase
	MethodName          string
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*ReceiverlessMethodCallNode) IsStatic() bool {
	return false
}

// Create a function call node eg. `to_string(123)`
func NewReceiverlessMethodCallNode(span *position.Span, methodName string, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ReceiverlessMethodCallNode {
	return &ReceiverlessMethodCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		MethodName:          methodName,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a generic function-like call eg. `foo::[Int](123)`
type GenericReceiverlessMethodCallNode struct {
	TypedNodeBase
	MethodName          string
	TypeArguments       []TypeNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (*GenericReceiverlessMethodCallNode) IsStatic() bool {
	return false
}

// Create a generic function call node eg. `foo::[Int](123)`
func NewGenericReceiverlessMethodCallNode(span *position.Span, methodName string, typeArgs []TypeNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *GenericReceiverlessMethodCallNode {
	return &GenericReceiverlessMethodCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		MethodName:          methodName,
		TypeArguments:       typeArgs,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a symbol value expression eg. `foo: bar`
type SymbolKeyValueExpressionNode struct {
	NodeBase
	Key   string
	Value ExpressionNode
}

func (s *SymbolKeyValueExpressionNode) IsStatic() bool {
	return s.Value.IsStatic()
}

// Create a symbol key value node eg. `foo: bar`
func NewSymbolKeyValueExpressionNode(span *position.Span, key string, val ExpressionNode) *SymbolKeyValueExpressionNode {
	return &SymbolKeyValueExpressionNode{
		NodeBase: NodeBase{span: span},
		Key:      key,
		Value:    val,
	}
}

// Represents a key value expression eg. `foo => bar`
type KeyValueExpressionNode struct {
	NodeBase
	Key    ExpressionNode
	Value  ExpressionNode
	static bool
}

func (k *KeyValueExpressionNode) IsStatic() bool {
	return k.static
}

// Create a key value expression node eg. `foo => bar`
func NewKeyValueExpressionNode(span *position.Span, key, val ExpressionNode) *KeyValueExpressionNode {
	return &KeyValueExpressionNode{
		NodeBase: NodeBase{span: span},
		Key:      key,
		Value:    val,
		static:   areExpressionsStatic(key, val),
	}
}

// Represents a ArrayList literal eg. `[1, 5, -6]`
type ArrayListLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	Capacity ExpressionNode
	static   bool
}

func (l *ArrayListLiteralNode) IsStatic() bool {
	return l.static
}

// Create a ArrayList literal node eg. `[1, 5, -6]`
func NewArrayListLiteralNode(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) *ArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = isExpressionSliceStatic(elements) && capacity.IsStatic()
	} else {
		static = isExpressionSliceStatic(elements)
	}
	return &ArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewArrayListLiteralNode] but returns an interface
func NewArrayListLiteralNodeI(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewArrayListLiteralNode(span, elements, capacity)
}

// Represents a word ArrayList literal eg. `\w[foo bar]`
type WordArrayListLiteralNode struct {
	TypedNodeBase
	Elements []WordCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (w *WordArrayListLiteralNode) IsStatic() bool {
	return w.static
}

// Create a word ArrayList literal node eg. `\w[foo bar]`
func NewWordArrayListLiteralNode(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) *WordArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &WordArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewWordArrayListLiteralNode] but returns an interface.
func NewWordArrayListLiteralExpressionNode(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewWordArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewWordArrayListLiteralNode] but returns an interface.
func NewWordArrayListLiteralPatternExpressionNode(span *position.Span, elements []WordCollectionContentNode) PatternExpressionNode {
	return NewWordArrayListLiteralNode(span, elements, nil)
}

// Represents a word ArrayTuple literal eg. `%w[foo bar]`
type WordArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []WordCollectionContentNode
}

func (*WordArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a word ArrayTuple literal node eg. `%w[foo bar]`
func NewWordArrayTupleLiteralNode(span *position.Span, elements []WordCollectionContentNode) *WordArrayTupleLiteralNode {
	return &WordArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewWordArrayTupleLiteralNode] but returns an interface.
func NewWordArrayTupleLiteralExpressionNode(span *position.Span, elements []WordCollectionContentNode) ExpressionNode {
	return NewWordArrayTupleLiteralNode(span, elements)
}

// Same as [NewWordArrayTupleLiteralNode] but returns an interface.
func NewWordArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []WordCollectionContentNode) PatternExpressionNode {
	return NewWordArrayTupleLiteralNode(span, elements)
}

// Represents a word HashSet literal eg. `^w[foo bar]`
type WordHashSetLiteralNode struct {
	TypedNodeBase
	Elements []WordCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (w *WordHashSetLiteralNode) IsStatic() bool {
	return w.static
}

// Create a word HashSet literal node eg. `^w[foo bar]`
func NewWordHashSetLiteralNode(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) *WordHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &WordHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewWordHashSetLiteralNode] but returns an interface.
func NewWordHashSetLiteralNodeI(span *position.Span, elements []WordCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewWordHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewWordHashSetLiteralNode] but returns an interface.
func NewWordHashSetLiteralPatternExpressionNode(span *position.Span, elements []WordCollectionContentNode) PatternExpressionNode {
	return NewWordHashSetLiteralNode(span, elements, nil)
}

// Represents a symbol ArrayList literal eg. `\s[foo bar]`
type SymbolArrayListLiteralNode struct {
	TypedNodeBase
	Elements []SymbolCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (s *SymbolArrayListLiteralNode) IsStatic() bool {
	return s.static
}

// Create a symbol ArrayList literal node eg. `\s[foo bar]`
func NewSymbolArrayListLiteralNode(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) *SymbolArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &SymbolArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewSymbolArrayListLiteralNode] but returns an interface.
func NewSymbolArrayListLiteralExpressionNode(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewSymbolArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewSymbolArrayListLiteralNode] but returns an interface.
func NewSymbolArrayListLiteralPatternExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) PatternExpressionNode {
	return NewSymbolArrayListLiteralNode(span, elements, nil)
}

// Represents a symbol ArrayTuple literal eg. `%s[foo bar]`
type SymbolArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []SymbolCollectionContentNode
}

func (*SymbolArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a symbol arrayTuple literal node eg. `%s[foo bar]`
func NewSymbolArrayTupleLiteralNode(span *position.Span, elements []SymbolCollectionContentNode) *SymbolArrayTupleLiteralNode {
	return &SymbolArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewSymbolArrayTupleLiteralNode] but returns an interface.
func NewSymbolArrayTupleLiteralExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) ExpressionNode {
	return NewSymbolArrayTupleLiteralNode(span, elements)
}

// Same as [NewSymbolArrayTupleLiteralNode] but returns an interface.
func NewSymbolArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) PatternExpressionNode {
	return NewSymbolArrayTupleLiteralNode(span, elements)
}

// Represents a symbol HashSet literal eg. `^s[foo bar]`
type SymbolHashSetLiteralNode struct {
	TypedNodeBase
	Elements []SymbolCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (s *SymbolHashSetLiteralNode) IsStatic() bool {
	return s.static
}

// Create a symbol HashSet literal node eg. `^s[foo bar]`
func NewSymbolHashSetLiteralNode(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) *SymbolHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &SymbolHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewSymbolHashSetLiteralNode] but returns an interface.
func NewSymbolHashSetLiteralNodeI(span *position.Span, elements []SymbolCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewSymbolHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewSymbolHashSetLiteralNode] but returns an interface.
func NewSymbolHashSetLiteralPatternExpressionNode(span *position.Span, elements []SymbolCollectionContentNode) PatternExpressionNode {
	return NewSymbolHashSetLiteralNode(span, elements, nil)
}

// Represents a hex ArrayList literal eg. `\x[ff ee]`
type HexArrayListLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (h *HexArrayListLiteralNode) IsStatic() bool {
	return h.static
}

// Create a hex ArrayList literal node eg. `\x[ff ee]`
func NewHexArrayListLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *HexArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &HexArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHexArrayListLiteralNode] but returns an interface.
func NewHexArrayListLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewHexArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewHexArrayListLiteralNode] but returns an interface.
func NewHexArrayListLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewHexArrayListLiteralNode(span, elements, nil)
}

// Represents a hex ArrayTuple literal eg. `%x[ff ee]`
type HexArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
}

func (*HexArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a hex ArrayTuple literal node eg. `%x[ff ee]`
func NewHexArrayTupleLiteralNode(span *position.Span, elements []IntCollectionContentNode) *HexArrayTupleLiteralNode {
	return &HexArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewHexArrayTupleLiteralNode] but returns an interface.
func NewHexArrayTupleLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode) ExpressionNode {
	return NewHexArrayTupleLiteralNode(span, elements)
}

// Same as [NewHexArrayTupleLiteralNode] but returns an interface.
func NewHexArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewHexArrayTupleLiteralNode(span, elements)
}

// Represents a hex HashSet literal eg. `^x[ff ee}]`
type HexHashSetLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (h *HexHashSetLiteralNode) IsStatic() bool {
	return h.static
}

// Create a hex HashSet literal node eg. `^x[ff ee]`
func NewHexHashSetLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *HexHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &HexHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHexHashSetLiteralNode] but returns an interface.
func NewHexHashSetLiteralNodeI(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewHexHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewHexHashSetLiteralNode] but returns an interface.
func NewHexHashSetLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewHexHashSetLiteralNode(span, elements, nil)
}

// Represents a bin ArrayList literal eg. `\b[11 10]`
type BinArrayListLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (b *BinArrayListLiteralNode) IsStatic() bool {
	return b.static
}

// Create a bin ArrayList literal node eg. `\b[11 10]`
func NewBinArrayListLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *BinArrayListLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &BinArrayListLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewBinArrayListLiteralNode] but returns an interface.
func NewBinArrayListLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewBinArrayListLiteralNode(span, elements, capacity)
}

// Same as [NewBinArrayListLiteralNode] but returns an interface.
func NewBinArrayListLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewBinArrayListLiteralNode(span, elements, nil)
}

// Represents a bin ArrayTuple literal eg. `%b[11 10]`
type BinArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
}

func (*BinArrayTupleLiteralNode) IsStatic() bool {
	return true
}

// Create a bin ArrayList literal node eg. `%b[11 10]`
func NewBinArrayTupleLiteralNode(span *position.Span, elements []IntCollectionContentNode) *BinArrayTupleLiteralNode {
	return &BinArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
	}
}

// Same as [NewBinArrayTupleLiteralNode] but returns an interface.
func NewBinArrayTupleLiteralExpressionNode(span *position.Span, elements []IntCollectionContentNode) ExpressionNode {
	return NewBinArrayTupleLiteralNode(span, elements)
}

// Same as [NewBinArrayTupleLiteralNode] but returns an interface.
func NewBinArrayTupleLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewBinArrayTupleLiteralNode(span, elements)
}

// Represents a bin HashSet literal eg. `^b[11 10]`
type BinHashSetLiteralNode struct {
	TypedNodeBase
	Elements []IntCollectionContentNode
	Capacity ExpressionNode
	static   bool
}

func (b *BinHashSetLiteralNode) IsStatic() bool {
	return b.static
}

// Create a bin HashSet literal node eg. `^b[11 10]`
func NewBinHashSetLiteralNode(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) *BinHashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = capacity.IsStatic()
	} else {
		static = true
	}
	return &BinHashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewBinHashSetLiteralNode] but returns an interface.
func NewBinHashSetLiteralNodeI(span *position.Span, elements []IntCollectionContentNode, capacity ExpressionNode) ExpressionNode {
	return NewBinHashSetLiteralNode(span, elements, capacity)
}

// Same as [NewBinHashSetLiteralNode] but returns an interface.
func NewBinHashSetLiteralPatternExpressionNode(span *position.Span, elements []IntCollectionContentNode) PatternExpressionNode {
	return NewBinHashSetLiteralNode(span, elements, nil)
}

// Represents a ArrayTuple literal eg. `%[1, 5, -6]`
type ArrayTupleLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	static   bool
}

func (t *ArrayTupleLiteralNode) IsStatic() bool {
	return t.static
}

// Create a ArrayTuple literal node eg. `%[1, 5, -6]`
func NewArrayTupleLiteralNode(span *position.Span, elements []ExpressionNode) *ArrayTupleLiteralNode {
	return &ArrayTupleLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		static:        isExpressionSliceStatic(elements),
	}
}

// Same as [NewArrayTupleLiteralNode] but returns an interface
func NewArrayTupleLiteralNodeI(span *position.Span, elements []ExpressionNode) ExpressionNode {
	return NewArrayTupleLiteralNode(span, elements)
}

// Represents a HashSet literal eg. `^[1, 5, -6]`
type HashSetLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	Capacity ExpressionNode
	static   bool
}

func (s *HashSetLiteralNode) IsStatic() bool {
	return s.static
}

// Create a HashSet literal node eg. `^[1, 5, -6]`
func NewHashSetLiteralNode(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) *HashSetLiteralNode {
	var static bool
	if capacity != nil {
		static = isExpressionSliceStatic(elements) && capacity.IsStatic()
	} else {
		static = isExpressionSliceStatic(elements)
	}
	return &HashSetLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHashSetLiteralNode] but returns an interface
func NewHashSetLiteralNodeI(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewHashSetLiteralNode(span, elements, capacity)
}

// Represents a HashMap literal eg. `{ foo: 1, 'bar' => 5, baz }`
type HashMapLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	Capacity ExpressionNode
	static   bool
}

func (m *HashMapLiteralNode) IsStatic() bool {
	return m.static
}

// Create a HashMap literal node eg. `{ foo: 1, 'bar' => 5, baz }`
func NewHashMapLiteralNode(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) *HashMapLiteralNode {
	var static bool
	if capacity != nil {
		static = isExpressionSliceStatic(elements) && capacity.IsStatic()
	} else {
		static = isExpressionSliceStatic(elements)
	}
	return &HashMapLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		Capacity:      capacity,
		static:        static,
	}
}

// Same as [NewHashMapLiteralNode] but returns an interface
func NewHashMapLiteralNodeI(span *position.Span, elements []ExpressionNode, capacity ExpressionNode) ExpressionNode {
	return NewHashMapLiteralNode(span, elements, capacity)
}

// Represents a Record literal eg. `%{ foo: 1, 'bar' => 5, baz }`
type HashRecordLiteralNode struct {
	TypedNodeBase
	Elements []ExpressionNode
	static   bool
}

func (r *HashRecordLiteralNode) IsStatic() bool {
	return r.static
}

// Create a Record literal node eg. `%{ foo: 1, 'bar' => 5, baz }`
func NewHashRecordLiteralNode(span *position.Span, elements []ExpressionNode) *HashRecordLiteralNode {
	return &HashRecordLiteralNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Elements:      elements,
		static:        isExpressionSliceStatic(elements),
	}
}

// Same as [NewHashRecordLiteralNode] but returns an interface
func NewHashRecordLiteralNodeI(span *position.Span, elements []ExpressionNode) ExpressionNode {
	return NewHashRecordLiteralNode(span, elements)
}

// Represents a Range literal eg. `1...5`
type RangeLiteralNode struct {
	NodeBase
	From   ExpressionNode
	To     ExpressionNode
	Op     *token.Token
	static bool
}

func (r *RangeLiteralNode) IsStatic() bool {
	return r.static
}

// Create a Range literal node eg. `1...5`
func NewRangeLiteralNode(span *position.Span, op *token.Token, from, to ExpressionNode) *RangeLiteralNode {
	return &RangeLiteralNode{
		NodeBase: NodeBase{span: span},
		Op:       op,
		From:     from,
		To:       to,
		static:   areExpressionsStatic(from, to),
	}
}
