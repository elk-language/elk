package value

import (
	"github.com/elk-language/elk/regex/flag"
)

var ElkASTModule *Module // Std::Elk::AST

var NodeMixin *Mixin                        // Std::Elk::AST::Node
var StructBodyStatementNodeMixin *Mixin     // Std::Elk::AST::StructBodyStatementNode
var StatementNodeMixin *Mixin               // Std::Elk::AST::StatementNode
var ParameterNodeMixin *Mixin               // Std::Elk::AST::ParameterNode
var TypeNodeMixin *Mixin                    // Std::Elk::AST::TypeNode
var PatternNodeMixin *Mixin                 // Std::Elk::AST::PatternNode
var ExpressionNodeMixin *Mixin              // Std::Elk::AST::ExpressionNode
var SymbolLiteralNodeMixin *Mixin           // Std::Elk::AST::SymbolLiteralNode
var ComplexConstantNodeMixin *Mixin         // Std::Elk::AST::ComplexConstantNode
var ConstantNodeMixin *Mixin                // Std::Elk::AST::ConstantNode
var UsingEntryNodeMixin *Mixin              // Std::Elk::AST::UsingEntryNode
var UsingSubentryNodeMixin *Mixin           // Std::Elk::AST::UsingSubentryNode
var NamedArgumentNodeMixin *Mixin           // Std::Elk::AST::NamedArgumentNode
var WordCollectionContentNodeMixin *Mixin   // Std::Elk::AST::WordCollectionContentNode
var SymbolCollectionContentNodeMixin *Mixin // Std::Elk::AST::SymbolCollectionContentNode
var IntCollectionContentNodeMixin *Mixin    // Std::Elk::AST::IntCollectionContentNode
var IdentifierNodeMixin *Mixin              // Std::Elk::AST::IdentifierNode
var RegexLiteralContentNodeMixin *Mixin     // Std::Elk::AST::RegexLiteralContentNode
var RegexLiteralNodeMixin *Mixin            // Std::Elk::AST::RegexLiteralNode
var StringLiteralContentNodeMixin *Mixin    // Std::Elk::AST::StringLiteralContentNode
var StringLiteralNodeMixin *Mixin           // Std::Elk::AST::StringLiteralNode
var SimpleStringLiteralNodeMixin *Mixin     // Std::Elk::AST::SimpleStringLiteralNode
var TypeParameterNodeMixin *Mixin           // Std::Elk::AST::TypeParameterNode
var PatternExpressionNodeMixin *Mixin       // Std::Elk::AST::PatternExpressionNode
var InstanceVariableNodeMixin *Mixin        // Std::Elk::AST::InstanceVariableNode

var NodeFormatErrorClass *Class                  // Std::Elk::AST::Node::FormatError
var IdentifierNodeFormatErrorClass *Class        // Std::Elk::AST::IdentifierNode::FormatError
var PrivateIdentifierNodeFormatErrorClass *Class // Std::Elk::AST::PrivateIdentifierNode::FormatError
var ConstantNodeFormatErrorClass *Class          // Std::Elk::AST::ConstantNode::FormatError
var PrivateConstantNodeFormatErrorClass *Class   // Std::Elk::AST::PrivateConstantNode::FormatError
var IntLiteralNodeFormatErrorClass *Class        // Std::Elk::AST::IntLiteralNode::FormatError
var Int8LiteralNodeFormatErrorClass *Class       // Std::Elk::AST::Int8LiteralNode::FormatError
var Int16LiteralNodeFormatErrorClass *Class      // Std::Elk::AST::Int16LiteralNode::FormatError
var Int32LiteralNodeFormatErrorClass *Class      // Std::Elk::AST::Int32LiteralNode::FormatError
var Int64LiteralNodeFormatErrorClass *Class      // Std::Elk::AST::Int64LiteralNode::FormatError
var UInt8LiteralNodeFormatErrorClass *Class      // Std::Elk::AST::UInt8LiteralNode::FormatError
var UInt16LiteralNodeFormatErrorClass *Class     // Std::Elk::AST::UInt16LiteralNode::FormatError
var UInt32LiteralNodeFormatErrorClass *Class     // Std::Elk::AST::UInt32LiteralNode::FormatError
var UInt64LiteralNodeFormatErrorClass *Class     // Std::Elk::AST::UInt64LiteralNode::FormatError
var FloatLiteralNodeFormatErrorClass *Class      // Std::Elk::AST::FloatLiteralNode::FormatError
var BigFloatLiteralNodeFormatErrorClass *Class   // Std::Elk::AST::BigFloatLiteralNode::FormatError
var Float32LiteralNodeFormatErrorClass *Class    // Std::Elk::AST::Float32LiteralNode::FormatError
var Float64LiteralNodeFormatErrorClass *Class    // Std::Elk::AST::Float64LiteralNode::FormatError

var ExpressionStatementNodeClass *Class // Std::Elk::AST::ExpressionStatementNode
var EmptyStatementNodeClass *Class      // Std::Elk::AST::EmptyStatementNode
var ImportStatementNodeClass *Class     // Std::Elk::AST::ImportStatementNode
var ParameterStatementNodeClass *Class  // Std::Elk::AST::ParameterStatementNode

