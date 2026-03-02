package value

import "github.com/cespare/xxhash/v2"

var BoolClass *Class // ::Std::Bool

type Bool bool

var _ ValueInterface = Bool(true)

const True = Bool(true)
const False = Bool(false)

func BoolVal(val bool) Value {
	return Bool(val).ToValue()
}

// Converts an Elk Value to an Elk Bool.
func ToBool(val Value) Bool {
	if val.IsReference() {
		return true
	}

	switch val.ValueFlag() {
	case BOOL_FLAG:
		if val.AsBool() {
			return true
		} else {
			return false
		}
	case NIL_FLAG:
		return false
	default:
		return true
	}
}

// Converts an Elk Value to an Elk Bool
// and negates it.
func ToNotBool(val Value) Bool {
	return !ToBool(val)
}

func (b Bool) ToValue() Value {
	var data uintptr
	if b {
		data = 1
	} else {
		data = 0
	}

	return Value{
		flag: BOOL_FLAG,
		data: data,
	}
}

func (b Bool) Class() *Class {
	if b {
		return TrueClass
	} else {
		return FalseClass
	}
}

func (b Bool) DirectClass() *Class {
	return b.Class()
}

func (b Bool) Hash() UInt64 {
	d := xxhash.New()
	var val byte
	if b {
		val = 1
	} else {
		val = 0
	}

	d.Write([]byte{val})
	return UInt64(d.Sum64())
}

func (Bool) SingletonClass() *Class {
	return nil
}

func (b Bool) Inspect() string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func (b Bool) Error() string {
	return b.Inspect()
}

func (Bool) InstanceVariables() *InstanceVariables {
	return nil
}

func initBool() {
	BoolClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	StdModule.AddConstantString("Bool", Ref(BoolClass))
	RegisterNativeClass("Std::Bool", "value.BoolClass")
}
