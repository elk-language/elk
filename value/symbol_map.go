package value

import (
	"fmt"
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

func (s SymbolMap) Inspect() string {
	if len(s) == 0 {
		return "{}"
	}

	var buff strings.Builder
	buff.WriteString("{ ")
	firstIteration := true

	for symbol, val := range s {
		if !firstIteration {
			buff.WriteString(", ")
		}
		buff.WriteString(fmt.Sprintf("%s: %s", symbol.InspectContent(), val.Inspect()))
		firstIteration = false
	}

	buff.WriteString(" }")
	return buff.String()
}

// Get a value stored using the given key.
func (s SymbolMap) Get(key Symbol) (Value, bool) {
	val, ok := s[key]
	return val, ok
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
func (s SymbolMap) GetString(key string) (Value, bool) {
	symbol, ok := SymbolTable.Get(key)
	if !ok {
		return nil, false
	}

	val, ok := s[symbol]
	return val, ok
}

// Check if the given key exists.
func (s SymbolMap) HasString(key string) bool {
	_, ok := s.GetString(key)
	return ok
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
	SymbolMapClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("SymbolMap", SymbolMapClass)
}
