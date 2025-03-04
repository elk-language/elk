package ast

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// All nodes that should be valid in parameter declaration lists
// of methods or functions should implement this interface.
type ParameterNode interface {
	Node
	parameterNode()
	IsOptional() bool
}

func (*InvalidNode) parameterNode()            {}
func (*FormalParameterNode) parameterNode()    {}
func (*MethodParameterNode) parameterNode()    {}
func (*SignatureParameterNode) parameterNode() {}
func (*AttributeParameterNode) parameterNode() {}

// checks whether the given parameter is a positional rest parameter.
func IsPositionalRestParam(p ParameterNode) bool {
	switch param := p.(type) {
	case *MethodParameterNode:
		return param.Kind == PositionalRestParameterKind
	case *FormalParameterNode:
		return param.Kind == PositionalRestParameterKind
	case *SignatureParameterNode:
		return param.Kind == PositionalRestParameterKind
	default:
		return false
	}
}

// checks whether the given parameter is a named rest parameter.
func IsNamedRestParam(p ParameterNode) bool {
	switch param := p.(type) {
	case *MethodParameterNode:
		return param.Kind == NamedRestParameterKind
	case *FormalParameterNode:
		return param.Kind == NamedRestParameterKind
	case *SignatureParameterNode:
		return param.Kind == NamedRestParameterKind
	default:
		return false
	}
}

// Indicates whether the parameter is a rest param
type ParameterKind uint8

const (
	NormalParameterKind ParameterKind = iota
	PositionalRestParameterKind
	NamedRestParameterKind
)

// Represents a formal parameter in function or struct declarations eg. `foo: String = 'bar'`
type FormalParameterNode struct {
	TypedNodeBase
	Name        string         // name of the variable
	TypeNode    TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
	Kind        ParameterKind
}

func (*FormalParameterNode) IsStatic() bool {
	return false
}

func (f *FormalParameterNode) IsOptional() bool {
	return f.Initialiser != nil
}

// Create a new formal parameter node eg. `foo: String = 'bar'`
func NewFormalParameterNode(span *position.Span, name string, typ TypeNode, init ExpressionNode, kind ParameterKind) *FormalParameterNode {
	return &FormalParameterNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Name:          name,
		TypeNode:      typ,
		Initialiser:   init,
		Kind:          kind,
	}
}

func (*FormalParameterNode) Class() *value.Class {
	return value.FormalParameterNodeClass
}

func (*FormalParameterNode) DirectClass() *value.Class {
	return value.FormalParameterNodeClass
}

func (n *FormalParameterNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::FormalParameterNode{\n  &: %p", n)

	buff.WriteString(",\n  name: ")
	buff.WriteString(n.Name)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString(",\n  initialiser: ")
	indentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)

	buff.WriteString(",\n  kind: ")
	buff.WriteString(value.UInt8(n.Kind).Inspect())

	buff.WriteString("\n}")

	return buff.String()
}

func (n *FormalParameterNode) Error() string {
	return n.Inspect()
}

// Represents a formal parameter in method declarations eg. `foo: String = 'bar'`
type MethodParameterNode struct {
	TypedNodeBase
	Name                string         // name of the variable
	TypeNode            TypeNode       // type of the variable
	Initialiser         ExpressionNode // value assigned to the variable
	SetInstanceVariable bool           // whether an instance variable with this name gets automatically assigned
	Kind                ParameterKind
}

func (*MethodParameterNode) IsStatic() bool {
	return false
}

func (f *MethodParameterNode) IsOptional() bool {
	return f.Initialiser != nil
}

func (*MethodParameterNode) Class() *value.Class {
	return value.MethodParameterNodeClass
}

func (*MethodParameterNode) DirectClass() *value.Class {
	return value.MethodParameterNodeClass
}

