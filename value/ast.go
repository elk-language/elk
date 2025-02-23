package value

var ASTModule *Module // Std::AST
var NodeClass *Class  // Std::AST::Node

var StatementNodeInterface *Interface   // Std::AST::StatementNode
var ExpressionStatementNodeClass *Class // Std::AST::ExpressionStatementNode
var EmptyStatementNodeClass *Class      // Std::AST::EmptyStatementNode
var ImportStatementNodeClass *Class     // Std::AST::ImportStatementNode

var ProgramNodeClass *Class    // Std::AST::ProgramNode
var ExpressionNodeMixin *Mixin // Std::AST::ExpressionNode
var InvalidNodeClass *Class    // Std::AST::InvalidNode

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
	ASTModule.AddConstantString("InvalidNode", Ref(InvalidNodeClass))
}
