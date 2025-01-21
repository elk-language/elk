package value

import (
	"fmt"
	"maps"
	"strings"
)

var SymbolMapClass *Class // ::Std::SymbolMap

// Simple map with symbol keys.
type SymbolMap map[Symbol]Value

func (SymbolMap) Class() *Class {
	return SymbolMapClass
}

func (SymbolMap) DirectClass() *Class {
	return SymbolMapClass
}

func (SymbolMap) SingletonClass() *Class {
	return nil
}

func (s SymbolMap) InstanceVariables() SymbolMap {
	return nil
}

func (s SymbolMap) Copy() Reference {
	newMap := make(SymbolMap, len(s))
	maps.Copy(newMap, s)
	return newMap
}

func (s SymbolMap) Error() string {
	return s.Inspect()
}

func (s SymbolMap) Inspect() string {
	if len(s) == 0 {
		return "{}"
	}

	var buff strings.Builder
	buff.WriteRune('{')
	firstIteration := true

	for symbol, val := range s {
		if !firstIteration {
			buff.WriteString(", ")
		}
		fmt.Fprintf(&buff, "%s: %s", symbol.InspectContent(), val.Inspect())
		firstIteration = false
	}

	buff.WriteRune('}')
	return buff.String()
}

// Get a value stored using the given key.
func (s SymbolMap) Get(key Symbol) Value {
	return s[key]
}

func (s SymbolMap) Delete(key Symbol) {
	delete(s, key)
}

// Check if the given key exists.
func (s SymbolMap) Has(key Symbol) bool {
	_, ok := s[key]
	return ok
}

// Set the passed value under the given key.
func (s SymbolMap) Set(key Symbol, val Value) {
	s[key] = val
}

// Get a value stored using the given key.
func (s SymbolMap) GetString(key string) (val Value) {
	symbol, ok := SymbolTable.Get(key)
	if !ok {
		return val
	}

	return s[symbol]
}

// Check if the given key exists.
func (s SymbolMap) HasString(key string) bool {
	return !s.GetString(key).IsUndefined()
}

func (s SymbolMap) DeleteString(key string) {
	symbol, ok := SymbolTable.Get(key)
	if !ok {
		return
	}

	delete(s, symbol)
}

// Set the passed value under the given key.
func (s SymbolMap) SetString(key string, val Value) {
	s[SymbolTable.Add(key)] = val
}

func initSymbolMap() {
	SymbolMapClass = NewClass()
	StdModule.AddConstantString("SymbolMap", Ref(SymbolMapClass))
}
