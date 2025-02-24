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

}
