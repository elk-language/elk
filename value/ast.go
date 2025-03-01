package value

var ElkASTModule *Module // Std::ElkAST
var NodeClass *Class     // Std::ElkAST::Node

var StructBodyStatementNodeMixin *Class // Std::ElkAST::StructBodyStatementNode
var StatementNodeMixin *Mixin           // Std::ElkAST::StatementNode
var ParameterNodeMixin *Mixin           // Std::ElkAST::ParameterNode
var TypeNodeMixin *Mixin                // Std::ElkAST::TypeNode
var PatternNodeMixin *Mixin             // Std::ElkAST::PatternNode

var ExpressionStatementNodeClass *Class // Std::ElkAST::ExpressionStatementNode
var EmptyStatementNodeClass *Class      // Std::ElkAST::EmptyStatementNode
var ImportStatementNodeClass *Class     // Std::ElkAST::ImportStatementNode
var ParameterStatementNodeClass *Class  // Std::ElkAST::ParameterStatementNode

var ProgramNodeClass *Class                       // Std::ElkAST::ProgramNode
var ExpressionNodeMixin *Mixin                    // Std::ElkAST::ExpressionNode
var InvalidNodeClass *Class                       // Std::ElkAST::InvalidNode
var TypeExpressionNodeClass *Class                // Std::ElkAST::TypeExpressionNode
var InstanceVariableDeclarationNodeClass *Class   // Std::ElkAST::InstanceVariableDeclarationNode
var VariablePatternDeclarationNodeClass *Class    // Std::ElkAST::VariablePatternDeclarationNode
var VariableDeclarationNodeClass *Class           // Std::ElkAST::VariableDeclarationNode
var ValuePatternDeclarationNodeClass *Class       // Std::ElkAST::ValuePatternDeclarationNode
var ValueDeclarationNodeClass *Class              // Std::ElkAST::ValueDeclarationNode
var PostfixExpressionNodeClass *Class             // Std::ElkAST::PostfixExpressionNode
var ModifierNodeClass *Class                      // Std::ElkAST::ModifierNode
var ModifierIfElseNodeClass *Class                // Std::ElkAST::ModifierIfElseNode
var ModifierForInNodeClass *Class                 // Std::ElkAST::ModifierForInNode
var AssignmentExpressionNodeClass *Class          // Std::ElkAST::AssignmentExpressionNode
var BinaryExpressionNodeClass *Class              // Std::ElkAST::BinaryExpressionNode
var LogicalExpressionNodeClass *Class             // Std::ElkAST::LogicalExpressionNode
var UnaryExpressionNodeClass *Class               // Std::ElkAST::UnaryExpressionNode
var TrueLiteralNodeClass *Class                   // Std::ElkAST::TrueLiteralNode
var FalseLiteralNodeClass *Class                  // Std::ElkAST::FalseLiteralNode
var NilLiteralNodeClass *Class                    // Std::ElkAST::NilLiteralNode
var UndefinedLiteralNodeClass *Class              // Std::ElkAST::UndefinedLiteralNode
var SelfLiteralNodeClass *Class                   // Std::ElkAST::SelfLiteralNode
var InstanceVariableNodeClass *Class              // Std::ElkAST::InstanceVariableNode
var SimpleSymbolLiteralNodeClass *Class           // Std::ElkAST::SimpleSymbolLiteralNode
var IntLiteralNodeClass *Class                    // Std::ElkAST::IntLiteralNode
var Int64LiteralNodeClass *Class                  // Std::ElkAST::Int64LiteralNode
var Int32LiteralNodeClass *Class                  // Std::ElkAST::Int32LiteralNode
var Int16LiteralNodeClass *Class                  // Std::ElkAST::Int16LiteralNode
var Int8LiteralNodeClass *Class                   // Std::ElkAST::Int8LiteralNode
var UInt64LiteralNodeClass *Class                 // Std::ElkAST::UInt64LiteralNode
var UInt32LiteralNodeClass *Class                 // Std::ElkAST::UInt32LiteralNode
var UInt16LiteralNodeClass *Class                 // Std::ElkAST::UInt16LiteralNode
var UInt8LiteralNodeClass *Class                  // Std::ElkAST::UInt8LiteralNode
var FloatLiteralNodeClass *Class                  // Std::ElkAST::FloatLiteralNode
var BigFloatLiteralNodeClass *Class               // Std::ElkAST::BigFloatLiteralNode
var Float32LiteralNodeClass *Class                // Std::ElkAST::Float32LiteralNode
var Float64LiteralNodeClass *Class                // Std::ElkAST::Float64LiteralNode
var UninterpolatedRegexLiteralNodeClass *Class    // Std::ElkAST::UninterpolatedRegexLiteralNode
var RegexLiteralContentSectionNodeClass *Class    // Std::ElkAST::RegexLiteralContentSection
var RegexInterpolationNodeClass *Class            // Std::ElkAST::RegexInterpolationNode
var InterpolatedRegexLiteralNodeClass *Class      // Std::ElkAST::InterpolatedRegexLiteralNode
var CharLiteralNodeClass *Class                   // Std::ElkAST::CharLiteralNode
var RawCharLiteralNodeClass *Class                // Std::ElkAST::RawCharLiteralNode
var RawStringLiteralNodeClass *Class              // Std::ElkAST::RawStringLiteralNode
var StringLiteralContentSectionNodeClass *Class   // Std::ElkAST::StringLiteralContentSectionNode
var StringInspectInterpolationNodeClass *Class    // Std::ElkAST::StringInspectInterpolationNode
var StringInterpolationNodeClass *Class           // Std::ElkAST::StringInterpolationNode
var InterpolatedStringLiteralNodeClass *Class     // Std::ElkAST::InterpolatedStringLiteralNode
var DoubleQuotedStringLiteralNodeClass *Class     // Std::ElkAST::DoubleQuotedStringLiteralNode
var InterpolatedSymbolLiteralNodeClass *Class     // Std::ElkAST::InterpolatedSymbolLiteralNode
var ConstantAsNodeClass *Class                    // Std::ElkAST::ConstantAsNode
var MethodLookupAsNodeClass *Class                // Std::ElkAST::MethodLookupAsNode
var PublicIdentifierNodeClass *Class              // Std::ElkAST::PublicIdentifierNode
var PublicIdentifierAsNodeClass *Class            // Std::ElkAST::PublicIdentifierAsNode
var PrivateIdentifierNodeClass *Class             // Std::ElkAST::PrivateIdentifierNode
var PublicConstantNodeClass *Class                // Std::ElkAST::PublicConstantNode
var PrivateConstantNodeClass *Class               // Std::ElkAST::PrivateConstantNode
var AsExpressionNodeClass *Class                  // Std::ElkAST::AsExpressionNode
var DoExpressionNodeClass *Class                  // Std::ElkAST::DoExpressionNode
var SingletonBlockExpressionNodeClass *Class      // Std::ElkAST::SingletonBlockExpressionNode
var SwitchExpressionNodeClass *Class              // Std::ElkAST::SwitchExpressionNode
var IfExpressionNodeClass *Class                  // Std::ElkAST::IfExpressionNode
var UnlessExpressionNodeClass *Class              // Std::ElkAST::UnlessExpressionNode
var WhileExpressionNodeClass *Class               // Std::ElkAST::WhileExpressionNode
var UntilExpressionNodeClass *Class               // Std::ElkAST::UntilExpressionNode
var LoopExpressionNodeClass *Class                // Std::ElkAST::LoopExpressionNode
var NumericForExpressionNodeClass *Class          // Std::ElkAST::NumericForExpressionNode
var ForInExpressionNodeClass *Class               // Std::ElkAST::ForInExpressionNode
var BreakExpressionNodeClass *Class               // Std::ElkAST::BreakExpressionNode
var LabeledExpressionNodeClass *Class             // Std::ElkAST::LabeledExpressionNode
var GoExpressionNodeClass *Class                  // Std::ElkAST::GoExpressionNode
var ReturnExpressionNodeClass *Class              // Std::ElkAST::ReturnExpressionNode
var YieldExpressionNodeClass *Class               // Std::ElkAST::YieldExpressionNode
var ContinueExpressionNodeClass *Class            // Std::ElkAST::ContinueExpressionNode
var ThrowExpressionNodeClass *Class               // Std::ElkAST::ThrowExpressionNode
var MustExpressionNodeClass *Class                // Std::ElkAST::MustExpressionNode
var TryExpressionNodeClass *Class                 // Std::ElkAST::TryExpressionNode
var AwaitExpressionNodeClass *Class               // Std::ElkAST::AwaitExpressionNode
var TypeofExpressionNodeClass *Class              // Std::ElkAST::TypeofExpressionNode
var ConstantLookupNodeClass *Class                // Std::ElkAST::ConstantLookupNode
var MethodLookupNodeClass *Class                  // Std::ElkAST::MethodLookupNode
var UsingEntryWithSubentriesNodeClass *Class      // Std::ElkAST::UsingEntryWithSubentriesNode
var UsingAllEntryNodeClass *Class                 // Std::ElkAST::UsingAllEntryNode
var ClosureLiteralNodeClass *Class                // Std::ElkAST::ClosureLiteralNode
var ClassDeclarationNodeClass *Class              // Std::ElkAST::ClassDeclarationNode
var ModuleDeclarationNodeClass *Class             // Std::ElkAST::ModuleDeclarationNode
var MixinDeclarationNodeClass *Class              // Std::ElkAST::MixinDeclarationNode
var InterfaceDeclarationNodeClass *Class          // Std::ElkAST::InterfaceDeclarationNode
var StructDeclarationNodeClass *Class             // Std::ElkAST::StructDeclarationNode
var MethodDefinitionNodeClass *Class              // Std::ElkAST::MethodDefinitionNode
var InitDefinitionNodeClass *Class                // Std::ElkAST::InitDefinitionNode
var MethodSignatureDefinitionNodeClass *Class     // Std::ElkAST::MethodSignatureDefinitionNode
var GenericConstantNodeClass *Class               // Std::ElkAST::GenericConstantNode
var GenericTypeDefinitionNodeClass *Class         // Std::ElkAST::GenericTypeDefinitionNode
var TypeDefinitionNodeClass *Class                // Std::ElkAST::TypeDefinitionNode
var AliasDeclarationEntryClass *Class             // Std::ElkAST::AliasDeclarationEntry
var AliasDeclarationNodeClass *Class              // Std::ElkAST::AliasDeclarationNode
var GetterDeclarationNodeClass *Class             // Std::ElkAST::GetterDeclarationNode
var SetterDeclarationNodeClass *Class             // Std::ElkAST::SetterDeclarationNode
var AttrDeclarationNodeClass *Class               // Std::ElkAST::AttrDeclarationNode
var UsingExpressionNodeClass *Class               // Std::ElkAST::UsingExpressionNode
var IncludeExpressionNodeClass *Class             // Std::ElkAST::IncludeExpressionNode
var ExtendWhereBlockExpressionNodeClass *Class    // Std::ElkAST::ExtendWhereBlockExpressionNode
var ImplementExpressionNodeClass *Class           // Std::ElkAST::ImplementExpressionNode
var NewExpressionNodeClass *Class                 // Std::ElkAST::NewExpressionNode
var GenericConstructorCallNodeClass *Class        // Std::ElkAST::GenericConstructorCallNode
var ConstructorCallNodeClass *Class               // Std::ElkAST::ConstructorCallNode
var NamedCallArgumentNodeClass *Class             // Std::ElkAST::NamedCallArgumentNode
var DoubleSplatExpressionNodeClass *Class         // Std::ElkAST::DoubleSplatExpressionNode
var AttributeAccessNodeClass *Class               // Std::ElkAST::AttributeAccessNode
var SubscriptExpressionNodeClass *Class           // Std::ElkAST::SubscriptExpressionNode
var NilSafeSubscriptExpressionNodeClass *Class    // Std::ElkAST::NilSafeSubscriptExpressionNode
var CallNodeClass *Class                          // Std::ElkAST::CallNode
var GenericMethodCallNodeClass *Class             // Std::ElkAST::GenericMethodCallNode
var MethodCallNodeClass *Class                    // Std::ElkAST::MethodCallNode
var ReceiverlessMethodCallNodeClass *Class        // Std::ElkAST::ReceiverlessMethodCallNode
var GenericReceiverlessMethodCallNodeClass *Class // Std::ElkAST::GenericReceiverlessMethodCallNode
var SplatExpressionNodeClass *Class               // Std::ElkAST::SplatExpressionNode
var KeyValueExpressionNodeClass *Class            // Std::ElkAST::KeyValueExpressionNode
var SymbolKeyValueExpressionNodeClass *Class      // Std::ElkAST::SymbolKeyValueExpressionNode
var WordArrayListLiteralNodeClass *Class          // Std::ElkAST::WordArrayListLiteralNode
var ArrayListLiteralNodeClass *Class              // Std::ElkAST::ArrayListLiteralNode
var SymbolArrayListLiteralNodeClass *Class        // Std::ElkAST::SymbolArrayListLiteralNode
var HexArrayListLiteralNodeClass *Class           // Std::ElkAST::HexArrayListLiteralNode
var BinArrayListLiteralNodeClass *Class           // Std::ElkAST::BinArrayListLiteralNode
var ArrayTupleLiteralNodeClass *Class             // Std::ElkAST::ArrayTupleLiteralNode
var WordArrayTupleLiteralNodeClass *Class         // Std::ElkAST::WordArrayTupleLiteralNode
var SymbolArrayTupleLiteralNodeClass *Class       // Std::ElkAST::SymbolArrayTupleLiteralNode
var HexArrayTupleLiteralNodeClass *Class          // Std::ElkAST::HexArrayTupleLiteralNode
var BinArrayTupleLiteralNodeClass *Class          // Std::ElkAST::BinArrayTupleLiteralNode
var HashSetLiteralNodeClass *Class                // Std::ElkAST::HashSetLiteralNode
var WordHashSetLiteralNodeClass *Class            // Std::ElkAST::WordHashSetLiteralNode
var SymbolHashSetLiteralNodeClass *Class          // Std::ElkAST::SymbolHashSetLiteralNode
var HexHashSetLiteralNodeClass *Class             // Std::ElkAST::HexHashSetLiteralNode
var BinHashSetLiteralNodeClass *Class             // Std::ElkAST::BinHashSetLiteralNode
var HashMapLiteralNodeClass *Class                // Std::ElkAST::HashMapLiteralNode
var HashRecordLiteralNodeClass *Class             // Std::ElkAST::HashRecordLiteralNode
var RangeLiteralNodeClass *Class                  // Std::ElkAST::RangeLiteralNode
var VariantTypeParameterNodeClass *Class          // Std::ElkAST::VariantTypeParameterNode
var FormalParameterNodeClass *Class               // Std::ElkAST::FormalParameterNode
var MethodParameterNodeClass *Class               // Std::ElkAST::MethodParameterNode
var SignatureParameterNodeClass *Class            // Std::ElkAST::SignatureParameterNode
var AttributeParameterNodeClass *Class            // Std::ElkAST::AttributeParameterNode
var BoolLiteralNodeClass *Class                   // Std::ElkAST::BoolLiteralNode
var VoidTypeNodeClass *Class                      // Std::ElkAST::VoidTypeNode
var NeverTypeNodeClass *Class                     // Std::ElkAST::NeverTypeNode
var AnyTypeNodeClass *Class                       // Std::ElkAST::AnyTypeNode
var UnionTypeNodeClass *Class                     // Std::ElkAST::UnionTypeNode
var IntersectionTypeNodeClass *Class              // Std::ElkAST::IntersectionTypeNode
var BinaryTypeNodeClass *Class                    // Std::ElkAST::BinaryTypeNode
var NilableTypeNodeClass *Class                   // Std::ElkAST::NilableTypeNode
var InstanceOfTypeNodeClass *Class                // Std::ElkAST::InstanceOfTypeNode
var SingletonTypeNodeClass *Class                 // Std::ElkAST::SingletonTypeNode
var NotTypeNodeClass *Class                       // Std::ElkAST::NotTypeNode
var ClosureTypeNodeClass *Class                   // Std::ElkAST::ClosureTypeNode
var UnaryTypeNodeClass *Class                     // Std::ElkAST::UnaryTypeNode

