package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

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

func (*NewExpressionNode) Class() *value.Class {
	return value.NewExpressionNodeClass
}

func (*NewExpressionNode) DirectClass() *value.Class {
	return value.NewExpressionNodeClass
}

func (n *NewExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::NewExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  positional_arguments: %%[\n")
	for i, element := range n.PositionalArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  named_arguments: %%[\n")
	for i, element := range n.NamedArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NewExpressionNode) Error() string {
	return n.Inspect()
}

// Represents a constructor call eg. `ArrayList::[Int](1, 2, 3)`
type GenericConstructorCallNode struct {
	TypedNodeBase
	ClassNode           ComplexConstantNode // class that is being instantiated
	TypeArguments       []TypeNode
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

// Equal checks if the given GenericConstructorCallNode is equal to another value.
func (n *GenericConstructorCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*GenericConstructorCallNode)
	if !ok {
		return false
	}

	if !n.ClassNode.Equal(value.Ref(o.ClassNode)) {
		return false
	}

	if len(n.TypeArguments) != len(o.TypeArguments) ||
		len(n.NamedArguments) != len(o.NamedArguments) ||
		len(n.PositionalArguments) != len(o.PositionalArguments) {
		return false
	}

	for i, arg := range n.TypeArguments {
		if !arg.Equal(value.Ref(o.TypeArguments[i])) {
			return false
		}
	}

	for i, arg := range n.PositionalArguments {
		if !arg.Equal(value.Ref(o.PositionalArguments[i])) {
			return false
		}
	}

	for i, arg := range n.NamedArguments {
		if !arg.Equal(value.Ref(o.NamedArguments[i])) {
			return false
		}
	}

	return n.span.Equal(o.span)
}

// String returns a string representation of the GenericConstructorCallNode.
func (n *GenericConstructorCallNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.ClassNode.String())
	buff.WriteString("::[")

	for i, arg := range n.TypeArguments {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(arg.String())
	}

	buff.WriteString("](")

	var hasMultilineArgs bool
	argCount := len(n.PositionalArguments) + len(n.NamedArguments)
	argStrings := make([]string, 0, argCount)

	for _, arg := range n.PositionalArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}
	for _, arg := range n.NamedArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}

	if hasMultilineArgs || argCount > 6 {
		buff.WriteRune('\n')
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, argStr, 1)
		}
		buff.WriteRune('\n')
	} else {
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(argStr)
		}
	}

	buff.WriteString(")")

	return buff.String()
}

func (*GenericConstructorCallNode) IsStatic() bool {
	return false
}

func (*GenericConstructorCallNode) Class() *value.Class {
	return value.GenericConstructorCallNodeClass
}

func (*GenericConstructorCallNode) DirectClass() *value.Class {
	return value.GenericConstructorCallNodeClass
}

func (n *GenericConstructorCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::GenericConstructorCallNode{\n  &: %p", n)

	buff.WriteString(",\n  class_node: ")
	indent.IndentStringFromSecondLine(&buff, n.ClassNode.Inspect(), 1)

	buff.WriteString(",\n  type_arguments: %%[\n")
	for i, element := range n.TypeArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  positional_arguments: %%[\n")
	for i, element := range n.PositionalArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  named_arguments: %%[\n")
	for i, element := range n.NamedArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *GenericConstructorCallNode) Error() string {
	return n.Inspect()
}

// Create a constructor call node eg. `ArrayList::[Int](1, 2, 3)`
func NewGenericConstructorCallNode(span *position.Span, class ComplexConstantNode, typeArgs []TypeNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *GenericConstructorCallNode {
	return &GenericConstructorCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		ClassNode:           class,
		TypeArguments:       typeArgs,
		PositionalArguments: posArgs,
		NamedArguments:      namedArgs,
	}
}

// Represents a constructor call eg. `String(123)`
type ConstructorCallNode struct {
	TypedNodeBase
	ClassNode           ComplexConstantNode // class that is being instantiated
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

// Check if this node equals another node.
func (n *ConstructorCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*ConstructorCallNode)
	if !ok {
		return false
	}

	if !n.ClassNode.Equal(value.Ref(o.ClassNode)) {
		return false
	}

	if len(n.PositionalArguments) != len(o.PositionalArguments) ||
		len(n.NamedArguments) != len(o.NamedArguments) {
		return false
	}

	for i, arg := range n.PositionalArguments {
		if !arg.Equal(value.Ref(o.PositionalArguments[i])) {
			return false
		}
	}

	for i, arg := range n.NamedArguments {
		if !arg.Equal(value.Ref(o.NamedArguments[i])) {
			return false
		}
	}

	return n.span.Equal(o.span)
}

// Return a string representation of the node.
func (n *ConstructorCallNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.ClassNode.String())
	buff.WriteString("(")

	var hasMultilineArgs bool
	argCount := len(n.PositionalArguments) + len(n.NamedArguments)
	argStrings := make([]string, 0, argCount)

	for _, arg := range n.PositionalArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}
	for _, arg := range n.NamedArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}

	if hasMultilineArgs || argCount > 6 {
		buff.WriteRune('\n')
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, argStr, 1)
		}
		buff.WriteRune('\n')
	} else {
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(argStr)
		}
	}

	buff.WriteString(")")

	return buff.String()
}

