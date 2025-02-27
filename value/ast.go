package value

var ElkASTModule *Module // Std::ElkAST
var NodeClass *Class     // Std::ElkAST::Node

var StatementNodeInterface *Interface   // Std::ElkAST::StatementNode
var ExpressionStatementNodeClass *Class // Std::ElkAST::ExpressionStatementNode
var EmptyStatementNodeClass *Class      // Std::ElkAST::EmptyStatementNode
var ImportStatementNodeClass *Class     // Std::ElkAST::ImportStatementNode

var ProgramNodeClass *Class                     // Std::ElkAST::ProgramNode
var ExpressionNodeMixin *Mixin                  // Std::ElkAST::ExpressionNode
var InvalidNodeClass *Class                     // Std::ElkAST::InvalidNode
var TypeExpressionNodeClass *Class              // Std::ElkAST::TypeExpressionNode
var InstanceVariableDeclarationNodeClass *Class // Std::ElkAST::InstanceVariableDeclarationNode
var VariablePatternDeclarationNodeClass *Class  // Std::ElkAST::VariablePatternDeclarationNode
var VariableDeclarationNodeClass *Class         // Std::ElkAST::VariableDeclarationNode
var ValuePatternDeclarationNodeClass *Class     // Std::ElkAST::ValuePatternDeclarationNode
var ValueDeclarationNodeClass *Class            // Std::ElkAST::ValueDeclarationNode
var PostfixExpressionNodeClass *Class           // Std::ElkAST::PostfixExpressionNode
var ModifierNodeClass *Class                    // Std::ElkAST::ModifierNode
var ModifierIfElseNodeClass *Class              // Std::ElkAST::ModifierIfElseNode
var ModifierForInNodeClass *Class               // Std::ElkAST::ModifierForInNode
var AssignmentExpressionNodeClass *Class        // Std::ElkAST::AssignmentExpressionNode
var BinaryExpressionNodeClass *Class            // Std::ElkAST::BinaryExpressionNode
var LogicalExpressionNodeClass *Class           // Std::ElkAST::LogicalExpressionNode
var UnaryExpressionNodeClass *Class             // Std::ElkAST::UnaryExpressionNode
var TrueLiteralNodeClass *Class                 // Std::ElkAST::TrueLiteralNode
var FalseLiteralNodeClass *Class                // Std::ElkAST::FalseLiteralNode
var NilLiteralNodeClass *Class                  // Std::ElkAST::NilLiteralNode
var UndefinedLiteralNodeClass *Class            // Std::ElkAST::UndefinedLiteralNode
var SelfLiteralNodeClass *Class                 // Std::ElkAST::SelfLiteralNode
var InstanceVariableNodeClass *Class            // Std::ElkAST::InstanceVariableNode
var SimpleSymbolLiteralNodeClass *Class         // Std::ElkAST::SimpleSymbolLiteralNode
var IntLiteralNodeClass *Class                  // Std::ElkAST::IntLiteralNode
var Int64LiteralNodeClass *Class                // Std::ElkAST::Int64LiteralNode
var Int32LiteralNodeClass *Class                // Std::ElkAST::Int32LiteralNode
var Int16LiteralNodeClass *Class                // Std::ElkAST::Int16LiteralNode
var Int8LiteralNodeClass *Class                 // Std::ElkAST::Int8LiteralNode
var UInt64LiteralNodeClass *Class               // Std::ElkAST::UInt64LiteralNode
var UInt32LiteralNodeClass *Class               // Std::ElkAST::UInt32LiteralNode
var UInt16LiteralNodeClass *Class               // Std::ElkAST::UInt16LiteralNode
var UInt8LiteralNodeClass *Class                // Std::ElkAST::UInt8LiteralNode
var FloatLiteralNodeClass *Class                // Std::ElkAST::FloatLiteralNode
var BigFloatLiteralNodeClass *Class             // Std::ElkAST::BigFloatLiteralNode
var Float32LiteralNodeClass *Class              // Std::ElkAST::Float32LiteralNode
var Float64LiteralNodeClass *Class              // Std::ElkAST::Float64LiteralNode
var UninterpolatedRegexLiteralNodeClass *Class  // Std::ElkAST::UninterpolatedRegexLiteralNode
var RegexLiteralContentSectionNodeClass *Class  // Std::ElkAST::RegexLiteralContentSection
var RegexInterpolationNodeClass *Class          // Std::ElkAST::RegexInterpolationNode
var InterpolatedRegexLiteralNodeClass *Class    // Std::ElkAST::InterpolatedRegexLiteralNode
var CharLiteralNodeClass *Class                 // Std::ElkAST::CharLiteralNode
var RawCharLiteralNodeClass *Class              // Std::ElkAST::RawCharLiteralNode
var RawStringLiteralNodeClass *Class            // Std::ElkAST::RawStringLiteralNode
var StringLiteralContentSectionNodeClass *Class // Std::ElkAST::StringLiteralContentSectionNode
var StringInspectInterpolationNodeClass *Class  // Std::ElkAST::StringInspectInterpolationNode
var StringInterpolationNodeClass *Class         // Std::ElkAST::StringInterpolationNode
var InterpolatedStringLiteralNodeClass *Class   // Std::ElkAST::InterpolatedStringLiteralNode
var DoubleQuotedStringLiteralNodeClass *Class   // Std::ElkAST::DoubleQuotedStringLiteralNode
var ConstantAsNodeClass *Class                  // Std::ElkAST::ConstantAsNode
var MethodLookupAsNodeClass *Class              // Std::ElkAST::MethodLookupAsNode
var PublicIdentifierNodeClass *Class            // Std::ElkAST::PublicIdentifierNode
var PublicIdentifierAsNodeClass *Class          // Std::ElkAST::PublicIdentifierAsNode
var PrivateIdentifierNodeClass *Class           // Std::ElkAST::PrivateIdentifierNode
var PublicConstantNodeClass *Class              // Std::ElkAST::PublicConstantNode
var PrivateConstantNodeClass *Class             // Std::ElkAST::PrivateConstantNode
var AsExpressionNodeClass *Class                // Std::ElkAST::AsExpressionNode
var DoExpressionNodeClass *Class                // Std::ElkAST::DoExpressionNode
var SingletonBlockExpressionNodeClass *Class    // Std::ElkAST::SingletonBlockExpressionNode
var SwitchExpressionNodeClass *Class            // Std::ElkAST::SwitchExpressionNode
var IfExpressionNodeClass *Class                // Std::ElkAST::IfExpressionNode
var UnlessExpressionNodeClass *Class            // Std::ElkAST::UnlessExpressionNode
var WhileExpressionNodeClass *Class             // Std::ElkAST::WhileExpressionNode
var UntilExpressionNodeClass *Class             // Std::ElkAST::UntilExpressionNode
var LoopExpressionNodeClass *Class              // Std::ElkAST::LoopExpressionNode
var NumericForExpressionNodeClass *Class        // Std::ElkAST::NumericForExpressionNode
var ForInExpressionNodeClass *Class             // Std::ElkAST::ForInExpressionNode
var BreakExpressionNodeClass *Class             // Std::ElkAST::BreakExpressionNode
var LabeledExpressionNodeClass *Class           // Std::ElkAST::LabeledExpressionNode
var GoExpressionNodeClass *Class                // Std::ElkAST::GoExpressionNode
var ReturnExpressionNodeClass *Class            // Std::ElkAST::ReturnExpressionNode
var YieldExpressionNodeClass *Class             // Std::ElkAST::YieldExpressionNode
var ContinueExpressionNodeClass *Class          // Std::ElkAST::ContinueExpressionNode
var ThrowExpressionNodeClass *Class             // Std::ElkAST::ThrowExpressionNode
var MustExpressionNodeClass *Class              // Std::ElkAST::MustExpressionNode
var TryExpressionNodeClass *Class               // Std::ElkAST::TryExpressionNode
var AwaitExpressionNodeClass *Class             // Std::ElkAST::AwaitExpressionNode
var TypeofExpressionNodeClass *Class            // Std::ElkAST::TypeofExpressionNode
var ConstantLookupNodeClass *Class              // Std::ElkAST::ConstantLookupNode
var MethodLookupNodeClass *Class                // Std::ElkAST::MethodLookupNode
var UsingEntryWithSubentriesNodeClass *Class    // Std::ElkAST::UsingEntryWithSubentriesNode
var UsingAllEntryNodeClass *Class               // Std::ElkAST::UsingAllEntryNode
var ClosureLiteralNodeClass *Class              // Std::ElkAST::ClosureLiteralNode
var ClassDeclarationNodeClass *Class            // Std::ElkAST::ClassDeclarationNode
var ModuleDeclarationNodeClass *Class           // Std::ElkAST::ModuleDeclarationNode
var MixinDeclarationNodeClass *Class            // Std::ElkAST::MixinDeclarationNode
var InterfaceDeclarationNodeClass *Class        // Std::ElkAST::InterfaceDeclarationNode
var StructDeclarationNodeClass *Class           // Std::ElkAST::StructDeclarationNode
var MethodDefinitionNodeClass *Class            // Std::ElkAST::MethodDefinitionNode
var InitDefinitionNodeClass *Class              // Std::ElkAST::InitDefinitionNode
var MethodSignatureDefinitionNodeClass *Class   // Std::ElkAST::MethodSignatureDefinitionNode
var GenericConstantNodeClass *Class             // Std::ElkAST::GenericConstantNode
var GenericTypeDefinitionNodeClass *Class       // Std::ElkAST::GenericTypeDefinitionNode
var TypeDefinitionNodeClass *Class              // Std::ElkAST::TypeDefinitionNode
var AliasDeclarationEntryClass *Class           // Std::ElkAST::AliasDeclarationEntry
var AliasDeclarationNodeClass *Class            // Std::ElkAST::AliasDeclarationNode
var GetterDeclarationNodeClass *Class           // Std::ElkAST::GetterDeclarationNode
var SetterDeclarationNodeClass *Class           // Std::ElkAST::SetterDeclarationNode
var AttrDeclarationNodeClass *Class             // Std::ElkAST::AttrDeclarationNode
var UsingExpressionNodeClass *Class             // Std::ElkAST::UsingExpressionNode
var IncludeExpressionNodeClass *Class           // Std::ElkAST::IncludeExpressionNode
var ExtendWhereBlockExpressionNodeClass *Class  // Std::ElkAST::ExtendWhereBlockExpressionNode
var ImplementExpressionNodeClass *Class         // Std::ElkAST::ImplementExpressionNode
var NewExpressionNodeClass *Class               // Std::ElkAST::NewExpressionNode