func initAST() {
	ElkASTModule = NewModule()
	StdModule.AddConstantString("ElkAST", Ref(ElkASTModule))

	NodeClass = NewClass()
	ElkASTModule.AddConstantString("Node", Ref(NodeClass))

	StatementNodeMixin = NewMixin()
	ElkASTModule.AddConstantString("StatementNode", Ref(StatementNodeMixin))

	StructBodyStatementNodeMixin = NewMixin()
	ElkASTModule.AddConstantString("StructBodyStatementNode", Ref(StructBodyStatementNodeMixin))

	ParameterNodeMixin = NewMixin()
	ElkASTModule.AddConstantString("ParameterNode", Ref(ParameterNodeMixin))

	ExpressionStatementNodeClass = NewClass()
	ExpressionStatementNodeClass.IncludeMixin(StatementNodeMixin)
	ElkASTModule.AddConstantString("ExpressionStatementNode", Ref(ExpressionStatementNodeClass))

	EmptyStatementNodeClass = NewClass()
	EmptyStatementNodeClass.IncludeMixin(StatementNodeMixin)
	EmptyStatementNodeClass.IncludeMixin(StructBodyStatementNodeMixin)
	ElkASTModule.AddConstantString("EmptyStatementNode", Ref(EmptyStatementNodeClass))

	ImportStatementNodeClass = NewClass()
	ImportStatementNodeClass.IncludeMixin(StatementNodeMixin)
	ElkASTModule.AddConstantString("ImportStatementNode", Ref(ImportStatementNodeClass))

	ParameterStatementNodeClass = NewClass()
	ParameterStatementNodeClass.IncludeMixin(StatementNodeMixin)
	ParameterStatementNodeClass.IncludeMixin(StructBodyStatementNodeMixin)
	ElkASTModule.AddConstantString("ParameterStatementNode", Ref(ParameterStatementNodeClass))

	ProgramNodeClass = NewClass()
	ElkASTModule.AddConstantString("ProgramNode", Ref(ProgramNodeClass))

	ExpressionNodeMixin = NewMixin()
	ElkASTModule.AddConstantString("ExpressionNode", Ref(ExpressionNodeMixin))

	InvalidNodeClass = NewClass()
	InvalidNodeClass.IncludeMixin(StatementNodeMixin)
	InvalidNodeClass.IncludeMixin(ExpressionNodeMixin)
	InvalidNodeClass.IncludeMixin(StructBodyStatementNodeMixin)
	InvalidNodeClass.IncludeMixin(ParameterNodeMixin)
	InvalidNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("InvalidNode", Ref(InvalidNodeClass))

	TypeExpressionNodeClass = NewClass()
	TypeExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("TypeExpressionNode", Ref(TypeExpressionNodeClass))

	InstanceVariableDeclarationNodeClass = NewClass()
	InstanceVariableDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InstanceVariableDeclarationNode", Ref(InstanceVariableDeclarationNodeClass))

	VariablePatternDeclarationNodeClass = NewClass()
	VariablePatternDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("VariablePatternDeclarationNode", Ref(VariablePatternDeclarationNodeClass))

	VariableDeclarationNodeClass = NewClass()
	VariableDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("VariableDeclarationNode", Ref(VariableDeclarationNodeClass))

	ValuePatternDeclarationNodeClass = NewClass()
	ValuePatternDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ValuePatternDeclarationNode", Ref(ValuePatternDeclarationNodeClass))

	ValueDeclarationNodeClass = NewClass()
	ValueDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ValueDeclarationNode", Ref(ValueDeclarationNodeClass))

	PostfixExpressionNodeClass = NewClass()
	PostfixExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("PostfixExpressionNode", Ref(PostfixExpressionNodeClass))

	ModifierNodeClass = NewClass()
	ModifierNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ModifierNode", Ref(ModifierNodeClass))

	ModifierIfElseNodeClass = NewClass()
	ModifierIfElseNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ModifierIfElseNode", Ref(ModifierIfElseNodeClass))

	ModifierForInNodeClass = NewClass()
	ModifierForInNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ModifierForInNode", Ref(ModifierForInNodeClass))

	AssignmentExpressionNodeClass = NewClass()
	AssignmentExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AssignmentExpressionNode", Ref(AssignmentExpressionNodeClass))

	BinaryExpressionNodeClass = NewClass()
	BinaryExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("BinaryExpressionNode", Ref(BinaryExpressionNodeClass))

	LogicalExpressionNodeClass = NewClass()
	LogicalExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("LogicalExpressionNode", Ref(LogicalExpressionNodeClass))

	UnaryExpressionNodeClass = NewClass()
	UnaryExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UnaryExpressionNode", Ref(UnaryExpressionNodeClass))

	TrueLiteralNodeClass = NewClass()
	TrueLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	TrueLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("TrueLiteralNode", Ref(TrueLiteralNodeClass))

	FalseLiteralNodeClass = NewClass()
	FalseLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	FalseLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("FalseLiteralNode", Ref(FalseLiteralNodeClass))

	NilLiteralNodeClass = NewClass()
	NilLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	NilLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("NilLiteralNode", Ref(NilLiteralNodeClass))

	UndefinedLiteralNodeClass = NewClass()
	UndefinedLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UndefinedLiteralNode", Ref(UndefinedLiteralNodeClass))

	SelfLiteralNodeClass = NewClass()
	SelfLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	SelfLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("SelfLiteralNode", Ref(SelfLiteralNodeClass))

	InstanceVariableNodeClass = NewClass()
	InstanceVariableNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InstanceVariableNode", Ref(InstanceVariableNodeClass))

	SimpleSymbolLiteralNodeClass = NewClass()
	SimpleSymbolLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	SimpleSymbolLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("SimpleSymbolLiteralNode", Ref(SimpleSymbolLiteralNodeClass))

	IntLiteralNodeClass = NewClass()
	IntLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	IntLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("IntLiteralNode", Ref(IntLiteralNodeClass))

	Int64LiteralNodeClass = NewClass()
	Int64LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Int64LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("Int64LiteralNode", Ref(Int64LiteralNodeClass))

	Int32LiteralNodeClass = NewClass()
	Int32LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Int32LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("Int32LiteralNode", Ref(Int32LiteralNodeClass))

	Int16LiteralNodeClass = NewClass()
	Int16LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Int16LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("Int16LiteralNode", Ref(Int16LiteralNodeClass))

	Int8LiteralNodeClass = NewClass()
	Int8LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Int8LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("Int8LiteralNode", Ref(Int8LiteralNodeClass))

	UInt64LiteralNodeClass = NewClass()
	UInt64LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	UInt64LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UInt64LiteralNode", Ref(UInt64LiteralNodeClass))

	UInt32LiteralNodeClass = NewClass()
	UInt32LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	UInt32LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UInt32LiteralNode", Ref(UInt32LiteralNodeClass))

	UInt16LiteralNodeClass = NewClass()
	UInt16LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	UInt16LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UInt16LiteralNode", Ref(UInt16LiteralNodeClass))

	UInt8LiteralNodeClass = NewClass()
	UInt8LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	UInt8LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UInt8LiteralNode", Ref(UInt8LiteralNodeClass))

	FloatLiteralNodeClass = NewClass()
	FloatLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	FloatLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("FloatLiteralNode", Ref(FloatLiteralNodeClass))

	BigFloatLiteralNodeClass = NewClass()
	BigFloatLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	BigFloatLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("BigFloatLiteralNode", Ref(BigFloatLiteralNodeClass))

	Float64LiteralNodeClass = NewClass()
	Float64LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Float64LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("Float64LiteralNode", Ref(Float64LiteralNodeClass))

	Float32LiteralNodeClass = NewClass()
	Float32LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	Float32LiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("Float32LiteralNode", Ref(Float32LiteralNodeClass))

	UninterpolatedRegexLiteralNodeClass = NewClass()
	UninterpolatedRegexLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UninterpolatedRegexLiteralNode", Ref(UninterpolatedRegexLiteralNodeClass))

	RegexLiteralContentSectionNodeClass = NewClass()
	ElkASTModule.AddConstantString("RegexLiteralContentSectionNode", Ref(RegexLiteralContentSectionNodeClass))

	RegexInterpolationNodeClass = NewClass()
	ElkASTModule.AddConstantString("RegexInterpolationNodeClass", Ref(RegexInterpolationNodeClass))

	InterpolatedRegexLiteralNodeClass = NewClass()
	InterpolatedRegexLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InterpolatedRegexLiteralNode", Ref(InterpolatedRegexLiteralNodeClass))

	CharLiteralNodeClass = NewClass()
	CharLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	CharLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("CharLiteralNode", Ref(CharLiteralNodeClass))

	RawCharLiteralNodeClass = NewClass()
	RawCharLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	RawCharLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("RawCharLiteralNode", Ref(RawCharLiteralNodeClass))

	RawStringLiteralNodeClass = NewClass()
	RawStringLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	RawStringLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("RawStringLiteralNode", Ref(RawStringLiteralNodeClass))

	StringLiteralContentSectionNodeClass = NewClass()
	ElkASTModule.AddConstantString("StringLiteralContentSectionNode", Ref(StringLiteralContentSectionNodeClass))

	StringInspectInterpolationNodeClass = NewClass()
	ElkASTModule.AddConstantString("StringInspectInterpolationNode", Ref(StringInspectInterpolationNodeClass))

	StringInterpolationNodeClass = NewClass()
	ElkASTModule.AddConstantString("StringInterpolationNode", Ref(StringInterpolationNodeClass))

	InterpolatedStringLiteralNodeClass = NewClass()
	InterpolatedStringLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	InterpolatedStringLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("InterpolatedStringLiteralNode", Ref(InterpolatedStringLiteralNodeClass))

	DoubleQuotedStringLiteralNodeClass = NewClass()
	DoubleQuotedStringLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	DoubleQuotedStringLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("DoubleQuotedStringLiteralNode", Ref(DoubleQuotedStringLiteralNodeClass))

	InterpolatedSymbolLiteralNodeClass = NewClass()
	InterpolatedSymbolLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	InterpolatedSymbolLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("InterpolatedSymbolLiteralNode", Ref(InterpolatedSymbolLiteralNodeClass))

	ConstantAsNodeClass = NewClass()
	ConstantAsNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ConstantAsNode", Ref(ConstantAsNodeClass))

	MethodLookupAsNodeClass = NewClass()
	MethodLookupAsNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MethodLookupAsNode", Ref(MethodLookupAsNodeClass))

	PublicIdentifierNodeClass = NewClass()
	PublicIdentifierNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("PublicIdentifierNode", Ref(PublicIdentifierNodeClass))

	PublicIdentifierNodeClass = NewClass()
	PublicIdentifierNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("PublicIdentifierNode", Ref(PublicIdentifierNodeClass))

	PublicIdentifierAsNodeClass = NewClass()
	PublicIdentifierAsNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("PublicIdentifierAsNode", Ref(PublicIdentifierAsNodeClass))

	PrivateIdentifierNodeClass = NewClass()
	PrivateIdentifierNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("PrivateIdentifierNode", Ref(PrivateIdentifierNodeClass))

	PublicConstantNodeClass = NewClass()
	PublicConstantNodeClass.IncludeMixin(ExpressionNodeMixin)
	PublicConstantNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("PublicConstantNode", Ref(PublicConstantNodeClass))

	PrivateConstantNodeClass = NewClass()
	PrivateConstantNodeClass.IncludeMixin(ExpressionNodeMixin)
	PrivateConstantNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("PrivateConstantNode", Ref(PrivateConstantNodeClass))

	AsExpressionNodeClass = NewClass()
	AsExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AsExpressionNode", Ref(AsExpressionNodeClass))

	DoExpressionNodeClass = NewClass()
	DoExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("DoExpressionNode", Ref(DoExpressionNodeClass))

	SingletonBlockExpressionNodeClass = NewClass()
	SingletonBlockExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SingletonBlockExpressionNode", Ref(SingletonBlockExpressionNodeClass))

	SwitchExpressionNodeClass = NewClass()
	SwitchExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SwitchExpressionNode", Ref(SwitchExpressionNodeClass))

	IfExpressionNodeClass = NewClass()
	IfExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("IfExpressionNode", Ref(IfExpressionNodeClass))

	UnlessExpressionNodeClass = NewClass()
	UnlessExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UnlessExpressionNode", Ref(UnlessExpressionNodeClass))

	WhileExpressionNodeClass = NewClass()
	WhileExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("WhileExpressionNode", Ref(WhileExpressionNodeClass))

	UntilExpressionNodeClass = NewClass()
	UntilExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UntilExpressionNode", Ref(UntilExpressionNodeClass))

	LoopExpressionNodeClass = NewClass()
	LoopExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("LoopExpressionNode", Ref(LoopExpressionNodeClass))

	ForInExpressionNodeClass = NewClass()
	ForInExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ForInExpressionNode", Ref(ForInExpressionNodeClass))

	BreakExpressionNodeClass = NewClass()
	BreakExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("BreakExpressionNode", Ref(BreakExpressionNodeClass))

	NumericForExpressionNodeClass = NewClass()
	NumericForExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("NumericForExpressionNode", Ref(NumericForExpressionNodeClass))

	LabeledExpressionNodeClass = NewClass()
	LabeledExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("LabeledExpressionNode", Ref(LabeledExpressionNodeClass))

	GoExpressionNodeClass = NewClass()
	GoExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GoExpressionNode", Ref(GoExpressionNodeClass))

	ReturnExpressionNodeClass = NewClass()
	ReturnExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ReturnExpressionNode", Ref(ReturnExpressionNodeClass))

	YieldExpressionNodeClass = NewClass()
	YieldExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("YieldExpressionNode", Ref(YieldExpressionNodeClass))

	ContinueExpressionNodeClass = NewClass()
	ContinueExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ContinueExpressionNode", Ref(ContinueExpressionNodeClass))

	ThrowExpressionNodeClass = NewClass()
	ThrowExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ThrowExpressionNode", Ref(ThrowExpressionNodeClass))

	MustExpressionNodeClass = NewClass()
	MustExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MustExpressionNode", Ref(MustExpressionNodeClass))

	TryExpressionNodeClass = NewClass()
	TryExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("TryExpressionNode", Ref(TryExpressionNodeClass))

	AwaitExpressionNodeClass = NewClass()
	AwaitExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AwaitExpressionNode", Ref(AwaitExpressionNodeClass))

	TypeofExpressionNodeClass = NewClass()
	TypeofExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("TypeofExpressionNode", Ref(TypeofExpressionNodeClass))

	ConstantLookupNodeClass = NewClass()
	ConstantLookupNodeClass.IncludeMixin(ExpressionNodeMixin)
	ConstantLookupNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("ConstantLookupNode", Ref(ConstantLookupNodeClass))

	MethodLookupNodeClass = NewClass()
	MethodLookupNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MethodLookupNode", Ref(MethodLookupNodeClass))

	UsingEntryWithSubentriesNodeClass = NewClass()
	UsingEntryWithSubentriesNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UsingEntryWithSubentriesNode", Ref(UsingEntryWithSubentriesNodeClass))

	UsingAllEntryNodeClass = NewClass()
	UsingAllEntryNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UsingAllEntryNode", Ref(UsingAllEntryNodeClass))

	ClosureLiteralNodeClass = NewClass()
	ClosureLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ClosureLiteralNode", Ref(ClosureLiteralNodeClass))

	ClassDeclarationNodeClass = NewClass()
	ClassDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ClassDeclarationNode", Ref(ClassDeclarationNodeClass))

	ModuleDeclarationNodeClass = NewClass()
	ModuleDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ModuleDeclarationNode", Ref(ModuleDeclarationNodeClass))

	MixinDeclarationNodeClass = NewClass()
	MixinDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MixinDeclarationNode", Ref(MixinDeclarationNodeClass))

	InterfaceDeclarationNodeClass = NewClass()
	InterfaceDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InterfaceDeclarationNode", Ref(InterfaceDeclarationNodeClass))

	StructDeclarationNodeClass = NewClass()
	StructDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("StructDeclarationNode", Ref(StructDeclarationNodeClass))

	MethodDefinitionNodeClass = NewClass()
	MethodDefinitionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MethodDefinitionNode", Ref(MethodDefinitionNodeClass))

	InitDefinitionNodeClass = NewClass()
	InitDefinitionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InitDefinitionNode", Ref(InitDefinitionNodeClass))

	MethodSignatureDefinitionNodeClass = NewClass()
	MethodSignatureDefinitionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MethodSignatureDefinitionNode", Ref(MethodSignatureDefinitionNodeClass))

	GenericConstantNodeClass = NewClass()
	GenericConstantNodeClass.IncludeMixin(ExpressionNodeMixin)
	GenericConstantNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("GenericConstantNode", Ref(GenericConstantNodeClass))

	GenericTypeDefinitionNodeClass = NewClass()
	GenericTypeDefinitionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GenericTypeDefinitionNode", Ref(GenericTypeDefinitionNodeClass))

	TypeDefinitionNodeClass = NewClass()
	TypeDefinitionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("TypeDefinitionNode", Ref(TypeDefinitionNodeClass))

	AliasDeclarationEntryClass = NewClass()
	ElkASTModule.AddConstantString("AliasDeclarationEntry", Ref(AliasDeclarationEntryClass))

	AliasDeclarationNodeClass = NewClass()
	AliasDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AliasDeclarationNode", Ref(AliasDeclarationNodeClass))

	GetterDeclarationNodeClass = NewClass()
	GetterDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GetterDeclarationNode", Ref(GetterDeclarationNodeClass))

	SetterDeclarationNodeClass = NewClass()
	SetterDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SetterDeclarationNode", Ref(SetterDeclarationNodeClass))

	AttrDeclarationNodeClass = NewClass()
	AttrDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AttrDeclarationNode", Ref(AttrDeclarationNodeClass))

	UsingExpressionNodeClass = NewClass()
	UsingExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UsingExpressionNode", Ref(UsingExpressionNodeClass))

	IncludeExpressionNodeClass = NewClass()
	IncludeExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("IncludeExpressionNode", Ref(IncludeExpressionNodeClass))

	ExtendWhereBlockExpressionNodeClass = NewClass()
	ExtendWhereBlockExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ExtendWhereBlockExpressionNode", Ref(ExtendWhereBlockExpressionNodeClass))

	ImplementExpressionNodeClass = NewClass()
	ImplementExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ImplementExpressionNode", Ref(ImplementExpressionNodeClass))

	NewExpressionNodeClass = NewClass()
	NewExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("NewExpressionNode", Ref(NewExpressionNodeClass))

	GenericConstructorCallNodeClass = NewClass()
	GenericConstructorCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GenericConstructorCallNode", Ref(GenericConstructorCallNodeClass))

	ConstructorCallNodeClass = NewClass()
	ConstructorCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ConstructorCallNode", Ref(ConstructorCallNodeClass))

	NamedCallArgumentNodeClass = NewClass()
	ElkASTModule.AddConstantString("NamedCallArgumentNode", Ref(NamedCallArgumentNodeClass))

	DoubleSplatExpressionNodeClass = NewClass()
	DoubleSplatExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("DoubleSplatExpressionNode", Ref(DoubleSplatExpressionNodeClass))

	AttributeAccessNodeClass = NewClass()
	AttributeAccessNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("AttributeAccessNode", Ref(AttributeAccessNodeClass))

	SubscriptExpressionNodeClass = NewClass()
	SubscriptExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SubscriptExpressionNode", Ref(SubscriptExpressionNodeClass))

	NilSafeSubscriptExpressionNodeClass = NewClass()
	NilSafeSubscriptExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("NilSafeSubscriptExpressionNode", Ref(NilSafeSubscriptExpressionNodeClass))

	CallNodeClass = NewClass()
	CallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("CallNode", Ref(CallNodeClass))

	GenericMethodCallNodeClass = NewClass()
	GenericMethodCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GenericMethodCallNode", Ref(GenericMethodCallNodeClass))

	MethodCallNodeClass = NewClass()
	MethodCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("MethodCallNode", Ref(MethodCallNodeClass))

	ReceiverlessMethodCallNodeClass = NewClass()
	ReceiverlessMethodCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ReceiverlessMethodCallNode", Ref(ReceiverlessMethodCallNodeClass))

	GenericReceiverlessMethodCallNodeClass = NewClass()
	GenericReceiverlessMethodCallNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("GenericReceiverlessMethodCallNode", Ref(GenericReceiverlessMethodCallNodeClass))

	SplatExpressionNodeClass = NewClass()
	SplatExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SplatExpressionNode", Ref(SplatExpressionNodeClass))

	KeyValueExpressionNodeClass = NewClass()
	KeyValueExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("KeyValueExpressionNode", Ref(KeyValueExpressionNodeClass))

	SymbolKeyValueExpressionNodeClass = NewClass()
	SymbolKeyValueExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SymbolKeyValueExpressionNode", Ref(SymbolKeyValueExpressionNodeClass))

	WordArrayListLiteralNodeClass = NewClass()
	WordArrayListLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("WordArrayListLiteralNode", Ref(WordArrayListLiteralNodeClass))

	ArrayListLiteralNodeClass = NewClass()
	ArrayListLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ArrayListLiteralNode", Ref(ArrayListLiteralNodeClass))

	SymbolArrayListLiteralNodeClass = NewClass()
	SymbolArrayListLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SymbolArrayListLiteralNode", Ref(SymbolArrayListLiteralNodeClass))

	HexArrayListLiteralNodeClass = NewClass()
	HexArrayListLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HexArrayListLiteralNode", Ref(HexArrayListLiteralNodeClass))

	BinArrayListLiteralNodeClass = NewClass()
	BinArrayListLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("BinArrayListLiteralNode", Ref(BinArrayListLiteralNodeClass))

	ArrayTupleLiteralNodeClass = NewClass()
	ArrayTupleLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("ArrayTupleLiteralNode", Ref(ArrayTupleLiteralNodeClass))

	WordArrayTupleLiteralNodeClass = NewClass()
	WordArrayTupleLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("WordArrayTupleLiteralNode", Ref(WordArrayTupleLiteralNodeClass))

	SymbolArrayTupleLiteralNodeClass = NewClass()
	SymbolArrayTupleLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SymbolArrayTupleLiteralNode", Ref(SymbolArrayTupleLiteralNodeClass))

	HexArrayTupleLiteralNodeClass = NewClass()
	HexArrayTupleLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HexArrayTupleLiteralNode", Ref(HexArrayTupleLiteralNodeClass))

	BinArrayTupleLiteralNodeClass = NewClass()
	BinArrayTupleLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("BinArrayTupleLiteralNode", Ref(BinArrayTupleLiteralNodeClass))

	HashSetLiteralNodeClass = NewClass()
	HashSetLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HashSetLiteralNode", Ref(HashSetLiteralNodeClass))

	WordHashSetLiteralNodeClass = NewClass()
	WordHashSetLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("WordHashSetLiteralNode", Ref(WordHashSetLiteralNodeClass))

	SymbolHashSetLiteralNodeClass = NewClass()
	SymbolHashSetLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SymbolHashSetLiteralNode", Ref(SymbolHashSetLiteralNodeClass))

	HexHashSetLiteralNodeClass = NewClass()
	HexHashSetLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HexHashSetLiteralNode", Ref(HexHashSetLiteralNodeClass))

	BinHashSetLiteralNodeClass = NewClass()
	BinHashSetLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("BinHashSetLiteralNode", Ref(BinHashSetLiteralNodeClass))

	HashMapLiteralNodeClass = NewClass()
	HashMapLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HashMapLiteralNode", Ref(HashMapLiteralNodeClass))

	HashRecordLiteralNodeClass = NewClass()
	HashRecordLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("HashRecordLiteralNode", Ref(HashRecordLiteralNodeClass))

	RangeLiteralNodeClass = NewClass()
	RangeLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("RangeLiteralNode", Ref(RangeLiteralNodeClass))

	VariantTypeParameterNodeClass = NewClass()
	ElkASTModule.AddConstantString("VariantTypeParameterNode", Ref(VariantTypeParameterNodeClass))

	FormalParameterNodeClass = NewClass()
	FormalParameterNodeClass.IncludeMixin(ParameterNodeMixin)
	ElkASTModule.AddConstantString("FormalParameterNode", Ref(FormalParameterNodeClass))

	MethodParameterNodeClass = NewClass()
	MethodParameterNodeClass.IncludeMixin(ParameterNodeMixin)
	ElkASTModule.AddConstantString("MethodParameterNode", Ref(MethodParameterNodeClass))

	SignatureParameterNodeClass = NewClass()
	SignatureParameterNodeClass.IncludeMixin(ParameterNodeMixin)
	ElkASTModule.AddConstantString("SignatureParameterNode", Ref(SignatureParameterNodeClass))

	AttributeParameterNodeClass = NewClass()
	AttributeParameterNodeClass.IncludeMixin(ParameterNodeMixin)
	ElkASTModule.AddConstantString("AttributeParameterNode", Ref(AttributeParameterNodeClass))

	BoolLiteralNodeClass = NewClass()
	BoolLiteralNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("BoolLiteralNode", Ref(BoolLiteralNodeClass))

	VoidTypeNodeClass = NewClass()
	VoidTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("VoidTypeNode", Ref(VoidTypeNodeClass))

	NeverTypeNodeClass = NewClass()
	NeverTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("NeverTypeNode", Ref(NeverTypeNodeClass))

	AnyTypeNodeClass = NewClass()
	AnyTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("AnyTypeNode", Ref(AnyTypeNodeClass))

	UnionTypeNodeClass = NewClass()
	UnionTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UnionTypeNode", Ref(UnionTypeNodeClass))

	IntersectionTypeNodeClass = NewClass()
	IntersectionTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("IntersectionTypeNode", Ref(IntersectionTypeNodeClass))

	BinaryTypeNodeClass = NewClass()
	BinaryTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("BinaryTypeExpressionNode", Ref(BinaryTypeNodeClass))

	NilableTypeNodeClass = NewClass()
	NilableTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("NilableTypeNode", Ref(NilableTypeNodeClass))

	InstanceOfTypeNodeClass = NewClass()
	InstanceOfTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("InstanceOfTypeNode", Ref(InstanceOfTypeNodeClass))

	SingletonTypeNodeClass = NewClass()
	SingletonTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("SingletonTypeNode", Ref(SingletonTypeNodeClass))

	NotTypeNodeClass = NewClass()
	NotTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("NotTypeNode", Ref(NotTypeNodeClass))

	ClosureTypeNodeClass = NewClass()
	ClosureTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("ClosureTypeNode", Ref(ClosureTypeNodeClass))

	UnaryTypeNodeClass = NewClass()
	UnaryTypeNodeClass.IncludeMixin(TypeNodeMixin)
	ElkASTModule.AddConstantString("UnaryTypeNode", Ref(UnaryTypeNodeClass))
}
