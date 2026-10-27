package vm

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/elk-language/elk/indent"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

const HashMapOfValueMaxLoad = 0.75

type HashMapOfValue struct {
	Table         []value.PairOfValue // underlying data container, it's `len` is always equal to it's `cap` and it also serves as the capacity of the HashMap
	OccupiedSlots int                 // number of slots taken by active pairs and those left by deleted pairs
	Elements      int                 // number of slots occupied by active pairs
	version       int                 // version of the map, each mutation increments this counter, used for guarding against concurrent mutation during iteration
}

var _ HashMap = &HashMapOfValue{}

func HashMapOfValueConstructor(class *value.Class) value.Value {
	return value.Ref(&HashMapOfValue{})
}

func NewHashMapOfValue(capacity int) *HashMapOfValue {
	return &HashMapOfValue{
		Table: make([]value.PairOfValue, capacity),
	}
}

func (h *HashMapOfValue) Keys() iter.Seq[value.Value] {
	return func(yield func(value.Value) bool) {
		for _, pair := range h.Table {
			if pair.Key().IsUndefined() {
				continue
			}

			if !yield(pair.Key()) {
				return
			}
		}
	}
}

func (h *HashMapOfValue) All() iter.Seq[value.PairOfValue] {
	return func(yield func(value.PairOfValue) bool) {
		for _, pair := range h.Table {
			if pair.Key().IsUndefined() {
				continue
			}

			if !yield(pair) {
				return
			}
		}
	}
}

func (h *HashMapOfValue) Iterate() iter.Seq2[value.Value, value.Value] {
	return func(yield func(value.Value, value.Value) bool) {
		originalVersion := h.version
		for _, pair := range h.Table {
			if pair.Key().IsUndefined() {
				continue
			}

			if originalVersion != h.version {
				yield(value.Undefined, value.NewMutationDuringIterationError(h.Class().Name()).ToValue())
				return
			}

			if !yield(pair.ToValue(), value.Undefined) {
				return
			}
		}
	}
}

func (h *HashMapOfValue) CloneHashMap(thread *Thread, capacity int) (HashMap, value.Value) {
	newMap := NewHashMapOfValue(capacity)
	err := HashMapOfValueCopy(thread, newMap, h)
	if !err.IsUndefined() {
		return nil, err
	}

	return newMap, value.Undefined
}

func (h *HashMapOfValue) NewHashMap(capacity int) HashMap {
	return NewHashMapOfValue(capacity)
}

func (h *HashMapOfValue) NewHashRecord(capacity int) HashRecord {
	return h.NewHashMap(capacity)
}

func (h *HashMapOfValue) CloneHashRecord(thread *Thread, capacity int) (HashRecord, value.Value) {
	return h.CloneHashMap(thread, capacity)
}

func (h *HashMapOfValue) IterNative() *HashMapOfValueIterator {
	return NewHashMapOfValueIterator(h)
}

func (h *HashMapOfValue) Iter() value.NativeIterator {
	return h.IterNative()
}

func (h *HashMapOfValue) IterMap() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *HashMapOfValue) IterRecord() value.NativeResettableIterator {
	return h.IterNative()
}

func (h *HashMapOfValue) GetPair(thread *Thread, key value.Value) (value.PairOfValue, value.Value) {
	return HashMapOfValueGetPair(thread, h, key)
}

func (h *HashMapOfValue) GetValUndefined(thread *Thread, key value.Value) (value.Value, value.Value) {
	return HashMapOfValueGet(thread, h, key)
}

func (h *HashMapOfValue) GetValNil(thread *Thread, key value.Value) (value.Value, value.Value) {
	val, err := h.GetValUndefined(thread, key)
	if err.IsNotUndefined() {
		return value.Undefined, err
	}
	if val.IsUndefined() {
		return value.Nil, value.Undefined
	}
	return val, value.Undefined
}

