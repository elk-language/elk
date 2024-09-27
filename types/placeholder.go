package types

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

// Used during typechecking as a placeholder for a future
// constant or type in using statements
type Placeholder struct {
	AsName    value.Symbol
	FullName  string
	Container ConstantMap
	Location  *position.Location
	Sibling   *Placeholder
	Checked   bool
	Replaced  bool
}

func IsPlaceholder(typ Type) bool {
	_, ok := typ.(*Placeholder)
	return ok
}

func NewPlaceholder(asName value.Symbol, fullName string, container ConstantMap, location *position.Location) *Placeholder {
	return &Placeholder{
		AsName:    asName,
		FullName:  fullName,
		Container: container,
		Location:  location,
	}
}

func (p *Placeholder) ToNonLiteral(env *GlobalEnvironment) Type {
	return p
}

func (*Placeholder) IsLiteral() bool {
	return false
}

func (p *Placeholder) inspect() string {
	return fmt.Sprintf("<placeholder: %s>", p.FullName)
}