func (*ConstructorCallNode) IsStatic() bool {
	return false
}

func (*ConstructorCallNode) Class() *value.Class {
	return value.ConstructorCallNodeClass
}

func (*ConstructorCallNode) DirectClass() *value.Class {
	return value.ConstructorCallNodeClass
}

func (n *ConstructorCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ConstructorCallNode{\n  &: %p", n)

	buff.WriteString(",\n  class_node: ")
	indent.IndentStringFromSecondLine(&buff, n.ClassNode.Inspect(), 1)

	buff.WriteString(",\n  positional_arguments: %%[\n")
	for i, element := range n.PositionalArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  named_arguments: %%[\n")
	for i, element := range n.NamedArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ConstructorCallNode) Error() string {
	return n.Inspect()
}

// Create a constructor call node eg. `String(123)`
func NewConstructorCallNode(span *position.Span, class ComplexConstantNode, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *ConstructorCallNode {
	return &ConstructorCallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		ClassNode:           class,
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

func (*AttributeAccessNode) Class() *value.Class {
	return value.AttributeAccessNodeClass
}

func (*AttributeAccessNode) DirectClass() *value.Class {
	return value.AttributeAccessNodeClass
}
func (n *AttributeAccessNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::AttributeAccessNode{\n  &: %p", n)

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  attribute_name: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.AttributeName).Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AttributeAccessNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*AttributeAccessNode)
	if !ok {
		return false
	}

	return n.Span().Equal(o.Span()) &&
		n.Receiver.Equal(value.Ref(o.Receiver)) &&
		n.AttributeName == o.AttributeName
}

func (n *AttributeAccessNode) String() string {
	var buff strings.Builder

	parens := ExpressionPrecedence(n) > ExpressionPrecedence(n.Receiver)
	if parens {
		buff.WriteRune('(')
	}
	buff.WriteString(n.Receiver.String())
	if parens {
		buff.WriteRune(')')
	}

	buff.WriteString(".")
	buff.WriteString(n.AttributeName)

	return buff.String()
}

func (n *AttributeAccessNode) Error() string {
	return n.Inspect()
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

func (*SubscriptExpressionNode) Class() *value.Class {
	return value.SubscriptExpressionNodeClass
}

func (*SubscriptExpressionNode) DirectClass() *value.Class {
	return value.SubscriptExpressionNodeClass
}

func (n *SubscriptExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SubscriptExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  key: ")
	indent.IndentStringFromSecondLine(&buff, n.Key.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SubscriptExpressionNode) Error() string {
	return n.Inspect()
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

func (*NilSafeSubscriptExpressionNode) Class() *value.Class {
	return value.SubscriptExpressionNodeClass
}

func (*NilSafeSubscriptExpressionNode) DirectClass() *value.Class {
	return value.SubscriptExpressionNodeClass
}

func (n *NilSafeSubscriptExpressionNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::NilSafeSubscriptExpressionNode{\n  &: %p", n)

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  key: ")
	indent.IndentStringFromSecondLine(&buff, n.Key.Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (n *NilSafeSubscriptExpressionNode) Error() string {
	return n.Inspect()
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
	TypedNodeBase
	Receiver            ExpressionNode
	NilSafe             bool
	PositionalArguments []ExpressionNode
	NamedArguments      []NamedArgumentNode
}

func (n *CallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*CallNode)
	if !ok {
		return false
	}

	if n.NilSafe != o.NilSafe {
		return false
	}

	if !n.Receiver.Equal(value.Ref(o.Receiver)) {
		return false
	}

	if len(n.PositionalArguments) != len(o.PositionalArguments) ||
		len(n.NamedArguments) != len(o.NamedArguments) {
		return false
	}

	for i, arg := range n.PositionalArguments {
		if !arg.Equal(value.Ref(o.PositionalArguments[i])) {
			return false
		}
	}

	for i, arg := range n.NamedArguments {
		if !arg.Equal(value.Ref(o.NamedArguments[i])) {
			return false
		}
	}

	return n.Span().Equal(o.Span())
}

func (n *CallNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Receiver.String())

	if n.NilSafe {
		buff.WriteRune('?')
	}
	buff.WriteString(".(")

	var hasMultilineArgs bool
	argCount := len(n.PositionalArguments) + len(n.NamedArguments)
	argStrings := make([]string, 0, argCount)

	for _, arg := range n.PositionalArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}
	for _, arg := range n.NamedArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}

	if hasMultilineArgs || argCount > 6 {
		buff.WriteRune('\n')
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, argStr, 1)
		}
		buff.WriteRune('\n')
	} else {
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(argStr)
		}
	}
	buff.WriteString(")")

	return buff.String()
}

func (*CallNode) IsStatic() bool {
	return false
}

func (*CallNode) Class() *value.Class {
	return value.CallNodeClass
}

func (*CallNode) DirectClass() *value.Class {
	return value.CallNodeClass
}

func (n *CallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::CallNode{\n  &: %p", n)

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	fmt.Fprintf(&buff, ",\n  nil_safe: %t", n.NilSafe)

	buff.WriteString(",\n  positional_arguments: %%[\n")
	for i, element := range n.PositionalArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  named_arguments: %%[\n")
	for i, element := range n.NamedArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *CallNode) Error() string {
	return n.Inspect()
}

// Create a method call node eg. `'123'.to_int()`
func NewCallNode(span *position.Span, recv ExpressionNode, nilSafe bool, posArgs []ExpressionNode, namedArgs []NamedArgumentNode) *CallNode {
	return &CallNode{
		TypedNodeBase:       TypedNodeBase{span: span},
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
	TailCall            bool
}

// Equal checks if the given GenericMethodCallNode is equal to another value.
func (n *GenericMethodCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*GenericMethodCallNode)
	if !ok {
		return false
	}

	if n.MethodName != o.MethodName ||
		!n.Receiver.Equal(value.Ref(o.Receiver)) ||
		!n.Op.Equal(o.Op) ||
		n.span.Equal(o.span) {
		return false
	}

	if len(n.TypeArguments) != len(o.TypeArguments) ||
		len(n.NamedArguments) != len(o.NamedArguments) ||
		len(n.PositionalArguments) != len(o.PositionalArguments) {
		return false
	}

	for i, arg := range n.TypeArguments {
		if !arg.Equal(value.Ref(o.TypeArguments[i])) {
			return false
		}
	}

	for i, arg := range n.PositionalArguments {
		if !arg.Equal(value.Ref(o.PositionalArguments[i])) {
			return false
		}
	}

	for i, arg := range n.NamedArguments {
		if !arg.Equal(value.Ref(o.NamedArguments[i])) {
			return false
		}
	}

	return true
}

// String returns a string representation of the GenericMethodCallNode.
func (n *GenericMethodCallNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.Receiver.String())
	buff.WriteString(n.Op.String())
	buff.WriteString(n.MethodName)
	buff.WriteString("::[")

	for i, arg := range n.TypeArguments {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(arg.String())
	}

	buff.WriteString("](")

	var hasMultilineArgs bool
	argCount := len(n.PositionalArguments) + len(n.NamedArguments)
	argStrings := make([]string, 0, argCount)

	for _, arg := range n.PositionalArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}
	for _, arg := range n.NamedArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}

	if hasMultilineArgs || argCount > 6 {
		buff.WriteRune('\n')
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, argStr, 1)
		}
		buff.WriteRune('\n')
	} else {
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(argStr)
		}
	}

	buff.WriteString(")")

	return buff.String()
}

