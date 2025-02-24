package value

var ASTModule *Module // Std::AST
var NodeClass *Class  // Std::AST::Node

var StatementNodeInterface *Interface   // Std::AST::StatementNode
var ExpressionStatementNodeClass *Class // Std::AST::ExpressionStatementNode
var EmptyStatementNodeClass *Class      // Std::AST::EmptyStatementNode
var ImportStatementNodeClass *Class     // Std::AST::ImportStatementNode

var ProgramNodeClass *Class                     // Std::AST::ProgramNode
var ExpressionNodeMixin *Mixin                  // Std::AST::ExpressionNode
var InvalidNodeClass *Class                     // Std::AST::InvalidNode
var TypeExpressionNodeClass *Class              // Std::AST::TypeExpressionNode
var InstanceVariableDeclarationNodeClass *Class // Std::AST::InstanceVariableDeclarationNode
var VariablePatternDeclarationNodeClass *Class  // Std::AST::VariablePatternDeclarationNode
var VariableDeclarationNodeClass *Class         // Std::AST::VariableDeclarationNode
var ValuePatternDeclarationNodeClass *Class     // Std::AST::ValuePatternDeclarationNode
var ValueDeclarationNodeClass *Class            // Std::AST::ValueDeclarationNode
var PostfixExpressionNodeClass *Class           // Std::AST::PostfixExpressionNode
var ModifierNodeClass *Class                    // Std::AST::ModifierNode
var ModifierIfElseNodeClass *Class              // Std::AST::ModifierIfElseNode
var ModifierForInNodeClass *Class               // Std::AST::ModifierForInNode
var AssignmentExpressionNodeClass *Class        // Std::AST::AssignmentExpressionNode
var BinaryExpressionNodeClass *Class            // Std::AST::BinaryExpressionNode
var LogicalExpressionNodeClass *Class           // Std::AST::LogicalExpressionNode
var UnaryExpressionNodeClass *Class             // Std::AST::UnaryExpressionNode
var TrueLiteralNodeClass *Class                 // Std::AST::TrueLiteralNode
var FalseLiteralNodeClass *Class                // Std::AST::FalseLiteralNode
var NilLiteralNodeClass *Class                  // Std::AST::NilLiteralNode
var UndefinedLiteralNodeClass *Class            // Std::AST::UndefinedLiteralNode
var SelfLiteralNodeClass *Class                 // Std::AST::SelfLiteralNode

func initAST() {
	ASTModule = NewModule()
	StdModule.AddConstantString("AST", Ref(ASTModule))

	NodeClass = NewClass()
	ASTModule.AddConstantString("Node", Ref(NodeClass))

	StatementNodeInterface = NewInterface()
	ASTModule.AddConstantString("StatementNode", Ref(StatementNodeInterface))

	ExpressionStatementNodeClass = NewClass()
	ASTModule.AddConstantString("ExpressionStatementNode", Ref(ExpressionStatementNodeClass))

	EmptyStatementNodeClass = NewClass()
	ASTModule.AddConstantString("EmptyStatementNode", Ref(EmptyStatementNodeClass))

	ImportStatementNodeClass = NewClass()
	ASTModule.AddConstantString("ImportStatementNode", Ref(ImportStatementNodeClass))

	ProgramNodeClass = NewClass()
	ASTModule.AddConstantString("ProgramNode", Ref(ProgramNodeClass))

	ExpressionNodeMixin = NewMixin()
	ASTModule.AddConstantString("ExpressionNode", Ref(ExpressionNodeMixin))

	InvalidNodeClass = NewClass()
	InvalidNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("InvalidNode", Ref(InvalidNodeClass))

	TypeExpressionNodeClass = NewClass()
	TypeExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("TypeExpressionNode", Ref(TypeExpressionNodeClass))

	InstanceVariableDeclarationNodeClass = NewClass()
	InstanceVariableDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("InstanceVariableDeclarationNode", Ref(InstanceVariableDeclarationNodeClass))

	VariablePatternDeclarationNodeClass = NewClass()
	VariablePatternDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("VariablePatternDeclarationNode", Ref(VariablePatternDeclarationNodeClass))

	VariableDeclarationNodeClass = NewClass()
	VariableDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("VariableDeclarationNode", Ref(VariableDeclarationNodeClass))

	ValuePatternDeclarationNodeClass = NewClass()
	ValuePatternDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("ValuePatternDeclarationNode", Ref(ValuePatternDeclarationNodeClass))

	ValueDeclarationNodeClass = NewClass()
	ValueDeclarationNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("ValueDeclarationNode", Ref(ValueDeclarationNodeClass))

	PostfixExpressionNodeClass = NewClass()
	PostfixExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("PostfixExpressionNode", Ref(PostfixExpressionNodeClass))

	ModifierNodeClass = NewClass()
	ModifierNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("ModifierNode", Ref(ModifierNodeClass))

	ModifierIfElseNodeClass = NewClass()
	ModifierIfElseNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("ModifierIfElseNode", Ref(ModifierIfElseNodeClass))

	ModifierForInNodeClass = NewClass()
	ModifierForInNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("ModifierForInNode", Ref(ModifierForInNodeClass))

	AssignmentExpressionNodeClass = NewClass()
	AssignmentExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("AssignmentExpressionNode", Ref(AssignmentExpressionNodeClass))

	BinaryExpressionNodeClass = NewClass()
	BinaryExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("BinaryExpressionNode", Ref(BinaryExpressionNodeClass))

	LogicalExpressionNodeClass = NewClass()
	LogicalExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("LogicalExpressionNode", Ref(LogicalExpressionNodeClass))

	UnaryExpressionNodeClass = NewClass()
	UnaryExpressionNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("UnaryExpressionNode", Ref(UnaryExpressionNodeClass))

	TrueLiteralNodeClass = NewClass()
	TrueLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("TrueLiteralNode", Ref(TrueLiteralNodeClass))

	FalseLiteralNodeClass = NewClass()
	FalseLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("FalseLiteralNode", Ref(FalseLiteralNodeClass))

	NilLiteralNodeClass = NewClass()
	NilLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("NilLiteralNode", Ref(NilLiteralNodeClass))

	UndefinedLiteralNodeClass = NewClass()
	UndefinedLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("UndefinedLiteralNode", Ref(UndefinedLiteralNodeClass))

	SelfLiteralNodeClass = NewClass()
	SelfLiteralNodeClass.IncludeMixin(ExpressionNodeMixin)
	ASTModule.AddConstantString("SelfLiteralNode", Ref(SelfLiteralNodeClass))
}
