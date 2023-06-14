package object

import "sync"

// Data structure that holds Elk Symbols.
type symbolTableStruct struct {
	table  map[string]*Symbol
	lastId SymbolId
	mutex  sync.RWMutex
}

type symbolTableOption func(*symbolTableStruct)

func symbolTableWithTable(table map[string]*Symbol) symbolTableOption {
	return func(s *symbolTableStruct) {
		s.table = table
	}
}

func symbolTableWithLastId(lastId SymbolId) symbolTableOption {
	return func(s *symbolTableStruct) {
		s.lastId = lastId
	}
}

func newSymbolTable(opts ...symbolTableOption) *symbolTableStruct {
	s := &symbolTableStruct{table: make(map[string]*Symbol)}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Check if a symbol with the given name exists.
func (s *symbolTableStruct) Exists(str string) bool {
	return s.Get(str) != nil
}

// Get the Symbol with the specified name.
// This function is thread-safe.
func (s *symbolTableStruct) Get(name string) *Symbol {
	s.mutex.RLock()
	val := s.table[name]
	s.mutex.RUnlock()

	return val
}

// Add a new symbol with the specified name.
// This function is idempotent, if the Symbol already exists
// nothing happens and a pointer to it gets returned.
// This function is thread-safe.
func (s *symbolTableStruct) Add(name string) *Symbol {
	s.mutex.Lock()
	val := s.table[name]
	if val != nil {
		s.mutex.Unlock()
		return val
	}

	s.lastId++
	val = newSymbol(name, s.lastId)
	s.table[name] = val
	s.mutex.Unlock()
	return val
}

// The global Symbol Table of the Elk interpreter.
var SymbolTable *symbolTableStruct = newSymbolTable()