func (*GenericMethodCallNode) IsStatic() bool {
	return false
}

func (*GenericMethodCallNode) Class() *value.Class {
	return value.GenericMethodCallNodeClass
}

func (*GenericMethodCallNode) DirectClass() *value.Class {
	return value.GenericMethodCallNodeClass
}

func (n *GenericMethodCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::GenericMethodCallNode{\n  &: %p", n)

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  method_name: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.MethodName).Inspect(), 1)

	buff.WriteString(",\n  type_arguments: %%[\n")
	for i, element := range n.TypeArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  positional_arguments: %%[\n")
	for i, element := range n.PositionalArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  named_arguments: %%[\n")
	for i, element := range n.NamedArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *GenericMethodCallNode) Error() string {
	return n.Inspect()
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
	TailCall            bool
}

func (*MethodCallNode) IsStatic() bool {
	return false
}

func (*MethodCallNode) Class() *value.Class {
	return value.MethodCallNodeClass
}

func (*MethodCallNode) DirectClass() *value.Class {
	return value.MethodCallNodeClass
}

func (n *MethodCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MethodCallNode{\n  &: %p", n)

	buff.WriteString(",\n  receiver: ")
	indent.IndentStringFromSecondLine(&buff, n.Receiver.Inspect(), 1)

	buff.WriteString(",\n  op: ")
	indent.IndentStringFromSecondLine(&buff, n.Op.Inspect(), 1)

	buff.WriteString(",\n  method_name: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.MethodName).Inspect(), 1)

	buff.WriteString(",\n  positional_arguments: %%[\n")
	for i, element := range n.PositionalArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  named_arguments: %%[\n")
	for i, element := range n.NamedArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MethodCallNode) Error() string {
	return n.Inspect()
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
	TailCall            bool
}

func (*ReceiverlessMethodCallNode) IsStatic() bool {
	return false
}

func (*ReceiverlessMethodCallNode) Class() *value.Class {
	return value.ReceiverlessMethodCallNodeClass
}

func (*ReceiverlessMethodCallNode) DirectClass() *value.Class {
	return value.ReceiverlessMethodCallNodeClass
}

func (n *ReceiverlessMethodCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::ReceiverlessMethodCallNode{\n  &: %p", n)

	buff.WriteString(",\n  method_name: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.MethodName).Inspect(), 1)

	buff.WriteString(",\n  positional_arguments: %%[\n")
	for i, element := range n.PositionalArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  named_arguments: %%[\n")
	for i, element := range n.NamedArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *ReceiverlessMethodCallNode) Error() string {
	return n.Inspect()
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
	TailCall            bool
}

// Equal checks if the given GenericReceiverlessMethodCallNode is equal to another value.
func (n *GenericReceiverlessMethodCallNode) Equal(other value.Value) bool {
	o, ok := other.SafeAsReference().(*GenericReceiverlessMethodCallNode)
	if !ok {
		return false
	}

	if n.MethodName != o.MethodName || !n.span.Equal(o.span) {
		return false
	}

	if len(n.TypeArguments) != len(o.TypeArguments) ||
		len(n.NamedArguments) != len(o.NamedArguments) ||
		len(n.PositionalArguments) != len(o.PositionalArguments) {
		return false
	}

	for i, arg := range n.TypeArguments {
		if !arg.Equal(value.Ref(o.TypeArguments[i])) {
			return false
		}
	}

	for i, arg := range n.PositionalArguments {
		if !arg.Equal(value.Ref(o.PositionalArguments[i])) {
			return false
		}
	}

	for i, arg := range n.NamedArguments {
		if !arg.Equal(value.Ref(o.NamedArguments[i])) {
			return false
		}
	}

	return true
}

// String returns a string representation of the GenericReceiverlessMethodCallNode.
func (n *GenericReceiverlessMethodCallNode) String() string {
	var buff strings.Builder

	buff.WriteString(n.MethodName)
	buff.WriteString("::[")

	for i, arg := range n.TypeArguments {
		if i > 0 {
			buff.WriteString(", ")
		}
		buff.WriteString(arg.String())
	}

	buff.WriteString("](")

	var hasMultilineArgs bool
	argCount := len(n.PositionalArguments) + len(n.NamedArguments)
	argStrings := make([]string, 0, argCount)

	for _, arg := range n.PositionalArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}
	for _, arg := range n.NamedArguments {
		argString := arg.String()
		argStrings = append(argStrings, argString)
		if strings.ContainsRune(argString, '\n') {
			hasMultilineArgs = true
		}
	}

	if hasMultilineArgs || argCount > 6 {
		buff.WriteRune('\n')
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, argStr, 1)
		}
		buff.WriteRune('\n')
	} else {
		for i, argStr := range argStrings {
			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(argStr)
		}
	}

	buff.WriteString(")")

	return buff.String()
}

func (*GenericReceiverlessMethodCallNode) IsStatic() bool {
	return false
}

func (*GenericReceiverlessMethodCallNode) Class() *value.Class {
	return value.GenericReceiverlessMethodCallNodeClass
}

func (*GenericReceiverlessMethodCallNode) DirectClass() *value.Class {
	return value.GenericReceiverlessMethodCallNodeClass
}

func (n *GenericReceiverlessMethodCallNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::GenericReceiverlessMethodCallNode{\n  &: %p", n)

	buff.WriteString(",\n  method_name: ")
	indent.IndentStringFromSecondLine(&buff, value.String(n.MethodName).Inspect(), 1)

	buff.WriteString(",\n  type_arguments: %%[\n")
	for i, element := range n.TypeArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  positional_arguments: %%[\n")
	for i, element := range n.PositionalArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString(",\n  named_arguments: %%[\n")
	for i, element := range n.NamedArguments {
		if i != 0 {
			buff.WriteString(",\n")
		}
		indent.IndentString(&buff, element.Inspect(), 2)
	}
	buff.WriteString("\n  ]")

	buff.WriteString("\n}")

	return buff.String()
}

func (n *GenericReceiverlessMethodCallNode) Error() string {
	return n.Inspect()
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