func (h *HashMapOfValue) SetVal(thread *Thread, key, val value.Value) value.Value {
	return HashMapOfValueSet(thread, h, key, val)
}

func (h *HashMapOfValue) RemoveKey(thread *Thread, key value.Value) (bool, value.Value) {
	return HashMapOfValueDelete(thread, h, key)
}

func (h *HashMapOfValue) Merge(thread *Thread, other value.Value) (*HashMapOfValue, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashMapOfValue:
		return HashMapOfValueMerge(thread, h, o)
	case HashRecord:
		return HashMapOfValueMergeInterface(thread, h, o)
	case Record:
		newMap := NewHashMapOfValue(h.Length() + o.Length())
		err := HashMapOfValueCopyInterface(thread, newMap, h)
		if err.IsNotUndefined() {
			return nil, err
		}

		err = HashMapOfValueCopyInterface(thread, newMap, o)
		if err.IsNotUndefined() {
			return nil, err
		}

		return newMap, value.Undefined
	default:
		result := NewHashMapOfValue(h.Length())

		err := HashMapOfValueCopyInterface(thread, result, h)
		if err.IsNotUndefined() {
			return nil, err
		}

		for element, err := range Iterate(thread, other) {
			if err.IsNotUndefined() {
				return nil, err
			}

			entry := element.AsReference().(value.Pair)
			result.SetVal(thread, entry.Key(), entry.Value())
		}

		return result, value.Undefined
	}
}

func (h *HashMapOfValue) MergeVal(thread *Thread, other value.Value) (value.Value, value.Value) {
	return value.RefErr(h.Merge(thread, other))
}

func (h *HashMapOfValue) Contains(thread *Thread, other value.Pair) (bool, value.Value) {
	return HashMapOfValueContains(thread, h, other)
}

func (h *HashMapOfValue) ContainsValue(thread *Thread, val value.Value) (bool, value.Value) {
	return HashMapOfValueContainsValue(thread, h, val)
}

func (h *HashMapOfValue) ContainsKey(thread *Thread, key value.Value) (bool, value.Value) {
	return HashMapOfValueContainsKey(thread, h, key)
}

func (h *HashMapOfValue) Equal(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashMapOfValue:
		return HashMapOfValueEqual(thread, h, o)
	case HashMap:
		return HashMapOfValueEqualInterface(thread, h, o)
	}

	return false, value.Undefined
}

func (h *HashMapOfValue) LaxEqual(thread *Thread, other value.Value) (bool, value.Value) {
	switch o := other.SafeAsReference().(type) {
	case *HashMapOfValue:
		return HashMapOfValueLaxEqual(thread, h, o)
	case HashMap:
		return HashMapOfValueLaxEqualInterface(thread, h, o)
	case HashRecord:
		return HashMapOfValueLaxEqualInterface(thread, h, o)
	default:
		if !value.IsA(other, value.MapMixin) {
			return false, value.Undefined
		}

		for element, err := range Iterate(thread, other) {
			if err.IsNotUndefined() {
				return false, err
			}

			pair := element.AsReference().(value.Pair)
			val, err := h.GetValUndefined(thread, pair.Key())
			if err.IsNotUndefined() {
				return false, err
			}
			if val.IsUndefined() {
				return false, value.Undefined
			}

			eq, err := Equal(thread, val, pair.Value())
			if err.IsNotUndefined() {
				return false, err
			}

			if value.Falsy(eq) {
				return false, value.Undefined
			}
		}

		return true, value.Undefined
	}
}

func (h *HashMapOfValue) Grow(thread *Thread, newSlots int) value.Value {
	return HashMapOfValueGrow(thread, h, newSlots)
}

func (*HashMapOfValue) Class() *value.Class {
	return value.HashMapClass
}

func (*HashMapOfValue) DirectClass() *value.Class {
	return value.HashMapClass
}

func (*HashMapOfValue) SingletonClass() *value.Class {
	return nil
}

