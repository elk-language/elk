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
	case DefaultValueParameterKind:
		return true
	default:
		return false
	}
}

type Method struct {
	Name               string
	Params             []*Parameter
	OptionalParamCount int
	PostParamCount     int
	HasNamedRestParam  bool
	ReturnType         Type
	ThrowType          Type
	DefinedUnder       *ConstantMap
}

func NewMethod(name string, params []*Parameter, returnType Type, throwType Type, definedUnder *ConstantMap) *Method {
	var optParamCount int
	var hasNamedRestParam bool
	postParamCount := -1
	for _, param := range params {
		switch param.Kind {
		case NormalParameterKind:
			if postParamCount != -1 {
				postParamCount++
			}
		case DefaultValueParameterKind:
			optParamCount++
			if postParamCount != -1 {
				postParamCount++
			}
		case PositionalRestParameterKind:
			postParamCount++
		case NamedRestParameterKind:
			hasNamedRestParam = true
		}
	}

	return &Method{
		Name:               name,
		Params:             params,
		ReturnType:         returnType,
		ThrowType:          throwType,
		DefinedUnder:       definedUnder,
		HasNamedRestParam:  hasNamedRestParam,
		OptionalParamCount: optParamCount,
		PostParamCount:     postParamCount,
	}
}

func (m *Method) RequiredParamCount() int {
	requiredParamCount := len(m.Params) - m.OptionalParamCount
	if m.HasNamedRestParam {
		requiredParamCount--
	}
	if m.HasPositionalRestParam() {
		requiredParamCount--
	}
	return requiredParamCount
}

func (m *Method) ExpectedParamCountString() string {
	requiredParamCount := m.RequiredParamCount()
	if m.HasNamedRestParam || m.HasPositionalRestParam() {
		return fmt.Sprintf("%d...", requiredParamCount)
	}

	if requiredParamCount == len(m.Params) {
		return fmt.Sprintf("%d", requiredParamCount)
	}

	return fmt.Sprintf("%d...%d", requiredParamCount, len(m.Params))
}

func (m *Method) NamedRestParam() *Parameter {
	if m.HasNamedRestParam {
		return m.Params[len(m.Params)-1]
	}
	return nil
}

func (m *Method) HasPositionalRestParam() bool {
	return m.PostParamCount != -1
}

func (m *Method) PositionalRestParamIndex() int {
	if m.PostParamCount == -1 {
		return -1
	}

	index := len(m.Params) - 1 - m.PostParamCount
	if m.HasNamedRestParam {
		index--
	}
	return index
}

func (m *Method) PositionalRestParam() *Parameter {
	index := m.PositionalRestParamIndex()
	if index == -1 {
		return nil
	}
	return m.Params[index]
}

func (m *Method) IsDefinedUnder(constantMap *ConstantMap) bool {
	return m.DefinedUnder == constantMap
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