var ProgramNodeClass *Class                       // Std::Elk::AST::ProgramNode
var InvalidNodeClass *Class                       // Std::Elk::AST::InvalidNode
var TypeExpressionNodeClass *Class                // Std::Elk::AST::TypeExpressionNode
var InstanceVariableDeclarationNodeClass *Class   // Std::Elk::AST::InstanceVariableDeclarationNode
var VariablePatternDeclarationNodeClass *Class    // Std::Elk::AST::VariablePatternDeclarationNode
var VariableDeclarationNodeClass *Class           // Std::Elk::AST::VariableDeclarationNode
var ValuePatternDeclarationNodeClass *Class       // Std::Elk::AST::ValuePatternDeclarationNode
var ValueDeclarationNodeClass *Class              // Std::Elk::AST::ValueDeclarationNode
var PostfixExpressionNodeClass *Class             // Std::Elk::AST::PostfixExpressionNode
var ModifierNodeClass *Class                      // Std::Elk::AST::ModifierNode
var ModifierIfElseNodeClass *Class                // Std::Elk::AST::ModifierIfElseNode
var ModifierForInNodeClass *Class                 // Std::Elk::AST::ModifierForInNode
var AssignmentExpressionNodeClass *Class          // Std::Elk::AST::AssignmentExpressionNode
var BinaryExpressionNodeClass *Class              // Std::Elk::AST::BinaryExpressionNode
var LogicalExpressionNodeClass *Class             // Std::Elk::AST::LogicalExpressionNode
var UnaryExpressionNodeClass *Class               // Std::Elk::AST::UnaryExpressionNode
var TrueLiteralNodeClass *Class                   // Std::Elk::AST::TrueLiteralNode
var FalseLiteralNodeClass *Class                  // Std::Elk::AST::FalseLiteralNode
var NilLiteralNodeClass *Class                    // Std::Elk::AST::NilLiteralNode
var UndefinedLiteralNodeClass *Class              // Std::Elk::AST::UndefinedLiteralNode
var SelfLiteralNodeClass *Class                   // Std::Elk::AST::SelfLiteralNode
var PublicInstanceVariableNodeClass *Class        // Std::Elk::AST::PublicInstanceVariableNode
var SimpleSymbolLiteralNodeClass *Class           // Std::Elk::AST::SimpleSymbolLiteralNode
var IntLiteralNodeClass *Class                    // Std::Elk::AST::IntLiteralNode
var Int64LiteralNodeClass *Class                  // Std::Elk::AST::Int64LiteralNode
var Int32LiteralNodeClass *Class                  // Std::Elk::AST::Int32LiteralNode
var Int16LiteralNodeClass *Class                  // Std::Elk::AST::Int16LiteralNode
var Int8LiteralNodeClass *Class                   // Std::Elk::AST::Int8LiteralNode
var UInt64LiteralNodeClass *Class                 // Std::Elk::AST::UInt64LiteralNode
var UInt32LiteralNodeClass *Class                 // Std::Elk::AST::UInt32LiteralNode
var UInt16LiteralNodeClass *Class                 // Std::Elk::AST::UInt16LiteralNode
var UInt8LiteralNodeClass *Class                  // Std::Elk::AST::UInt8LiteralNode
var FloatLiteralNodeClass *Class                  // Std::Elk::AST::FloatLiteralNode
var BigFloatLiteralNodeClass *Class               // Std::Elk::AST::BigFloatLiteralNode
var Float32LiteralNodeClass *Class                // Std::Elk::AST::Float32LiteralNode
var Float64LiteralNodeClass *Class                // Std::Elk::AST::Float64LiteralNode
var UninterpolatedRegexLiteralNodeClass *Class    // Std::Elk::AST::UninterpolatedRegexLiteralNode
var RegexLiteralContentSectionNodeClass *Class    // Std::Elk::AST::RegexLiteralContentSection
var RegexInterpolationNodeClass *Class            // Std::Elk::AST::RegexInterpolationNode
var InterpolatedRegexLiteralNodeClass *Class      // Std::Elk::AST::InterpolatedRegexLiteralNode
var CharLiteralNodeClass *Class                   // Std::Elk::AST::CharLiteralNode
var RawCharLiteralNodeClass *Class                // Std::Elk::AST::RawCharLiteralNode
var RawStringLiteralNodeClass *Class              // Std::Elk::AST::RawStringLiteralNode
var StringLiteralContentSectionNodeClass *Class   // Std::Elk::AST::StringLiteralContentSectionNode
var StringInspectInterpolationNodeClass *Class    // Std::Elk::AST::StringInspectInterpolationNode
var StringInterpolationNodeClass *Class           // Std::Elk::AST::StringInterpolationNode
var InterpolatedStringLiteralNodeClass *Class     // Std::Elk::AST::InterpolatedStringLiteralNode
var DoubleQuotedStringLiteralNodeClass *Class     // Std::Elk::AST::DoubleQuotedStringLiteralNode
var InterpolatedSymbolLiteralNodeClass *Class     // Std::Elk::AST::InterpolatedSymbolLiteralNode
var ConstantAsNodeClass *Class                    // Std::Elk::AST::ConstantAsNode
var MethodLookupAsNodeClass *Class                // Std::Elk::AST::MethodLookupAsNode
var PublicIdentifierNodeClass *Class              // Std::Elk::AST::PublicIdentifierNode
var PublicIdentifierAsNodeClass *Class            // Std::Elk::AST::PublicIdentifierAsNode
var PrivateIdentifierNodeClass *Class             // Std::Elk::AST::PrivateIdentifierNode
var PublicConstantAsNodeClass *Class              // Std::Elk::AST::PublicConstantAsNode
var PublicConstantNodeClass *Class                // Std::Elk::AST::PublicConstantNode
var PrivateConstantNodeClass *Class               // Std::Elk::AST::PrivateConstantNode
var AsExpressionNodeClass *Class                  // Std::Elk::AST::AsExpressionNode
var DoExpressionNodeClass *Class                  // Std::Elk::AST::DoExpressionNode
var MacroBoundaryNodeClass *Class                 // Std::Elk::AST::MacroBoundaryNode
var QuoteExpressionNodeClass *Class               // Std::Elk::AST::QuoteExpressionNode
var UnquoteNodeClass *Class                       // Std::Elk::AST::UnquoteNode
var SingletonBlockExpressionNodeClass *Class      // Std::Elk::AST::SingletonBlockExpressionNode
var SwitchExpressionNodeClass *Class              // Std::Elk::AST::SwitchExpressionNode
var IfExpressionNodeClass *Class                  // Std::Elk::AST::IfExpressionNode
var UnlessExpressionNodeClass *Class              // Std::Elk::AST::UnlessExpressionNode
var WhileExpressionNodeClass *Class               // Std::Elk::AST::WhileExpressionNode
var UntilExpressionNodeClass *Class               // Std::Elk::AST::UntilExpressionNode
var LoopExpressionNodeClass *Class                // Std::Elk::AST::LoopExpressionNode
var NumericForExpressionNodeClass *Class          // Std::Elk::AST::NumericForExpressionNode
var ForInExpressionNodeClass *Class               // Std::Elk::AST::ForInExpressionNode
var BreakExpressionNodeClass *Class               // Std::Elk::AST::BreakExpressionNode
var LabeledExpressionNodeClass *Class             // Std::Elk::AST::LabeledExpressionNode
var GoExpressionNodeClass *Class                  // Std::Elk::AST::GoExpressionNode
var ReturnExpressionNodeClass *Class              // Std::Elk::AST::ReturnExpressionNode
var YieldExpressionNodeClass *Class               // Std::Elk::AST::YieldExpressionNode
var ContinueExpressionNodeClass *Class            // Std::Elk::AST::ContinueExpressionNode
var ThrowExpressionNodeClass *Class               // Std::Elk::AST::ThrowExpressionNode
var MustExpressionNodeClass *Class                // Std::Elk::AST::MustExpressionNode
var TryExpressionNodeClass *Class                 // Std::Elk::AST::TryExpressionNode
var AwaitExpressionNodeClass *Class               // Std::Elk::AST::AwaitExpressionNode
var TypeofExpressionNodeClass *Class              // Std::Elk::AST::TypeofExpressionNode
var ConstantLookupNodeClass *Class                // Std::Elk::AST::ConstantLookupNode
var MethodLookupNodeClass *Class                  // Std::Elk::AST::MethodLookupNode
var InstanceMethodLookupNodeClass *Class          // Std::Elk::AST::InstanceMethodLookupNode
var UsingEntryWithSubentriesNodeClass *Class      // Std::Elk::AST::UsingEntryWithSubentriesNode
var UsingAllEntryNodeClass *Class                 // Std::Elk::AST::UsingAllEntryNode
var ClosureLiteralNodeClass *Class                // Std::Elk::AST::ClosureLiteralNode
var ClassDeclarationNodeClass *Class              // Std::Elk::AST::ClassDeclarationNode
var ModuleDeclarationNodeClass *Class             // Std::Elk::AST::ModuleDeclarationNode
var MixinDeclarationNodeClass *Class              // Std::Elk::AST::MixinDeclarationNode
var InterfaceDeclarationNodeClass *Class          // Std::Elk::AST::InterfaceDeclarationNode
var StructDeclarationNodeClass *Class             // Std::Elk::AST::StructDeclarationNode
var MacroDefinitionNodeClass *Class               // Std::Elk::AST::MacroDefinitionNode
var MethodDefinitionNodeClass *Class              // Std::Elk::AST::MethodDefinitionNode
var InitDefinitionNodeClass *Class                // Std::Elk::AST::InitDefinitionNode
var MethodSignatureDefinitionNodeClass *Class     // Std::Elk::AST::MethodSignatureDefinitionNode
var ConstantDeclarationNodeClass *Class           // Std::Elk::AST::ConstantDeclarationNode
var GenericConstantNodeClass *Class               // Std::Elk::AST::GenericConstantNode
var GenericTypeDefinitionNodeClass *Class         // Std::Elk::AST::GenericTypeDefinitionNode
var TypeDefinitionNodeClass *Class                // Std::Elk::AST::TypeDefinitionNode
var AliasDeclarationEntryClass *Class             // Std::Elk::AST::AliasDeclarationEntry
var AliasDeclarationNodeClass *Class              // Std::Elk::AST::AliasDeclarationNode
var GetterDeclarationNodeClass *Class             // Std::Elk::AST::GetterDeclarationNode
var SetterDeclarationNodeClass *Class             // Std::Elk::AST::SetterDeclarationNode
var AttrDeclarationNodeClass *Class               // Std::Elk::AST::AttrDeclarationNode
var UsingExpressionNodeClass *Class               // Std::Elk::AST::UsingExpressionNode
var IncludeExpressionNodeClass *Class             // Std::Elk::AST::IncludeExpressionNode
var ExtendWhereBlockExpressionNodeClass *Class    // Std::Elk::AST::ExtendWhereBlockExpressionNode
var ImplementExpressionNodeClass *Class           // Std::Elk::AST::ImplementExpressionNode
var NewExpressionNodeClass *Class                 // Std::Elk::AST::NewExpressionNode
var GenericConstructorCallNodeClass *Class        // Std::Elk::AST::GenericConstructorCallNode
var ConstructorCallNodeClass *Class               // Std::Elk::AST::ConstructorCallNode
var NamedCallArgumentNodeClass *Class             // Std::Elk::AST::NamedCallArgumentNode
var DoubleSplatExpressionNodeClass *Class         // Std::Elk::AST::DoubleSplatExpressionNode
var AttributeAccessNodeClass *Class               // Std::Elk::AST::AttributeAccessNode
var SubscriptExpressionNodeClass *Class           // Std::Elk::AST::SubscriptExpressionNode
var NilSafeSubscriptExpressionNodeClass *Class    // Std::Elk::AST::NilSafeSubscriptExpressionNode
var CallNodeClass *Class                          // Std::Elk::AST::CallNode
var GenericMethodCallNodeClass *Class             // Std::Elk::AST::GenericMethodCallNode
var MethodCallNodeClass *Class                    // Std::Elk::AST::MethodCallNode
var ReceiverlessMethodCallNodeClass *Class        // Std::Elk::AST::ReceiverlessMethodCallNode
var GenericReceiverlessMethodCallNodeClass *Class // Std::Elk::AST::GenericReceiverlessMethodCallNode
var ScopedMacroCallNodeClass *Class               // Std::Elk::AST::ScopedMacroCallNode
var MacroCallNodeClass *Class                     // Std::Elk::AST::MacroCallNode
var ReceiverlessMacroCallNodeClass *Class         // Std::Elk::AST::ReceiverlessMacroCallNode
var SplatExpressionNodeClass *Class               // Std::Elk::AST::SplatExpressionNode
var KeyValueExpressionNodeClass *Class            // Std::Elk::AST::KeyValueExpressionNode
var SymbolKeyValueExpressionNodeClass *Class      // Std::Elk::AST::SymbolKeyValueExpressionNode
var WordArrayListLiteralNodeClass *Class          // Std::Elk::AST::WordArrayListLiteralNode
var ArrayListLiteralNodeClass *Class              // Std::Elk::AST::ArrayListLiteralNode
var SymbolArrayListLiteralNodeClass *Class        // Std::Elk::AST::SymbolArrayListLiteralNode
var HexArrayListLiteralNodeClass *Class           // Std::Elk::AST::HexArrayListLiteralNode
var BinArrayListLiteralNodeClass *Class           // Std::Elk::AST::BinArrayListLiteralNode
var ArrayTupleLiteralNodeClass *Class             // Std::Elk::AST::ArrayTupleLiteralNode
var WordArrayTupleLiteralNodeClass *Class         // Std::Elk::AST::WordArrayTupleLiteralNode
var SymbolArrayTupleLiteralNodeClass *Class       // Std::Elk::AST::SymbolArrayTupleLiteralNode
var HexArrayTupleLiteralNodeClass *Class          // Std::Elk::AST::HexArrayTupleLiteralNode
var BinArrayTupleLiteralNodeClass *Class          // Std::Elk::AST::BinArrayTupleLiteralNode
var HashSetLiteralNodeClass *Class                // Std::Elk::AST::HashSetLiteralNode
var WordHashSetLiteralNodeClass *Class            // Std::Elk::AST::WordHashSetLiteralNode
var SymbolHashSetLiteralNodeClass *Class          // Std::Elk::AST::SymbolHashSetLiteralNode
var HexHashSetLiteralNodeClass *Class             // Std::Elk::AST::HexHashSetLiteralNode
var BinHashSetLiteralNodeClass *Class             // Std::Elk::AST::BinHashSetLiteralNode
var HashMapLiteralNodeClass *Class                // Std::Elk::AST::HashMapLiteralNode
var HashRecordLiteralNodeClass *Class             // Std::Elk::AST::HashRecordLiteralNode
var RangeLiteralNodeClass *Class                  // Std::Elk::AST::RangeLiteralNode
var CaseNodeClass *Class                          // Std::Elk::AST::CaseNode
var CatchNodeClass *Class                         // Std::Elk::AST::CatchNode
var VariantTypeParameterNodeClass *Class          // Std::Elk::AST::VariantTypeParameterNode
var FormalParameterNodeClass *Class               // Std::Elk::AST::FormalParameterNode
var MethodParameterNodeClass *Class               // Std::Elk::AST::MethodParameterNode
var SignatureParameterNodeClass *Class            // Std::Elk::AST::SignatureParameterNode
var AttributeParameterNodeClass *Class            // Std::Elk::AST::AttributeParameterNode
var BoolLiteralNodeClass *Class                   // Std::Elk::AST::BoolLiteralNode
var VoidTypeNodeClass *Class                      // Std::Elk::AST::VoidTypeNode
var NeverTypeNodeClass *Class                     // Std::Elk::AST::NeverTypeNode
var AnyTypeNodeClass *Class                       // Std::Elk::AST::AnyTypeNode
var UnionTypeNodeClass *Class                     // Std::Elk::AST::UnionTypeNode
var IntersectionTypeNodeClass *Class              // Std::Elk::AST::IntersectionTypeNode
var BinaryTypeNodeClass *Class                    // Std::Elk::AST::BinaryTypeNode
var NilableTypeNodeClass *Class                   // Std::Elk::AST::NilableTypeNode
var InstanceOfTypeNodeClass *Class                // Std::Elk::AST::InstanceOfTypeNode
var SingletonTypeNodeClass *Class                 // Std::Elk::AST::SingletonTypeNode
var NotTypeNodeClass *Class                       // Std::Elk::AST::NotTypeNode
var ClosureTypeNodeClass *Class                   // Std::Elk::AST::ClosureTypeNode
var UnaryTypeNodeClass *Class                     // Std::Elk::AST::UnaryTypeNode
var AsPatternNodeClass *Class                     // Std::Elk::AST::AsPatternNode
var SymbolKeyValuePatternNodeClass *Class         // Std::Elk::AST::SymbolKeyValuePatternNode
var KeyValuePatternNodeClass *Class               // Std::Elk::AST::KeyValuePatternNode
var ObjectPatternNodeClass *Class                 // Std::Elk::AST::ObjectPatternNode
var RecordPatternNodeClass *Class                 // Std::Elk::AST::RecordPatternNode
var MapPatternNodeClass *Class                    // Std::Elk::AST::MapPatternNode
var RestPatternNodeClass *Class                   // Std::Elk::AST::RestPatternNode
var SetPatternNodeClass *Class                    // Std::Elk::AST::SetPatternNode
var ListPatternNodeClass *Class                   // Std::Elk::AST::ListPatternNode
var TuplePatternNodeClass *Class                  // Std::Elk::AST::TuplePatternNode
var BinaryPatternNodeClass *Class                 // Std::Elk::AST::BinaryPatternNode

