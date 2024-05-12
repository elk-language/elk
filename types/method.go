package types

import (
	"strings"

	"github.com/elk-language/elk/value"
)

type ParameterKind uint8

const (
	NormalParameterKind ParameterKind = iota
	DefaultValueParameterKind
	PositionalRestParameterKind
	NamedRestParameterKind
)

type Parameter struct {
	Type             Type
	Kind             ParameterKind
	InstanceVariable bool
}

func (p *Parameter) IsPositionalRest() bool {
	return p.Kind == PositionalRestParameterKind
}

func (p *Parameter) IsNamedRest() bool {
	return p.Kind == NamedRestParameterKind
}

func (p *Parameter) HasDefaultValue() bool {
	return p.Kind == DefaultValueParameterKind
}

func (p *Parameter) IsOptional() bool {
	switch p.Kind {
	case DefaultValueParameterKind, PositionalRestParameterKind, NamedRestParameterKind:
		return true
	default:
		return false
	}
}

type ParameterMap map[value.Symbol]*Parameter

type Method struct {
	Name       string
	Params     ParameterMap
	ReturnType Type
	ThrowType  Type
}

func NewMethod(name string, params ParameterMap, returnType Type, throwType Type) *Method {
	return &Method{
		Name:       name,
		Params:     params,
		ReturnType: returnType,
		ThrowType:  throwType,
	}
}

func (m *Method) Inspect() string {
	buffer := new(strings.Builder)
	buffer.WriteString("sig ")
	buffer.WriteString(m.Name)
	buffer.WriteRune('(')
	firstIteration := true
	for name, param := range m.Params {
		if !firstIteration {
			buffer.WriteString(", ")
		}
		if param.IsPositionalRest() {
			buffer.WriteRune('*')
		} else if param.IsNamedRest() {
			buffer.WriteString("**")
		}
		buffer.WriteString(name.String())
		if param.HasDefaultValue() {
			buffer.WriteRune('?')
		}
		buffer.WriteString(": ")
		buffer.WriteString(param.Type.Inspect())
	}
	buffer.WriteRune(')')

	return buffer.String()
}

func (m *Method) IsSupertypeOf(other Type) bool {
	otherMethod, ok := other.(*Method)
	if !ok {
		return false
	}

	return m == otherMethod
}
