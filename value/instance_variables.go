package value

// A slice containing instance variables of an Elk value
type InstanceVariables []Value

func (i InstanceVariables) Length() int {
	return len(i)
}

func (i *InstanceVariables) BoxOf(index int) *Box {
	i.ExpandUpToIndex(index)

	return (*Box)(&(*i)[index])
}

func (i InstanceVariables) Get(index int) Value {
	if i.Length() <= index {
		return Undefined
	}

	return i[index]
}

func (i *InstanceVariables) ExpandUpToIndex(index int) {
	missingSlots := index + 1 - i.Length()
	if missingSlots <= 0 {
		return
	}

	newSlice := make(InstanceVariables, i.Length()+missingSlots)
	copy(newSlice, *i)
	*i = newSlice
}

func (i *InstanceVariables) Set(index int, val Value) {
	i.ExpandUpToIndex(index)
	(*i)[index] = val
}
