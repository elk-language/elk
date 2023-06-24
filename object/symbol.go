package object

import (
	"fmt"
)

// Numerical ID of a particular symbol.
type SymbolId int

var SymbolClass *Class // ::Std::Symbol

// Represents a symbol eg. `:foo`
type Symbol struct {
	Name string
	Id   SymbolId
}

// Create a new Symbol.
func newSymbol(name string, id SymbolId) *Symbol {
	return &Symbol{
		Name: name,
		Id:   id,
	}
}

func (s *Symbol) Class() *Class {
	return StringClass
}

func (s *Symbol) IsFrozen() bool {
	return true
}

func (s *Symbol) SetFrozen() {}

func (s *Symbol) Inspect() string {
	return fmt.Sprintf(":%q", s)
}

func (s *Symbol) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initSymbol() {
	SymbolClass = NewClass(
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("Symbol", SymbolClass)
}