func (h *HashMapOfValue) Clone() *HashMapOfValue {
	return &HashMapOfValue{
		Table:         slices.Clone(h.Table),
		OccupiedSlots: h.OccupiedSlots,
		Elements:      h.Elements,
	}
}

func (h *HashMapOfValue) Copy() value.Reference {
	return h.Clone()
}

func (h *HashMapOfValue) ToValue() value.Value {
	return value.Ref(h)
}

func (h *HashMapOfValue) Error() string {
	return h.Inspect()
}

const MAX_HASH_MAP_ELEMENTS_IN_INSPECT = 300

func (h *HashMapOfValue) Inspect() string {
	var hasMultilineElements bool
	keyStrings := make(
		[]string,
		0,
		min(MAX_HASH_MAP_ELEMENTS_IN_INSPECT, h.Length()),
	)
	valStrings := make(
		[]string,
		0,
		min(MAX_HASH_MAP_ELEMENTS_IN_INSPECT, h.Length()),
	)

	i := 0
	for _, entry := range h.Table {
		if entry.Key().IsUndefined() {
			continue
		}

		keyString := entry.Key().Inspect()
		keyStrings = append(keyStrings, keyString)

		valString := entry.Value().Inspect()
		valStrings = append(valStrings, valString)

		if strings.ContainsRune(keyString, '\n') ||
			strings.ContainsRune(valString, '\n') {
			hasMultilineElements = true
		}

		if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
			break
		}
		i++
	}

	var buff strings.Builder

	buff.WriteRune('{')
	if hasMultilineElements || h.Length() > 15 {
		buff.WriteRune('\n')
		for i := range len(keyStrings) {
			keyString := keyStrings[i]
			valString := valStrings[i]

			if i != 0 {
				buff.WriteString(",\n")
			}
			indent.IndentString(&buff, keyString, 1)
			buff.WriteString(" => ")
			indent.IndentStringFromSecondLine(&buff, valString, 1)

			if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(",\n  ...")
				break
			}
		}
		buff.WriteRune('\n')
	} else {
		for i := range len(keyStrings) {
			keyString := keyStrings[i]
			valString := valStrings[i]

			if i != 0 {
				buff.WriteString(", ")
			}
			buff.WriteString(keyString)
			buff.WriteString(" => ")
			buff.WriteString(valString)

			if i >= MAX_HASH_MAP_ELEMENTS_IN_INSPECT-1 {
				buff.WriteString(", ...")
				break
			}
		}
	}
	buff.WriteRune('}')

	return buff.String()
}

func (*HashMapOfValue) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *HashMapOfValue) Capacity() int {
	return len(h.Table)
}

func (h *HashMapOfValue) LeftCapacity() int {
	return h.Capacity() - h.Length()
}

func (h *HashMapOfValue) Length() int {
	return h.Elements
}

type HashMapOfValueIterator struct {
	HashMap  *HashMapOfValue
	Index    int
	snapshot []value.PairOfValue
	version  int
}

var _ value.NativeResettableIterator = &HashMapOfValueIterator{}

func NewHashMapOfValueIterator(hmap *HashMapOfValue) *HashMapOfValueIterator {
	iterator := &HashMapOfValueIterator{
		HashMap: hmap,
		version: hmap.version,
	}
	iterator.captureSnapshot()
	return iterator
}

func (h *HashMapOfValueIterator) captureSnapshot() {
	snapshot := make([]value.PairOfValue, 0, h.HashMap.Length())
	for pair := range h.HashMap.All() {
		snapshot = append(snapshot, pair)
	}
	h.snapshot = snapshot
}

func (*HashMapOfValueIterator) Class() *value.Class {
	return value.HashMapIteratorClass
}

func (*HashMapOfValueIterator) DirectClass() *value.Class {
	return value.HashMapIteratorClass
}

func (*HashMapOfValueIterator) SingletonClass() *value.Class {
	return nil
}

