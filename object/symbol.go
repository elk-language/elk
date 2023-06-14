package object

// Numerical ID of a particular symbol.
type SymbolId int

// Represent a symbol eg. `:foo`
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
