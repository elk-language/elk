package value

import "unsafe"

// Elk Reference Value
type Reference interface {
	ValueInterface
	Copy() Reference // Creates a shallow copy of the reference. If the value is immutable, no copying should be done, the same value should be returned.
}

// Convert a Reference to a Value
func Ref(ref Reference) Value {
	i := *(*iface)(unsafe.Pointer(&ref))

	return Value{
		data: i.tab,
		ptr:  i.ptr,
		flag: REFERENCE_FLAG,
	}
}
