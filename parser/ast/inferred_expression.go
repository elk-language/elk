package ast

import (
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
)

// An AST node that represents an expression (like the receiver in `.foo()`)
// that is absent and will be inferred.
type InferredExpressionNode struct {
	AbsentNodeBase
}

func IsInferredExpressionNode(node Node) bool {
	_, ok := node.(InferredExpressionNode)
	return ok
}

func (n InferredExpressionNode) splice(loc *position.Location, args *[]Node, unquote bool) Node {
	return InferredExpressionNode{}
}

func (n InferredExpressionNode) MacroType(env *types.GlobalEnvironment) types.Type {
	return types.NameToType("Std::Elk::AST::InferredExpressionNode", env)
}

func (n InferredExpressionNode) traverse(parent Node, enter func(node, parent Node) TraverseOption, leave func(node, parent Node) TraverseOption) TraverseOption {
	switch enter(n, parent) {
	case TraverseBreak:
		return TraverseBreak
	case TraverseSkip:
		return leave(n, parent)
	}

	return leave(n, parent)
}

func (n InferredExpressionNode) Equal(other value.Value) bool {
	_, ok := other.SafeAsReference().(InferredExpressionNode)
	return ok
}

func (n InferredExpressionNode) String() string {
	return "InferredExpressionNode{}"
}

func (InferredExpressionNode) IsStatic() bool {
	return true
}

func (InferredExpressionNode) Type(globalEnv *types.GlobalEnvironment) types.Type {
	return types.Untyped{}
}

func (InferredExpressionNode) Class() *value.Class {
	return value.TrueLiteralNodeClass
}

func (InferredExpressionNode) DirectClass() *value.Class {
	return value.TrueLiteralNodeClass
}

func (n InferredExpressionNode) Inspect() string {
	return "Std::Elk::AST::InferredExpressionNode{}"
}

func (n InferredExpressionNode) ToValue() value.Value {
	return value.Ref(n)
}

func (n InferredExpressionNode) Error() string {
	return n.Inspect()
}
