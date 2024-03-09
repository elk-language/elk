package vm

import (
	"github.com/elk-language/elk/value"
	"github.com/google/go-cmp/cmp"
)

// ::Std::HashRecord
func init() {
	// Instance methods
	c := &value.HashRecordClass.MethodContainer
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashRecord)
			iterator := value.NewHashRecordIterator(self)
			return iterator, nil
		},
	)
	Def(
		c,
		"length",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashRecord)
			return value.SmallInt(self.Length()), nil
		},
	)
	Def(
		c,
		"[]",
		func(vm *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashRecord)
			key := args[1]
			result, err := HashRecordGet(vm, self, key)
			if err != nil {
				return nil, err
			}
			if result == nil {
				return value.Nil, nil
			}
			return result, nil
		},
		DefWithParameters("key"),
		DefWithSealed(),
	)
}

// ::Std::HashRecord::Iterator
func init() {
	// Instance methods
	c := &value.HashRecordIteratorClass.MethodContainer
	Def(
		c,
		"next",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			self := args[0].(*value.HashRecordIterator)
			return self.Next()
		},
	)
	Def(
		c,
		"iterator",
		func(_ *VM, args []value.Value) (value.Value, value.Value) {
			return args[0], nil
		},
	)

}

// Create a new hash record with the given entries.
func NewHashRecordWithElements(vm *VM, elements ...value.Pair) (*value.HashRecord, value.Value) {
	return NewHashRecordWithCapacityAndElements(vm, len(elements), elements...)
}

// Create a new hash record with the given entries.
func MustNewHashRecordWithElements(vm *VM, elements ...value.Pair) *value.HashRecord {
	hrec, err := NewHashRecordWithElements(vm, elements...)
	if err != nil {
		panic(err)
	}

	return hrec
}

func NewHashRecordWithCapacityAndElements(vm *VM, capacity int, elements ...value.Pair) (*value.HashRecord, value.Value) {
	h := value.NewHashRecord(capacity)
	for _, element := range elements {
		err := HashRecordSet(vm, h, element.Key, element.Value)
		if err != nil {
			return nil, err
		}
	}

	return h, nil
}

func MustNewHashRecordWithCapacityAndElements(vm *VM, capacity int, elements ...value.Pair) *value.HashRecord {
	hrec, err := NewHashRecordWithCapacityAndElements(vm, capacity, elements...)
	if err != nil {
		panic(err)
	}

	return hrec
}

// Delete the given key from the hashMap
func HashRecordDelete(vm *VM, hashRecord *value.HashRecord, key value.Value) (bool, value.Value) {
	return HashMapDelete(vm, (*value.HashMap)(hashRecord), key)
}

// Get the element under the given key.
func HashRecordGet(vm *VM, hashRecord *value.HashRecord, key value.Value) (value.Value, value.Value) {
	return HashMapGet(vm, (*value.HashMap)(hashRecord), key)
}

func HashRecordCopyTable(vm *VM, target *value.HashRecord, source []value.Pair) value.Value {
	return HashMapCopyTable(vm, (*value.HashMap)(target), source)
}

// Copy the pairs of one hash record to the other.
func HashRecordCopy(vm *VM, target *value.HashRecord, source *value.HashRecord) value.Value {
	return HashMapCopy(vm, (*value.HashMap)(target), (*value.HashMap)(source))
}

// Add additional n empty slots for new elements.
func HashRecordGrow(vm *VM, hashRecord *value.HashRecord, newSlots int) value.Value {
	return HashMapGrow(vm, (*value.HashMap)(hashRecord), newSlots)
}

// Resize the given hash record to the desired capacity.
func HashRecordSetCapacity(vm *VM, hashRecord *value.HashRecord, capacity int) value.Value {
	return HashMapSetCapacity(vm, (*value.HashMap)(hashRecord), capacity)
}

func HashRecordSetWithMaxLoad(vm *VM, hashRecord *value.HashRecord, key, val value.Value, maxLoad float64) value.Value {
	return HashMapSetWithMaxLoad(vm, (*value.HashMap)(hashRecord), key, val, maxLoad)
}

// Set a value under the given key.
func HashRecordSet(vm *VM, hashRecord *value.HashRecord, key, val value.Value) value.Value {
	return HashMapSet(vm, (*value.HashMap)(hashRecord), key, val)
}

func NewHashRecordComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *value.HashRecord) bool {
		if x == y {
			return true
		}
		if x.Length() != y.Length() {
			return false
		}

		v := New()
		for _, xPair := range x.Table {
			if xPair.Key == nil {
				continue
			}

			yVal, err := HashRecordGet(v, y, xPair.Key)
			if err != nil {
				return false
			}

			if !cmp.Equal(xPair.Value, yVal, *opts...) {
				return false
			}

		}

		return true
	})
}
