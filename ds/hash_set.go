package ds

import (
	"fmt"
	"slices"
)

type hashSetEntryType uint8

const (
	hashSetEntryEmptyType hashSetEntryType = iota
	hashSetEntryOccupiedType
	hashSetEntryDeletedType
)

type hashSetEntry[V Hashable] struct {
	val V
	typ hashSetEntryType
}

// A set with custom hashing and equality methods
type HashSet[V Hashable] struct {
	table         []hashSetEntry[V] // underlying data container, it's `len` is always equal to it's `cap` and it also serves as the capacity of the HashSet
	occupiedSlots int               // number of slots taken by active values and those left by deleted values
	elements      int               // number of slots occupied by active values
	version       int               // version of the set, each mutation increments this counter, used for guarding against concurrent mutation during iteration
}

func MakeHashSet[V Hashable](capacity int) HashSet[V] {
	return HashSet[V]{
		table: make([]hashSetEntry[V], capacity),
	}
}

func MakeHashSetWithValues[V Hashable](values []V) HashSet[V] {
	s := MakeHashSet[V](len(values))
	for _, val := range values {
		s.Add(val)
	}
	return s
}

func NewHashSet[V Hashable](capacity int) *HashSet[V] {
	return &HashSet[V]{
		table: make([]hashSetEntry[V], capacity),
	}
}

func NewHashSetWithValues[V Hashable](values []V) *HashSet[V] {
	s := NewHashSet[V](len(values))
	for _, val := range values {
		s.Add(val)
	}
	return s
}

func (s *HashSet[V]) Copy() *HashSet[V] {
	return &HashSet[V]{
		table:         slices.Clone(s.table),
		occupiedSlots: s.occupiedSlots,
		elements:      s.elements,
		version:       s.version,
	}
}

func (s HashSet[V]) CopyVal() HashSet[V] {
	return HashSet[V]{
		table:         slices.Clone(s.table),
		occupiedSlots: s.occupiedSlots,
		elements:      s.elements,
		version:       s.version,
	}
}

func (s *HashSet[V]) Len() int {
	return s.elements
}

func (s *HashSet[V]) Capacity() int {
	return len(s.table)
}

func (s *HashSet[V]) LeftCapacity() int {
	return s.Capacity() - s.Len()
}

const HashSetMaxLoad = 0.75

func (s *HashSet[V]) Add(v V) bool {
	return s.AppendWithMaxLoad(v, HashSetMaxLoad)
}

func (s *HashSet[V]) Index(v V) int {
	hash := v.HashUint64()
	deletedIndex := -1

	capacity := s.Capacity()
	index := int(hash % uint64(capacity))
	startIndex := index

	for {
		entry := s.table[index]
		switch entry.typ {
		case hashSetEntryEmptyType:
			// empty bucket
			if deletedIndex != -1 {
				return deletedIndex
			}
			return index
		case hashSetEntryDeletedType:
			if deletedIndex == -1 {
				// deleted entry
				deletedIndex = index
			}
		case hashSetEntryOccupiedType:
			// present entry
			if entry.val.EqualAny(v) {
				return index
			}
		default:
			panic(fmt.Sprintf("invalid hashset entry type: %d", entry.typ))
		}

		if index == capacity-1 {
			index = 0
		} else {
			index++
		}

		// when we reach the start index
		// all slots are checked
		if index == startIndex {
			return -1
		}
	}
}

func (s *HashSet[V]) SetCapacity(capacity int) {
	if s.Capacity() == capacity {
		return
	}

	oldTable := s.table
	newTable := make([]hashSetEntry[V], capacity)
	tmpHashSet := &HashSet[V]{
		table: newTable,
	}

tableLoop:
	for _, entry := range oldTable {
		switch entry.typ {
		case hashSetEntryDeletedType, hashSetEntryEmptyType:
			continue tableLoop
		}

		i := tmpHashSet.Index(entry.val)
		if i == -1 {
			panic("no room in target hashset during resizing")
		}
		newTable[i] = entry
		tmpHashSet.occupiedSlots++
		tmpHashSet.elements++
	}

	s.occupiedSlots = tmpHashSet.occupiedSlots
	s.elements = tmpHashSet.elements
	s.table = newTable
	s.version++
}

func (s *HashSet[V]) AppendWithMaxLoad(v V, maxLoad float64) bool {
	if s.Capacity() == 0 {
		s.SetCapacity(5)
	} else if float64(s.occupiedSlots) >= float64(s.Capacity())*maxLoad {
		s.SetCapacity(s.occupiedSlots * 2)
	}

	index := s.Index(v)
	if index == -1 {
		panic(fmt.Sprintf("no room in target hashset when trying to add a new value: %#v", s))
	}
	entry := s.table[index]

	var newValue bool

	switch entry.typ {
	case hashSetEntryEmptyType:
		// the slot is empty
		s.occupiedSlots++
		s.elements++
		newValue = true
	case hashSetEntryDeletedType:
		// this is a zombie slot, just overwrite it's content
		s.elements++
		newValue = true
	}

	s.table[index] = hashSetEntry[V]{
		val: v,
		typ: hashSetEntryOccupiedType,
	}
	s.version++

	return newValue
}

func (s *HashSet[V]) Remove(v V) bool {
	if s.Len() == 0 {
		return false
	}

	index := s.Index(v)
	if index < 0 {
		return false
	}
	existingVal := s.table[index]
	switch existingVal.typ {
	case hashSetEntryDeletedType, hashSetEntryEmptyType:
		return false
	}

	// mark as deleted
	s.table[index] = hashSetEntry[V]{typ: hashSetEntryDeletedType}
	s.elements--
	s.version++

	return true
}

func (s *HashSet[V]) Contains(v V) bool {
	if s.Len() == 0 {
		return false
	}

	index := s.Index(v)
	if index == -1 {
		return false
	}

	valInSlot := s.table[index]
	switch valInSlot.typ {
	case hashSetEntryDeletedType, hashSetEntryEmptyType:
		return false
	case hashSetEntryOccupiedType:
		return true
	default:
		panic(fmt.Sprintf("invalid hashset entry type: %d", valInSlot.typ))
	}
}

func (s *HashSet[V]) Get(v V) (result V, ok bool) {
	if s.Len() == 0 {
		return result, false
	}

	index := s.Index(v)
	if index == -1 {
		return result, false
	}

	valInSlot := s.table[index]
	switch valInSlot.typ {
	case hashSetEntryDeletedType, hashSetEntryEmptyType:
		return result, false
	case hashSetEntryOccupiedType:
		return valInSlot.val, true
	default:
		panic(fmt.Sprintf("invalid hashset entry type: %d", valInSlot.typ))
	}
}
