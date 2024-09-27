package types

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Used during typechecking as a placeholder for a future
// constant or type in using statements
type ConstantPlaceholder struct {
	AsName    value.Symbol
	FullName  string
	Container ConstantMap
	Location  *position.Location
	Sibling   *ConstantPlaceholder
	Checked   bool
	Replaced  bool
}

func IsConstantPlaceholder(typ Type) bool {
	_, ok := typ.(*ConstantPlaceholder)
	return ok
}

func NewConstantPlaceholder(asName value.Symbol, fullName string, container ConstantMap, location *position.Location) *ConstantPlaceholder {
	return &ConstantPlaceholder{
		AsName:    asName,
		FullName:  fullName,
		Container: container,
		Location:  location,
	}
}

func (p *ConstantPlaceholder) ToNonLiteral(env *GlobalEnvironment) Type {
	return p
}

func (*ConstantPlaceholder) IsLiteral() bool {
	return false
}

func (p *ConstantPlaceholder) inspect() string {
	return fmt.Sprintf("<ConstantPlaceholder: %s>", p.FullName)
}