func initElkAST() {
	ElkASTModule = NewModule()
	ElkModule.AddConstantString("AST", Ref(ElkASTModule))

	NodeMixin = NewMixin()
	ElkASTModule.AddConstantString("Node", Ref(NodeMixin))

	NodeMixin.AddConstantString("Convertible", Ref(NewInterface()))

	NodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(FormatErrorClass))
	NodeMixin.AddConstantString("FormatError", Ref(NodeFormatErrorClass))

	StatementNodeMixin = NewMixin()
	StatementNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("StatementNode", Ref(StatementNodeMixin))

	StructBodyStatementNodeMixin = NewMixin()
	StructBodyStatementNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("StructBodyStatementNode", Ref(StructBodyStatementNodeMixin))

	ParameterNodeMixin = NewMixin()
	ParameterNodeMixin.IncludeMixin(NodeMixin)
	ParameterNodeMixin.AddConstantString("NORMAL_KIND", UInt8(0).ToValue())
	ParameterNodeMixin.AddConstantString("POSITIONAL_REST_KIND", UInt8(1).ToValue())
	ParameterNodeMixin.AddConstantString("NAMED_REST_KIND", UInt8(2).ToValue())
	ElkASTModule.AddConstantString("ParameterNode", Ref(ParameterNodeMixin))

	PatternNodeMixin = NewMixin()
	PatternNodeMixin.IncludeMixin(NodeMixin)
	PatternNodeMixin.AddConstantString("Convertible", Ref(NewInterface()))
	ElkASTModule.AddConstantString("PatternNode", Ref(PatternNodeMixin))

	TypeNodeMixin = NewMixin()
	TypeNodeMixin.IncludeMixin(NodeMixin)
	TypeNodeMixin.AddConstantString("Convertible", Ref(NewInterface()))
	ElkASTModule.AddConstantString("TypeNode", Ref(TypeNodeMixin))

	ExpressionStatementNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ExpressionStatementNodeClass.IncludeMixin(StatementNodeMixin)
	ElkASTModule.AddConstantString("ExpressionStatementNode", Ref(ExpressionStatementNodeClass))

	EmptyStatementNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	EmptyStatementNodeClass.IncludeMixin(StatementNodeMixin)
	EmptyStatementNodeClass.IncludeMixin(StructBodyStatementNodeMixin)
	ElkASTModule.AddConstantString("EmptyStatementNode", Ref(EmptyStatementNodeClass))

	ImportStatementNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ImportStatementNodeClass.IncludeMixin(StatementNodeMixin)
	ElkASTModule.AddConstantString("ImportStatementNode", Ref(ImportStatementNodeClass))

	ParameterStatementNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ParameterStatementNodeClass.IncludeMixin(StructBodyStatementNodeMixin)
	ElkASTModule.AddConstantString("ParameterStatementNode", Ref(ParameterStatementNodeClass))

	ProgramNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ProgramNodeClass.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("ProgramNode", Ref(ProgramNodeClass))

	ExpressionNodeMixin = NewMixin()
	ExpressionNodeMixin.IncludeMixin(NodeMixin)
	ExpressionNodeMixin.AddConstantString("Convertible", Ref(NewInterface()))
	ElkASTModule.AddConstantString("ExpressionNode", Ref(ExpressionNodeMixin))

	PatternExpressionNodeMixin = NewMixin()
	PatternExpressionNodeMixin.IncludeMixin(ExpressionNodeMixin)
	PatternExpressionNodeMixin.IncludeMixin(PatternNodeMixin)
	PatternExpressionNodeMixin.AddConstantString("Convertible", Ref(NewInterface()))
	ElkASTModule.AddConstantString("PatternExpressionNode", Ref(PatternExpressionNodeMixin))

	SymbolLiteralNodeMixin = NewMixin()
	SymbolLiteralNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SymbolLiteralNode", Ref(SymbolLiteralNodeMixin))

	UsingEntryNodeMixin = NewMixin()
	UsingEntryNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UsingEntryNode", Ref(UsingEntryNodeMixin))

	UsingSubentryNodeMixin = NewMixin()
	UsingSubentryNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UsingSubentryNode", Ref(UsingSubentryNodeMixin))

	NamedArgumentNodeMixin = NewMixin()
	NamedArgumentNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("NamedArgumentNode", Ref(NamedArgumentNodeMixin))

	WordCollectionContentNodeMixin = NewMixin()
	WordCollectionContentNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("WordCollectionContentNode", Ref(WordCollectionContentNodeMixin))

	SymbolCollectionContentNodeMixin = NewMixin()
	SymbolCollectionContentNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SymbolCollectionContentNode", Ref(SymbolCollectionContentNodeMixin))

	IntCollectionContentNodeMixin = NewMixin()
	IntCollectionContentNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("IntCollectionContentNode", Ref(IntCollectionContentNodeMixin))

	IdentifierNodeMixin = NewMixin()
	IdentifierNodeMixin.IncludeMixin(PatternExpressionNodeMixin)
	IdentifierNodeMixin.AddConstantString("Convertible", Ref(NewInterface()))
	ElkASTModule.AddConstantString("IdentifierNode", Ref(IdentifierNodeMixin))

	IdentifierNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(NodeFormatErrorClass))
	IdentifierNodeMixin.AddConstantString("FormatError", Ref(IdentifierNodeFormatErrorClass))

	RegexLiteralContentNodeMixin = NewMixin()
	RegexLiteralContentNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("RegexLiteralContentNode", Ref(RegexLiteralContentNodeMixin))

	RegexLiteralNodeMixin = NewMixin()
	RegexLiteralNodeMixin.IncludeMixin(PatternExpressionNodeMixin)
	RegexLiteralNodeMixin.AddConstantString("CASE_INSENSITIVE_FLAG", UInt8(flag.CaseInsensitiveFlag).ToValue())
	RegexLiteralNodeMixin.AddConstantString("MULTILINE_FLAG", UInt8(flag.MultilineFlag).ToValue())
	RegexLiteralNodeMixin.AddConstantString("DOT_ALL_FLAG", UInt8(flag.DotAllFlag).ToValue())
	RegexLiteralNodeMixin.AddConstantString("UNGREEDY_FLAG", UInt8(flag.UngreedyFlag).ToValue())
	RegexLiteralNodeMixin.AddConstantString("EXTENDED_FLAG", UInt8(flag.ExtendedFlag).ToValue())
	RegexLiteralNodeMixin.AddConstantString("ASCII_FLAG", UInt8(flag.ASCIIFlag).ToValue())
	ElkASTModule.AddConstantString("RegexLiteralNode", Ref(RegexLiteralNodeMixin))

	StringLiteralContentNodeMixin = NewMixin()
	StringLiteralContentNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("StringLiteralContentNode", Ref(StringLiteralContentNodeMixin))

	StringLiteralNodeMixin = NewMixin()
	StringLiteralNodeMixin.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("StringLiteralNode", Ref(StringLiteralNodeMixin))

	SimpleStringLiteralNodeMixin = NewMixin()
	SimpleStringLiteralNodeMixin.IncludeMixin(StringLiteralNodeMixin)
	SimpleStringLiteralNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SimpleStringLiteralNode", Ref(SimpleStringLiteralNodeMixin))

	TypeParameterNodeMixin = NewMixin()
	TypeParameterNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("TypeParameterNode", Ref(TypeParameterNodeMixin))

	InstanceVariableNodeMixin = NewMixin()
	InstanceVariableNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InstanceVariableNode", Ref(InstanceVariableNodeMixin))

	ComplexConstantNodeMixin = NewMixin()
	ComplexConstantNodeMixin.IncludeMixin(PatternExpressionNodeMixin)
	ComplexConstantNodeMixin.IncludeMixin(TypeNodeMixin)
	ComplexConstantNodeMixin.IncludeMixin(UsingEntryNodeMixin)
	ElkASTModule.AddConstantString("ComplexConstantNode", Ref(ComplexConstantNodeMixin))

	ConstantNodeMixin = NewMixin()
	ConstantNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ConstantNodeMixin.IncludeMixin(TypeNodeMixin)
	ConstantNodeMixin.IncludeMixin(ComplexConstantNodeMixin)
	ConstantNodeMixin.IncludeMixin(UsingEntryNodeMixin)
	ConstantNodeMixin.AddConstantString("Convertible", Ref(NewInterface()))
	ElkASTModule.AddConstantString("ConstantNode", Ref(ConstantNodeMixin))

	ConstantNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(FormatErrorClass))
	ConstantNodeMixin.AddConstantString("FormatError", Ref(ConstantNodeFormatErrorClass))

	InvalidNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InvalidNodeClass.IncludeMixin(StatementNodeMixin)
	InvalidNodeClass.IncludeMixin(StructBodyStatementNodeMixin)
	InvalidNodeClass.IncludeMixin(ParameterNodeMixin)
	InvalidNodeClass.IncludeMixin(TypeNodeMixin)
	InvalidNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	InvalidNodeClass.IncludeMixin(SymbolLiteralNodeMixin)
	InvalidNodeClass.IncludeMixin(UsingEntryNodeMixin)
	InvalidNodeClass.IncludeMixin(UsingSubentryNodeMixin)
	InvalidNodeClass.IncludeMixin(ComplexConstantNodeMixin)
	InvalidNodeClass.IncludeMixin(ConstantNodeMixin)
	InvalidNodeClass.IncludeMixin(NamedArgumentNodeMixin)
	InvalidNodeClass.IncludeMixin(WordCollectionContentNodeMixin)
	InvalidNodeClass.IncludeMixin(SymbolCollectionContentNodeMixin)
	InvalidNodeClass.IncludeMixin(IntCollectionContentNodeMixin)
	InvalidNodeClass.IncludeMixin(IdentifierNodeMixin)
	InvalidNodeClass.IncludeMixin(RegexLiteralContentNodeMixin)
	InvalidNodeClass.IncludeMixin(RegexLiteralNodeMixin)
	InvalidNodeClass.IncludeMixin(StringLiteralContentNodeMixin)
	InvalidNodeClass.IncludeMixin(StringLiteralNodeMixin)
	InvalidNodeClass.IncludeMixin(SimpleStringLiteralNodeMixin)
	InvalidNodeClass.IncludeMixin(TypeParameterNodeMixin)
	ElkASTModule.AddConstantString("InvalidNode", Ref(InvalidNodeClass))

	TypeExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	TypeExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("TypeExpressionNode", Ref(TypeExpressionNodeClass))

	InstanceVariableDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InstanceVariableDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InstanceVariableDeclarationNode", Ref(InstanceVariableDeclarationNodeClass))

	VariablePatternDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	VariablePatternDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("VariablePatternDeclarationNode", Ref(VariablePatternDeclarationNodeClass))

	VariableDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	VariableDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("VariableDeclarationNode", Ref(VariableDeclarationNodeClass))

	ValuePatternDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ValuePatternDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ValuePatternDeclarationNode", Ref(ValuePatternDeclarationNodeClass))

	ValueDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ValueDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ValueDeclarationNode", Ref(ValueDeclarationNodeClass))

	PostfixExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PostfixExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("PostfixExpressionNode", Ref(PostfixExpressionNodeClass))

	ModifierNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ModifierNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ModifierNode", Ref(ModifierNodeClass))

	ModifierIfElseNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ModifierIfElseNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ModifierIfElseNode", Ref(ModifierIfElseNodeClass))

	ModifierForInNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ModifierForInNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ModifierForInNode", Ref(ModifierForInNodeClass))

	AssignmentExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	AssignmentExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AssignmentExpressionNode", Ref(AssignmentExpressionNodeClass))

	BinaryExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BinaryExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("BinaryExpressionNode", Ref(BinaryExpressionNodeClass))

	LogicalExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	LogicalExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("LogicalExpressionNode", Ref(LogicalExpressionNodeClass))

	UnaryExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UnaryExpressionNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("UnaryExpressionNode", Ref(UnaryExpressionNodeClass))

	TrueLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	TrueLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	TrueLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("TrueLiteralNode", Ref(TrueLiteralNodeClass))

	FalseLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	FalseLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	FalseLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("FalseLiteralNode", Ref(FalseLiteralNodeClass))

	NilLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	NilLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	NilLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	NilLiteralNodeClass.IncludeMixin(UsingEntryNodeMixin)
	NilLiteralNodeClass.IncludeMixin(ComplexConstantNodeMixin)
	ElkASTModule.AddConstantString("NilLiteralNode", Ref(NilLiteralNodeClass))

	UndefinedLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UndefinedLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UndefinedLiteralNode", Ref(UndefinedLiteralNodeClass))

	SelfLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SelfLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	SelfLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("SelfLiteralNode", Ref(SelfLiteralNodeClass))

	PublicInstanceVariableNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PublicInstanceVariableNodeClass.IncludeMixin(ExpressionNodeMixin)
	PublicInstanceVariableNodeClass.IncludeMixin(InstanceVariableNodeMixin)
	ElkASTModule.AddConstantString("PublicInstanceVariableNode", Ref(PublicInstanceVariableNodeClass))

	SimpleSymbolLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SimpleSymbolLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	SimpleSymbolLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	SimpleSymbolLiteralNodeClass.IncludeMixin(SymbolLiteralNodeMixin)
	SimpleSymbolLiteralNodeClass.IncludeMixin(SymbolCollectionContentNodeMixin)
	ElkASTModule.AddConstantString("SimpleSymbolLiteralNode", Ref(SimpleSymbolLiteralNodeClass))

	IntLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	IntLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	IntLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	IntLiteralNodeClass.IncludeMixin(IntCollectionContentNodeMixin)
	ElkASTModule.AddConstantString("IntLiteralNode", Ref(IntLiteralNodeClass))

	IntLiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	IntLiteralNodeClass.AddConstantString("FormatError", Ref(IntLiteralNodeFormatErrorClass))

	Int64LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Int64LiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	Int64LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("Int64LiteralNode", Ref(Int64LiteralNodeClass))

	Int64LiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	Int64LiteralNodeClass.AddConstantString("FormatError", Ref(Int64LiteralNodeFormatErrorClass))

	Int32LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Int32LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Int32LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	Int32LiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("Int32LiteralNode", Ref(Int32LiteralNodeClass))

	Int32LiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	Int32LiteralNodeClass.AddConstantString("FormatError", Ref(Int32LiteralNodeFormatErrorClass))

	Int16LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Int16LiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	Int16LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("Int16LiteralNode", Ref(Int16LiteralNodeClass))

	Int16LiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	Int16LiteralNodeClass.AddConstantString("FormatError", Ref(Int16LiteralNodeFormatErrorClass))

	Int8LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Int8LiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	Int8LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("Int8LiteralNode", Ref(Int8LiteralNodeClass))

	Int8LiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	Int8LiteralNodeClass.AddConstantString("FormatError", Ref(Int8LiteralNodeFormatErrorClass))

	UInt64LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UInt64LiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	UInt64LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UInt64LiteralNode", Ref(UInt64LiteralNodeClass))

	UInt64LiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	UInt64LiteralNodeClass.AddConstantString("FormatError", Ref(UInt64LiteralNodeFormatErrorClass))

	UInt32LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UInt32LiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	UInt32LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UInt32LiteralNode", Ref(UInt32LiteralNodeClass))

	UInt32LiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	UInt32LiteralNodeClass.AddConstantString("FormatError", Ref(UInt32LiteralNodeFormatErrorClass))

	UInt16LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UInt16LiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	UInt16LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UInt16LiteralNode", Ref(UInt16LiteralNodeClass))

	UInt16LiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	UInt16LiteralNodeClass.AddConstantString("FormatError", Ref(UInt16LiteralNodeFormatErrorClass))

	UInt8LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UInt8LiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	UInt8LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UInt8LiteralNode", Ref(UInt8LiteralNodeClass))

	UInt8LiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	UInt8LiteralNodeClass.AddConstantString("FormatError", Ref(UInt8LiteralNodeFormatErrorClass))

	FloatLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	FloatLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	FloatLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("FloatLiteralNode", Ref(FloatLiteralNodeClass))

	FloatLiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	FloatLiteralNodeClass.AddConstantString("FormatError", Ref(FloatLiteralNodeFormatErrorClass))

	BigFloatLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BigFloatLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	BigFloatLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("BigFloatLiteralNode", Ref(BigFloatLiteralNodeClass))

	BigFloatLiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	BigFloatLiteralNodeClass.AddConstantString("FormatError", Ref(BigFloatLiteralNodeFormatErrorClass))

	Float64LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Float64LiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	Float64LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("Float64LiteralNode", Ref(Float64LiteralNodeClass))

	Float64LiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	Float64LiteralNodeClass.AddConstantString("FormatError", Ref(Float64LiteralNodeFormatErrorClass))

	Float32LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Float32LiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	Float32LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("Float32LiteralNode", Ref(Float32LiteralNodeClass))

	Float32LiteralNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	Float32LiteralNodeClass.AddConstantString("FormatError", Ref(Float32LiteralNodeFormatErrorClass))

	UninterpolatedRegexLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UninterpolatedRegexLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	UninterpolatedRegexLiteralNodeClass.IncludeMixin(RegexLiteralNodeMixin)
	ElkASTModule.AddConstantString("UninterpolatedRegexLiteralNode", Ref(UninterpolatedRegexLiteralNodeClass))

	RegexLiteralContentSectionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	RegexLiteralContentSectionNodeClass.IncludeMixin(NodeMixin)
	RegexLiteralContentSectionNodeClass.IncludeMixin(RegexLiteralContentNodeMixin)
	ElkASTModule.AddConstantString("RegexLiteralContentSectionNode", Ref(RegexLiteralContentSectionNodeClass))

	RegexInterpolationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	RegexInterpolationNodeClass.IncludeMixin(NodeMixin)
	RegexInterpolationNodeClass.IncludeMixin(RegexLiteralContentNodeMixin)
	ElkASTModule.AddConstantString("RegexInterpolationNodeClass", Ref(RegexInterpolationNodeClass))

	InterpolatedRegexLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InterpolatedRegexLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	InterpolatedRegexLiteralNodeClass.IncludeMixin(RegexLiteralNodeMixin)
	ElkASTModule.AddConstantString("InterpolatedRegexLiteralNode", Ref(InterpolatedRegexLiteralNodeClass))

	CharLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	CharLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	CharLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("CharLiteralNode", Ref(CharLiteralNodeClass))

	RawCharLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	RawCharLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	RawCharLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("RawCharLiteralNode", Ref(RawCharLiteralNodeClass))

	RawStringLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	RawStringLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	RawStringLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	RawStringLiteralNodeClass.IncludeMixin(WordCollectionContentNodeMixin)
	RawStringLiteralNodeClass.IncludeMixin(StringLiteralNodeMixin)
	RawStringLiteralNodeClass.IncludeMixin(SimpleStringLiteralNodeMixin)
	ElkASTModule.AddConstantString("RawStringLiteralNode", Ref(RawStringLiteralNodeClass))

	StringLiteralContentSectionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StringLiteralContentSectionNodeClass.IncludeMixin(NodeMixin)
	StringLiteralContentSectionNodeClass.IncludeMixin(StringLiteralContentNodeMixin)
	ElkASTModule.AddConstantString("StringLiteralContentSectionNode", Ref(StringLiteralContentSectionNodeClass))

	StringInspectInterpolationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StringInspectInterpolationNodeClass.IncludeMixin(NodeMixin)
	StringInspectInterpolationNodeClass.IncludeMixin(StringLiteralContentNodeMixin)
	ElkASTModule.AddConstantString("StringInspectInterpolationNode", Ref(StringInspectInterpolationNodeClass))

	StringInterpolationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StringInterpolationNodeClass.IncludeMixin(NodeMixin)
	StringInterpolationNodeClass.IncludeMixin(StringLiteralContentNodeMixin)
	ElkASTModule.AddConstantString("StringInterpolationNode", Ref(StringInterpolationNodeClass))

	InterpolatedStringLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InterpolatedStringLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	InterpolatedStringLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	InterpolatedStringLiteralNodeClass.IncludeMixin(StringLiteralNodeMixin)
	ElkASTModule.AddConstantString("InterpolatedStringLiteralNode", Ref(InterpolatedStringLiteralNodeClass))

	DoubleQuotedStringLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	DoubleQuotedStringLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	DoubleQuotedStringLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	DoubleQuotedStringLiteralNodeClass.IncludeMixin(StringLiteralNodeMixin)
	DoubleQuotedStringLiteralNodeClass.IncludeMixin(SimpleStringLiteralNodeMixin)
	ElkASTModule.AddConstantString("DoubleQuotedStringLiteralNode", Ref(DoubleQuotedStringLiteralNodeClass))

	InterpolatedSymbolLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InterpolatedSymbolLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	InterpolatedSymbolLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	InterpolatedSymbolLiteralNodeClass.IncludeMixin(SymbolLiteralNodeMixin)
	ElkASTModule.AddConstantString("InterpolatedSymbolLiteralNode", Ref(InterpolatedSymbolLiteralNodeClass))

	ConstantAsNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ConstantAsNodeClass.IncludeMixin(ExpressionNodeMixin)
	ConstantAsNodeClass.IncludeMixin(UsingEntryNodeMixin)
	ElkASTModule.AddConstantString("ConstantAsNode", Ref(ConstantAsNodeClass))

	MethodLookupAsNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MethodLookupAsNodeClass.IncludeMixin(ExpressionNodeMixin)
	MethodLookupAsNodeClass.IncludeMixin(UsingEntryNodeMixin)
	ElkASTModule.AddConstantString("MethodLookupAsNode", Ref(MethodLookupAsNodeClass))

	PublicIdentifierNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PublicIdentifierNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	PublicIdentifierNodeClass.IncludeMixin(UsingSubentryNodeMixin)
	PublicIdentifierNodeClass.IncludeMixin(IdentifierNodeMixin)
	ElkASTModule.AddConstantString("PublicIdentifierNode", Ref(PublicIdentifierNodeClass))

	PublicIdentifierAsNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PublicIdentifierAsNodeClass.IncludeMixin(ExpressionNodeMixin)
	PublicIdentifierAsNodeClass.IncludeMixin(UsingSubentryNodeMixin)
	ElkASTModule.AddConstantString("PublicIdentifierAsNode", Ref(PublicIdentifierAsNodeClass))

	PrivateIdentifierNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PrivateIdentifierNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	PrivateIdentifierNodeClass.IncludeMixin(IdentifierNodeMixin)
	ElkASTModule.AddConstantString("PrivateIdentifierNode", Ref(PrivateIdentifierNodeClass))

	PrivateIdentifierNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(IdentifierNodeFormatErrorClass))
	PrivateIdentifierNodeClass.AddConstantString("FormatError", Ref(PrivateIdentifierNodeFormatErrorClass))

	PublicConstantNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PublicConstantNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	PublicConstantNodeClass.IncludeMixin(TypeNodeMixin)
	PublicConstantNodeClass.IncludeMixin(UsingEntryNodeMixin)
	PublicConstantNodeClass.IncludeMixin(UsingSubentryNodeMixin)
	PublicConstantNodeClass.IncludeMixin(ComplexConstantNodeMixin)
	PublicConstantNodeClass.IncludeMixin(ConstantNodeMixin)
	ElkASTModule.AddConstantString("PublicConstantNode", Ref(PublicConstantNodeClass))

	PublicConstantAsNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PublicConstantAsNodeClass.IncludeMixin(ExpressionNodeMixin)
	PublicConstantAsNodeClass.IncludeMixin(UsingSubentryNodeMixin)
	ElkASTModule.AddConstantString("PublicConstantAsNode", Ref(PublicConstantAsNodeClass))

	PrivateConstantNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PrivateConstantNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	PrivateConstantNodeClass.IncludeMixin(TypeNodeMixin)
	PrivateConstantNodeClass.IncludeMixin(UsingEntryNodeMixin)
	PrivateConstantNodeClass.IncludeMixin(ComplexConstantNodeMixin)
	PrivateConstantNodeClass.IncludeMixin(ConstantNodeMixin)
	ElkASTModule.AddConstantString("PrivateConstantNode", Ref(PrivateConstantNodeClass))

	PrivateConstantNodeFormatErrorClass = NewClassWithOptions(ClassWithSuperclass(ConstantNodeMixin))
	PrivateConstantNodeClass.AddConstantString("FormatError", Ref(PrivateConstantNodeFormatErrorClass))

	AsExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	AsExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AsExpressionNode", Ref(AsExpressionNodeClass))

	DoExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	DoExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("DoExpressionNode", Ref(DoExpressionNodeClass))

	MacroBoundaryNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MacroBoundaryNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MacroBoundaryNode", Ref(MacroBoundaryNodeClass))

	QuoteExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	QuoteExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("QuoteExpressionNode", Ref(QuoteExpressionNodeClass))

	UnquoteNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UnquoteNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	UnquoteNodeClass.IncludeMixin(TypeNodeMixin)
	UnquoteNodeClass.IncludeMixin(ConstantNodeMixin)
	UnquoteNodeClass.IncludeMixin(IdentifierNodeMixin)
	ElkASTModule.AddConstantString("UnquoteNode", Ref(UnquoteNodeClass))

	SingletonBlockExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SingletonBlockExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SingletonBlockExpressionNode", Ref(SingletonBlockExpressionNodeClass))

	SwitchExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SwitchExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SwitchExpressionNode", Ref(SwitchExpressionNodeClass))

	IfExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	IfExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("IfExpressionNode", Ref(IfExpressionNodeClass))

	UnlessExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UnlessExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UnlessExpressionNode", Ref(UnlessExpressionNodeClass))

	WhileExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	WhileExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("WhileExpressionNode", Ref(WhileExpressionNodeClass))

	UntilExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UntilExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UntilExpressionNode", Ref(UntilExpressionNodeClass))

	LoopExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	LoopExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("LoopExpressionNode", Ref(LoopExpressionNodeClass))

	ForInExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ForInExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ForInExpressionNode", Ref(ForInExpressionNodeClass))

	BreakExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BreakExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("BreakExpressionNode", Ref(BreakExpressionNodeClass))

	NumericForExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	NumericForExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("NumericForExpressionNode", Ref(NumericForExpressionNodeClass))

	LabeledExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	LabeledExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("LabeledExpressionNode", Ref(LabeledExpressionNodeClass))

	GoExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	GoExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GoExpressionNode", Ref(GoExpressionNodeClass))

	ReturnExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ReturnExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ReturnExpressionNode", Ref(ReturnExpressionNodeClass))

	YieldExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	YieldExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("YieldExpressionNode", Ref(YieldExpressionNodeClass))

	ContinueExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ContinueExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ContinueExpressionNode", Ref(ContinueExpressionNodeClass))

	ThrowExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ThrowExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ThrowExpressionNode", Ref(ThrowExpressionNodeClass))

	MustExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MustExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MustExpressionNode", Ref(MustExpressionNodeClass))

	TryExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	TryExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("TryExpressionNode", Ref(TryExpressionNodeClass))

	AwaitExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	AwaitExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AwaitExpressionNode", Ref(AwaitExpressionNodeClass))

	TypeofExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	TypeofExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("TypeofExpressionNode", Ref(TypeofExpressionNodeClass))

	ConstantLookupNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ConstantLookupNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ConstantLookupNodeClass.IncludeMixin(TypeNodeMixin)
	ConstantLookupNodeClass.IncludeMixin(UsingEntryNodeMixin)
	ConstantLookupNodeClass.IncludeMixin(ComplexConstantNodeMixin)
	ElkASTModule.AddConstantString("ConstantLookupNode", Ref(ConstantLookupNodeClass))

	MethodLookupNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MethodLookupNodeClass.IncludeMixin(ExpressionNodeMixin)
	MethodLookupNodeClass.IncludeMixin(UsingEntryNodeMixin)
	ElkASTModule.AddConstantString("MethodLookupNode", Ref(MethodLookupNodeClass))

	InstanceMethodLookupNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InstanceMethodLookupNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InstanceMethodLookupNode", Ref(InstanceMethodLookupNodeClass))

	UsingEntryWithSubentriesNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UsingEntryWithSubentriesNodeClass.IncludeMixin(ExpressionNodeMixin)
	UsingEntryWithSubentriesNodeClass.IncludeMixin(UsingEntryNodeMixin)
	ElkASTModule.AddConstantString("UsingEntryWithSubentriesNode", Ref(UsingEntryWithSubentriesNodeClass))

	UsingAllEntryNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UsingAllEntryNodeClass.IncludeMixin(ExpressionNodeMixin)
	UsingAllEntryNodeClass.IncludeMixin(UsingEntryNodeMixin)
	ElkASTModule.AddConstantString("UsingAllEntryNode", Ref(UsingAllEntryNodeClass))

	ClosureLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ClosureLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ClosureLiteralNode", Ref(ClosureLiteralNodeClass))

	ClassDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ClassDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ClassDeclarationNode", Ref(ClassDeclarationNodeClass))

	ModuleDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ModuleDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ModuleDeclarationNode", Ref(ModuleDeclarationNodeClass))

	MixinDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MixinDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MixinDeclarationNode", Ref(MixinDeclarationNodeClass))

	InterfaceDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InterfaceDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InterfaceDeclarationNode", Ref(InterfaceDeclarationNodeClass))

	StructDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StructDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("StructDeclarationNode", Ref(StructDeclarationNodeClass))

	MethodDefinitionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MethodDefinitionNodeClass.IncludeMixin(ExpressionNodeMixin)
	MethodDefinitionNodeClass.AddConstantString("ABSTRACT_FLAG", UInt8(1).ToValue())
	MethodDefinitionNodeClass.AddConstantString("SEALED_FLAG", UInt8(2).ToValue())
	MethodDefinitionNodeClass.AddConstantString("GENERATOR_FLAG", UInt8(4).ToValue())
	MethodDefinitionNodeClass.AddConstantString("ASYNC_FLAG", UInt8(8).ToValue())
	ElkASTModule.AddConstantString("MethodDefinitionNode", Ref(MethodDefinitionNodeClass))

	MacroDefinitionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MacroDefinitionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MacroDefinitionNode", Ref(MacroDefinitionNodeClass))

	InitDefinitionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InitDefinitionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InitDefinitionNode", Ref(InitDefinitionNodeClass))

	MethodSignatureDefinitionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MethodSignatureDefinitionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MethodSignatureDefinitionNode", Ref(MethodSignatureDefinitionNodeClass))

	ConstantDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ConstantDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ConstantDeclarationNode", Ref(ConstantDeclarationNodeClass))

	GenericConstantNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	GenericConstantNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	GenericConstantNodeClass.IncludeMixin(TypeNodeMixin)
	GenericConstantNodeClass.IncludeMixin(UsingEntryNodeMixin)
	GenericConstantNodeClass.IncludeMixin(ComplexConstantNodeMixin)
	ElkASTModule.AddConstantString("GenericConstantNode", Ref(GenericConstantNodeClass))

	GenericTypeDefinitionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	GenericTypeDefinitionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GenericTypeDefinitionNode", Ref(GenericTypeDefinitionNodeClass))

	TypeDefinitionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	TypeDefinitionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("TypeDefinitionNode", Ref(TypeDefinitionNodeClass))

	AliasDeclarationEntryClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	AliasDeclarationEntryClass.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("AliasDeclarationEntry", Ref(AliasDeclarationEntryClass))

	AliasDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	AliasDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AliasDeclarationNode", Ref(AliasDeclarationNodeClass))

	GetterDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	GetterDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GetterDeclarationNode", Ref(GetterDeclarationNodeClass))

	SetterDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SetterDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SetterDeclarationNode", Ref(SetterDeclarationNodeClass))

	AttrDeclarationNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	AttrDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AttrDeclarationNode", Ref(AttrDeclarationNodeClass))

	UsingExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UsingExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UsingExpressionNode", Ref(UsingExpressionNodeClass))

	IncludeExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	IncludeExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("IncludeExpressionNode", Ref(IncludeExpressionNodeClass))

	ExtendWhereBlockExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ExtendWhereBlockExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ExtendWhereBlockExpressionNode", Ref(ExtendWhereBlockExpressionNodeClass))

	ImplementExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ImplementExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ImplementExpressionNode", Ref(ImplementExpressionNodeClass))

	NewExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	NewExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("NewExpressionNode", Ref(NewExpressionNodeClass))

	GenericConstructorCallNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	GenericConstructorCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GenericConstructorCallNode", Ref(GenericConstructorCallNodeClass))

	ConstructorCallNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ConstructorCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ConstructorCallNode", Ref(ConstructorCallNodeClass))

	NamedCallArgumentNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	NamedCallArgumentNodeClass.IncludeMixin(NodeMixin)
	NamedCallArgumentNodeClass.IncludeMixin(NamedArgumentNodeMixin)
	ElkASTModule.AddConstantString("NamedCallArgumentNode", Ref(NamedCallArgumentNodeClass))

	DoubleSplatExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	DoubleSplatExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	DoubleSplatExpressionNodeClass.IncludeMixin(NamedArgumentNodeMixin)
	ElkASTModule.AddConstantString("DoubleSplatExpressionNode", Ref(DoubleSplatExpressionNodeClass))

	AttributeAccessNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	AttributeAccessNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AttributeAccessNode", Ref(AttributeAccessNodeClass))

	SubscriptExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SubscriptExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SubscriptExpressionNode", Ref(SubscriptExpressionNodeClass))

	NilSafeSubscriptExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	NilSafeSubscriptExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("NilSafeSubscriptExpressionNode", Ref(NilSafeSubscriptExpressionNodeClass))

	CallNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	CallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("CallNode", Ref(CallNodeClass))

	GenericMethodCallNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	GenericMethodCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GenericMethodCallNode", Ref(GenericMethodCallNodeClass))

	MethodCallNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MethodCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MethodCallNode", Ref(MethodCallNodeClass))

	ReceiverlessMethodCallNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ReceiverlessMethodCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ReceiverlessMethodCallNode", Ref(ReceiverlessMethodCallNodeClass))

	GenericReceiverlessMethodCallNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	GenericReceiverlessMethodCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GenericReceiverlessMethodCallNode", Ref(GenericReceiverlessMethodCallNodeClass))

	ScopedMacroCallNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ScopedMacroCallNodeClass.IncludeMixin(UsingEntryNodeMixin)
	ScopedMacroCallNodeClass.IncludeMixin(TypeNodeMixin)
	ScopedMacroCallNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ScopedMacroCallNodeClass.IncludeMixin(ComplexConstantNodeMixin)
	ElkASTModule.AddConstantString("ScopedMacroCallNode", Ref(ScopedMacroCallNodeClass))

	MacroCallNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MacroCallNodeClass.IncludeMixin(UsingEntryNodeMixin)
	MacroCallNodeClass.IncludeMixin(TypeNodeMixin)
	MacroCallNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	MacroCallNodeClass.IncludeMixin(ComplexConstantNodeMixin)
	ElkASTModule.AddConstantString("MacroCallNode", Ref(MacroCallNodeClass))

	ReceiverlessMacroCallNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ReceiverlessMacroCallNodeClass.IncludeMixin(UsingEntryNodeMixin)
	ReceiverlessMacroCallNodeClass.IncludeMixin(TypeNodeMixin)
	ReceiverlessMacroCallNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ReceiverlessMacroCallNodeClass.IncludeMixin(ComplexConstantNodeMixin)
	ElkASTModule.AddConstantString("ReceiverlessMacroCallNode", Ref(ReceiverlessMacroCallNodeClass))

	SplatExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SplatExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SplatExpressionNode", Ref(SplatExpressionNodeClass))

	KeyValueExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	KeyValueExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("KeyValueExpressionNode", Ref(KeyValueExpressionNodeClass))

	SymbolKeyValueExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SymbolKeyValueExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SymbolKeyValueExpressionNode", Ref(SymbolKeyValueExpressionNodeClass))

	WordArrayListLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	WordArrayListLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("WordArrayListLiteralNode", Ref(WordArrayListLiteralNodeClass))

	ArrayListLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ArrayListLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ArrayListLiteralNode", Ref(ArrayListLiteralNodeClass))

	SymbolArrayListLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SymbolArrayListLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("SymbolArrayListLiteralNode", Ref(SymbolArrayListLiteralNodeClass))

	HexArrayListLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HexArrayListLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("HexArrayListLiteralNode", Ref(HexArrayListLiteralNodeClass))

	BinArrayListLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BinArrayListLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("BinArrayListLiteralNode", Ref(BinArrayListLiteralNodeClass))

	ArrayTupleLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ArrayTupleLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ArrayTupleLiteralNode", Ref(ArrayTupleLiteralNodeClass))

	WordArrayTupleLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	WordArrayTupleLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("WordArrayTupleLiteralNode", Ref(WordArrayTupleLiteralNodeClass))

	SymbolArrayTupleLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SymbolArrayTupleLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("SymbolArrayTupleLiteralNode", Ref(SymbolArrayTupleLiteralNodeClass))

	HexArrayTupleLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HexArrayTupleLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("HexArrayTupleLiteralNode", Ref(HexArrayTupleLiteralNodeClass))

	BinArrayTupleLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BinArrayTupleLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("BinArrayTupleLiteralNode", Ref(BinArrayTupleLiteralNodeClass))

	HashSetLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HashSetLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HashSetLiteralNode", Ref(HashSetLiteralNodeClass))

	WordHashSetLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	WordHashSetLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("WordHashSetLiteralNode", Ref(WordHashSetLiteralNodeClass))

	SymbolHashSetLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SymbolHashSetLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("SymbolHashSetLiteralNode", Ref(SymbolHashSetLiteralNodeClass))

	HexHashSetLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HexHashSetLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("HexHashSetLiteralNode", Ref(HexHashSetLiteralNodeClass))

	BinHashSetLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BinHashSetLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("BinHashSetLiteralNode", Ref(BinHashSetLiteralNodeClass))

	HashMapLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HashMapLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HashMapLiteralNode", Ref(HashMapLiteralNodeClass))

	HashRecordLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HashRecordLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HashRecordLiteralNode", Ref(HashRecordLiteralNodeClass))

	RangeLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	RangeLiteralNodeClass.IncludeMixin(PatternExpressionNodeMixin)
	ElkASTModule.AddConstantString("RangeLiteralNode", Ref(RangeLiteralNodeClass))

	CaseNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	CaseNodeClass.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("CaseNode", Ref(CaseNodeClass))

	CatchNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	CatchNodeClass.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("CatchNode", Ref(CatchNodeClass))

	VariantTypeParameterNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	VariantTypeParameterNodeClass.IncludeMixin(TypeParameterNodeMixin)
	ElkASTModule.AddConstantString("VariantTypeParameterNode", Ref(VariantTypeParameterNodeClass))

	FormalParameterNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	FormalParameterNodeClass.IncludeMixin(ParameterNodeMixin)
	ElkASTModule.AddConstantString("FormalParameterNode", Ref(FormalParameterNodeClass))

	MethodParameterNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MethodParameterNodeClass.IncludeMixin(ParameterNodeMixin)
	ElkASTModule.AddConstantString("MethodParameterNode", Ref(MethodParameterNodeClass))

	SignatureParameterNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SignatureParameterNodeClass.IncludeMixin(ParameterNodeMixin)
	ElkASTModule.AddConstantString("SignatureParameterNode", Ref(SignatureParameterNodeClass))

	AttributeParameterNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	AttributeParameterNodeClass.IncludeMixin(ParameterNodeMixin)
	ElkASTModule.AddConstantString("AttributeParameterNode", Ref(AttributeParameterNodeClass))

	BoolLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BoolLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("BoolLiteralNode", Ref(BoolLiteralNodeClass))

	VoidTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	VoidTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("VoidTypeNode", Ref(VoidTypeNodeClass))

	NeverTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	NeverTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("NeverTypeNode", Ref(NeverTypeNodeClass))

	AnyTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	AnyTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("AnyTypeNode", Ref(AnyTypeNodeClass))

	UnionTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UnionTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UnionTypeNode", Ref(UnionTypeNodeClass))

	IntersectionTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	IntersectionTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("IntersectionTypeNode", Ref(IntersectionTypeNodeClass))

	BinaryTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BinaryTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("BinaryTypeExpressionNode", Ref(BinaryTypeNodeClass))

	NilableTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	NilableTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("NilableTypeNode", Ref(NilableTypeNodeClass))

	InstanceOfTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InstanceOfTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("InstanceOfTypeNode", Ref(InstanceOfTypeNodeClass))

	SingletonTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SingletonTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("SingletonTypeNode", Ref(SingletonTypeNodeClass))

	NotTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	NotTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("NotTypeNode", Ref(NotTypeNodeClass))

	ClosureTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ClosureTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("ClosureTypeNode", Ref(ClosureTypeNodeClass))

	UnaryTypeNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UnaryTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UnaryTypeNode", Ref(UnaryTypeNodeClass))

	AsPatternNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	AsPatternNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("AsPatternNode", Ref(AsPatternNodeClass))

	SymbolKeyValuePatternNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SymbolKeyValuePatternNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("SymbolKeyValuePatternNode", Ref(SymbolKeyValuePatternNodeClass))

	KeyValuePatternNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	KeyValuePatternNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("KeyValuePatternNode", Ref(KeyValuePatternNodeClass))

	ObjectPatternNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ObjectPatternNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("ObjectPatternNode", Ref(ObjectPatternNodeClass))

	RecordPatternNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	RecordPatternNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("RecordPatternNode", Ref(RecordPatternNodeClass))

	MapPatternNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MapPatternNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("MapPatternNode", Ref(MapPatternNodeClass))

	RestPatternNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	RestPatternNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("RestPatternNode", Ref(RestPatternNodeClass))

	SetPatternNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SetPatternNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("SetPatternNode", Ref(SetPatternNodeClass))

	ListPatternNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ListPatternNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("ListPatternNode", Ref(ListPatternNodeClass))

	TuplePatternNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	TuplePatternNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("TuplePatternNode", Ref(TuplePatternNodeClass))

	BinaryPatternNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BinaryPatternNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("BinaryPatternNode", Ref(BinaryPatternNodeClass))
}