func (h *HashMapOfValueIterator) Copy() value.Reference {
	return &HashMapOfValueIterator{
		HashMap:  h.HashMap,
		Index:    h.Index,
		snapshot: h.snapshot,
		version:  h.version,
	}
}

func (i *HashMapOfValueIterator) ToValue() value.Value {
	return value.Ref(i)
}

func (h *HashMapOfValueIterator) Error() string {
	return h.Inspect()
}

func (h *HashMapOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::HashMap::Iterator{hash_map: %s}", h.HashMap.Inspect())
}

func (*HashMapOfValueIterator) InstanceVariables() *value.InstanceVariables {
	return nil
}

func (h *HashMapOfValueIterator) NextValue() (value.Value, value.Value) {
	if h.version != h.HashMap.version {
		return value.Undefined, value.NewMutationDuringIterationError(h.HashMap.Class().Name()).ToValue()
	}

	if h.Index >= len(h.snapshot) {
		return value.Undefined, value.S(symbol.L_stop_iteration).ToValue()
	}

	pair := &h.snapshot[h.Index]
	h.Index++
	return value.Ref(pair), value.Undefined
}

func (h *HashMapOfValueIterator) Reset() {
	h.Index = 0
	h.version = h.HashMap.version
}

// Create a new hashmap with the given entries.
func NewHashMapOfValueWithElements(vm *Thread, elements ...value.PairOfValue) (*HashMapOfValue, value.Value) {
	return NewHashMapOfValueWithCapacityAndElements(vm, len(elements), elements...)
}

