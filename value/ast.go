package value

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
var InstanceVariableNodeClass *Class              // Std::Elk::AST::InstanceVariableNode
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
var UsingEntryWithSubentriesNodeClass *Class      // Std::Elk::AST::UsingEntryWithSubentriesNode
var UsingAllEntryNodeClass *Class                 // Std::Elk::AST::UsingAllEntryNode
var ClosureLiteralNodeClass *Class                // Std::Elk::AST::ClosureLiteralNode
var ClassDeclarationNodeClass *Class              // Std::Elk::AST::ClassDeclarationNode
var ModuleDeclarationNodeClass *Class             // Std::Elk::AST::ModuleDeclarationNode
var MixinDeclarationNodeClass *Class              // Std::Elk::AST::MixinDeclarationNode
var InterfaceDeclarationNodeClass *Class          // Std::Elk::AST::InterfaceDeclarationNode
var StructDeclarationNodeClass *Class             // Std::Elk::AST::StructDeclarationNode
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

	StatementNodeMixin = NewMixin()
	StatementNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("StatementNode", Ref(StatementNodeMixin))

	StructBodyStatementNodeMixin = NewMixin()
	StructBodyStatementNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("StructBodyStatementNode", Ref(StructBodyStatementNodeMixin))

	ParameterNodeMixin = NewMixin()
	ParameterNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("ParameterNode", Ref(ParameterNodeMixin))

	PatternNodeMixin = NewMixin()
	PatternNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("PatternNode", Ref(PatternNodeMixin))

	TypeNodeMixin = NewMixin()
	TypeNodeMixin.IncludeMixin(NodeMixin)
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
	ElkASTModule.AddConstantString("ExpressionNode", Ref(ExpressionNodeMixin))

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
	IdentifierNodeMixin.IncludeMixin(ExpressionNodeMixin)
	IdentifierNodeMixin.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("IdentifierNode", Ref(IdentifierNodeMixin))

	RegexLiteralContentNodeMixin = NewMixin()
	RegexLiteralContentNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("RegexLiteralContentNode", Ref(RegexLiteralContentNodeMixin))

	RegexLiteralNodeMixin = NewMixin()
	RegexLiteralNodeMixin.IncludeMixin(PatternNodeMixin)
	RegexLiteralNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("RegexLiteralNode", Ref(RegexLiteralNodeMixin))

	StringLiteralContentNodeMixin = NewMixin()
	StringLiteralContentNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("StringLiteralContentNode", Ref(StringLiteralContentNodeMixin))

	StringLiteralNodeMixin = NewMixin()
	StringLiteralNodeMixin.IncludeMixin(PatternNodeMixin)
	StringLiteralNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("StringLiteralNode", Ref(StringLiteralNodeMixin))

	SimpleStringLiteralNodeMixin = NewMixin()
	SimpleStringLiteralNodeMixin.IncludeMixin(StringLiteralNodeMixin)
	SimpleStringLiteralNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SimpleStringLiteralNode", Ref(SimpleStringLiteralNodeMixin))

	TypeParameterNodeMixin = NewMixin()
	TypeParameterNodeMixin.IncludeMixin(NodeMixin)
	ElkASTModule.AddConstantString("TypeParameterNode", Ref(TypeParameterNodeMixin))

	ComplexConstantNodeMixin = NewMixin()
	ComplexConstantNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ComplexConstantNodeMixin.IncludeMixin(TypeNodeMixin)
	ComplexConstantNodeMixin.IncludeMixin(PatternNodeMixin)
	ComplexConstantNodeMixin.IncludeMixin(UsingEntryNodeMixin)
	ElkASTModule.AddConstantString("ComplexConstantNode", Ref(ComplexConstantNodeMixin))

	ConstantNodeMixin = NewMixin()
	ConstantNodeMixin.IncludeMixin(ExpressionNodeMixin)
	ConstantNodeMixin.IncludeMixin(TypeNodeMixin)
	ConstantNodeMixin.IncludeMixin(ComplexConstantNodeMixin)
	ConstantNodeMixin.IncludeMixin(UsingEntryNodeMixin)
	ElkASTModule.AddConstantString("ConstantNode", Ref(ConstantNodeMixin))

	InvalidNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InvalidNodeClass.IncludeMixin(StatementNodeMixin)
	InvalidNodeClass.IncludeMixin(ExpressionNodeMixin)
	InvalidNodeClass.IncludeMixin(StructBodyStatementNodeMixin)
	InvalidNodeClass.IncludeMixin(ParameterNodeMixin)
	InvalidNodeClass.IncludeMixin(TypeNodeMixin)
	InvalidNodeClass.IncludeMixin(PatternNodeMixin)
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
	UnaryExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	UnaryExpressionNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("UnaryExpressionNode", Ref(UnaryExpressionNodeClass))

	TrueLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	TrueLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	TrueLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	TrueLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("TrueLiteralNode", Ref(TrueLiteralNodeClass))

	FalseLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	FalseLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	FalseLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	FalseLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("FalseLiteralNode", Ref(FalseLiteralNodeClass))

	NilLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	NilLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	NilLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	NilLiteralNodeClass.IncludeMixin(PatternNodeMixin)
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

	InstanceVariableNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InstanceVariableNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InstanceVariableNode", Ref(InstanceVariableNodeClass))

	SimpleSymbolLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SimpleSymbolLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	SimpleSymbolLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	SimpleSymbolLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	SimpleSymbolLiteralNodeClass.IncludeMixin(SymbolLiteralNodeMixin)
	SimpleSymbolLiteralNodeClass.IncludeMixin(SymbolCollectionContentNodeMixin)
	ElkASTModule.AddConstantString("SimpleSymbolLiteralNode", Ref(SimpleSymbolLiteralNodeClass))

	IntLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	IntLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	IntLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	IntLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	IntLiteralNodeClass.IncludeMixin(IntCollectionContentNodeMixin)
	ElkASTModule.AddConstantString("IntLiteralNode", Ref(IntLiteralNodeClass))

	Int64LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Int64LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Int64LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	Int64LiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("Int64LiteralNode", Ref(Int64LiteralNodeClass))

	Int32LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Int32LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Int32LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	Int32LiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("Int32LiteralNode", Ref(Int32LiteralNodeClass))

	Int16LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Int16LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Int16LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	Int16LiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("Int16LiteralNode", Ref(Int16LiteralNodeClass))

	Int8LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Int8LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Int8LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	Int8LiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("Int8LiteralNode", Ref(Int8LiteralNodeClass))

	UInt64LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UInt64LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	UInt64LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	UInt64LiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("UInt64LiteralNode", Ref(UInt64LiteralNodeClass))

	UInt32LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UInt32LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	UInt32LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	UInt32LiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("UInt32LiteralNode", Ref(UInt32LiteralNodeClass))

	UInt16LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UInt16LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	UInt16LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	UInt16LiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("UInt16LiteralNode", Ref(UInt16LiteralNodeClass))

	UInt8LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UInt8LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	UInt8LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	UInt8LiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("UInt8LiteralNode", Ref(UInt8LiteralNodeClass))

	FloatLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	FloatLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	FloatLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	FloatLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("FloatLiteralNode", Ref(FloatLiteralNodeClass))

	BigFloatLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BigFloatLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	BigFloatLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	BigFloatLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("BigFloatLiteralNode", Ref(BigFloatLiteralNodeClass))

	Float64LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Float64LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Float64LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	Float64LiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("Float64LiteralNode", Ref(Float64LiteralNodeClass))

	Float32LiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	Float32LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Float32LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	Float32LiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("Float32LiteralNode", Ref(Float32LiteralNodeClass))

	UninterpolatedRegexLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	UninterpolatedRegexLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	UninterpolatedRegexLiteralNodeClass.IncludeMixin(PatternNodeMixin)
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
	InterpolatedRegexLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	InterpolatedRegexLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	InterpolatedRegexLiteralNodeClass.IncludeMixin(RegexLiteralNodeMixin)
	ElkASTModule.AddConstantString("InterpolatedRegexLiteralNode", Ref(InterpolatedRegexLiteralNodeClass))

	CharLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	CharLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	CharLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	CharLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("CharLiteralNode", Ref(CharLiteralNodeClass))

	RawCharLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	RawCharLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	RawCharLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	RawCharLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("RawCharLiteralNode", Ref(RawCharLiteralNodeClass))

	RawStringLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	RawStringLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	RawStringLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	RawStringLiteralNodeClass.IncludeMixin(PatternNodeMixin)
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
	InterpolatedStringLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	InterpolatedStringLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	InterpolatedStringLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	InterpolatedStringLiteralNodeClass.IncludeMixin(StringLiteralNodeMixin)
	ElkASTModule.AddConstantString("InterpolatedStringLiteralNode", Ref(InterpolatedStringLiteralNodeClass))

	DoubleQuotedStringLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	DoubleQuotedStringLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	DoubleQuotedStringLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	DoubleQuotedStringLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	DoubleQuotedStringLiteralNodeClass.IncludeMixin(StringLiteralNodeMixin)
	DoubleQuotedStringLiteralNodeClass.IncludeMixin(SimpleStringLiteralNodeMixin)
	ElkASTModule.AddConstantString("DoubleQuotedStringLiteralNode", Ref(DoubleQuotedStringLiteralNodeClass))

	InterpolatedSymbolLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	InterpolatedSymbolLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	InterpolatedSymbolLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	InterpolatedSymbolLiteralNodeClass.IncludeMixin(PatternNodeMixin)
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
	PublicIdentifierNodeClass.IncludeMixin(ExpressionNodeMixin)
	PublicIdentifierNodeClass.IncludeMixin(PatternNodeMixin)
	PublicIdentifierNodeClass.IncludeMixin(UsingSubentryNodeMixin)
	PublicIdentifierNodeClass.IncludeMixin(IdentifierNodeMixin)
	ElkASTModule.AddConstantString("PublicIdentifierNode", Ref(PublicIdentifierNodeClass))

	PublicIdentifierNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PublicIdentifierNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("PublicIdentifierNode", Ref(PublicIdentifierNodeClass))

	PublicIdentifierAsNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PublicIdentifierAsNodeClass.IncludeMixin(ExpressionNodeMixin)
	PublicIdentifierAsNodeClass.IncludeMixin(UsingSubentryNodeMixin)
	ElkASTModule.AddConstantString("PublicIdentifierAsNode", Ref(PublicIdentifierAsNodeClass))

	PrivateIdentifierNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PrivateIdentifierNodeClass.IncludeMixin(ExpressionNodeMixin)
	PrivateIdentifierNodeClass.IncludeMixin(PatternNodeMixin)
	PrivateIdentifierNodeClass.IncludeMixin(IdentifierNodeMixin)
	ElkASTModule.AddConstantString("PrivateIdentifierNode", Ref(PrivateIdentifierNodeClass))

	PublicConstantNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	PublicConstantNodeClass.IncludeMixin(ExpressionNodeMixin)
	PublicConstantNodeClass.IncludeMixin(TypeNodeMixin)
	PublicConstantNodeClass.IncludeMixin(PatternNodeMixin)
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
	PrivateConstantNodeClass.IncludeMixin(ExpressionNodeMixin)
	PrivateConstantNodeClass.IncludeMixin(TypeNodeMixin)
	PrivateConstantNodeClass.IncludeMixin(PatternNodeMixin)
	PrivateConstantNodeClass.IncludeMixin(UsingEntryNodeMixin)
	PrivateConstantNodeClass.IncludeMixin(ComplexConstantNodeMixin)
	PrivateConstantNodeClass.IncludeMixin(ConstantNodeMixin)
	ElkASTModule.AddConstantString("PrivateConstantNode", Ref(PrivateConstantNodeClass))

	AsExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	AsExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AsExpressionNode", Ref(AsExpressionNodeClass))

	DoExpressionNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	DoExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("DoExpressionNode", Ref(DoExpressionNodeClass))

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
	ConstantLookupNodeClass.IncludeMixin(ExpressionNodeMixin)
	ConstantLookupNodeClass.IncludeMixin(TypeNodeMixin)
	ConstantLookupNodeClass.IncludeMixin(PatternNodeMixin)
	ConstantLookupNodeClass.IncludeMixin(UsingEntryNodeMixin)
	ConstantLookupNodeClass.IncludeMixin(ComplexConstantNodeMixin)
	ElkASTModule.AddConstantString("ConstantLookupNode", Ref(ConstantLookupNodeClass))

	MethodLookupNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	MethodLookupNodeClass.IncludeMixin(ExpressionNodeMixin)
	MethodLookupNodeClass.IncludeMixin(UsingEntryNodeMixin)
	ElkASTModule.AddConstantString("MethodLookupNode", Ref(MethodLookupNodeClass))

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
	ElkASTModule.AddConstantString("MethodDefinitionNode", Ref(MethodDefinitionNodeClass))

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
	GenericConstantNodeClass.IncludeMixin(ExpressionNodeMixin)
	GenericConstantNodeClass.IncludeMixin(TypeNodeMixin)
	GenericConstantNodeClass.IncludeMixin(PatternNodeMixin)
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
	WordArrayListLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	WordArrayListLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("WordArrayListLiteralNode", Ref(WordArrayListLiteralNodeClass))

	ArrayListLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ArrayListLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ArrayListLiteralNode", Ref(ArrayListLiteralNodeClass))

	SymbolArrayListLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SymbolArrayListLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	SymbolArrayListLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("SymbolArrayListLiteralNode", Ref(SymbolArrayListLiteralNodeClass))

	HexArrayListLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HexArrayListLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	HexArrayListLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("HexArrayListLiteralNode", Ref(HexArrayListLiteralNodeClass))

	BinArrayListLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BinArrayListLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	BinArrayListLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("BinArrayListLiteralNode", Ref(BinArrayListLiteralNodeClass))

	ArrayTupleLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ArrayTupleLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ArrayTupleLiteralNode", Ref(ArrayTupleLiteralNodeClass))

	WordArrayTupleLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	WordArrayTupleLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	WordArrayTupleLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("WordArrayTupleLiteralNode", Ref(WordArrayTupleLiteralNodeClass))

	SymbolArrayTupleLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SymbolArrayTupleLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	SymbolArrayTupleLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("SymbolArrayTupleLiteralNode", Ref(SymbolArrayTupleLiteralNodeClass))

	HexArrayTupleLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HexArrayTupleLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	HexArrayTupleLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("HexArrayTupleLiteralNode", Ref(HexArrayTupleLiteralNodeClass))

	BinArrayTupleLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BinArrayTupleLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	BinArrayTupleLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("BinArrayTupleLiteralNode", Ref(BinArrayTupleLiteralNodeClass))

	HashSetLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HashSetLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HashSetLiteralNode", Ref(HashSetLiteralNodeClass))

	WordHashSetLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	WordHashSetLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	WordHashSetLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("WordHashSetLiteralNode", Ref(WordHashSetLiteralNodeClass))

	SymbolHashSetLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	SymbolHashSetLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	SymbolHashSetLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("SymbolHashSetLiteralNode", Ref(SymbolHashSetLiteralNodeClass))

	HexHashSetLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HexHashSetLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	HexHashSetLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("HexHashSetLiteralNode", Ref(HexHashSetLiteralNodeClass))

	BinHashSetLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	BinHashSetLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	BinHashSetLiteralNodeClass.IncludeMixin(PatternNodeMixin)
	ElkASTModule.AddConstantString("BinHashSetLiteralNode", Ref(BinHashSetLiteralNodeClass))

	HashMapLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HashMapLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HashMapLiteralNode", Ref(HashMapLiteralNodeClass))

	HashRecordLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	HashRecordLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HashRecordLiteralNode", Ref(HashRecordLiteralNodeClass))

	RangeLiteralNodeClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	RangeLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	RangeLiteralNodeClass.IncludeMixin(PatternNodeMixin)
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