func (n *MethodParameterNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::MethodParameterNode{\n  &: %p", n)

	buff.WriteString(",\n  name: ")
	buff.WriteString(n.Name)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString(",\n  initialiser: ")
	indentStringFromSecondLine(&buff, n.Initialiser.Inspect(), 1)

	buff.WriteString(",\n  kind: ")
	buff.WriteString(value.UInt8(n.Kind).Inspect())

	buff.WriteString("\n}")

	return buff.String()
}

func (n *MethodParameterNode) Error() string {
	return n.Inspect()
}

// Create a new formal parameter node eg. `foo: String = 'bar'`
func NewMethodParameterNode(span *position.Span, name string, setIvar bool, typ TypeNode, init ExpressionNode, kind ParameterKind) *MethodParameterNode {
	return &MethodParameterNode{
		TypedNodeBase:       TypedNodeBase{span: span},
		SetInstanceVariable: setIvar,
		Name:                name,
		TypeNode:            typ,
		Initialiser:         init,
		Kind:                kind,
	}
}

// Represents a signature parameter in method and function signatures eg. `foo?: String`
type SignatureParameterNode struct {
	TypedNodeBase
	Name     string   // name of the variable
	TypeNode TypeNode // type of the variable
	Optional bool     // whether this parameter is optional
	Kind     ParameterKind
}

func (*SignatureParameterNode) IsStatic() bool {
	return false
}

func (f *SignatureParameterNode) IsOptional() bool {
	return f.Optional
}

func (*SignatureParameterNode) Class() *value.Class {
	return value.SignatureParameterNodeClass
}

func (*SignatureParameterNode) DirectClass() *value.Class {
	return value.SignatureParameterNodeClass
}

func (n *SignatureParameterNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::SignatureParameterNode{\n  &: %p", n)

	buff.WriteString(",\n  name: ")
	buff.WriteString(n.Name)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	fmt.Fprintf(&buff, ",\n  optional: %t", n.Optional)

	buff.WriteString(",\n  kind: ")
	buff.WriteString(value.UInt8(n.Kind).Inspect())

	buff.WriteString("\n}")

	return buff.String()
}

func (n *SignatureParameterNode) Error() string {
	return n.Inspect()
}

// Create a new signature parameter node eg. `foo?: String`
func NewSignatureParameterNode(span *position.Span, name string, typ TypeNode, opt bool, kind ParameterKind) *SignatureParameterNode {
	return &SignatureParameterNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Name:          name,
		TypeNode:      typ,
		Optional:      opt,
		Kind:          kind,
	}
}

// Represents an attribute declaration in getters, setters and accessors eg. `foo: String`
type AttributeParameterNode struct {
	TypedNodeBase
	Name        string         // name of the variable
	TypeNode    TypeNode       // type of the variable
	Initialiser ExpressionNode // value assigned to the variable
}

func (*AttributeParameterNode) IsStatic() bool {
	return false
}

func (a *AttributeParameterNode) IsOptional() bool {
	return a.Initialiser != nil
}

func (*AttributeParameterNode) Class() *value.Class {
	return value.AttributeParameterNodeClass
}

func (*AttributeParameterNode) DirectClass() *value.Class {
	return value.AttributeParameterNodeClass
}

func (n *AttributeParameterNode) Inspect() string {
	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Elk::AST::AttributeParameterNode{\n  &: %p", n)

	buff.WriteString(",\n  name: ")
	buff.WriteString(n.Name)

	buff.WriteString(",\n  type_node: ")
	indentStringFromSecondLine(&buff, n.TypeNode.Inspect(), 1)

	buff.WriteString(",\n  initialiser: ")
	buff.WriteString(n.Initialiser.Inspect())

	buff.WriteString("\n}")

	return buff.String()
}

func (n *AttributeParameterNode) Error() string {
	return n.Inspect()
}

// Create a new attribute declaration in getters, setters and accessors eg. `foo: String`
func NewAttributeParameterNode(span *position.Span, name string, typ TypeNode, init ExpressionNode) *AttributeParameterNode {
	return &AttributeParameterNode{
		TypedNodeBase: TypedNodeBase{span: span},
		Name:          name,
		TypeNode:      typ,
		Initialiser:   init,
	}
}