// Create a new hashmap with the given entries.
func MustNewHashMapOfValueWithElements(vm *Thread, elements ...value.PairOfValue) *HashMapOfValue {
	hmap, err := NewHashMapOfValueWithElements(vm, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hmap
}

func NewHashMapOfValueWithCapacityAndElements(vm *Thread, capacity int, elements ...value.PairOfValue) (*HashMapOfValue, value.Value) {
	h := NewHashMapOfValue(capacity)
	for _, element := range elements {
		err := HashMapOfValueSet(vm, h, element.Key(), element.Value())
		if !err.IsUndefined() {
			return nil, err
		}
	}

	return h, value.Undefined
}

func MustNewHashMapOfValueWithCapacityAndElements(vm *Thread, capacity int, elements ...value.PairOfValue) *HashMapOfValue {
	hmap, err := NewHashMapOfValueWithCapacityAndElements(vm, capacity, elements...)
	if !err.IsUndefined() {
		panic(err)
	}

	return hmap
}

// Checks whether two hash maps are equal (lax)
func HashMapOfValueLaxEqual(vm *Thread, x *HashMapOfValue, y *HashMapOfValue) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xPair := range x.Table {
		if xPair.Key().IsUndefined() {
			continue
		}

		yVal, err := HashMapOfValueGet(vm, y, xPair.Key())
		if !err.IsUndefined() {
			return false, err
		}
		if yVal.IsUndefined() {
			return false, value.Undefined
		}
		eqVal, err := LaxEqual(vm, xPair.Value(), yVal)
		if !err.IsUndefined() {
			return false, err
		}
		equal := value.Truthy(eqVal)
		if !equal {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}

func HashMapOfValueLaxEqualInterface(vm *Thread, x *HashMapOfValue, y Record) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xPair := range x.Table {
		if xPair.Key().IsUndefined() {
			continue
		}

		yVal, err := y.GetValUndefined(vm, xPair.Key())
		if !err.IsUndefined() {
			return false, err
		}
		if yVal.IsUndefined() {
			return false, value.Undefined
		}
		eqVal, err := LaxEqual(vm, xPair.Value(), yVal)
		if !err.IsUndefined() {
			return false, err
		}
		equal := value.Truthy(eqVal)
		if !equal {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}

// Create a new map containing the pairs of both maps.
func HashMapOfValueMerge(vm *Thread, x, y *HashMapOfValue) (*HashMapOfValue, value.Value) {
	result := x.Clone()
	err := HashMapOfValueCopy(vm, result, y)
	if !err.IsUndefined() {
		return nil, err
	}
	return result, value.Undefined
}

func HashMapOfValueMergeInterface(vm *Thread, x *HashMapOfValue, y HashRecord) (*HashMapOfValue, value.Value) {
	result := x.Clone()
	err := HashMapOfValueCopyInterface(vm, result, y)
	if !err.IsUndefined() {
		return nil, err
	}
	return result, value.Undefined
}

// Checks whether two hash maps are equal
func HashMapOfValueEqual(vm *Thread, x *HashMapOfValue, y *HashMapOfValue) (bool, value.Value) {
	if x == y {
		return true, value.Undefined
	}
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xPair := range x.Table {
		if xPair.Key().IsUndefined() {
			continue
		}

		yVal, err := HashMapOfValueGet(vm, y, xPair.Key())
		if !err.IsUndefined() {
			return false, err
		}
		if yVal.IsUndefined() {
			return false, value.Undefined
		}
		eqVal, err := Equal(vm, xPair.Value(), yVal)
		if !err.IsUndefined() {
			return false, err
		}
		equal := value.Truthy(eqVal)
		if !equal {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}

func HashMapOfValueEqualInterface(vm *Thread, x *HashMapOfValue, y HashRecord) (bool, value.Value) {
	if x.Length() != y.Length() {
		return false, value.Undefined
	}

	for _, xPair := range x.Table {
		if xPair.Key().IsUndefined() {
			continue
		}

		yVal, err := y.GetValUndefined(vm, xPair.Key())
		if !err.IsUndefined() {
			return false, err
		}
		if yVal.IsUndefined() {
			return false, value.Undefined
		}
		eqVal, err := Equal(vm, xPair.Value(), yVal)
		if !err.IsUndefined() {
			return false, err
		}
		equal := value.Truthy(eqVal)
		if !equal {
			return false, value.Undefined
		}
	}

	return true, value.Undefined
}

// Delete the given key from the hashMap
func HashMapOfValueDelete(vm *Thread, hashMap *HashMapOfValue, key value.Value) (bool, value.Value) {
	if hashMap.Length() == 0 {
		return false, value.Undefined
	}

	index, err := HashMapOfValueIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return false, err
	}
	if index < 0 {
		return false, value.Undefined
	}
	if hashMap.Table[index].Key().IsUndefined() {
		return false, value.Undefined
	}

	hashMap.Table[index] = value.MakePairOfValue(
		value.Undefined,
		value.True.ToValue(),
	)
	// decrement only the count of active pairs
	// leave `OccupiedSlots` with the same value, cause
	// it represents the sum of slots that are active
	// and "zombie" slots left over by deleted pairs
	hashMap.Elements--
	hashMap.version++

	return true, value.Undefined
}

// Get the element under the given key.
// Returns (value, undefined) when the value has been found.
// Returns (undefined, undefined) when the key is not present.
// Returns (undefined, err) when there was an error.
func HashMapOfValueGet(vm *Thread, hashMap *HashMapOfValue, key value.Value) (value.Value, value.Value) {
	if hashMap.Length() == 0 {
		return value.Undefined, value.Undefined
	}

	index, err := HashMapOfValueIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return value.Undefined, err
	}
	if index == -1 {
		return value.Undefined, value.Undefined
	}

	return hashMap.Table[index].Value(), value.Undefined
}

// Get the pair with the given key.
// Returns (pair, undefined) when the pair has been found.
// Returns (undefined, undefined) when the key is not present.
// Returns (undefined, err) when there was an error.
func HashMapOfValueGetPair(vm *Thread, hashMap *HashMapOfValue, key value.Value) (p value.PairOfValue, err value.Value) {
	if hashMap.Length() == 0 {
		return p, err
	}

	index, err := HashMapOfValueIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return p, err
	}
	if index == -1 {
		return p, value.Undefined
	}

	return hashMap.Table[index], value.Undefined
}

// Check if the given pair is present in the map
func HashMapOfValueContains(vm *Thread, hashMap *HashMapOfValue, pair value.Pair) (bool, value.Value) {
	val, err := HashMapOfValueGet(vm, hashMap, pair.Key())
	if !err.IsUndefined() {
		return false, err
	}
	if val.IsUndefined() {
		return false, value.Undefined
	}

	equal, err := Equal(vm, val, pair.Value())
	if !err.IsUndefined() {
		return false, err
	}

	return value.Truthy(equal), value.Undefined
}

// Check if the given key is present in the map
func HashMapOfValueContainsKey(vm *Thread, hashMap *HashMapOfValue, key value.Value) (bool, value.Value) {
	if hashMap.Length() == 0 {
		return false, value.Undefined
	}

	index, err := HashMapOfValueIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return false, err
	}
	if index == -1 {
		return false, value.Undefined
	}

	pair := hashMap.Table[index]
	if pair.Key().IsUndefined() {
		return false, value.Undefined
	}

	return true, value.Undefined
}

// Check if the given value is present in the map
func HashMapOfValueContainsValue(vm *Thread, hashMap *HashMapOfValue, val value.Value) (bool, value.Value) {
	for _, pair := range hashMap.Table {
		// if the Key is undefined the entry is empty or deleted
		// so we skip it
		if pair.Key().IsUndefined() {
			continue
		}

		equal, err := Equal(vm, pair.Value(), val)
		if !err.IsUndefined() {
			return false, err
		}

		if value.Truthy(equal) {
			return true, value.Undefined
		}
	}

	return false, value.Undefined
}

func HashMapOfValueCopyTable(vm *Thread, target *HashMapOfValue, source []value.PairOfValue) value.Value {
	for _, entry := range source {
		if entry.Key().IsUndefined() {
			continue
		}

		err := HashMapOfValueSetWithMaxLoad(vm, target, entry.Key(), entry.Value(), 1)
		if !err.IsUndefined() {
			return err
		}
	}

	target.version++
	return value.Undefined
}

// Copy the pairs of one hashmap to the other.
func HashMapOfValueCopy(vm *Thread, target *HashMapOfValue, source *HashMapOfValue) value.Value {
	requiredCapacity := target.Length() + source.Length()
	if target.Capacity() < requiredCapacity {
		HashMapOfValueSetCapacity(vm, target, requiredCapacity)
	}

	target.version++
	for _, entry := range source.Table {
		if entry.Key().IsUndefined() {
			continue
		}

		i, err := HashMapOfValueIndex(vm, target, entry.Key())
		if !err.IsUndefined() {
			return err
		}
		if i == -1 {
			panic("no room in target hashmap during copy")
		}
		target.Table[i] = entry
		target.OccupiedSlots++
		target.Elements++
	}

	return value.Undefined
}

func HashMapOfValueCopyInterface(vm *Thread, target *HashMapOfValue, source Record) value.Value {
	requiredCapacity := target.Length() + source.Length()
	if target.Capacity() < requiredCapacity {
		HashMapOfValueSetCapacity(vm, target, requiredCapacity)
	}

	target.version++
	for entry := range source.All() {
		i, err := HashMapOfValueIndex(vm, target, entry.Key())
		if !err.IsUndefined() {
			return err
		}
		if i == -1 {
			panic("no room in target hashmap during copy")
		}
		target.Table[i] = entry
		target.OccupiedSlots++
		target.Elements++
	}

	return value.Undefined
}

// Add additional n empty slots for new elements.
func HashMapOfValueGrow(vm *Thread, hashMap *HashMapOfValue, newSlots int) value.Value {
	return HashMapOfValueSetCapacity(vm, hashMap, hashMap.Capacity()+newSlots)
}

// Resize the given hashmap to the desired capacity.
func HashMapOfValueSetCapacity(vm *Thread, hashMap *HashMapOfValue, capacity int) value.Value {
	if hashMap.Capacity() == capacity {
		return value.Undefined
	}

	oldTable := hashMap.Table
	newTable := make([]value.PairOfValue, capacity)
	tmpHashMap := &HashMapOfValue{
		Table: newTable,
	}
	hashMap.version++

	for _, entry := range oldTable {
		if entry.Key().IsUndefined() {
			continue
		}

		i, err := HashMapOfValueIndex(vm, tmpHashMap, entry.Key())
		if !err.IsUndefined() {
			return err
		}
		if i == -1 {
			panic("no room in target hashmap during resizing")
		}
		newTable[i] = entry
		tmpHashMap.OccupiedSlots++
		tmpHashMap.Elements++
	}

	hashMap.OccupiedSlots = tmpHashMap.OccupiedSlots
	hashMap.Elements = tmpHashMap.Elements
	hashMap.Table = newTable
	return value.Undefined
}

func HashMapOfValueSetWithMaxLoad(vm *Thread, hashMap *HashMapOfValue, key, val value.Value, maxLoad float64) value.Value {
	if hashMap.Capacity() == 0 {
		HashMapOfValueSetCapacity(vm, hashMap, 5)
	} else if float64(hashMap.OccupiedSlots) >= float64(hashMap.Capacity())*maxLoad {
		HashMapOfValueSetCapacity(vm, hashMap, hashMap.OccupiedSlots*2)
	}

	index, err := HashMapOfValueIndex(vm, hashMap, key)
	if !err.IsUndefined() {
		return err
	}
	if index == -1 {
		panic(fmt.Sprintf("no room in target hashmap when trying to add a new key: %s", hashMap.Inspect()))
	}
	entry := hashMap.Table[index]
	if entry.Key().IsUndefined() {
		hashMap.Elements++
		if entry.Value().IsUndefined() {
			// increment OccupiedSlots only when the slot was completely empty
			// leave the count the same when a zombie slot is overwritten
			hashMap.OccupiedSlots++
		}
	}

	hashMap.Table[index] = value.MakePairOfValue(
		key,
		val,
	)
	hashMap.version++

	return value.Undefined
}

// Set a value under the given key.
func HashMapOfValueSet(vm *Thread, hashMap *HashMapOfValue, key, val value.Value) value.Value {
	return HashMapOfValueSetWithMaxLoad(vm, hashMap, key, val, HashMapOfValueMaxLoad)
}

// Get the index that the key should be inserted into.
// Returns (nil, err) when an error has been encountered.
// Returns (-1, nil) when there's no room for new values.
func HashMapOfValueIndex(vm *Thread, hashMap *HashMapOfValue, key value.Value) (int, value.Value) {
	hash, err := Hash(vm, key)
	if !err.IsUndefined() {
		return 0, err
	}
	deletedIndex := -1

	capacity := hashMap.Capacity()
	index := int(hash % value.UInt64(capacity))
	startIndex := index

	for {
		entry := hashMap.Table[index]
		if entry.Key().IsUndefined() {
			// empty or deleted entry

			if entry.Value().IsUndefined() {
				// empty bucket
				if deletedIndex != -1 {
					return deletedIndex, value.Undefined
				}
				return index, value.Undefined
			} else if deletedIndex == -1 {
				// deleted entry
				deletedIndex = index
			}
		} else {
			// present entry
			equal, err := Equal(vm, entry.Key(), key)
			if !err.IsUndefined() {
				return 0, err
			}
			if value.Truthy(equal) {
				return index, value.Undefined
			}
		}

		if index == capacity-1 {
			index = 0
		} else {
			index++
		}

		// when we reach the start index
		// all slots are checked
		if index == startIndex {
			return -1, value.Undefined
		}
	}
}
