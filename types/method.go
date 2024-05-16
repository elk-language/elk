package types

import (
	"fmt"
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
	Name             value.Symbol
	Type             Type
	Kind             ParameterKind
	InstanceVariable bool
}

func NewParameter(name value.Symbol, typ Type, kind ParameterKind, instanceVariable bool) *Parameter {
	return &Parameter{
		Name:             name,
		Type:             typ,
		Kind:             kind,
		InstanceVariable: instanceVariable,
	}
}

func (p *Parameter) NameWithKind() string {
	name := p.Name.String()
	if p.InstanceVariable {
		name = "@" + name
	}
	switch p.Kind {
	case NormalParameterKind:
		return name
	case DefaultValueParameterKind:
		return fmt.Sprintf("%s = x", name)
	case PositionalRestParameterKind:
		return fmt.Sprintf("*%s", name)
	case NamedRestParameterKind:
		return fmt.Sprintf("**%s", name)
	default:
		panic("invalid parameter kind")
	}
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

type Method struct {
	Name       string
	Params     []*Parameter
	ReturnType Type
	ThrowType  Type
}

func NewMethod(name string, params []*Parameter, returnType Type, throwType Type) *Method {
	return &Method{
		Name:       name,
		Params:     params,
		ReturnType: returnType,
		ThrowType:  throwType,
	}
}

func (m *Method) inspect() string {
	buffer := new(strings.Builder)
	buffer.WriteString("sig ")
	buffer.WriteString(m.Name)
	buffer.WriteRune('(')
	firstIteration := true
	for _, param := range m.Params {
		if !firstIteration {
			buffer.WriteString(", ")
		}
		if param.IsPositionalRest() {
			buffer.WriteRune('*')
		} else if param.IsNamedRest() {
			buffer.WriteString("**")
		}
		buffer.WriteString(param.Name.String())
		if param.HasDefaultValue() {
			buffer.WriteRune('?')
		}
		buffer.WriteString(": ")
		buffer.WriteString(Inspect(param.Type))
	}
	buffer.WriteRune(')')

	return buffer.String()
}

func (m *Method) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}
