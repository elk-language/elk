package ast

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
