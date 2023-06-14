package object

// Map with symbol keys.
type SymbolMap map[SymbolId]Object

// Get a value stored using the given key.
func (s SymbolMap) Get(key *Symbol) Object {
	return s[key.Id]
}

// Set the passed value under the given key.
func (s SymbolMap) Set(key *Symbol, val Object) {
	s[key.Id] = val
}

// Get a value stored using the given key.
func (s SymbolMap) GetString(key string) Object {
	symbol := SymbolTable.Get(key)
	if symbol == nil {
		return nil
	}

	return s[symbol.Id]
}

// Set the passed value under the given key.
func (s SymbolMap) SetString(key string, val Object) {
	s[SymbolTable.Add(key).Id] = val
}

// Get a value stored using the given ID.
func (s SymbolMap) GetId(id SymbolId) Object {
	return s[id]
}

// Set the passed value under the given ID.
func (s SymbolMap) SetId(id SymbolId, val Object) {
	s[id] = val
}
