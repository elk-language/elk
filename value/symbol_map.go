package value

import (
	"fmt"
	"strings"
)

var SymbolMapClass *Class // ::Std::SymbolMap

// Map with symbol keys.
type SymbolMap struct {
	SimpleSymbolMap
	frozen bool
}

func (s *SymbolMap) IsFrozen() bool {
	return s.frozen
}

func (s *SymbolMap) SetFrozen() {
	s.frozen = true
}

func (s *SymbolMap) Class() *Class {
	return SymbolMapClass
}

func (s *SymbolMap) InstanceVariables() SimpleSymbolMap {
	return nil
}

// Simple map with symbol keys.
type SimpleSymbolMap map[Symbol]Value

func (s SimpleSymbolMap) Inspect() string {
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
func (s SimpleSymbolMap) Get(key Symbol) (Value, bool) {
	val, ok := s[key]
	return val, ok
}

func (s SimpleSymbolMap) Delete(key Symbol) {
	delete(s, key)
}

// Check if the given key exists.
func (s SimpleSymbolMap) Has(key Symbol) bool {
	_, ok := s[key]
	return ok
}

// Set the passed value under the given key.
func (s SimpleSymbolMap) Set(key Symbol, val Value) {
	s[key] = val
}

// Get a value stored using the given key.
func (s SimpleSymbolMap) GetString(key string) (Value, bool) {
	symbol, ok := SymbolTable.Get(key)
	if !ok {
		return nil, false
	}

	val, ok := s[symbol]
	return val, ok
}

// Check if the given key exists.
func (s SimpleSymbolMap) HasString(key string) bool {
	_, ok := s.GetString(key)
	return ok
}

func (s SimpleSymbolMap) DeleteString(key string) {
	symbol, ok := SymbolTable.Get(key)
	if !ok {
		return
	}

	delete(s, symbol)
}

// Set the passed value under the given key.
func (s SimpleSymbolMap) SetString(key string, val Value) {
	s[SymbolTable.Add(key)] = val
}

func initSymbolMap() {
	SymbolMapClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("SymbolMap", SymbolMapClass)
}
