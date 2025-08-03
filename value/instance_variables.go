package value

// A slice containing instance variables of an Elk value
type InstanceVariables []Value

func (i InstanceVariables) Length() int {
	return len(i)
}

func (i InstanceVariables) Get(index int) Value {
	if i.Length() <= index {
		return Undefined
	}

	return i[index]
}

func (i *InstanceVariables) Set(index int, val Value) {
	missingSlots := index + 1 - i.Length()
	if missingSlots > 0 {
		newSlice := make(InstanceVariables, i.Length()+missingSlots)
		copy(newSlice, *i)
		*i = newSlice
	}

	(*i)[index] = val
}
