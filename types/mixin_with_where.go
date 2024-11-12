package types

import (
	"strings"

	"github.com/elk-language/elk/lexer"
)

// Represents a mixin parent with a where clause
// `extend where Bar < Baz`
type MixinWithWhere struct {
	*MixinProxy
	Namespace Namespace
	Where     []*TypeParameter
}

func IsMixinWithWhere(typ Type) bool {
	_, ok := typ.(*MixinWithWhere)
	return ok
}

func NewMixinWithWhere(mixin *MixinProxy, namespace Namespace, where []*TypeParameter) *MixinWithWhere {
	return &MixinWithWhere{
		MixinProxy: mixin,
		Namespace:  namespace,
		Where:      where,
	}
}

func (m *MixinWithWhere) ToNonLiteral(env *GlobalEnvironment) Type {
	return m
}

func (m *MixinWithWhere) inspect() string {
	return m.Namespace.inspect()
}

func (m *MixinWithWhere) InspectExtend() string {
	buffer := new(strings.Builder)
	buffer.WriteString("extend where ")
	firstIteration := true
	for _, whereElement := range m.Where {
		if !firstIteration {
			buffer.WriteString(", ")
		} else {
			firstIteration = false
		}

		buffer.WriteString(whereElement.InspectSignature())
	}

	return buffer.String()
}

func (m *MixinWithWhere) InspectExtendWithColor() string {
	return lexer.Colorize(m.InspectExtend())
}

func (m *MixinWithWhere) Copy() *MixinWithWhere {
	return &MixinWithWhere{
		MixinProxy: m.MixinProxy,
		Namespace:  m.Namespace,
		Where:      m.Where,
	}
}

func (m *MixinWithWhere) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *MixinWithWhere {
	newMixin := m.Copy()
	newMixin.MixinProxy = DeepCopyEnv(m.MixinProxy, oldEnv, newEnv).(*MixinProxy)
	newMixin.Namespace = DeepCopyEnv(m.Namespace, oldEnv, newEnv).(Namespace)
	newMixin.Where = TypeParametersDeepCopyEnv(m.Where, oldEnv, newEnv)

	return newMixin
}