func initAST() {
	ElkASTModule = NewModule()
	StdModule.AddConstantString("ElkAST", Ref(ElkASTModule))

	NodeClass = NewClass()
	ElkASTModule.AddConstantString("Node", Ref(NodeClass))

	StatementNodeInterface = NewInterface()
	ElkASTModule.AddConstantString("StatementNode", Ref(StatementNodeInterface))

	ExpressionStatementNodeClass = NewClass()
	ElkASTModule.AddConstantString("ExpressionStatementNode", Ref(ExpressionStatementNodeClass))

	EmptyStatementNodeClass = NewClass()
	ElkASTModule.AddConstantString("EmptyStatementNode", Ref(EmptyStatementNodeClass))

	ImportStatementNodeClass = NewClass()
	ElkASTModule.AddConstantString("ImportStatementNode", Ref(ImportStatementNodeClass))

	ProgramNodeClass = NewClass()
	ElkASTModule.AddConstantString("ProgramNode", Ref(ProgramNodeClass))

	ExpressionNodeMixin = NewMixin()
	ElkASTModule.AddConstantString("ExpressionNode", Ref(ExpressionNodeMixin))

	InvalidNodeClass = NewClass()
	InvalidNodeClass.IncludeMixin(ExpressionNodeMixin)
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
	ElkASTModule.AddConstantString("TrueLiteralNode", Ref(TrueLiteralNodeClass))

	FalseLiteralNodeClass = NewClass()
	FalseLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("FalseLiteralNode", Ref(FalseLiteralNodeClass))

	NilLiteralNodeClass = NewClass()
	NilLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("NilLiteralNode", Ref(NilLiteralNodeClass))

	UndefinedLiteralNodeClass = NewClass()
	UndefinedLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UndefinedLiteralNode", Ref(UndefinedLiteralNodeClass))

	SelfLiteralNodeClass = NewClass()
	SelfLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SelfLiteralNode", Ref(SelfLiteralNodeClass))

	InstanceVariableNodeClass = NewClass()
	InstanceVariableNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("InstanceVariableNode", Ref(InstanceVariableNodeClass))

	SimpleSymbolLiteralNodeClass = NewClass()
	SimpleSymbolLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("SimpleSymbolLiteralNode", Ref(SimpleSymbolLiteralNodeClass))

	IntLiteralNodeClass = NewClass()
	IntLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("IntLiteralNode", Ref(IntLiteralNodeClass))

	Int64LiteralNodeClass = NewClass()
	Int64LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("Int64LiteralNode", Ref(Int64LiteralNodeClass))

	Int32LiteralNodeClass = NewClass()
	Int32LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("Int32LiteralNode", Ref(Int32LiteralNodeClass))

	Int16LiteralNodeClass = NewClass()
	Int16LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("Int16LiteralNode", Ref(Int16LiteralNodeClass))

	Int8LiteralNodeClass = NewClass()
	Int8LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("Int8LiteralNode", Ref(Int8LiteralNodeClass))

	UInt64LiteralNodeClass = NewClass()
	UInt64LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UInt64LiteralNode", Ref(UInt64LiteralNodeClass))

	UInt32LiteralNodeClass = NewClass()
	UInt32LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UInt32LiteralNode", Ref(UInt32LiteralNodeClass))

	UInt16LiteralNodeClass = NewClass()
	UInt16LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UInt16LiteralNode", Ref(UInt16LiteralNodeClass))

	UInt8LiteralNodeClass = NewClass()
	UInt8LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("UInt8LiteralNode", Ref(UInt8LiteralNodeClass))

	FloatLiteralNodeClass = NewClass()
	FloatLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("FloatLiteralNode", Ref(FloatLiteralNodeClass))

	BigFloatLiteralNodeClass = NewClass()
	BigFloatLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("BigFloatLiteralNode", Ref(BigFloatLiteralNodeClass))

	Float64LiteralNodeClass = NewClass()
	Float64LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("Float64LiteralNode", Ref(Float64LiteralNodeClass))

	Float32LiteralNodeClass = NewClass()
	Float32LiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
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
	ElkASTModule.AddConstantString("CharLiteralNode", Ref(CharLiteralNodeClass))

	RawCharLiteralNodeClass = NewClass()
	RawCharLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("RawCharLiteralNode", Ref(RawCharLiteralNodeClass))

	RawStringLiteralNodeClass = NewClass()
	RawStringLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ElkASTModule.AddConstantString("RawStringLiteralNode", Ref(RawStringLiteralNodeClass))

	StringLiteralContentSectionNodeClass = NewClass()
	ElkASTModule.AddConstantString("StringLiteralContentSectionNode", Ref(StringLiteralContentSectionNodeClass))

	StringInspectInterpolationNodeClass = NewClass()
	ElkASTModule.AddConstantString("StringInspectInterpolationNode", Ref(StringInspectInterpolationNodeClass))

	StringInterpolationNodeClass = NewClass()
	ElkASTModule.AddConstantString("StringInterpolationNode", Ref(StringInterpolationNodeClass))

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
	ElkASTModule.AddConstantString("PublicConstantNode", Ref(PublicConstantNodeClass))

	PrivateConstantNodeClass = NewClass()
	PrivateConstantNodeClass.IncludeMixin(ExpressionNodeMixin)
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
}
